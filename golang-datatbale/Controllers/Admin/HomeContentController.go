package Admin

import (
	"github.com/gin-gonic/gin"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var(
	HomeContent Mod.HomeContent
)

func EditHomeContent(c *gin.Context) {
	var user Mod.User
	var success bool
	user, success = M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var homeContent Mod.HomeContent
	id, _ := strconv.Atoi(c.Param("id"))
	homeContent.ID = uint(id)
	homeContent = R.HomeContent(homeContent)
	HomeContent = homeContent

	c.HTML(http.StatusOK, "edit-home-content.html",map[string]interface{}{"AppEnv":G.AppEnv, "User":user,
		"Nav":"settings", "Title":"Edit-Home-Content", "HomeContent":homeContent, "Msg":G.Msg})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func UpdateHomeContent(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var homeContent Mod.HomeContent
	homeContent = HomeContent

	id, _ := strconv.Atoi(c.PostForm("id"))
	homeContent.ID = uint(id)

	homeContent.TextContent = template.HTML(c.PostForm("text_content"))

	img, _ := c.FormFile("img")
	if img != nil {
		os.Remove("." + homeContent.Image)

		ext := filepath.Ext(img.Filename)
		imgName := H.RandomString(60) + ext
		dst := "./Storage/Images/" + imgName
		err := c.SaveUploadedFile(img, dst)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Image Upload Failed. Try Again Later."
			c.Redirect(http.StatusFound, "/website-settings/home-content")
			return
		}
		imgUrl := []byte(dst)
		homeContent.Image = string(imgUrl[1:])
	}

	if R.UpdateHomeContent(homeContent) {
		G.Msg.Success = "Updated Successfully"
		c.Redirect(http.StatusFound, "/website-settings/home-content")
	} else {
		G.Msg.Fail = "Some Error Occurred, Update Failed. Please Reload And Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/home-content")
	}
}
