package model

import (
	sys "github.com/mrtomyum/nava-sys/model"
	"time"
)

type Order struct {
	sys.Base
	Number string       `json:"number"`
	Date   time.Time    `json:"order_status"`
	Status sys.Status `json:"status"`
}

type OrderItemStatus int

const (
	HOLD_ITEM OrderItemStatus = 1 + iota
	PICK_ITEM
	BACK_ORDER_ITEM  // สินค้าไม่เพียงพอใน AVL ส่วนเกินจะรอสร้าง PO DRAFT
	PO_DRAFTED_ITEM  // เอกสารกำลังรออนุมัติ อาจไม่ใช้ โดยอาจรวมกับ BACKORDER
	PO_APPROVED_ITEM // PO อนุมัติ และส่งอีเมล์แล้ว
	SHIPPING_ITEM    // Vendor ส่งสินค้าออกมาแล้ว สินค้าอยู่ระหว่างจัดส่ง
	RECEIVED_ITEM    // ได้รับสินค้าแล้ว
)

type OrderItem struct {
	sys.Base
	OrderID  uint64
	Date     time.Time       `json:"date"`
	Status   OrderItemStatus `json:"status"`
	ItemID   uint64          `json:"itemID"`
	ItemName string          `json:"itemName"` // may differ from master data from time to time so it be recorded
	UnitSale Unit            `json:"unitSale"`
	Price    uint64          `json:"price"`
	Discount uint64          `json:"discount"`
	Total    uint64          `json:"total"`
} // Todo: Next stage Document that generate from this must be refer back to this Order Item
