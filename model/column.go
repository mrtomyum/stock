package model

import (
	"github.com/shopspring/decimal"
	"log"
	"github.com/jmoiron/sqlx"
)

// MachineColumn เก็บยอด Counter ล่าสุดของแต่ละ column ในแต่ละ Machine
type MachineColumn struct {
	Base
	MachineID   uint64          `json:"machine_id" db:"machine_id"`
	ColumnNo    int             `json:"column_no" db:"column_no"`
	ItemId      *uint64          `json:"item_id" db:"item_id"`
	Price       decimal.Decimal `json:"price"`
	LastCounter int             `json:"last_counter" db:"last_counter"`
	CurrCounter int             `json:"curr_counter" db:"curr_counter"`
	Size        ColumnSize      `json:"size"`
	Status      ColumnStatus    `json:"status"`
	MaxStock    int            `json:"max_stock" db:"max_stock"` // จำนวนสต๊อคสูงสุดในแต่ละคอลัมน์ แปรผันตาม Size ของสินค้าที่นำมาใส่
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

func (m *Machine) ColumnExist(db *sqlx.DB) (bool) {
	sql := `SELECT * FROM machine_column WHERE machine_id = ?`
	rows, err := db.Queryx(sql, m.Id)
	if err != nil {
		return false
	}
	if rows.Next() {
		return false
	}
	return true
}

func (mc *MachineColumn) Update(db *sqlx.DB) error {
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
	res, err := db.Exec(sql,
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

// ChangeItem() จะถูกเรียกเฉพาะตอน mc.Fulfill() ที่มีการเปลี่ยนสินค้าใหม่
func (mc *MachineColumn) ChangeItem(item Item) error {
	mc.ItemId = &item.Id
	// mc.MaxStock = ต้องทำตารางบันทึกผลทดสอบการหยอดสินค้าเข้าเก็บในตู้ทีละ Item ซึ่งต้องใช้เวลามาก ดังนั้นใส่ค่าประมาณการค่าเดียวไปก่อนคือ 30 ชิ้น
	mc.MaxStock = 35
	return nil
}