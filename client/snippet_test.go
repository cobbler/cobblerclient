package client_test

import (
	"testing"

	"github.com/ContainerSolutions/go-utils"

	"github.com/cobbler/cobblerclient/client"

	cobbler_testing "github.com/cobbler/cobblerclient/internal/testing"
)

func TestCreateSnippet(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "create-snippet-req.xml", "create-snippet-res.xml")
	snippet := client.Snippet{
		Name: "/var/lib/cobbler/snippets/some-snippet",
		Body: "sample content",
	}
	err := c.CreateSnippet(snippet)
	utils.FailOnError(t, err)
}

func TestGetSnippet(t *testing.T) {
	snippetName := "/var/lib/cobbler/snippets/some-snippet"
	c := cobbler_testing.CreateStubHTTPClient(t, "get-snippet-req.xml", "get-snippet-res.xml")
	expectedSnippet := client.Snippet{
		Name: snippetName,
		Body: "sample content",
	}
	returnedSnippet, err := c.GetSnippet(snippetName)
	utils.FailOnError(t, err)
	if returnedSnippet.Body != expectedSnippet.Body {
		t.Errorf("Snippet Body did not match.")
	}
}
