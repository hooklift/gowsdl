package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type GoWsdl struct {
	file, pkg string
	logger    *log.Logger
	wsdl      *Wsdl
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

func (g *GoWsdl) Start() (map[string][]byte, error) {
	gocode := make(map[string][]byte)

	err := g.Unmarshal()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		gocode["types"], err = g.GenTypes()
		if err != nil {
			g.logger.Println(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		gocode["operations"], err = g.GenOperations()
		if err != nil {
			g.logger.Println(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		gocode["proxy"], err = g.GenSoapProxy()
		if err != nil {
			g.logger.Println(err)
		}
	}()

	wg.Wait()

	return gocode, nil
}

func (g *GoWsdl) Unmarshal() error {
	g.logger.Printf("Using %s...\n", g.file)

	//g.file is URL or local file?
	//if URL, download!

	data, err := ioutil.ReadFile(g.file)
	if err != nil {
		return err
	}

	g.wsdl = &Wsdl{}
	err = xml.Unmarshal(data, g.wsdl)
	if err != nil {
		return err
	}

	//Resolve wsdl imports
	//Resolve xsd includes
	//Resolve xsd imports

	return nil
}

func (g *GoWsdl) GenTypes() ([]byte, error) {
	if g.wsdl == nil {
		g.logger.Fatalln("You have to unmarshal the WSDL file first")
	}

	//element > complexType

	for _, schema := range g.wsdl.Types.Schemas {
		g.logger.Println(schema.XMLName)
		// for _, element := range schema.Elements {
		// 	g.logger.Printf("Type: %s\n", element.Name)
		// }
		g.logger.Printf("Total types: %d\n", len(schema.Elements))
	}

	return nil, nil
}

func (g *GoWsdl) GenOperations() ([]byte, error) {
	if g.wsdl == nil {
		g.logger.Fatalln("You have to unmarshal the WSDL file first")
	}

	for _, pt := range g.wsdl.PortTypes {
		// for _, o := range pt.Operations {
		// 	g.logger.Printf("Operation: %s", o.Name)
		// }
		g.logger.Printf("Total ops: %d\n", len(pt.Operations))
	}

	return nil, nil
}

func (g *GoWsdl) GenSoapProxy() ([]byte, error) {
	if g.wsdl == nil {
		g.logger.Fatalln("You have to unmarshal the WSDL file first")
	}
	return nil, nil
}
