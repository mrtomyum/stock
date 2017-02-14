package model

import (
	sys "github.com/mrtomyum/sys/model"
	"github.com/shopspring/decimal"
	"log"
)

// Stock คือ Stock Card นั่นเอง เก็บรายการสินค้าและ Loc ที่เดินรายการ ผ่าน Doc
// ทั้งนี้ระบบสต๊อคบัญชีคู่ หมายความว่าจะต้องมี Stock อย่างน้อย 2 รายการ
// โดยมีจำนวน Qty และ Value เป็น + และ - คู่กัน
// ทั้งนี้ BaseCost จะถูกคำนวณกลับจาก Lot.Cost ที่ผูกกันผ่าน StockLot
type Stock struct {
	Base
	DocSubId uint64          `json:"doc_sub_id" db:"doc_sub_id"` //
	DocDate  sys.Date        `json:"doc_date" db:"doc_date"`
	ItemId   uint64          `json:"item_id" db:"item_id"`
	LocId    uint64          `json:"loc_id" db:"loc_id"`
	BaseUnit Unit            `json:"base_unit" db:"base_unit"`   // หน่วยนับฐาน
	BaseQty  int64           `json:"base_qty" db:"base_qty"`     // จำนวนตามหน่วยนับฐาน
	BaseCost decimal.Decimal `json:"base_cost" db:"base_cost"`   // มูลค่าทุน ต่อหน่วยตาม BaseUnit โดยไม่รวม VAT
	LotIds   []uint64
}

// StockLot บันทึกการอ้างอิง Lot ของ item ที่เดินรายการใน Stock ตามจำนววนจริง
// ทั้งนี้ StockLot มีไว้เพื่อการบันทึก Persistent Data
type StockLot struct {
	StockId uint64
	LotId   uint64
}

// Lot เก็บบันทึกต้นทุนสินค้าแต่ละชิ้น ที่อยู่ในแต่ละ Loc
// มีลักษณะเหมือน Lot/Serial Control แต่พิเศษคือ 1 Record = 1 ชิ้น  ที่ BaseUnit
// ประโยชน์คือเพื่อให้ทราบต้นทุนต่อหน่วยฐานที่เล็กที่สุดได้แบบ ชิ้นต่อชิ้น
type Lot struct {
	Base
	ItemId    uint64 `json:"item_id" db:"item_id"`
	Cost      decimal.Decimal
	CurrLocId uint64 `json:"curr_loc_id" db:"curr_loc_id"`
}

// StockItem เก็บยอดสินค้าคงเหลือ(Quantity) และมูลค่า (Value) ณ ปัจจุบัน
// ณ Location หนึ่งๆ
// ของสินค้าแต่ละ Item ซึ่งจะถูก Update ทุกครั้งที่มีการเคลื่อนไหวของ Stock
type StockItem struct {
	Base
	ItemId uint64 `json:"item_id" db:"item_id"`
	LocId  uint64 `json:"loc_id" db:"loc_id"`
	Qty    int64  `json:"qty" db:"qty"`
	Value  decimal.Decimal
}

func (s *StockItem) GetAll() ([]*StockItem, error) {
	log.Println("call model.Stock.All()")
	sql := `SELECT * FROM stock`
	var stocks []*StockItem
	err := DB.Select(&stocks, sql)
	if err != nil {
		log.Println("Error: model.Stock.All() db.Select...", err)
		return nil, err
	}
	return stocks, nil
}

type Unit struct {
	Base
	Name   string `json:"name" db:"name"`
	NameEn string `json:"name_en" db:"name_en"`
}