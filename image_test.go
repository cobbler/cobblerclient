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

func TestNewImage(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewImage()
}

func TestGetImages(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-images")
	images, err := c.GetImages()
	FailOnError(t, err)

	if len(images) != 1 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestGetImage(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-image")
	c.CachedVersion = CobblerVersion{3, 3, 2}

	// Act
	image, err := c.GetImage("testimage", false, false)

	// Assert
	FailOnError(t, err)
	if image.Name != "testimage" {
		t.Errorf("Wrong image returned.")
	}
}

func TestDeleteImage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-image")
	err := c.DeleteImage("test")
	FailOnError(t, err)
}

func TestDeleteImageRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-image")
	err := c.DeleteImageRecursive("test", false)
	FailOnError(t, err)
}

func TestListImageNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-image")
	images, err := c.ListImageNames()
	FailOnError(t, err)

	if len(images) != 1 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestGetImagesSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-images-since")
	images, err := c.GetImagesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(images) != 1 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestFindImage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-image")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testimage"
	images, err := c.FindImage(criteria)
	FailOnError(t, err)

	if len(images) != 1 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestFindImageNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-image-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testimage"
	images, err := c.FindImageNames(criteria)
	FailOnError(t, err)

	if len(images) != 1 {
		t.Error("Wrong number of images returned.")
	}
}

func TestSaveImage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-image")
	err := c.SaveImage("image::testimage", "bypass")
	FailOnError(t, err)
}

func TestCopyImage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-image")
	err := c.CopyImage("image::testimage", "testimage2")
	FailOnError(t, err)
}

func TestRenameImage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-image")
	err := c.RenameImage("image::testimage2", "testimage1")
	FailOnError(t, err)
}

func TestGetImageHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-image-handle")
	res, err := c.GetImageHandle("testimage")
	FailOnError(t, err)

	if res != "image::testimage" {
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
