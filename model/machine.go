package model

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/guregu/null"
	"database/sql"
)

type Machine struct {
	Base
	LocId        uint64           `json:"loc_id" db:"loc_id"`
	Code         string           `json:"code"`
	Type         machineType      `json:"type"`
	Brand        machineBrand     `json:"brand"`
	ProfileId    uint64           `json:"profile_id" db:"profile_id"` // ref. parent MachineId
	SerialNumber null.String      `json:"serial_number" db:"serial_number"`
	Selection    int              `json:"selection"` //จำนวน Column หรือช่องเก็บสินค้า
	PlaceId      uint64           `json:"place_id" db:"place_id"`
	Status       MachineStatus    `json:"status"`
	Note         null.String      `json:"note"`
	IsProfile    bool             `json:"is_profile"` // Profile template for Initialize New Machine each column data such as item, price
	PriceLevel   int              `json:"price_level" db:"price_level"`
	Sub          []*MachineColumn `json:"sub"`
}

type machineType uint8

const (
	NO_TYPE          machineType = iota
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
	NO_BRAND      machineBrand = iota
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

type MachineStatus int

const (
	NO_STATUS MachineStatus = iota
	OFFLINE
	ONLINE
	ALARM
)

// Transaction row Batch data received from mobile app daily.

// เก็บ Transaction ที่มีความผิดปกติทั้งหมด เช่น  ข้อมูลที่ส่งมาหา Column ไม่เจอ ไปจนถึง Error ที่แจ้งจาก Machine
type MachineErrType int

const (
	X                 MachineErrType = iota // UNIDENTIFIED ERROR
	COLUMN_NOT_FOUND
	COUNTER_OVER_SALE
)

type MachineErrLog struct {
	Base
	MachineID uint64         `json:"machine_id"`
	ColumnNo  int            `json:"column_no"`
	Type      MachineErrType `json:"type"`
	Message   string         `json:"message"`
}

func MachineBatchSaleIsErr() bool {
	// Check Error Code in Transaction
	return false
}

func (m *Machine) GetAll() ([]*Machine, error) {
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

// rowExists ใช้ตรวจสอบว่าตารางที่กำลังจะ insert มีข้อมูลอยู่หรือไม่? https://snippets.aktagon.com/snippets/756-checking-if-a-row-exists-in-go-database-sql-and-sqlx-
func rowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("error checking if row exists '%s' %v", args, err)
	}
	log.Println("exists =", exists)
	return exists
}

func (m *Machine) New() (*Machine, error) {
	// rowExists ตรวจสอบรหัสตู้ m.Code ว่าซ้ำอยู่หรือไม่?
	log.Println("call model.Machine.New()")
	if rowExists("SELECT * FROM machine WHERE code = ?", m.Code) {
		return nil, errors.New("มี Machine นี้อยู่แล้วใน Database กรุณาลบของเดิมทิ้งก่อนเพิ่มใหม่")
	}
	sql := `INSERT INTO machine(
		loc_id,
		code,
		type,
		brand,
		profile_id,
		serial_number,
		selection,
		place_id
		) VALUES(?,?,?,?,?,?,?,?)`
	res, err := db.Exec(sql,
		m.LocId,
		m.Code,
		m.Type,
		m.Brand,
		m.ProfileId,
		m.SerialNumber,
		m.Selection,
		m.PlaceId,
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
	_, err = newMachine.InitMachineColumn()
	if err != nil {
		return nil, err
	}
	log.Println("New Machine:", newMachine)
	return &newMachine, nil
}

func (m *Machine) Get() (*Machine, error) {
	log.Println("call model.Machine.Get()")
	sql := `SELECT * FROM machine WHERE id = ? AND deleted IS NULL`
	err := db.Get(m, sql, m.Id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Machine) GetColumns() ([]*MachineColumn, error) {
	log.Println("call model.Machine.Columns()")
	var mcs []*MachineColumn
	sql := `SELECT * FROM machine_column WHERE machine_id = ?`
	err := db.Select(&mcs, sql, m.Id)
	if err != nil {
		return nil, err
	}
	return mcs, nil
}

func (m *Machine) Init(template *Machine) error {
	// Check If this Machine is blank no data
	// Copy MachineColumn From Template to This Machine
	//for i:=1; i=m.Selection; i++ {
	//
	//}
	return nil
}

func (m *Machine) GetTemplate() ([]*Machine, error) {
	var templates []*Machine
	sql := `SELECT * FROM machine
	WHERE is_template = true`
	err := db.Select(templates, sql)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (m *Machine) InitMachineColumn() (count int, err error) {
	// Create only missing one.
	sql1 := `INSERT INTO machine_column(machine_id, column_no) VALUES(?, ?)`
	sql2 := `SELECT * FROM machine_column WHERE id =?`
	fmt.Println("m.Selection =", m.Selection)

	for n := 1; n <= m.Selection; n++ {
		fmt.Println("Loop times= ", n)
		if rowExists("SELECT * FROM machine_column WHERE machine_id = ? AND column_no = ?", m.Id, n) {
			fmt.Printf("Machine: %v Column: %v exist.\n", m.Id, n)
			continue
		}
		fmt.Printf("Machine: %v Column: %v Not Exist.\n", m.Id, n)
		res, err := db.Exec(sql1, m.Id, n)
		if err != nil {
			return 0, err
		}
		id, _ := res.LastInsertId()
		fmt.Println("LastInsertId:", id)
		count++
		col := new(MachineColumn)
		err = db.Get(col, sql2, id)
		if err != nil {
			return 0, err
		}
		m.Sub = append(m.Sub, col)
		fmt.Println("count=", count)
	}
	fmt.Println("New MachineColumn initiated = ", count)
	return count, nil
}

func (m *Machine) GetMachineColumn(columnNo int) (*MachineColumn, error) {
	column := new(MachineColumn)
	sql := `SELECT * FROM machine_column WHERE machine_id = ? AND column_no = ? LIMIT 1`
	err := db.Get(column, sql, m.Id, columnNo)
	if err != nil {
		return nil, errors.New("Wrong column number in this machine:" + err.Error())
	}
	return column, nil
}

// NewColumn เพิ่มคอลัมน์ให้ครบตามจำนวน Selection ที่กำหนด ระวัง!! ถ้า Machine มี column ใดๆอยู่จะ Error ต้องลบ Column เดิมทิ้งก่อน
func (m *Machine) NewColumn(selection int) error {
	sql := `INSERT INTO machine_column(
		machine_id,
		column_no,
		max_qty
		) VALUES(?,?,?)
	`
	for col := 1; col == selection; col++ {
		// ตรวจสอบก่อนว่ามี ColumnNo ซ้ำอยู่หรือไม่?
		if rowExists("SELECT * FROM machine_column WHERE machine_id = ? AND column_no = ?", m.Id, col) {
			continue
		}
		// Todo: max_qty ควร Refactor โดยบันทึกตอน Fulfill เปลี่ยนสินค้าในคอลัมน์
		res, err := db.Exec(sql, m.Id, col, 30)
		if err != nil {
			return err
		}
		number, _ := res.RowsAffected()
		log.Println("Inserted", number, "row in MachineColumn so far", col, " from", selection)
	}
	return nil
}

func (m *Machine) Update() (*Machine, error) {
	// rowExists ตรวจสอบรหัสตู้ m.Code ว่าซ้ำอยู่หรือไม่?
	log.Println("call model.Machine.New()")
	if !rowExists("SELECT * FROM machine WHERE code = ?", m.Code) {
		return nil, errors.New("ไม่มี Machine รหัสนี้ นี้อยู่ใน Database กรุณาเพิ่มใหม่")
	}
	sql := `UPDATE machine(
		loc_id,
		code,
		type,
		brand,
		profile_id,
		serial_number,
		selection,
		place_id
		) VALUES(?,?,?,?,?,?,?,?)`
	res, err := db.Exec(sql,
		m.LocId,
		m.Code,
		m.Type,
		m.Brand,
		m.ProfileId,
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

func GetMachineIdFromCode(code string) (uint64, error) {
	sql1 := `SELECT id FROM machine WHERE code = ?`
	var id uint64
	err := db.Get(&id, sql1, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//func (m *Machine) UpdateColumnCounter(columnNo int, counter int) error {
//	mc, err := m.GetMachineColumn(columnNo)
//	if err != nil {
//		return err
//	}
//	sql := `
//			UPDATE machine_column
//			SET last_counter = ?, curr_counter = ?
//			WHERE machine_id = ?
//			AND column_no = ?
//	//		`
//	_, err = DB.Exec(sql,
//		mc.CurrCounter,
//		counter,
//		m.ID,
//		columnNo,
//	)
//	if err != nil {
//		log.Println("Error>> Exec() machine_column = ", err)
//		return err
//	}
//	log.Println("Update MachineColumn 'MachineID':", m.ID, "ColumnNo:", columnNo)
//	log.Println("Pass>>3.tx.Exec() UPDATE machine_column")
//	return nil
//}
