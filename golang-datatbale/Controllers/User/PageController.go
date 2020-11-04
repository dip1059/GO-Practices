package User

import (
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func AboutPage(c *gin.Context) {
	var user Mod.User
	var authUser, guest bool
	user, guest = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	lenWish := R.CountWishlist(user)
	session, _ := G.Store.Get(c.Request, "cart")
	lenCart := len(session.Values)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	var page Mod.Page
	page = R.Page(page, "status=1 and url='about-us'")

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	c.HTML(http.StatusOK, "about-us.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": page.Title, "Carts": lenCart,"Lang":lang,
		"Msg": G.Msg, "Adm": G.Adm, "Active":page.Title, "Title": page.Title, "LenWish": lenWish, "Wbsts": wbsts, "Menus":menus, "Page":page})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func ContactPage(c *gin.Context) {
	var user Mod.User
	var authUser, guest bool
	user, guest = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	lenWish := R.CountWishlist(user)
	session, _ := G.Store.Get(c.Request, "cart")
	lenCart := len(session.Values)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	var page Mod.Page
	page = R.Page(page, "status=1 and url='contact-us'")

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	c.HTML(http.StatusOK, "contact-us.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": page.Title, "Carts": lenCart,"Lang":lang,
		"Msg": G.Msg, "Adm": G.Adm, "Active":page.Title, "Title": page.Title, "LenWish": lenWish, "Wbsts": wbsts, "Menus":menus, "Page":page})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func ContactMessage(c *gin.Context) {
	var cm Mod.ContactMessage
	var success bool
	err := c.ShouldBind(&cm)
	if err != nil {
		log.Println(err.Error())
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	_, err = strconv.Atoi(cm.Phone)
	if err != nil {
		log.Println(err.Error())
	}

	if cm.Email == "" || !re.MatchString(cm.Email) || cm.Phone == "" || err != nil || cm.FullName == "" || cm.Message == "" {
		G.Msg.Fail = "All Fields Are Required And Inputs Must Be Valid."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	cm, success = R.AddContactMessage(cm)
	if success {
		G.Msg.Success = "Message Sent Successfully."
	} else {
		G.Msg.Fail = "Some Error Occurred. Please Try Again."
	}
	c.Redirect(http.StatusFound, "/contact-us")
}


func SubsNewsletter(c *gin.Context) {
	var nl Mod.Newsletter
	var success bool
	err := c.ShouldBind(&nl)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Invalid Email Address."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	//re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	/*nl.Email = c.PostForm("email")
	if nl.Email == "" || !re.MatchString(nl.Email) {
		G.Msg.Fail = "Invalid Email Address."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}*/

	nl2 := R.Newsletter(nl, "email=?", nl.Email)
	if nl2.ID != 0 {
		G.Msg.Fail = "Already Subscribed."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	nl, success = R.AddNewsletter(nl)
	if success {
		G.Msg.Success = "You Subscribed Successfully."
	} else {
		G.Msg.Fail = "Some Error Occurred,Please Try Again."
	}
	c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])

}


func AlertasPage(c *gin.Context) {
	var user Mod.User
	var authUser, guest bool
	user, guest = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	lenWish := R.CountWishlist(user)
	session, _ := G.Store.Get(c.Request, "cart")
	lenCart := len(session.Values)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	var page Mod.Page
	page = R.Page(page, "status=1 and url='alerts'")

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	c.HTML(http.StatusOK, "alertas.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": page.Title, "Carts": lenCart,"Lang":lang,
		"Msg": G.Msg, "Adm": G.Adm, "Active":page.Title, "Title": page.Title, "LenWish": lenWish, "Wbsts": wbsts, "Menus":menus, "Page":page})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func PrivacidadePage(c *gin.Context) {
	var user Mod.User
	var authUser, guest bool
	user, guest = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	lenWish := R.CountWishlist(user)
	session, _ := G.Store.Get(c.Request, "cart")
	lenCart := len(session.Values)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	var page Mod.Page
	page = R.Page(page, "status=1 and url='privacypolicy'")

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	c.HTML(http.StatusOK, "politicas-de-privacidade.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": page.Title, "Carts": lenCart, "Lang":lang,
		"Msg": G.Msg, "Adm": G.Adm, "Active":page.Title, "Title": page.Title, "LenWish": lenWish, "Wbsts": wbsts, "Menus":menus, "Page":page})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func CookiesPage(c *gin.Context) {
	var user Mod.User
	var authUser, guest bool
	user, guest = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	lenWish := R.CountWishlist(user)
	session, _ := G.Store.Get(c.Request, "cart")
	lenCart := len(session.Values)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	var page Mod.Page
	page = R.Page(page, "status=1 and url='cookiespolicies'")

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	c.HTML(http.StatusOK, "politicas-de-cookies.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": page.Title, "Carts": lenCart,"Lang":lang,
		"Msg": G.Msg, "Adm": G.Adm, "Active":page.Title, "Title": page.Title, "LenWish": lenWish, "Wbsts": wbsts, "Menus":menus, "Page":page})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func TermosPage(c *gin.Context) {
	var user Mod.User
	var authUser, guest bool
	user, guest = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	lenWish := R.CountWishlist(user)
	session, _ := G.Store.Get(c.Request, "cart")
	lenCart := len(session.Values)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	var page Mod.Page
	page = R.Page(page, "status=1 and url='termsofuse'")

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	c.HTML(http.StatusOK, "termos-de-uso.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": page.Title, "Carts": lenCart,"Lang":lang,
		"Msg": G.Msg, "Adm": G.Adm, "Active":page.Title, "Title": page.Title, "LenWish": lenWish, "Wbsts": wbsts, "Menus":menus, "Page":page})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func DynamicPage(c *gin.Context) {
	var user Mod.User
	var authUser, guest bool
	user, guest = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	lenWish := R.CountWishlist(user)
	session, _ := G.Store.Get(c.Request, "cart")
	lenCart := len(session.Values)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	url := c.Param("url")
	var page Mod.Page
	page = R.Page(page, "status=1 and url=?", url)

	if url == "about-us" {
		c.Redirect(http.StatusFound, "/about-us")
		return
	} else if url == "contact-us" {
		c.Redirect(http.StatusFound, "/contact-us")
		return
	}

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	c.HTML(http.StatusOK, "dynamic-page.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": page.Title, "Carts": lenCart,"Lang":lang,
		"Msg": G.Msg, "Adm": G.Adm, "Active":page.Title, "Title": page.Title, "LenWish": lenWish, "Wbsts": wbsts, "Menus":menus, "Page":page})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}