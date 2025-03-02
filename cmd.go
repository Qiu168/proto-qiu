package main

import (
	"flag"
	"fmt"
	"os"
)

type Cmd struct {
	javaOutput string
	version    bool
	protocPath []string
}

func parseCmd() *Cmd {
	cmd := &Cmd{}
	flag.Usage = printUsage
	flag.BoolVar(&cmd.version, "version", false, "show version")
	flag.StringVar(&cmd.javaOutput, "java_out", "\\", "java output file")
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		cmd.protocPath = args[1:]
	}
	return cmd
}
func printUsage() {
	fmt.Printf("Usage: %s -java_out=[Path] [args...]\n", os.Args[0])
}
