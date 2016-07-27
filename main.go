package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	c "github.com/mrtomyum/nava-stock/controller"
	m "github.com/mrtomyum/nava-stock/model"
	"encoding/json"
	"os"
)

type Config struct {
	DBHost string `json:"db_host"`
	DBName string `json:"db_name"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
}

func main() {
	// Read configuration file from "cofig.json"
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Println("error:", err)
	}

	var dsn = config.DBUser + ":" + config.DBPass + "@" + config.DBHost + "/" + config.DBName + "?parseTime=true"

	db := m.NewDB(dsn)
	c := c.Env{DB:db}
	defer db.Close()

	r := mux.NewRouter().StrictSlash(true)

	// ## Item
	s := r.PathPrefix("/api/v1/item").Subrouter()
	s.HandleFunc("/", c.AllItem).Methods("GET"); log.Println("/api/v1/item GET AllItem")
	s.HandleFunc("/", c.NewItem).Methods("POST"); log.Println("/api/v1/item POST NewItem")
	s.HandleFunc("/{id:[0-9]+}", c.ShowItem).Methods("GET"); log.Println("/api/v1/item/:id GET ShowItem")
	s.HandleFunc("/{id:[0-9]+}", c.UpdateItem).Methods("PUT"); log.Println("/api/v1/item/:id PUT UpdateItem ")
	s.HandleFunc("/search", c.FindItem).Methods("POST"); log.Println("/api/v1/item/search POST FindItem")
	s.HandleFunc("/{id:[0-9]+}", c.DelItem).Methods("DELETE"); log.Println("/api/v1/item/:id DELETE ItemDelete")
	s.HandleFunc("/{id:[0-9]+}/undelete", c.UndelItem).Methods("PUT"); log.Println("/api/v1/item/:id/undelete PUT ItemUndelete")

	// # Stock


	// ## Location

	http.Handle("/", r)
	http.ListenAndServe(":8001", nil)
}
