package model

import (
	sys "github.com/mrtomyum/nava-sys/model"
	"go/types"
)

type Client struct {
	sys.Base
	Code   string `json:"code"`
	NameTH string `json:"name_th" db:"name_th"`
	NameEN string `json:"name_en" db:"name_en"`
	//Type       ClientType
	//PriceLevel int
}

// สถานที่วางตู้ ///น่าจะรวมกับ Client เสียดีไหม?

type PlaceType int

const (
	NOT_SPECIFY PlaceType = iota
	FACTORY
	SCHOOL
	DORM
	PATHWAY
	MARKET_PLACE
	DEPARTMENT_STORE
	MODERN_TRADE
)

type Place struct {
	sys.Base
	ClientID uint64    `json:"client_id"`
	NameTh   string    `json:"name_th" db:"name_th"`
	NameEn   string    `json:"name_en" db:"name_en"`
	Address  string `json:"address"`
	Type     PlaceType `json:"type"`
	Lat      float64   `json:"lat"`
	Lng      float64   `json:"lng"`
}
