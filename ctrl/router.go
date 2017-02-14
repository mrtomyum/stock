package ctrl

import "github.com/gin-gonic/gin"

func Router() *gin.Engine {
	r := gin.Default()
	itemV1 := r.Group("/v1/item")
	{
		itemV1.GET("/", GetAllItem)
		itemV1.GET("/:id", GetItem)
		itemV1.POST("/", PostNewItem)
		itemV1.PUT("/", UpdateItem)
		//itemV1.GET("/:id/prices", c.GetItemPriceByID)
	}

	machineV1 := r.Group("/v1/machine")
	{
		machineV1.POST("/", PostNewMachine)
		machineV1.GET("/", GetAllMachines)
		machineV1.GET("/:id", GetThisMachine)
		machineV1.GET("/:id/columns", GetMachineColumns)
		machineV1.GET("/:id/template", GetMachineTemplate)
	}

	columnV1 := r.Group("/v1/column")
	{
		columnV1.PUT("/:id", PutMachineColumn)
	}

	counterV1 := r.Group("/v1/counter/")
	{
		counterV1.GET("/", GetAllCounter)
		counterV1.POST("/", PostCounter)

		//counterV1.GET("/:id/", GetCounter)
		//counterV1.PUT("/:id", PutCounter)
		//counterV1.DELETE("/:id", DeleteCounter)
		//counterByMachineCode := r.Group("/machineCode")
		//{
		counterV1.GET("/machineCode/:code", GetCounterByMachineCode)
		//}
		//counterByRouteMan := r.Group("/routeman")
		//{
		//	counterByRouteMan.GET("/:id", GetByRoutemanId)
		//}
	}

	locationV1 := r.Group("/v1/location")
	{
		locationV1.GET("/", GetAllLocationTree)
		locationV1.GET("/:id", GetLocationTreeByID)
		locationV1.POST("/", PostNewLocation)
	}

	stockV1 := r.Group("/v1/stock")
	{
		stockV1.GET("/", GetAllStock)
	}

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
