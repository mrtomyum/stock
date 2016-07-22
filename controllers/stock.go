package controllers

import (
	"net/http"
	"github.com/mrtomyum/nava-stock/models"
	"fmt"
	"log"
	"encoding/json"
)

func (e *Env) ItemIndex(w http.ResponseWriter, r *http.Request){
	i := models.Item{}
	items, err := i.Index(e.DB)
	rs := models.APIResponse{}
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

func (e *Env) ItemShow(w http.ResponseWriter, r *http.Request){

}

func (e *Env) ItemInsert(w http.ResponseWriter, r *http.Request){

}

func (e *Env) ItemUpdate(w http.ResponseWriter, r *http.Request){

}

func (e *Env) ItemDelete(w http.ResponseWriter, r *http.Request){

}

func (e *Env) ItemSearch(w http.ResponseWriter, r *http.Request){

}

func (e *Env) ItemUndelete(w http.ResponseWriter, r *http.Request) {

}