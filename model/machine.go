package model

import (
	"encoding/json"
	sys "github.com/mrtomyum/nava-sys/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"github.com/guregu/null"
)

type machineType uint8

const (
	CAN25 machineType = 1 + iota
	CAN30
	CUP_HOT_COLD
	CUP_FRESH_COFFEE
	CUP_NOODLE
	SEE_THROUGH
)

func (t machineType) MarshalJSON() ([]byte, error) {
	typeStr, ok := map[machineType]string{
		CAN25:            "CAN25",
		CAN30:            "CAN30",
		CUP_HOT_COLD:     "CUP_HOT_COLD",
		CUP_FRESH_COFFEE: "CUP_FRESH_COFFEE",
		CUP_NOODLE:       "CUP_NOODLE",
		SEE_THROUGH:      "SEE_THROUGH",
	}[t]
	if !ok {
		return nil, fmt.Errorf("invalid Machine Type value %v", t)
	}
	return json.Marshal(typeStr)
}

// Todo: implement UnmarshalJSON for MachineType
//func (t machineType) UnmarshalJSON(data []byte) error {
//
//}

type machineBrand int

const (
	NATIONAL machineBrand = 1 + iota
	SANDEN
	FUJI_ELECTRIC
	CIRBOX
)

func (b machineBrand) MarshalJSON() ([]byte, error) {
	brandStr, ok := map[machineBrand]string{
		NATIONAL:      "NATIONAL",
		SANDEN:        "SANDEN",
		FUJI_ELECTRIC: "FUJI_ELECTRIC",
		CIRBOX:        "CIRBOX",
	}[b]
	if !ok {
		return nil, fmt.Errorf("invalid Brand value %v", b)
	}
	return json.Marshal(brandStr)
}

type Machine struct {
	sys.Base
	LocID        uint64       `json:"loc_id" db:"loc_id"`
	Code         string       `json:"code"`
	Type         machineType  `json:"type"`
	Brand        machineBrand `json:"brand"`
	SerialNumber null.String  `json:"serial_number" db:"serial_number"`
	//จำนวน Column หรือช่องเก็บสินค้า
	Selection    int          `json:"selection"`
	//LocRow int  	//จำนวนแถว และคอลัมน์ไว้ทำ Schematic Profile  หน้าตู้
	//LocCol int  //ควรจะเป็น 2 Dimension Array
}

type ColumnType int

const (
	FREE ColumnType = iota //สินค้าไม่มีตัวตน หรือต้องส่งข้อมูลสั่งขายไปยังระบบอื่น
	TICKET                   // สินค้าที่ต้องพิมพ์ตั๋ว
	CAN_S                    // กระป๋องหรือขวดสั้น
	CAN_L                    // กระป๋องหรือขวดยาว
	SPRING_S
	SPRING_M
	SPRING_L
)

type ColumnStatus int

const (
	OK ColumnStatus = iota
	FAIL
)

// MachineColumn เก็บยอด Counter ล่าสุดของแต่ละ column ในแต่ละ Machine
type MachineColumn struct {
	sys.Base
	MachineID   uint64       `json:"machine_id"`
	ColumnNo    int          `json:"column_no"`
	CurrCounter int          `json:"curr_counter" db:"curr_counter"`
	Type        ColumnType   `json:"type"`
	Status      ColumnStatus `json:"status"`
}


// Transaction row Batch data received from mobile app daily.

// เก็บ Transaction ที่มีความผิดปกติทั้งหมด เช่น  ข้อมูลที่ส่งมาหา Column ไม่เจอ ไปจนถึง Error ที่แจ้งจาก Machine
type MachineErrType int

const (
	X MachineErrType = iota // UNIDENTIFIED ERROR
	COLUMN_NOT_FOUND
	COUNTER_OVER_SALE
)

type MachineErrLog struct {
	sys.Base
	MachineID uint64         `json:"machine_id"`
	ColumnNo  int            `json:"column_no"`
	Type      MachineErrType `json:"type"`
	Message   string         `json:"message"`
}

func MachineBatchSaleIsErr() bool {
	// Check Error Code in Transaction
	return false
}

func (m *Machine) All(db *sqlx.DB) ([]*Machine, error) {
	log.Println("call Machine.All()")
	machines := []*Machine{}
	sql := `SELECT * FROM machine`
	err := db.Select(&machines, sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(machines)
	return machines, nil

}