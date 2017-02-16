package model

import (
	"github.com/shopspring/decimal"
	"log"
	"errors"
	"fmt"
)

type Counter struct {
	Base
	//RecDate    *time.Time `json:"rec_date" db:"rec_date"`
	RecDate    Date   `json:"rec_date" db:"rec_date"`
	MachineId  uint64 `json:"machine_id" db:"machine_id"`
	CounterSum int    `json:"counter_sum" db:"counter_sum"`
	Sub        []*SubCounter `json:"sub"`
}

type SubCounter struct {
	Base
	CounterId  uint64          `json:"-" db:"counter_id"`    // FK
	ColumnNo   int             `json:"column_no" db:"column_no"`
	Counter    int             `json:"counter" db:"counter"`
	ItemId     uint64          `json:"item_id" db:"item_id"` // Record as history data.
	Price      decimal.Decimal `json:"price" db:"price"`     // from Last updated Price of this Machine.Column
	MaxCounter int `json:"max_counter"`                      // ทดยอดสูงสุดที่จะสต๊อคสินค้าขายได้ในแต่ละ Column
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
func (c *Counter) Insert() (*Counter, error) {
	var machine Machine
	machine.Id = c.MachineId
	columns, err := machine.GetColumns() // เอาข้อมูล lastCounter จาก columns ล่าสุดออกมา
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
	res, err := DB.Exec(sql,
		c.RecDate.Time,
		c.MachineId,
		c.CounterSum,
	)
	if err != nil {
		log.Println("Error>>1.tx.Exec() INSERT INTO counter Error: ", err)
		return nil, err
	}
	log.Println("Pass>>1.tx.Exec() INSERT INTO counter")
	log.Println(c)

	id, _ := res.LastInsertId()
	c.Id = uint64(id)
	for _, sub := range c.Sub {
		// Update MachineColumn.LastCounter and CurrCounter
		column, err := machine.GetMachineColumn(sub.ColumnNo)
		if err != nil {
			log.Println("Error machine.GetMachineColumn:", err)
			return nil, err
		}
		column.LastCounter = column.CurrCounter
		column.CurrCounter = sub.Counter
		// update SubCounter.ItemId with last fulfill ItemId
		sub.ItemId = column.ItemId
		// update SubCounter.Price with last update Price
		sub.Price = column.Price
		err = column.Update()
		if err != nil {
			log.Println("Error mc.Update(db)", err)
			return nil, errors.New("Error Update MachineColumn" + err.Error())
		}
		// Insert new SubCounter
		sub.CounterId = c.Id
		err = sub.Insert()
		if err != nil {
			log.Println("Error sub.Insert(db)", err)
			return nil, err
		}
	}
	newCounter, err := c.Get()
	return newCounter, nil
}

//--------------------------------------------------
// Return New CounterSub of this Counter
//--------------------------------------------------
func (c *Counter) GetSub() ([]*SubCounter, error) {
	var sub []*SubCounter
	sql := `SELECT * FROM counter_sub WHERE counter_id = ? AND deleted IS NULL`
	err := DB.Select(&sub, sql, c.Id)
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
func (c *Counter) GetAll() ([]*Counter, error) {
	log.Println("call model.Counter.All()")
	// กรอง WHERE deleted <> null
	sql := `SELECT * FROM counter WHERE deleted IS NULL`
	//err := db.Select(&counters, sql)
	row, err := DB.Queryx(sql)
	defer row.Close()
	counters := []*Counter{} //<<-- น่าจะต้องไม่ใช่ Pointer นะ
	if row.Next() {
		err = row.Scan(
			&c.Id, &c.Created, &c.Updated, &c.Deleted,
			&c.RecDate, &c.MachineId, &c.CounterSum,
		)
		counters = append(counters, c)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(counters)
	return counters, nil
}

//-----------------------------------------------------------------
// model.Counter.Get() will return single Counter with []CounterSub
//-----------------------------------------------------------------
func (c *Counter) Get() (*Counter, error) {
	log.Println("call model.Counter.Get() c.ID=", c.Id)
	sql := `SELECT * FROM counter WHERE deleted IS NULL AND id = ?`
	err := DB.Get(c, sql, c.Id)
	if err != nil {
		log.Println("Fail>>1.db.Get()", err)
		//log.Println("Fail>>1.db.QueryRowx", err)
		return nil, err
	}
	log.Println("Success>>1.db.QueryRowx")
	subCounter, err := c.GetSub()
	log.Println("subCounter=", subCounter)
	c.Sub = subCounter
	log.Println(c.Sub)
	return c, nil
}

//-----------------------------------------------------------------
// Counter.Update()
// ระวัง การ model.Counter.Update() จะต้องไม่ update last_counter
// เราจะ update last_counter เฉพาะตอน Insert() เท่านั้น
//-----------------------------------------------------------------
func (c *Counter) Update() (*Counter, error) {
	var updatedCounter Counter
	return &updatedCounter, nil
}
//-----------------------------------------------------------------
// Counter.Delete()
// การยกเลิกบันทึก Counter โดยทำการ Update Counter.Deleted เพื่อลบรายการ และ
// ต้องเอา Counter ก่อนหน้า กลับมาใหม่ จาก CounterSub.Counter ก่อนหน้าด้วย
// โดยเขียนกลับลงไปใน MachineColumn.CurrCounter และ .LastCounter ตามลำดับ
//-----------------------------------------------------------------
func (c *Counter) Delete() error {
	return nil
}

func (sub *SubCounter) Insert() error {
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
	res, err := DB.Exec(sql,
		sub.CounterId,
		sub.ColumnNo,
		sub.ItemId,
		sub.Price,
		sub.Counter,
	)
	if err != nil {
		log.Println("Error>> Exec() INSERT counter_sub = ", err)
		return err
	}
	id, _ := res.LastInsertId()
	log.Println("Pass>> Exec() INSERT counter_sub", id)
	return nil
}

func (c *Counter) GetAllByMachineCode(code string) (counters []*Counter, err error) {
	// หา machine_id จาก machine code ที่ได้รับ
	sql1 := `SELECT id FROM machine WHERE code = ?`
	var id uint64
	err = DB.Get(&id, sql1, code)
	if err != nil {
		return nil, err
	}

	sql2 := `SELECT * FROM counter WHERE machine_id = ? ORDER BY created DESC`
	err = DB.Select(&counters, sql2, id)
	if err != nil {
		return nil, err
	}
	for _, counter := range counters {
		sql3 := `SELECT * FROM counter_sub WHERE counter_id = ?`
		err = DB.Select(&counter.Sub, sql3, counter.Id)
		if err != nil {
			return nil, err
		}
		fmt.Println("add counter sub")
	}
	return counters, nil
}

func (c *Counter) GetLastByMachineCode(code string) error {
	// จากตู้ไหน? หา machine_id จาก machine code ที่ได้รับ
	sql1 := `SELECT id FROM machine WHERE code = ?`
	var id uint64
	err := DB.Get(&id, sql1, code)
	if err != nil {
		return err
	}
	// คัดเคาท์เตอร์ที่บันทึกตามเวลาล่าสุดของตู้นี้
	sql2 := `SELECT * FROM counter WHERE machine_id = ? ORDER BY created DESC LIMIT 1 `
	err = DB.Get(c, sql2, id)
	if err != nil {
		log.Println("Error when select last counter: ", err)
		return err
	}
	// ดึงเอา Sub Counter จากแต่ละคอลัมน์ขึ้นมา
	sql3 := `SELECT * FROM counter_sub WHERE counter_id = ?`
	err = DB.Select(&c.Sub, sql3, id)
	if err != nil {
		log.Println("Error when select counter_sub: ", err)
		return err
	}
	// Todo: แปะ MaxCounter ให้ SubCounter โดยดึงจาก MachineColumn นั้นๆ
	sql4 := `SELECT max_stock FROM machine_column WHERE machine_id = ?`
	for _, sub := range c.Sub {
		//todo: Refactor เพิ่มฟิลด์ max_stock
	}
	return nil
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