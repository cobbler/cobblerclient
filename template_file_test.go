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

	"github.com/ContainerSolutions/go-utils"
)

func TestCreateTemplateFile(t *testing.T) {
	c := createStubHTTPClient(t, "create-template-file-req.xml", "create-template-file-res.xml")
	ks := TemplateFile{
		Name: "/var/lib/cobbler/templates/foo.ks",
		Body: "sample content",
	}
	err := c.CreateTemplateFile(ks)
	utils.FailOnError(t, err)
}

func TestGetTemplateFile(t *testing.T) {
	ksName := "/var/lib/cobbler/templates/foo.ks"
	c := createStubHTTPClient(t, "get-template-file-req.xml", "get-template-file-res.xml")
	expectedKS := TemplateFile{
		Name: ksName,
		Body: "sample content",
	}
	returnedKS, err := c.GetTemplateFile(ksName)
	utils.FailOnError(t, err)
	if returnedKS.Body != expectedKS.Body {
		t.Errorf("Template Body did not match.")
	}
}
