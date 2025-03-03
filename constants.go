package main

// 类型相关常量
const (
	// 基本类型
	TypeInt32    = "int32"
	TypeInt64    = "int64"
	TypeUint32   = "uint32"
	TypeUint64   = "uint64"
	TypeSint32   = "sint32"
	TypeSint64   = "sint64"
	TypeBool     = "bool"
	TypeFloat    = "float"
	TypeDouble   = "double"
	TypeString   = "string"
	TypeBytes    = "bytes"
	TypeFixed32  = "fixed32"
	TypeFixed64  = "fixed64"
	TypeSfixed32 = "sfixed32"
	TypeSfixed64 = "sfixed64"
	TypeEnum     = "enum"

	// 关键字
	KeywordMessage  = "message"
	KeywordEnum     = "enum"
	KeywordService  = "service"
	KeywordRepeated = "repeated"
	KeywordMap      = "map"
	KeywordOneof    = "oneof"
	KeywordSyntax   = "syntax"
	KeywordPackage  = "package"
	KeywordImport   = "import"
	KeywordPublic   = "public"
	KeywordStream   = "stream"
	KeywordReturns  = "returns"

	// 符号
	SymbolSemicolon    = ";"
	SymbolEqual        = "="
	SymbolLeftBrace    = "{"
	SymbolRightBrace   = "}"
	SymbolLeftBracket  = "["
	SymbolRightBracket = "]"
	SymbolLeftParen    = "("
	SymbolRightParen   = ")"
	SymbolComma        = ","
	SymbolLessThan     = "<"
	SymbolGreaterThan  = ">"

	// 选项
	OptionDeprecated = "deprecated"
	OptionPacked     = "packed"

	// 默认值
	DefaultTrue  = "true"
	DefaultFalse = "false"
)

// 错误信息常量
const (
	ErrUnterminatedString = "unterminated string"
	ErrUnexpectedToken    = "unexpected token: %v"
	ErrInvalidOptionValue = "invalid option value: %v"
)

// Java 相关常量
const (
	JavaByteArrayOutputStream = "java.io.ByteArrayOutputStream"
	JavaByteArrayInputStream  = "java.io.ByteArrayInputStream"
	JavaMapEntry              = "java.util.Map.Entry"
	JavaHashMap               = "java.util.HashMap"
	JavaArrayList             = "java.util.ArrayList"
	JavaList                  = "java.util.List"
	JavaMap                   = "java.util.Map"
)
