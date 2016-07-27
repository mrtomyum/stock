package model

type Unit struct {
	ID uint64 `json:"id"`
	TH string `json:"th"`
	EN string `json:"en"`
}

type ItemUnit struct {
	ID     uint64 `json:"id"`
	ItemID uint64 `json:"Item_id"`
	UnitID uint64 `json:"unit_id"` // Base Unit is smallest and always = 1
	Ratio  int    `json:"ratio"`   // เป็นตัวคูณ ratio in times of BaseUnit.
	IsSale bool   `json:"is_sale"` // show for user choice to select in Sale Order.
	IsBuy  bool   `json:"is_buy"`  // show for user choice to select in Buy Order
}