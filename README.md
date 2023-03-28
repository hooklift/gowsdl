# WSDL to Go

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/hooklift/gowsdl?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![GoDoc](https://godoc.org/github.com/hooklift/gowsdl?status.svg)](https://godoc.org/github.com/hooklift/gowsdl)
[![Build Status](https://travis-ci.org/hooklift/gowsdl.svg?branch=master)](https://travis-ci.org/hooklift/gowsdl)

Generates Go code from a WSDL file.

### Install

* [Download release](https://github.com/hooklift/gowsdl/releases)
* Download and build locally
    * 1.15: `go get github.com/hooklift/gowsdl/...`
    * 1.20: `go install github.com/hooklift/gowsdl/cmd/gowsdl@latest`
* Install from Homebrew: `brew install gowsdl`

### Goals
* Generate idiomatic Go code as much as possible
* Support only Document/Literal wrapped services, which are [WS-I](http://ws-i.org/) compliant
* Support:
	* WSDL 1.1
	* XML Schema 1.0
	* SOAP 1.1
* Resolve external XML Schemas
* Support external and local WSDL

### Caveats
* Please keep in mind that the generated code is just a reflection of what the WSDL is like. If your WSDL has duplicated type definitions, your Go code is going to have the same and may not compile.

### Usage
```
Usage: gowsdl [options] myservice.wsdl
  -o string
        File where the generated code will be saved (default "myservice.go")
  -p string
        Package under which code will be generated (default "myservice")
  -i    Skips TLS Verification
  -v    Shows gowsdl version
  ```
