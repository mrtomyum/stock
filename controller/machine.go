package controller

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/mrtomyum/nava-sys/api"
	m "github.com/mrtomyum/nava-stock/model"
	"log"
)

func (e *Env) AllMachine(w http.ResponseWriter, r *http.Request) {
	log.Println("call AllMachine()")
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	m := m.Machine{}
	rs := api.Response{}
	machines, err := m.All(e.DB)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		w.WriteHeader(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = machines
		w.WriteHeader(http.StatusOK)
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}

func (e *Env) NewMachine(w http.ResponseWriter, r *http.Request) {
	log.Println("call controller.NewMachine()")
	w.Header().Set("Server", "nava Stock")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	m := m.Machine{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&m)
	log.Println(m)
	rs := api.Response{}
	newMachine, err := m.New(e.DB)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Data = newMachine
	}
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(output))
}