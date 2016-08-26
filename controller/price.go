package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mrtomyum/nava-sys/api"
	"github.com/mrtomyum/nava-stock/model"
	"log"
	"net/http"
	"strconv"
)

func (e *Env) GetItemPriceByID(w http.ResponseWriter, r *http.Request) {
	log.Println("call ItemPrice")
	v := mux.Vars(r)
	id := v["id"]
	ip := new(model.ItemPrice)
	ips := []*model.ItemPrice{}
	ip.ItemID, _ = strconv.ParseUint(id, 10, 64)
	log.Println("item.ID = ", ip.ItemID)

	ips, err := ip.All(e.DB)

	rs := new(api.Response)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = "NO CONTENT"
		w.WriteHeader(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = ips
		w.WriteHeader(http.StatusOK)
	}
	o, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(o))
}
//=================
// บันทึกราคาจากหน้าตู้
//=================

func (e *Env) NewBatchPrice(w http.ResponseWriter, r *http.Request) {
	log.Println("call NewBatchPrice()")
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	bp := &model.BatchPrice{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&bp)

	rs := api.Response{}
	newPrice, err := bp.New(e.DB)
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

func (e *Env) AllBatchPrice(w http.ResponseWriter, r *http.Request) {
	log.Println("call AllMachineBatchSale()")
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rs := api.Response{}
	p := model.BatchPrice{}
	prices, err := p.All(e.DB)
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