package User

import (
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	S "gold-store/Services"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"math/big"
	"net/http"
	"strconv"
)

func AddCart(c *gin.Context) {
	_, authUser := M.IsGuest(c, G.FStore)
	_, guest := M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	id := c.Param("id")
	quantity := c.Param("quantity")

	data := c.Param("data")
	success := S.AddCart(c, id, quantity)

	/*lang = H.GetCookie("secret", nil,"lang", c)
	if lang == "" {
		lang = "en"
		H.SetCookie("secret", nil, lang, "lang", 60*60*24*365, c)
	}*/
	fromAjax := c.Query("ajax")
	if fromAjax == "" {
		fromAjax = "false"
	}

	process := c.Query("process")

	if success {
		if process == "update" {
			G.Msg.Success = "Product quantity updated successfully."
		} else {
			G.Msg.Success = "Product Added To Cart Successfully."// template.HTML(gotrans.Tr(lang, "Product Added To Cart Successfully."))
		}

		if data == "cart" {
			c.Redirect(http.StatusFound, "/show-cart?ajax="+fromAjax)
			return
		} else if wishId, err:=strconv.Atoi(data); err == nil && wishId > 0 {
			var wish Mod.Wishlist
			wish.ID = uint(wishId)
			R.RemoveFromWishlist(wish)
			c.Redirect(http.StatusFound, "/account/wish")
			return
		}
		c.Redirect(http.StatusFound, "/")
		return

	} else {
		G.Msg.Fail = "Some Error Occurred. Please Try Again."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0]+"?ajax="+fromAjax)
		//c.Redirect(http.StatusFound, "/product-details/"+strconv.Itoa(id))
	}
}


func AddAllToCart(c *gin.Context) {
	user, authUser := M.IsGuest(c, G.FStore)
	_, guest := M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	user.Wishlist = R.Wishlist(user.Wishlist, user)

	var success bool
	for _, wish := range user.Wishlist {
		success = S.AddCart(c, strconv.Itoa(int(wish.ProductID)), "1")
	}

	if success {
		G.Msg.Success = "All Products Added To Cart Successfully."// template.HTML(gotrans.Tr(lang, "Product Added To Cart Successfully."))
		c.Redirect(http.StatusFound, "/show-cart?ajax=false")
		return

	} else {
		G.Msg.Fail = "Some Error Occurred. Please Try Again."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
	}
}


func AddCartCustom(c *gin.Context) {
	_, guest := M.IsGuest(c, G.FStore)
	user, authUser := M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	id := c.PostForm("id")
	var product Mod.Product
	uid, _ := strconv.Atoi(id)
	product.ID = uint(uid)
	product = R.Product(product, user.ID)

	var quantity string

	fromAjax := c.Query("ajax")
	if fromAjax == "" {
		fromAjax = "false"
	}

	amountType := c.PostForm("amount_type")
	if amountType == "2" {
		price, err := strconv.ParseFloat(c.PostForm("amount"), 64)
		if err != nil {
			G.Msg.Fail = "You Did Something Wrong. Invalid Amount."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0]+"?ajax="+fromAjax)
			return
		}
		price = math.Floor(price * 100.0) / 100.0
		log.Println(price)
		log.Println((product.GrmAmount / product.Price) * price)
		grmAmount := math.Floor((product.GrmAmount / product.Price) * price)
		quantity = fmt.Sprintf("%.1f", grmAmount)
		log.Println(quantity)
	} else {
		//quantity = c.PostForm("amount")
		quan, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
		quantity = fmt.Sprintf("%.1f",quan)
	}

	bigQuan, success := new(big.Float).SetString(quantity)
	bigMax, _ := new(big.Float).SetString(product.Max)
	bigMin, _ := new(big.Float).SetString(product.Min)

	if !success {
		G.Msg.Fail = "You Did Something Wrong. Invalid Amount."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0]+"?ajax="+fromAjax)
		return
	} else {
		if bigQuan.Cmp(bigMax) == 1 || bigQuan.Cmp(bigMin) == -1 {
			G.Msg.Fail = "Invalid amount, doesn't meet minimum gram amount requirement."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0]+"?ajax="+fromAjax)
			return
		}
	}
	data := c.PostForm("data")
	success = S.AddCart(c, id, quantity)

	/*lang := H.GetCookie("secret", nil,"lang", c)
	if lang == "" {
		lang = "en"
		H.SetCookie("secret", nil, lang, "lang", 60*60*24*365, c)
	}*/
	process := c.Query("process")

	if success {
		if process == "update" {
			G.Msg.Success = "Product quantity updated successfully."
		} else {
			G.Msg.Success = "Product Added To Cart Successfully."// template.HTML(gotrans.Tr(lang, "Product Added To Cart Successfully."))
		}
		if data == "cart" {
			c.Redirect(http.StatusFound, "/show-cart?ajax="+fromAjax)
			return
		} else if wishId, err:=strconv.Atoi(data); err == nil && wishId > 0 {
			var wish Mod.Wishlist
			wish.ID = uint(wishId)
			R.RemoveFromWishlist(wish)
			c.Redirect(http.StatusFound, "/account/wish")
			return
		}
		c.Redirect(http.StatusFound, "/")
		return

	} else {
		G.Msg.Fail = "Some Error Occurred. Please Try Again."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0]+"?ajax="+fromAjax)
		//c.Redirect(http.StatusFound, "/product-details/"+strconv.Itoa(id))
	}
}

func ShowCart(c *gin.Context) {
	var user Mod.User
	var authUser, guest bool
	user, authUser = M.IsGuest(c, G.FStore)
	user, guest = M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	lenWish := R.CountWishlist(user)
	var wishlist []Mod.Wishlist
	wishlist = R.Wishlist(wishlist, user)
	//session, _ := G.Store.Get(c.Request, "cart")
	//lenCart := len(session.Values)

	finalCart := S.ProcessCart(c)
	lenCart := len(finalCart.Carts)

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	c.HTML(http.StatusOK, "show-cart.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": "Cart", "Carts": lenCart,"Lang":lang,
		"Msg": G.Msg, "Adm": G.Adm, "Title": "Cart", "FinalCart": finalCart, "Wishlist": wishlist,
		"LenWish": lenWish, "Wbsts":wbsts, "Menus":menus, "Active": "Cart"})


	fromAjax := c.Query("ajax")
	if fromAjax != "true" {
		G.Msg.Success = ""
		G.Msg.Fail = ""
	}
}

func DeleteFromCart(c *gin.Context) {
	_, authUser := M.IsGuest(c, G.FStore)
	_, guest := M.IsAuthUser(c, G.FStore)
	if !authUser && !guest {
		return
	}
	key := c.Param("key")
	success := S.DeleteFromCart(c, key)
	if success {
		G.Msg.Success = "Product Removed From Cart Successfully."
		c.Redirect(http.StatusFound, "/show-cart")
	} else {
		G.Msg.Fail = "Some Error Occurred. Please Try Again."
		c.Redirect(http.StatusFound, "/show-cart")
	}
}
