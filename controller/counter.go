package controller

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"github.com/mrtomyum/stock/model"
	"github.com/mrtomyum/sys/api"
	"github.com/gin-gonic/gin"
	"strconv"
)

func PostCounter(ctx *gin.Context) {
	log.Println("call ctrl.Counter()")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")

	c := &model.Counter{}
	rs := api.Response{}
	if err := ctx.BindJSON(&c); err != nil {
		log.Println("ctx.BindJSON decode Error: ", err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, rs)
	} else {
		newCounter, err := c.Insert()
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = "CANNOT_INSERT New Counter >>" + err.Error()
			ctx.JSON(http.StatusConflict, rs)
		} else {
			rs.Status = api.SUCCESS
			rs.Data = newCounter
			ctx.JSON(http.StatusOK, rs)
		}
	}
	return
}

func GetAllCounter(ctx *gin.Context) {
	log.Println("call ctrl.Counter.GetAllCounter()")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")

	c := model.Counter{}
	rs := api.Response{}
	counters, err := c.GetAll()
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		ctx.Status(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = counters
		ctx.JSON(http.StatusOK, rs)
	}
	return
}

//====================================
// ขอข้อมูลเคาทเตอร์ เฉพาะรายการตาม id
//====================================
func GetCounter(ctx *gin.Context) {
	log.Println("call ctrl.Counter.GetCounterById()")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Origin", "*")
	//ctx.Header("Authorization", "Baerer" + jwt)

	id := ctx.Param("id")
	c := model.Counter{}
	c.Id, _ = strconv.ParseUint(id, 10, 64)
	rs := api.Response{}
	counters, err := c.Get()
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		ctx.Status(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = counters
		ctx.JSON(http.StatusOK, rs)
	}
	return
}

func PutCounter(ctx *gin.Context) {

}

func DeleteCounter(ctx *gin.Context) {

}

//====================================
// บันทึกเคาทเตอร์ขาย จากหน้าตู้ แบบส่งเป็นชุด
//====================================
//func (e *Env) NewArrayCounter(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Server", "nava Stock")
//	w.Header().Set("Content-Type", "application/json")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//
//	cs := []*model.Counter{}
//	d := json.NewDecoder(r.Body)
//	err := d.Decode(&cs)
//	if err != nil {
//		log.Println("Decode Error: ", err)
//	}
//	rs := api.Response{}
//	newCounters, err := model.NewArrayCounter(e.DB, cs)
//	if err != nil {
//		rs.Status = api.ERROR
//		rs.Message = err.Error()
//		w.WriteHeader(http.StatusConflict)
//	} else {
//		rs.Status = api.SUCCESS
//		rs.Data = newCounters
//		w.WriteHeader(http.StatusOK)
//	}
//	output, _ := json.Marshal(rs)
//	fmt.Fprintf(w, string(output))
//}