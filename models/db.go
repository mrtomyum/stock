package models

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

func NewDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic("sql.Open() Error>>", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Panic("db.Ping() Error>>", err)
		return nil, err
	}
	log.Println("db = ", db)
	return db, nil //return db so in main can call defer db.Close()
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
type Status int
const (
	ACTIVE Status = 1 + iota
	HOLD
	SUSPEND
)

// Base structure contains fields that are common to objects
// returned by the nava's REST API.
type Base struct {
	ID        uint64         `json:"id"`
	CreatedAt mysql.NullTime `json:"created_at"` //todo: change datatype to sql.NullTime
	UpdatedAt mysql.NullTime `json:"updated_at"`
	DeletedAt mysql.NullTime `json:"deleted_at"`
}