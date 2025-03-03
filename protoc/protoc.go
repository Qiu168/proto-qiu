package protoc

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
