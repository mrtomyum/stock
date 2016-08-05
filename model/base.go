package model

import (
	"github.com/go-sql-driver/mysql"
	"github.com/mrtomyum/nava-api3/model"
	"database/sql"
	"encoding/json"
	"time"
)

type Status int

const (
	DRAFT Status = iota // ฉบับร่าง ค่า int = 0
	OPEN // เอกสารบันทึกเข้าระบบงานแล้ว ยังสามารถแก้ไขได้
	HOLD // เอกสารถูกพักรอดำเนินการ
	POST // เอกสารถูกบันทึกบัญชีแล้ว ห้ามแก้ไข
	CANCEL // เอกสารถูกยกเลิกแล้ว ห้ามแก้ไข
)

// Base structure contains fields that are common to objects
// returned by the nava's REST API.
type Base struct {
	ID      uint64         `json:"id"`
	//Created mysql.NullTime `json:"created"`
	Created JsonNullTime `json:"created"`
	Updated JsonNullTime `json:"updated"`
	Deleted JsonNullTime `json:"deleted"`
}

type Doc struct {
	CreatedBy  model.User
	ApprovedBy model.User
	UpdatedBy  model.User
	DeletedBy  model.User
}

type JsonNullInt struct {
	sql.NullInt64
}

func (v JsonNullInt) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

type JsonNullString struct {
	sql.NullString
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v JsonNullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.String = *x
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullTime struct {
	mysql.NullTime
}

func (v JsonNullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (v JsonNullTime) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Time = *x
	} else {
		v.Valid = false
	}
	return nil
}