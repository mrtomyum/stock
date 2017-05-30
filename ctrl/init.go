package ctrl

import (
	"github.com/jmoiron/sqlx"
	"github.com/mrtomyum/stock/model"
	"log"
)

var db *sqlx.DB
//var rs *Response

func init() {
	// Read configuration file from "config.json"
	//dsn := GetConfig("./model/config.json") // เปิดใช้งานจริงเมื่อ Docker Container run --link ตรงเข้า mariadb เท่านั้น
	dsn := model.GetConfig("./model/config.json")
	db = sqlx.MustConnect("mysql", dsn)
	log.Println("Connected db: ", db)
}

