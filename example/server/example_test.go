package main

import (
	"net/http/httptest"
	"testing"

	"github.com/hooklift/gowsdl/example/server/gen"
	"github.com/hooklift/gowsdl/soap"
)

func TestExampleServer(t *testing.T) {
	httpServer := httptest.NewServer(gen.CreateEndpoint(&server{}))
	t.Cleanup(httpServer.Close)

	httpClient := httpServer.Client()
	soapClient := soap.NewClient(httpServer.URL, soap.WithHTTPClient(httpClient))
	service := gen.NewMNBArfolyamServiceType(soapClient)

	id := "veryrandomid"
	resp, err := service.GetInfoSoap(&gen.GetInfo{
		Id: id,
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := "gowsdl, " + id
	if resp.GetInfoResult != expected {
		t.Fatalf("got %s, want %s", resp.GetInfoResult, expected)
	}
}
