package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	cmd := parseCmd()
	if cmd.version {
		fmt.Println("version: 1.0.0")
	} else {
		var protoPaths []string
		for _, path := range cmd.protocPath {
			if strings.HasSuffix(path, ".proto") {
				protoPaths = append(protoPaths, path)
			} else {
				// 获取path目录下的所有文件
				paths := getDirFiles(path)
				for _, p := range paths {
					if strings.HasSuffix(p, ".proto") {
						protoPaths = append(protoPaths, p)
					}
				}
			}
		}
		for _, path := range protoPaths {
			protoc, err := newJavaProtoc(cmd.javaOutput, path)
			if err != nil {
				panic(fmt.Errorf("parse protoc error: %v", err))
			}
			err = protoc.generate()
		}

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
