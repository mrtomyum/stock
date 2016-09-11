package controller

import (

	"github.com/mrtomyum/nava-sys/api"
	m "github.com/mrtomyum/nava-stock/model"
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func (e *Env) GetAllItem(c *gin.Context) {
	log.Println("call GET AllItem")
	c.Header("Server", "NAVA Stock")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	i := m.Item{}
	items, err := i.All(e.DB)
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
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	i := new(m.Item)
	rs := new(api.Response)
	if err := c.BindJSON(&i); err != nil {
		log.Println("NewItem: Error decode.Decode(&i) >>", err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
	} else {
		newItem, err := i.Insert(e.DB)
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

func (e *Env) GetItemByID(c *gin.Context) {
	log.Println("call FindItem")
	c.Header("Server", "NAVA Stock")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	var i m.Item
	//id, _ := strconv.Atoi(c.Param("id"))
	id := c.Param("id")
	i.ID, _ = strconv.ParseUint(id, 10, 64)
	rs := api.Response{}
	iv, err := i.GetItemView(e.DB)
	log.Println("return from FindItemByID()")
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
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	var i m.Item
	rs := api.Response{}
	if c.BindJSON(&i) != nil {
		c.JSON(http.StatusBadRequest, i)
	} else {
		updatedItem, err := i.Update(e.DB)
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

func (e *Env) DelItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call NewItem")
}

func (e *Env) UndelItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call NewItem")
}
