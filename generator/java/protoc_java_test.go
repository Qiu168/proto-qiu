package java

import (
	"fmt"
	"os"
	"proto-qiu/protoc"
	"strings"
	"testing"
)

func TestToJavaType(t *testing.T) {
	fields := []*protoc.Field{
		{
			TypeName: "string",
			Repeated: true,
		},
		{
			TypeName: "int32",
			Repeated: false,
		},
		{
			TypeName: "int64",
			Repeated: true,
		},
		{
			TypeName: "bool",
			Repeated: false,
		},
	}
	ans := []string{"java.util.List<java.lang.String>", "int", "java.util.List<java.lang.Long>", "boolean"}
	for i, field := range fields {
		if toJavaType(field) != ans[i] {
			t.Errorf("toJavaType(%v) = %v, want %v", field, toJavaType(field), ans[i])
		}
	}
}

func TestGenerateOuterClass(t *testing.T) {
	os.Chdir("../../")
	proto, err := NewJavaProtoc("", "./proto/example.proto")
	if err != nil {
		t.Fatal(err)
	}
	outerClass := proto.generateOuterClass(strings.Builder{})
	// print outerClass
	fmt.Println(outerClass)
}

func TestGenerateMessageClass(t *testing.T) {
	os.Chdir("../../")
	proto, err := NewJavaProtoc("", "./proto/example.proto")
	if err != nil {
		t.Fatal(err)
	}
	messageClass := proto.generateMessageClass(proto.Messages[0], false)
	// print messageClass
	fmt.Println(messageClass)
}

func TestGenerateMessageClass2(t *testing.T) {
	os.Chdir("../../")
	proto, err := NewJavaProtoc("", "./proto/example.proto")
	if err != nil {
		t.Fatal(err)
	}
	messageClass := proto.generateMessageClass(proto.Messages[1], false)
	// print messageClass
	fmt.Println(messageClass)
}

func TestGenerateEnum(t *testing.T) {
	os.Chdir("../../")
	proto, err := NewJavaProtoc("", "./proto/example.proto")
	if err != nil {
		t.Fatal(err)
	}
	enum := proto.Enums[0]
	enumClass := proto.generateEnum(enum)
	// print enumClass
	fmt.Println(enumClass)
}

func TestGenerate(t *testing.T) {
	os.Chdir("../../")
	fmt.Println(os.Getwd())
	proto, err := NewJavaProtoc("", "./proto/example.proto")
	if err != nil {
		t.Fatal(err)
	}
	err = proto.Generate()
	if err != nil {
		t.Fatal(err)
	}
}
