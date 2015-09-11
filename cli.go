// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"flag"
	"go/format"
	"log"
	"os"

	gen "github.com/sanbornm/gowsdl/generator"
)

const version = "v0.0.1"

var vers = flag.Bool("v", false, "Shows gowsdl version")
var pkg = flag.String("p", "myservice", "Package under which code will be generated")
var outFile = flag.String("o", "myservice.go", "File where the generated code will be saved")

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.SetPrefix("🍀  ")
}

func main() {
	flag.Parse()

	if *vers {
		log.Println(version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		log.Fatalln("WSDL file is required to start the party")
	}

	if *outFile == os.Args[1] {
		log.Fatalln("Output file cannot be the same WSDL file")
	}

	gowsdl, err := gen.NewGoWsdl(os.Args[1], *pkg, false)
	if err != nil {
		log.Fatalln(err)
	}

	gocode, err := gowsdl.Start()
	if err != nil {
		log.Fatalln(err)
	}

	pkg := "./" + *pkg
	err = os.Mkdir(pkg, 0744)

	fd, err := os.Create(pkg + "/" + *outFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer fd.Close()

	data := new(bytes.Buffer)
	data.Write(gocode["header"])
	data.Write(gocode["types"])
	data.Write(gocode["operations"])

	source, err := format.Source(data.Bytes())
	if err != nil {
		fd.Write(data.Bytes())
		log.Fatalln(err)
	}

	fd.Write(source)

	log.Println("Done 💩")
}
