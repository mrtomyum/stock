package ctrl

import (
	"net/http"
	"github.com/mrtomyum/stock/model"
	//log "github.com/Sirupsen/logrus"
	"log"
	"github.com/gin-gonic/gin"
	"strconv"
)

// PostNewMachine สั่งเพิ่มตู้ใหม่ โปรแกรมจะดู MachineType เพื่อสร้าง MachineColumn ให้ตาม Type โดยอัตโนมัติ
func PostNewMachine(c *gin.Context) {
	//log.Info(log.Fields{"func":"ctrl.Machine.PostNewMachine()"})
	log.Println("ctrl.Machine.PostNewMachine()")
	c.Header("Server", "NAVA Stock")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json")

	var m model.Machine
	rs := Response{}
	if err := c.BindJSON(&m); err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
		return
	}
	//log.Info(m)
	newMachine, err := m.New(db) // TODO: ให้ดัก New() ที่ m เป็นค่าว่างด้วย ต้อง Error
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusConflict, rs)
		return
	}
	// ตรวจสอบ MachineType เพื่อ Insert Column ให้อัตโนมัติ
	switch newMachine.Type {
	case model.CAN, model.SEE_THROUGH:
		m.NewColumn(db, m.Selection)
	}
	rs.Status = SUCCESS
	rs.Data = newMachine
	c.JSON(http.StatusOK, rs)
}

func GetAllMachines(ctx *gin.Context) {
	//log.Info(log.Fields{"func":"ctrl.GetAllMachines()"})
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Content-Type", "application/json")

	m := model.Machine{}
	rs := Response{}
	machines, err := m.GetAll(db)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		ctx.Status(http.StatusNoContent)
		return
	}
	rs.Status = SUCCESS
	rs.Data = machines
	ctx.JSON(http.StatusOK, rs)
	return
}

// GetThisMachine คืน JSON ข้อมูล Machine แต่ละตู้พร้อม Sub Columns ทั้งหมด
func GetThisMachine(ctx *gin.Context) {
	//log.Info(log.Fields{"func":"ctrl.GetThisMachine()"})
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")

	id := ctx.Param("id")
	m := model.Machine{}
	m.Id, _ = strconv.ParseUint(id, 10, 64)
	rs := Response{}
	machine, err := m.Get(db)
	// todo: แนบ MachineColumn มาด้วย
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		return
	}
	rs.Status = SUCCESS
	rs.Data = machine
	ctx.JSON(http.StatusOK, rs)
	return
}

func GetMachineColumns(c *gin.Context) {
	//log.Info(log.Fields{"func":"ctrl.Machine.GetMachineColumns()"})
	c.Header("Server", "NAVA Stock")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	id := c.Param("id")
	m := model.Machine{}
	m.Id, _ = strconv.ParseUint(id, 10, 64)
	rs := Response{}
	machineColumns, err := m.GetColumns(db)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		c.Status(http.StatusNoContent)
		return
	}
	rs.Status = SUCCESS
	rs.Data = machineColumns
	c.JSON(http.StatusOK, rs)
	return
}

func GetMachineTemplate(c *gin.Context) {
	var m *model.Machine
	rs := Response{}
	templates, err := m.GetTemplate(db)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusNotFound, rs)
		return
	}
	rs.Status = SUCCESS
	rs.Self = "nava.work/v1/machine/template"
	rs.Data = templates
	c.JSON(http.StatusOK, rs)
}

func PutMachineColumn(c *gin.Context) {
	//log.Info(log.Fields{"func":"ctrl.Machine.PutMachineColumn()"})
	c.Header("Server", "NAVA Stock")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json")

	var col model.MachineColumn
	rs := Response{}
	if err := c.BindJSON(&col); err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
		return
	}
	//TODO: ให้ดัก Update() ที่ mc เป็นค่าว่างด้วย ต้อง Error
	//log.Info(mc)
	switch col.ColumnNo {
	case 0:
		rs.Status = ERROR
		rs.Message = "No data in ColumnNo."
		c.JSON(http.StatusOK, rs)
		return
	default:
		err := col.Update(db)
		if err != nil {
			rs.Status = ERROR
			rs.Message = err.Error()
			c.JSON(http.StatusOK, rs)
			return
		}
		rs.Status = SUCCESS
		rs.Data = col
	}
	c.JSON(http.StatusOK, rs)
	return
}

// สั่งเพิ่ม Column ที่ขาดให้กับ Machine จนครบตาม Machine.Selection
func PostMachineColumnInit(ctx *gin.Context) {
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Origin", "*")
	id := ctx.Param("id")
	m := new(model.Machine)
	m.Id, _ = strconv.ParseUint(id, 10, 64)
	rs := Response{}
	m, err := m.Get(db)
	count, err := m.InitMachineColumn(db)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusConflict, rs)
		return
	}
	rs.Status = SUCCESS
	sCount := strconv.Itoa(count)
	rs.Message = "New Column Count = " + sCount
	rs.Data = m
	ctx.JSON(http.StatusOK, rs)
}
