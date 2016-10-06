package model

import (
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/sys/model"
	"log"
	"time"
)

type Stock struct {
	sys.Base
	LocID  uint64 `json:"loc_id" db:"loc_id"`
	ItemID uint64 `json:"item_id" db:"item_id"`
	Qty    int64  `json:"qty" db:"qty"`
}

type StockCard struct {
	sys.Base
	//DocID uint64
	ItemID    uint64
	LocID     uint64
	TransUnit Unit
	TransQty  int64
	BaseUnit  Unit
	BaseQty   int64
}

type Vehicle struct {
	sys.Base
	Name      string // V1, V2,...
	NamePlate string // ทะเบียนรถ
	Brand     string // ยี่ห้อ
}

// ข้อมูลฝ่ายบริการเติมสินค้า ลงทะเบียนเบิกกุญแจรถตอนเช้า จะมีผลกับการขาย VanSale  ไม่ต้องระบุรหัสรถ และผู้ขับ RouteMan อีก
type RouteMan struct {
	sys.Base
	Driver    sys.User
	Recorded  *time.Time
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
