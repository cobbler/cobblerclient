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
	"time"
)

func TestNewTemplate(t *testing.T) {
	tpl := NewTemplate()
	if tpl.URI.Schema != TemplateSchemaFile {
		t.Errorf("default URI.Schema = %v, want File", tpl.URI.Schema)
	}
	if tpl.Tags == nil {
		t.Errorf("Tags should be non-nil")
	}
}

func TestTemplateSchemaString(t *testing.T) {
	// TemplateSchema is a plain string alias (see template_types.go) so the XML-RPC codec, which picks its wire
	// representation from reflect.Kind, sends these as <string> rather than <int> — this just pins the wire
	// values the constants must keep sending.
	cases := []struct {
		s        TemplateSchema
		expected string
	}{
		{TemplateSchemaFile, "file"},
		{TemplateSchemaImportlib, "importlib"},
		{TemplateSchemaEnvironment, "environment"},
	}
	for _, tc := range cases {
		if tc.s != tc.expected {
			t.Errorf("TemplateSchema = %q, want %q", tc.s, tc.expected)
		}
	}
}

func TestGetTemplateFileForProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-template-file-for-profile")
	res, err := c.GetTemplateFileForProfile("testprof", "/etc/motd")
	FailOnError(t, err)

	if res == "" {
		t.Error("Expected non-empty template file content.")
	}
}

func TestGetTemplateFileForSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-template-file-for-system")
	res, err := c.GetTemplateFileForSystem("testsys", "/etc/motd")
	FailOnError(t, err)

	if res == "" {
		t.Error("Expected non-empty template file content.")
	}
}

func TestGetTemplates(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-templates")
	templates, err := c.GetTemplates()
	FailOnError(t, err)
	if len(templates) != 1 {
		t.Errorf("Wrong number of templates returned.")
	}
}

func TestGetTemplate(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-template")
	tpl, err := c.GetTemplate("testtemplate", false, false)
	FailOnError(t, err)
	if tpl.Name != "testtemplate" {
		t.Errorf("Wrong template name returned: %v", tpl.Name)
	}
}

func TestGetTemplateHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-template-handle")
	handle, err := c.GetTemplateHandle("testtemplate")
	FailOnError(t, err)
	if handle != "template::testtemplate" {
		t.Errorf("wrong handle: %q", handle)
	}
}

func TestGetTemplatesSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-templates-since")
	templates, err := c.GetTemplatesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)
	if len(templates) != 1 {
		t.Errorf("Wrong number of templates returned.")
	}
}

func TestFindTemplate(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-template")
	criteria := map[string]interface{}{"name": "testtemplate"}
	res, err := c.FindTemplate(criteria)
	FailOnError(t, err)
	if len(res) != 1 {
		t.Errorf("Expected 1 template, got %d.", len(res))
	}
}

func TestFindTemplateNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-template-names")
	criteria := map[string]interface{}{"name": "testtemplate"}
	names, err := c.FindTemplateNames(criteria)
	FailOnError(t, err)
	if len(names) != 1 {
		t.Errorf("Expected 1 template name, got %d.", len(names))
	}
}

func TestListTemplateNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-template")
	names, err := c.ListTemplateNames()
	FailOnError(t, err)
	if len(names) != 1 {
		t.Errorf("Expected 1 template name, got %d.", len(names))
	}
}

func TestGetTemplateContent(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-template-content")
	content, err := c.GetTemplateContent("template-uid-1")
	FailOnError(t, err)
	if content == "" {
		t.Error("Expected non-empty template content.")
	}
}

func TestTemplatesRefreshContent(t *testing.T) {
	c := createStubHTTPClientSingle(t, "templates-refresh-content")
	err := c.TemplatesRefreshContent(nil)
	FailOnError(t, err)
}

func TestBackgroundTemplatesRefreshContent(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-templates-refresh-content")
	eventID, err := c.BackgroundTemplatesRefreshContent(nil)
	FailOnError(t, err)
	if eventID == "" {
		t.Error("Expected non-empty event ID.")
	}
}

func TestSaveTemplate(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-template")
	err := c.SaveTemplate("template::testtemplate", true, true, "bypass")
	FailOnError(t, err)
}

func TestCopyTemplate(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-template")
	err := c.CopyTemplate("template::testtemplate", "testtemplate-copy")
	FailOnError(t, err)
}

func TestDeleteTemplate(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-template")
	err := c.DeleteTemplate("testtemplate")
	FailOnError(t, err)
}

func TestDeleteTemplateRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-template")
	err := c.DeleteTemplateRecursive("testtemplate", false)
	FailOnError(t, err)
}

func TestRenameTemplate(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-template")
	err := c.RenameTemplate("template::testtemplate", "testtemplate-new")
	FailOnError(t, err)
}

func TestCreateTemplate(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"new-template",
		"new-template-modify-name",
		"new-template-modify-comment",
		"new-template-modify-kernel-options",
		"new-template-modify-kernel-options-post",
		"new-template-modify-autoinstall-meta",
		"new-template-modify-fetchable-files",
		"new-template-modify-boot-files",
		"new-template-modify-template-files",
		"new-template-modify-owners",
		"new-template-modify-template-type",
		"new-template-modify-uri",
		"new-template-modify-tags",
		"new-template-modify-content",
		"new-template-save",
		"new-template-get",
	})
	tpl := NewTemplate()
	tpl.Name = "testtemplate"
	result, err := c.CreateTemplate(tpl)
	FailOnError(t, err)
	if result.Name != "testtemplate" {
		t.Errorf("Wrong template name returned: %v", result.Name)
	}
}

func TestUpdateTemplate(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"update-template-handle",
		"update-template-modify-name",
		"update-template-modify-comment",
		"update-template-modify-kernel-options",
		"update-template-modify-kernel-options-post",
		"update-template-modify-autoinstall-meta",
		"update-template-modify-fetchable-files",
		"update-template-modify-boot-files",
		"update-template-modify-template-files",
		"update-template-modify-owners",
		"update-template-modify-template-type",
		"update-template-modify-uri",
		"update-template-modify-tags",
		"update-template-modify-content",
		"update-template-save",
	})
	tpl := NewTemplate()
	tpl.Name = "testtemplate"
	err := c.UpdateTemplate(&tpl)
	FailOnError(t, err)
}
