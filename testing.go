// SPDX-LLicense-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright SUSE LLC

package cobblerclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

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
	for _, a := range s.answers {
		if !bytes.Equal(a.Expected, a.Actual) {
			spit("/tmp/expected", a.Expected)
			spit("/tmp/actual", a.Actual)
			s.t.Errorf("expected:\n%sgot:\n%s", a.Expected, a.Actual)
		}
	}
}

func (s *StubHTTPClient) Post(uri, bodyType string, req io.Reader) (*http.Response, error) {
	b, err := io.ReadAll(req)
	if err != nil {
		s.t.Fatal(err)
	}

	s.answers[s.requestCounter].Actual = b
	if s.ShouldVerify {
		s.Verify()
	}
	res := &http.Response{Body: io.NopCloser(bytes.NewBuffer(s.answers[s.requestCounter].Response))}
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
