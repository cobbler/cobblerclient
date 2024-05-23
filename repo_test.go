/*
Copyright 2017 HomeAway, Inc

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
	"time"

	"github.com/ContainerSolutions/go-utils"
)

func TestGetRepos(t *testing.T) {
	c := createStubHTTPClient(t, "get-repos-req.xml", "get-repos-res.xml")
	repos, err := c.GetRepos()
	utils.FailOnError(t, err)
	if len(repos) != 1 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestGetRepo(t *testing.T) {
	c := createStubHTTPClient(t, "get-repo-req.xml", "get-repo-res.xml")
	repo, err := c.GetRepo("rhel-7-x86_64")
	utils.FailOnError(t, err)
	if repo.Name != "rhel-7-x86_64" {
		t.Errorf("Wrong repo returned.")
	}
}

func TestListRepoNames(t *testing.T) {
	c := createStubHTTPClient(t, "get-item-names-repo-req.xml", "get-item-names-repo-res.xml")
	repos, err := c.ListRepoNames()
	utils.FailOnError(t, err)

	if len(repos) != 0 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestGetReposSince(t *testing.T) {
	c := createStubHTTPClient(t, "get-repo-since-req.xml", "get-repo-since-res.xml")
	repos, err := c.GetReposSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	utils.FailOnError(t, err)

	if len(repos) != 1 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestFindRepo(t *testing.T) {
	c := createStubHTTPClient(t, "find-repo-req.xml", "find-repo-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	repos, err := c.FindRepo(criteria)
	utils.FailOnError(t, err)

	if len(repos) != 1 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestFindRepoNames(t *testing.T) {
	c := createStubHTTPClient(t, "find-repo-names-req.xml", "find-repo-names-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	repos, err := c.FindRepoNames(criteria)
	utils.FailOnError(t, err)

	if len(repos) != 1 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestSaveRepo(t *testing.T) {
	c := createStubHTTPClient(t, "save-repo-req.xml", "save-repo-res.xml")
	err := c.SaveRepo("repo::testrepo", "bypass")
	utils.FailOnError(t, err)
}

func TestCopyRepo(t *testing.T) {
	c := createStubHTTPClient(t, "copy-repo-req.xml", "copy-repo-res.xml")
	err := c.CopyRepo("repo::testrepo", "testrepo2")
	utils.FailOnError(t, err)
}

func TestRenameRepo(t *testing.T) {
	c := createStubHTTPClient(t, "rename-repo-req.xml", "rename-repo-res.xml")
	err := c.RenameRepo("repo::testrepo2", "testrepo1")
	utils.FailOnError(t, err)
}

func TestGetRepoHandle(t *testing.T) {
	c := createStubHTTPClient(t, "get-repo-handle-req.xml", "get-repo-handle-res.xml")
	res, err := c.GetRepoHandle("testrepo")
	utils.FailOnError(t, err)

	if res != "repo::testrepo" {
		t.Error("Wrong object id returned.")
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
