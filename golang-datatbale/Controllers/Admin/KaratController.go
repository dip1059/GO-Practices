package Admin

import (
	"github.com/gin-gonic/gin"
	G "gold-store/Globals"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"log"
	"net/http"
	"strconv"
)

var (
	Karat = make(map[uint]Mod.Karat)
)

func AddKaratGet(c *gin.Context) {
	if user, success := M.IsAuthAdminUser(c, G.FStore); success {
		c.HTML(http.StatusOK, "add-karat.html", map[string]interface{}{"AppEnv":G.AppEnv,
			"User":user, "Msg":G.Msg, "Nav":"products", "Title":"Add-Karat"})
		G.Msg.Success = ""
		G.Msg.Fail = ""
	}
}

func AddKaratPost(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var karat Mod.Karat
	err := c.ShouldBind(&karat)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill up all the input fields with valid values."
		c.Redirect(http.StatusFound, "/add-karat")
		return
	}

	if R.AddKarat(karat) {
		G.Msg.Success = "Karat Added Successfully."
		c.Redirect(http.StatusFound, "/all-karat")
		return
	} else {
		G.Msg.Fail = "Some error occurred. Please reload and try again later."
		c.Redirect(http.StatusFound, "/add-karat")
		return
	}
}


func AllKarat(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var karats []Mod.Karat
	karats = R.Karats(karats)
	for _, karat := range karats {
		Karat[karat.ID] = karat
	}
	c.HTML(http.StatusOK, "all-karat.html", map[string]interface{}{"AppEnv":G.AppEnv,
		"User":user, "Nav":"products", "Title":"All-Karat", "Karats":karats, "Msg":G.Msg})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func UpdateStatusKarat(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var karat Mod.Karat
	id, _ := strconv.Atoi(c.Param("id"))
	status, _ := strconv.Atoi(c.Param("status"))
	karat = Karat[uint(id)]
	karat.Status = status
	if R.SaveKarat(karat) {
		G.Msg.Success = "Status Updated Successfully."
		c.Redirect(http.StatusFound, "/all-karat")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Reload And Try Again Later."
		c.Redirect(http.StatusFound, "/all-karat")
	}
}

func EditKarat(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var karat Mod.Karat
	id, _ := strconv.Atoi(c.Param("id"))
	karat = Karat[uint(id)]
	c.HTML(http.StatusOK, "edit-karat.html",map[string]interface{}{"AppEnv":G.AppEnv,
		"User":user,  "Msg":G.Msg, "Nav":"products", "Title":"Edit-Karat", "Karat":karat})
}


func UpdateKarat(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var karat Mod.Karat
	id, _ := strconv.Atoi(c.PostForm("id"))
	karat = Karat[uint(id)]
	err := c.ShouldBind(&karat)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill up all the input fields with valid values."
		c.Redirect(http.StatusFound, "/all-karat")
		return
	}

	if R.SaveKarat(karat) {
		G.Msg.Success = "Updated Successfully"
		c.Redirect(http.StatusFound, "/all-karat")
	} else {
		G.Msg.Fail = "Some Error Occurred, Update Failed. Please Reload And Try Again Later."
		c.Redirect(http.StatusFound, "/all-karat")
	}
}


func DeleteKarat(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var karat Mod.Karat
	id, _ := strconv.Atoi(c.Param("id"))
	karat = Karat[uint(id)]
	if R.DeleteKarat(karat) {
		G.Msg.Success = "Deleted Successfully"
		c.Redirect(http.StatusFound, "/all-karat")
	} else {
		G.Msg.Fail = "Some Error Occurred, Deletion Failed. Please Reload And Try Again Later."
		c.Redirect(http.StatusFound, "/all-karat")
	}
}