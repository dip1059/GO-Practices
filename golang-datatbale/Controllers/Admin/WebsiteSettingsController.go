package Admin

import (
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

var(
	Wbst = make(map[uint]Mod.WebsiteSetting)
	Page = make(map[uint]Mod.Page)
	Menu = make(map[uint]Mod.Menu)
)

func WebsiteSettings(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	data := c.Param("data")
	//primary
	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts)
	for _, wbst := range wbsts{
		Wbst[wbst.ID] = wbst
	}
	//

	//page
	var pages []Mod.Page
	pages = R.Pages(pages)
	for _, page := range pages{
		Page[page.ID] = page
	}
	//

	//menu
	var menus []Mod.Menu
	menus = R.Menus(menus)
	for _, menu := range menus{
		Menu[menu.ID] = menu
	}
	//

	//home contents
	var homeContents []Mod.HomeContent
	homeContents = R.HomeContents(homeContents)
	//

	c.HTML(http.StatusOK, "website-settings.html", map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user, "Nav":"settings", "Title":"Website-Settings", "Msg":G.Msg,
		"Wbsts":wbsts, "Pages":pages,  "Menus":menus,"HomeContents":homeContents, "Data":data})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func UpdateWebsiteSettings(c * gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	
	var wbst Mod.WebsiteSetting
	hdLogo, _ := c.FormFile("hd_logo")
	wbst = Wbst[1]
	if hdLogo != nil {
		//os.Remove("."+string(wbst.Content))

		ext := filepath.Ext(hdLogo.Filename)
		imgName := H.RandomString(60) + ext
		dst := "./Storage/Images/" + imgName
		err := c.SaveUploadedFile(hdLogo, dst)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Image Upload Failed. Try Again Later."
			c.Redirect(http.StatusFound, "/website-settings/prm")
			return
		}
		imgUrl := []byte(dst)
		wbst.Content = template.HTML(imgUrl[1:])
		if !R.SaveWebsiteSetting(wbst) {
			G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
			c.Redirect(http.StatusFound, "/website-settings/prm")
			return
		}
		
	}

	ftLogo, _ := c.FormFile("ft_logo")
	wbst = Wbst[2]
	if ftLogo != nil {
		//os.Remove("."+string(wbst.Content))

		ext := filepath.Ext(ftLogo.Filename)
		imgName := H.RandomString(60) + ext
		dst := "./Storage/Images/" + imgName
		err := c.SaveUploadedFile(ftLogo, dst)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Image Upload Failed. Try Again Later."
			c.Redirect(http.StatusFound, "/website-settings/prm")
			return
		}
		imgUrl := []byte(dst)
		wbst.Content = template.HTML(imgUrl[1:])
		if !R.SaveWebsiteSetting(wbst) {
			G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
			c.Redirect(http.StatusFound, "/website-settings/prm")
			return
		}

	}

	ftText := c.PostForm("footer_text")
	wbst = Wbst[3]
	if ftText != "" {
		wbst.Content = template.HTML(ftText)
		if !R.SaveWebsiteSetting(wbst) {
			G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
			c.Redirect(http.StatusFound, "/website-settings/prm")
			return
		}
	}

	cprText := c.PostForm("cpr_text")
	wbst = Wbst[4]
	if cprText != "" {
		wbst.Content = template.HTML(cprText)
		if !R.SaveWebsiteSetting(wbst) {
			G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
			c.Redirect(http.StatusFound, "/website-settings/prm")
			return
		}
	}
	G.Msg.Success = "Changes Saved Successfully."
	c.Redirect(http.StatusFound, "/website-settings/prm")
}


func AddPage(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var page Mod.Page
	page.Title = c.PostForm("title")
	page.TextContent = template.HTML(c.PostForm("text_content"))
	page.Status, _ = strconv.Atoi(c.PostForm("status"))
	page.Url = H.MakeUrl(page.Title)
	if R.AddPage(page) {
		G.Msg.Success = "Page Added Successfully."
		c.Redirect(http.StatusFound, "/website-settings/page")
		return
	} else {
		G.Msg.Fail = "Some Error Occured. Please Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/page")
		return
	}
}


func MakePageInactive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var page Mod.Page
	id, _ := strconv.Atoi(c.Param("id"))
	page = Page[uint(id)]
	page.Status = 0
	if R.SavePage(page) {
		G.Msg.Success = "Status Updated Successfully"
		c.Redirect(http.StatusFound, "/website-settings/page")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/page")
	}
}


func MakePageActive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var page Mod.Page
	id, _ := strconv.Atoi(c.Param("id"))
	page = Page[uint(id)]
	page.Status = 1
	if R.SavePage(page) {
		G.Msg.Success = "Status Updated Successfully."
		c.Redirect(http.StatusFound, "/website-settings/page")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/page")
	}
}


func EditPage(c *gin.Context) {
	var user Mod.User
	var success bool
	user, success = M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var page Mod.Page
	id, _ := strconv.Atoi(c.Param("id"))
	page = Page[uint(id)]

	c.HTML(http.StatusOK, "edit-page.html",map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user,  "Nav":"settings", "Title":"Edit-Page", "Page":page, "Msg":G.Msg})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func UpdatePage(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var page Mod.Page
	id, _ := strconv.Atoi(c.PostForm("id"))
	page = Page[uint(id)]
	page.Title = c.PostForm("title")
	page.TextContent = template.HTML(c.PostForm("text_content"))

	img, _ := c.FormFile("img")
	if img != nil {
		//os.Remove("."+page.ImgUrl.String)
		ext := filepath.Ext(img.Filename)
		imgName := H.RandomString(60) + ext
		dst := "./Storage/Images/" + imgName
		err := c.SaveUploadedFile(img, dst)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Image Upload Failed. Try Again Later."
			c.Redirect(http.StatusFound, "/add-page")
			return
		}
		imgUrl := []byte(dst)
		page.ImgUrl.String = string(imgUrl[1:])
		page.ImgUrl = H.NullStringProcess(page.ImgUrl)
	}

	if R.SavePage(page) {
		G.Msg.Success = "Updated Successfully"
		c.Redirect(http.StatusFound, "/website-settings/page")
	} else {
		G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/page")
	}
}


func DeletePage(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var page Mod.Page
	id, _ := strconv.Atoi(c.Param("id"))
	page = Page[uint(id)]
	if R.DeletePage(page) {
		G.Msg.Success = "Deleted Successfully"
		c.Redirect(http.StatusFound, "/website-settings/page")
	} else {
		G.Msg.Fail = "Some Error Occurred, Deletion Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/page")
	}
}


func AddMenu(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var menu Mod.Menu
	menu.Title = c.PostForm("title")
	parentID, _ := strconv.Atoi(c.PostForm("parent_id"))
	menu.ParentID = uint(parentID)
	menu.Status, _ = strconv.Atoi(c.PostForm("status"))
	menu.Level = 2
	menu.Type, _ = strconv.Atoi(c.PostForm("type"))
	if menu.Type == 1 {
		menu.PageUrl = c.PostForm("url")
	} else if menu.Type == 2 {
		menu.PageUrl = c.PostForm("page")
	}
	menu.Position = Menu[menu.ParentID].Position
	if R.AddMenu(menu) {
		G.Msg.Success = "Menu Added Successfully."
		c.Redirect(http.StatusFound, "/website-settings/menu")
		return
	} else {
		G.Msg.Fail = "Some Error Occured. Please Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/menu")
		return
	}
}


func MakeMenuInactive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var menu Mod.Menu
	id, _ := strconv.Atoi(c.Param("id"))
	menu = Menu[uint(id)]
	menu.Status = 0
	if R.SaveMenu(menu) {
		G.Msg.Success = "Status Updated Successfully"
		c.Redirect(http.StatusFound, "/website-settings/menu")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/menu")
	}
}


func MakeMenuActive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var menu Mod.Menu
	id, _ := strconv.Atoi(c.Param("id"))
	menu = Menu[uint(id)]
	menu.Status = 1
	if R.SaveMenu(menu) {
		G.Msg.Success = "Status Updated Successfully."
		c.Redirect(http.StatusFound, "/website-settings/menu")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/menu")
	}
}


func EditMenu(c *gin.Context) {
	var user Mod.User
	var success bool
	user, success = M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var menu Mod.Menu
	id, _ := strconv.Atoi(c.Param("id"))
	menu = Menu[uint(id)]

	c.HTML(http.StatusOK, "edit-menu.html",map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user,  "Nav":"settings", "Title":"Edit-Menu", "Menu":menu, "Msg":G.Msg, "Pages":Page})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func UpdateMenu(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var menu Mod.Menu
	id, _ := strconv.Atoi(c.PostForm("id"))
	menu = Menu[uint(id)]
	if c.PostForm("title") != "" {
		menu.Title = c.PostForm("title")
	}
	parentID, _ := strconv.Atoi(c.PostForm("parent_id"))
	if parentID != 0 {
		menu.ParentID = uint(parentID)
	}
	if c.PostForm("page_url") != "" {
		menu.PageUrl = c.PostForm("page_url")
	}

	if R.SaveMenu(menu) {
		G.Msg.Success = "Updated Successfully"
		c.Redirect(http.StatusFound, "/website-settings/menu")
	} else {
		G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/menu")
	}
}


func DeleteMenu(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var menu Mod.Menu
	id, _ := strconv.Atoi(c.Param("id"))
	menu = Menu[uint(id)]
	if R.DeleteMenu(menu) {
		G.Msg.Success = "Deleted Successfully"
		c.Redirect(http.StatusFound, "/website-settings/menu")
	} else {
		G.Msg.Fail = "Some Error Occurred, Deletion Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/website-settings/menu")
	}
}