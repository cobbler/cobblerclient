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
	"fmt"
	"reflect"
	"time"
)

// Repo is a created repo.
// Get the fileds from cobbler/items/repo.py
type Repo struct {
	Item `mapstructure:",squash"`

	// These are internal fields and cannot be modified.
	TreeBuildTime string `mapstructure:"tree_build_time" cobbler:"noupdate"`

	AptComponents   []string          `mapstructure:"apt_components"`
	AptDists        []string          `mapstructure:"apt_dists"`
	Arch            string            `mapstructure:"arch"`
	Breed           string            `mapstructure:"breed"`
	CreateRepoFlags string            `mapstructure:"createrepo_flags"`
	Environment     map[string]string `mapstructure:"environment"`
	KeepUpdated     bool              `mapstructure:"keep_updated"`
	Mirror          string            `mapstructure:"mirror"`
	MirrorLocally   bool              `mapstructure:"mirror_locally"`
	MirrorType      string            `mapstructure:"mirror_type"`
	Priority        int               `mapstructure:"priority"`
	Proxy           string            `mapstructure:"proxy" cobbler:"newfield"`
	RsyncOpts       map[string]string `mapstructure:"rsyncopts"`
	RpmList         []string          `mapstructure:"rpm_list"`
	YumOpts         map[string]string `mapstructure:"yumopts"`
}

func convertRawRepo(name string, xmlrpcResult interface{}) (*Repo, error) {
	var repo Repo

	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("profile %s not found", name)
	}

	decodeResult, err := decodeCobblerItem(xmlrpcResult, &repo)
	if err != nil {
		return nil, err
	}

	return decodeResult.(*Repo), nil
}

func convertRawReposList(xmlrpcResult interface{}) ([]*Repo, error) {
	var repos []*Repo

	for _, r := range xmlrpcResult.([]interface{}) {
		repo, err := convertRawRepo("unknown", r)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

// GetRepos returns all repos in Cobbler.
func (c *Client) GetRepos() ([]*Repo, error) {
	result, err := c.Call("get_repos", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawReposList(result)
}

// GetRepo returns a single repo obtained by its name.
func (c *Client) GetRepo(name string) (*Repo, error) {
	result, err := c.Call("get_repo", name, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawRepo(name, result)
}

// CreateRepo creates a repo.
func (c *Client) CreateRepo(repo Repo) (*Repo, error) {
	// Make sure a repo with the same name does not already exist
	if _, err := c.GetRepo(repo.Name); err == nil {
		return nil, fmt.Errorf("a Repo with the name %s already exists", repo.Name)
	}

	result, err := c.Call("new_repo", c.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	item := reflect.ValueOf(&repo).Elem()
	if err := c.updateCobblerFields("repo", item, newID); err != nil {
		return nil, err
	}

	if err := c.SaveRepo(newID, "new"); err != nil {
		return nil, err
	}

	return c.GetRepo(repo.Name)
}

// UpdateRepo updates a single repo.
func (c *Client) UpdateRepo(repo *Repo) error {
	item := reflect.ValueOf(repo).Elem()
	id, err := c.GetItemHandle("repo", repo.Name)
	if err != nil {
		return err
	}

	if err := c.updateCobblerFields("repo", item, id); err != nil {
		return err
	}

	return c.SaveRepo(id, "bypass")
}

// SaveRepo saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveRepo(objectId, editmode string) error {
	_, err := c.Call("save_repo", objectId, c.Token, editmode)
	return err
}

// CopyRepo duplicates a given repository on the server with a new name.
func (c *Client) CopyRepo(objectId, newName string) error {
	_, err := c.Call("copy_repo", objectId, newName, c.Token)
	return err
}

// DeleteRepo deletes a single Repo by its name.
func (c *Client) DeleteRepo(name string) error {
	return c.DeleteRepoRecursive(name, false)
}

// DeleteRepoRecursive deletes a single Repo by its name with the option to do so recursively.
func (c *Client) DeleteRepoRecursive(name string, recursive bool) error {
	_, err := c.Call("remove_repo", name, c.Token, recursive)
	return err
}

// ListRepoNames returns the names of all repositories that exist in Cobbler.
func (c *Client) ListRepoNames() ([]string, error) {
	return c.GetItemNames("repo")
}

// FindRepo searches for one or more repositories by any of its attributes.
func (c *Client) FindRepo(criteria map[string]interface{}) ([]*Repo, error) {
	result, err := c.Call("find_repo", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}
	return convertRawReposList(result)
}

// FindRepoNames searches for one or more repositories by any of its attributes.
func (c *Client) FindRepoNames(criteria map[string]interface{}) ([]string, error) {
	resultUnmarshalled, err := c.Call("find_repo", criteria, false, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetReposSince returns all repositories which were created after the specified date.
func (c *Client) GetReposSince(mtime time.Time) ([]*Repo, error) {
	result, err := c.Call("get_repos_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}
	return convertRawReposList(result)
}

// RenameRepo renames a repository with a given object id.
func (c *Client) RenameRepo(objectId, newName string) error {
	_, err := c.Call("rename_repo", objectId, newName, c.Token)
	return err
}

// GetRepoHandle gets the internal ID of a Cobbler item.
func (c *Client) GetRepoHandle(name string) (string, error) {
	res, err := c.Call("get_repo_handle", name, c.Token)
	return returnString(res, err)
}
