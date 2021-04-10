package soap

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

type AttachmentRequest struct {
	XMLName xml.Name `xml:"http://example.com/service.xsd attachmentRequest"`

	Name      string `xml:"name,omitempty"`
	ContentID string `xml:"contentID,omitempty"`
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

func TestClient_Attachments_WithAttachmentResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			w.Header().Set(k, v[0])
		}
		bodyBuf, _ := ioutil.ReadAll(r.Body)
		_, err := w.Write(bodyBuf)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	// GIVEN
	firstAtt := MIMEMultipartAttachment{
		Name: "First_Attachment",
		Data: []byte(`foobar`),
	}
	secondAtt := MIMEMultipartAttachment{
		Name: "Second_Attachment",
		Data: []byte(`tl;tr`),
	}
	client := NewClient(ts.URL, WithMIMEMultipartAttachments())
	client.AddMIMEMultipartAttachment(firstAtt)
	client.AddMIMEMultipartAttachment(secondAtt)
	req := &AttachmentRequest{
		Name:      "UploadMyFilePlease",
		ContentID: "First_Attachment",
	}
	reply := new(AttachmentRequest)
	retAttachments := make([]MIMEMultipartAttachment, 0)

	// WHEN
	if err := client.CallContextWithAttachmentsAndFaultDetail(context.TODO(), "''", req,
		reply, nil, &retAttachments); err != nil {
		t.Fatalf("couln't call service: %v", err)
	}

	// THEN
	assert.Equal(t, req.ContentID, reply.ContentID)
	assert.Len(t, retAttachments, 2)
	assert.Equal(t, retAttachments[0], firstAtt)
	assert.Equal(t, retAttachments[1], secondAtt)
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

type SimpleNode struct {
	Detail string      `xml:"Detail,omitempty"`
	Num    float64     `xml:"Num,omitempty"`
	Nested *SimpleNode `xml:"Nested,omitempty"`
}

func (s SimpleNode) ErrorString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%.2f: %s", s.Num, s.Detail))
	if s.Nested != nil {
		sb.WriteString("\n" + s.Nested.ErrorString())
	}
	return sb.String()
}

func (s SimpleNode) HasData() bool {
	return true
}

type Wrapper struct {
	Item    interface{} `xml:"SimpleNode"`
	hasData bool
}

func (w *Wrapper) HasData() bool {
	return w.hasData
}

func (w *Wrapper) ErrorString() string {
	switch w.Item.(type) {
	case FaultError:
		return w.Item.(FaultError).ErrorString()
	}
	return "default error"
}

func Test_SimpleNode(t *testing.T) {
	input := `<SimpleNode>
  <Name>SimpleNode</Name>
  <Detail>detail message</Detail>
  <Num>6.005</Num>
</SimpleNode>`
	decoder := xml.NewDecoder(strings.NewReader(input))
	var simple interface{}
	simple = &SimpleNode{}
	if err := decoder.Decode(&simple); err != nil {
		t.Fatalf("error decoding: %v", err)
	}
	assert.EqualValues(t, &SimpleNode{
		Detail: "detail message",
		Num:    6.005,
	}, simple)
}

func Test_Client_FaultDefault(t *testing.T) {
	tests := []struct {
		name          string
		hasData       bool
		wantErrString string
		fault         interface{}
		emptyFault    interface{}
	}{
		{
			name:          "Empty-WithFault",
			wantErrString: "default error",
			hasData:       true,
		},
		{
			name:          "Empty-NoFaultDetail",
			wantErrString: "Custom error message.",
			hasData:       false,
		},
		{
			name:          "SimpleNode",
			wantErrString: "7.70: detail message",
			hasData:       true,
			fault: &SimpleNode{
				Detail: "detail message",
				Num:    7.7,
			},
			emptyFault: &SimpleNode{},
		},
		{
			name:          "ArrayOfNode",
			wantErrString: "default error",
			hasData:       true,
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
			wantErrString: "0.00: detail-1\n" +
				"0.00: nested-2",
			hasData: true,
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

			faultErrString := tt.wantErrString

			client := NewClient(ts.URL)
			req := &Ping{Request: &PingRequest{Message: "Hi"}}
			var reply PingResponse
			fault := Wrapper{
				Item:    tt.emptyFault,
				hasData: tt.hasData,
			}
			if err := client.CallWithFaultDetail("GetData", req, &reply, &fault); err != nil {
				assert.EqualError(t, err, faultErrString)
				assert.EqualValues(t, tt.fault, fault.Item)
			} else {
				t.Fatalf("call to ping() should have failed, but succeeded.")
			}
		})
	}
}

// TestXsdDateTime checks the marshalled xsd datetime
func TestXsdDateTime(t *testing.T) {
	type TestDateTime struct {
		XMLName  xml.Name `xml:"TestDateTime"`
		Datetime XSDDateTime
	}
	// test marshalling
	{
		// without nanosecond
		testDateTime := TestDateTime{
			Datetime: CreateXsdDateTime(time.Date(1951, time.October, 22, 1, 2, 3, 0, time.FixedZone("UTC-8", -8*60*60)), true),
		}
		if output, err := xml.MarshalIndent(testDateTime, "", ""); err != nil {
			t.Error(err)
		} else {
			outputstr := string(output)
			expected := "<TestDateTime><Datetime>1951-10-22T01:02:03-08:00</Datetime></TestDateTime>"
			if outputstr != expected {
				t.Errorf("Got: %v\nExpected: %v", outputstr, expected)
			}
		}
	}
	{
		// with nanosecond
		testDateTime := TestDateTime{
			Datetime: CreateXsdDateTime(time.Date(1951, time.October, 22, 1, 2, 3, 4, time.FixedZone("UTC-8", -8*60*60)), true),
		}
		if output, err := xml.MarshalIndent(testDateTime, "", ""); err != nil {
			t.Error(err)
		} else {
			outputstr := string(output)
			expected := "<TestDateTime><Datetime>1951-10-22T01:02:03.000000004-08:00</Datetime></TestDateTime>"
			if outputstr != expected {
				t.Errorf("Got: %v\nExpected: %v", outputstr, expected)
			}
		}
	}

	// test marshalling of UTC
	{
		testDateTime := TestDateTime{
			Datetime: CreateXsdDateTime(time.Date(1951, time.October, 22, 1, 2, 3, 4, time.UTC), true),
		}
		if output, err := xml.MarshalIndent(testDateTime, "", ""); err != nil {
			t.Error(err)
		} else {
			outputstr := string(output)
			expected := "<TestDateTime><Datetime>1951-10-22T01:02:03.000000004Z</Datetime></TestDateTime>"
			if outputstr != expected {
				t.Errorf("Got:      %v\nExpected: %v", outputstr, expected)
			}
		}
	}

	// test marshalling of XsdDateTime without TZ
	{
		testDateTime := TestDateTime{
			Datetime: CreateXsdDateTime(time.Date(1951, time.October, 22, 1, 2, 3, 4, time.UTC), false),
		}
		if output, err := xml.MarshalIndent(testDateTime, "", ""); err != nil {
			t.Error(err)
		} else {
			outputstr := string(output)
			expected := "<TestDateTime><Datetime>1951-10-22T01:02:03.000000004</Datetime></TestDateTime>"
			if outputstr != expected {
				t.Errorf("Got:      %v\nExpected: %v", outputstr, expected)
			}
		}
	}

	// test unmarshalling
	{
		dateTimes := map[string]time.Time{
			"<TestDateTime><Datetime>1951-10-22T01:02:03.000000004-08:00</Datetime></TestDateTime>": time.Date(1951, time.October, 22, 1, 2, 3, 4, time.FixedZone("-0800", -8*60*60)),
			"<TestDateTime><Datetime>1951-10-22T01:02:03Z</Datetime></TestDateTime>":                time.Date(1951, time.October, 22, 1, 2, 3, 0, time.UTC),
			"<TestDateTime><Datetime>1951-10-22T01:02:03</Datetime></TestDateTime>":                 time.Date(1951, time.October, 22, 1, 2, 3, 0, time.Local),
		}
		for dateTimeStr, dateTimeObj := range dateTimes {
			parsedDt := TestDateTime{}
			if err := xml.Unmarshal([]byte(dateTimeStr), &parsedDt); err != nil {
				t.Error(err)
			} else {
				if !parsedDt.Datetime.ToGoTime().Equal(dateTimeObj) {
					t.Errorf("Got:      %#v\nExpected: %#v", parsedDt.Datetime.ToGoTime(), dateTimeObj)
				}
			}
		}
	}
}

// TestXsdDateTime checks the marshalled xsd datetime
func TestXsdDate(t *testing.T) {
	type TestDate struct {
		XMLName xml.Name `xml:"TestDate"`
		Date    XSDDate
	}

	// test marshalling
	{
		testDate := TestDate{
			Date: CreateXsdDate(time.Date(1951, time.October, 22, 0, 0, 0, 0, time.FixedZone("UTC-8", -8*60*60)), false),
		}
		if output, err := xml.MarshalIndent(testDate, "", ""); err != nil {
			t.Error(err)
		} else {
			outputstr := string(output)
			expected := "<TestDate><Date>1951-10-22</Date></TestDate>"
			if outputstr != expected {
				t.Errorf("Got: %v\nExpected: %v", outputstr, expected)
			}
		}
	}

	// test marshalling
	{
		testDate := TestDate{
			Date: CreateXsdDate(time.Date(1951, time.October, 22, 0, 0, 0, 0, time.FixedZone("UTC-8", -8*60*60)), true),
		}
		if output, err := xml.MarshalIndent(testDate, "", ""); err != nil {
			t.Error(err)
		} else {
			outputstr := string(output)
			expected := "<TestDate><Date>1951-10-22-08:00</Date></TestDate>"
			if outputstr != expected {
				t.Errorf("Got: %v\nExpected: %v", outputstr, expected)
			}
		}
	}

	// test marshalling of UTC
	{
		testDate := TestDate{
			Date: CreateXsdDate(time.Date(1951, time.October, 22, 0, 0, 0, 0, time.UTC), true),
		}
		if output, err := xml.MarshalIndent(testDate, "", ""); err != nil {
			t.Error(err)
		} else {
			outputstr := string(output)
			expected := "<TestDate><Date>1951-10-22Z</Date></TestDate>"
			if outputstr != expected {
				t.Errorf("Got:      %v\nExpected: %v", outputstr, expected)
			}
		}
	}

	// test unmarshalling
	{
		dates := map[string]time.Time{
			"<TestDate><Date>1951-10-22</Date></TestDate>":       time.Date(1951, time.October, 22, 0, 0, 0, 0, time.Local),
			"<TestDate><Date>1951-10-22Z</Date></TestDate>":      time.Date(1951, time.October, 22, 0, 0, 0, 0, time.UTC),
			"<TestDate><Date>1951-10-22-08:00</Date></TestDate>": time.Date(1951, time.October, 22, 0, 0, 0, 0, time.FixedZone("UTC-8", -8*60*60)),
		}
		for dateStr, dateObj := range dates {
			parsedDate := TestDate{}
			if err := xml.Unmarshal([]byte(dateStr), &parsedDate); err != nil {
				t.Error(dateStr, err)
			} else {
				if !parsedDate.Date.ToGoTime().Equal(dateObj) {
					t.Errorf("Got:      %#v\nExpected: %#v", parsedDate.Date.ToGoTime(), dateObj)
				}
			}
		}
	}
}

// TestXsdTime checks the marshalled xsd datetime
func TestXsdTime(t *testing.T) {
	type TestTime struct {
		XMLName xml.Name `xml:"TestTime"`
		Time    XSDTime
	}

	// test marshalling
	{
		testTime := TestTime{
			Time: CreateXsdTime(12, 13, 14, 4, time.FixedZone("Test", -19800)),
		}
		if output, err := xml.MarshalIndent(testTime, "", ""); err != nil {
			t.Error(err)
		} else {
			outputstr := string(output)
			expected := "<TestTime><Time>12:13:14.000000004-05:30</Time></TestTime>"
			if outputstr != expected {
				t.Errorf("Got: %v\nExpected: %v", outputstr, expected)
			}
		}
	}
	{
		testTime := TestTime{
			Time: CreateXsdTime(12, 13, 14, 0, time.FixedZone("UTC-8", -8*60*60)),
		}
		if output, err := xml.MarshalIndent(testTime, "", ""); err != nil {
			t.Error(err)
		} else {
			outputstr := string(output)
			expected := "<TestTime><Time>12:13:14-08:00</Time></TestTime>"
			if outputstr != expected {
				t.Errorf("Got: %v\nExpected: %v", outputstr, expected)
			}
		}
	}
	{
		testTime := TestTime{
			Time: CreateXsdTime(12, 13, 14, 0, nil),
		}
		if output, err := xml.MarshalIndent(testTime, "", ""); err != nil {
			t.Error(err)
		} else {
			outputstr := string(output)
			expected := "<TestTime><Time>12:13:14</Time></TestTime>"
			if outputstr != expected {
				t.Errorf("Got: %v\nExpected: %v", outputstr, expected)
			}
		}
	}

	// test unmarshalling without TZ
	{
		timeStr := "<TestTime><Time>12:13:14.000000004</Time></TestTime>"
		parsedTime := TestTime{}
		if err := xml.Unmarshal([]byte(timeStr), &parsedTime); err != nil {
			t.Error(err)
		} else {
			if parsedTime.Time.Hour() != 12 {
				t.Errorf("Got hour %#v\nExpected: %#v", parsedTime.Time.Hour(), 12)
			}
			if parsedTime.Time.Minute() != 13 {
				t.Errorf("Got minute %#v\nExpected: %#v", parsedTime.Time.Minute(), 13)
			}
			if parsedTime.Time.Second() != 14 {
				t.Errorf("Got second %#v\nExpected: %#v", parsedTime.Time.Second(), 14)
			}
			if parsedTime.Time.Nanosecond() != 4 {
				t.Errorf("Got nsec %#v\nExpected: %#v", parsedTime.Time.Nanosecond(), 4)
			}
			if parsedTime.Time.Location() != nil {
				t.Errorf("Got location %v\nExpected: Nil/Undetermined", parsedTime.Time.Location().String())
			}
		}
	}
	// test unmarshalling with UTC
	{
		timeStr := "<TestTime><Time>12:13:14Z</Time></TestTime>"
		parsedTime := TestTime{}
		if err := xml.Unmarshal([]byte(timeStr), &parsedTime); err != nil {
			t.Error(err)
		} else {
			if parsedTime.Time.Hour() != 12 {
				t.Errorf("Got hour %#v\nExpected: %#v", parsedTime.Time.Hour(), 12)
			}
			if parsedTime.Time.Minute() != 13 {
				t.Errorf("Got minute %#v\nExpected: %#v", parsedTime.Time.Minute(), 13)
			}
			if parsedTime.Time.Second() != 14 {
				t.Errorf("Got second %#v\nExpected: %#v", parsedTime.Time.Second(), 14)
			}
			if parsedTime.Time.Nanosecond() != 0 {
				t.Errorf("Got nsec %#v\nExpected: %#v", parsedTime.Time.Nanosecond(), 0)
			}
			if parsedTime.Time.Location().String() != "UTC" {
				t.Errorf("Got location %v\nExpected: UTC", parsedTime.Time.Location().String())
			}
		}
	}
	// test unmarshalling with non-UTC Tz
	{
		timeStr := "<TestTime><Time>12:13:14-08:00</Time></TestTime>"
		parsedTime := TestTime{}
		if err := xml.Unmarshal([]byte(timeStr), &parsedTime); err != nil {
			t.Error(err)
		} else {
			if parsedTime.Time.Hour() != 12 {
				t.Errorf("Got hour %#v\nExpected: %#v", parsedTime.Time.Hour(), 12)
			}
			if parsedTime.Time.Minute() != 13 {
				t.Errorf("Got minute %#v\nExpected: %#v", parsedTime.Time.Minute(), 13)
			}
			if parsedTime.Time.Second() != 14 {
				t.Errorf("Got second %#v\nExpected: %#v", parsedTime.Time.Second(), 14)
			}
			if parsedTime.Time.Nanosecond() != 0 {
				t.Errorf("Got nsec %#v\nExpected: %#v", parsedTime.Time.Nanosecond(), 0)
			}
			_, tzOffset := parsedTime.Time.innerTime.Zone()
			if tzOffset != -8*3600 {
				t.Errorf("Got location offset %v\nExpected: %v", tzOffset, -8*3600)
			}
		}
	}
}

func TestHTTPError(t *testing.T) {
	type httpErrorTest struct {
		name         string
		responseCode int
		responseBody string
		wantErr      bool
		wantErrMsg   string
	}

	tests := []httpErrorTest{
		{
			name:         "should error if server returns 500",
			responseCode: http.StatusInternalServerError,
			responseBody: "internal server error",
			wantErr:      true,
			wantErrMsg:   "HTTP Status 500: internal server error",
		},
		{
			name:         "should error if server returns 403",
			responseCode: http.StatusForbidden,
			responseBody: "forbidden",
			wantErr:      true,
			wantErrMsg:   "HTTP Status 403: forbidden",
		},
		{
			name:         "should not error if server returns 200",
			responseCode: http.StatusOK,
			responseBody: `<?xml version="1.0" encoding="utf-8"?>
							<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
								<soap:Body>
									<PingResponse xmlns="http://example.com/service.xsd">
										<PingResult>
											<Message>Pong hi</Message>
										</PingResult>
									</PingResponse>
								</soap:Body>
							</soap:Envelope>`,
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.responseCode)
				w.Write([]byte(test.responseBody))
			}))
			defer ts.Close()
			client := NewClient(ts.URL)
			gotErr := client.Call("GetData", &Ping{}, &PingResponse{})
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("Expected an error from call.  Received none")
				}
				requestError, ok := gotErr.(*HTTPError)
				if !ok {
					t.Fatalf("Expected a HTTPError.  Received: %s", gotErr.Error())
				}

				if requestError.StatusCode != test.responseCode {
					t.Fatalf("Unexpected StatusCode.  Got %d", requestError.StatusCode)
				}

				if string(requestError.ResponseBody) != test.responseBody {
					t.Fatalf("Unexpected ResponseBody.  Got %s", requestError.ResponseBody)
				}

				if requestError.Error() != test.wantErrMsg {
					t.Fatalf("Unexpected Error message.  Got %s", requestError.Error())
				}
			} else if gotErr != nil {
				t.Fatalf("Expected no error from call.  Received: %s", gotErr.Error())
			}
		})
	}

}
