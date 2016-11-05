package model

import (
	"time"
	"database/sql/driver"
	"strings"
	"log"
	"errors"
)

const DateFormat = "2006-01-02" // yyyy-mm-dd

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) error {
	log.Println("json.Unmashaller == Overide UnmarshalJSON()", string(data))
	var err error
	d.Time, err = time.Parse(DateFormat, strings.Trim(string(data), `"`)) // << ตรงนี้ต้องทำการ Trim(") ออก
	if err != nil {
		return err
	}
	return nil
}

func (d *Date) Value() (driver.Value, error) {
	return d.Time, nil
}

func (d *Date) Scan(src interface{}) error {
	if date, ok := src.(time.Time); ok {
		d.Time = date
		return nil
	}
	//d.Time = Date(src.(time.Time))
	return errors.New("wrong type it's not time.Time")
}
