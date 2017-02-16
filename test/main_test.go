package test

import (
	"github.com/mrtomyum/stock/model"
	"github.com/jmoiron/sqlx"
	"log"
)

var mockDB *sqlx.DB

func init() {
	dsn := model.GetConfig("../model/config_test.json")
	mockDB = sqlx.MustConnect("mysql", dsn)
	log.Println("Connected db: ", mockDB)
}
