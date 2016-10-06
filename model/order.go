package model

import (
	sys "github.com/mrtomyum/sys/model"
	"time"
)

type Order struct {
	sys.Base
	sys.Doc
	Type   OrderType  `json:"type"`
	Number string     `json:"number"`
	Date   time.Time  `json:"order_status"`
	Status sys.Status `json:"status"`
}

type OrderType int

const (
	DRAFT OrderType = iota
	SO              // Sale Order
	RO              // Reserve Order
	PR              // Purchase Requisition
	PO              // Purchase Order
	IV              // Invoice
)

type ItemStatus int

const (
	HOLD_ITEM ItemStatus = iota
	PICK_ITEM
	BACK_ORDER_ITEM  // สินค้าไม่เพียงพอใน AVL ส่วนเกินจะรอสร้าง PO DRAFT
	PO_DRAFTED_ITEM  // เอกสารกำลังรออนุมัติ อาจไม่ใช้ โดยอาจรวมกับ BACKORDER
	PO_APPROVED_ITEM // PO อนุมัติ และส่งอีเมล์แล้ว
	SHIPPING_ITEM    // Vendor ส่งสินค้าออกมาแล้ว สินค้าอยู่ระหว่างจัดส่ง
	RECEIVED_ITEM    // ได้รับสินค้าแล้ว
)

type OrderItem struct {
	sys.Base
	OrderID   uint64
	Status    ItemStatus `json:"status"`
	Date      time.Time  `json:"date"`
	ItemID    uint64     `json:"itemID"`
	ItemName  string     `json:"itemName"` // may differ from master data from time to time so it be recorded
	Price     uint64     `json:"price"`
	Discount  uint64     `json:"discount"`
	Total     uint64     `json:"total"`
	TransUnit Unit       `json:"trans_unit"`
	BaseUnit  Unit       `json:"base_unit"`
	TransQty  int64      `json:"trans_qty"`
	BaseQty   int64      `json:"base_qty"`
} // Todo: Next stage Document that generate from this must be refer back to this Order Item
