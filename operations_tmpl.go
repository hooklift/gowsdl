// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var opsTmpl = `
{{range .}}
	{{$privateType := .Name | makePrivate}}
	{{$exportType := .Name | makePublic}}

	type {{$exportType}} interface {
		{{if eq $exportType "NetSuitePortType"}}
			//---
			Customer_Add(request *Customer_AddRequest) (*Customer_AddResponse, error)
			Customer_AddContext(ctx context.Context, request *Customer_AddRequest) (*Customer_AddResponse, error)
	
			Customer_Get(request *Customer_GetRequest) (*Customer_GetResponse, error)
			Customer_GetContext(ctx context.Context, request *Customer_GetRequest) (*Customer_GetResponse, error)
			
			//---
			CreditMemo_Add(request *CreditMemo_AddRequest) (*CreditMemo_AddResponse, error)
			CreditMemo_AddContext(ctx context.Context, request *CreditMemo_AddRequest) (*CreditMemo_AddResponse, error)
		
			CreditMemo_Get(request *CreditMemo_GetRequest) (*CreditMemo_GetResponse, error)
			CreditMemo_GetContext(ctx context.Context, request *CreditMemo_GetRequest) (*CreditMemo_GetResponse, error)
		
			//---
			CashSale_Add(request *CashSale_AddRequest) (*CashSale_AddResponse, error)
			CashSale_AddContext(ctx context.Context, request *CashSale_AddRequest) (*CashSale_AddResponse, error)
		
			CashSale_Get(request *CashSale_GetRequest) (*CashSale_GetResponse, error)
			CashSale_GetContext(ctx context.Context, request *CashSale_GetRequest) (*CashSale_GetResponse, error)
		
			//---
			Invoice_Get(request *Invoice_GetRequest) (*Invoice_GetResponse, error)
			Invoice_GetContext(ctx context.Context, request *Invoice_GetRequest) (*Invoice_GetResponse, error)
		
			Invoice_Add(request *Invoice_AddRequest) (*Invoice_AddResponse, error)
			Invoice_AddContext(ctx context.Context, request *Invoice_AddRequest) (*Invoice_AddResponse, error)
		
			//---
			SalesOrder_Get(request *SalesOrder_GetRequest) (*SalesOrder_GetResponse, error)
			SalesOrder_GetContext(ctx context.Context, request *SalesOrder_GetRequest) (*SalesOrder_GetResponse, error)
		
			SalesOrder_Add(request *SalesOrder_AddRequest) (*SalesOrder_AddResponse, error)
			SalesOrder_AddContext(ctx context.Context, request *SalesOrder_AddRequest) (*SalesOrder_AddResponse, error)

			//---
			CustomerDeposit_Get(request *CustomerDeposit_GetRequest) (*CustomerDeposit_GetResponse, error)
			CustomerDeposit_GetContext(ctx context.Context, request *CustomerDeposit_GetRequest) (*CustomerDeposit_GetResponse, error)
		
			CustomerDeposit_Add(request *CustomerDeposit_AddRequest) (*CustomerDeposit_AddResponse, error)
			CustomerDeposit_AddContext(ctx context.Context, request *CustomerDeposit_AddRequest) (*CustomerDeposit_AddResponse, error)

		{{end}}
		{{range .Operations}}
			{{$faults := len .Faults}}
			{{$soapAction := findSOAPAction .Name $privateType}}
			{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}}
			{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}

			{{/*if ne $soapAction ""*/}}
			{{if gt $faults 0}}
			// Error can be either of the following types:
			// {{range .Faults}}
			//   - {{.Name}} {{.Doc}}{{end}}{{end}}
			{{if ne .Doc ""}}/* {{.Doc}} */{{end}}
			{{makePublic .Name | replaceReservedWords}} ({{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error)
			{{/*end*/}}
			{{makePublic .Name | replaceReservedWords}}Context (ctx context.Context, {{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error)
			{{/*end*/}}
		{{end}}
	}

	type {{$privateType}} struct {
		client *soap.Client
	}

	func New{{$exportType}}(client *soap.Client) {{$exportType}} {
		/* Here */
		return &{{$privateType}}{
			client: client,
		}
	}
	
	{{range .Operations}}
		{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}}
		{{$soapAction := findSOAPAction .Name $privateType}}
		{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}
		func (service *{{$privateType}}) {{makePublic .Name | replaceReservedWords}}Context (ctx context.Context, {{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error) {
			{{if ne $responseType ""}}response := new({{$responseType}}){{end}}
			err := service.client.CallContext(ctx, "{{if ne $soapAction ""}}{{$soapAction}}{{else}}''{{end}}", {{if ne $requestType ""}}request{{else}}nil{{end}}, {{if ne $responseType ""}}response{{else}}struct{}{}{{end}})
			if err != nil {
				return {{if ne $responseType ""}}nil, {{end}}err
			}

			return {{if ne $responseType ""}}response, {{end}}nil
		}

		func (service *{{$privateType}}) {{makePublic .Name | replaceReservedWords}} ({{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error) {
			return service.{{makePublic .Name | replaceReservedWords}}Context(
				context.Background(),
				{{if ne $requestType ""}}request,{{end}}
			)
		}

	{{end}}
{{end}}
`
