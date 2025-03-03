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
		var protoPaths []string
		for _, path := range cmd.protocPath {
			if strings.HasSuffix(path, constant.ProtoFileSuffix) {
				protoPaths = append(protoPaths, path)
			} else {
				// 获取path目录下的所有文件
				paths := getDirFiles(path)
				for _, p := range paths {
					if strings.HasSuffix(p, constant.ProtoFileSuffix) {
						protoPaths = append(protoPaths, p)
					}
				}
			}
		}
		for _, path := range protoPaths {
			javaProto, err := java.NewJavaProtoc(cmd.javaOutput, path)
			if err != nil {
				panic(fmt.Errorf("parse protoc error: %v", err))
			}
			err = javaProto.Generate()
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
