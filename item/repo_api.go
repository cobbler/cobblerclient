package item

import (
	"fmt"
	"reflect"
)

// Create creates a repo.
func (r *Repo) Create() (*Repo, error) {
	// Make sure a repo with the same name does not already exist
	if _, err := GetRepo(r.Client, r.Name.Get()); err == nil {
		return nil, fmt.Errorf("a Repo with the name %s already exists", r.Name)
	}

	result, err := r.Client.Call("new_repo", r.Client.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	item := reflect.ValueOf(r).Elem()
	if err := r.UpdateCobblerFields("repo", item, newID); err != nil {
		return nil, err
	}

	if _, err := r.Client.Call("save_repo", newID, r.Client.Token); err != nil {
		return nil, err
	}

	return GetRepo(r.Client, r.Name.Get())
}

// Update updates a single repo.
func (r *Repo) Update() error {
	item := reflect.ValueOf(r).Elem()
	id, err := r.Client.GetItemHandle("repo", r.Name.Get())
	if err != nil {
		return err
	}

	if err := r.UpdateCobblerFields("repo", item, id); err != nil {
		return err
	}

	if _, err := r.Client.Call("save_repo", id, r.Client.Token); err != nil {
		return err
	}

	return nil
}

// Delete deletes a single repo by its name.
func (r *Repo) Delete() error {
	_, err := r.Client.Call("remove_repo", r.Name, r.Client.Token)
	return err
}

// GetHandle gets the internal ID of a Cobbler item.
func (r *Repo) GetHandle() (string, error) {
	result, err := r.Client.Call("get_repo_handle", r.Name, r.Client.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// Copy is ...
func (r *Repo) Copy(newName string) error {
	_, err := r.Client.Call("copy_repo", r.Handle, newName, r.Client.Token)
	return err
}

// Rename is ...
func (r *Repo) Rename(objectId, newName string) error {
	_, err := r.Client.Call("rename_repo", r.Handle, newName, r.Client.Token)
	return err
}

// GetAsRendered is ...
func (r *Repo) GetAsRendered() error {
	_, err := r.Client.Call("get_repo_as_rendered", r.Name)
	return err
}

// Save is ...
func (r *Repo) Save(editmode string) error {
	_, err := r.Client.Call("save_repo", r.Handle, r.Client.Token, editmode)
	return err
}
