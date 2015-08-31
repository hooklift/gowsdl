package main

import (
	"fmt"
	"os"

	gowsdl "github.com/hooklift/gowsdl/generator"
	"gopkg.in/inconshreveable/log15.v2"
)

func main() {
	gowsdl.Log.SetHandler(log15.StreamHandler(os.Stdout, log15.TerminalFormat()))

	service := NewAwdbWebService("http://www.wcc.nrcs.usda.gov/awdbWebService/services", false, nil)
	amIthere, err := service.AreYouThere(&AreYouThere{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Alive?: %t\n", amIthere.Return_)

	request := &GetStations{NetworkCds: []string{"SNTL"}, LogicalAnd: true}
	stations, err := service.GetStations(request)
	if err != nil {
		fmt.Printf("\n->%s\n", err.(*gowsdl.SoapFault).Faultstring)
		return
	}
	fmt.Printf("Stations: %+v\n", stations.Return_)
}
