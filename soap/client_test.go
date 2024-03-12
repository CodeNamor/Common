package soap

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CodeNamor/Common/logging"
	"github.com/CodeNamor/Common/logging/logfields"
	"github.com/stretchr/testify/assert"
)

type Ping struct {
	XMLName xml.Name     `xml:"http://example.com/service.xsd Ping"`
	Request *PingRequest `xml:"request,omitempty"`
}

type PingRequest struct {
	// XMLName xml.Name `xml:"http://example.com/service.xsd PingRequest"`
	Message string `xml:"Message,omitempty"`
}

type PingResponse struct {
	XMLName    xml.Name   `xml:"http://example.com/service.xsd PingResponse"`
	PingResult *PingReply `xml:"PingResult,omitempty"`
}

type PingReply struct {
	// XMLName xml.Name `xml:"http://example.com/service.xsd PingReply"`
	Message string `xml:"Message,omitempty"`
}

type MockRequestClient struct {
	client *http.Client
}

func (r MockRequestClient) Do(req *http.Request) (*http.Response, error) {
	return r.client.Do(req)
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

	mockClient := MockRequestClient{client: &http.Client{}}
	client := NewClient(mockClient, ts.URL, logging.WithField(logfields.RequestId, ""))
	req := &Ping{Request: &PingRequest{Message: "Hi"}}
	reply := &PingResponse{}

	if err := client.Call("GetData", req, reply); err != nil {
		t.Fatalf("couln't call service: %v", err)
	}

	assert.Equal(t, "Pong hi", reply.PingResult.Message)
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
	logEntry := logging.WithField(logfields.RequestId, "")
	for _, test := range tests {
		client := NewClient(MockRequestClient{client: &http.Client{}}, ts.URL, logEntry, WithHTTPHeaders(test.reqHeaders))
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
