package model

import (
	sys "github.com/mrtomyum/sys/model"
	"github.com/shopspring/decimal"
)

type Profile struct {
	sys.Base
	MachineID uint64 `json:"machine_id" db:"machine_id"` // one to one relate to Machine
	ItemRow   int `json:"planogram_row" db:"planogram_row"`
	ItemCol   int `json:"planogram_col" db:"planogram_col"`
	Matrix    [][]MachineColumn
                                                         // สร้าง MachineColumn Data Templat เพื่ออ่านทีละแถวแสดงผลทั้งภาพ และชื่อสินค้าและราคา
}

type ProfileItem struct {
	sys.Base
	ProfileId uint64          `json:"profile_id" db:"profile_id"`
	ColumnNo  int             `json:"column_no" db:"column_no"`
	ItemId    uint64          `json:"item_id" db:"item_id"`
	Price     decimal.Decimal `json:"price" db:"price"`
}
