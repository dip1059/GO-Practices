package User

import (
	"github.com/gin-gonic/gin"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	S "gold-store/Services"
	"net/http"
	"strconv"
)

var (
	Order   = make(map[uint]Mod.Order)
	//lang string
)

type Country struct {
	Key   string
	Value string
}

func Home(c *gin.Context) {
	var user Mod.User
	//var authUser, guest bool
	user, _ = M.IsGuest(c, G.FStore)
	/*user, authUser = M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}*/
	if user.RoleID == 1 {
		return
	}

	G.LangVal = H.GetCookie("secret", nil, "lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	var products []Mod.Product

	lenWish := R.CountWishlist(user)
	session, _ := G.Store.Get(c.Request, "cart")
	lenCart := len(session.Values)

	fixedProducts := R.ProductsWithOthers(products, user.ID,"status = ? and type = ?", 1, 1)
	customProducts := R.ProductsWithOthers(products, user.ID,"status = ? and type = ?", 1, 2)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	var homeContents []Mod.HomeContent
	homeContents = R.HomeContents(homeContents, "status=?", 1)

	c.HTML(http.StatusOK, "home.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": "All Products", "Carts": lenCart, "Lang":lang,
		"Msg": G.Msg, "Adm": G.Adm, "Title": "Home", "FixedProducts": fixedProducts, "LenWish": lenWish,
		"Wbsts": wbsts, "Menus": menus, "Active":"Home","CustomProducts": customProducts, "HomeContents": homeContents,
	})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func ProductDetails(c *gin.Context) {
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

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	var product Mod.Product
	id, _ := strconv.Atoi(c.Param("id"))
	product.ID = uint(id)
	product = R.Product(product, user.ID)
	if product.Type == 2 {
		return
	}

	var products []Mod.Product
	products = R.ProductsWithOthers(products, user.ID, "id <> ? and type = ? and status = ?", product.ID, 1, 1)
	if len(products) > 4 {
		products = products[0:4]
	}

	c.HTML(http.StatusOK, "product-details.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": "Product Details", "Carts": lenCart,"Lang":lang,
		"Msg": G.Msg, "Adm": G.Adm, "Title": "Product Details", "Product": product, "LenWish": lenWish,
		"Wbsts": wbsts, "Menus": menus, "Active": "Product-Details", "Products":products})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func ChangeLang(c *gin.Context) {
	lang := c.PostForm("lang")
	H.SetCookie("secret",nil, lang,"lang",60*60*24*365, c)
	c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
}


func BuyNow(c *gin.Context) {
	user, guest := M.IsGuest(c, G.FStore)
	user, authUser := M.IsAuthUser(c, G.FStore)
	id := c.Param("id")
	if authUser {
		proID, _ := strconv.Atoi(id)
		finalCart := S.BuyNowProcess(uint(proID), c)

		lenWish := R.CountWishlist(user)

		session, _ := G.Store.Get(c.Request, "cart")
		lenCart := len(session.Values)

		DataAmount := int64(finalCart.GrandTotal * 100.0)
		stripePK := G.Adm["Stripe_Publishable_Key"].Value.String
		//brainTkzKey := G.Adm["Brain_Tree_Tokenization_Key"].Value.String
		bankRefCode := "REF"+H.RandomString(11)

		var payMethods []Mod.PayMethod
		payMethods = R.AllMethod(payMethods, "status=1")

		var banks []Mod.Bank
		banks = R.AllBank(banks, "status=1")

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		G.LangVal = H.GetCookie("secret", nil,"lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal


		redsysOrdId := H.RandomNumber(10)

		c.HTML(http.StatusOK, "checkout.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Checkout", "Carts": lenCart,"Lang":lang,"Banks": banks,
			"Msg": G.Msg, "Adm": G.Adm, "Title": "Checkout", "DataAmount": DataAmount, "FinalCart": finalCart,"FreeBids":false,
			/*"SecretKey":sessionId, "RedsysUrl":G.Redsys.Url, */"StripePK":stripePK, "ProductID":id,
			"LenWish": lenWish, "PayMethods": payMethods,  "UniqueCode":bankRefCode, "RedsysOrdId":redsysOrdId,
			"Wbsts": wbsts, "Menus":menus,/*"WalletApp":G.WalletApp, */})
		G.Msg.Success = ""
		G.Msg.Fail = ""
	} else if guest {
		if user.ID > 0 && user.ActiveStatus == 0 {
			G.Msg.Fail = "Please Activate Your Account."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		}else {
			c.Redirect(http.StatusFound, "/login")
		}
	}
}


func GetCustomProductAjax(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	success := true
	message := ""
	var pro Mod.Product

	if err != nil {
		success = false
		message = "Invalid product"
	} else {
		pro.ID = uint(id)
		pro = R.Product(pro, 0)
	}

	c.JSON(http.StatusOK,gin.H {
		"success": success,
		"message": message,
		"product": pro,
	})

}