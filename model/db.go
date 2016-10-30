package model

import (
	"log"
	"github.com/jmoiron/sqlx"
	"encoding/json"
	"os"
)

var DB *sqlx.DB

type Config struct {
	DBHost string `json:"db_host"`
	DBName string `json:"db_name"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
	Port   string `json:"port"`
}

func getConfig(fileName string) string {
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
	dsn := getConfig("./config.json")
	DB = sqlx.MustConnect("mysql", dsn)
	log.Println("Connected db: ", DB)
}