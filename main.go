package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/mrtomyum/nava-stock/controllers"
	"github.com/mrtomyum/nava-stock/models"
)

//TODO: Move Config to JSON file and Create Config{} to handle DB const.
const (
	DB_HOST = "tcp(nava.work:3306)"
	//TODO: เมื่อรันจริงต้องเปลี่ยนเป็น Docker Network Bridge IP เช่น 172.17.0.3 เป็นต้น
	DB_NAME = "system"
	DB_USER = "root"
	DB_PASS = "mypass"
)

var dsn = DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"


func main() {
	db, err := models.NewDB(dsn)
	if err != nil {
		log.Panic("NewDB() Error:", err)
	}

	c := controllers.Env{DB:db}
	defer db.Close()

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1/user").Subrouter()
	s.HandleFunc("/", c.ItemIndex).Methods("GET"); log.Println("/api/v1/item GET ItemIndex")
	s.HandleFunc("/", c.ItemInsert).Methods("POST"); log.Println("/api/v1/item POST ItemInsert")
	s.HandleFunc("/{id:[0-9]+}", c.ItemShow).Methods("GET"); log.Println("/api/v1/item/:id GET ItemShow")
	s.HandleFunc("/{id:[0-9]+}", c.ItemUpdate).Methods("PUT"); log.Println("/api/v1/item/:id PUT ItemUpdate ")
	s.HandleFunc("/search", c.ItemSearch).Methods("POST"); log.Println("/api/v1/item/search POST ItemSearch")
	s.HandleFunc("/{id:[0-9]+}", c.ItemDelete).Methods("DELETE"); log.Println("start '/api/v1/item/:id' DELETE UserDelete")
	s.HandleFunc("/{id:[0-9]+}/undelete", c.ItemUndelete).Methods("PUT"); log.Println("start '/api/v1/item/:id/undelete' PUT UserUndelete")
	// Menu
	// # Stock

	// ## Item

	// ## Location

	http.Handle("/", r)
	http.ListenAndServe(":8001", nil)
}
