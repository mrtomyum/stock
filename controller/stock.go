package controller

import (
	"net/http"
	m "github.com/mrtomyum/stock/model"
	"fmt"
	"encoding/json"
	"github.com/mrtomyum/sys/api"
)

func (e *Env) AllStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	s := m.Stock{}
	rs := api.Response{}
	items, err := s.GetAll()
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		w.WriteHeader(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = items
		w.WriteHeader(http.StatusOK)
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}

func (e *Env) FindStockByID(w http.ResponseWriter, r *http.Request) {

}

func (e *Env) NewStock(w http.ResponseWriter, r *http.Request) {

}

func (e *Env) UpdateStock(w http.ResponseWriter, r *http.Request) {

}

func (e *Env) DelStock(w http.ResponseWriter, r *http.Request) {

}

func (e *Env) UndelStock(w http.ResponseWriter, r *http.Request) {

}
