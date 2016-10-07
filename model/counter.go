package model

import (
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/sys/model"
	"github.com/shopspring/decimal"
	"log"
	"strings"
	"time"
	"database/sql/driver"
	"errors"
)

const DateFormat = "2006-01-02" // yyyy-mm-dd

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) error {
	log.Println("json.Unmashaller == Overide UnmarshalJSON()", string(data))
	var err error
	d.Time, err = time.Parse(DateFormat, strings.Trim(string(data), `"`)) // << ตรงนี้ต้องทำการ Trim(") ออก
	if err != nil {
		return err
	}
	return nil
}

func (d *Date) Value() (driver.Value, error) {
	return d.Time, nil
}

func (d *Date) Scan(src interface{}) error {
	if date, ok := src.(time.Time); ok {
		d.Time = date
		return nil
	}
	//d.Time = Date(src.(time.Time))
	return errors.New("wrong type it's not time.Time")
}

type Counter struct {
	sys.Base
	//RecDate    *time.Time `json:"rec_date" db:"rec_date"`
	RecDate    Date   `json:"rec_date" db:"rec_date"`
	MachineId  uint64 `json:"machine_id" db:"machine_id"`
	CounterSum int    `json:"counter_sum" db:"counter_sum"`
	Sub        []*CounterSub `json:"sub"`
}

type CounterSub struct {
	sys.Base
	CounterId uint64          `json:"-" db:"counter_id"`    // FK
	ColumnNo  int             `json:"column_no" db:"column_no"`
	Counter   int             `json:"counter" db:"counter"`
	ItemId    uint64          `json:"item_id" db:"item_id"` // Record as history data.
	Price     decimal.Decimal `json:"price"`                // from Last updated Price of this Machine.Column
}

//--------------------------------------------------
// Check mc.LastCounter ต้องน้อยกว่าหรือเท่ากับ CurrCounter
//--------------------------------------------------
func (c *Counter) FoundCounterLessThan(mcs []MachineColumn) bool {
	for _, sub := range c.Sub {
		for _, mc := range mcs {
			if sub.Counter < mc.LastCounter {
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
	//-------------------------------------------------
	// Load MachineColumn table and validate new counter must greater than last counter.
	//-------------------------------------------------
	var m Machine
	m.ID = c.MachineId
	mcs, err := m.GetColumns(db)
	if err != nil {
		return nil, err
	}
	if c.FoundCounterLessThan(mcs) {
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
	log.Println("Pass>>1.tx.Exec() INSERT INTO counter")
	//--------------------------------------------------
	// Loop for range Counter.Sub
	// Update MachineColumn.LastCounter and CurrCounter
	//--------------------------------------------------
	id, _ := res.LastInsertId()
	c.ID = uint64(id)
	for _, sub := range c.Sub {
		sql = `
			UPDATE machine_column
			SET last_counter = ?, curr_counter = ?
			WHERE machine_id = ? AND column_no = ?
			`
		for _, mc := range mcs {
			_, err := db.Exec(sql,
				mc.CurrCounter, sub.Counter,
				c.MachineId, mc.ColumnNo,
			)
			if err != nil {
				log.Println("Error>>3.tx.Exec() machine_column = ", err)
				return nil, err
			}
			sub.ItemId = mc.ItemId
			sub.Price = mc.Price
		}
		log.Println("Update MachineColumn 'MachineID':", c.MachineId, "ColumnNo:", sub.ColumnNo)
		log.Println("Pass>>3.tx.Exec() UPDATE machine_column")
		sub.CounterId = c.ID
		err = sub.Insert(db)
	}
	newCounter, err := c.Get(db)
	return newCounter, nil
}

func (c *Counter) GetSub(db *sqlx.DB) ([]*CounterSub, error) {
	//--------------------------------------------------
	// Return New CounterSub of this Counter
	//--------------------------------------------------
	var sub []*CounterSub
	//id := int8(c.ID)
	err := db.Select(&sub, `SELECT * FROM counter_sub WHERE counter_id = ?`, c.ID)
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
func (c *Counter) All(db *sqlx.DB) ([]*Counter, error) {
	log.Println("call model.Counter.All()")
	// กรอง WHERE deleted <> null
	sql := `SELECT * FROM counter WHERE deleted IS NULL`
	//err := db.Select(&counters, sql)
	row, err := db.Queryx(sql)
	defer row.Close()
	counters := []*Counter{} //<<-- น่าจะต้องไม่ใช่ Pointer นะ
	if row.Next() {
		err = row.Scan(
			&c.ID, &c.Created, &c.Updated, &c.Deleted,
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
func (c *Counter) Get(db *sqlx.DB) (*Counter, error) {
	log.Println("call model.Counter.Get() c.ID=", c.ID)
	sql := `SELECT * FROM counter WHERE deleted IS NULL AND id = ?`
	err := db.Get(c, sql, c.ID)
	if err != nil {
		log.Println("Fail>>1.db.Get()", err)
		//log.Println("Fail>>1.db.QueryRowx", err)
		return nil, err
	}
	log.Println("Success>>1.db.QueryRowx")
	//var counterSub []*CounterSub
	//sql = `SELECT * FROM counter_sub WHERE deleted IS NULL AND counter_id = ?`
	//err = db.Select(&counterSub, sql, c.ID)
	counterSub, err := c.GetSub(db)
	log.Println("counterSub=", counterSub)
	c.Sub = counterSub
	log.Println(c.Sub)
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