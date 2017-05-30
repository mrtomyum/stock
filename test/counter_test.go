package test

import (
	"testing"
)

func Test_CounterInsert(t *testing.T) {
	for _, counter := range mockCounter {
		newCounter, err := counter.Insert(mockDB)
		if err != nil {
			t.Error(err)
		}
		println("Inserted counter id: ", newCounter.Id)
	}
	t.Log("Success Insert Counter to DB")
	for _, counter := range mockCounter {
		err := counter.Delete(mockDB)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log("Success Deleted Counter from DB")
}

//func TestCounter_GetLastByMachineCode(t *testing.T) {
//	for _, c := range mockCounter {
//		err := c.GetLastByMachineCode(mockDB, c.MachineCode)
//		if err != nil {
//			t.Error(err.Error())
//		}
//		t.Logf("Get Last Counter from MachineId:", c.Id)
//	}
//}

//func Test_JSONCounterInsert(t *testing.T) {
//	var c model.Counter
//	file, err := os.Open("counter.json")
//	if err != nil {
//		t.Error(err)
//	}
//	d := json.NewDecoder(file)
//	err = d.Decode(&c) // Todo: decoding unsuccessful
//	if err != nil {
//		t.Error(err)
//	}
//	log.Println("Decoded JSON to Struct: ", c)
//	newCounter, err := c.Insert()
//	if err != nil {
//		t.Error(err.Error())
//	}
//	log.Println("Inserted Counter: ", newCounter)
//}

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