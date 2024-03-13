package soap

import (
	"bytes"
	"encoding/xml"
	"io/ioutil" //TODO: migrate to package io
	"net/http"

	"github.com/CodeNamor/Common/requestclient"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Envelope represents the top level of the SOAP message.
type Envelope struct {
	NSxsdAttr string   `xml:"xmlns:xsd,attr"`
	NSxsiAttr string   `xml:"xmlns:xsi,attr"`
	Attr      xml.Attr `xml:",attr,omitempty"`
	XMLName   xml.Name
	Header    *Header
	Body      Body
}

// Header encapsulates the SOAP header.
type Header struct {
	XMLName xml.Name      `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Items   []interface{} `xml:",omitempty"`
}

type Body struct {
	XMLName xml.Name
	Fault   *Fault      `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

// UnmarshalXML unmarshals SOAPBody xml.
func (b *Body) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}
		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &Fault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}
				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}
				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}
	return nil
}

type Fault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`
	Code    string   `xml:"faultcode,omitempty"`
	String  string   `xml:"faultstring,omitempty"`
	Actor   string   `xml:"faultactor,omitempty"`
	Detail  string   `xml:"detail,omitempty"`
}

func (f *Fault) Error() string {
	return f.String
}

const (
	// Predefined WSS namespaces to be used in
	WssNsWSSE string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	WssNsWSU  string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	WssNsType string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText"
)

type WSSSecurityHeader struct {
	XMLName        xml.Name          `xml:"http://schemas.xmlsoap.org/soap/envelope/ wsse:Security"`
	XmlNSWsse      string            `xml:"xmlns:wsse,attr"`
	MustUnderstand string            `xml:"mustUnderstand,attr,omitempty"`
	Token          *WSSUsernameToken `xml:",omitempty"`
}

type WSSUsernameToken struct {
	XMLName   xml.Name     `xml:"wsse:UsernameToken"`
	XmlNSWsu  string       `xml:"xmlns:wsu,attr"`
	XmlNSWsse string       `xml:"xmlns:wsse,attr"`
	Id        string       `xml:"wsu:Id,attr,omitempty"`
	Username  *WSSUsername `xml:",omitempty"`
	Password  *WSSPassword `xml:",omitempty"`
}

type WSSUsername struct {
	XMLName   xml.Name `xml:"wsse:Username"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`
	Data      string   `xml:",chardata"`
}

type WSSPassword struct {
	XMLName   xml.Name `xml:"wsse:Password"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`
	XmlNSType string   `xml:"Type,attr"`
	Data      string   `xml:",chardata"`
}

// NewWSSSecurityHeader creates WSSSecurityHeader instance
func NewWSSSecurityHeader(user, pass, tokenID, mustUnderstand string) *WSSSecurityHeader {
	return &WSSSecurityHeader{
		XmlNSWsse:      WssNsWSSE,
		MustUnderstand: mustUnderstand,
		Token: &WSSUsernameToken{
			XmlNSWsu:  WssNsWSU,
			XmlNSWsse: WssNsWSSE,
			Id:        tokenID,
			Username:  &WSSUsername{XmlNSWsse: WssNsWSSE, Data: user},
			Password:  &WSSPassword{XmlNSWsse: WssNsWSSE, XmlNSType: WssNsType, Data: pass},
		},
	}
}

type basicAuth struct {
	Login    string
	Password string
}

// Most httpClient options are already being set on httpClient that is passed in
type options struct {
	auth        *basicAuth
	httpHeaders map[string]string
}

var defaultOptions = options{}

// A Option sets options such as credentials, tls, etc.
type Option func(*options)

// WithBasicAuth is an Option to set BasicAuth
func WithBasicAuth(login, password string) Option {
	return func(o *options) {
		o.auth = &basicAuth{Login: login, Password: password}
	}
}

// WithHTTPHeaders is an Option to set global HTTP headers for all requests
func WithHTTPHeaders(headers map[string]string) Option {
	return func(o *options) {
		o.httpHeaders = headers
	}
}

// Client is soap client
type Client struct {
	httpClient requestclient.RequestClient
	url        string
	logEntry   *logrus.Entry
	opts       *options
	headers    []interface{}
}

// NewClient creates new SOAP client instance
func NewClient(httpClient requestclient.RequestClient, url string, logEntry *logrus.Entry, opt ...Option) *Client {

	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	return &Client{
		httpClient: httpClient,
		url:        url,
		logEntry:   logEntry,
		opts:       &opts,
	}
}

// AddHeader adds envelope header
func (s *Client) AddHeader(header interface{}) {
	s.headers = append(s.headers, header)
}

// Call performs HTTP POST request
func (s *Client) Call(soapAction string, request, response interface{}) error {
	var envelope Envelope
	soapRequest, ok := request.(Request)
	if ok {
		envelope = soapRequest.GetSoapEnvelope()
	} else {
		envelope = Envelope{
			NSxsdAttr: "http://www.w3.org/2001/XMLSchema",
			NSxsiAttr: "http://www.w3.org/2001/XMLSchema-instance",
			XMLName: xml.Name{
				Space: "http://schemas.xmlsoap.org/soap/envelope/",
				Local: "Envelope",
			},
			Body: Body{
				XMLName: xml.Name{
					Space: "http://schemas.xmlsoap.org/soap/envelope/",
					Local: "Body",
				},
				Content: request,
			},
		}

	}

	if s.headers != nil && len(s.headers) > 0 {
		soapHeader := &Header{Items: make([]interface{}, len(s.headers))}
		copy(soapHeader.Items, s.headers)
		envelope.Header = soapHeader
	}

	requestBodyBuffer, err := encodeEnvelopeIntoBuffer(&envelope)
	if err != nil {
		return err
	}

	// we log.info request (and response if available) on errors already
	// fmt.Println("buffer", requestBodyBuffer.String()) // raw soap request

	req, err := http.NewRequest("POST", s.url, requestBodyBuffer)
	if err != nil {
		return err
	}
	if s.opts.auth != nil {
		req.SetBasicAuth(s.opts.auth.Login, s.opts.auth.Password)
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", soapAction)
	req.Header.Set("User-Agent", "gowsdl/0.1")
	if s.opts.httpHeaders != nil {
		for k, v := range s.opts.httpHeaders {
			req.Header.Set(k, v)
		}
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		logSoapRequest(s.logEntry, &envelope)
		return errors.Wrapf(err, "Error calling making soap request url: %s", s.url)
	}
	if res == nil {
		logSoapRequest(s.logEntry, &envelope)
		return errors.Wrapf(err, "Response was nil: %s", s.url)
	}

	// read all of body and close as quickly as possible to
	// return the connection to the pool
	rawResponseBody, err := readAllOfResponseAndClose(res)
	if err != nil {
		logSoapRequest(s.logEntry, &envelope)
		return err
	}
	if len(rawResponseBody) == 0 {
		return nil
	}

	// we log.info request (and response if available) on errors already
	// fmt.Println("response rawbody", string(rawResponseBody))  // raw response

	respEnvelope := new(Envelope)
	respEnvelope.Body = Body{Content: response}
	err = xml.Unmarshal(rawResponseBody, respEnvelope)
	if err != nil {
		logSoapResponse(s.logEntry, rawResponseBody)
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		logSoapResponse(s.logEntry, rawResponseBody)
		return fault
	}

	// successful so only do this extra work of
	// creating log output if level is trace
	if s.logEntry.Logger.IsLevelEnabled(logrus.TraceLevel) {
		logSoapRequest(s.logEntry, &envelope)
		logSoapResponse(s.logEntry, rawResponseBody)
	}

	return nil
}

func encodeEnvelopeIntoBuffer(envelope *Envelope) (*bytes.Buffer, error) {
	requestBodyBuffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(requestBodyBuffer)
	if err := encoder.Encode(*envelope); err != nil {
		return nil, err
	}

	if err := encoder.Flush(); err != nil {
		return nil, err
	}
	return requestBodyBuffer, nil
}

// readAllOfResponseAndClose is designed to read all of the httpResponse body and
// close as quickly as possible. If there are any errors, panics, or timeouts during
// the read, the response body is still properly closed.
func readAllOfResponseAndClose(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return []byte{}, errors.Errorf("Soap call returned status %v", res.Status)
	}
	return ioutil.ReadAll(res.Body)
}

// logSoapRequest is used to log the soap request to log which is typically
// only done if there is an error. Since the buffer is already gone by this
// point, we recreate it from the envelope.
func logSoapRequest(logEntry *logrus.Entry, envelope *Envelope) {
	requestBodyBuffer, err := encodeEnvelopeIntoBuffer(envelope)
	if err == nil {
		logEntry.Info("soapRequest:", requestBodyBuffer)
	}
}

func logSoapResponse(logEntry *logrus.Entry, response []byte) {
	logEntry.Info("soapResponse:", string(response))
}
