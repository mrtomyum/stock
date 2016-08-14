package model

import "github.com/mrtomyum/nava-sys/model"

type Stock struct {
	model.Base
	LocationID uint64 `json:"location_id"`
	ItemID     uint64 `json:"item_id"`
	Quantity   int64  `json:"quantity"`
}

type ClientType int

const (
	FACTORY ClientType = 1 + iota
	EDUCATION
	OFFICE
)

type Client struct {
	model.Base
	Name string
	Type ClientType
}

type Place struct {
	model.Base
	ClientID uint64
	Name     string
	Lat      float64
	Long     float64
}

type carBrand int

const (
	SUZUKI carBrand = 1 + iota
	TATA
)

type Vehicle struct {
	model.Base
	Name      string   // V1, V2,...
	NamePlate string   // ทะเบียนรถ
	Brand     carBrand // ยี่ห้อ
}

type RouteMan struct {
	model.Base
	Name      string
	VehicleID uint64
}
