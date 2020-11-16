// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
/*

Gowsdl generates Go code from a WSDL file.

This project is originally intended to generate Go clients for WS-* services.

Usage: gowsdl [options] myservice.wsdl
  -o string
        File where the generated code will be saved (default "myservice.go")
  -p string
        Package under which code will be generated (default "myservice")
  -v    Shows gowsdl version

Features

Supports only Document/Literal wrapped services, which are WS-I (http://ws-i.org/) compliant.

Attempts to generate idiomatic Go code as much as possible.

Supports WSDL 1.1, XML Schema 1.0, SOAP 1.1.

Resolves external XML Schemas

Supports providing WSDL HTTP URL as well as a local WSDL file.

Not supported

UDDI.

TODO

Add support for filters to allow the user to change the generated code.

If WSDL file is local, resolve external XML schemas locally too instead of failing due to not having a URL to download them from.

Resolve XSD element references.

Support for generating namespaces.

Make code generation agnostic so generating code to other programming languages is feasible through plugins.

*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"

	gen "github.com/hooklift/gowsdl"
)

// Version is initialized in compilation time by go build.
var Version string

// Name is initialized in compilation time by go build.
var Name string

var vers = flag.Bool("v", false, "Shows gowsdl version")
var pkg = flag.String("p", "myservice", "Package under which code will be generated")
var outFile = flag.String("o", "myservice.go", "File where the generated code will be saved")
var dir = flag.String("d", "./", "Directory under which package directory will be created")
var insecure = flag.Bool("i", false, "Skips TLS Verification")
var makePublic = flag.Bool("make-public", true, "Make the generated types public/exported")

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.SetPrefix("🍀  ")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] myservice.wsdl\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// Show app version
	if *vers {
		log.Println(Version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	wsdlPath := os.Args[len(os.Args)-1]

	if *outFile == wsdlPath {
		log.Fatalln("Output file cannot be the same WSDL file")
	}

	// load wsdl
	gowsdl, err := gen.NewGoWSDL(wsdlPath, *pkg, *insecure, *makePublic)
	if err != nil {
		log.Fatalln(err)
	}

	// generate code
	gocode, err := gowsdl.Start()
	if err != nil {
		log.Fatalln(err)
	}

	pkg := filepath.Join(*dir, *pkg)
	err = os.Mkdir(pkg, 0744)

	file, err := os.Create(filepath.Join(pkg, *outFile))
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	data := new(bytes.Buffer)
	data.Write(gocode["header"])
	data.Write(gocode["types"])
	data.Write(gocode["operations"])
	data.Write(gocode["soap"])

	// go fmt the generated code
	source, err := format.Source(data.Bytes())
	if err != nil {
		file.Write(data.Bytes())
		log.Fatalln(err)
	}

	file.Write(source)

	log.Println("Done 👍")
}
