package constant

// proto相关常量
const (
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

	OptionDeprecated = "deprecated"
	OptionPacked     = "packed"

	DefaultTrue  = "true"
	DefaultFalse = "false"

	ProtoFileSuffix = ".proto"

	ProtoUsage = "Usage: %s -java_out=[Path] [args...]\n"
)
