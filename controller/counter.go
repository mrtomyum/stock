package controller

import (
	"fmt"
	"encoding/json"
	"net/http"
	log "github.com/Sirupsen/logrus"
	"github.com/mrtomyum/nava-stock/model"
	"github.com/mrtomyum/nava-sys/api"
	"github.com/gin-gonic/gin"
)

func (e *Env) PostNewCounter(ctx *gin.Context) {
	log.Println("call AllMachineBatchSale()")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Origin", "*")

	c := &model.Counter{}
	rs := api.Response{}
	if err := ctx.BindJSON(&c); err != nil {
		log.Println("Decode Error: ", err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, rs)
	} else {
		newCounter, err := c.Insert(e.DB)
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = "CANNOT_INSERT New Counter >>" + err.Error()
			ctx.JSON(http.StatusConflict, rs)
		} else {
			rs.Status = api.SUCCESS
			rs.Data = newCounter
			ctx.JSON(http.StatusOK, rs)
		}
	}
	return
}


func (e *Env) GetAllCounter(w http.ResponseWriter, r *http.Request) {
	log.Println("call AllMachineBatchSale()")
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	c := model.Counter{}
	rs := api.Response{}
	counters, err := c.All(e.DB)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Data = counters
	}
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}

//====================================
// บันทึกเคาทเตอร์ขาย จากหน้าตู้ แบบส่งเป็นชุด
//====================================
//func (e *Env) NewArrayCounter(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Server", "nava Stock")
//	w.Header().Set("Content-Type", "application/json")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//
//	cs := []*model.Counter{}
//	d := json.NewDecoder(r.Body)
//	err := d.Decode(&cs)
//	if err != nil {
//		log.Println("Decode Error: ", err)
//	}
//	rs := api.Response{}
//	newCounters, err := model.NewArrayCounter(e.DB, cs)
//	if err != nil {
//		rs.Status = api.ERROR
//		rs.Message = err.Error()
//		w.WriteHeader(http.StatusConflict)
//	} else {
//		rs.Status = api.SUCCESS
//		rs.Data = newCounters
//		w.WriteHeader(http.StatusOK)
//	}
//	output, _ := json.Marshal(rs)
//	fmt.Fprintf(w, string(output))
//}