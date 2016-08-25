package model

import (
	"encoding/json"
	"math"
	"log"
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/nava-sys/model"
	"golang.org/x/text/currency"

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
	Value     currency.Amount
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

func (ip *ItemPrice) AllPrice(db *sqlx.DB) ([]*ItemPrice, error) {
	ips := []*ItemPrice{}
	return ips, nil
}

// บันทึกราคาจากหน้าตู้ ด้วยมือถือ
type BatchPrice struct {
	sys.Base
	Recorded  *time.Time      `json:"recorded"`
	MachineID uint64          `json:"machine_id"`
	ColumnNo  int             `json:"column_no"`
	Price     currency.Amount `json:"price"`
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