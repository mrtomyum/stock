package model

import (
	sys "github.com/mrtomyum/sys/model"
	"github.com/shopspring/decimal"
	"time"
)

// RouteMan สามารถขายสินค้านอกตู้ได้ในหลายๆกรณี เช่นยังเติมของไม่เสร็จ และเก็บเงินสดจากการขายนำส่งต่างหากได้
// พฤติกรรมงานจะเหมือนกับการขาย POS แต่สรุปยอดบิลวันละ 1 ใบ หลายๆรายการรวมกัน ตัดสต๊อคจากท้ายรถ ไม่ใช่จากตู้
// ลองออกแบบ Type นี้ไว้รับข้อมูลดิบก่อนหาข้อมูลสินค้า ตรวจความถูกต้องก่อนใช้ VanSale เก็บ เขียนลง DB
type VanSaleRawData struct {
	Recorded   *time.Time      `json:"recorded"`
	Barcode    string          `json:"barcode"`
	Qty        int             `json:"qty"`
	PriceUnit  decimal.Decimal `json:"price_unit"`
	PriceTotal decimal.Decimal `json:"price_total"`
}

// ใช้ type นี้  map DB
type VanSale struct {
	sys.Base
	Recorded   *time.Time      `json:"recorded"`
	Barcode    string          `json:"barcode"`
	ItemId     uint64          `json:"item_id"`
	Qty        int             `json:"qty"`
	UnitPrice  decimal.Decimal `json:"unit_price"`
	TotalPrice decimal.Decimal `json:"total_price"`
}

// Design this struct for data from VMC telemetry system.
//type SaleStatus int
//const (
//	COMPLETED SaleStatus = iota
//	INCOMPLETED
//)
//type RealTimeSale struct {
//	sys.Base
//	Recorded  *time.Time      `json:"recorded"`
//	MachineID uint64          `json:"machine_id"`
//	ColumnNo  int             `json:"column_no"`
//	ItemID    uint64          `json:"item_id"`
//	Price     decimal.Decimal `json:"price"`
//	//Status SaleStatus `json:"status"`
//}
