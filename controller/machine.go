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
	log.Println("call AllMachineBatchSale()")
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