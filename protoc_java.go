package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var _ generator = (*JavaProtoc)(nil)

type JavaProtoc struct {
	Protoc
	JavaOutput    string
	ProtoFilePath string
}

func newJavaProtoc(javaOutput, protoFilePath string) (*JavaProtoc, error) {
	// read file
	content, err := os.ReadFile(protoFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read proto file: %v", err)
	}

	// 解析.proto文件
	parser := NewParser(strings.NewReader(string(content)))
	protoc, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("proto parse error: %v", err)
	}

	return &JavaProtoc{
		Protoc:        *protoc,
		JavaOutput:    javaOutput,
		ProtoFilePath: protoFilePath,
	}, nil

}

// generate .java file
func (jp *JavaProtoc) generate() error {
	// 创建包对应的目录
	var fileStr strings.Builder
	packagePath := filepath.Join(jp.JavaOutput, strings.Replace(jp.PackageName, ".", "/", -1))
	if err := os.MkdirAll(packagePath, 0755); err != nil {
		return fmt.Errorf("failed to create package dir: %v", err)
	}

	// 为每个 Message 生成 Java 类
	for _, msg := range jp.Messages {
		fileStr.WriteString(jp.generateMessageClass(msg, false))
	}

	// 为每个 Enum 生成 Java 枚举
	for _, enum := range jp.Enums {
		fileStr.WriteString(jp.generateEnum(packagePath, enum))
	}
	// get proto file Name
	protoFileName := strings.Split(filepath.Base(jp.ProtoFilePath), ".")[0]
	err := os.WriteFile(jp.JavaOutput+protoFileName+".java", []byte(jp.JavaOutput), 0644)

	if err != nil {
		return fmt.Errorf("failed to write proto file: %v", err)
	}

	return nil
}

func (jp *JavaProtoc) generateMessageClass(msg *Message, inner bool) string {
	className := toCamelCase(msg.Name, true)

	var builder strings.Builder

	// 生成包声明
	if !inner {
		builder.WriteString(fmt.Sprintf("package %s;\n\n", jp.PackageName))
	}

	// 类声明
	if inner {
		builder.WriteString(fmt.Sprintf("public static class %s {\n", className))
	} else {
		builder.WriteString(fmt.Sprintf("public final class %s {\n", className))
	}

	// 生成字段声明
	for _, field := range msg.Fields {
		javaType := toJavaType(field)
		builder.WriteString(fmt.Sprintf("    private %s %s;\n", javaType, toCamelCase(field.Name, false)))
	}

	// 生成构造方法
	builder.WriteString("\n    public " + className + "() {\n")
	for _, field := range msg.Fields {
		if field.Repeated {
			builder.WriteString(fmt.Sprintf("        this.%s = new java.util.ArrayList<>();\n", toCamelCase(field.Name, false)))
		}
		if field.MapInfo != nil {
			builder.WriteString(fmt.Sprintf("        this.%s = new java.util.HashMap<>();\n", toCamelCase(field.Name, false)))
		}
	}
	builder.WriteString("    }\n")

	// 生成 Getter/Setter
	for _, field := range msg.Fields {
		javaType := toJavaType(field)
		// Getter
		builder.WriteString(fmt.Sprintf("\n    public %s get%s() {\n", javaType, toCamelCase(field.Name, true)))
		builder.WriteString(fmt.Sprintf("        return this.%s;\n    }\n", field.Name))
		// Setter
		builder.WriteString(fmt.Sprintf("\n    public void set%s(%s %s) {\n",
			toCamelCase(field.Name, true), javaType, field.Name))
		builder.WriteString(fmt.Sprintf("        this.%s = %s;\n    }\n", field.Name, field.Name))
	}

	// 处理 oneof 字段
	for _, oneof := range msg.OneOfs {
		builder.WriteString("\n    // OneOf: " + oneof.Name + "\n")
		builder.WriteString("    private Object " + oneof.Name + ";\n")
		// 生成枚举类型表示当前激活的字段
		builder.WriteString(fmt.Sprintf("    public enum %sCase {\n", toCamelCase(oneof.Name, true)))
		for _, f := range oneof.Fields {
			builder.WriteString(fmt.Sprintf("        %s,\n", strings.ToUpper(f.Name)))
		}
		builder.WriteString("        NOT_SET\n    }\n")
	}

	// inner class
	for _, innerMsg := range msg.InnerMessages {
		builder.WriteString(jp.generateMessageClass(innerMsg, true))
	}

	builder.WriteString("}\n")

	return builder.String()
}

// Proto 类型到 Java 类型映射
func boxed(s string) string {
	switch s {
	case "int":
		return "java.lang.Integer"
	case "float":
		return "java.lang.Float"
	case "bool":
		return "java.lang.Boolean"
	case "double":
		return "java.lang.Double"
	case "long":
		return "java.lang.Long"
	default:
		return s
	}
}

func toJavaType(field *Field) string {
	var javaTypeName string
	if field.MapInfo != nil {
		// 处理 Map 类型
		keyType := toJavaType(&Field{TypeName: field.MapInfo.KeyType})
		valueType := toJavaType(&Field{TypeName: field.MapInfo.ValueType})
		return fmt.Sprintf("java.util.Map<%s, %s>", boxed(keyType), boxed(valueType))
	}
	switch field.TypeName {
	case "string":
		javaTypeName = "java.lang.String"
	case "int32", "uint32", "sint32":
		javaTypeName = "int"
	case "int64", "uint64", "sint64":
		javaTypeName = "long"
	case "bool":
		javaTypeName = "boolean"
	case "double", "fixed64", "sfixed64":
		javaTypeName = "double"
	case "float", "fixed32", "sfixed32":
		javaTypeName = "float"
	case "bytes":
		javaTypeName = "com.google.protobuf.ByteString"
	default:
		javaTypeName = toCamelCase(field.TypeName, true) // 自定义类型
	}
	if field.Repeated {
		// to boxed type
		return fmt.Sprintf("java.util.List<%s>", boxed(javaTypeName))
	}
	return javaTypeName
}

// 生成枚举类
func (jp *JavaProtoc) generateEnum(packagePath string, enum *Enum) string {
	className := toCamelCase(enum.Name, true)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("package %s;\n\n", jp.PackageName))
	builder.WriteString(fmt.Sprintf("public enum %s {\n", className))

	for i, value := range enum.Values {
		if i > 0 {
			builder.WriteString(",\n")
		}
		builder.WriteString(fmt.Sprintf("    %s(%d)", strings.ToUpper(value.Name), value.Value))
	}

	builder.WriteString(";\n\n")
	builder.WriteString("    private final int value;\n\n")
	builder.WriteString(fmt.Sprintf("    %s(int value) {\n", className))
	builder.WriteString("        this.Value = Value;\n    }\n\n")
	builder.WriteString("    public int getNumber() {\n        return value;\n    }\n")
	builder.WriteString("}\n")

	return builder.String()
}

// 工具函数：下划线命名转驼峰
func toCamelCase(s string, firstUpper bool) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if firstUpper && i == 0 {
			parts[i] = strings.ToUpper(parts[i][0:1]) + parts[i][1:]
		}
		if i > 0 {
			parts[i] = strings.ToUpper(parts[i][0:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
