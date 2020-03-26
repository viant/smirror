package cmd

import (
	"context"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"smirror/cmd/build"
	"smirror/cmd/mirror"
	"smirror/cmd/option"
	"smirror/cmd/validate"
	"smirror/shared"
)

//RunClient run client
func RunClient(Version string, args []string) {
	options := &option.Options{}
	_, err := flags.ParseArgs(options, args)
	if isHelOption(args) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	if options.Version {
		shared.LogF("SMirror: Version: %v\n", Version)
		return
	}
	canBuildRule :=  options.DestinationURL != ""
	canMirror := options.SourceURL != ""
	if !(canMirror || options.Validate || canBuildRule) && len(args) == 1 {
		os.Exit(1)
	}

	srv, err := New(options.ProjectID)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if options.RuleURL == "" || canBuildRule {
		err = srv.Build(ctx, &build.Request{Options: options})
		if err != nil {
			log.Fatal(err)
		}
	}
	if options.Validate {
		err = srv.Validate(ctx, &validate.Request{Options: options})
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	response, err := srv.Mirror(ctx, &mirror.Request{options})
	if err != nil {
		log.Fatal(err)
	}
	shared.LogLn(response)
	if len(response.Errors) > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}


func isHelOption(args []string) bool {
	for _, arg := range args {
		if arg == "-h" {
			return true
		}
	}
	return false
}

