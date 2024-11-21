package example

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/hooklift/gowsdl/example/gen"
	"github.com/hooklift/gowsdl/soap"
)

//go:generate go run ../cmd/gowsdl -p gen -o gen.go ../fixtures/stock.wsdl

func ExampleBasicUsage() {
	client := soap.NewClient("http://svc.asmx")
	service := gen.NewStockQuotePortType(client)
	reply, err := service.GetLastTradePrice(&gen.TradePriceRequest{})
	if err != nil {
		log.Fatalf("could't get trade prices: %v", err)
	}
	log.Println(reply)
}

func ExampleWithOptions() {
	client := soap.NewClient(
		"http://svc.asmx",
		soap.WithTimeout(time.Second*5),
		soap.WithBasicAuth("usr", "psw"),
		soap.WithTLS(&tls.Config{InsecureSkipVerify: true}),
	)
	service := gen.NewStockQuotePortType(client)
	reply, err := service.GetLastTradePrice(&gen.TradePriceRequest{})
	if err != nil {
		log.Fatalf("could't get trade prices: %v", err)
	}
	log.Println(reply)
}
