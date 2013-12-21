package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"log"
	"os"
)

const version = "v0.0.1"

var opts struct {
	//Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Version bool   `short:"v" long:"version" description:"Shows gowsdl version"`
	Package string `short:"p" long:"package" description:"Package under which code will be generated" default:"main"`
	Output  string `short:"o" long:"output" description:"Directory where the code will be generated" default:"."`
}

func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	if len(args) == 0 {
		fmt.Println("WSDL file is required to start the party")
		os.Exit(1)
	}

	logger := log.New(os.Stdout, "", 0)

	gowsdl, err := NewGoWsdl(args[0], opts.Package, opts.Output, logger)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gowsdl.start()
}
