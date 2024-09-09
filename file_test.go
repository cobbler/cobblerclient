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

func TestNewFile(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewFile()
}

func TestGetFiles(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-files")
	files, err := c.GetFiles()
	FailOnError(t, err)

	if len(files) != 1 {
		t.Errorf("Wrong number of files returned.")
	}
}

func TestGetFile(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-file")
	c.CachedVersion = CobblerVersion{3, 3, 2}

	// Act
	file, err := c.GetFile("testfile", false, false)

	// Assert
	FailOnError(t, err)
	if file.Name != "testfile" {
		t.Errorf("Wrong file returned.")
	}
}

func TestDeleteFile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-file")
	err := c.DeleteFile("test")
	FailOnError(t, err)
}

func TestDeleteFileRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-file")
	err := c.DeleteFileRecursive("test", false)
	FailOnError(t, err)
}

func TestListFileNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-file")
	files, err := c.ListFileNames()
	FailOnError(t, err)

	if len(files) != 1 {
		t.Errorf("Wrong number of files returned.")
	}
}

func TestGetFilesSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-files-since")
	files, err := c.GetFilesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(files) != 1 {
		t.Errorf("Wrong number of files returned.")
	}
}

func TestFindFile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-file")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testfile"
	files, err := c.FindFile(criteria)
	FailOnError(t, err)

	if len(files) != 1 {
		t.Errorf("Wrong number of files returned.")
	}
}

func TestFindFileNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-file-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testfile"
	files, err := c.FindFileNames(criteria)
	FailOnError(t, err)

	if len(files) != 1 {
		t.Error("Wrong number of files returned.")
	}
}

func TestSaveFile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-file")
	err := c.SaveFile("file::testfile", "bypass")
	FailOnError(t, err)
}

func TestCopyFile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-file")
	err := c.CopyFile("file::testfile", "testfile2")
	FailOnError(t, err)
}

func TestRenameFile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-file")
	err := c.RenameFile("file::testfile2", "testfile1")
	FailOnError(t, err)
}

func TestGetFileHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-file-handle")
	res, err := c.GetFileHandle("testfile")
	FailOnError(t, err)

	if res != "file::testfile" {
		t.Error("Wrong object id returned.")
	}
}

func TestGetFileAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-file-as-rendered")
	res, err := c.GetFileAsRendered("testfile")
	FailOnError(t, err)

	if res["name"] != "testfile" {
		t.Error("Wrong object name returned.")
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
