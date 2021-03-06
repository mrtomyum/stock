package ctrl

import (
	m "github.com/mrtomyum/stock/model"
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func GetAllItem(ctx *gin.Context) {
	log.Println("call GET AllItem")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Host", "nava.work:8001")
	ctx.Header("Access-Control-Allow-Origin", "*")

	i := m.Item{}
	items, err := i.GetAll(db)
	rs := Response{}
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, rs)
		return
	}
	rs.Status = SUCCESS
	rs.Data = items
	ctx.JSON(http.StatusFound, rs)
}

func PostNewItem(ctx *gin.Context) {
	log.Println("call POST NewItem")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")

	i := new(m.Item)
	rs := new(Response)
	if ctx.BindJSON(&i) != nil {
		//log.Println("NewItem: Error decode.Decode(&i) >>", err)
		rs.Status = ERROR
		rs.Message = "Cannot Bind JSON requested."
		ctx.JSON(http.StatusBadRequest, rs)
		return
	}
	newItem, err := i.Insert(db)
	log.Println("i= ", i)
	if err != nil {
		rs.Status = ERROR
		rs.Message = "CANNOT_UPDATE >>" + err.Error()
		ctx.JSON(http.StatusConflict, rs)
		return
	}
	rs.Status = SUCCESS
	rs.Data = newItem
	ctx.JSON(http.StatusOK, rs)
	return
}

func GetItem(ctx *gin.Context) {
	log.Println("call FindItem")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")
	var i m.Item
	id := ctx.Param("id")
	i.Id, _ = strconv.ParseUint(id, 10, 64)
	rs := Response{}
	iv, err := i.GetItemView(db)
	log.Println("return from GetItemView()")
	if err != nil {
		log.Println(err)
		rs.Status = ERROR
		rs.Message = err.Error()
		return
	}
	rs.Status = SUCCESS
	rs.Data = iv
	ctx.JSON(200, rs)
	return
}

func UpdateItem(ctx *gin.Context) {
	log.Println("call UpdateItem")
	ctx.Header("Server", "NAVA Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")
	var i m.Item
	rs := Response{}
	if ctx.BindJSON(&i) != nil {
		ctx.JSON(http.StatusBadRequest, i)
		return
	}
	updatedItem, err := i.Update(db)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusNotModified, rs)
		return
	}
	rs.Status = SUCCESS
	rs.Data = updatedItem
	ctx.JSON(http.StatusOK, rs)
	return
}

func DelItem(ctx *gin.Context) {
	log.Println("call NewItem")
}

func UndelItem(ctx *gin.Context) {
	log.Println("call NewItem")
}
