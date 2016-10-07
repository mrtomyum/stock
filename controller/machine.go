package controller

import (
	"net/http"
	"github.com/mrtomyum/sys/api"
	m "github.com/mrtomyum/nava-stock/model"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"strconv"
)

//func (e *Env) AllMachine(w http.ResponseWriter, r *http.Request) {
func (e *Env) GetAllMachines(ctx *gin.Context) {
	log.Info(log.Fields{"func":"controller.GetAllMachines()"})
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Origin", "*")

	m := m.Machine{}
	rs := api.Response{}
	machines, err := m.All(e.DB)
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

func (e *Env) GetThisMachine(ctx *gin.Context) {
	log.Info(log.Fields{"func":"controller.GetThisMachine()"})
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Origin", "*")

	id := ctx.Param("id")
	m := m.Machine{}
	m.ID, _ = strconv.ParseUint(id, 10, 64)
	rs := api.Response{}
	machine, err := m.Get(e.DB)
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

func (e *Env) PostNewMachine(c *gin.Context) {
	log.Info(log.Fields{"func":"controller.Machine.PostNewMachine()"})
	c.Header("Server", "NAVA Stock")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	var m m.Machine
	rs := api.Response{}
	if err := c.BindJSON(&m); err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
	} else {
		log.Info(m)
		newMachine, err := m.New(e.DB) // TODO: ให้ดัก New() ที่ m เป็นค่าว่างด้วย ต้อง Error
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = err.Error()
		} else {
			rs.Status = api.SUCCESS
			rs.Data = newMachine
		}
		c.JSON(http.StatusOK, rs)
	}
}

func (e *Env) GetMachineColumns(c *gin.Context) {
	log.Info(log.Fields{"func":"controller.Machine.GetMachineColumns()"})
	c.Header("Server", "NAVA Stock")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	var m m.Machine
	rs := api.Response{}
	machineColumns, err := m.Columns(e.DB)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.Status(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = machineColumns
		c.JSON(http.StatusOK, rs)
	}
	return
}

func (e *Env) PutMachineColumn(c *gin.Context) {
	log.Info(log.Fields{"func":"controller.Machine.PutMachineColumn()"})
	c.Header("Server", "NAVA Stock")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	var mc m.MachineColumn
	rs := api.Response{}
	if err := c.BindJSON(&mc); err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
	} else {
		//TODO: ให้ดัก Update() ที่ mc เป็นค่าว่างด้วย ต้อง Error
		log.Info(mc)
		switch mc.ColumnNo {
		case 0:
			rs.Status = api.ERROR
			rs.Message = "No data in ColumnNo."
		default:
			updatedColumn, err := mc.Update(e.DB)
			if err != nil {
				rs.Status = api.ERROR
				rs.Message = err.Error()
			} else {
				rs.Status = api.SUCCESS
				rs.Data = updatedColumn
			}
			c.JSON(http.StatusOK, rs)
		}
	}
	return
}