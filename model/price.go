package model

import (
	"encoding/json"
	sys "github.com/mrtomyum/sys/model"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"time"
)

type MyPrice struct {
	Base
	Value    int64
	Digit    float64
	Currency decimal.Decimal
}

type ItemPrice struct {
	Base
	ItemId    uint64
	Locations []*Location
	Value     decimal.Decimal
}

type JsonPrice struct {
	Base
	MyPrice
}

func (p JsonPrice) MashalJSON() ([]byte, error) {
	dec := float64(p.Value) * math.Pow(10, -p.Digit)
	output, err := json.Marshal(dec)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return output, nil
}

func (ip *ItemPrice) GetAll() ([]*ItemPrice, error) {
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
	Price     decimal.Decimal `json:"price"`
}

func (bp *BatchPrice) New() (*BatchPrice, error) {
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
	res, err := DB.Exec(sql,
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
	err = DB.Get(
		&newBatchPrice,
		`SELECT * FROM batch_price WHERE id = ?`,
		id)
	if err != nil {
		log.Println("Error db.Get()=", err)
		return nil, err
	}
	return &newBatchPrice, nil
}

func (s *BatchPrice) GetAll() ([]*BatchPrice, error) {
	log.Println("call model.BatchPrice.All()")
	prices := []*BatchPrice{}
	sql := `SELECT * FROM batch_price`
	err := DB.Select(&prices, sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(prices)
	return prices, nil
}
