package protoc

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestNewJavaProtoc(t *testing.T) {
	os.Chdir("../")
	proto, err := NewProtoc("./proto/example.proto")
	if err != nil {
		t.Fatal(err)
	}
	// print protoc
	protoJson, _ := json.MarshalIndent(proto, "", "\t")
	fmt.Println(string(protoJson))
}
