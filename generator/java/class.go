// Package java is java language generator
package java

import (
	"fmt"
	"proto-qiu/protoc"
	"strings"
)

func (jp *JavaProtoc) generateOuterClass(innerStr strings.Builder) string {
	var fileStr strings.Builder
	fileStr.WriteString(fmt.Sprintf("package %s;\n\n", jp.PackageName))
	var b bool
	for _, message := range jp.Protoc.Messages {
		if message.Name == jp.ProtoName {
			b = true
		}
	}

	fileStr.WriteString("@javax.annotation.Generated(\"by proto-qiu\")\n")
	if b {
		fileStr.WriteString("public final class " + toCamelCase(jp.ProtoName, true) + "Outer {\n")
	} else {
		fileStr.WriteString("public final class " + toCamelCase(jp.ProtoName, true) + " {\n")
	}
	fileStr.WriteString("private static final com.protoc.qiu.GeneratedMessage message = new com.protoc.qiu.GeneratedMessage();\n")
	fileStr.WriteString(innerStr.String())
	fileStr.WriteString("}\n")

	return fileStr.String()
}

func (jp *JavaProtoc) generateMessageClass(msg *protoc.Message, inner bool) string {
	className := toCamelCase(msg.Name, true)

	var builder strings.Builder

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
		if field.MapInfo != nil {
			builder.WriteString(fmt.Sprintf("        this.%s = new java.util.HashMap<>();\n", toCamelCase(field.Name, false)))
		} else if field.Repeated {
			builder.WriteString(fmt.Sprintf("        this.%s = new java.util.ArrayList<>();\n", toCamelCase(field.Name, false)))
		}

	}
	builder.WriteString("    }\n")

	// 生成 Getter/Setter
	for _, field := range msg.Fields {
		javaType := toJavaType(field)
		// Getter
		builder.WriteString(fmt.Sprintf("\n    public %s get%s() {\n", javaType, toCamelCase(field.Name, true)))
		builder.WriteString(fmt.Sprintf("        return this.%s;\n    }\n", toCamelCase(field.Name, false)))
		// Setter
		builder.WriteString(fmt.Sprintf("\n    public void set%s(%s %s) {\n",
			toCamelCase(field.Name, true), javaType, toCamelCase(field.Name, false)))
		builder.WriteString(fmt.Sprintf("        this.%s = %s;\n    }\n", toCamelCase(field.Name, false), toCamelCase(field.Name, false)))
	}

	// 处理 oneof 字段
	for _, oneof := range msg.OneOfs {
		builder.WriteString("\n    // OneOf: " + oneof.Name + "\n")
		builder.WriteString("    private Object " + toCamelCase(oneof.Name, false) + ";\n")
		builder.WriteString("    private int " + toCamelCase(oneof.Name, false) + "Case = 0;\n")

		// 生成枚举类型表示当前激活的字段
		builder.WriteString(fmt.Sprintf("    public enum %sCase {\n", toCamelCase(oneof.Name, true)))
		for _, f := range oneof.Fields {
			builder.WriteString(fmt.Sprintf("        %s(%d),\n", strings.ToUpper(f.Name), f.FieldNumber))
		}
		builder.WriteString("        NOT_SET(0);\n")

		builder.WriteString("        private final int value;\n")
		builder.WriteString(fmt.Sprintf("        private %sCase(int value) {\n", toCamelCase(oneof.Name, true)))
		builder.WriteString("            this.value = value;\n")
		builder.WriteString("        }\n")
		builder.WriteString("        public int getValue() { return value; }\n")
		builder.WriteString("    }\n\n")

		// 生成 getter/setter
		for _, f := range oneof.Fields {
			javaType := toJavaType(f)
			fieldName := toCamelCase(f.Name, true)

			// Getter
			builder.WriteString(fmt.Sprintf("    public %s get%s() {\n", javaType, fieldName))
			builder.WriteString(fmt.Sprintf("        if (%sCase == %d) {\n", toCamelCase(oneof.Name, false), f.FieldNumber))
			builder.WriteString(fmt.Sprintf("            return (%s) %s;\n", javaType, toCamelCase(oneof.Name, false)))
			builder.WriteString("        }\n")
			builder.WriteString("        return " + getDefaultValue(f) + ";\n")
			builder.WriteString("    }\n\n")

			// Setter
			builder.WriteString(fmt.Sprintf("    public void set%s(%s value) {\n", fieldName, javaType))
			builder.WriteString(fmt.Sprintf("        %s = value;\n", toCamelCase(oneof.Name, false)))
			builder.WriteString(fmt.Sprintf("        %sCase = %d;\n", toCamelCase(oneof.Name, false), f.FieldNumber))
			builder.WriteString("    }\n\n")
		}

		// Case getter
		builder.WriteString(fmt.Sprintf("    public %sCase get%sCase() {\n", toCamelCase(oneof.Name, true), toCamelCase(oneof.Name, true)))
		builder.WriteString(fmt.Sprintf("        switch (%sCase) {\n", toCamelCase(oneof.Name, false)))
		for _, f := range oneof.Fields {
			builder.WriteString(fmt.Sprintf("            case %d: return %sCase.%s;\n",
				f.FieldNumber, toCamelCase(oneof.Name, true), strings.ToUpper(f.Name)))
		}
		builder.WriteString("            default: return " + toCamelCase(oneof.Name, true) + "Case.NOT_SET;\n")
		builder.WriteString("        }\n")
		builder.WriteString("    }\n\n")

		// Clear method
		builder.WriteString(fmt.Sprintf("    public void clear%s() {\n", toCamelCase(oneof.Name, true)))
		builder.WriteString(fmt.Sprintf("        %s = null;\n", toCamelCase(oneof.Name, false)))
		builder.WriteString(fmt.Sprintf("        %sCase = 0;\n", toCamelCase(oneof.Name, false)))
		builder.WriteString("    }\n\n")
	}

	// inner class
	for _, innerMsg := range msg.InnerMessages {
		builder.WriteString(jp.generateMessageClass(innerMsg, true))
	}
	// 添加序列化和反序列化方法
	builder.WriteString(jp.generateToByteArray(msg))
	builder.WriteString(jp.generateParseFrom(msg))

	builder.WriteString("}\n")

	return builder.String()
}
