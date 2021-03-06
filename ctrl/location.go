package ctrl

import (
	m "github.com/mrtomyum/stock/model"
	"net/http"
	"log"
	"strconv"
	"github.com/gin-gonic/gin"
)

func CreateLocationTree(locations []*m.Location) *m.Location {
	tree := new(m.Location)
	for _, l := range locations {
		tree.AddTree(l)
	}
	return tree
}

func GetAllLocationTree(ctx *gin.Context) {
	log.Println("call ShowLocationTree()")
	ctx.Header("Access-Control-Allow-Origin", "*")
	loc := new(m.Location)
	rs := new(Response)
	locations, err := loc.All(db)
	if err != nil {
		log.Fatal("Error LocationsTreeAll()", err)
		rs.Status = ERROR
		rs.Message = "Location not found or Error."
		ctx.JSON(http.StatusBadRequest, rs)
		return
	}
	tree := CreateLocationTree(locations)
	rs.Status = SUCCESS
	rs.Data = tree.Child
	ctx.JSON(http.StatusOK, rs)
}

func GetLocationTreeByID(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	l := new(m.Location)
	id := ctx.Param("id")
	l.Id, _ = strconv.ParseUint(id, 10, 64)

	locations, err := l.Get(db)
	rs := new(Response)
	if err != nil {
		log.Fatal("Error LocationsTreeByID()", err)
		rs.Status = ERROR
		rs.Message = "Location not found or Error."
		ctx.JSON(http.StatusNoContent, rs)
	}
	tree := CreateLocationTree(locations)
	rs.Status = SUCCESS
	rs.Data = tree.Child
	ctx.JSON(http.StatusOK, rs)
}

func PostNewLocation(ctx *gin.Context) {
	log.Println("call ShowLocationTree()")
	ctx.Header("Access-Control-Allow-Origin", "*")
	l := m.Location{}
	if ctx.BindJSON(&l) != nil {
		ctx.JSON(http.StatusBadRequest, l)
		log.Println("Error in Decoded request body.")
	}
	log.Println("Success decode JSON -> :", l, " Result location decoded -> ", l.Code)

	newLoc, err := l.Insert(db)
	rs := new(Response)
	if err != nil {
		rs.Status = ERROR
		rs.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, rs)
		return
	}
	rs.Status = SUCCESS
	rs.Data = newLoc
	ctx.JSON(http.StatusOK, rs)
}
