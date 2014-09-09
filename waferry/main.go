package main

import "fmt"

func main() {
	service := NewWSF_x0020_ScheduleSoap("http://b2b.wsdot.wa.gov/ferries/schedule/Default.asmx", false)
	seasons, err := service.GetActiveScheduledSeasons(&GetActiveScheduledSeasons{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Seasons: %+v\n", seasons)

}
