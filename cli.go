package main

import (
	flags "github.com/jessevdk/go-flags"
	"log"
	"os"
)

const version = "v0.0.1"

var opts struct {
	Version    bool   `short:"v" long:"version" description:"Shows gowsdl version"`
	Package    string `short:"p" long:"package" description:"Package under which code will be generated" default:"myservice"`
	OutputFile string `short:"o" long:"output" description:"File where the generated code will be saved" default:"myservice.go"`
}

func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	logger := log.New(os.Stdout, "🍀  ", 0)

	if opts.Version {
		logger.Println(version)
		os.Exit(0)
	}

	if len(args) == 0 {
		logger.Fatalln("WSDL file is required to start the party")
	}

	if opts.OutputFile == args[0] {
		logger.Fatalln("Output file cannot be the same WSDL file")
	}

	gowsdl, err := NewGoWsdl(args[0], opts.Package, logger)
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	gocode, err := gowsdl.Start()
	if err != nil {
		logger.Fatalln(err)
	}

	pkg := "./" + opts.Package
	err = os.Mkdir(pkg, 0744)

	if perr, ok := err.(*os.PathError); ok && os.IsExist(perr.Err) {
		logger.Println("Package directory already exist, skipping creation")
	} else {
		if err != nil {
			logger.Fatalln(err)
		}
	}

	fd, err := os.Create(pkg + "/" + opts.OutputFile)
	if err != nil {
		logger.Fatalln(err)
	}
	defer fd.Close()

	fd.Write(gocode["types"])
	fd.Write(gocode["proxy"])
	fd.Write(gocode["operations"])

	logger.Println("Done 💩")
}
