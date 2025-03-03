package java

import (
	"fmt"
	"proto-qiu/protoc"
	"strings"
)

// 生成枚举类
func (jp *JavaProtoc) generateEnum(enum *protoc.Enum) string {
	className := toCamelCase(enum.Name, true)

	var builder strings.Builder
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
	builder.WriteString("        this.value = value;\n    }\n\n")
	builder.WriteString("    public int getNumber() {\n        return value;\n    }\n")
	builder.WriteString(" 	public static " + className + " forNumber(int value) {\n")
	builder.WriteString("        for (" + className + " e : values()) {\n")
	builder.WriteString("            if (e.value == value) {\n")
	builder.WriteString("                return e;\n")
	builder.WriteString("            }\n")
	builder.WriteString("        }\n")
	builder.WriteString("        return null;\n")
	builder.WriteString("    }\n")
	builder.WriteString("}\n")

	return builder.String()
}
