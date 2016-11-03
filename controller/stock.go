package controller

import (
	"net/http"
	m "github.com/mrtomyum/stock/model"
	"github.com/mrtomyum/sys/api"
	"github.com/gin-gonic/gin"
)

func GetAllStock(ctx *gin.Context) {
	ctx.Header("Server", "nava Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")
	s := m.StockItem{}
	rs := api.Response{}
	items, err := s.GetAll()
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusNoContent, rs)
	}
	rs.Status = api.SUCCESS
	rs.Data = items
	ctx.JSON(http.StatusOK, rs)
}

func NewStock(ctx *gin.Context) {

}

func FindStockByID(ctx *gin.Context) {

}

func UpdateStock(ctx *gin.Context) {

}

func DelStock(ctx *gin.Context) {

}

func UndelStock(ctx *gin.Context) {

}
