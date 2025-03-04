package protoc

import (
	"fmt"
	"os"
	"strings"
)

type WireType int

type Type int

func str2WireType(s string) WireType {
	switch s {
	case "int32", "int64", "uint32", "uint64", "bool", "enum", "sint32", "sint64":
		return Varint
	case "fixed64", "sfixed64", "double":
		return Fixed64
	case "fixed32", "sfixed32", "float":
		return Fixed32
	default:
		return LengthDelimited
	}
}

const (
	Varint WireType = iota
	Fixed64
	LengthDelimited
	// StartGroup Deprecated
	StartGroup
	// EndGroup Deprecated
	EndGroup
	Fixed32
)
const (
	// BASE 基本类型
	BASE Type = iota
	// ENUM 枚举类型
	ENUM
	// CUSTOM 自定义类型
	CUSTOM
	// ONE_OF oneof类型
	ONE_OF
	// MAP 映射类型
	MAP
	// MESSAGE 消息类型
	MESSAGE
)

type MapInfo struct {
	KeyType   string
	ValueType string
}

type Field struct {
	Name        string
	TypeName    string // 基础类型或自定义类型（如 "int32", "MyMessage"）
	Type        Type
	WireType    WireType
	FieldNumber int
	Repeated    bool
	MapInfo     *MapInfo
	Options     *FieldOptions
}

type FieldOptions struct {
	Deprecated bool
	Packed     bool
}

type OneOf struct {
	Name   string
	Fields []*Field
}

type Message struct {
	Name          string
	SuperMessage  *Message `json:"-"`
	InnerMessages []*Message
	OneOfs        []*OneOf
	Fields        []*Field
	Enums         []*Enum
}

type Enum struct {
	Name         string
	SuperMessage *Message `json:"-"`
	Values       []*EnumValue
}

type EnumValue struct {
	Name  string
	Value int
}

type Service struct {
	Name    string
	Methods []*Method
}

type Method struct {
	Name            string
	InputType       string
	OutputType      string
	ClientStreaming bool
	ServerStreaming bool
}

type Protoc struct {
	ProtoName     string
	SyntaxVersion string
	PackageName   string
	Imports       []*Import
	Messages      []*Message
	Enums         []*Enum
	Services      []*Service
}

type Import struct {
	Path   string
	Public bool
}

func NewProtoc(protoFilePath string) (*Protoc, error) {
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

	protoc.fillFieldType()

	return protoc, nil
}

func (p *Protoc) fillFieldType() {
	fillMessageFieldType(p.Messages, p.Enums)
}

func fillMessageFieldType(messages []*Message, enums []*Enum) {
	for _, message := range messages {
		// 保存原始 message 引用，因为后面会修改 message 变量
		originalMessage := message

		// 处理当前消息的字段
		for _, field := range message.Fields {
			field.Type = getFieldType(field, originalMessage, enums)
		}

		// 递归处理内部消息
		fillMessageFieldType(message.InnerMessages, nil)
	}
}

func getFieldType(field *Field, message *Message, enums []*Enum) Type {
	// 处理基本类型
	switch field.TypeName {
	case "int32", "uint32", "sint32", "fixed32", "sfixed32",
		"int64", "uint64", "sint64", "fixed64", "sfixed64",
		"string", "double", "float", "bytes", "bool":
		return BASE
	}

	// 处理 map 类型
	if field.MapInfo != nil {
		return MAP
	}

	// 在消息层次结构中查找枚举类型
	currentMsg := message
	for currentMsg != nil {
		for _, enum := range currentMsg.Enums {
			if enum.Name == field.TypeName {
				return ENUM
			}
		}
		currentMsg = currentMsg.SuperMessage
	}
	if enums != nil && len(enums) > 0 {
		for _, enum := range enums {
			if enum.Name == field.TypeName {
				return ENUM
			}
		}
	}

	// 如果都不匹配，则为自定义类型
	return CUSTOM
}
