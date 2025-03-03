package java

import (
	"fmt"
	"proto-qiu/protoc"
	"strconv"
	"strings"
)

func (jp *JavaProtoc) generateToByteArray(msg *protoc.Message) string {
	var builder strings.Builder
	writeMethodHeader(&builder)
	writeFields(&builder, msg.Fields)
	writeOneofs(&builder, msg.OneOfs)
	writeMethodFooter(&builder)
	return builder.String()
}

func writeMethodHeader(builder *strings.Builder) {
	builder.WriteString("\n    public byte[] toByteArray() {\n")
	builder.WriteString("        java.io.ByteArrayOutputStream stream = new java.io.ByteArrayOutputStream();\n")
	builder.WriteString("        try {\n")
}

func writeFields(builder *strings.Builder, fields []*protoc.Field) {
	for _, field := range fields {
		fieldName := toCamelCase(field.Name, false)
		if field.MapInfo != nil {
			writeMapField(builder, fieldName, field)
		} else if field.Repeated {
			writeRepeatedField(builder, fieldName, field)
		} else {
			writeSimpleField(builder, fieldName, field)
		}
	}
}

func writeSimpleField(builder *strings.Builder, fieldName string, field *protoc.Field) {
	builder.WriteString(fmt.Sprintf("            if (%s != %s) {\n",
		fieldName, getDefaultValue(field)))
	builder.WriteString(generateWriteField(fieldName, field))
	builder.WriteString("            }\n")
}

func writeOneofs(builder *strings.Builder, oneofs []*protoc.OneOf) {
	for _, oneof := range oneofs {
		writeOneofField(builder, oneof)
	}
}

func writeMethodFooter(builder *strings.Builder) {
	builder.WriteString("            return stream.toByteArray();\n")
	builder.WriteString("        } catch (Exception e) {\n")
	builder.WriteString("            throw new RuntimeException(\"Failed to serialize message\", e);\n")
	builder.WriteString("        } finally {\n")
	builder.WriteString("            try {\n")
	builder.WriteString("                stream.close();\n")
	builder.WriteString("            } catch (Exception e) {\n")
	builder.WriteString("                // Ignore close exception\n")
	builder.WriteString("            }\n")
	builder.WriteString("        }\n")
	builder.WriteString("    }\n")
}

// 基本类型序列化
func writeBasicField(builder *strings.Builder, fieldName, typeName string, fieldNumber int) {
	builder.WriteString(fmt.Sprintf("        message.write%s(stream, %d, %s);\n",
		toCamelCase(typeName, true), fieldNumber, fieldName))
}

func writeMessageField(builder *strings.Builder, fieldName string, fieldNumber int) {
	builder.WriteString(fmt.Sprintf("        if (%s != null) {\n", fieldName))
	builder.WriteString(fmt.Sprintf("            byte[] bytes = %s.toByteArray();\n", fieldName))
	builder.WriteString(fmt.Sprintf("            message.writeTag(stream, %d, message.WIRETYPE_LENGTH_DELIMITED);\n", fieldNumber))
	builder.WriteString("            message.writeBytes(stream, bytes);\n")
	builder.WriteString("        }\n")
}

// 处理 repeated 字段序列化
func writeRepeatedField(builder *strings.Builder, fieldName string, field *protoc.Field) {
	builder.WriteString(fmt.Sprintf("        if (%s != null) {\n", fieldName))
	builder.WriteString(fmt.Sprintf("            for (%s item : %s) {\n",
		getElementType(field), fieldName))
	builder.WriteString(generateWriteField("item", field))
	builder.WriteString("            }\n")
	builder.WriteString("        }\n")
}

// 处理 map 字段序列化
func writeMapField(builder *strings.Builder, fieldName string, field *protoc.Field) {
	builder.WriteString(fmt.Sprintf("        if (%s != null) {\n", fieldName))
	builder.WriteString(fmt.Sprintf("            for (java.util.Map.Entry<%s, %s> entry : %s.entrySet()) {\n",
		boxed(toJavaType(&protoc.Field{TypeName: field.MapInfo.KeyType})),
		boxed(toJavaType(&protoc.Field{TypeName: field.MapInfo.ValueType})),
		fieldName))

	// Map 键值序列化
	writeMapKeyValue(builder, field)
	builder.WriteString("            }\n")
	builder.WriteString("        }\n")
}

// 处理 map 键值序列化
func writeMapKeyValue(builder *strings.Builder, field *protoc.Field) {
	builder.WriteString("                message.writeTag(stream, " +
		strconv.Itoa(field.FieldNumber) + ", message.WIRETYPE_LENGTH_DELIMITED);\n")
	builder.WriteString("                java.io.ByteArrayOutputStream mapStream = new java.io.ByteArrayOutputStream();\n")

	keyField := &protoc.Field{TypeName: field.MapInfo.KeyType, FieldNumber: 1}
	valueField := &protoc.Field{TypeName: field.MapInfo.ValueType, FieldNumber: 2}

	builder.WriteString(generateWriteField("entry.getKey()", keyField))
	builder.WriteString(generateWriteField("entry.getValue()", valueField))
	builder.WriteString("                byte[] mapBytes = mapStream.toByteArray();\n")
	builder.WriteString("                message.writeBytes(stream, mapBytes);\n")
}

// 处理 oneof 字段序列化
func writeOneofField(builder *strings.Builder, oneof *protoc.OneOf) {
	builder.WriteString(fmt.Sprintf("            switch (%sCase) {\n",
		toCamelCase(oneof.Name, false)))

	for _, f := range oneof.Fields {
		builder.WriteString(fmt.Sprintf("                case %d:\n", f.FieldNumber))
		builder.WriteString(generateWriteField(
			fmt.Sprintf("(%s)%s", toJavaType(f), toCamelCase(oneof.Name, false)), f))
		builder.WriteString("                    break;\n")
	}

	builder.WriteString("            }\n")
}

func generateWriteField(varName string, field *protoc.Field) string {
	var builder strings.Builder
	switch field.TypeName {
	case "int32", "uint32":
		builder.WriteString(fmt.Sprintf("        message.writeInt32(stream, %d, %s);\n", field.FieldNumber, varName))
	case "int64", "uint64":
		builder.WriteString(fmt.Sprintf("        message.writeInt64(stream, %d, %s);\n", field.FieldNumber, varName))
	case "sint32":
		builder.WriteString(fmt.Sprintf("        message.writeSint32(stream, %d, %s);\n", field.FieldNumber, varName))
	case "sint64":
		builder.WriteString(fmt.Sprintf("        message.writeSint64(stream, %d, %s);\n", field.FieldNumber, varName))
	case "fixed32", "sfixed32":
		builder.WriteString(fmt.Sprintf("        message.writeFixed32(stream, %d, %s);\n", field.FieldNumber, varName))
	case "fixed64", "sfixed64":
		builder.WriteString(fmt.Sprintf("        message.writeFixed64(stream, %d, %s);\n", field.FieldNumber, varName))
	case "bool":
		builder.WriteString(fmt.Sprintf("        message.writeBool(stream, %d, %s);\n", field.FieldNumber, varName))
	case "string":
		builder.WriteString(fmt.Sprintf("        message.writeString(stream, %d, %s);\n", field.FieldNumber, varName))
	case "float":
		builder.WriteString(fmt.Sprintf("        message.writeFloat(stream, %d, %s);\n", field.FieldNumber, varName))
	case "double":
		builder.WriteString(fmt.Sprintf("        message.writeDouble(stream, %d, %s);\n", field.FieldNumber, varName))
	case "bytes":
		builder.WriteString(fmt.Sprintf("        message.writeTag(stream, %d, message.WIRETYPE_LENGTH_DELIMITED);\n", field.FieldNumber))
		builder.WriteString(fmt.Sprintf("        message.writeBytes(stream, %s);\n", varName))
	default:
		// 处理嵌套消息
		builder.WriteString(fmt.Sprintf("        if (%s != null) {\n", varName))
		builder.WriteString(fmt.Sprintf("            byte[] bytes = %s.toByteArray();\n", varName))
		builder.WriteString(fmt.Sprintf("            message.writeTag(stream, %d, message.WIRETYPE_LENGTH_DELIMITED);\n", field.FieldNumber))
		builder.WriteString("            message.writeBytes(stream, bytes);\n")
		builder.WriteString("        }\n")
	}
	return builder.String()
}

func (jp *JavaProtoc) generateParseFrom(msg *protoc.Message) string {
	var builder strings.Builder
	writeParseFromHeader(&builder, msg.Name)
	writeParseFromBody(&builder, msg)
	writeParseFromFooter(&builder)
	return builder.String()
}

func writeParseFromHeader(builder *strings.Builder, msgName string) {
	builder.WriteString("\n    public static " + toCamelCase(msgName, true) + " parseFrom(byte[] data) {\n")
	builder.WriteString("        java.io.ByteArrayInputStream stream = new java.io.ByteArrayInputStream(data);\n")
	builder.WriteString(fmt.Sprintf("        %s result = new %s();\n",
		toCamelCase(msgName, true), toCamelCase(msgName, true)))
	builder.WriteString("        byte[] bytes;\n")
	builder.WriteString("        try {\n")
	writeParseLoop(builder)
}

func writeParseLoop(builder *strings.Builder) {
	builder.WriteString("            while (stream.available() > 0) {\n")
	builder.WriteString("                int tag = message.readTag(stream);\n")
	builder.WriteString("                int fieldNumber = message.getFieldNumberFromTag(tag);\n")
	builder.WriteString("                int wireType = message.getWireTypeFromTag(tag);\n")
	builder.WriteString("                switch (fieldNumber) {\n")
}

func writeParseFromBody(builder *strings.Builder, msg *protoc.Message) {
	writeFieldCases(builder, msg.Fields)
	writeOneofCases(builder, msg.OneOfs)
	writeDefaultCase(builder)
}

func writeFieldCases(builder *strings.Builder, fields []*protoc.Field) {
	for _, field := range fields {
		builder.WriteString(fmt.Sprintf("                    case %d:\n", field.FieldNumber))
		if field.MapInfo != nil {
			builder.WriteString(generateReadMapField(field))
		} else {
			builder.WriteString(generateReadField(field))
		}
		builder.WriteString("                        break;\n")
	}
}

func writeOneofCases(builder *strings.Builder, oneofs []*protoc.OneOf) {
	for _, oneof := range oneofs {
		for _, f := range oneof.Fields {
			builder.WriteString(fmt.Sprintf("                    case %d: // oneof %s\n",
				f.FieldNumber, oneof.Name))
			builder.WriteString(generateOneofReadField(f, oneof))
			builder.WriteString("                        break;\n")
		}
	}
}

func writeDefaultCase(builder *strings.Builder) {
	builder.WriteString("                    default:\n")
	builder.WriteString("                        if (wireType == message.WIRETYPE_LENGTH_DELIMITED) {\n")
	builder.WriteString("                            message.readBytes(stream);\n")
	builder.WriteString("                        }\n")
	builder.WriteString("                        break;\n")
	builder.WriteString("                }\n")
	builder.WriteString("            }\n")
}

func writeParseFromFooter(builder *strings.Builder) {
	builder.WriteString("        } catch (Exception e) {\n")
	builder.WriteString("            throw new RuntimeException(\"Failed to parse message\", e);\n")
	builder.WriteString("        }\n")
	builder.WriteString("        return result;\n")
	builder.WriteString("    }\n")
}
func generateOneofReadField(field *protoc.Field, oneof *protoc.OneOf) string {
	var builder strings.Builder
	oneofName := toCamelCase(oneof.Name, false)

	switch field.TypeName {
	case "int32", "uint32":
		builder.WriteString(fmt.Sprintf("                        result.%s = message.readInt32(stream);\n", oneofName))
	case "int64", "uint64":
		builder.WriteString(fmt.Sprintf("                        result.%s = message.readInt64(stream);\n", oneofName))
	case "sint32":
		builder.WriteString(fmt.Sprintf("                        result.%s = message.readSint32(stream);\n", oneofName))
	case "sint64":
		builder.WriteString(fmt.Sprintf("                        result.%s = message.readSint64(stream);\n", oneofName))
	case "bool":
		builder.WriteString(fmt.Sprintf("                        result.%s = message.readBool(stream);\n", oneofName))
	case "string":
		builder.WriteString(fmt.Sprintf("                        result.%s = message.readString(stream);\n", oneofName))
	default:
		builder.WriteString(fmt.Sprintf("                        byte[] bytes = message.readBytes(stream);\n"))
		builder.WriteString(fmt.Sprintf("                        result.%s = %s.parseFrom(bytes);\n",
			oneofName, toCamelCase(field.TypeName, true)))
	}
	builder.WriteString(fmt.Sprintf("                        result.%sCase = %d;\n", oneofName, field.FieldNumber))

	return builder.String()
}
func generateMapField(varName string, field *protoc.Field) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("        if (%s != null) {\n", varName))
	builder.WriteString(fmt.Sprintf("            for (java.util.Map.Entry<%s, %s> entry : %s.entrySet()) {\n",
		boxed(toJavaType(&protoc.Field{TypeName: field.MapInfo.KeyType})),
		boxed(toJavaType(&protoc.Field{TypeName: field.MapInfo.ValueType})),
		varName))
	builder.WriteString("                message.writeTag(stream, " + strconv.Itoa(field.FieldNumber) + ", message.WIRETYPE_LENGTH_DELIMITED);\n")
	builder.WriteString("                java.io.ByteArrayOutputStream mapStream = new java.io.ByteArrayOutputStream();\n")

	// 写入 key
	keyField := &protoc.Field{TypeName: field.MapInfo.KeyType, FieldNumber: 1}
	builder.WriteString(generateWriteField("entry.getKey()", keyField))

	// 写入 value
	valueField := &protoc.Field{TypeName: field.MapInfo.ValueType, FieldNumber: 2}
	builder.WriteString(generateWriteField("entry.getValue()", valueField))

	builder.WriteString("                byte[] mapBytes = mapStream.toByteArray();\n")
	builder.WriteString("                message.writeBytes(stream, mapBytes);\n")
	builder.WriteString("            }\n")
	builder.WriteString("        }\n")
	return builder.String()
}

func generateReadField(field *protoc.Field) string {
	var builder strings.Builder
	fieldName := toCamelCase(field.Name, false)
	if field.Repeated {
		switch field.TypeName {
		case "int32", "uint32":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(message.readInt32(stream));\n", fieldName))
		case "int64", "uint64":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(message.readInt64(stream));\n", fieldName))
		case "sint32":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(message.readSint32(stream));\n", fieldName))
		case "sint64":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(message.readSint64(stream));\n", fieldName))
		case "fixed32", "sfixed32":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(message.readFixed32(stream));\n", fieldName))
		case "fixed64", "sfixed64":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(message.readFixed64(stream));\n", fieldName))
		case "bool":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(message.readBool(stream));\n", fieldName))
		case "string":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(message.readString(stream));\n", fieldName))
		case "float":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(message.readFloat(stream));\n", fieldName))
		case "double":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(message.readDouble(stream));\n", fieldName))

		default:
			// 处理嵌套消息
			builder.WriteString(fmt.Sprintf("                    bytes = message.readBytes(stream);\n"))
			builder.WriteString(fmt.Sprintf("                    result.%s.add(%s.parseFrom(bytes));\n", fieldName, toCamelCase(field.TypeName, true)))
		}
	} else {
		switch field.TypeName {
		case "int32", "uint32":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readInt32(stream);\n", fieldName))
		case "int64", "uint64":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readInt64(stream);\n", fieldName))
		case "sint32":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readSint32(stream);\n", fieldName))
		case "sint64":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readSint64(stream);\n", fieldName))
		case "fixed32", "sfixed32":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readFixed32(stream);\n", fieldName))
		case "fixed64", "sfixed64":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readFixed64(stream);\n", fieldName))
		case "bool":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readBool(stream);\n", fieldName))
		case "string":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readString(stream);\n", fieldName))
		case "float":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readFloat(stream);\n", fieldName))
		case "double":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readDouble(stream);\n", fieldName))
		case "bytes":
			builder.WriteString(fmt.Sprintf("                    result.%s = message.readBytes(stream);\n", fieldName))
		default:
			// 处理嵌套消息
			builder.WriteString(fmt.Sprintf("                    bytes = message.readBytes(stream);\n"))
			builder.WriteString(fmt.Sprintf("                    result.%s = %s.parseFrom(bytes);\n", fieldName, toCamelCase(field.TypeName, true)))
		}
	}
	return builder.String()
}

func generateReadMapField(field *protoc.Field) string {
	var builder strings.Builder
	fieldName := toCamelCase(field.Name, false)

	builder.WriteString("                    byte[] mapBytes = message.readBytes(stream);\n")
	builder.WriteString("                    java.io.ByteArrayInputStream mapStream = new java.io.ByteArrayInputStream(mapBytes);\n")

	keyType := toJavaType(&protoc.Field{TypeName: field.MapInfo.KeyType})
	valueType := toJavaType(&protoc.Field{TypeName: field.MapInfo.ValueType})

	builder.WriteString(fmt.Sprintf("                    %s key = %s;\n", keyType, getDefaultValueByStr(field.MapInfo.KeyType)))
	builder.WriteString(fmt.Sprintf("                    %s value = %s;\n", valueType, getDefaultValueByStr(field.MapInfo.ValueType)))

	builder.WriteString("                    while (mapStream.available() > 0) {\n")
	builder.WriteString("                        int mapTag = message.readTag(mapStream);\n")
	builder.WriteString("                        int mapFieldNumber = message.getFieldNumberFromTag(mapTag);\n")
	builder.WriteString("                        switch (mapFieldNumber) {\n")
	builder.WriteString("                            case 1: // key\n")
	builder.WriteString(generateMapKeyRead(field.MapInfo.KeyType))
	builder.WriteString("                                break;\n")
	builder.WriteString("                            case 2: // value\n")
	builder.WriteString(generateMapValueRead(field.MapInfo.ValueType))
	builder.WriteString("                                break;\n")
	builder.WriteString("                            default:\n")
	builder.WriteString("                                break;\n")
	builder.WriteString("                        }\n")
	builder.WriteString("                    }\n")
	builder.WriteString(fmt.Sprintf("                    result.%s.put(key, value);\n", fieldName))

	return builder.String()
}

func generateMapKeyRead(keyType string) string {
	switch keyType {
	case "int32", "uint32":
		return "                                key = message.readInt32(mapStream);\n"
	case "int64", "uint64":
		return "                                key = message.readInt64(mapStream);\n"
	case "sint32":
		return "                                key = message.readSint32(mapStream);\n"
	case "sint64":
		return "                                key = message.readSint64(mapStream);\n"
	case "fixed32", "sfixed32":
		return "                                key = message.readFixed32(mapStream);\n"
	case "fixed64", "sfixed64":
		return "                                key = message.readFixed64(mapStream);\n"
	case "string":
		return "                                key = message.readString(mapStream);\n"
	case "bool":
		return "                                key = message.readBool(mapStream);\n"
	default:
		return fmt.Sprintf("                                // Unsupported map key type: %s\n", keyType)
	}
}

func generateMapValueRead(valueType string) string {
	switch valueType {
	case "int32", "uint32":
		return "                                value = message.readInt32(mapStream);\n"
	case "int64", "uint64":
		return "                                value = message.readInt64(mapStream);\n"
	case "sint32":
		return "                                value = message.readSint32(mapStream);\n"
	case "sint64":
		return "                                value = message.readSint64(mapStream);\n"
	case "fixed32", "sfixed32":
		return "                                value = message.readFixed32(mapStream);\n"
	case "fixed64", "sfixed64":
		return "                                value = message.readFixed64(mapStream);\n"
	case "bool":
		return "                                value = message.readBool(mapStream);\n"
	case "string":
		return "                                value = message.readString(mapStream);\n"
	case "float":
		return "                                value = message.readFloat(mapStream);\n"
	case "double":
		return "                                value = message.readDouble(mapStream);\n"
	default:
		// 处理嵌套消息
		return fmt.Sprintf("                                byte[] bytes = message.readBytes(mapStream);\n"+
			"                                value = %s.parseFrom(bytes);\n", toCamelCase(valueType, true))
	}
}

func getElementType(field *protoc.Field) string {
	if field.MapInfo != nil {
		keyType := toJavaType(&protoc.Field{TypeName: field.MapInfo.KeyType})
		valueType := toJavaType(&protoc.Field{TypeName: field.MapInfo.ValueType})
		return fmt.Sprintf("java.util.Map.Entry<%s, %s>", boxed(keyType), boxed(valueType))
	}
	return boxed(toJavaType(&protoc.Field{TypeName: field.TypeName}))
}
