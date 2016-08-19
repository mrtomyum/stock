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
	Port   string `json:"port"`
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
	c := &c.Env{DB:db}
	defer db.Close()

	r := SetupRoute(c)

	http.Handle("/", r)
	http.ListenAndServe(":8001", nil)
}

func SetupRoute(c *c.Env) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// # Item
	r.HandleFunc("/v1/items", c.AllItem).Methods("GET"); log.Println("/v1/items GET AllItem")
	r.HandleFunc("/v1/items/{id:[0-9]+}", c.ShowItem).Methods("GET"); log.Println("/v1/items/:id GET ShowItem")
	r.HandleFunc("/v1/items", c.NewItem).Methods("POST"); log.Println("/v1/items POST NewItem")
	r.HandleFunc("/v1/locations/{id:[0-9]+}", c.LocationTreeByID).Methods("GET"); log.Println("/v1/locations/:id GET Location tree by ID")
	r.HandleFunc("/v1/locations", c.LocationTreeAll).Methods("GET"); log.Println("/v1/locations GET All Location tree")
	r.HandleFunc("/v1/locations", c.NewLocation).Methods("POST"); log.Println("/v1/locations POST New Location")


	//s.HandleFunc("/{id:[0-9]+}", c.UpdateItem).Methods("PUT"); log.Println("/api/v1/item/:id PUT UpdateItem ")
	//s.HandleFunc("/search", c.FindItem).Methods("POST"); log.Println("/api/v1/item/search POST FindItem")
	//s.HandleFunc("/{id:[0-9]+}", c.DelItem).Methods("DELETE"); log.Println("/api/v1/item/:id DELETE ItemDelete")
	//s.HandleFunc("/{id:[0-9]+}/undelete", c.UndelItem).Methods("PUT"); log.Println("/api/v1/item/:id/undelete PUT ItemUndelete")
	// ## ItemPrice
	//s.HandleFunc("/{id:[0-9]+}/price", c.ItemPrice).Methods("GET"); log.Println("/api/v1/item/:id/price GET PriceByItemID")
	// # Stock
	s := r.PathPrefix("/v1/stocks").Subrouter()
	s.HandleFunc("/", c.AllStock).Methods("GET"); log.Println("/v1/stocks/")

	// ## Location
	// ## Machine
	s = r.PathPrefix("/v1/machines").Subrouter()
	s.HandleFunc("/", c.AllMachine).Methods("GET"); log.Println("/v1/machines/ GET AllMachine")
	//s.HandleFunc("/", c.NewMachine).Methods("POST"); log.Println("/v1/machines/ POST NewMachine")

	// ## Batch
	s = r.PathPrefix("/v1/batchs/").Subrouter()
	s.HandleFunc("/sales", c.AllBatchSale).Methods("GET"); log.Println("/v1/machines/batchSales GET ALlMachineBatchSale")
	s.HandleFunc("/sales", c.NewBatchSale).Methods("POST"); log.Println("/v1/machines/batchSales POST NewMachineBatchSale")
	s.HandleFunc("/prices", c.AllBatchPrice).Methods("POST"); log.Println("/v1/machines/batchSales POST NewMachineBatchSale")

	return r
}