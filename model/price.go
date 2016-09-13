package model

import (
	"encoding/json"
	"math"
	"log"
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/nava-sys/model"
	"golang.org/x/text/currency"

	"time"
	"github.com/shopspring/decimal"
)

type MyPrice struct {
	sys.Base
	Value    int64
	Digit    float64
	Currency currency.Amount
}

type ItemPrice struct {
	sys.Base
	ItemID    uint64
	Locations []*Location
	Value     decimal.Decimal
}

type JsonPrice struct {
	sys.Base
	MyPrice
}

func (p JsonPrice) MashalJSON() ([]byte, error) {
	decimal := float64(p.Value) * math.Pow(10, - p.Digit)
	output, err := json.Marshal(decimal)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return output, nil
}

func (ip *ItemPrice) All(db *sqlx.DB) ([]*ItemPrice, error) {
	ips := []*ItemPrice{}
	//todo: ยังไม่เสร็จ
	return ips, nil
}

// =======================
// บันทึกราคาจากหน้าตู้ ด้วยมือถือ
// =======================
type BatchPrice struct {
	sys.Base
	Recorded  *time.Time      `json:"recorded"`
	MachineID uint64          `json:"machine_id"`
	ColumnNo  int             `json:"column_no"`
	Price     currency.Amount `json:"price"`
}

func (bp *BatchPrice) New(db *sqlx.DB) (*BatchPrice, error) {
	log.Println("call model.BatchPrice.New()")
	sql := `
		INSERT INTO batch_price(
			recorded,
			machine_id,
			column_no,
			price
		)
		VALUES(?,?,?,?)
		`
	res, err := db.Exec(sql,
		bp.Recorded,
		bp.MachineID,
		bp.ColumnNo,
		bp.Price,
	)
	if err != nil {
		log.Println("Error db.Exec()=", err)
		return nil, err
	}
	id, _ := res.LastInsertId()
	newBatchPrice := BatchPrice{}
	err = db.Get(
		&newBatchPrice,
		`SELECT * FROM batch_price WHERE id = ?`,
		id)
	if err != nil {
		log.Println("Error db.Get()=", err)
		return nil, err
	}
	return &newBatchPrice, nil
}

func (s *BatchPrice) All(db *sqlx.DB) ([]*BatchPrice, error) {
	log.Println("call model.BatchPrice.All()")
	prices := []*BatchPrice{}
	sql := `SELECT * FROM batch_price`
	err := db.Select(&prices, sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(prices)
	return prices, nil
}