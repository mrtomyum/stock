package model

import (
	sys "github.com/mrtomyum/nava-sys/model"
)

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
	sys.Base
	LocID     uint64       `json:"loc_id"`
	Name      string       `json:"name"`
	Type      machineType  `json:"type"`
	Brand     MachineBrand `json:"brand"`
	Model     string       `json:"model"`
	Selection int          `json:"selection"`  //จำนวน Column หรือช่องเก็บสินค้า
	LocNumber []int        `json:"loc_number"` // Slice of number
											   //LocRow int  	//จำนวนแถว และคอลัมน์ไว้ทำ Schematic Profile  หน้าตู้
											   //LocCol int  //ควรจะเป็น 2 Dimension Array
}

type ColumnType int

const (
	CAN1 ColumnType = iota
	CAN2
	SPRING_S
	SPRING_M
	SPRING_L
)

type ColumnStatus int

const (
	OK ColumnStatus = iota
	FAIL
)

// เก็บยอด Counter ล่าสุดของแต่ละ column ในแต่ละ Machine
type MachineColumn struct {
	sys.Base
	MachineID uint64       `json:"machine_id"`
	ColumnNo  int          `json:"column_no"`
	SaleCount int          `json:"sale_count"`
	Type      ColumnType   `json:"type"`
	Status    ColumnStatus `json:"status"`
}

// Design this struct for data from VMC telemetry system.
//type SaleStatus int
//const (
//	COMPLETED SaleStatus = iota
//	INCOMPLETED
//)

type MachineSale struct {
	sys.Base
	MachineID uint64 `json:"machine_id"`
	ColumnNo  int    `json:"column_no"`
	ItemID    uint64 `json:"item_id"`
	//Status SaleStatus `json:"status"`
}
