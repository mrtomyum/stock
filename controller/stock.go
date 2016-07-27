package controller

import (
	"net/http"
	m "github.com/mrtomyum/nava-stock/model"
	"fmt"
	"log"
	"encoding/json"
)

func (e *Env) AllStock(w http.ResponseWriter, r *http.Request) {
	i := m.Item{}
	items, err := i.All(e.DB)
	rs := m.APIResponse{}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		rs.Status = "500"
		rs.Message = err.Error()
	} else {
		rs.Status = "200"
		rs.Message = "SUCCESS"
		rs.Result = items
	}
	output, err := json.Marshal(rs)
	if err != nil {
		log.Println("Error json.Marshal:", err)
	}
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
