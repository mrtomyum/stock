package test

import (
	"testing"
	"github.com/jmoiron/sqlx"
	"github.com/mrtomyum/stock/model"
	"encoding/json"
	"os"
	"log"
)

var url string = "http://localhost:8001"

const (
	DB_HOST = "tcp(nava.work:3306)"
	DB_NAME = "test_stock"
	DB_USER = "root"
	DB_PASS = "mypass"
)

var db *sqlx.DB

func init() {
	var testDSN = DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?parseTime=true"
	db = sqlx.MustConnect("mysql", testDSN)
}

func Test_ModelCounterInsert(t *testing.T) {
	var c model.Counter
	file, err := os.Open("counter.json")
	if err != nil {
		t.Error(err)
	}
	d := json.NewDecoder(file)
	err = d.Decode(&c) // Todo: decoding unsuccessful
	if err != nil {
		t.Error(err)
	}
	log.Println("Decoded JSON to Struct: ", c)
	newCounter, err := c.Insert(db)
	if err != nil {
		t.Error(err.Error())
	}
	log.Println("Inserted Counter: ", newCounter)
}

//func Test_PostCounterAPI(t *testing.T) {
//	client := &http.Client{Timeout:time.Second * 10}
//	file, err := os.Open("counter.json")
//	b := strings.NewReader(file)
//	res, err := client.Post(
//		"http://localhost:8001/counters",
//		"application/json",
//		b)
//	if err != nil {
//		t.Error(err)
//	}
//	var r api.Response
//	d := json.NewDecoder(res.Body)
//	err = d.Decode(&r)
//	if err != nil {
//		t.Error("Error in Decode response body.")
//	}
//	if r.Status != api.SUCCESS {
//		t.Error("Expected SUCCESS got", r.Status)
//		t.Fail()
//	}
//	return
//}