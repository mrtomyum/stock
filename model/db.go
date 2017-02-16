package model

import (
	"log"
	"github.com/jmoiron/sqlx"
	"encoding/json"
	"os"
	"time"
	"strings"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

type Config struct {
	DBHost string `json:"db_host"`
	DBName string `json:"db_name"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
	Port   string `json:"port"`
}

func GetConfig(fileName string) string {
	file, _ := os.Open(fileName)
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Println("error:", err)
	}
	var dsn = config.DBUser + ":" + config.DBPass + "@" + config.DBHost + "/" + config.DBName + "?parseTime=true"
	return dsn
}

func init() {
	// Read configuration file from "cofig.json"
	//dsn := GetConfig("./model/config.json") // เปิดใช้งานจริงเมื่อ Docker Container run --link ตรงเข้า mariadb เท่านั้น
	dsn := GetConfig("./model/config_debug.json")
	DB = sqlx.MustConnect("mysql", dsn)
	log.Println("Connected db: ", DB)
}

// ใช้สำหรับล้างตารางทดสอบ Mock การเขียนลง DB Table ใดๆ พร้อม Reset Auto Increment index ให้ด้วย
func ResetTable(tableName string) error {
	sql1 := `TRUNCATE  ` + tableName
	res, err := DB.Exec(sql1)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	// Reset auto increment
	sql2 := `ALTER TABLE ` + tableName + ` AUTO_INCREMENT = 1`
	_, err = DB.Exec(sql2)
	if err != nil {
		return err
	}
	log.Println("Truncate Table:", tableName, "rows = ", rows)
	return nil
}

//type Status int
//
//const (
//	DRAFT Status = iota // ฉบับร่าง ค่า int = 0
//	OPEN // เอกสารบันทึกเข้าระบบงานแล้ว ยังสามารถแก้ไขได้
//	HOLD // เอกสารถูกพักรอดำเนินการ
//	POST // เอกสารถูกบันทึกบัญชีแล้ว ห้ามแก้ไข
//	CANCEL // เอกสารถูกยกเลิกแล้ว ห้ามแก้ไข
//)

// Base structure contains fields that are common to objects
// returned by the nava's REST API.
type Base struct {
	Id      uint64       `json:"id"`
	Created JsonNullTime `json:"-"`
	Updated JsonNullTime `json:"-"`
	Deleted JsonNullTime `json:"-"`
}

type JsonNullTime struct {
	mysql.NullTime
}

func (v *JsonNullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullTime) UnmarshalJSON(data []byte) error {
	var err error
	if data != nil {
		v.Time, err = time.Parse(time.RFC3339, strings.Trim(string(data), `"`))
		if err != nil {
			log.Println("Error in time.Parse()", err.Error())
			return err
		}
		v.Valid = true
		log.Println("data = ", string(data), "v.Time = ", v.Time)
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int8
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = int64(*x)
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullDate struct {
	mysql.NullTime
}

func (v *JsonNullDate) MarshalJSON() ([]byte, error) {
	if v.Valid {
		log.Println("MarshalJSON() v.Valid")
		return json.Marshal(v.Time)
	}
	log.Println("MarshalJSON() Invalid")
	return json.Marshal(nil)
}

func (v *JsonNullDate) UnmarshalJSON(data []byte) error {
	const LAYOUT = "2006-01-02"
	var err error
	if data != nil {
		v.Time, err = time.Parse(LAYOUT, strings.Trim(string(data), `"`))
		if err != nil {
			log.Println("Error in time.Parse()", err.Error())
			return err
		}
		v.Valid = true
		log.Println("data = ", string(data), "v.Time = ", v.Time)
	} else {
		v.Valid = false
	}
	return nil
}
