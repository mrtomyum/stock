package model

import (
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/nava-sys/model"
	"log"
)

type Stock struct {
	sys.Base
	LocationID uint64 `json:"location_id"`
	ItemID     uint64 `json:"item_id"`
	Quantity   int64  `json:"quantity"`
}

type StockTrans struct {
	sys.Base
	//DocID uint64
	ItemID     uint64
	LocationID uint64
	TransUnit  Unit
	BaseUnit   Unit
	TransQty   int64
	BaseQty    int64
}

type ClientType int

const (
	FACTORY ClientType = 1 + iota
	EDUCATION
	OFFICE
)

type Client struct {
	sys.Base
	Name string
	Type ClientType
}

type Place struct {
	sys.Base
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
	sys.Base
	Name      string   // V1, V2,...
	NamePlate string   // ทะเบียนรถ
	Brand     carBrand // ยี่ห้อ
}

type RouteMan struct {
	sys.Base
	Name      string
	VehicleID uint64
}

func (s *Stock) All(db *sqlx.DB) ([]*Stock, error) {
	log.Println("call model.Stock.All()")
	sql := `SELECT * FROM stock`
	var stocks []*Stock
	err := db.Select(&stocks, sql)
	if err != nil {
		log.Println("Error: model.Stock.All() db.Select...", err)
		return nil, err
	}
	return stocks, nil
}
