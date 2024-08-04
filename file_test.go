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

	"github.com/ContainerSolutions/go-utils"
)

func TestGetFiles(t *testing.T) {
	c := createStubHTTPClient(t, "get-files-req.xml", "get-files-res.xml")
	files, err := c.GetFiles()
	utils.FailOnError(t, err)

	if len(files) != 1 {
		t.Errorf("Wrong number of files returned.")
	}
}

func TestGetFile(t *testing.T) {
	c := createStubHTTPClient(t, "get-file-req.xml", "get-file-res.xml")
	file, err := c.GetFile("testfile")
	utils.FailOnError(t, err)

	if file.Name != "testfile" {
		t.Errorf("Wrong file returned.")
	}
}

func TestDeleteFile(t *testing.T) {
	c := createStubHTTPClient(t, "delete-file-req.xml", "delete-file-res.xml")
	err := c.DeleteFile("test")
	utils.FailOnError(t, err)
}

func TestDeleteFileRecursive(t *testing.T) {
	c := createStubHTTPClient(t, "delete-file-req.xml", "delete-file-res.xml")
	err := c.DeleteFileRecursive("test", false)
	utils.FailOnError(t, err)
}

func TestListFileNames(t *testing.T) {
	c := createStubHTTPClient(t, "get-item-names-file-req.xml", "get-item-names-file-res.xml")
	files, err := c.ListFileNames()
	utils.FailOnError(t, err)

	if len(files) != 1 {
		t.Errorf("Wrong number of files returned.")
	}
}

func TestGetFilesSince(t *testing.T) {
	c := createStubHTTPClient(t, "get-files-since-req.xml", "get-files-since-res.xml")
	files, err := c.GetFilesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	utils.FailOnError(t, err)

	if len(files) != 1 {
		t.Errorf("Wrong number of files returned.")
	}
}

func TestFindFile(t *testing.T) {
	c := createStubHTTPClient(t, "find-file-req.xml", "find-file-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testfile"
	files, err := c.FindFile(criteria)
	utils.FailOnError(t, err)

	if len(files) != 1 {
		t.Errorf("Wrong number of files returned.")
	}
}

func TestFindFileNames(t *testing.T) {
	c := createStubHTTPClient(t, "find-file-names-req.xml", "find-file-names-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testfile"
	files, err := c.FindFileNames(criteria)
	utils.FailOnError(t, err)

	if len(files) != 1 {
		t.Error("Wrong number of files returned.")
	}
}

func TestSaveFile(t *testing.T) {
	c := createStubHTTPClient(t, "save-file-req.xml", "save-file-res.xml")
	err := c.SaveFile("file::testfile", "bypass")
	utils.FailOnError(t, err)
}

func TestCopyFile(t *testing.T) {
	c := createStubHTTPClient(t, "copy-file-req.xml", "copy-file-res.xml")
	err := c.CopyFile("file::testfile", "testfile2")
	utils.FailOnError(t, err)
}

func TestRenameFile(t *testing.T) {
	c := createStubHTTPClient(t, "rename-file-req.xml", "rename-file-res.xml")
	err := c.RenameFile("file::testfile2", "testfile1")
	utils.FailOnError(t, err)
}

func TestGetFileHandle(t *testing.T) {
	c := createStubHTTPClient(t, "get-file-handle-req.xml", "get-file-handle-res.xml")
	res, err := c.GetFileHandle("testfile")
	utils.FailOnError(t, err)

	if res != "file::testfile" {
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
