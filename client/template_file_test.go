package client_test

import (
	"testing"

	"github.com/ContainerSolutions/go-utils"

	"github.com/cobbler/cobblerclient/client"

	cobbler_testing "github.com/cobbler/cobblerclient/internal/testing"
)

func TestCreateTemplateFile(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "create-template-file-req.xml", "create-template-file-res.xml")
	ks := client.TemplateFile{
		Name: "/var/lib/cobbler/templates/foo.ks",
		Body: "sample content",
	}
	err := c.CreateTemplateFile(ks)
	utils.FailOnError(t, err)
}

func TestGetTemplateFile(t *testing.T) {
	ksName := "/var/lib/cobbler/templates/foo.ks"
	c := cobbler_testing.CreateStubHTTPClient(t, "get-template-file-req.xml", "get-template-file-res.xml")
	expectedKS := client.TemplateFile{
		Name: ksName,
		Body: "sample content",
	}
	returnedKS, err := c.GetTemplateFile(ksName)
	utils.FailOnError(t, err)
	if returnedKS.Body != expectedKS.Body {
		t.Errorf("Template Body did not match.")
	}
}
