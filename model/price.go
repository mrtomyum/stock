package model

import (
	"encoding/json"
	"math"
	"log"
)

type Price struct {
	Value int64
	Digit float64
}

type ItemPrice struct {
	ItemID    uint64
	Locations []*Location
	Value     Price
}

type JsonPrice struct {
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