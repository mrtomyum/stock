package controller

import (
	m "github.com/mrtomyum/nava-stock/model"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
)

func CreateLocationTree(locations []*m.Location) *m.Location {
	tree := new(m.Location)
	for _, l := range locations {
		tree.Add(l)
	}
	return tree
}

func (e *Env) ShowLocationTree(w http.ResponseWriter, r *http.Request) {
	log.Println("call AllLocationTree()")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	v := mux.Vars(r)
	id := v["id"]
	locations, err := m.ShowLocations(id)
	if err != nil {
		log.Fatal("ShowLocations()", err)
		w.WriteHeader(http.StatusNotFound)
	}
	tree := CreateLocationTree(locations)
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(tree.Child)
	fmt.Fprintf(w, string(output))
}

