package ctrl

import (
	"net/http"
	//log "github.com/Sirupsen/logrus"
	"log"
	"github.com/mrtomyum/stock/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

// PostNewCounter สร้างรายการจด Counter ใหม่ของแต่ละตู้
// โดยร้องขอมาเป็น JSON counter + counter_sub
// ผลลัพธ์นี้จะต้องสร้างใบนำส่งเงินรอไว้ (cash/matching)
// หากส่งซ้ำ ระบบจะสนใจเฉพาะรายการล่าสุด และแก้ไขใบนำส่งเงินที่สร้างไว้
func PostNewCounter(ctx *gin.Context) {
	log.Println("call ctrl.Counter()")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")

	c := &model.Counter{}
	rs := Response{}
	if err := ctx.BindJSON(&c); err != nil {
		log.Println("ctx.BindJSON decode Error: ", err)
		rs.Status = ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, rs)
		return
	}
	newCounter, err := c.Insert(db)
	if err != nil {
		rs.Status = ERROR
		rs.Message = "CANNOT_INSERT New Counter >>" + err.Error()
		ctx.JSON(http.StatusConflict, rs)
	}
	rs.Status = SUCCESS
	rs.Data = newCounter
	ctx.JSON(http.StatusOK, rs)
	return
}

func GetAllCounter(ctx *gin.Context) {
	log.Println("call ctrl.Counter.GetAllCounter()")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Origin", "*")

	c := model.Counter{}
	rs := Response{}
	counters, err := c.GetAll(db)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		ctx.Status(http.StatusNoContent)
	} else {
		rs.Status = SUCCESS
		rs.Data = counters
		ctx.JSON(http.StatusOK, rs)
	}
	return
}

//====================================
// ขอข้อมูลเคาทเตอร์ เฉพาะรายการตาม id ของ Counter
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
	counters, err := c.Get(db)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		ctx.Status(http.StatusNoContent)
		return
	}
	rs.Status = SUCCESS
	rs.Data = counters
	ctx.JSON(http.StatusOK, rs)
	return
}

//====================================
// GetLastCounterByMachineCode ขอข้อมูลเคาท์เตอร์ล่าสุดของแต่ละตู้ตาม MachineCode
//====================================
func GetLastCounterByMachineCode(ctx *gin.Context) {
	log.Println("call ctrl.Counter.GetCounterByMachineCode()")
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Origin", "*")
	machineCode := ctx.Param("code")
	c := model.Counter{}
	err := c.GetLastByMachineCode(db, machineCode)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusNotFound, c)
		return
	}
	rs.Status = SUCCESS
	rs.Data = c
	ctx.JSON(http.StatusOK, rs)
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
//	rs := Response{}
//	newCounters, err := model.NewArrayCounter(e.DB, cs)
//	if err != nil {
//		rs.Status = ERROR
//		rs.Message = err.Error()
//		w.WriteHeader(http.StatusConflict)
//	} else {
//		rs.Status = SUCCESS
//		rs.Data = newCounters
//		w.WriteHeader(http.StatusOK)
//	}
//	output, _ := json.Marshal(rs)
//	fmt.Fprintf(w, string(output))
//}
