package model

import (
	"encoding/json"
	"math"
	"log"
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/nava-sys/model"
	"golang.org/x/text/currency"
)

type Currency int


type Price struct {
	sys.Base
	Value    int64
	Digit    float64
	Currency Currency
}

type ItemPrice struct {
	sys.Base
	ItemID    uint64
	Locations []*Location
	Value     Price
}

type JsonPrice struct {
	sys.Base
	Price
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