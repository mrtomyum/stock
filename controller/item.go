package controller

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	m "github.com/mrtomyum/nava-stock/model"
)

func (e *Env) AllItem(w http.ResponseWriter, r *http.Request) {
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

func (e *Env) ShowItem(w http.ResponseWriter, r *http.Request) {

}

func (e *Env) NewItem(w http.ResponseWriter, r *http.Request) {

}

func (e *Env) UpdateItem(w http.ResponseWriter, r *http.Request) {

}

func (e *Env) FindItem(w http.ResponseWriter, r *http.Request) {

}

func (e *Env) DelItem(w http.ResponseWriter, r *http.Request) {

}

func (e *Env) UndelItem(w http.ResponseWriter, r *http.Request) {

}
