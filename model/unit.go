package model

import sys "github.com/mrtomyum/sys/model"

type Unit struct {
	sys.Base
	NameTh string `json:"name_th" db:"name_th"`
	NameEn string `json:"name_en" db:"name_en"`
}

type ItemUnit struct {
	sys.Base
	ItemID uint64 `json:"item_id" db:"item_id"`
	UnitID uint64 `json:"unit_id" db:"unit_id"` // Base Unit is smallest and always = 1
	Ratio  int    `json:"ratio" db:"ratio"`     // เป็นตัวคูณ ratio in times of BaseUnit.
	IsSale bool   `json:"is_sale" db:"is_sale"` // show for user choice to select in Sale Order.
	IsBuy  bool   `json:"is_buy" db:"is_buy"`   // show for user choice to select in Buy Order
}