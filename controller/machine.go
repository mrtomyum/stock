package controller

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/mrtomyum/nava-sys/api"
	m "github.com/mrtomyum/nava-stock/model"
	"log"
)

func (e *Env) AllMachineBatchSale(w http.ResponseWriter, r *http.Request) {
	log.Println("call AllMachineBatchSale()")
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rs := api.Response{}
	mbs := m.BatchSale{}
	sales, err := mbs.All(e.DB)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		w.WriteHeader(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = sales
		w.WriteHeader(http.StatusOK)
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}

func (e *Env) NewMachineBatchSale(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sales := []*m.BatchSale{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&sales)
	if err != nil {
		log.Println("Decode Error: ", err)
	}
	rs := api.Response{}
	err = m.NewBatchSale(e.DB, sales)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		w.WriteHeader(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = sales
		w.WriteHeader(http.StatusOK)
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}