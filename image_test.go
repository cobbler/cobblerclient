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

func TestGetImages(t *testing.T) {
	c := createStubHTTPClient(t, "get-images-req.xml", "get-images-res.xml")
	images, err := c.GetImages()
	utils.FailOnError(t, err)

	if len(images) != 1 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestGetImage(t *testing.T) {
	c := createStubHTTPClient(t, "get-image-req.xml", "get-image-res.xml")
	image, err := c.GetImage("testimage")
	utils.FailOnError(t, err)

	if image.Name != "testimage" {
		t.Errorf("Wrong image returned.")
	}
}

func TestDeleteImage(t *testing.T) {
	c := createStubHTTPClient(t, "delete-image-req.xml", "delete-image-res.xml")
	err := c.DeleteImage("test")
	utils.FailOnError(t, err)
}

func TestDeleteImageRecursive(t *testing.T) {
	c := createStubHTTPClient(t, "delete-image-req.xml", "delete-image-res.xml")
	err := c.DeleteImageRecursive("test", false)
	utils.FailOnError(t, err)
}

func TestListImageNames(t *testing.T) {
	c := createStubHTTPClient(t, "get-item-names-image-req.xml", "get-item-names-image-res.xml")
	images, err := c.ListImageNames()
	utils.FailOnError(t, err)

	if len(images) != 1 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestGetImagesSince(t *testing.T) {
	c := createStubHTTPClient(t, "get-images-since-req.xml", "get-images-since-res.xml")
	images, err := c.GetImagesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	utils.FailOnError(t, err)

	if len(images) != 1 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestFindImage(t *testing.T) {
	c := createStubHTTPClient(t, "find-image-req.xml", "find-image-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testimage"
	images, err := c.FindImage(criteria)
	utils.FailOnError(t, err)

	if len(images) != 1 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestFindImageNames(t *testing.T) {
	c := createStubHTTPClient(t, "find-image-names-req.xml", "find-image-names-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testimage"
	images, err := c.FindImageNames(criteria)
	utils.FailOnError(t, err)

	if len(images) != 1 {
		t.Error("Wrong number of images returned.")
	}
}

func TestSaveImage(t *testing.T) {
	c := createStubHTTPClient(t, "save-image-req.xml", "save-image-res.xml")
	err := c.SaveImage("image::testimage", "bypass")
	utils.FailOnError(t, err)
}

func TestCopyImage(t *testing.T) {
	c := createStubHTTPClient(t, "copy-image-req.xml", "copy-image-res.xml")
	err := c.CopyImage("image::testimage", "testimage2")
	utils.FailOnError(t, err)
}

func TestRenameImage(t *testing.T) {
	c := createStubHTTPClient(t, "rename-image-req.xml", "rename-image-res.xml")
	err := c.RenameImage("image::testimage2", "testimage1")
	utils.FailOnError(t, err)
}

func TestGetImageHandle(t *testing.T) {
	c := createStubHTTPClient(t, "get-image-handle-req.xml", "get-image-handle-res.xml")
	res, err := c.GetImageHandle("testimage")
	utils.FailOnError(t, err)

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
