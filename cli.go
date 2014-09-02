package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"runtime"

	gen "github.com/cloudscape/gowsdl/generator"
	flags "github.com/jessevdk/go-flags"
)

const version = "v0.0.1"

var opts struct {
	Version    bool   `short:"v" long:"version" description:"Shows gowsdl version"`
	Package    string `short:"p" long:"package" description:"Package under which code will be generated" default:"myservice"`
	OutputFile string `short:"o" long:"output" description:"File where the generated code will be saved" default:"myservice.go"`
	IgnoreTls  bool   `short:"i" long:"ignore-tls" description:"Ignores invalid TLS certificates. It is not recomended for production. Use at your own risk" default:"false"`
}

func init() {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.SetPrefix("🍀  ")
}

func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		log.Println(version)
		os.Exit(0)
	}

	if len(args) == 0 {
		log.Fatalln("WSDL file is required to start the party")
	}

	if opts.OutputFile == args[0] {
		log.Fatalln("Output file cannot be the same WSDL file")
	}

	gowsdl, err := gen.NewGoWsdl(args[0], opts.Package, opts.IgnoreTls)
	if err != nil {
		log.Fatalln(err)
	}

	gocode, err := gowsdl.Start()
	if err != nil {
		log.Fatalln(err)
	}

	pkg := "./" + opts.Package
	err = os.Mkdir(pkg, 0744)

	if perr, ok := err.(*os.PathError); ok && os.IsExist(perr.Err) {
		log.Printf("Package directory %s already exist, skipping creation\n", pkg)
	} else {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fd, err := os.Create(pkg + "/" + opts.OutputFile)
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
