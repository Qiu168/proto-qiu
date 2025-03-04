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
	writeOneOfs(&builder, msg.OneOfs)
	writeMethodFooter(&builder)
	return builder.String()
}

// ParseFrom
func (jp *JavaProtoc) generateParseFrom(msg *protoc.Message) string {
	var builder strings.Builder
	writeParseFromHeader(&builder, msg.Name)
	writeParseFromBody(&builder, msg)
	writeParseFromFooter(&builder)
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
		if field.Type == protoc.ENUM {
			writeEnumField(builder, fieldName, field)
		} else if field.MapInfo != nil {
			writeMapField(builder, fieldName, field)
		} else if field.Repeated {
			writeRepeatedField(builder, fieldName, field)
		} else {
			writeSimpleField(builder, fieldName, field)
		}
	}
}

func writeEnumField(builder *strings.Builder, fieldName string, field *protoc.Field) {
	builder.WriteString(fmt.Sprintf("            if (%s != null) {\n", fieldName))
	builder.WriteString(fmt.Sprintf("                writeInt32(stream, %d, %s.getNumber());\n",
		field.FieldNumber, fieldName))
	builder.WriteString("            }\n")
}

func writeSimpleField(builder *strings.Builder, fieldName string, field *protoc.Field) {
	builder.WriteString(fmt.Sprintf("            if (%s != %s) {\n",
		fieldName, getDefaultValue(field)))
	builder.WriteString(generateWriteField(fieldName, field))
	builder.WriteString("            }\n")
}

func writeOneOfs(builder *strings.Builder, oneofs []*protoc.OneOf) {
	for _, oneOf := range oneofs {
		writeOneofField(builder, oneOf)
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
	builder.WriteString(fmt.Sprintf("        write%s(stream, %d, %s);\n",
		toCamelCase(typeName, true), fieldNumber, fieldName))
}

func writeMessageField(builder *strings.Builder, fieldName string, fieldNumber int) {
	builder.WriteString(fmt.Sprintf("        if (%s != null) {\n", fieldName))
	builder.WriteString(fmt.Sprintf("            byte[] bytes = %s.toByteArray();\n", fieldName))
	builder.WriteString(fmt.Sprintf("            writeTag(stream, %d, WIRETYPE_LENGTH_DELIMITED);\n", fieldNumber))
	builder.WriteString("            writeBytes(stream, bytes);\n")
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
	builder.WriteString("                writeTag(stream, " +
		strconv.Itoa(field.FieldNumber) + ", WIRETYPE_LENGTH_DELIMITED);\n")
	builder.WriteString("                java.io.ByteArrayOutputStream mapStream = new java.io.ByteArrayOutputStream();\n")

	keyField := &protoc.Field{TypeName: field.MapInfo.KeyType, FieldNumber: 1}
	valueField := &protoc.Field{TypeName: field.MapInfo.ValueType, FieldNumber: 2}

	builder.WriteString(generateWriteField("entry.getKey()", keyField))
	builder.WriteString(generateWriteField("entry.getValue()", valueField))
	builder.WriteString("                byte[] mapBytes = mapStream.toByteArray();\n")
	builder.WriteString("                writeBytes(stream, mapBytes);\n")
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
		builder.WriteString(fmt.Sprintf("        writeInt32(stream, %d, %s);\n", field.FieldNumber, varName))
	case "int64", "uint64":
		builder.WriteString(fmt.Sprintf("        writeInt64(stream, %d, %s);\n", field.FieldNumber, varName))
	case "sint32":
		builder.WriteString(fmt.Sprintf("        writeSint32(stream, %d, %s);\n", field.FieldNumber, varName))
	case "sint64":
		builder.WriteString(fmt.Sprintf("        writeSint64(stream, %d, %s);\n", field.FieldNumber, varName))
	case "fixed32", "sfixed32":
		builder.WriteString(fmt.Sprintf("        writeFixed32(stream, %d, %s);\n", field.FieldNumber, varName))
	case "fixed64", "sfixed64":
		builder.WriteString(fmt.Sprintf("        writeFixed64(stream, %d, %s);\n", field.FieldNumber, varName))
	case "bool":
		builder.WriteString(fmt.Sprintf("        writeBool(stream, %d, %s);\n", field.FieldNumber, varName))
	case "string":
		builder.WriteString(fmt.Sprintf("        writeString(stream, %d, %s);\n", field.FieldNumber, varName))
	case "float":
		builder.WriteString(fmt.Sprintf("        writeFloat(stream, %d, %s);\n", field.FieldNumber, varName))
	case "double":
		builder.WriteString(fmt.Sprintf("        writeDouble(stream, %d, %s);\n", field.FieldNumber, varName))
	case "bytes":
		builder.WriteString(fmt.Sprintf("        writeTag(stream, %d, WIRETYPE_LENGTH_DELIMITED);\n", field.FieldNumber))
		builder.WriteString(fmt.Sprintf("        writeBytes(stream, %s);\n", varName))
	default:
		// 处理嵌套消息
		//builder.WriteString(fmt.Sprintf("        if (%s != null) {\n", varName))
		builder.WriteString(fmt.Sprintf("            byte[] bytes = %s.toByteArray();\n", varName))
		builder.WriteString(fmt.Sprintf("            writeTag(stream, %d, WIRETYPE_LENGTH_DELIMITED);\n", field.FieldNumber))
		builder.WriteString("            writeBytes(stream, bytes);\n")
		//builder.WriteString("        }\n")
	}
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
	builder.WriteString("                int tag = readTag(stream);\n")
	builder.WriteString("                int fieldNumber = getFieldNumberFromTag(tag);\n")
	builder.WriteString("                int wireType = getWireTypeFromTag(tag);\n")
	builder.WriteString("                switch (fieldNumber) {\n")
}

func writeParseFromBody(builder *strings.Builder, msg *protoc.Message) {
	writeFieldCases(builder, msg.Fields)
	writeOneOfCases(builder, msg.OneOfs)
	writeDefaultCase(builder)
}

func writeFieldCases(builder *strings.Builder, fields []*protoc.Field) {
	for _, field := range fields {
		builder.WriteString(fmt.Sprintf("                    case %d:\n", field.FieldNumber))
		if field.Type == protoc.ENUM {
			builder.WriteString(generateReadEnumField(field))
		} else if field.MapInfo != nil {
			builder.WriteString(generateReadMapField(field))
		} else {
			builder.WriteString(generateReadField(field))
		}
		builder.WriteString("                        break;\n")
	}
}

func generateReadEnumField(field *protoc.Field) string {
	var builder strings.Builder
	fieldName := toCamelCase(field.Name, false)

	if field.Repeated {
		builder.WriteString(fmt.Sprintf("                    result.%s.add(%s.forNumber(readInt32(stream)));\n",
			fieldName, toCamelCase(field.TypeName, true)))
	} else {
		builder.WriteString(fmt.Sprintf("                    result.%s = %s.forNumber(readInt32(stream));\n",
			fieldName, toCamelCase(field.TypeName, true)))
	}

	return builder.String()
}

func writeOneOfCases(builder *strings.Builder, oneofs []*protoc.OneOf) {
	for _, oneOf := range oneofs {
		for _, f := range oneOf.Fields {
			builder.WriteString(fmt.Sprintf("                    case %d: // oneof %s\n",
				f.FieldNumber, oneOf.Name))
			builder.WriteString(generateOneofReadField(f, oneOf))
			builder.WriteString("                        break;\n")
		}
	}
}

func writeDefaultCase(builder *strings.Builder) {
	builder.WriteString("                    default:\n")
	builder.WriteString("                        if (wireType == WIRETYPE_LENGTH_DELIMITED) {\n")
	builder.WriteString("                            readBytes(stream);\n")
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
		builder.WriteString(fmt.Sprintf("                        result.%s = readInt32(stream);\n", oneofName))
	case "int64", "uint64":
		builder.WriteString(fmt.Sprintf("                        result.%s = readInt64(stream);\n", oneofName))
	case "sint32":
		builder.WriteString(fmt.Sprintf("                        result.%s = readSint32(stream);\n", oneofName))
	case "sint64":
		builder.WriteString(fmt.Sprintf("                        result.%s = readSint64(stream);\n", oneofName))
	case "bool":
		builder.WriteString(fmt.Sprintf("                        result.%s = readBool(stream);\n", oneofName))
	case "string":
		builder.WriteString(fmt.Sprintf("                        result.%s = readString(stream);\n", oneofName))
	default:
		builder.WriteString(fmt.Sprintf("                        byte[] bytes = readBytes(stream);\n"))
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
	builder.WriteString("                writeTag(stream, " + strconv.Itoa(field.FieldNumber) + ", WIRETYPE_LENGTH_DELIMITED);\n")
	builder.WriteString("                java.io.ByteArrayOutputStream mapStream = new java.io.ByteArrayOutputStream();\n")

	// 写入 key
	keyField := &protoc.Field{TypeName: field.MapInfo.KeyType, FieldNumber: 1}
	builder.WriteString(generateWriteField("entry.getKey()", keyField))

	// 写入 value
	valueField := &protoc.Field{TypeName: field.MapInfo.ValueType, FieldNumber: 2}
	builder.WriteString(generateWriteField("entry.getValue()", valueField))

	builder.WriteString("                byte[] mapBytes = mapStream.toByteArray();\n")
	builder.WriteString("                writeBytes(stream, mapBytes);\n")
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
			builder.WriteString(fmt.Sprintf("                    result.%s.add(readInt32(stream));\n", fieldName))
		case "int64", "uint64":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(readInt64(stream));\n", fieldName))
		case "sint32":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(readSint32(stream));\n", fieldName))
		case "sint64":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(readSint64(stream));\n", fieldName))
		case "fixed32", "sfixed32":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(readFixed32(stream));\n", fieldName))
		case "fixed64", "sfixed64":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(readFixed64(stream));\n", fieldName))
		case "bool":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(readBool(stream));\n", fieldName))
		case "string":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(readString(stream));\n", fieldName))
		case "float":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(readFloat(stream));\n", fieldName))
		case "double":
			builder.WriteString(fmt.Sprintf("                    result.%s.add(readDouble(stream));\n", fieldName))

		default:
			// 处理嵌套消息
			builder.WriteString(fmt.Sprintf("                    bytes = readBytes(stream);\n"))
			builder.WriteString(fmt.Sprintf("                    result.%s.add(%s.parseFrom(bytes));\n", fieldName, toCamelCase(field.TypeName, true)))
		}
	} else {
		switch field.TypeName {
		case "int32", "uint32":
			builder.WriteString(fmt.Sprintf("                    result.%s = readInt32(stream);\n", fieldName))
		case "int64", "uint64":
			builder.WriteString(fmt.Sprintf("                    result.%s = readInt64(stream);\n", fieldName))
		case "sint32":
			builder.WriteString(fmt.Sprintf("                    result.%s = readSint32(stream);\n", fieldName))
		case "sint64":
			builder.WriteString(fmt.Sprintf("                    result.%s = readSint64(stream);\n", fieldName))
		case "fixed32", "sfixed32":
			builder.WriteString(fmt.Sprintf("                    result.%s = readFixed32(stream);\n", fieldName))
		case "fixed64", "sfixed64":
			builder.WriteString(fmt.Sprintf("                    result.%s = readFixed64(stream);\n", fieldName))
		case "bool":
			builder.WriteString(fmt.Sprintf("                    result.%s = readBool(stream);\n", fieldName))
		case "string":
			builder.WriteString(fmt.Sprintf("                    result.%s = readString(stream);\n", fieldName))
		case "float":
			builder.WriteString(fmt.Sprintf("                    result.%s = readFloat(stream);\n", fieldName))
		case "double":
			builder.WriteString(fmt.Sprintf("                    result.%s = readDouble(stream);\n", fieldName))
		case "bytes":
			builder.WriteString(fmt.Sprintf("                    result.%s = readBytes(stream);\n", fieldName))
		default:
			// 处理嵌套消息
			builder.WriteString(fmt.Sprintf("                    bytes = readBytes(stream);\n"))
			builder.WriteString(fmt.Sprintf("                    result.%s = %s.parseFrom(bytes);\n", fieldName, toCamelCase(field.TypeName, true)))
		}
	}
	return builder.String()
}

func generateReadMapField(field *protoc.Field) string {
	var builder strings.Builder
	fieldName := toCamelCase(field.Name, false)

	builder.WriteString("                    byte[] mapBytes = readBytes(stream);\n")
	builder.WriteString("                    java.io.ByteArrayInputStream mapStream = new java.io.ByteArrayInputStream(mapBytes);\n")

	keyType := toJavaType(&protoc.Field{TypeName: field.MapInfo.KeyType})
	valueType := toJavaType(&protoc.Field{TypeName: field.MapInfo.ValueType})

	builder.WriteString(fmt.Sprintf("                    %s key = %s;\n", keyType, getDefaultValueByStr(field.MapInfo.KeyType)))
	builder.WriteString(fmt.Sprintf("                    %s value = %s;\n", valueType, getDefaultValueByStr(field.MapInfo.ValueType)))

	builder.WriteString("                    while (mapStream.available() > 0) {\n")
	builder.WriteString("                        int mapTag = readTag(mapStream);\n")
	builder.WriteString("                        int mapFieldNumber = getFieldNumberFromTag(mapTag);\n")
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
		return "                                key = readInt32(mapStream);\n"
	case "int64", "uint64":
		return "                                key = readInt64(mapStream);\n"
	case "sint32":
		return "                                key = readSint32(mapStream);\n"
	case "sint64":
		return "                                key = readSint64(mapStream);\n"
	case "fixed32", "sfixed32":
		return "                                key = readFixed32(mapStream);\n"
	case "fixed64", "sfixed64":
		return "                                key = readFixed64(mapStream);\n"
	case "string":
		return "                                key = readString(mapStream);\n"
	case "bool":
		return "                                key = readBool(mapStream);\n"
	default:
		return fmt.Sprintf("                                // Unsupported map key type: %s\n", keyType)
	}
}

func generateMapValueRead(valueType string) string {
	switch valueType {
	case "int32", "uint32":
		return "                                value = readInt32(mapStream);\n"
	case "int64", "uint64":
		return "                                value = readInt64(mapStream);\n"
	case "sint32":
		return "                                value = readSint32(mapStream);\n"
	case "sint64":
		return "                                value = readSint64(mapStream);\n"
	case "fixed32", "sfixed32":
		return "                                value = readFixed32(mapStream);\n"
	case "fixed64", "sfixed64":
		return "                                value = readFixed64(mapStream);\n"
	case "bool":
		return "                                value = readBool(mapStream);\n"
	case "string":
		return "                                value = readString(mapStream);\n"
	case "float":
		return "                                value = readFloat(mapStream);\n"
	case "double":
		return "                                value = readDouble(mapStream);\n"
	default:
		// 处理嵌套消息
		return fmt.Sprintf("                                byte[] bytes = readBytes(mapStream);\n"+
			"                                value = %s.parseFrom(bytes);\n", toCamelCase(valueType, true))
	}
}
