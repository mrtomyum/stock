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

type APIResponse struct {
	Status  string
	Message string
	Result  interface{}
}

// Structure for collection of search string for frontend request.
type APISearch struct {
	Name string
}