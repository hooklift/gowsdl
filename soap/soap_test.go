package soap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Ping struct {
	XMLName xml.Name `xml:"http://example.com/service.xsd Ping"`

	Request *PingRequest `xml:"request,omitempty"`
}

type PingRequest struct {
	// XMLName xml.Name `xml:"http://example.com/service.xsd PingRequest"`

	Message    string  `xml:"Message,omitempty"`
	Attachment *Binary `xml:"Attachment,omitempty"`
}

type PingResponse struct {
	XMLName xml.Name `xml:"http://example.com/service.xsd PingResponse"`

	PingResult *PingReply `xml:"PingResult,omitempty"`
}

type PingReply struct {
	// XMLName xml.Name `xml:"http://example.com/service.xsd PingReply"`

	Message    string `xml:"Message,omitempty"`
	Attachment []byte `xml:"Attachment,omitempty"`
}

func TestClient_Call(t *testing.T) {
	var pingRequest = new(Ping)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xml.NewDecoder(r.Body).Decode(pingRequest)
		rsp := `<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
			<soap:Body>
				<PingResponse xmlns="http://example.com/service.xsd">
					<PingResult>
						<Message>Pong hi</Message>
					</PingResult>
				</PingResponse>
			</soap:Body>
		</soap:Envelope>`
		w.Write([]byte(rsp))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	req := &Ping{Request: &PingRequest{Message: "Hi"}}
	reply := &PingResponse{}
	if err := client.Call("GetData", req, reply); err != nil {
		t.Fatalf("couln't call service: %v", err)
	}

	wantedMsg := "Pong hi"
	if reply.PingResult.Message != wantedMsg {
		t.Errorf("got msg %s wanted %s", reply.PingResult.Message, wantedMsg)
	}
}

func TestClient_Send_Correct_Headers(t *testing.T) {
	tests := []struct {
		action          string
		reqHeaders      map[string]string
		expectedHeaders map[string]string
	}{
		// default case when no custom headers are set
		{
			"GetTrade",
			map[string]string{},
			map[string]string{
				"User-Agent":   "gowsdl/0.1",
				"SOAPAction":   "GetTrade",
				"Content-Type": "text/xml; charset=\"utf-8\"",
			},
		},
		// override default User-Agent
		{
			"SaveTrade",
			map[string]string{"User-Agent": "soap/0.1"},
			map[string]string{
				"User-Agent": "soap/0.1",
				"SOAPAction": "SaveTrade",
			},
		},
		// override default Content-Type
		{
			"SaveTrade",
			map[string]string{"Content-Type": "text/xml; charset=\"utf-16\""},
			map[string]string{"Content-Type": "text/xml; charset=\"utf-16\""},
		},
	}

	var gotHeaders http.Header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header
	}))
	defer ts.Close()

	for _, test := range tests {
		client := NewClient(ts.URL, WithHTTPHeaders(test.reqHeaders))
		req := struct{}{}
		reply := struct{}{}
		client.Call(test.action, req, reply)

		for k, v := range test.expectedHeaders {
			h := gotHeaders.Get(k)
			if h != v {
				t.Errorf("got %s wanted %s", h, v)
			}
		}
	}
}

func TestClient_MTOM(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			w.Header().Set(k, v[0])
		}
		bodyBuf, _ := ioutil.ReadAll(r.Body)
		w.Write(bodyBuf)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, WithMTOM())
	req := &PingRequest{Attachment: NewBinary([]byte("Attached data")).SetContentType("text/plain")}
	reply := &PingRequest{}
	if err := client.Call("GetData", req, reply); err != nil {
		t.Fatalf("couln't call service: %v", err)
	}

	if !bytes.Equal(reply.Attachment.Bytes(), req.Attachment.Bytes()) {
		t.Errorf("got %s wanted %s", reply.Attachment.Bytes(), req.Attachment.Bytes())
	}

	if reply.Attachment.ContentType() != req.Attachment.ContentType() {
		t.Errorf("got %s wanted %s", reply.Attachment.Bytes(), req.Attachment.ContentType())
	}
}

type SimpleDetail struct {
	Node interface{} `xml:"SimpleNode"`
}

type SimpleNode struct {
	Detail string      `xml:"Detail,omitempty"`
	Num    float64     `xml:"Num,omitempty"`
	Nested *SimpleNode `xml:"Nested,omitempty"`
}

func Test_CustomNestedFaul(t *testing.T) {
	input := `<SimpleNode>
  <Name>SimpleNode</Name>
  <Detail>detail message</Detail>
  <Num>6.005</Num>
</SimpleNode>`
	decoder := xml.NewDecoder(strings.NewReader(input))
	var nested interface{}
	nested = &SimpleNode{}
	if err := decoder.Decode(&nested); err != nil {
		t.Fatalf("error decoding: %v", err)
	}
	assert.EqualValues(t, &SimpleNode{
		Detail: "detail message",
		Num:    6.005,
	}, nested)
}

func Test_Client_FaultDefault(t *testing.T) {
	tests := []struct {
		name       string
		fault      interface{}
		emptyFault interface{}
	}{
		{
			name: "Empty",
		},
		{
			name: "SimpleNode",
			fault: &SimpleNode{
				Detail: "detail message",
				Num:    7.7,
			},
			emptyFault: &SimpleNode{},
		},
		{
			name: "ArrayOfNode",
			fault: &[]SimpleNode{
				{
					Detail: "detail message-1",
					Num:    7.7,
				}, {
					Detail: "detail message-2",
					Num:    7.8,
				},
			},
			emptyFault: &[]SimpleNode{},
		},
		{
			name: "NestedNode",
			fault: &SimpleNode{
				Detail: "detail-1",
				Num:    .003,
				Nested: &SimpleNode{
					Detail: "nested-2",
					Num:    .004,
				},
			},
			emptyFault: &SimpleNode{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.MarshalIndent(tt.fault, "\t\t\t\t", "\t")
			if err != nil {
				t.Fatalf("Failed to encode input as XML: %v", err)
			}

			var pingRequest = new(Ping)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				xml.NewDecoder(r.Body).Decode(pingRequest)
				rsp := fmt.Sprintf(`
<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
	<soap:Body>
		<soap:Fault>
			<faultcode>soap:Server</faultcode>
			<faultstring>Custom error message.</faultstring>
			<detail>
%v
			</detail>
		</soap:Fault>
	</soap:Body>
</soap:Envelope>`, string(data))
				w.Write([]byte(rsp))
			}))
			defer ts.Close()

			faultErrString := "soap:Server: Custom error message."

			client := NewClient(ts.URL)
			req := &Ping{Request: &PingRequest{Message: "Hi"}}
			var reply PingResponse
			fault := SimpleDetail{Node: tt.emptyFault}
			if err := client.CallWithFault("GetData", req, &reply, &fault); err != nil {
				assert.EqualError(t, err, faultErrString)
				assert.EqualValues(t, tt.fault, fault.Node)
			} else {
				t.Fatalf("call to ping() should have failed, but succeeded.")
			}
		})
	}
}
