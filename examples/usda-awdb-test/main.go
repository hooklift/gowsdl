package main

import "fmt"

func main() {
	service := NewAwdbWebService("http://www.wcc.nrcs.usda.gov/awdbWebService/services", false)
	amIthere, err := service.AreYouThere(&areYouThere{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Alive?: %t\n", amIthere.Return_)

	stations, err := service.GetStations(&getStations{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Stations: %v\n", stations.Return_)
}
