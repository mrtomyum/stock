package test

import (
	"testing"
	"net/http"
	"encoding/json"
	"github.com/mrtomyum/nava-sys/api"
	"time"
	"strings"
)

var url string = "http://localhost:8001"
var jsonStr string = `{
	"access_token":"zx;lja;ldf",
  	"job_hash":"3235a5b3de8",
  	"job":{
		"type":"counter",
		"recorded": "2016-09-05T12:00:00:00Z07:00",
		"machine_id":1,
		"counter_sum":45,
		"columns":
		[
		  {
			"column_no":1,
			"item_id":233,
			"price":15,
			"last_counter":0,
			"curr_counter":10
		  },
		  {
			"no":2,
			"item_id":42,
			"price":10,
			"last_counter":0,
			"curr_counter":8
		  },
		  {
			"no":3,
			"item_id":33,
			"price":12,
			"last_counter":3,
			"curr_counter":15
		  }
		]
  	}
}`

func Test_PostCounterAPI(t *testing.T) {
	client := &http.Client{Timeout:time.Second * 10}
	b := strings.NewReader(jsonStr)
	res, err := client.Post(
		"http://localhost:8001/counters",
		"application/json",
		b)
	if err != nil {
		t.Error(err)
	}
	var r api.Response
	d := json.NewDecoder(res.Body)
	err = d.Decode(&r)
	if err != nil {
		t.Error("Error in Decode response body.")
	}
	if r.Status != api.SUCCESS {
		t.Error("Expected SUCCESS got", r.Status)
		t.Fail()
	}
}

func init() {

}