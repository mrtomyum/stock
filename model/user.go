package model

import (
	"log"
)

type User struct {
	Base
	Name   string `json:"name" db:"name"`
	Secret []byte `json:"-" db:"secret"`
	Title  UserTitle `json:"title" db:"title"`
}

type UserTitle int

const (
	ADMIN    UserTitle = iota
	ROUTEMAN
	STOREMAN
	CASHIER
)

func (u *User) New() error {
	sql := `INSERT INTO user(name, title) VALUES(?,?)`
	_, err := DB.Exec(sql, u.Name, u.Title)
	if err != nil {
		return err
	}
	sql2 := `SELECT * FROM user WHERE name = ? LIMIT 1`
	err = DB.Get(u, sql2, u.Name)
	if err != nil {
		log.Println("error DB.Get:", err)
	}
	log.Println("User:", u)
	return nil
}

func (u *User) Delete() error {
	sql := `DELETE FROM user WHERE id = ?`
	res, err := DB.Exec(sql, u.Id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	log.Println("Deleted:", u, rows)
	return nil
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
