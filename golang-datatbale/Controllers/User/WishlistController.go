package User

import (
	G "gold-store/Globals"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddWish(c *gin.Context) {
	var user Mod.User
	var authUser, guestUser bool
	user, authUser = M.IsAuthUser(c, G.FStore)
	user, guestUser = M.IsGuest(c, G.FStore)
	if authUser {
		proId, _ := strconv.Atoi(c.Param("proId"))
		var wishlist Mod.Wishlist
		wishlist.UserID = user.ID
		wishlist.ProductID = uint(proId)
		success := R.AddWish(wishlist)
		if success {
			G.Msg.Success = "Added To Wishlist Successfully."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		} else {
			G.Msg.Fail = "Some Error Occurred. Please Try Again."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		}
	} else if guestUser {
		c.Redirect(http.StatusFound, "/login")
	}
}

func RemoveFromWishlist(c *gin.Context) {
	var guestUser bool
	user, authUser := M.IsAuthUser(c, G.FStore)
	_, guestUser = M.IsGuest(c, G.FStore)
	if authUser {
		var success bool
		id, _ := strconv.Atoi(c.Param("id"))
		var wishlist Mod.Wishlist
		if id != 0 {
			wishlist.ID = uint(id)
			success = R.RemoveFromWishlist(wishlist)
		} else {
			proId, _ := strconv.Atoi(c.Param("proId"))
			success = R.RemoveFromWishlist(wishlist, "user_id = ? and product_id = ?", user.ID, proId)
		}
		if success {
			G.Msg.Success = "Removed Successfully."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		} else {
			G.Msg.Fail = "Some Error Occurred. Please Try Again."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		}
	} else if guestUser {
		c.Redirect(http.StatusFound, "/login")
	}
}
