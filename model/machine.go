package model

type machineType int

const (
	CAN machineType = 1 + iota
	CUP
	CUP_FRESH_COFFEE
	CUP_NOODLE
	SEE_THROUGH
)

type MachineBrand int

const (
	NATIONAL MachineBrand = 1 + iota
	SANDEN
	FUJI_ELECTRIC
	CIRBOX
)

type Machine struct {
	ID        uint64 `json:"id"`
	LocID     uint64 `json:"loc_id"`
	Name      string `json:"name"`
	Type      machineType `json:"type"`
	Brand     MachineBrand `json:"brand"`
	Model     string `json:"model"`
	Selection int `json:"selection"`    //จำนวน Column หรือช่องเก็บสินค้า
	LocNumber []int `json:"loc_number"` // Slice of number
										//LocRow int  	//จำนวนแถว และคอลัมน์ไว้ทำ Schematic Profile  หน้าตู้
										//LocCol int  //ควรจะเป็น 2 Dimension Array
}

