package model

import (
	"github.com/go-sql-driver/mysql"
	"github.com/mrtomyum/nava-api3/model"
)

// Base structure contains fields that are common to objects
// returned by the nava's REST API.
type Base struct {
	ID      uint64         `json:"id"`
	Created mysql.NullTime `json:"created"`
	Updated mysql.NullTime `json:"updated"`
	Deleted mysql.NullTime `json:"deleted"`
}

type Doc struct {
	CreatedBy  model.User
	ApprovedBy model.User
	UpdatedBy  model.User
	DeletedBy  model.User
}

type Status int

const (
	DRAFT Status = iota // ฉบับร่าง ค่า int = 0
	OPEN // เอกสารบันทึกเข้าระบบงานแล้ว ยังสามารถแก้ไขได้
	HOLD // เอกสารถูกพักรอดำเนินการ
	POST // เอกสารถูกบันทึกบัญชีแล้ว ห้ามแก้ไข
	CANCEL // เอกสารถูกยกเลิกแล้ว ห้ามแก้ไข
)
