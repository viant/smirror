package main

import (
	_ "github.com/viant/afsc/gs"
	_ "github.com/viant/afsc/s3"
	"os"
	"smirror/cmd"
)

var Version string

func main() {
	cmd.RunClient(Version, os.Args[1:])
}

