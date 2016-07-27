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
	ID        uint64
	PlaceID   uint64
	Name      string
	Type      machineType
	Brand     MachineBrand
	Model     string
	Selection int //จำนวน Column หรือช่องเก็บสินค้า
	LocNumber []int
				  //LocRow int  	//จำนวนแถว และคอลัมน์ไว้ทำ Schematic Profile  หน้าตู้
				  //LocCol int
}

