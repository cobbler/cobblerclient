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

/*
 * NOTE: We're skipping the testing of CREATE, UPDATE, DELETE methods for now because
 *       the current implementation of the StubHTTPClient does not allow
 *       buffered mock responses so as soon as the method makes the second
 *       call to Cobbler it'll fail.
 *       This is a system test, so perhaps we can run Cobbler in a Docker container
 *       and take it from there.
 */
