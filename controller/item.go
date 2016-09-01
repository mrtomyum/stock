package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mrtomyum/nava-sys/api"
	m "github.com/mrtomyum/nava-stock/model"
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func (e *Env) AllItem(c *gin.Context) {
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

func (e *Env) GetItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET ShowItem (by id)")
	w.Header().Set("Content-Type", "application/json")

	v := mux.Vars(r)
	id := v["id"]
	i := new(m.Item)
	i.ID, _ = strconv.ParseUint(id, 10, 64)
	log.Println("item.ID = ", i.ID)

	//var iv *m.ItemView
	rs := new(api.Response)
	iv, err := i.FindItemByID(e.DB)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		rs.Status = api.ERROR
		rs.Message = "NOT_FOUND>> " + err.Error()
	} else {
		w.WriteHeader(http.StatusFound)
		rs.Status = api.SUCCESS
		rs.Data = iv
	}
	o, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(o))
}

func (e *Env) NewItem(c *gin.Context) {
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
		newItem, err := i.New(e.DB)
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

func (e *Env) UpdateItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call UpdateItem")
}

func (e *Env) FindItemByID(c *gin.Context) {
	log.Println("call FindItem")
	c.Header("Server", "NAVA Stock")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	id, _ := strconv.Atoi(c.Param("id"))
	//id := 1
	var i m.Item
	i.ID = uint64(id)
	rs := api.Response{}
	iv, err := i.FindItemByID(e.DB)
	log.Println("return from FindItemByID()")
	if err != nil {
		log.Println(err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Data = iv
		c.JSON(200, rs)
	}
	return
}

func (e *Env) DelItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call NewItem")
}

func (e *Env) UndelItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call NewItem")
}
