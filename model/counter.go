package model

import (
	"github.com/shopspring/decimal"
	"log"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Counter struct {
	Base
	//RecDate    *time.Time `json:"rec_date" db:"rec_date"`
	RecDate    Date   `json:"rec_date" db:"rec_date"`
	MachineId  uint64 `json:"machine_id" db:"machine_id"`
	CounterSum int    `json:"counter_sum" db:"counter_sum"`
	Sub        []*CounterSub `json:"sub"`
}

type CounterSub struct {
	Base
	CounterId uint64            `json:"-" db:"counter_id"` // FK
	ColumnNo  int                `json:"column_no" db:"column_no"`
	Counter   int                `json:"counter" db:"counter"`
	// ฟิลด์หลังจากนี้ Method New, Update etc ควรจะกรอกให้
	ItemId *uint64          `json:"item_id" db:"item_id"` // Record as history data.
	Price  decimal.Decimal  `json:"price" db:"price"`     // from Last updated Price of this Machine.Column
	Max    int              `json:"max" db:"max"`         // ทดยอดสูงสุดที่จะสต๊อคสินค้าขายได้ในแต่ละ Column
}

//-------------------------------------------------
// Check mc.LastCounter ต้องน้อยกว่าหรือเท่ากับ CurrCounter
// Load MachineColumn table and validate new counter must greater than last counter.
//-------------------------------------------------
func (c *Counter) LessThanLastCount(columns []*MachineColumn) bool {

	for _, sub := range c.Sub {
		for _, column := range columns {
			if sub.ColumnNo == column.ColumnNo && sub.Counter < column.LastCounter {
				return true
			}
		}
	}
	return false
}

//---------------------------------------------------------------------------
// model.Counter.Insert
// ทำการเก็บผลการบันทึก Counter โดยมีการบันทึก LastCounter และ CurrCounter ลงใน
// MachineColumn ด้วย โดยต้องระวังการ Update จะไม่บันทึก LastCounter
// และถ้ามีการยกเลิก Counter ที่บันทึกไปแล้วต้องคืนค่า LastCounter และ CurrCounter ด้วย
//---------------------------------------------------------------------------
func (c *Counter) Insert(db *sqlx.DB) (*Counter, error) {
	var machine Machine
	machine.Id = c.MachineId
	columns, err := machine.GetColumns(db) // เอาข้อมูล lastCounter จาก columns ล่าสุดออกมา
	if err != nil {
		return nil, err
	}
	// ตรวจว่า เคาท์เตอร์ที่จดมาหากน้อยกว่าเคาท์เตอร์ล่าสุด ให้ออกและแจ้ง Error
	if c.LessThanLastCount(columns) {
		return nil, errors.New("Found error input counter: New counter < Last counter in the same Machine-Column.")
	}
	//---------------------
	// Insert to table Counter.
	//---------------------
	sql := `
		INSERT INTO counter (
		rec_date,
		machine_id,
		counter_sum
		) VALUES(?,?,?)
	`
	res, err := db.Exec(sql,
		c.RecDate.Time,
		c.MachineId,
		c.CounterSum,
	)
	if err != nil {
		log.Println("Error>>1.tx.Exec() INSERT INTO counter Error: ", err)
		return nil, err
	}
	log.Println("Pass>>1.tx.Exec() INSERT INTO counter:", c)

	id, _ := res.LastInsertId()
	c.Id = uint64(id)
	for _, sub := range c.Sub {
		// Update MachineColumn.LastCounter and CurrCounter
		column, err := machine.GetMachineColumn(db, sub.ColumnNo)
		if err != nil {
			log.Println("Error machine.GetMachineColumn:", err)
			return nil, err
		}
		column.LastCounter = column.CurrCounter
		column.CurrCounter = sub.Counter
		// update CounterSub.ItemId with last fulfill ItemId
		sub.ItemId = column.ItemId
		// update CounterSub.Price with last update Price
		sub.Price = column.Price
		err = column.Update(db)
		if err != nil {
			log.Println("Error mc.Update(db)", err)
			return nil, errors.New("Error Update MachineColumn" + err.Error())
		}
		// Insert new CounterSub
		sub.CounterId = c.Id
		err = sub.Insert(db)
		if err != nil {
			log.Println("Error sub.Insert(db)", err)
			return nil, err
		}
	}
	newCounter, err := c.Get(db)
	return newCounter, nil
}

//--------------------------------------------------
// Return New CounterSub of this Counter
//--------------------------------------------------
func (c *Counter) GetSub(db *sqlx.DB) ([]*CounterSub, error) {
	var sub []*CounterSub
	sql := `SELECT * FROM counter_sub WHERE counter_id = ? AND deleted IS NULL`
	err := db.Select(&sub, sql, c.Id)
	if err != nil {
		log.Println("Error>>5. model.Counter.GetSub() = ", err)
		return nil, err
	}
	log.Println("Pass>>5. model.Counter.GetSub()")
	return sub, nil
}

//--------------------------------------------------
// All will return only []Counter
// ถ้าต้องการ CounterSub จะต้องคอล Counter.Get() ทีละตัว
//--------------------------------------------------
func (c *Counter) GetAll(db *sqlx.DB) (counters []*Counter, err error) {
	log.Println("call model.Counter.GetAll()")
	// กรอง WHERE deleted <> null
	sql := `SELECT * FROM counter WHERE deleted IS NULL`
	err = db.Select(&counters, sql)
	if err != nil {
		return nil, err
	}
	sqlSub := `SELECT * FROM counter_sub WHERE counter_id = ? AND deleted IS NULL`
	for _, counter := range counters {
		err = db.Select(&counter.Sub, sqlSub, counter.Id)
		if err != nil {
			return nil, err
		}
		fmt.Println("Get:", counter)
	}
	log.Println(counters)
	return counters, nil
}

//-----------------------------------------------------------------
// model.Counter.Get() will return single Counter with []CounterSub
//-----------------------------------------------------------------
func (c *Counter) Get(db *sqlx.DB) (*Counter, error) {
	log.Println("call model.Counter.Get() c.ID=", c.Id)
	sql := `SELECT * FROM counter WHERE deleted IS NULL AND id = ?`
	err := db.Get(c, sql, c.Id)
	if err != nil {
		log.Println("Fail>>1.db.Get()", err)
		//log.Println("Fail>>1.db.QueryRowx", err)
		return nil, err
	}
	log.Println("Success>>1.db.QueryRowx")
	counterSub, err := c.GetSub(db)
	log.Println("counterSub=", counterSub)
	c.Sub = counterSub
	return c, nil
}

//-----------------------------------------------------------------
// Counter.Update()
// ระวัง การ model.Counter.Update() จะต้องไม่ update last_counter
// เราจะ update last_counter เฉพาะตอน Insert() เท่านั้น
//-----------------------------------------------------------------
func (c *Counter) Update(db *sqlx.DB) (*Counter, error) {
	var updatedCounter Counter
	return &updatedCounter, nil
}

//-----------------------------------------------------------------
// Counter.Delete()
// การยกเลิกบันทึก Counter โดยทำการ Update Counter.Deleted เพื่อลบรายการ และ
// ต้องเอา Counter ก่อนหน้า กลับมาใหม่ จาก CounterSub.Counter ก่อนหน้าด้วย
// โดยเขียนกลับลงไปใน MachineColumn.CurrCounter และ .LastCounter ตามลำดับ
//-----------------------------------------------------------------
func (c *Counter) Delete(db *sqlx.DB) error {
	sql := `DELETE counter WHERE id =?`
	_, err := db.Exec(sql, c.Id)
	if err != nil {
		return err
	}
	sql2 := `DELETE counter_sub WHERE id =?`
	for _, sub := range c.Sub {
		_, err := db.Exec(sql2, sub.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sub *CounterSub) Insert(db *sqlx.DB) error {
	//--------------------------------------------------
	// Insert CounterSub{}
	//--------------------------------------------------
	sql := `
		INSERT INTO counter_sub (
			counter_id,
			column_no,
			item_id,
			price,
			counter
		) VALUES(?,?,?,?,?)`
	res, err := db.Exec(sql,
		&sub.CounterId,
		&sub.ColumnNo,
		&sub.ItemId,
		&sub.Price,
		&sub.Counter,
	)
	if err != nil {
		log.Println("Error>> Exec() INSERT counter_sub = ", err)
		return err
	}
	id, _ := res.LastInsertId()
	log.Println("Pass>> Exec() INSERT counter_sub", id)
	return nil
}

func (c *Counter) GetAllByMachineCode(db *sqlx.DB, code string) (counters []*Counter, err error) {
	// หา machine_id จาก machine code ที่ได้รับ
	// todo: ควรเปลี่ยน args จาก code string -> m Machine
	id, err := GetMachineIdFromCode(db, code)
	if err != nil {
		return nil, err
	}
	sql2 := `SELECT * FROM counter WHERE machine_id = ? AND deleted IS NULL ORDER BY created DESC`
	err = db.Select(&counters, sql2, id)
	if err != nil {
		return nil, err
	}
	for _, counter := range counters {
		sql3 := `SELECT * FROM counter_sub WHERE counter_id = ?`
		err = db.Select(&counter.Sub, sql3, counter.Id)
		if err != nil {
			return nil, err
		}
		fmt.Println("add counter sub")
	}
	return counters, nil
}

func (c *Counter) GetLastByMachineCode(db *sqlx.DB, code string) (lastCounter *Counter, err error) {
	// จากตู้ไหน? หา machine_id จาก machine code ที่ได้รับ
	// todo: ควรเปลี่ยน args จาก code string -> m Machine ไหม?
	machineId, err := GetMachineIdFromCode(db, code)
	if err != nil {
		return nil, err
	}
	// คัดเคาท์เตอร์ที่บันทึกตามเวลาล่าสุดของตู้นี้
	sql2 := `SELECT * FROM counter WHERE machine_id = ? AND deleted IS NULL ORDER BY created DESC LIMIT 1 `
	err = db.Get(lastCounter, sql2, machineId)
	if err != nil {
		log.Println("Error when select last counter: ", err)
		return nil, err
	}
	// ดึงเอา CounterSub จากแต่ละคอลัมน์ขึ้นมา
	sql3 := `SELECT * FROM counter_sub WHERE counter_id = ? AND deleted IS NULL`
	err = db.Select(&lastCounter.Sub, sql3, lastCounter.Id)
	if err != nil {
		log.Println("Error when select counter_sub: ", err)
		return nil, err
	}
	// Todo: แปะ Max ให้ CounterSub โดยดึงจาก MachineColumn นั้นๆ
	//sql4 := `SELECT max_stock FROM machine_column WHERE machine_id =? AND column_no =?`
	//col := MachineColumn{}
	//for _, sub := range c.Sub {
	//	//todo: Refactor เพิ่มฟิลด์ max_stock
	//	err := db.Get(&col, sql4, c.MachineId, sub.ColumnNo)
	//	if err != nil {
	//		fmt.Println("Error DB.Get column from sub")
	//	}
	//	sub.Max = col.MaxStock
	//}
	return lastCounter, err
}

//func GetCounter(db *sqlx.DB, ids []uint64) ([]*Counter, error) {
//	sql := `SELECT * FROM batch_counter WHERE id = ?`
//	sales := []*Counter{}
//	for _, id := range ids {
//		var s Counter
//		log.Print("id: ", id)
//		err := db.Get(&s, sql, id)
//		if err != nil {
//			return nil, err
//		}
//		sales = append(sales, &s)
//		log.Print(sales)
//	}
//	return sales, nil
//}

//func NewArrayCounter(db *sqlx.DB, sales []*Counter) ([]*Counter, error) {
//	// Call from controller.PostMachineBatchSale()
//	tx, err := db.Beginx()
//	if err != nil {
//		return nil, err
//	}
//	sql := `INSERT INTO counter (
//		rec_date,
//		machine_id,
//		counter_sum
//		) VALUES(?,?,?)
//	`
//	var ids []uint64
//	for _, c := range sales {
//		res, err := tx.Exec(sql,
//			c.RecDate,
//			c.MachineId,
//			c.CounterSum,
//		)
//		if err != nil {
//			log.Println("error in tx.Exec(), res =", res, "Error: ", err)
//			errRollback := tx.Rollback()
//			if errRollback != nil {
//				log.Println("errRollback", errRollback)
//				return nil, errRollback
//			}
//			log.Println("tx.Rollback()", err)
//			return nil, err
//		}
//		id, err := res.LastInsertId()
//		log.Println("id = ", id, "err = ", err)
//		ids = append(ids, uint64(id))
//	}
//	err = tx.Commit()
//	if err != nil {
//		return nil, err
//	}
//	////Read from written DB.
//	log.Println(ids)
//	readSales, err := GetCounter(db, ids)
//	if err != nil {
//		log.Println("Error in GetBatchCounter() = ", err)
//		return nil, err
//	}
//	return readSales, nil
//}

// mysqlDriverErr
//if driverErr, ok := err.(*mysql.MySQLError); ok {
//	if driverErr.Number == mysqlerr.ER_DUP_INDEX {
//		// err handling here
//	}
//}
