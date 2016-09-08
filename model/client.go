package model

import (
	sys "github.com/mrtomyum/nava-sys/model"
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
type Place struct {
	sys.Base
	ClientID uint64
	Name     string
	Lat      float64
	Long     float64
}