package gen

import (
	"encoding/xml"
	"time"

	"github.com/hooklift/gowsdl/soap"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type TradePriceRequest struct {
	XMLName xml.Name `xml:"http://example.com/stockquote.xsd TradePriceRequest"`

	TickerSymbol string `xml:"tickerSymbol,omitempty"`
}

type TradePrice struct {
	XMLName xml.Name `xml:"http://example.com/stockquote.xsd TradePrice"`

	Price float32 `xml:"price,omitempty"`
}

type StockQuotePortType interface {
	GetLastTradePrice(request *TradePriceRequest) (*TradePrice, error)
}

type stockQuotePortType struct {
	client *soap.Client
}

func NewStockQuotePortType(client *soap.Client) StockQuotePortType {
	return &stockQuotePortType{
		client: client,
	}
}

func (service *stockQuotePortType) GetLastTradePrice(request *TradePriceRequest) (*TradePrice, error) {
	response := new(TradePrice)
	err := service.client.Call("http://example.com/GetLastTradePrice", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
