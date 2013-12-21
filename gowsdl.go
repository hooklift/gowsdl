package main

import (
	"log"
	"os"
	"strings"
)

type GoWsdl struct {
	file, pkg, output string
	logger            *log.Logger
}

func NewGoWsdl(file, pkg, output string, logger *log.Logger) (*GoWsdl, error) {
	file = strings.TrimSpace(file)
	if file == "" {
		logger.Fatalln("WSDL file is required to generate Go proxy")
	}

	pkg = strings.TrimSpace(pkg)
	if pkg == "" {
		pkg = "main"
	}

	output = strings.TrimSpace(output)
	if output == "" {
		output = "."
	}

	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}

	log.Println(file)

	return &GoWsdl{
		file:   file,
		pkg:    pkg,
		output: output,
		logger: logger,
	}, nil
}

func (g *GoWsdl) start() error {
	return nil
}
