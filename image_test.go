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

	if len(images) != 3 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestGetImage(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-image")

	// Act
	image, err := c.GetImage("00000000000000000000000000000040", false, false)

	// Assert
	FailOnError(t, err)
	if image.Name != "testimage" {
		t.Errorf("Wrong image returned.")
	}
}

func TestDeleteImage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-image")
	err := c.DeleteImage("00000000000000000000000000000041")
	FailOnError(t, err)
}

func TestDeleteImageRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-image")
	err := c.DeleteImageRecursive("00000000000000000000000000000041", false)
	FailOnError(t, err)
}

func TestListImageNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-image")
	images, err := c.ListImageNames()
	FailOnError(t, err)

	if len(images) != 3 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestGetImagesSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-images-since")
	images, err := c.GetImagesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(images) != 3 {
		t.Errorf("Wrong number of images returned.")
	}
}

func TestFindImage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-image")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testimage"
	images, err := c.FindImage(criteria, false)
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
	err := c.SaveImage("000000000000000000000000000000ce", true, true, "bypass")
	FailOnError(t, err)
}

func TestCopyImage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-image")
	err := c.CopyImage("000000000000000000000000000000ce", "testimage2")
	FailOnError(t, err)
}

func TestRenameImage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-image")
	err := c.RenameImage("000000000000000000000000000000cf", "testimage1")
	FailOnError(t, err)
}

func TestGetImageHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-image-handle")
	res, err := c.GetImageHandle("testimage")
	FailOnError(t, err)

	if res != "000000000000000000000000000000ce" {
		t.Error("Wrong object id returned.")
	}
}

func TestCreateImage(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"create-image-name-check",
		"new-image",
		"new-image-modify-name",
		"new-image-modify-comment",
		"new-image-modify-kernel-options",
		"new-image-modify-kernel-options-post",
		"new-image-modify-autoinstall-meta",
		"new-image-modify-template-files",
		"new-image-modify-owners",
		"new-image-modify-arch",
		"new-image-modify-autoinstall",
		"new-image-modify-breed",
		"new-image-modify-file",
		"new-image-modify-image-type",
		"new-image-modify-network-count",
		"new-image-modify-os-version",
		"new-image-modify-boot-loaders",
		"new-image-modify-menu",
		"new-image-modify-virt-auto-boot",
		"new-image-modify-virt-cpus",
		"new-image-modify-virt-disk-driver",
		"new-image-modify-virt-file-size",
		"new-image-modify-virt-path",
		"new-image-modify-virt-pxe-boot",
		"new-image-modify-virt-ram",
		"new-image-modify-virt-type",
		"new-image-modify-virt-uefi",
		"new-image-modify-virt-bridge",
		"new-image-save",
		"new-image-get",
	})
	image := NewImage()
	image.Name = "testimage"
	// Menu must be the referenced menu's real uid, not its name - see the doc comment on
	// Image.Menu.
	image.Menu = "00000000000000000000000000000020"
	image.Virt.Path = "/var/lib/libvirt/images"
	image.Virt.Cpus = Value[int]{Data: 2}

	result, err := c.CreateImage(image)
	FailOnError(t, err)
	if result.Name != "testimage" {
		t.Errorf("Wrong image name returned: %v", result.Name)
	}
	if result.Menu != "00000000000000000000000000000020" {
		t.Errorf("Menu uid was not passed through unmodified, got %q.", result.Menu)
	}
}

func TestUpdateImage(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"update-image-handle",
		"update-image-modify-name",
		"update-image-modify-comment",
		"update-image-modify-kernel-options",
		"update-image-modify-kernel-options-post",
		"update-image-modify-autoinstall-meta",
		"update-image-modify-template-files",
		"update-image-modify-owners",
		"update-image-modify-arch",
		"update-image-modify-autoinstall",
		"update-image-modify-breed",
		"update-image-modify-file",
		"update-image-modify-image-type",
		"update-image-modify-network-count",
		"update-image-modify-os-version",
		"update-image-modify-boot-loaders",
		"update-image-modify-menu",
		"update-image-modify-virt-auto-boot",
		"update-image-modify-virt-cpus",
		"update-image-modify-virt-disk-driver",
		"update-image-modify-virt-file-size",
		"update-image-modify-virt-path",
		"update-image-modify-virt-pxe-boot",
		"update-image-modify-virt-ram",
		"update-image-modify-virt-type",
		"update-image-modify-virt-uefi",
		"update-image-modify-virt-bridge",
		"update-image-save",
	})
	image := NewImage()
	image.Name = "testimage"
	image.Virt.Path = "/var/lib/libvirt/images"
	image.Virt.Cpus = Value[int]{Data: 2}

	err := c.UpdateImage(&image)
	FailOnError(t, err)
}
