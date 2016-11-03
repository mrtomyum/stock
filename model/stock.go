package model

import (
	sys "github.com/mrtomyum/sys/model"
	"github.com/shopspring/decimal"
	"log"
)

// DocType แยกประเภทเอกสาร ที่ใช้เดินรายการ Stock แต่ละรายการ
// โดยใช้โครงสร้างข้อมูลร่วมกันได้ หากมี Property ใดของเอกสารแต่ละประเภท
// พิศดารออกไป ให้แยกตารางอ้างอิงเพิ่มเอา
type DocType int

const (
	DOC_NO_TYPE DocType = iota // 0 = ไม่กำหนด Type
	DOC_BUY                // ซื้อ
	DOC_SELL               // ขาย
	DOC_PAYMENT            // จ่ายเงิน
	DOC_RECEIVE            // รับเงิน
	DOC_STOCK_MOVE         // เคลื่อนย้ายสินค้า
)

//===================================
//           ทดสอบ
//===================================
type DocCode int

const (
	NO_CODE DocCode = iota
	PO // Purchase Order ใบสั่งซื้อ
	RV // Receiving ใบรับสินค้า
	IV // Invoice บิลขาย
	IC
)

// TaxType ประเภทภาษี ที่จะระบุในหัวเอกสาร มีผลกับการบันทึกเพื่อคำนวณภาษี ให้ตรงฟิลด์
// โดย API จะตรวจว่าเอกสารระบุราคาแบบ รวม หรือ แยก ภาษี แล้วบันทึกมูลค่าให้ถูกต้อง
// เช่น เอกสารซื้อมักให้ราคาแยก VAT (EXCLUDE_VAT) method DocGen() จะบันทึก
// มูลค่าสินค้า ในฟิลด์กลุ่มที่ไม่มี _Tax ต่อท้าย และคำนวณยอดรวมภาษีไว้ให้ในฟิลด์  _Tax
type TaxType int

const (
	NO_TAX TaxType = iota // 0 = เอกสารที่ไม่บันทึกภาษีขาย
	INCLUDE_TAX                // 1 = เอกสารที่สินค้าทั้งหมดราคารวม VAT
	EXCLUDE_TAX                // 2 = เอกสารที่สินค้าทั้งหมดราคาไม่รวม VAT
)

// Doc คือ หัวเอกสารที่เกี่ยวกับการเงิน และเคลื่อนไหวของสินค้าทั้งหมด
// โดยส่วนหัวเอกสารนี้จะเก็บข้อมูลหลักของ Transaction และยอดสรุปทางการเงิน และภาษี
// แยกประเภทด้วย Type แต่ละ Type จะมี Template/Method ควบคุมพฤติกรรม ในการบันทึกบัญชีคู่ เช่น
// ซื้อ จะต้อง ลดคลัง Supplier และ เพิ่มคลัง Inventory เสมอ เป็นต้น
// ปกติจะใช้ Method Doc.New() สร้างเอกสารขึ้น ซึ่งจะเรียกใช้ Doc.Get()
// ช่วยรันเลขที่เอกสารล่าสุดโดยต้องระบ DocCode กำหนดกลุ่มด้วยรหัสหัวเอกสารเพื่อรันเลขที่ด
// และ BranchId กำหนดรหัสสาขาด้วย
type Doc struct {
	sys.Base
	Num       string                                                  // เลขที่เอกสาร ที่สร้างจาก func DocGen(DocType, Company, Branch)
	Type      DocType                                                 // ประเภทเอกสาร
	DocDate   Date            `json:"doc_date" db:"doc_date"`         // วันที่ของเอกสาร (ไม่มีเวลา สร้างจาก Base.Created)
	BranchId  uint64          `json:"branch_id" db:"branch_id"`       // สาขาที่ทำรายการ
	AccountId uint64          `json:"account_id" db:"account_id"`     // บัญชีลูกหนี้ หรือเจ้าหนี้ที่เดินรายการ
	LocOut    uint64          `json:"loc_out" db:"loc_out"`           // LocId ที่ต้อง - Stock
	LocIn     uint64          `json:"loc_in" db:"loc_in"`             // LocId ที่ต้อง + Stock
	TaxType   TaxType         `json:"tax_type" db:"tax_type"`         // ประเภทภาษี
	TaxRate   decimal.Decimal `json:"tax_rate" db:"tax_rate"`         // อัตราภาษีมูลค่าเพิ่ม เป็นเปอร์เซนต์
	Amount    decimal.Decimal `json:"amount" db:"amount"`             // มูลค่ารวมสินค้า ที่มี VAT รวมอยู่แล้ว
	Discount  decimal.Decimal `json:"discount"`                       // ส่วนลดรวม VAT
	Total     decimal.Decimal `json:"total" db:"total"`               // มูลค่าเอกสารรวม VAT
	NetAmount decimal.Decimal `json:"net_amount" db:"net_amount"`     // มูลค่ารวมสินค้าที่ไม่มี VAT
	NetDisc   decimal.Decimal `json:"net_discount" db:"net_discount"` // ส่วนลดไม่รวม VAT
	NetTotal  decimal.Decimal `json:"net_total" db:"net_total"`       // มูลค่าเอกสารไม่รวม VAT
	Tax       decimal.Decimal `json:"tax"`                            // มูลค่าภาษี VAT สุทธิของเอกสารนี้
	DocSubs   []*DocSub       `json:"doc_subs" db:"doc_subs"`
}

// DocSub เป็นตารางลูกของ Doc เก็บรายละเอียดสินค้าที่เดินรายการปกติเอาไว้
// โดยจะแสดงรายละเอียดสำคัญด้านภาษี และการเงินไว้ด้วย
type DocSub struct {
	Id        uint64
	DocId     uint64                                              // อ้างถึง StockId
	ItemId    uint64                                              // Item ที่เดินรายการ
	UnitId    uint64                                              // หน่วยนับที่เดินรายการ
	Qty       int64                                               // จำนวนที่เดินรายการ
	Value     decimal.Decimal                                     // มูลค่าต่อหน่วยที่เดินรายการ
	Disc      decimal.Decimal                                     // ส่วนลดย่อย
	Amount    decimal.Decimal `json:"amount" db:"amount"`         // ผลรวมหลังหักส่วนลด รวมภาษี
	NetAmount decimal.Decimal `json:"net_amount" db:"net_amount"` // ผลรวมหลังหักส่วนลด หักภาษีแล้ว
}

// Stock คือ Stock Card นั่นเอง เก็บรายการสินค้าและ Loc ที่เดินรายการ ผ่าน Doc
// ทั้งนี้ระบบสต๊อคบัญชีคู่ หมายความว่าจะต้องมี Stock อย่างน้อย 2 รายการ
// โดยมีจำนวน Qty และ Value เป็น + และ - คู่กัน
// ทั้งนี้ BaseCost จะถูกคำนวณกลับจาก Lot.Cost ที่ผูกกันผ่าน StockLot
type Stock struct {
	Id       uint64
	DocSubId uint64          `json:"doc_sub_id" db:"doc_sub_id"` //
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
	Id        uint64
	ItemId    uint64 `json:"item_id" db:"item_id"`
	CurrLocId uint64 `json:"curr_loc_id" db:"curr_loc_id"`
	Cost      decimal.Decimal
}

// StockItem เก็บยอดสินค้าคงเหลือ(Quantity) และมูลค่า (Value) ณ ปัจจุบัน
// ณ Location หนึ่งๆ
// ของสินค้าแต่ละ Item ซึ่งจะถูก Update ทุกครั้งที่มีการเคลื่อนไหวของ Stock
type StockItem struct {
	sys.Base
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

// Construct New Document
// Parameter t DocType, data []byte from UnmarshalJSON
func (d *Doc) New(t DocType, data []byte) error {
	// Call Doc.Gen()
	return nil
}

func (d *Doc) Gen(head DocCode)
