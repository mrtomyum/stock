package model

import (
	"golang.org/x/text/currency"
	sys "github.com/mrtomyum/nava-sys/model"
)

type Profile struct {
	sys.Base
	Code        string      `json:"code"`
	MachineType machineType `json:"machine_type" db:"machine_type"`
	PriceLevel  int         `json:"price_level" db:"price_level"`
}

type ProfileItem struct {
	sys.Base
	ProfileID uint64          `json:"profile_id" db:"profile_id"`
	ColumnNo  int             `json:"column_no" db:"column_no"`
	ItemID    uint64          `json:"item_id" db:"item_id"`
	Price     currency.Amount `json:"price" db:"price"`
}

