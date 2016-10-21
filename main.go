package main

import (
	//"net/http"
	log "github.com/Sirupsen/logrus"
	c "github.com/mrtomyum/stock/controller"
	"github.com/mrtomyum/stock/model"
	"encoding/json"
	"os"
	"github.com/sebest/logrusly"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/mrtomyum/sys/api"
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

//func NewDB(dsn string) (*sqlx.DB) {
//	db := sqlx.MustConnect("mysql", dsn)
//	log.Println("Connected db: ", db)
//	return db //return db so in main can call defer db.Close()
//}

func SetupRoute(c *c.Env) *gin.Engine {
	r := gin.Default()
	itemV1 := r.Group("/v1/items")
	{
		itemV1.GET("/", c.GetAllItem)
		itemV1.GET("/:id", c.GetItem)
		itemV1.POST("/", c.PostNewItem)
		itemV1.PUT("/", c.UpdateItem)
		//itemV1.GET("/:id/prices", c.GetItemPriceByID)
	}

	machineV1 := r.Group("/v1/machines")
	{
		machineV1.POST("/", c.PostNewMachine)
		machineV1.GET("/", c.GetAllMachines)
		machineV1.GET("/:id", c.GetThisMachine)
		machineV1.GET("/:id/columns", c.GetMachineColumns)
		machineV1.GET("/:id/templates", GetMachineTemplate)
	}

	columnV1 := r.Group("/v1/columns")
	{
		columnV1.PUT("/:id", c.PutMachineColumn)
	}

	counterV1 := r.Group("/v1/counters")
	{
		counterV1.POST("/", c.PostCounter)
		counterV1.GET("/", c.GetAllCounter)
		counterV1.GET("/:id", c.GetCounter)
		counterV1.PUT("/:id", c.PutCounter)
		counterV1.DELETE("/:id", c.DeleteCounter)
	}

	//locationV1 := r.Group("/v1/locations")
	//{
	//	locationV1.GET("/", c.GetAllLocationTree)
	//	locationV1.GET("/:id", c.GetLocationTreeByID)
	//	locationV1.POST("/", c.PostNewLocation)
	//}

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
	db := model.NewDB(dsn)
	c := &c.Env{DB:db}
	defer db.Close()

	r := SetupRoute(c)

	r.Run(":8001")
}

func GetMachineTemplate(c *gin.Context) {
	var m *model.Machine
	rs := api.Response{}
	templates, err := m.GetTemplate()
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
	}
	rs.Status = api.SUCCESS
	rs.Self = "api.nava.work/v1/machine/template"
	rs.Data = templates
	c.JSON(http.StatusOK, rs)
}
