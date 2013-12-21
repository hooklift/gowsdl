package main

import (
	flags "github.com/jessevdk/go-flags"
	"log"
	"os"
	"sync"
)

const version = "v0.0.1"

var opts struct {
	Version    bool   `short:"v" long:"version" description:"Shows gowsdl version"`
	Package    string `short:"p" long:"package" description:"Package under which code will be generated" default:"gowsdl"`
	OutputFile string `short:"o" long:"output" description:"File where the generated code will be saved" default:"./myservice.go"`
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

	err = gowsdl.Unmarshal()
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	var types, operations, proxy []byte
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		types, err = gowsdl.GenTypes()
		if err != nil {
			logger.Fatalln(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		operations, err = gowsdl.GenOperations()
		if err != nil {
			logger.Fatalln(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		proxy, err = gowsdl.GenSoapProxy()
		if err != nil {
			logger.Fatalln(err)
		}
	}()

	wg.Wait()

	fd, err := os.Create(opts.OutputFile)
	if err != nil {
		logger.Fatalln(err)
	}
	defer fd.Close()

	fd.Write(types)
	fd.Write(proxy)
	fd.Write(operations)
}
