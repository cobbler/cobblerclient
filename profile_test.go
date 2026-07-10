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

func TestNewProfile(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewProfile()
}

func TestGetProfiles(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profiles")
	profiles, err := c.GetProfiles()
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestGetProfile(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-profile")

	// Act
	profile, err := c.GetProfile("Ubuntu-20.04-x86_64", false, false)

	// Assert
	FailOnError(t, err)
	if profile.Name != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong profile returned.")
	}
}

func TestDeleteProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-profile")
	err := c.DeleteProfile("test")
	FailOnError(t, err)
}

func TestDeleteProfileRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-profile")
	err := c.DeleteProfileRecursive("test", false)
	FailOnError(t, err)
}

func TestListProfileNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-profile")
	profiles, err := c.ListProfileNames()
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestGetProfilesSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profiles-since")
	profiles, err := c.GetProfilesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestFindProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-profile")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	profiles, err := c.FindProfile(criteria)
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestFindProfileNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-profile-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	profiles, err := c.FindProfileNames(criteria)
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestSaveProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-profile")
	err := c.SaveProfile("profile::testprof", true, true, "bypass")
	FailOnError(t, err)
}

func TestCopyProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-profile")
	err := c.CopyProfile("profile::testprof", "testprof2")
	FailOnError(t, err)
}

func TestRenameProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-profile")
	err := c.RenameProfile("profile::testprof2", "testprof1")
	FailOnError(t, err)
}

func TestGetProfileHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-handle")
	res, err := c.GetProfileHandle("testprof")
	FailOnError(t, err)

	if res != "profile::testprof" {
		t.Error("Wrong object id returned.")
	}
}

func TestGetValidProfileBootLoaders(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-valid-profile-boot-loaders")
	res, err := c.GetValidProfileBootLoaders("Ubuntu-20.04-x86_64")
	FailOnError(t, err)

	if len(res) < 1 {
		t.Error("Expected at least one boot loader.")
	}
}

func TestGetProfileAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-as-rendered")
	res, err := c.GetProfileAsRendered("Ubuntu-20.04-x86_64")
	FailOnError(t, err)

	if res["name"] != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong profile name returned: %v", res["name"])
	}
}

func TestNewSubprofile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "new-subprofile")
	res, err := c.NewSubprofile()
	FailOnError(t, err)
	if res == "" {
		t.Error("Expected a non-empty object ID from NewSubprofile.")
	}
}
