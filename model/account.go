package model

import (
	sys "github.com/mrtomyum/sys/model"
	"github.com/shopspring/decimal"
)

type Account struct {
	sys.Base
	Code     string `json:"code" db:"code"`
	Name     string `json:"name" db:"name"`
	NameEn   string `json:"name_en" db:"name_en"`
	ApValue  decimal.Decimal `json:"ap_value" db:"ap_value"`   // มูลค่าเจ้าหนี้ บริษัทค้างเงินเค้าเท่าใด
	ArValue  decimal.Decimal `json:"ar_value" db:"ar_value"`   // มูลค่าลูกหนี้ เค้าค้างเงินเราเท่าใด
	ApCredit decimal.Decimal `json:"ap_credit" db:"ap_credit"` // เครดิตเจ้าหนี้คงค้างจ่าย
	ArCredit decimal.Decimal `json:"ar_credit" db:"ar_credit"` // เครดติลูกหนี้คงค้างรับ
}
