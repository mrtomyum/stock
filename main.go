package main

import (
	//"net/http"
	log "github.com/Sirupsen/logrus"
	c "github.com/mrtomyum/nava-stock/controller"
	"encoding/json"
	"os"
	"github.com/jmoiron/sqlx"
	"github.com/sebest/logrusly"
	"github.com/gin-gonic/gin"
)

var logglyToken string = "4cd7bdfb-0345-4205-aeee-53e85a030eda"

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

func NewDB(dsn string) (*sqlx.DB) {
	db := sqlx.MustConnect("mysql", dsn)
	log.Println("Connected db: ", db)
	return db //return db so in main can call defer db.Close()
}

func SetupRoute(c *c.Env) *gin.Engine {
	r := gin.Default()
	machineV1 := r.Group("/v1/machines")
	{
		machineV1.GET("/", c.AllMachine)
		machineV1.POST("/", c.NewMachine)
	}

	itemV1 := r.Group("/v1/items")
	{
		itemV1.GET("/", c.AllItem)
		itemV1.GET("/:id", c.GetItemByID)
		itemV1.POST("/", c.NewItem)
		itemV1.PUT("/", c.UpdateItem)
	}

	// ## Location
	//s := r.PathPrefix("/v1/locations").Subrouter()
	//r.HandleFunc("/", c.GetAllLocationTree).Methods("GET"); log.Println("/v1/locations GET All Location tree")
	//r.HandleFunc("/", c.NewLocation).Methods("POST"); log.Println("/v1/locations POST New Location")
	//r.HandleFunc("/{id:[0-9]+}", c.GetLocationTreeByID).Methods("GET"); log.Println("/v1/locations/:id GET Location tree by ID")


	//s.HandleFunc("/{id:[0-9]+}", c.UpdateItem).Methods("PUT"); log.Println("/api/v1/item/:id PUT UpdateItem ")
	//s.HandleFunc("/search", c.FindItem).Methods("POST"); log.Println("/api/v1/item/search POST FindItem")
	//s.HandleFunc("/{id:[0-9]+}", c.DelItem).Methods("DELETE"); log.Println("/api/v1/item/:id DELETE ItemDelete")
	//s.HandleFunc("/{id:[0-9]+}/undelete", c.UndelItem).Methods("PUT"); log.Println("/api/v1/item/:id/undelete PUT ItemUndelete")
	// ## ItemPrice
	//s.HandleFunc("/{id:[0-9]+}/price", c.ItemPrice).Methods("GET"); log.Println("/api/v1/item/:id/price GET PriceByItemID")
	// # Stock
	//s = r.PathPrefix("/v1/stocks").Subrouter()
	//s.HandleFunc("/", c.AllStock).Methods("GET"); log.Println("/v1/stocks/")


	// ## Batch
	//s = r.PathPrefix("/v1/batchs/counters/").Subrouter()
	//s.HandleFunc("/", c.GetAllCounter).Methods("GET"); log.Println("/v1/machines/batchSales GET All Batch Sale")
	//s.HandleFunc("/", c.NewCounter).Methods("POST"); log.Println("/v1/machines/batchSales POST New Batch Sale")
	//s.HandleFunc("/", c.NewArrayCounter).Methods("POST"); log.Println("/v1/machines/batchSales POST New Batch Array Sale")
	//
	//s = r.PathPrefix("/v1/batchs/prices/").Subrouter()
	//s.HandleFunc("/", c.AllBatchPrice).Methods("POST"); log.Println("/v1/machines/batchSales POST New Batch Price")
	//
	//s = r.PathPrefix("/v1/fulfill/").Subrouter()
	////s.HandleFunc("/", c.NewFulfill).Methods("POST")
	//s.HandleFunc("/", c.GetAllFulfill).Methods("POST")
	//s.HandleFunc("/{id:[0-9+]}", c.GetFulfillByID).Methods("POST")
	return r
}

func main() {
	// Log
	logrus := log.New()
	hook := logrusly.NewLogglyHook(logglyToken, "http://logs-01.loggly.com/inputs/", log.InfoLevel, "info")
	logrus.Hooks.Add(hook)
	defer hook.Flush()
	log.WithFields(log.Fields{
		"name": "Tom NAVA Stock",
	}).Info("Start Logrus")
	// Read configuration file from "cofig.json"
	dsn := getConfig("config.json")
	db := NewDB(dsn)
	c := &c.Env{DB:db}
	defer db.Close()

	r := SetupRoute(c)

	r.Run(":8001")
	//http.Handle("/", r)
	//http.ListenAndServe(":8001", nil)
}
