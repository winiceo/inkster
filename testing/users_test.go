package testing

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/dominik-zeglen/inkster/core"
)

func testUsers(t *testing.T) {
	dummyID := 0
	t.Run("Test setters", func(t *testing.T) {
		t.Run("Add user", func(t *testing.T) {
			defer resetDatabase()
			user := core.User{
				Email:  "newuser@example.com",
				Active: false,
			}
			user.CreatePassword("passwd")
			result, err := dataSource.AddUser(user)
			if err != nil {
				t.Fatal(err)
			}

			result.ID = dummyID
			result.Password = []byte("")
			result.Salt = []byte("")
			data, err := ToJSON(result)
			if err != nil {
				t.Error(err)
			}
			cupaloy.SnapshotT(t, data)
		})
		t.Run("Add user without email", func(t *testing.T) {
			user := core.User{
				Active: false,
			}
			_, err := dataSource.AddUser(user)
			if err == nil {
				t.Error(ErrNoError)
			}
		})
		t.Run("Update user's password", func(t *testing.T) {
			defer resetDatabase()

			passwd := "thisisnewpassword"
			result, err := dataSource.UpdateUser(Users[0].ID, core.UserInput{
				Password: &passwd,
			})
			if err != nil {
				t.Error(err)
			}

			if !result.AuthPassword(passwd) {
				t.Fatal()
			}

			result.Password = []byte("")
			result.Salt = []byte("")
			data, err := ToJSON(result)
			if err != nil {
				t.Error(err)
			}

			cupaloy.SnapshotT(t, data)
		})
		t.Run("Update user's name", func(t *testing.T) {
			defer resetDatabase()

			email := "updated.email@example.com"
			result, err := dataSource.UpdateUser(Users[0].ID, core.UserInput{
				Email: &email,
			})
			if err != nil {
				t.Error(err)
			}

			result.Password = []byte("")
			result.Salt = []byte("")
			data, err := ToJSON(result)
			if err != nil {
				t.Error(err)
			}

			cupaloy.SnapshotT(t, data)
		})
		t.Run("Remove user", func(t *testing.T) {
			defer resetDatabase()
			err := dataSource.RemoveUser(Users[3].ID)
			if err != nil {
				t.Error(err)
			}
		})
	})
	t.Run("Test getters", func(t *testing.T) {
		t.Run("Get user", func(t *testing.T) {
			result, err := dataSource.GetUser(Users[0].ID)
			if err != nil {
				t.Error(err)
			}
			result.Password = []byte("")
			result.Salt = []byte("")
			data, err := ToJSON(result)
			if err != nil {
				t.Error(err)
			}
			cupaloy.SnapshotT(t, data)
		})
		t.Run("Get user that does not exist", func(t *testing.T) {
			_, err := dataSource.GetUser(99)
			if err == nil {
				t.Error(ErrNoError)
			}
		})
		t.Run("Get all users", func(t *testing.T) {
			result, _ := dataSource.GetUserList()
			for index := range result {
				result[index].Password = []byte("")
				result[index].Salt = []byte("")
			}
			data, err := ToJSON(result)
			if err != nil {
				t.Error(err)
			}
			cupaloy.SnapshotT(t, data)
		})
	})
}
