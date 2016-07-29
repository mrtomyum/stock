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
	log.Println("call GET AllItem")

	if r.Method != "GET" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

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
	log.Println("call GET ShowItem (by id)")
	v := mux.Vars(r)
	id := v["id"]
	i := new(m.Item)
	i.ID, _ = strconv.ParseUint(id, 10, 64)
	log.Println("item.ID = ", i.ID)

	var item m.Item
	item, err := i.FindItemByID(e.DB)

	rs := new(api.Response)
	if err != nil {
		rs.Status = "204"
		rs.Message = "NOT_FOUND"
		rs.Result = nil
	} else {
		rs.Status = "200"
		rs.Message = "SUCCESS return Item."
		rs.Result = item
	}
	o, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(o))
}

func (e *Env) NewItem(w http.ResponseWriter, r *http.Request) {
	log.Println("call POST NewItem")

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
	err = i.New(e.DB)
	log.Println("i= ", i)
	rs := new(api.Response)
	if err != nil {
		rs.Status = "XXX"
		rs.Message = "CANNOT_UPDATE"
		rs.Result = nil
	} else {
		rs.Status = "200"
		rs.Message = "SUCCESS UPDATE ITEM"
		rs.Result = i
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
