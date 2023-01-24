package item_test

import (
	"github.com/cobbler/cobblerclient/item"
	"testing"

	"github.com/ContainerSolutions/go-utils"

	cobblerTesting "github.com/cobbler/cobblerclient/internal/testing"
)

func TestGetRepos(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-repos-req.xml", "get-repos-res.xml")
	repos, err := item.GetRepos(c)
	utils.FailOnError(t, err)
	if len(repos) != 1 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestGetItemNames(t *testing.T) {
	// Arrange
	c := cobblerTesting.CreateStubHTTPClient(t, "get-list-names-repo-req.xml", "get-list-names-repo-res.xml")

	// Act
	var repoNames, err = item.ListRepoNames(c)

	// Assert
	utils.FailOnError(t, err)
	if len(repoNames) != 0 {
		t.Errorf("Non-empty list of repositories detected.")
	}
}

func TestGetRepo(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-repo-req.xml", "get-repo-res.xml")
	r, err := item.GetRepo(c, "rhel-7-x86_64")
	utils.FailOnError(t, err)
	if r.Name.Get() != "rhel-7-x86_64" {
		t.Errorf("Wrong repo returned.")
	}
}

/*
 * NOTE: We're skipping the testing of CREATE, UPDATE, DELETE methods for now because
 *       the current implementation of the StubHTTPClient does not allow
 *       buffered mock responses so as soon as the method makes the second
 *       call to Cobbler it'll fail.
 *       This is a system test, so perhaps we can run Cobbler in a Docker container
 *       and take it from there.
 */
