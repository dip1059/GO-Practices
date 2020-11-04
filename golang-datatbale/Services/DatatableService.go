package Services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//PaginateDataStruct fields for the AJAX response to paginate
type PaginateDataStruct struct {
	Draw            interface{} `json:"draw"`
	RecordsTotal    interface{} `json:"recordsTotal"`
	RecordsFiltered interface{} `json:"recordsFiltered"`
	Data            interface{} `json:"data"`
}

func ProcessDatatableData(c *gin.Context) map[string]interface{} {
	var Params = make(map[string]interface{})

	draw := c.PostForm("draw")
	offset := c.PostForm("start")
	limit := c.PostForm("length")
	search := c.PostForm("search[value]")
	orderCol, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	sort := strconv.Itoa(orderCol+1) + " " + c.PostForm("order[0][dir]")
	arg := "%" + search + "%"

	Params["draw"] = draw
	Params["offset"] = offset
	Params["limit"] = limit
	Params["sort"] = sort
	Params["search"] = search
	Params["arg"] = arg

	return Params
}

func SendDatatableData(c *gin.Context, Params map[string]interface{}) {
	var paging PaginateDataStruct
	paging.Data = Params["data"]
	paging.RecordsTotal = Params["count"]
	paging.RecordsFiltered = Params["count"]
	paging.Draw = Params["draw"]

	c.JSON(http.StatusOK, paging)
}
