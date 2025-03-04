package main

import (
	"fmt"
	"os"
	"path/filepath"
	"proto-qiu/constant"
	"proto-qiu/generator/java"
	"strings"
)

func main() {
	cmd := parseCmd()
	if cmd.version {
		fmt.Println(constant.QiuProtoVersion)
	} else {
		fmt.Printf("Starting proto-qiu compiler version %s\n", constant.QiuProtoVersion)
		// getAbsolutePath, _ := filepath.Abs(cmd.javaOutput)
		dir, _ := filepath.Abs(cmd.javaOutput)
		fmt.Printf("Output directory: %s\n", dir)

		var protoPaths []string
		for _, path := range cmd.protocPath {
			dir, _ = filepath.Abs(path)
			if strings.HasSuffix(path, constant.ProtoFileSuffix) {
				protoPaths = append(protoPaths, path)
				fmt.Printf("Found proto file: %s\n", dir)
			} else {
				fmt.Printf("Scanning directory: %s\n", dir)
				paths := getDirFiles(path)
				for _, p := range paths {
					if strings.HasSuffix(p, constant.ProtoFileSuffix) {
						protoPaths = append(protoPaths, p)
						fmt.Printf("Found proto file: %s\n", p)
					}
				}
			}
		}
		if len(protoPaths) == 0 {
			fmt.Println("No .proto files found!")
			return
		}

		fmt.Printf("\nProcessing %d proto files...\n", len(protoPaths))
		for i, path := range protoPaths {
			fmt.Printf("\n[%d/%d] Compiling: %s\n", i+1, len(protoPaths), path)
			javaProto, err := java.NewJavaProtoc(cmd.javaOutput, path)
			if err != nil {
				fmt.Printf("Error parsing proto file: %v\n", err)
				panic(fmt.Errorf("parse protoc error: %v", err))
			}
			err = javaProto.Generate()
			if err != nil {
				fmt.Printf("Error generating Java code: %v\n", err)
				panic(err)
			}
			fmt.Printf("Successfully compiled: %s\n", path)
		}
		fmt.Println("\nCompilation completed successfully!")
	}
}

func getDirFiles(path string) []string {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
