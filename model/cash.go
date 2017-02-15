package model

import (
	"github.com/shopspring/decimal"
)

// CashCollection คือ ใบนำส่งเงิน ถูกสร้างจาก Counter
type CashCollection struct {
	Base
	UserId          uint64 `json:"user_id"`                // user ของ RouteMan ที่นำเงินมาส่ง
	CounterId       uint64 `json:"counter_id"`             // ใบนี้ถูกสร้างขึ้นจากการจดเคาท์เตอร์รายการใด
	CounterValue    decimal.Decimal `json:"counter_value"` // มูลค่าเงินนำส่งที่ประเมินจากยอดเคาท์เตอร์
	C1              int
	C2              int
	C5              int
	C10             int
	B20             int
	B50             int
	B100            int
	B500            int
	B1000           int
	CollectionValue decimal.Decimal `json:"collection_value"` // ยอดเงินที่นับได้
	Balance         decimal.Decimal                           // ผลต่างของ เงินที่นับได้ - เงินที่ต้องนำส่ง ติดลบคือเงินขาด ติดบวกคือเงินเกิน CollectionValue - CounterValue
}

