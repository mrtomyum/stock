package ctrl

import (
	"net/http"
	"github.com/mrtomyum/sys/api"
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

	var m model.Machine
	rs := api.Response{}
	if err := c.BindJSON(&m); err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
		return
	}
	//log.Info(m)
	newMachine, err := m.New() // TODO: ให้ดัก New() ที่ m เป็นค่าว่างด้วย ต้อง Error
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusConflict, rs)
		return
	}
	// ตรวจสอบ MachineType เพื่อ Insert Column ให้อัตโนมัติ
	switch newMachine.Type {
	case model.CAN, model.SEE_THROUGH:
		m.NewColumn(m.Selection)
	}
	rs.Status = api.SUCCESS
	rs.Data = newMachine
	c.JSON(http.StatusOK, rs)
}

func GetAllMachines(ctx *gin.Context) {
	//log.Info(log.Fields{"func":"ctrl.GetAllMachines()"})
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")

	m := model.Machine{}
	rs := api.Response{}
	machines, err := m.GetAll()
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		ctx.Status(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = machines
		ctx.JSON(http.StatusOK, rs)
	}
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
	rs := api.Response{}
	machine, err := m.Get()
	// todo: แนบ MachineColumn มาด้วย
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Data = machine
	}
	ctx.JSON(http.StatusOK, rs)
	return
}

func GetMachineColumns(c *gin.Context) {
	//log.Info(log.Fields{"func":"ctrl.Machine.GetMachineColumns()"})
	c.Header("Server", "NAVA Stock")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	var m model.Machine
	rs := api.Response{}

	machineColumns, err := m.GetColumns()
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.Status(http.StatusNoContent)
	}
	rs.Status = api.SUCCESS
	rs.Data = machineColumns
	c.JSON(http.StatusOK, rs)
	return
}

func GetMachineTemplate(c *gin.Context) {
	var m *model.Machine
	rs := api.Response{}
	templates, err := m.GetTemplate()
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
	}
	rs.Status = api.SUCCESS
	rs.Self = "api.nava.work/v1/machine/template"
	rs.Data = templates
	c.JSON(http.StatusOK, rs)
}

func PutMachineColumn(c *gin.Context) {
	//log.Info(log.Fields{"func":"ctrl.Machine.PutMachineColumn()"})
	c.Header("Server", "NAVA Stock")
	c.Header("Access-Control-Allow-Origin", "*")

	var col model.MachineColumn
	rs := api.Response{}
	if err := c.BindJSON(&col); err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
	} else {
		//TODO: ให้ดัก Update() ที่ mc เป็นค่าว่างด้วย ต้อง Error
		//log.Info(mc)
		switch col.ColumnNo {
		case 0:
			rs.Status = api.ERROR
			rs.Message = "No data in ColumnNo."
		default:
			err := col.Update()
			if err != nil {
				rs.Status = api.ERROR
				rs.Message = err.Error()
			} else {
				rs.Status = api.SUCCESS
				rs.Data = col
			}
			c.JSON(http.StatusOK, rs)
		}
	}
	return
}
