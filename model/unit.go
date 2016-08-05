package model

type Unit struct {
	ID uint64 `json:"id"`
	TH string `json:"th"`
	EN string `json:"en"`
}

type ItemUnit struct {
	ID     uint64 `json:"id"`
	ItemID uint64 `json:"itemID" db:"item_id"`
	UnitID uint64 `json:"unitID" db:"unit_id"` // Base Unit is smallest and always = 1
	Ratio  int    `json:"ratio" db:"ratio"`    // เป็นตัวคูณ ratio in times of BaseUnit.
	IsSale bool   `json:"isSale" db:"is_sale"` // show for user choice to select in Sale Order.
	IsBuy  bool   `json:"isBuy" db:"is_buy"`   // show for user choice to select in Buy Order
}