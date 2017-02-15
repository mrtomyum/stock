package test

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	for _, u := range mockUsers {
		err := u.New()
		if err != nil {
			t.Error(err)
		}
		t.Logf("Success Create new user:", u)
	}
}
