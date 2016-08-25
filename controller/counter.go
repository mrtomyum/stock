package controller

import (
	"fmt"
	"encoding/json"
	"net/http"
	"log"
	"github.com/mrtomyum/nava-stock/model"
	"github.com/mrtomyum/nava-sys/api"
)

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

func (e *Env) NewCounter(w http.ResponseWriter, r *http.Request) {
	log.Println("call AllMachineBatchSale()")
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	c := &model.Counter{}
	rs := api.Response{}
	counter, err := c.NewCounter(e.DB)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Data = counter
	}
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}

func (e *Env) NewArrayCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sales := []*model.Counter{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&sales)
	if err != nil {
		log.Println("Decode Error: ", err)
	}
	rs := api.Response{}
	newBS, err := model.NewArrayCounter(e.DB, sales)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		w.WriteHeader(http.StatusConflict)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = newBS
		w.WriteHeader(http.StatusOK)
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}