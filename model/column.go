package model

import (
	"github.com/shopspring/decimal"
	"log"
	"errors"
)

// MachineColumn เก็บยอด Counter ล่าสุดของแต่ละ column ในแต่ละ Machine
type MachineColumn struct {
	Base
	MachineID   uint64          `json:"machine_id" db:"machine_id"`
	ColumnNo    int             `json:"column_no" db:"column_no"`
	ItemId      uint64          `json:"item_id" db:"item_id"`
	Price       decimal.Decimal `json:"price"`
	LastCounter int             `json:"last_counter" db:"last_counter"`
	CurrCounter int             `json:"curr_counter" db:"curr_counter"`
	Size        ColumnSize      `json:"size"`
	Status      ColumnStatus    `json:"status"`
}

type ColumnSize int

const (
	NO_SIZE     ColumnSize = iota //สินค้าไม่มีตัวตน หรือต้องส่งข้อมูลสั่งขายไปยังระบบอื่น
	S
	L
	SPRING_5MM
	SPRING_10MM
	SPRING_15MM
)

type ColumnStatus int

const (
	COL_ERROR ColumnStatus = iota
	COL_OK
)

func (m *Machine) ColumnExist() (bool) {
	sql := `SELECT * FROM machine_column WHERE machine_id = ?`
	rows, err := DB.Queryx(sql, m.Id)
	if err != nil {
		return false
	}
	if rows.Next() {
		return false
	}
	return true
}

func (mc *MachineColumn) Update() error {
	log.Println("call model.MachineColumn.Update()")

	sql := `
	UPDATE machine_column
	SET
		machine_id 	= ?,
		column_no = ?,
		item_id = ?,
		price = ?,
		last_counter = ?,
		curr_counter = ?
	WHERE column_no = ?
	`
	res, err := DB.Exec(sql,
		mc.MachineID,
		mc.ColumnNo,
		mc.ItemId,
		mc.Price,
		mc.LastCounter,
		mc.CurrCounter,
		mc.ColumnNo,
	)
	if err != nil {
		return err
	}

	//var updatedMC MachineColumn
	id, _ := res.LastInsertId()
	log.Println("Insert MachineColumn_ID = ", id)
	return nil
}

func (m *Machine) InitMachineColumn() {
	// Select all MachineColumn from this Machine
	// Choose Delete old data or Create only missing one?
	//
}

func (m *Machine) GetMachineColumn(columnNo int) (*MachineColumn, error) {
	mc := new(MachineColumn)
	sql := `SELECT * FROM machine_column WHERE machine_id = ? AND column_no = ? LIMIT 1`
	err := DB.Get(mc, sql, m.Id, columnNo)
	if err != nil {
		return nil, errors.New("Wrong column number in this machine:" + err.Error())
	}
	return mc, nil
}

// NewColumn เพิ่มคอลัมน์ให้ครบตามจำนวน Selection ที่กำหนด ระวัง!! ถ้า Machine มี column ใดๆอยู่จะ Error ต้องลบ Column เดิมทิ้งก่อน
func (m *Machine) NewColumn(selection int) error {
	// ตรวจสอบก่อนว่ามี MachineColumn อยู่หรือไม่?
	//if m.ColumnExist() {
	if rowExists("SELECT * FROM machine_column WHERE machine_id = ?", m.Id) {
		return errors.New("ตู้นี้มี MachineColumn เดิมอยู่ กรุณาลบข้อมูลเดิมทิ้งก่อน")
	}
	sql := `INSERT INTO machine_column(
		machine_id,
		column_no
		) VALUES(?,?)
	`
	for col := 1; col == selection; col++ {
		res, err := DB.Exec(sql,
			m.Id,
			col,
		)
		if err != nil {
			return err
		}
		number, _ := res.RowsAffected()
		log.Println("Inserted", number, "row in MachineColumn so far", col, " from", selection)
	}
	return nil
}
