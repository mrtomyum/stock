package model

import (
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/nava-sys/model"
	"golang.org/x/text/currency"
	"log"
	"time"
)

type Counter struct {
	sys.Base
	RecDate    *time.Time `json:"rec_date" db:"rec_date"`
	MachineId  uint64     `json:"machine_id" db:"machine_id"`
	CounterSum int        `json:"counter_sum" db:"counter_sum"`
	Sub        []*CounterSub
}

type CounterSub struct {
	CounterId uint64          `json:"-" db:"counter_id"`    // FK
	ColumnNo  int             `json:"column_no" db:"column_no"`
	Counter   int             `json:"curr_counter" db:"curr_counter"`
	ItemId    uint64          `json:"item_id" db:"item_id"` // Record as history data.
	Price     currency.Amount `json:"price"`                // from Last updated Price of this Machine.Column
}

//-----------------------
// model.Counter.Insert
// ทำการเก็บผลการบันทึก Counter โดยมีการบันทึก LastCounter และ CurrCounter ลงใน
// MachineColumn ด้วย โดยต้องระวังการ Update จะไม่บันทึก LastCounter
// และถ้ามีการยกเลิก Counter ที่บันทึกไปแล้วต้องคืนค่า LastCounter และ CurrCounter ด้วย
//-----------------------
func (c *Counter) Insert(db *sqlx.DB) (*Counter, error) {
	tx, err := db.Beginx()
	sql := `INSERT INTO counter (
		rec_date,
		machine_id,
		counter_sum
		) VALUES(?,?,?)
	`
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec(sql,
		c.RecDate,
		c.MachineId,
		c.CounterSum,
	)
	if err != nil {
		// if err tx.Rollback
		log.Println("error in tx.Exec(), res =", res, "Error: ", err, "If duplicate entry please use model.Counter.Update()")
		errRollback := tx.Rollback()
		if errRollback != nil {
			log.Println("errRollback", errRollback)
			return nil, errRollback
		}
	}

	// Loop for range Counter.Sub
	counterId, _ := res.LastInsertId()
	var newSubs []*CounterSub
	for _, sub := range c.Sub {
		// Get relate data from other table.
		// Select related data from MachineColumn.
		var mc MachineColumn
		sql = `
			SELECT *
			FROM machine_column
			WHERE machine_id = ? AND number = ?
			LIMIT 1
			`
		err := db.Get(&mc, sql, c.MachineId, sub.ColumnNo)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Update MachineColumn.LastCounter and CurrCounter
		sql = `
			UPDATE machine_column
			SET last_counter = ?, curr_counter = ?
			WHERE machine_id = ? AND column_no = ?
			`
		res, err := tx.Exec(sql,
			mc.CurrCounter, sub.Counter,
			c.MachineId, sub.ColumnNo,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Insert CounterSub{}
		sql = `
			INSERT INTO counter_sub (
				counter_id,
				number,
				item_id,
				price,
				curr_counter
			)`
		sub.ItemId = mc.ItemId
		sub.Price = mc.Price
		res, _ = tx.Exec(sql,
			counterId,
			sub.ColumnNo,
			sub.ItemId,
			sub.Price,
			sub.Counter,
		)
		// Return New CounterSub to confirm
		var inserted CounterSub
		id, _ := res.LastInsertId()
		err = db.Get(&inserted, `SELECT * FROM counter_sub WHERE id = ?`, id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		newSubs = append(newSubs, &inserted)
	}

	// if success tx.Commit
	tx.Commit()
	var newCounter Counter
	sql = `SELECT * FROM counter WHERE id = ?`
	id, _ := res.LastInsertId()
	err = db.Get(newCounter, sql, uint64(id))
	if err != nil {
		log.Println("Error in db.Get() = ", err)
		return nil, err
	}
	newCounter.Sub = newSubs
	return &newCounter, nil
}

func (c *Counter) Update(db *sqlx.DB) (*Counter, error) {
	// ระวัง การ model.Counter.Update() จะไม่ update last_counter
	// เราจะ update last_counter เฉพาะตอน Insert() เท่านั้น
	var updatedCounter Counter
	return &updatedCounter, nil
}

func (c *Counter) Delete(db *sqlx.DB) error {
	// การยกเลิกบันทึก Counter โดยทำการ Update Counter.Deleted เพื่อลบรายการ และ
	// ต้องเอา Counter ก่อนหน้า กลับมาใหม่ จาก CounterSub.Counter ก่อนหน้าด้วย
	// โดยเขียนกลับลงไปใน MachineColumn.CurrCounter และ .LastCounter ตามลำดับ
	return nil
}

func (s *Counter) All(db *sqlx.DB) ([]*Counter, error) {
	log.Println("call model.BatchSale.All()")
	sales := []*Counter{}
	// Todo: เพิ่มกรอง WHERE deleted <> null
	sql := `SELECT * FROM batch_counter`
	err := db.Select(&sales, sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(sales)
	return sales, nil
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


