# WSDL to Go
Generates Go code from a WSDL file. This project is originally intended to generate Go clients for WS-* services.

### Features
* Supports only Document/Literal wrapped services, which are [WS-I](http://ws-i.org/) compliant
* Attempts to generate idiomatic Go code as much as possible
* Generates Go code in parallel: types, operations and soap proxy
* Supports: 
	* WSDL 1.1
	* XML Schema 1.0
	* SOAP 1.1
* Resolves external XML Schemas recursively, up to 5 recursions.
* Supports providing WSDL HTTP URL as well as a local WSDL file

### Usage
```
gowsdl [OPTIONS]

Application Options:
  -v, --version     Shows gowsdl version
  -p, --package=    Package under which code will be generated (myservice)
  -o, --output=     File where the generated code will be saved (myservice.go)
  -i, --ignore-tls  Ignores invalid TLS certificates. It is not recomended for production. Use at your own risk
                    (false)

Help Options:
  -h, --help        Show this help message
```

### TODO
* If WSDL file is local, resolve external XML schemas locally too instead of failing due to not having a URL to download them from.
* Resolve XSD element references
* Support for generating namespaces
* Make code generation agnostic so generating code to other programming languages is feasible through plugins


## License
Copyright 2014 Cloudescape. All rights reserved.

