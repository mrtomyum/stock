package test

import (
	"testing"
	"github.com/mrtomyum/stock/model"
)

func TestNewColumn(t *testing.T) {
	columns := []*model.MachineColumn{}
	m := model.Machine{
		Type:      model.CAN,
		Selection: 30,
		Sub:       columns,
	}
	err := m.NewColumn(m.Selection)
	if err != nil {
		t.Error(err)
	}
	t.Log("ok")
}
