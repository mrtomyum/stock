package controller

import (
	"encoding/json"
	"fmt"
	"github.com/mrtomyum/sys/api"
	"github.com/mrtomyum/stock/model"
	"log"
	"net/http"
)


//=================
// บันทึกราคาจากหน้าตู้
//=================

func NewBatchPrice(w http.ResponseWriter, r *http.Request) {
	log.Println("call NewBatchPrice()")
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	bp := &model.BatchPrice{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&bp)

	rs := api.Response{}
	newPrice, err := bp.New()
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		w.WriteHeader(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = newPrice
		w.WriteHeader(http.StatusOK)
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}

func AllBatchPrice(w http.ResponseWriter, r *http.Request) {
	log.Println("call AllMachineBatchSale()")
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rs := api.Response{}
	p := model.BatchPrice{}
	prices, err := p.GetAll()
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		w.WriteHeader(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = prices
		w.WriteHeader(http.StatusOK)
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}