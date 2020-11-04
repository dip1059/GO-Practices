package User

import (
	Cfg "gold-store/Config"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	S "gold-store/Services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func Orders(c *gin.Context) {
	var user Mod.User
	var authUser, guestUser bool
	user, guestUser = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	data := c.Param("data")
	if authUser {
		user.CountryCode = G.Country[user.CountryCode]
		var orders []Mod.Order
		orders = R.Orders(orders, "id desc", "user_id=?", user.ID)
		for _, order := range orders {
			order.OrderDetails = R.OrderDetails(order.OrderDetails)
			order.CouponOrder = R.CouponOrder(order.CouponOrder, "order_id=?", order.ID)
			if order.CouponOrder.Coupon.ID != 0 && order.CouponOrder.Coupon.Type == 1 {
				order.IsDiscountCoupon = true
			} else if order.CouponOrder.Coupon.ID != 0 && order.CouponOrder.Coupon.Type == 2 {
				order.IsFree = true
			}
			Order[order.ID] = order
		}
		lenOrder := len(orders)
		if lenOrder > 5 {
			orders = orders[0:5]
		}

		orderCounts := R.CountDiffTypeOrders(user.ID)

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		session, _ := G.Store.Get(c.Request, "cart")
		lenCart := len(session.Values)

		var countries []Country
		var country Country
		for key, value := range G.Country {
			country.Key = key
			country.Value = value
			countries = append(countries, country)
		}

		G.LangVal = H.GetCookie("secret", nil, "lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal

		c.HTML(http.StatusOK, "user-orders.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Account", "Carts": lenCart, "Lang": lang, "LenOrder":lenOrder,
			"Msg": G.Msg, "Adm": G.Adm, "Title": "Account",  "Orders": orders, "Data": data, "Active": "Account",
			"Countries": countries, "Wbsts": wbsts, "Menus": menus,
			"NavActive": "Orders", "OrderCounts": orderCounts,
		})
		G.Msg.Success = ""
		G.Msg.Fail = ""
		return
	} else if guestUser {
		if user.ID > 0 && user.ActiveStatus == 0 {
			G.Msg.Fail = "Please Activate Your Account."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		} else {
			c.Redirect(http.StatusFound, "/login?val=account/orders")
		}
	}
}


func MyWallet(c *gin.Context) {
	var user Mod.User
	var authUser, guestUser bool
	user, guestUser = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	data := c.Param("data")
	if authUser {
		user.GoldTransfers = R.GoldTransfers(user.GoldTransfers, "sender_wallet_id=?", user.Wallet.ID)

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		session, _ := G.Store.Get(c.Request, "cart")
		lenCart := len(session.Values)

		if user.DefaultAddress.ID == 0 {
			user.DefaultAddress.FirstName = user.FirstName
			user.DefaultAddress.LastName = user.LastName
			user.DefaultAddress.Phone = user.Phone
			user.DefaultAddress.Country = user.CountryCode
			user.DefaultAddress.City = user.City
			user.DefaultAddress.Address = user.Address
			user.DefaultAddress.ZipCode = user.ZipCode
			user.DefaultAddress.Country = G.Country[user.DefaultAddress.Country]
		}
		
		var countries []Country
		var country Country
		for key, value := range G.Country {
			country.Key = key
			country.Value = value
			countries = append(countries, country)
		}

		G.LangVal = H.GetCookie("secret", nil, "lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal

		c.HTML(http.StatusOK, "user-wallet.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Account", "Carts": lenCart, "Lang": lang,
			"Msg": G.Msg, "Adm": G.Adm, "Title": "Account", "Data": data, "Active": "Account",
			"Wbsts": wbsts, "Menus": menus, "NavActive": "Wallet","Countries": countries,
		})
		G.Msg.Success = ""
		G.Msg.Fail = ""
		return
	} else if guestUser {
		if user.ID > 0 && user.ActiveStatus == 0 {
			G.Msg.Fail = "Please Activate Your Account."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		} else {
			c.Redirect(http.StatusFound, "/login?val=account/address")
		}
	}
}


func MyAddress(c *gin.Context) {
	var user Mod.User
	var authUser, guestUser bool
	user, guestUser = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	data := c.Param("data")
	if authUser {
		user.ShippingAddresses = R.ShippingAddresses(user.ShippingAddresses, "user_id=?", user.ID)
		for i,_ := range user.ShippingAddresses {
			user.ShippingAddresses[i].Country = G.Country[user.ShippingAddresses[i].Country]
		}

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		session, _ := G.Store.Get(c.Request, "cart")
		lenCart := len(session.Values)

		G.LangVal = H.GetCookie("secret", nil, "lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal

		c.HTML(http.StatusOK, "my-address.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Account", "Carts": lenCart, "Lang": lang,
			"Msg": G.Msg, "Adm": G.Adm, "Title": "Account", "Data": data, "Active": "Account",
			"Wbsts": wbsts, "Menus": menus, "NavActive": "Address",
		})
		G.Msg.Success = ""
		G.Msg.Fail = ""
		return
	} else if guestUser {
		if user.ID > 0 && user.ActiveStatus == 0 {
			G.Msg.Fail = "Please Activate Your Account."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		} else {
			c.Redirect(http.StatusFound, "/login?val=account/address")
		}
	}
}

func AddAddress(c *gin.Context) {
	var user Mod.User
	var authUser, guestUser bool
	user, guestUser = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	data := c.Param("data")
	if authUser {
		

		var countries []Country
		var country Country
		for key, value := range G.Country {
			country.Key = key
			country.Value = value
			countries = append(countries, country)
		}

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		session, _ := G.Store.Get(c.Request, "cart")
		lenCart := len(session.Values)

		G.LangVal = H.GetCookie("secret", nil, "lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal

		c.HTML(http.StatusOK, "add-address.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Account", "Carts": lenCart, "Lang": lang,
			"Msg": G.Msg, "Adm": G.Adm, "Title": "Account", "Data": data, "Active": "Account",
			"Wbsts": wbsts, "Menus": menus, "NavActive": "Address","Countries": countries,
		})
		G.Msg.Success = ""
		G.Msg.Fail = ""
		return
	} else if guestUser {
		if user.ID > 0 && user.ActiveStatus == 0 {
			G.Msg.Fail = "Please Activate Your Account."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		} else {
			c.Redirect(http.StatusFound, "/login?val=account/address")
		}
	}
}


func EditAddress(c *gin.Context) {
	var user Mod.User
	var authUser, guestUser bool
	user, guestUser = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	data := c.Param("data")
	if authUser {
		var countries []Country
		var country Country
		for key, value := range G.Country {
			country.Key = key
			country.Value = value
			countries = append(countries, country)
		}

		id, _ := strconv.Atoi(c.Query("id"))
		var shipAdd Mod.ShippingAddress
		shipAdd.ID = uint(id)
		shipAdd = R.ShippingAddress(shipAdd, "user_id=?", user.ID)
		shipAdd.Country = G.Country[shipAdd.Country]

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		session, _ := G.Store.Get(c.Request, "cart")
		lenCart := len(session.Values)

		G.LangVal = H.GetCookie("secret", nil, "lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal

		c.HTML(http.StatusOK, "edit-address.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Account", "Carts": lenCart, "Lang": lang,
			"Msg": G.Msg, "Adm": G.Adm, "Title": "Account", "Data": data, "Active": "Account",
			"Wbsts": wbsts, "Menus": menus, "NavActive": "Address","Countries": countries,"ShipAdd":shipAdd,
		})
		G.Msg.Success = ""
		G.Msg.Fail = ""
		return
	} else if guestUser {
		if user.ID > 0 && user.ActiveStatus == 0 {
			G.Msg.Fail = "Please Activate Your Account."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		} else {
			c.Redirect(http.StatusFound, "/login?val=account/address")
		}
	}
}


func UpdateShippingAddress(c *gin.Context) {
	var user Mod.User
	user = M.GetAuthUser(c, G.FStore)

	if user.RoleID > 1 {
		id, _ := strconv.Atoi(c.PostForm("id"))
		var shipAdd Mod.ShippingAddress
		shipAdd.ID = uint(id)
		if id > 0 {
			shipAdd = R.ShippingAddress(shipAdd, "user_id=?", user.ID)
		}
		action := c.Query("action")
		shipAdd.UserID = user.ID

		err := c.ShouldBind(&shipAdd)
		if err != nil {
			G.Msg.Fail = "Please fill all the required fields with valid inputs."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		} else {
			log.Println(shipAdd.Status)

			phone, err := strconv.Atoi(c.PostForm("phone"))
			if err != nil {
				log.Println(err.Error())
				G.Msg.Fail = "Please fill all the input fields with valid inputs."
				c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
				return
			}
			shipAdd.Phone = strconv.Itoa(phone)
			isDefault := false
			if action == "add" && shipAdd.Status == 1 {
				isDefault = true
			}

			if R.SaveShippingAddress(shipAdd, isDefault) {
				if action == "add" {
					G.Msg.Success = "Added Successfully."
				} else {
					G.Msg.Success = "Updated Successfully."
				}
				c.Redirect(http.StatusFound, "/account/address")
				return
			} else {
				G.Msg.Fail = "Some Error Occurred, Please Try Again."
				c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
				return
			}

		}
	}
}


func SetAsDefaultShippingAddress(c *gin.Context) {
	var user Mod.User
	user = M.GetAuthUser(c, G.FStore)

	if user.RoleID > 1 {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			log.Println(err.Error())
			return
		}
		var shipAdd Mod.ShippingAddress
		shipAdd.ID = uint(id)
		shipAdd = R.ShippingAddress(shipAdd, "user_id=?", user.ID)
		if shipAdd.Address == "" {
			return
		}
		shipAdd.Status = 1

		if R.SaveShippingAddress(shipAdd, true) {
			G.Msg.Success = "Successfully set as default."
		} else {
			G.Msg.Fail = "Some Error Occurred, Please Try Again."
		}
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
	}
}

func UnsetAsDefaultShippingAddress(c *gin.Context) {
	var user Mod.User
	user = M.GetAuthUser(c, G.FStore)

	if user.RoleID > 1 {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			log.Println(err.Error())
			return
		}
		var shipAdd Mod.ShippingAddress
		shipAdd.ID = uint(id)
		shipAdd = R.ShippingAddress(shipAdd, "user_id=?", user.ID)
		if shipAdd.Address == "" {
			return
		}
		shipAdd.Status = 0

		if R.SaveShippingAddress(shipAdd, false) {
			G.Msg.Success = "Successfully unset as default."
		} else {
			G.Msg.Fail = "Some Error Occurred, Please Try Again."
		}
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
	}
}

func DeleteShippingAddress(c *gin.Context) {
	var user Mod.User
	user = M.GetAuthUser(c, G.FStore)

	if user.RoleID > 1 {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			log.Println(err.Error())
			return
		}
		var shipAdd Mod.ShippingAddress
		shipAdd.ID = uint(id)

		if R.DeleteShippingAddress(shipAdd, "user_id=?", user.ID) {
			G.Msg.Success = "Deleted successfully."
		} else {
			G.Msg.Fail = "Some Error Occurred, Please Try Again."
		}
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
	}
}


func MyWishlist(c *gin.Context) {
	var user Mod.User
	var authUser, guestUser bool
	user, guestUser = M.IsGuest(c, G.FStore)
	user, authUser = M.IsAuthUser(c, G.FStore)
	data := c.Param("data")
	if authUser {

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		session, _ := G.Store.Get(c.Request, "cart")
		lenCart := len(session.Values)

		var wishes []Mod.Wishlist
		wishes = R.Wishlist(wishes, user)

		G.LangVal = H.GetCookie("secret", nil, "lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal

		c.HTML(http.StatusOK, "my-wishlist.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Account", "Carts": lenCart, "Lang": lang,
			"Msg": G.Msg, "Adm": G.Adm, "Title": "Account", "Data": data, "Active": "Account",
			"Wbsts": wbsts, "Menus": menus, "NavActive": "Wishlist", "Wishlist":wishes,
		})
		G.Msg.Success = ""
		G.Msg.Fail = ""
		return
	} else if guestUser {
		if user.ID > 0 && user.ActiveStatus == 0 {
			G.Msg.Fail = "Please Activate Your Account."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		} else {
			c.Redirect(http.StatusFound, "/login?val=account/wishlist")
		}
	}
}


func DownloadOrderPDF(c *gin.Context) {
	user, _ :=M.IsAuthUser(c,G.FStore)

	id, _ := strconv.Atoi(c.Param("id"))
	order := Order[uint(id)]
	//order.OrderDetails = R.OrderDetails(order.OrderDetails)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts)

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	Cfg.LoadAdminSettings()

	pdf := S.NewRequestPdf("")

	err := pdf.ParseTemplate("Views/PdfTemp/order-bank-details-pdf.html", map[string]interface{}{
		"User":user, "Wbsts":wbsts, "Order":order, "Adm":G.Adm, "Lang":lang, "AppEnv":G.AppEnv})
	if err != nil {
		log.Println(err.Error())
		return
	}

	t := time.Now().Unix()
	fileName := strconv.FormatInt(int64(t), 10)
	success := pdf.GeneratePDF("Storage/Temp/"+fileName+".pdf")
	if success {
		log.Println("PDF Generated Successfully.")
		c.Redirect(http.StatusFound, "/assets/Storage/Temp/"+fileName+".pdf")
		/*err := os.Remove("Storage/Temp/"+fileName+".pdf")
		if err != nil {
			log.Println(err.Error())
		}*/
	}
}


func UserOrderDetails(c *gin.Context) {
	user, authUser := M.IsAuthUser(c, G.FStore)
	if !authUser {
		return
	}
	lenWish := R.CountWishlist(user)
	var order Mod.Order
	id, _ := strconv.Atoi(c.Param("id"))
	order = Order[uint(id)]
	//order.OrderDetails = R.OrderDetails(order.OrderDetails)
	session, _ := G.Store.Get(c.Request, "cart")
	lenCart := len(session.Values)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)
	lenDoc := len(order.OrderDocs)

	G.LangVal = H.GetCookie("secret", nil, "lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	c.HTML(http.StatusOK, "user-order-details.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Adm": G.Adm, "Title": "User-Order-Details", "Order": order, "Lang": lang, "Banner": "Account",
		"Msg": G.Msg, "Carts": lenCart, "LenWish": lenWish,  "Wbsts": wbsts, "Menus": menus, "LenDoc": lenDoc})
	G.Msg.Success = ""
	G.Msg.Fail = ""

}

func UpdateProfilePic(c *gin.Context) {
	user, success := M.IsAuthUser(c, G.FStore)
	if !success {
		return
	}

	img, _ := c.FormFile("img")
	if img == nil {
		G.Msg.Fail = "No image file found."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	} 
	file, err := img.Open()
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Invalid file."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	typ, _ := H.ValidateFile(file)
	if typ != 1 {
		G.Msg.Fail = "Invalid file."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}
	
	ext := filepath.Ext(img.Filename)
	imgName := H.RandomString(60) + ext
	dst := "./Storage/Images/" + imgName
	err = c.SaveUploadedFile(img, dst)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Image Upload Failed. Try Again Later."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}
	imgUrl := []byte(dst)
	user.ProfilePic = string(imgUrl[1:])
	user.ProfilePic = "/assets"+user.ProfilePic


	if R.SaveUserChanges(user) {
		G.Msg.Success = "Image uploaded successfully. "
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	} else {
		G.Msg.Fail = "Image Upload Failed. Try Again Later."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

}

func UpdateUser(c *gin.Context) {
	user := M.GetAuthUser(c, G.FStore)
	if user.ID == 0 || user.RegType != 1 {
		return
	}

	oldEmail := user.Email
	data := c.Param("data")
	form := c.PostForm("form")

	currentPassword := c.PostForm("current_password")
	//var success bool

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		if data == "user" {
			G.Msg.Fail = "Current Password Doesn't Match"
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		} else if data == "admin" {
			G.Msg.Fail = "Current Password Doesn't Match."
			c.Redirect(http.StatusFound, "/admin-settings/gs")
			return
		}
	}

	if form == "general" {
		err := c.ShouldBind(&user)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Please fill all the input fields with valid input."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}

		var emptyUser Mod.User
		emptyUser.Email = user.Email
		_, success := R.ReadUserWithEmail(emptyUser, "id <> "+strconv.Itoa(int(user.ID)))
		if success {
			G.Msg.Fail = "Email already exists."
			if data == "user" {
				c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			} else if data == "admin" {
				c.Redirect(http.StatusFound, "/admin-settings/gs")
			}
			return
		}

		emptyUser.Phone = user.Phone
		_, success = R.ReadUserWithPhone(emptyUser, "id <> "+strconv.Itoa(int(user.ID)))
		if success {
			G.Msg.Fail = "Phone number already exists."
			if data == "user" {
				c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			} else if data == "admin" {
				c.Redirect(http.StatusFound, "/admin-settings/gs")
			}
			return
		}

	} else if form == "password" {

		newPassword := c.PostForm("new_password")
		confirmPassword := c.PostForm("confirm_password")
		if newPassword != confirmPassword {
			if data == "user" {
				G.Msg.Fail = "Confirm Password Doesn't Match"
				c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
				return
			} else if data == "admin" {
				G.Msg.Fail = "Confirm Password Doesn't Match."
				c.Redirect(http.StatusFound, "/admin-settings/gs")
				return
			}
		}

		if newPassword != "" {
			cost := bcrypt.DefaultCost
			hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), cost)
			user.Password = string(hash)
		}
	} else {

		if data == "user" {
			G.Msg.Fail = "You Did Something Wrong."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		} else if data == "admin" {
			G.Msg.Fail = "You Did Something Wrong."
			c.Redirect(http.StatusFound, "/admin-settings/gs")
			return
		}
	}

	success := R.SaveUserChanges(user)
	if success {
		if form == "general" {
			session, _ := G.FStore.Get(c.Request, "login_token")
			session.Values["userEmail"] = user.Email
			session.Options.MaxAge = 60 * 60 * 24 * 5
			err := session.Save(c.Request, c.Writer)
			if err != nil {
				log.Println(err.Error())
			}
			if oldEmail != user.Email {
				var eCH Mod.EmailChangeHistory
				eCH.UserID = user.ID
				eCH.OldEmail = oldEmail
				eCH.NewEmail = user.Email
				ok := R.CreateEmailChangeHistory(eCH)
				log.Println("email change history save:", ok)

				user.OldEmail = oldEmail
				ok = S.SendEmailChangeMail(user)
				log.Println("email change mail sending:", ok)

			}

			G.Msg.Success = "Changes Saved Successfully"

		} else if form == "password" {
			G.Msg.Success = "Password changed successfully"
		}

		if data == "user" {
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		} else if data == "admin" {
			c.Redirect(http.StatusFound, "/admin-settings/gs")
			return
		}
	}
}


func UpdateSocialUser(c *gin.Context) {
	user := M.GetAuthUser(c, G.FStore)
	if user.ID == 0 || user.RegType < 2 || user.RegType > 3 {
		return
	}

	var emptyUser Mod.User
	emptyUser.Email = user.Email

	err := c.ShouldBind(&user)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill all the input fields with valid input."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}
	if emptyUser.Email != user.Email {
		G.Msg.Fail = "It is forbidden."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	emptyUser.Phone = user.Phone
	_, success := R.ReadUserWithPhone(emptyUser, "id <> "+strconv.Itoa(int(user.ID)))
	if success {
		G.Msg.Fail = "Phone number already exists."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	user.Email = emptyUser.Email
	user.Password = ""
	success = R.SaveUserChanges(user)
	if success {
		G.Msg.Success = "Changes Saved Successfully"
	} else {
		G.Msg.Fail = "Something went wrong."
	}
	c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
}


func AddOrderDocPost(c *gin.Context) {
	var authUser, guestUser bool
	_, authUser = M.IsAuthUser(c, G.FStore)
	_, guestUser = M.IsGuest(c, G.FStore)
	if authUser {
		orderId, _ := strconv.Atoi(c.PostForm("order_id"))
		var orderDoc Mod.OrderDoc
		orderDoc.OrderID = uint(orderId)
		doc, _ := c.FormFile("doc")
		if doc != nil {

			//file validation
			log.Println(doc.Size)
			if doc.Size > 2097152 {
				G.Msg.Fail = "Document Upload Failed, File Size Exceeds The Limit."
				c.Redirect(http.StatusFound, "/account/orders")
				return
			}
			file, errr := doc.Open()
			if errr !=nil {
				log.Println(errr.Error())
				return
			}
			fileType, success := H.ValidateFile(file)
			if !success {
				G.Msg.Fail = "Document Upload Failed, Not A Valid Type."
				c.Redirect(http.StatusFound, "/account/orders")
				return
			}
			//

			ext := filepath.Ext(doc.Filename)
			docName := H.RandomString(60) + ext
			dst := "./Storage/Images/" + docName
			err := c.SaveUploadedFile(doc, dst)
			if err != nil {
				log.Println(err.Error())
				G.Msg.Fail = "Document Upload Failed. Try Again Later."
				c.Redirect(http.StatusFound, "/account/orders")
				return
			}
			docUrl := []byte(dst)
			orderDoc.DocUrl.String = string(docUrl[1:])
			orderDoc.DocUrl = H.NullStringProcess(orderDoc.DocUrl)
			orderDoc.Type = fileType
		} else {
			G.Msg.Fail = "No File Selected."
			c.Redirect(http.StatusFound, "/account/orders")
			return
		}
		if R.AddOrderDoc(orderDoc) {
			G.Msg.Success = "Document Added Successfully. See The Order's Details."
		} else {
			G.Msg.Fail = "Some Error Occurred. Please Try Again."
		}
		c.Redirect(http.StatusFound, "/account/orders")
	} else if guestUser {
		c.Redirect(http.StatusFound, "/login")
	}
}

func DownloadOrderInvoice(c *gin.Context) {
	if _, success := M.IsAuthUser(c, G.FStore); !success {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return
	}

	var invoice Mod.Invoice
	invoice.OrderID = uint(id)
	invoice = R.AddInvoice(invoice)

	var order Mod.Order
	order = Order[uint(id)]

	//order = R.Order(order)
	if order.OrderStatus != 1 && order.PaymentStatus != 1 {
		return
	}

	//order.User = R.GetUserAddress(order.User)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var lang G.Lang
	lang.LangValue = G.LangVal

	address := template.HTML(G.Adm["Invoice_Company_Address"].Value.String)

	pdf := S.NewRequestPdf("")
	err = pdf.ParseTemplate("Views/PdfTemp/invoice.html", map[string]interface{}{
		"Wbsts":wbsts, "Adm":G.Adm, "Lang":lang, "Address":address,"AppEnv":G.AppEnv, "Order":order})
	if err != nil {
		log.Println(err.Error())
		return
	}

	t := time.Now().Unix()
	fileName := strconv.FormatInt(int64(t), 10)
	success := pdf.GeneratePDF("Storage/Temp/"+fileName+".pdf")
	if success {
		log.Println("PDF Generated Successfully.")
		c.Redirect(http.StatusFound, "/assets/Storage/Temp/"+fileName+".pdf")
	}
}
