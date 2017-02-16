package ctrl

import (
	"encoding/json"
	"fmt"
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

	rs := Response{}
	newPrice, err := bp.New(db)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		w.WriteHeader(http.StatusNoContent)
	} else {
		rs.Status = SUCCESS
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

	rs := Response{}
	p := model.BatchPrice{}
	prices, err := p.GetAll(db)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		w.WriteHeader(http.StatusNoContent)
	} else {
		rs.Status = SUCCESS
		rs.Data = prices
		w.WriteHeader(http.StatusOK)
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}