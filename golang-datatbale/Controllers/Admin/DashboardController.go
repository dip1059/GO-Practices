package Admin

import (
	G "gold-store/Globals"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

type DashData struct {
	TotalUsers, TotalProducts, TotalOrders int
	TotalSales float64
	NewUsers []Mod.User
	ThisYearSale []float64
	CountrySale []R.CountrySale
}

func Dashboard(c *gin.Context) {

	if user, success := M.IsAuthAdminUser(c, G.FStore); success {
		var dash DashData
		var total = make([]float64, 12)

		dash.TotalUsers = R.TotalUsers()
		dash.TotalOrders = R.TotalOrders()
		dash.TotalProducts = R.TotalProducts()
		dash.TotalSales = math.Ceil(R.TotalSales()*100)/100

		data := R.ThisYearSales()
		lenData := len(data)
		for i:=0; i<lenData; i++ {
			total[data[i].Month-1] = data[i].Total
		}
		dash.ThisYearSale = total
		dash.CountrySale = R.CountrySales()
		for i, _ := range dash.CountrySale {
			dash.CountrySale[i].Name = G.Country[dash.CountrySale[i].Name]
		}

		dash.NewUsers = R.NewUsers()

		c.HTML(http.StatusOK, "dashboard.html", map[string]interface{}{
			"AppEnv":G.AppEnv, "User":user,  "Nav":"dashboard", "Title":"Dashboard", "Msg":G.Msg, "DashData":dash})
		G.Msg.Success = ""
		G.Msg.Fail = ""
	}
	return
}
