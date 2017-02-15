package test

import (
	"testing"
	"github.com/mrtomyum/stock/model"
)

func TestNewUser(t *testing.T) {
	// Setup
	for _, u := range mockUsers {
		err := u.New()
		if err != nil {
			t.Error(err)
		}
		t.Logf("Success Create mock user:", u)
	}
	// Tear down
	err := model.ResetTable("user")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Success truncate table: User")
}
