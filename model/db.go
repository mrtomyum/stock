package model

import (
	"log"
	"github.com/jmoiron/sqlx"
)

func NewDB(dsn string) (*sqlx.DB) {
	db := sqlx.MustConnect("mysql", dsn)
	log.Println("Connected db: ", db)
	return db //return db so in main can call defer db.Close()
}
