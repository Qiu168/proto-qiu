package java

import (
	"fmt"
	"os"
	"path/filepath"
	"proto-qiu/generator"
	"proto-qiu/protoc"
	"strings"
)

var _ generator.Generator = (*JavaProtoc)(nil)

type JavaProtoc struct {
	protoc.Protoc
	JavaOutput    string
	ProtoFilePath string
}

func NewJavaProtoc(javaOutput, protoFilePath string) (*JavaProtoc, error) {
	// read file
	content, err := os.ReadFile(protoFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read proto file: %v", err)
	}

	// 解析.proto文件
	parser := protoc.NewParser(strings.NewReader(string(content)))
	protoc, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("proto parse error: %v", err)
	}

	fillFieldType(protoc)

	protoc.ProtoName = strings.Split(filepath.Base(protoFilePath), ".")[0]
	return &JavaProtoc{
		Protoc:        *protoc,
		JavaOutput:    javaOutput,
		ProtoFilePath: protoFilePath,
	}, nil

}

func fillFieldType(p *protoc.Protoc) {
	for _, message := range p.Messages {
	LOOP:
		for _, field := range message.Fields {
			switch field.TypeName {
			case "int32", "uint32", "sint32", "fixed32", "sfixed32", "int64", "uint64", "sint64", "fixed64", "sfixed64", "string", "double", "float", "bytes", "bool":
				field.Type = protoc.BASE
			default:
				if field.MapInfo != nil {
					field.Type = protoc.MAP
					continue LOOP
				}
				for message != nil {
					for _, enum := range message.Enums {
						if enum.Name == field.TypeName {
							field.Type = protoc.ENUM
							continue LOOP
						}
					}
					message = message.SuperMessage
				}
				field.Type = protoc.CUSTOM
			}
		}
	}
}

// Generate .java file
func (jp *JavaProtoc) Generate() error {
	// 创建包对应的目录
	packagePath := filepath.Join(jp.JavaOutput, strings.Replace(jp.PackageName, ".", "/", -1))
	if err := os.MkdirAll(packagePath, 0777); err != nil {
		return fmt.Errorf("failed to create package dir: %v", err)
	}

	var innerStr strings.Builder

	// 为每个 Message 生成 Java 类
	for _, msg := range jp.Messages {
		innerStr.WriteString(jp.generateMessageClass(msg, true))
	}

	// 为每个 Enum 生成 Java 枚举
	for _, enum := range jp.Enums {
		innerStr.WriteString(jp.generateEnum(enum))
	}

	// 生成外部类
	fileStr := jp.generateOuterClass(innerStr)

	// get proto file Name
	javaFilePath := packagePath + "\\" + toCamelCase(jp.ProtoName, true) + ".java"

	file, err := os.OpenFile(javaFilePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("无法打开或创建文件:", err)
		return err
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		fmt.Println("无法清空文件:", err)
		return err
	}
	_, err = file.WriteString(fileStr)

	if err != nil {
		return fmt.Errorf("failed to write proto file: %v", err)
	}

	return nil
}
