package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewJavaProtoc(t *testing.T) {
	protoc, err := newJavaProtoc("", "./proto/example.proto")
	if err != nil {
		t.Fatal(err)
	}
	// print protoc
	protoJson, _ := json.MarshalIndent(protoc, "", "\t")
	fmt.Println(string(protoJson))
}

func TestGenerateMessageClass(t *testing.T) {
	protoc, err := newJavaProtoc("", "./proto/example.proto")
	if err != nil {
		t.Fatal(err)
	}
	messageClass := protoc.generateMessageClass(protoc.Messages[0], false)
	// print messageClass
	fmt.Println(messageClass)
}

func TestGenerateMessageClass2(t *testing.T) {
	protoc, err := newJavaProtoc("", "./proto/example.proto")
	if err != nil {
		t.Fatal(err)
	}
	messageClass := protoc.generateMessageClass(protoc.Messages[1], false)
	// print messageClass
	fmt.Println(messageClass)
}

func TestGenerateEnum(t *testing.T) {
	protoc, err := newJavaProtoc("", "./proto/example.proto")
	if err != nil {
		t.Fatal(err)
	}
	enum := protoc.Enums[0]
	enumClass := protoc.generateEnum("", enum)
	// print enumClass
	fmt.Println(enumClass)
}

func TestToJavaType(t *testing.T) {
	fields := []*Field{
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
