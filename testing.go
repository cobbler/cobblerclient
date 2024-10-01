// SPDX-LLicense-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright SUSE LLC

package cobblerclient

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/go-test/deep"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"testing"
)

var config = ClientConfig{
	URL:      "http://localhost:8081/cobbler_api",
	Username: "cobbler",
	Password: "cobbler",
}

// FailOnError ...
func FailOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// Fixture implies that from the context where the test is being run a "fixtures" folder exists.
func Fixture(fn string) ([]byte, error) {
	// Disable semgrep (linter in Codacy) since this is testcode
	return os.ReadFile("./fixtures/" + fn) // nosemgrep
}

type APIResponsePair struct {
	Actual   []byte // This is the actual response that the client gets from the server side.
	Expected []byte // The payload that you expect to receive. This is to verify that your implementation is sending the proper payload to the server.
	Response []byte // The response you want to return.
}

type StubHTTPClient struct {
	t              *testing.T
	answers        []APIResponsePair
	ShouldVerify   bool // Make sure that the expected and the actual sent payload match.
	requestCounter int
}

func NewStubHTTPClient(t *testing.T) *StubHTTPClient {
	s := StubHTTPClient{t: t, ShouldVerify: true}
	return &s
}

func (s *StubHTTPClient) Verify() {
	a := s.answers[s.requestCounter]
	var expec XMLRPCMethodCall
	var obt XMLRPCMethodCall
	var err error

	err = xml.NewDecoder(bytes.NewReader(a.Expected)).Decode(&expec)
	if err != nil {
		s.t.Fatal(err)
	}
	err = xml.NewDecoder(bytes.NewReader(a.Actual)).Decode(&obt)
	if err != nil {
		s.t.Fatal(err)
	}

	sortMethodCall(expec)
	sortMethodCall(obt)
	comparison := deep.Equal(expec, obt)
	fmt.Println(expec)
	fmt.Println(obt)
	fmt.Println(comparison)
	if len(comparison) > 0 {
		s.t.Fatal(comparison)
	}
}

func (s *StubHTTPClient) Post(uri, bodyType string, req io.Reader) (*http.Response, error) {
	b, err := io.ReadAll(req)
	if err != nil {
		s.t.Fatal(err)
	}
	if s.requestCounter >= len(s.answers) {
		s.t.Errorf("Received unbuffered request: %s", b)
		s.t.Fatal("Not enough buffered answers!")
	}
	a := &s.answers[s.requestCounter]

	a.Actual = b
	if s.ShouldVerify {
		s.Verify()
	}
	res := &http.Response{Body: io.NopCloser(bytes.NewBuffer(a.Response))}
	s.requestCounter++
	return res, nil
}

// createStubHTTPClient ...
func createStubHTTPClient(t *testing.T, fixtures []string) Client {
	hc := NewStubHTTPClient(t)

	for _, fixture := range fixtures {
		if fixture != "" {
			rawRequest, err := Fixture(fixture + "-req.xml")
			FailOnError(t, err)
			response, err := Fixture(fixture + "-res.xml")
			FailOnError(t, err)

			// flatten the request so it matches the kolo generated xml
			r := regexp.MustCompile(`\s+<`)
			expectedReq := []byte(r.ReplaceAllString(string(rawRequest), "<"))
			hc.answers = append(hc.answers, APIResponsePair{
				Expected: expectedReq,
				Response: response,
			})
		}
	}

	c := NewClient(hc, config)
	c.Token = "securetoken99"
	return c
}

// createStubHTTPClientSingle ...
func createStubHTTPClientSingle(t *testing.T, fixture string) Client {
	return createStubHTTPClient(t, []string{fixture})
}

func sortMethodCall(methodCall XMLRPCMethodCall) {
	for _, param := range methodCall.Params {
		value := param.Value
		if len(value.Struct.Members) > 0 {
			sortStruct(value.Struct)
		}
		// Arrays should already be correctly handled since they are insertion ordered.
	}
}

func sortStruct(xmlrpcStruct XMLRPCStruct) {
	sort.Slice(xmlrpcStruct.Members, func(i, j int) bool {
		return xmlrpcStruct.Members[i].Name < xmlrpcStruct.Members[j].Name
	})
	for _, value := range xmlrpcStruct.Members {
		if len(value.StructValue.Struct.Members) > 0 {
			sortStruct(value.StructValue.Struct)
		}
	}
}

type XMLRPCStructMember struct {
	Name        string      `xml:"name"`
	StructValue XMLRPCValue `xml:"value"`
}

type XMLRPCStruct struct {
	Members []XMLRPCStructMember `xml:"member"`
}

type XMLRPCArray struct {
	ArrayValues []XMLRPCValue `xml:"data>value"`
}

// XMLRPCValue is a wrapper struct where only one of the fields will be set for a single instance of this struct.
type XMLRPCValue struct {
	XMLType string
	Int     int          `xml:"int"`
	String  string       `xml:"string"`
	Boolean bool         `xml:"boolean"`
	Double  float64      `xml:"double"`
	Base64  string       `xml:"base64"`
	Struct  XMLRPCStruct `xml:"struct"`
	Array   XMLRPCArray  `xml:"array"`
}

type XMLRPCParam struct {
	Value XMLRPCValue `xml:"value"`
}

type XMLRPCMethodCall struct {
	MethodName string        `xml:"methodName"`
	Params     []XMLRPCParam `xml:"params>param"`
}
