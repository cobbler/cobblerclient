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

func TestGetProfiles(t *testing.T) {
	c := createStubHTTPClient(t, "get-profiles-req.xml", "get-profiles-res.xml")
	profiles, err := c.GetProfiles()
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestGetProfile(t *testing.T) {
	c := createStubHTTPClient(t, "get-profile-req.xml", "get-profile-res.xml")
	profile, err := c.GetProfile("Ubuntu-20.04-x86_64")
	FailOnError(t, err)

	if profile.Name != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong profile returned.")
	}
}

func TestDeleteProfile(t *testing.T) {
	c := createStubHTTPClient(t, "delete-profile-req.xml", "delete-profile-res.xml")
	err := c.DeleteProfile("test")
	FailOnError(t, err)
}

func TestDeleteProfileRecursive(t *testing.T) {
	c := createStubHTTPClient(t, "delete-profile-req.xml", "delete-profile-res.xml")
	err := c.DeleteProfileRecursive("test", false)
	FailOnError(t, err)
}

func TestListProfileNames(t *testing.T) {
	c := createStubHTTPClient(t, "get-item-names-profile-req.xml", "get-item-names-profile-res.xml")
	profiles, err := c.ListProfileNames()
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestGetProfilesSince(t *testing.T) {
	c := createStubHTTPClient(t, "get-profiles-since-req.xml", "get-profiles-since-res.xml")
	profiles, err := c.GetProfilesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestFindProfile(t *testing.T) {
	c := createStubHTTPClient(t, "find-profile-req.xml", "find-profile-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	profiles, err := c.FindProfile(criteria)
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestFindProfileNames(t *testing.T) {
	c := createStubHTTPClient(t, "find-profile-names-req.xml", "find-profile-names-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	profiles, err := c.FindProfileNames(criteria)
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestSaveProfile(t *testing.T) {
	c := createStubHTTPClient(t, "save-profile-req.xml", "save-profile-res.xml")
	err := c.SaveProfile("profile::testprof", "bypass")
	FailOnError(t, err)
}

func TestCopyProfile(t *testing.T) {
	c := createStubHTTPClient(t, "copy-profile-req.xml", "copy-profile-res.xml")
	err := c.CopyProfile("profile::testprof", "testprof2")
	FailOnError(t, err)
}

func TestRenameProfile(t *testing.T) {
	c := createStubHTTPClient(t, "rename-profile-req.xml", "rename-profile-res.xml")
	err := c.RenameProfile("profile::testprof2", "testprof1")
	FailOnError(t, err)
}

func TestGetProfileHandle(t *testing.T) {
	c := createStubHTTPClient(t, "get-profile-handle-req.xml", "get-profile-handle-res.xml")
	res, err := c.GetProfileHandle("testprof")
	FailOnError(t, err)

	if res != "profile::testprof" {
		t.Error("Wrong object id returned.")
	}
}
