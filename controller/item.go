package controller

import (
	"github.com/mrtomyum/sys/api"
	m "github.com/mrtomyum/stock/model"
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func GetAllItem(ctx *gin.Context) {
	log.Println("call GET AllItem")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Host", "api.nava.work:8001")
	ctx.Header("Access-Control-Allow-Origin", "*")

	i := m.Item{}
	items, err := i.GetAll()
	rs := api.Response{}
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, rs)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = items
		ctx.JSON(http.StatusFound, rs)
	}
}

func PostNewItem(ctx *gin.Context) {
	log.Println("call POST NewItem")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")

	i := new(m.Item)
	rs := new(api.Response)
	if err := ctx.BindJSON(&i); err != nil {
		log.Println("NewItem: Error decode.Decode(&i) >>", err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, rs)
	} else {
		newItem, err := i.Insert()
		log.Println("i= ", i)
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = "CANNOT_UPDATE >>" + err.Error()
			ctx.JSON(http.StatusConflict, rs)
		} else {
			rs.Status = api.SUCCESS
			rs.Data = newItem
			ctx.JSON(http.StatusOK, rs)
		}
	}
	return
}

func GetItem(ctx *gin.Context) {
	log.Println("call FindItem")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")
	var i m.Item
	id := ctx.Param("id")
	i.ID, _ = strconv.ParseUint(id, 10, 64)
	rs := api.Response{}
	iv, err := i.GetItemView()
	log.Println("return from GetItemView()")
	if err != nil {
		log.Println(err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Data = iv
	}
	ctx.JSON(200, rs)
	return
}

func UpdateItem(ctx *gin.Context) {
	log.Println("call UpdateItem")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")
	var i m.Item
	rs := api.Response{}
	if ctx.BindJSON(&i) != nil {
		ctx.JSON(http.StatusBadRequest, i)
	} else {
		updatedItem, err := i.Update()
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = err.Error()
		} else {
			rs.Status = api.SUCCESS
			rs.Data = updatedItem
		}
		ctx.JSON(http.StatusOK, rs)
	}
}

func DelItem(ctx *gin.Context) {
	log.Println("call NewItem")
}

func UndelItem(ctx *gin.Context) {
	log.Println("call NewItem")
}