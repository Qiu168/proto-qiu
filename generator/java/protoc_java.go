package java

import (
	"fmt"
	"os"
	"path/filepath"
	"proto-qiu/constant"
	"proto-qiu/generator"
	"proto-qiu/protoc"
	"strings"
)

var _ generator.Generator = (*JavaProtoc)(nil)

type JavaProtoc struct {
	*protoc.Protoc
	JavaOutput    string
	ProtoFilePath string
}

func NewJavaProtoc(javaOutput, protoFilePath string) (*JavaProtoc, error) {
	proto, err := protoc.NewProtoc(protoFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proto file: %v", err)
	}

	proto.ProtoName = strings.Split(filepath.Base(protoFilePath), ".")[0]
	return &JavaProtoc{
		Protoc:        proto,
		JavaOutput:    javaOutput,
		ProtoFilePath: protoFilePath,
	}, nil

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
	javaFilePath := packagePath + "\\" + toCamelCase(jp.ProtoName, true) + constant.JavaFileSuffix

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
