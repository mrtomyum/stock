package model

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/nava-sys/model"
	"golang.org/x/text/currency"
)

type machineType uint8

const (
	NO_TYPE machineType = iota
	CAN
	CUP_HOT_COLD
	CUP_FRESH_COFFEE
	CUP_NOODLE
	SPRING
	SEE_THROUGH
	TICKET
)

func (t machineType) MarshalJSON() ([]byte, error) {
	typeStr, ok := map[machineType]string{
		NO_TYPE:          "NO_TYPE",
		CAN:              "CAN",
		CUP_HOT_COLD:     "CUP_HOT_COLD",
		CUP_FRESH_COFFEE: "CUP_FRESH_COFFEE",
		CUP_NOODLE:       "CUP_NOODLE",
		SPRING:           "SPRING",
		SEE_THROUGH:      "SEE_THROUGH",
		TICKET:           "TICKET",
	}[t]
	if !ok {
		return nil, fmt.Errorf("invalid Machine Type value %v", t)
	}
	return json.Marshal(typeStr)
}

type machineBrand int

const (
	NO_BRAND machineBrand = iota
	NATIONAL
	SANDEN
	FUJI_ELECTRIC
	CIRBOX
)

func (b machineBrand) MarshalJSON() ([]byte, error) {
	brandStr, ok := map[machineBrand]string{
		NO_BRAND:      "NO_BRAND",
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
	ProfileID    uint64       `json:"profile_id" db:"profile_id"`
	SerialNumber null.String  `json:"serial_number" db:"serial_number"`
	Selection    int          `json:"selection"` //จำนวน Column หรือช่องเก็บสินค้า
	ClientID     uint64       `json:"client_id" db:"client_id"`
}

type ColumnSize int

const (
	NO_SIZE ColumnSize = iota //สินค้าไม่มีตัวตน หรือต้องส่งข้อมูลสั่งขายไปยังระบบอื่น
	S
	L
	SPRING_5MM
	SPRING_10MM
	SPRING_15MM
)

type ColumnStatus int

const (
	OK ColumnStatus = iota
	FAIL
)

// MachineColumn เก็บยอด Counter ล่าสุดของแต่ละ column ในแต่ละ Machine
type MachineColumn struct {
	sys.Base
	MachineID   uint64          `json:"machine_id" db:"machine_id"`
	Number      int             `json:"column_no" db:"column_no"`
	ItemId      uint64          `json:"item_id" db:"item_id"`
	Price       currency.Amount `json:"price"`
	LastCounter int             `json:"last_counter" db:"last_counter"`
	CurrCounter int             `json:"curr_counter" db:"curr_counter"`
	Size        ColumnSize      `json:"size"`
	Status      ColumnStatus    `json:"status"`
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
	log.Info(log.Fields{"func": "Machine.All()"})
	machines := []*Machine{}
	sql := `SELECT * FROM machine`
	err := db.Select(&machines, sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("All Machine:", machines)
	return machines, nil

}

func (m *Machine) New(db *sqlx.DB) (*Machine, error) {
	log.Println("call Machine.New()")
	sql := `INSERT INTO machine(
		loc_id,
		code,
		type,
		brand,
		profile_id,
		serial_number,
		selection
		) VALUES(?,?,?,?,?,?,?)`
	res, err := db.Exec(sql,
		m.LocID,
		m.Code,
		m.Type,
		m.Brand,
		m.ProfileID,
		m.SerialNumber,
		m.Selection,
	)
	if err != nil {
		return nil, err
	}
	var newMachine Machine
	sql = `SELECT * FROM machine WHERE id = ?`
	id, _ := res.LastInsertId()
	err = db.Get(&newMachine, sql, uint64(id))
	if err != nil {
		return nil, err
	}
	log.Println("New Machine:", newMachine)
	return &newMachine, nil
}
