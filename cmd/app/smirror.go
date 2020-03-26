package main



import (
	_ "github.com/viant/afsc/s3"
	_ "github.com/viant/afsc/gs"
	"os"
	"smirror/cmd"
)

var Version string

func main() {
	cmd.RunClient(Version, os.Args[1:])
}

