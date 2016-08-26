package model

import (
	"log"
	"github.com/jmoiron/sqlx"
	"time"
	sys "github.com/mrtomyum/nava-sys/model"
)

type Counter struct {
	sys.Base
	Recorded  *time.Time `json:"recorded"`
	MachineID uint64     `json:"machine_id" db:"machine_id"`
	ColumnNo  int        `json:"column_no" db:"column_no"`
	Counter   int        `json:"counter"`
	//SalePrice currency.Amount `json:"-" db:"sale_price"` // SalePrice search data from Last update Price of this Machine.Column
}

func (s *Counter) All(db *sqlx.DB) ([]*Counter, error) {
	log.Println("call model.BatchSale.All()")
	sales := []*Counter{}
	sql := `SELECT * FROM batch_counter`
	err := db.Select(&sales, sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(sales)
	return sales, nil
}

func (c *Counter) NewCounter(db *sqlx.DB) (*Counter, error) {
	sql := `INSERT INTO batch_counter (
		recorded,
		machine_id,
		column_no,
		counter
		) VALUES(?,?,?,?)
	`
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec(sql,
		c.Recorded,
		c.MachineID,
		c.ColumnNo,
		c.Counter,
	)
	if err != nil {
		log.Println("error in tx.Exec(), res =", res, "Error: ", err)
		errRollback := tx.Rollback()
		if errRollback != nil {
			log.Println("errRollback", errRollback)
			return nil, errRollback
		}
	}
	newCounter := new(Counter)
	sql = `SELECT * FROM batch_counter WHERE id = ?`
	id, _ := res.LastInsertId()
	err = db.Get(*newCounter, sql, uint64(id))
	if err != nil {
		log.Println("Error in db.Get() = ", err)
		return nil, err
	}
	return newCounter, nil
}

func NewArrayCounter(db *sqlx.DB, sales []*Counter) ([]*Counter, error) {
	// Call from controller.PostMachineBatchSale()
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	sql := `INSERT INTO batch_counter (
		recorded,
		machine_id,
		column_no,
		counter
		) VALUES(?,?,?,?)
	`
	var ids []uint64
	for _, c := range sales {
		res, err := tx.Exec(sql,
			c.Recorded,
			c.MachineID,
			c.ColumnNo,
			c.Counter,
		)
		if err != nil {
			log.Println("error in tx.Exec(), res =", res, "Error: ", err)
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Println("errRollback", errRollback)
				return nil, errRollback
			}
			log.Println("tx.Rollback()", err)
			return nil, err
		}
		id, err := res.LastInsertId()
		log.Println("id = ", id, "err = ", err)
		ids = append(ids, uint64(id))
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	////Read from written DB.
	log.Println(ids)
	readSales, err := GetCounter(db, ids)
	if err != nil {
		log.Println("Error in GetBatchCounter() = ", err)
		return nil, err
	}
	return readSales, nil
}

func GetCounter(db *sqlx.DB, ids []uint64) ([]*Counter, error) {
	sql := `SELECT * FROM batch_counter WHERE id = ?`
	sales := []*Counter{}
	for _, id := range ids {
		var s Counter
		log.Print("id: ", id)
		err := db.Get(&s, sql, id)
		if err != nil {
			return nil, err
		}
		sales = append(sales, &s)
		log.Print(sales)
	}
	return sales, nil
}
