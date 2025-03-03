package java

import (
	"fmt"
	"proto-qiu/constant"
	"proto-qiu/protoc"
	"strings"
)

// Proto 类型到 Java 类型映射
func boxed(s string) string {
	switch s {
	case constant.JavaInt:
		return constant.JavaInteger
	case constant.JavaFloat:
		return constant.JavaBoxedFloat
	case constant.JavaBoolean:
		return constant.JavaBoxedBoolean
	case constant.JavaDouble:
		return constant.JavaBoxedDouble
	case constant.JavaLong:
		return constant.JavaBoxedLong
	default:
		return s
	}
}

func toJavaType(field *protoc.Field) string {
	var javaTypeName string
	if field.MapInfo != nil {
		// 处理 Map 类型
		keyType := toJavaType(&protoc.Field{TypeName: field.MapInfo.KeyType})
		valueType := toJavaType(&protoc.Field{TypeName: field.MapInfo.ValueType})
		return fmt.Sprintf("java.util.Map<%s, %s>", boxed(keyType), boxed(valueType))
	}
	switch field.TypeName {
	case constant.TypeString:
		javaTypeName = constant.JavaString
	case "int32", "uint32", "sint32", "fixed32", "sfixed32":
		javaTypeName = "int"
	case "int64", "uint64", "sint64", "fixed64", "sfixed64":
		javaTypeName = "long"
	case "bool":
		javaTypeName = "boolean"
	case "double":
		javaTypeName = "double"
	case "float":
		javaTypeName = "float"
	case "bytes":
		javaTypeName = constant.JavaByteArray
	default:
		javaTypeName = toCamelCase(field.TypeName, true) // 自定义类型
	}
	if field.Repeated {
		// to boxed type
		return fmt.Sprintf("java.util.List<%s>", boxed(javaTypeName))
	}
	return javaTypeName
}

func getDefaultValue(field *protoc.Field) string {
	if field.MapInfo != nil {
		return "new java.lang.HashMap<>()"
	}
	if field.Repeated {
		return "new java.lang.ArrayList<>()"
	}
	switch field.TypeName {
	case "int32", "uint32", "sint32", "fixed32", "sfixed32":
		return "0"
	case "int64", "uint64", "sint64", "fixed64", "sfixed64":
		return "0L"
	case "bool":
		return "false"
	case "double":
		return "0.0"
	case "float":
		return "0.0f"
	default:
		return "null"
	}
}
func getDefaultValueByStr(str string) string {
	switch str {
	case "int32", "uint32", "sint32", "fixed32", "sfixed32":
		return "0"
	case "int64", "uint64", "sint64", "fixed64", "sfixed64":
		return "0L"
	case "bool":
		return "false"
	case "double":
		return "0.0"
	case "float":
		return "0.0f"
	default:
		return "null"
	}
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

// element in list or map
func getElementType(field *protoc.Field) string {
	if field.MapInfo != nil {
		keyType := toJavaType(&protoc.Field{TypeName: field.MapInfo.KeyType})
		valueType := toJavaType(&protoc.Field{TypeName: field.MapInfo.ValueType})
		return fmt.Sprintf("java.util.Map.Entry<%s, %s>", boxed(keyType), boxed(valueType))
	}
	return boxed(toJavaType(&protoc.Field{TypeName: field.TypeName}))
}
