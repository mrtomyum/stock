package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mrtomyum/nava-stock/api"
	m "github.com/mrtomyum/nava-stock/model"
	"log"
	"net/http"
	"strconv"
)

func (e *Env) AllItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET AllItem123")
	if r.Method != "GET" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")

	i := m.Item{}
	items, err := i.All(e.DB)
	rs := m.APIResponse{}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		rs.Status = "500"
		rs.Message = err.Error()
	} else {
		w.WriteHeader(http.StatusFound)
		rs.Status = "302"
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
	log.Println("call GET ShowItem (by id)")
	w.Header().Set("Content-Type", "application/json")

	v := mux.Vars(r)
	id := v["id"]
	i := new(m.Item)
	i.ID, _ = strconv.ParseUint(id, 10, 64)
	log.Println("item.ID = ", i.ID)

	var iv m.ItemView
	iv, err := i.FindItemByID(e.DB)

	rs := new(api.Response)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		rs.Status = "404"
		rs.Message = "NOT_FOUND>> " + err.Error()
		rs.Result = nil
	} else {
		w.WriteHeader(http.StatusFound)
		rs.Status = "302"
		rs.Message = "SUCCESS return Item."
		rs.Result = iv
	}
	o, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(o))
}

func (e *Env) NewItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call POST NewItem")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	i := new(m.Item)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&i)
	if err != nil {
		log.Println("NewItem: Error decode.Decode(&i) >>", err)
	}
	newItem, err := i.New(e.DB)
	log.Println("i= ", i)
	rs := new(api.Response)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		rs.Status = "409"
		rs.Message = "CANNOT_UPDATE >>" + err.Error()
		rs.Result = nil
	} else {
		w.WriteHeader(http.StatusCreated)
		rs.Status = "201"
		rs.Message = "SUCCESS CREATED"
		rs.Result = newItem
	}
	o, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(o))
}

func (e *Env) UpdateItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call UpdateItem")
}

func (e *Env) FindItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call FindItem")
}

func (e *Env) DelItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call NewItem")
}

func (e *Env) UndelItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call NewItem")
}