package Admin

import (
	G "gold-store/Globals"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	Coupon = make(map[uint]Mod.Coupon)
)

func CouponSettings(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var coupons []Mod.Coupon
	coupons = R.Coupons(coupons)
	for _, coupon := range coupons{
		Coupon[coupon.ID] = coupon
	}

	c.HTML(http.StatusOK, "coupon-settings.html", map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user, "Nav":"coupon", "Title":"Coupon-Settings", "Msg":G.Msg,
		"Coupons":coupons, })
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func AddCouponGet(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}

	c.HTML(http.StatusOK, "add-coupon.html", map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user, "Nav":"coupon", "Title":"Add-Coupon", "Msg":G.Msg, })
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func AddCouponPost(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var coupon Mod.Coupon
	err := c.ShouldBind(&coupon)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Some Error Occured. Please Try Again Later."
		c.Redirect(http.StatusFound, "/add-coupon")
		return
	}
	coupon.Code = strings.TrimSpace(coupon.Code)

	if R.CouponExists(coupon.Code) {
		G.Msg.Fail = "Coupon Code Already Used. Set An Unique Code."
		c.Redirect(http.StatusFound, "/add-coupon")
		return
	}

	//coupon.StartDate, _ = time.Parse("2006-01-02 15:04", c.PostForm("start_date"))
	coupon.StartDate, _ = time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"))
	coupon.EndDate, _ = time.Parse("2006-01-02 15:04", c.PostForm("end_date"))
	
	if R.AddCoupon(coupon) {
		G.Msg.Success = "Coupon Added Successfully."
		c.Redirect(http.StatusFound, "/coupon-settings")
		return
	} else {
		G.Msg.Fail = "Some Error Occured. Please Try Again Later."
		c.Redirect(http.StatusFound, "/add-coupon")
		return
	}
}


func MakeCouponInactive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var coupon Mod.Coupon
	id, _ := strconv.Atoi(c.Param("id"))
	coupon = Coupon[uint(id)]
	coupon.Status = 0
	if R.SaveCoupon(coupon) {
		G.Msg.Success = "Status Updated Successfully"
		c.Redirect(http.StatusFound, "/coupon-settings")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/coupon-settings")
	}
}


func MakeCouponActive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var coupon Mod.Coupon
	id, _ := strconv.Atoi(c.Param("id"))
	coupon = Coupon[uint(id)]
	coupon.Status = 1
	if R.SaveCoupon(coupon) {
		G.Msg.Success = "Status Updated Successfully."
		c.Redirect(http.StatusFound, "/coupon-settings")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/coupon-settings")
	}
}


func EditCoupon(c *gin.Context) {
	var user Mod.User
	var success bool
	user, success = M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var coupon Mod.Coupon
	id, _ := strconv.Atoi(c.Param("id"))
	coupon = Coupon[uint(id)]

	c.HTML(http.StatusOK, "edit-coupon.html",map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user,  "Nav":"coupon", "Title":"Edit-Coupon",
		"Coupon":coupon, "Msg":G.Msg})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func UpdateCoupon(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var coupon Mod.Coupon
	id, _ := strconv.Atoi(c.PostForm("ID"))
	coupon = Coupon[uint(id)]
	err := c.ShouldBind(&coupon)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Some Error Occured. Please Try Again Later."
		c.Redirect(http.StatusFound, "/coupon-settings")
		return
	}


	//coupon.StartDate, _ = time.Parse("2006-01-02 15:04", c.PostForm("start_date"))
	coupon.StartDate, _ = time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"))
	coupon.EndDate, _ = time.Parse("2006-01-02 15:04", c.PostForm("end_date"))

	if R.SaveCoupon(coupon) {
		G.Msg.Success = "Updated Successfully"
		c.Redirect(http.StatusFound, "/coupon-settings")
	} else {
		G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/coupon-settings")
	}
}


func DeleteCoupon(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var coupon Mod.Coupon
	id, _ := strconv.Atoi(c.Param("id"))
	coupon = Coupon[uint(id)]
	if R.DeleteCoupon(coupon) {
		G.Msg.Success = "Deleted Successfully"
		c.Redirect(http.StatusFound, "/coupon-settings")
	} else {
		G.Msg.Fail = "Some Error Occurred, Deletion Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/coupon-settings")
	}
}
