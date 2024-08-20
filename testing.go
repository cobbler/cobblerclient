// SPDX-LLicense-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright SUSE LLC

package cobblerclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
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

func removeLineBreaks(input []byte) []byte {
	return []byte(strings.Replace(string(input), "\r", "", -1))
}

func (s *StubHTTPClient) Verify() {
	a := s.answers[s.requestCounter]
	if runtime.GOOS == "windows" {
		a.Expected = removeLineBreaks(a.Expected)
		a.Actual = removeLineBreaks(a.Actual)
	}
	if !bytes.Equal(a.Expected, a.Actual) {
		spit("/tmp/expected", a.Expected)
		spit("/tmp/actual", a.Actual)
		s.t.Errorf("expected:\n%sgot:\n%s", a.Expected, a.Actual)
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

func spit(path string, b []byte) {
	file, err := os.Create(path)
	if err != nil {
		return
	}

	n, err := file.Write(b)
	if err != nil {
		return
	}

	fmt.Printf("%v bytes written to %s\n", n, path)
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
