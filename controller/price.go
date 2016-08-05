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

func (e *Env) ItemPrice(w http.ResponseWriter, r *http.Request) {
	log.Println("call ItemPrice")
	v := mux.Vars(r)
	id := v["id"]
	ip := new(m.ItemPrice)
	ips := []*m.ItemPrice{}
	ip.ItemID, _ = strconv.ParseUint(id, 10, 64)
	log.Println("item.ID = ", ip.ItemID)

	ips, err := ip.AllPrice(e.DB)

	rs := new(api.Response)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = "NO CONTENT"
		w.WriteHeader(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = ips
		w.WriteHeader(http.StatusOK)
	}
	o, _ := json.Marshal(rs)
	fmt.Fprintf(w, string(o))
}
