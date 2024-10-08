/*
Copyright 2015 Container Solutions

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cobblerclient

import (
	"testing"
)

func TestCreateSnippet(t *testing.T) {
	c := createStubHTTPClientSingle(t, "create-snippet")
	snippet := Snippet{
		Name: "/var/lib/cobbler/snippets/some-snippet",
		Body: "sample content",
	}
	err := c.CreateSnippet(snippet)
	FailOnError(t, err)
}

func TestGetSnippet(t *testing.T) {
	snippetName := "/var/lib/cobbler/snippets/some-snippet"
	c := createStubHTTPClientSingle(t, "get-snippet")
	expectedSnippet := Snippet{
		Name: snippetName,
		Body: "sample content",
	}
	returnedSnippet, err := c.GetSnippet(snippetName)
	FailOnError(t, err)
	if returnedSnippet.Body != expectedSnippet.Body {
		t.Errorf("Snippet Body did not match.")
	}
}
