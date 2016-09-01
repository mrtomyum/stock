package controller

import (
	"net/http"
	"github.com/mrtomyum/nava-sys/api"
	m "github.com/mrtomyum/nava-stock/model"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

//func (e *Env) AllMachine(w http.ResponseWriter, r *http.Request) {
func (e *Env) AllMachine(c *gin.Context) {
	log.Info(log.Fields{"func":"controller.AllMachine()"})
	c.Header("Server", "NAVA Stock")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	m := m.Machine{}
	rs := api.Response{}
	machines, err := m.All(e.DB)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.Status(http.StatusNoContent)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = machines
		c.JSON(http.StatusOK, rs)
	}
	return
}

func (e *Env) NewMachine(c *gin.Context) {
	log.Info(log.Fields{"func":"controller.NewMachine()"})
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