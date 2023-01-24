package item

import (
	"fmt"
	"github.com/cobbler/cobblerclient/client"
	"github.com/cobbler/cobblerclient/internal/item/raw"

	clientInternals "github.com/cobbler/cobblerclient/internal/client"
)

// GetRepos returns all repos in Cobbler.
func GetRepos(c client.Client) ([]*Repo, error) {
	var repos []*Repo

	result, err := c.Call("get_repos", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	for _, r := range result.([]interface{}) {
		var repo = BuildRepo(c)
		var rawRepo raw.Repo
		decodedResult, err := clientInternals.DecodeCobblerItem(r, &rawRepo)
		if err != nil {
			return nil, err
		}
		repo.raw = decodedResult.(*raw.Repo)

		repos = append(repos, &repo)
	}

	return repos, nil
}

// ListRepoNames returns the names of all repositories that exist in Cobbler.
func ListRepoNames(c client.Client) ([]string, error) {
	return c.GetItemNames("repo")
}

// GetRepo returns a single repo obtained by its name.
func GetRepo(c client.Client, name string) (*Repo, error) {
	var repo = BuildRepo(c)
	var rawRepo raw.Repo

	result, err := c.Call("get_repo", name, c.Token)
	if result == "~" {
		return nil, fmt.Errorf("repo %s not found", name)
	}

	if err != nil {
		return nil, err
	}

	decodedResult, err := clientInternals.DecodeCobblerItem(result, &rawRepo)
	if err != nil {
		return nil, err
	}

	repo.raw = decodedResult.(*raw.Repo)
	refreshRepoPointers(&repo)
	return &repo, nil
}

// FindRepo is ...
func FindRepo(c client.Client) error {
	_, err := c.Call("find_repo")
	return err
}

// GetReposSince is ...
func GetReposSince(c client.Client) error {
	_, err := c.Call("get_repos_since")
	return err
}
