package Admin

import (
	G "gold-store/Globals"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)


var (
	User = make(map[uint]Mod.User)
)

type UserData struct {
	UserModel Mod.User
	CreatedAt string
}

func Users(c *gin.Context) {
	authUser, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var users []Mod.User
	var usersData []UserData
	var userData UserData
	users = R.UsersWithOthers(users, "role_id = 2")
	for _, user := range users {
		user.CountryCode = G.Country[user.CountryCode]
		User[user.ID] = user
		userData.UserModel = user
		userData.CreatedAt = strconv.Itoa(user.CreatedAt.Day())+"-"+user.CreatedAt.Month().String()+"-"+strconv.Itoa(user.CreatedAt.Year())
		usersData = append(usersData, userData)
	}

	c.HTML(http.StatusOK, "users.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": authUser,  "Nav":"users", "Title": "User-Management", "Msg": G.Msg, "Users": usersData})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func MakeUserSuspend(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var user Mod.User
	user.ID = uint(id)
	if R.UpdateUser(user, map[string]interface{}{"active_status": 2}) {
		G.Msg.Success = "Status Updated Successfully"
		c.Redirect(http.StatusFound, "/users")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/users")
	}

}


func MakeUserActive(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var user Mod.User
	user.ID = uint(id)
	if R.UpdateUser(user, map[string]interface{}{"active_status": 1}) {
		G.Msg.Success = "Status Updated Successfully"
		c.Redirect(http.StatusFound, "/users")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/users")
	}

}


func UserOrders(c *gin.Context) {
	authUser, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var orders []Mod.Order
	orders = R.Orders(orders, "created_at desc","user_id=?",id)
	for _, order := range orders {
		G.Order[order.ID] = order
	}
	c.HTML(http.StatusOK, "orders.html",map[string]interface{}{
		"AppEnv":G.AppEnv, "User":authUser,  "Nav":"users", "Title":"User-Orders", "Msg":G.Msg, "Orders":orders, "UserID":id})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func DeleteUser(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var user Mod.User
	user.ID = uint(id)
	if R.DeleteUser(user) {
		G.Msg.Success = "Successfully Deleted"
		c.Redirect(http.StatusFound, "/users")
	} else {
		G.Msg.Fail = "Some Error Occurred, Delete Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/users")
	}
}
