package testing

import (
	"regexp"
	"testing"

	"github.com/ContainerSolutions/go-utils"
	"github.com/cobbler/cobblerclient/client"
)

var config = client.ClientConfig{
	URL:      "http://localhost:8081/cobbler_api",
	Username: "cobbler",
	Password: "cobbler",
}

func CreateStubHTTPClient(t *testing.T, reqFixture string, resFixture string) client.Client {
	hc := utils.NewStubHTTPClient(t)

	if reqFixture != "" {
		rawRequest, err := utils.Fixture(reqFixture)
		utils.FailOnError(t, err)

		// flatten the request so it matches the kolo generated xml
		r := regexp.MustCompile(`\s+<`)
		expectedReq := []byte(r.ReplaceAllString(string(rawRequest), "<"))
		hc.Expected = expectedReq
	}

	if resFixture != "" {
		response, err := utils.Fixture(resFixture)
		utils.FailOnError(t, err)
		hc.Response = response
	}

	c := client.NewClient(hc, config)
	c.Token = "securetoken99"
	return c
}
