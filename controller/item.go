package controller

import (
	"github.com/mrtomyum/sys/api"
	m "github.com/mrtomyum/stock/model"
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func (e *Env) GetAllItem(c *gin.Context) {
	log.Println("call GET AllItem")
	c.Header("Server", "NAVA Stock")
	c.Header("Host", "api.nava.work:8001")
	c.Header("Access-Control-Allow-Origin", "*")

	i := m.Item{}
	items, err := i.GetAll()
	rs := api.Response{}
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = items
		c.JSON(http.StatusFound, rs)
	}
}

func (e *Env) PostNewItem(c *gin.Context) {
	log.Println("call POST NewItem")
	c.Header("Server", "NAVA Stock")
	c.Header("Access-Control-Allow-Origin", "*")

	i := new(m.Item)
	rs := new(api.Response)
	if err := c.BindJSON(&i); err != nil {
		log.Println("NewItem: Error decode.Decode(&i) >>", err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
	} else {
		newItem, err := i.Insert()
		log.Println("i= ", i)
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = "CANNOT_UPDATE >>" + err.Error()
			c.JSON(http.StatusConflict, rs)
		} else {
			rs.Status = api.SUCCESS
			rs.Data = newItem
			c.JSON(http.StatusOK, rs)
		}
	}
	return
}

func (e *Env) GetItem(c *gin.Context) {
	log.Println("call FindItem")
	c.Header("Server", "NAVA Stock")
	c.Header("Access-Control-Allow-Origin", "*")
	var i m.Item
	id := c.Param("id")
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
	c.JSON(200, rs)
	return
}

func (e *Env) UpdateItem(c *gin.Context) {
	log.Println("call UpdateItem")
	c.Header("Server", "NAVA Stock")
	c.Header("Access-Control-Allow-Origin", "*")
	var i m.Item
	rs := api.Response{}
	if c.BindJSON(&i) != nil {
		c.JSON(http.StatusBadRequest, i)
	} else {
		updatedItem, err := i.Update()
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = err.Error()
		} else {
			rs.Status = api.SUCCESS
			rs.Data = updatedItem
		}
		c.JSON(http.StatusOK, rs)
	}
}

func (e *Env) DelItem(c *gin.Context) {
	log.Println("call NewItem")
}

func (e *Env) UndelItem(c *gin.Context) {
	log.Println("call NewItem")
}