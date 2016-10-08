package model

import (
	"log"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func NewDB(dsn string) (*sqlx.DB) {
	DB = sqlx.MustConnect("mysql", dsn)
	log.Println("Connected db: ", DB)
	return DB //return db so in main can call defer db.Close()
}
