package mock

import (
	"fmt"

	"github.com/dominik-zeglen/inkster/core"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (adapter Adapter) findDirectory(id string) (int, error) {
	for index := range directories {
		if directories[index].ID == id {
			return index, nil
		}
	}
	return 0, mgo.ErrNotFound
}

// AddDirectory puts directory in a in-memory array
func (adapter Adapter) AddDirectory(directory core.Directory) (core.Directory, error) {
	errors := directory.Validate()
	if len(errors) > 0 {
		return core.Directory{}, core.ErrNotValidated
	}

	if directory.ID == "" {
		directory.ID = bson.NewObjectId().String()
	} else {
		_, err := adapter.findDirectory(directory.ID)
		if err == nil {
			return core.Directory{}, fmt.Errorf("Could not add directory with ID: %s; directory already exists", directory.ID)
		}
	}
	directory.CreatedAt = adapter.GetCurrentTime()
	directory.UpdatedAt = adapter.GetCurrentTime()

	directories = append(directories, directory)
	return directory, nil
}

// GetDirectory gets directory from a in-memory array
func (adapter Adapter) GetDirectory(id string) (core.Directory, error) {
	index, err := adapter.findDirectory(id)
	if err != nil {
		return core.Directory{}, err
	}
	return directories[index], nil
}

// GetDirectoryList gets all directories from a in-memory array
func (adapter Adapter) GetDirectoryList() ([]core.Directory, error) {
	return directories, nil
}

// GetRootDirectoryList gets only directories from a in-memory array that don't have parent
func (adapter Adapter) GetRootDirectoryList() ([]core.Directory, error) {
	result := []core.Directory{}
	for k := range directories {
		if directories[k].ParentID == "" {
			result = append(result, directories[k])
		}
	}
	return result, nil
}

// GetDirectoryChildrenList gets directories from a in-memory array which
// ParentID equals to function id parameter
func (adapter Adapter) GetDirectoryChildrenList(id string) ([]core.Directory, error) {
	result := []core.Directory{}
	for k := range directories {
		if directories[k].ParentID == id {
			result = append(result, directories[k])
		}
	}
	return result, nil
}

// UpdateDirectory updates directory with given properties
func (adapter Adapter) UpdateDirectory(id string, data core.DirectoryInput) error {
	errors := core.ValidateModel(data)
	if len(errors) > 0 {
		return core.ErrNotValidated
	}

	index, err := adapter.findDirectory(id)
	if err != nil {
		return err
	}
	if data.Name != nil {
		directories[index].Name = *data.Name
	}
	if data.ParentID != nil {
		directories[index].ParentID = *data.ParentID
	}
	if data.IsPublished != nil {
		directories[index].IsPublished = *data.IsPublished
	}
	directories[index].UpdatedAt = adapter.GetCurrentTime()
	return nil
}

// RemoveDirectory removes directory from a in-memory array
func (adapter Adapter) RemoveDirectory(id string) error {
	for index := range directories {
		if directories[index].ID == id {
			directories = append(directories[:index], directories[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not remove directory with ID: %s. Directory does not exist", id)
}
