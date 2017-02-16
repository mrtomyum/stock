package ctrl

import (
	"net/http"
	m "github.com/mrtomyum/stock/model"
	"github.com/gin-gonic/gin"
)

func GetAllStock(ctx *gin.Context) {
	ctx.Header("Server", "nava Stock")
	ctx.Header("Access-Control-Allow-Origin", "*")
	s := m.StockItem{}
	rs := Response{}
	items, err := s.GetAll()
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusNoContent, rs)
	}
	rs.Status = SUCCESS
	rs.Data = items
	ctx.JSON(http.StatusOK, rs)
}

func NewDoc(ctx *gin.Context) {
	var d m.Doc
	rs := Response{}
	if ctx.Bind(&d) != nil {
		rs.Message = "Cannot bind JSON requested."
		rs.Status = ERROR
		ctx.JSON(http.StatusBadRequest, rs)
		return
	}
	rs.Status = SUCCESS
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
