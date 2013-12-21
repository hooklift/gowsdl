package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type GoWsdl struct {
	file, pkg string
	logger    *log.Logger
}

func NewGoWsdl(file, pkg string, logger *log.Logger) (*GoWsdl, error) {
	file = strings.TrimSpace(file)
	if file == "" {
		logger.Fatalln("WSDL file is required to generate Go proxy")
	}

	pkg = strings.TrimSpace(pkg)
	if pkg == "" {
		pkg = "main"
	}

	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}

	return &GoWsdl{
		file:   file,
		pkg:    pkg,
		logger: logger,
	}, nil
}

func (g *GoWsdl) Unmarshal() error {
	g.logger.Printf("Using %s...\n", g.file)

	//URL or local file?
	//if URL, download!

	data, err := ioutil.ReadFile(g.file)
	if err != nil {
		return err
	}

	wsdl := Wsdl{}
	err = xml.Unmarshal(data, &wsdl)
	if err != nil {
		return err
	}

	return nil
}

func (g *GoWsdl) GenTypes() ([]byte, error) {
	return nil, nil
}

func (g *GoWsdl) GenOperations() ([]byte, error) {
	return nil, nil
}

func (g *GoWsdl) GenSoapProxy() ([]byte, error) {
	return nil, nil
}
