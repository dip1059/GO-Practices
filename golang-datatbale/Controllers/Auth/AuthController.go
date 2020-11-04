package Auth

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	S "gold-store/Services"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func LoginGet(c *gin.Context) {

	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	if user, success := M.IsGuest(c, G.FStore); success {
		//c.Redirect(http.StatusFound, WalletAppUrl+"/login?"+G.AppEnv.Name)
		//convRate,_, currCode, currSym := S.CurrencyProcessing(c)

		finalCart := S.ProcessCart(c)

		lenCart := len(finalCart.Carts)

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		G.LangVal = H.GetCookie("secret", nil, "lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal

		lenWish := R.CountWishlist(user)

		url := c.Query("val")

		c.HTML(http.StatusOK, "login.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Login", "Carts": lenCart, "Lang": lang,
			"Msg": G.Msg, "Title": "Login", "Wbsts": wbsts, "Url": url, "Menus": menus, "LenWish": lenWish,
			"Adm": G.Adm,
		})
		G.Msg.Success = ""
		G.Msg.Fail = ""
	}
}

func LoginPost(c *gin.Context) {
	user, success := M.IsGuest(c, G.FStore)
	if user.ID > 0 {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	url := c.Query("val")
	log.Println(url)

	user.Email = c.PostForm("email")
	password := c.PostForm("password")
	rememberMe, _ := strconv.Atoi(c.PostForm("remember_me"))

	user, success = R.Login(user)
	if success {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			G.Msg.Fail = "Email or password does not match."
			c.Redirect(http.StatusFound, "/login")
		} else {
			if user.ActiveStatus == 1 {
				H.SetCookie("secret", nil, "", "email", -1, c)

				session, _ := G.FStore.Get(c.Request, "login_token")
				session.Values["userEmail"] = user.Email
				session.Options.MaxAge = 60 * 60 * 24 * 5
				err := session.Save(c.Request, c.Writer)
				if err != nil {
					log.Println(err.Error())
				}

				if rememberMe == 1 {
					session.Options.MaxAge = 60 * 60 * 24 * 365
					err := session.Save(c.Request, c.Writer)
					if err != nil {
						log.Println(err.Error())
					}
				}
				if user.RoleID == 1 {
					c.Redirect(http.StatusFound, "/dashboard")
				} else if user.RoleID == 2 {
					if url != "" {
						log.Println(url)
						c.Redirect(http.StatusFound, "/"+url)
					} else {
						log.Println(url)
						c.Redirect(http.StatusFound, "/")
					}
				} /* else if user.RoleID == 3 {
					if url != "" {
						c.Redirect(http.StatusFound, "/"+url)
					} else {
						c.Redirect(http.StatusFound, "/vendor/dashboard")
					}
				}*/
			} else if user.ActiveStatus == 2 {
				if G.Msg.Fail == "" {
					G.Msg.Fail = "You Are Suspended. Contact With The Authority Quickly."
				}
				c.Redirect(http.StatusFound, "/login")
			} else {
				H.SetCookie("secret", nil, user.Email, "email", 60*60*5, c)
				var link template.HTML
				link = "<a href='" + template.HTML(G.AppEnv.Url) + "/resend-email-verification'>Click Here To Send Verification Email</a>"
				if G.Msg.Fail == "" {
					G.Msg.Fail = "Please Activate Your Account, " + link + "."
				}

				c.Redirect(http.StatusFound, "/login")
			}
		}

	} else {
		if G.Msg.Fail == "" {
			G.Msg.Fail = "Email does not exist."
		}
		c.Redirect(http.StatusFound, "/login")
	}
}


func SocialLogin(c *gin.Context) {
	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	driver := c.Query("driver")
	log.Println(driver)
	clientID := ""
	clientSecret := ""
	if driver == "facebook" {
		clientID = G.SocialEnv.FacebookClientID
		clientSecret = G.SocialEnv.FacebookClientSecret
		H.SetCookie("secret", nil, "2", "driver", 60*60*24, c)
	} else if driver == "google" {
		clientID = G.SocialEnv.GoogleClientID
		clientSecret = G.SocialEnv.GoogleClientSecret
		H.SetCookie("secret", nil, "3", "driver", 60*60*24, c)
	} else {
		G.Msg.Fail = "Driver not found"
		log.Println(driver)
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	authUrl, success := S.SocialLogin(driver, clientID, clientSecret, "/social-login-callback")
	if success {
		c.Redirect(http.StatusFound, authUrl)
	} else {
		G.Msg.Fail = "Something went wrong."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}
}


func SocialLoginCallback(c *gin.Context) {
	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	state := c.Query("state")
	code := c.Query("code")
	driver := H.GetCookie("secret", nil, "driver", c)
	userModel, _, success := S.SocialLoginCallback(state, code)

	if success {
		var user Mod.User
		user.Email = userModel.Email
		if user.Email == "" {
			G.Msg.Fail = "Login failed, no email found."
			c.Redirect(http.StatusFound, "/login")
			return
		}
		user, success := R.ReadUserWithEmail(user)
		log.Println(user.Email, user.RegType)

		if success && driver != strconv.Itoa(user.RegType) {
			G.Msg.Fail = "Login failed, email already used."
			c.Redirect(http.StatusFound, "/login")
			return
		} else if success && driver == strconv.Itoa(user.RegType) && user.SocialAuthID == userModel.ID {
			H.SetCookie("secret", nil, "", "driver", -1, c)

			//user.SocialAuthID = userModel.ID
			//_, success = R.AddSocialUser(user)
			//if !success {
			//	G.Msg.Fail = "Something went wrong."
			//	c.Redirect(http.StatusFound, "/login")
			//	return
			//}

			session, _ := G.FStore.Get(c.Request, "login_token")
			session.Values["userEmail"] = user.Email
			session.Options.MaxAge = 60 * 60 * 24 * 5
			err := session.Save(c.Request, c.Writer)
			if err != nil {
				log.Println(err.Error())
				G.Msg.Fail = "Something went wrong."
				c.Redirect(http.StatusFound, "/login")
				return
			}
			c.Redirect(http.StatusFound,"/")
			return
		} else {
			if driver == "" {
				G.Msg.Fail = "Login failed, it took too much time."
				c.Redirect(http.StatusFound, "/login")
				return
			} else if driver == "3" || driver == "2" {
				user.RegType, _ = strconv.Atoi(driver)
				user.SocialAuthID = userModel.ID
				user.FirstName = userModel.FirstName
				user.LastName = userModel.LastName
				user.ProfilePic = userModel.Avatar
				user.Email = userModel.Email

				H.SetCookie("secret", nil, user.Email, "abc", 60*60*24, c)
				S.RedirectToProfileComplete(c, user)
				return

			} else {
				return
			}
		}

	} else {
		H.SetCookie("secret", nil, "", "driver", -1, c)
		H.SetCookie("secret", nil, "", "abc", -1, c)
		G.Msg.Fail = "Something went wrong."
		c.Redirect(http.StatusFound, "/login")
		return
	}
}


func SocialRegister(c *gin.Context) {
	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}
	var user Mod.User
	var success bool

	driver := H.GetCookie("secret", nil, "driver", c)
	email := H.GetCookie("secret", nil, "abc", c)

	if email == "" || (driver != "3" && driver != "2") {
		G.Msg.Fail = "Login failed, it took too much time."
		c.Redirect(http.StatusFound, "/login")
		return
	}

	user.RegType, _ = strconv.Atoi(driver)
	user.SocialAuthID = c.PostForm("social")
	user.ProfilePic = c.PostForm("avatar")

	if user.SocialAuthID == "" || user.ProfilePic == "" {
		return
	}

	err := c.ShouldBind(&user)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill all the required fields with valid input."
		S.RedirectToProfileComplete(c, user)
		return
	}

	if user.FirstName == "" || user.LastName == "" || user.Phone == "" {
		G.Msg.Fail = "Please fill all the required fields with valid input."
		S.RedirectToProfileComplete(c, user)
		return
	}

	phone, err := strconv.Atoi(c.PostForm("phone"))
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill all the input fields with valid input."
		S.RedirectToProfileComplete(c, user)
		return
	}
	user.Phone = strconv.Itoa(phone)

	_, success = R.ReadUserWithEmail(user)
	if success {
		G.Msg.Fail = "Email already exists."
		c.Redirect(http.StatusFound, "/login")
		return
	}

	_, success = R.ReadUserWithPhone(user)
	if success {
		G.Msg.Fail = "Phone number already exists."
		S.RedirectToProfileComplete(c, user)
		return
	}

	user.ActiveStatus = 1
	user.RoleID = 2
	user.Email = email
	user.Password = ""
	user, success = R.AddSocialUser(user)

	if !success {
		H.SetCookie("secret", nil, "", "driver", -1, c)
		H.SetCookie("secret", nil, "", "abc", -1, c)

		G.Msg.Fail = "Something went wrong."
		c.Redirect(http.StatusFound, "/login")
		return
	}
	H.SetCookie("secret", nil, "", "driver", -1, c)
	H.SetCookie("secret", nil, "", "abc", -1, c)

	session, _ := G.FStore.Get(c.Request, "login_token")
	session.Values["userEmail"] = user.Email
	session.Options.MaxAge = 60 * 60 * 24 * 5
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Something went wrong."
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.Redirect(http.StatusFound,"/")
	return
}


func RegisterGet(c *gin.Context) {

	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	if user, success := M.IsGuest(c, G.FStore); success {
		//c.Redirect(http.StatusFound, WalletAppUrl+"/login?tab=registration&"+G.AppEnv.Name)
		//convRate,_, currCode, currSym := S.CurrencyProcessing(c)

		finalCart := S.ProcessCart(c)
		lenCart := len(finalCart.Carts)

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		G.LangVal = H.GetCookie("secret", nil, "lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal

		lenWish := R.CountWishlist(user)

		c.HTML(http.StatusOK, "register.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Sign Up", "Carts": lenCart, "Lang": lang,
			"Msg": G.Msg, "Title": "Sign Up", "Wbsts": wbsts, "Menus": menus, "LenWish": lenWish,
			"Adm": G.Adm})
		G.Msg.Success = ""
		G.Msg.Fail = ""

	}

}

func RegisterPost(c *gin.Context) {
	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var success bool
	var user Mod.User
	err := c.ShouldBind(&user)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill all the required fields with valid input."
		c.Redirect(http.StatusFound, "/register")
		return
	}

	if user.FirstName == "" || user.LastName == "" || user.Phone == "" {
		G.Msg.Fail = "Please fill all the required fields with valid input."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	phone, err := strconv.Atoi(c.PostForm("phone"))
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill all the input fields with valid input."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}
	user.Phone = strconv.Itoa(phone)

	_, success = R.ReadUserWithEmail(user)
	if success {
		G.Msg.Fail = "Email already exists."
		c.Redirect(http.StatusFound, "/register")
		return
	}

	_, success = R.ReadUserWithPhone(user)
	if success {
		G.Msg.Fail = "Phone number already exists."
		c.Redirect(http.StatusFound, "/register")
		return
	}
	password := c.PostForm("password")
	confirmPass := c.PostForm("confirm_password")
	if password != confirmPass {
		G.Msg.Fail = "Confirm Password Doesn't Match."
		c.Redirect(http.StatusFound, "/register")
		return
	}

	cost := bcrypt.DefaultCost
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), cost)
	user.Password = string(hash)
	user.ActiveStatus = 0

	country := c.PostForm("country")
	user.CountryCode = country
	user.EmailVerification = H.RandomString(60)
	user.EmailVerifyCode = H.RandomNumber(10)

	user.RoleID = 2
	user, success = R.Register(user)

	if success {
		success = S.SendVerificationEmail(user)
		log.Println("Verification Email Sending Success:", success)
		/*var link template.HTML
		link = "<a href='"+template.HTML(G.AppEnv.Url)+"'>Click Here To Resend</a>"*/
		G.Msg.Success = "Successfully Registered. Please Check Your Verification Email." /* If You Don't Get it " + link + ".*/

		/*session, _ := G.FStore.Get(c.Request, "login_token")
		session.Values["userEmail"] = user.Email
		//session.Values["remember_token"] = user.RememberToken.String
		err := session.Save(c.Request, c.Writer)
		if err != nil {
			log.Println(err.Error())
		}*/
		c.Redirect(http.StatusFound, "/login")

	} else {
		if G.Msg.Fail == "" {
			G.Msg.Fail = "Some Internal Server Error Occurred, Please Try Again."
		}
		c.Redirect(http.StatusFound, "/register")
	}
}

func ResendEmailVf(c *gin.Context) {
	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var user Mod.User
	user.Email = H.GetCookie("secret", nil, "email", c)
	if user.Email != "" {
		user, _ = R.ReadUserWithEmail(user)
		if user.ActiveStatus == 0 {
			success := S.SendVerificationEmail(user)
			log.Println("Verification Email Re Sending Success:", success)
			if success {
				G.Msg.Success = "Email Has Been Sent Successfully."
			}
		} else if user.ActiveStatus > 0 {
			G.Msg.Success = "Already Activated."
		} else {
			var link template.HTML
			link = "<a href='" + template.HTML(G.AppEnv.Url) + "'>Click Here To Resend</a>"
			G.Msg.Fail = "Email Sending Failed. " + link + "."
		}
	}
	c.Redirect(http.StatusFound, "/login")
}

func ActivateAccount(c *gin.Context) {
	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var user Mod.User
	encEmail := c.Query("encEmail")
	emailVf := c.Query("emailVf")
	var err error
	var decoded []byte
	decoded, err = base64.URLEncoding.DecodeString(encEmail)
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusOK, "404.html", nil)
		return
	}

	user.Email = string(decoded)
	user.EmailVerification = emailVf
	var success bool

	user, success = R.ActivateAccount(user, "reg_type=1 and email=? and email_verification=?", user.Email, user.EmailVerification)
	if success {
		G.Msg.Success = "Congratulations, Your Account Is Activated."
		c.Redirect(http.StatusFound, "/login")
	} else {
		c.HTML(http.StatusOK, "404.html", nil)
	}
}

func ForgotPassword(c *gin.Context) {

	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	if user, success := M.IsGuest(c, G.FStore); success {
		//convRate,_, currCode, currSym := S.CurrencyProcessing(c)

		finalCart := S.ProcessCart(c)

		lenCart := len(finalCart.Carts)

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		G.LangVal = H.GetCookie("secret", nil, "lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal

		lenWish := R.CountWishlist(user)

		c.HTML(http.StatusOK, "forgot-password.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Forgot Password", "Carts": lenCart, "Lang": lang,
			"Msg": G.Msg, "Title": "Forgot-Password", "Wbsts": wbsts, "LenWish": lenWish, "Menus": menus, "Adm": G.Adm,
		})
		G.Msg.Success = ""
		G.Msg.Fail = ""
	}
}

func SendPasswordResetLink(c *gin.Context) {
	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var user Mod.User
	var ps Mod.PasswordReset
	var success bool
	user.Email = c.PostForm("email")
	user, success = R.ReadUserWithEmail(user)
	if !success {
		G.Msg.Fail = "Email does not exist."
		c.Redirect(http.StatusFound, "/forgot-password")
		return
	}
	/*if user.ActiveStatus == 0 {
		H.SetCookie("secret", nil, user.Email, "email", 60*60*5, c)
		var link template.HTML
		link = "<a href='" + template.HTML(G.AppEnv.Url) + "/resend-email-verification'>Click Here To Send Verification Email</a>"
		G.Msg.Fail = "Please Activate Your Account, "+link+"."

		c.Redirect(http.StatusFound, "/login")
		return
	}*/
	if user.ActiveStatus == 2 {
		G.Msg.Fail = "User Is Suspended."
		c.Redirect(http.StatusFound, "/login")
		return
	}
	ps.Email = user.Email
	ps.Token = H.RandomString(60)
	ps.Code = H.RandomNumber(10)
	if !R.SendPasswordResetLink(ps) {
		return
	}
	if S.SendPasswordResetLinkEmail(user, ps) {
		G.Msg.Success = "Reset Password Link Sent Successfully. Check Your Email."
	}
	c.Redirect(http.StatusFound, "/login")
}

func ResetPasswordGet(c *gin.Context) {
	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	user, success := M.IsGuest(c, G.FStore)
	if !success {
		return
	}
	encEmail := c.Query("email")
	var err error
	var decoded []byte
	decoded, err = base64.URLEncoding.DecodeString(encEmail)
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusOK, "404.html", nil)
		return
	}
	var ps Mod.PasswordReset
	ps.Email = string(decoded)
	ps.Token = c.Query("token")
	if !R.ResetPasswordGet(ps) {
		c.HTML(http.StatusOK, "404.html", nil)
		return
	}
	//convRate,_, currCode, currSym := S.CurrencyProcessing(c)

	finalCart := S.ProcessCart(c)

	lenCart := len(finalCart.Carts)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	G.LangVal = H.GetCookie("secret", nil, "lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	lenWish := R.CountWishlist(user)

	c.HTML(http.StatusOK, "reset-password.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": "Reset Password", "Carts": lenCart, "LenWish": lenWish,
		"Msg": G.Msg, "Title": "Reset-Password", "PS": ps, "Wbsts": wbsts, "Menus": menus, "Lang": lang, "Adm": G.Adm,
	})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func ResetPasswordPost(c *gin.Context) {
	if user := M.GetAuthUser(c, G.FStore); user.ID > 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var user Mod.User
	var ps Mod.PasswordReset
	password := c.PostForm("password")
	confirmPass := c.PostForm("confirm_password")
	ps.Email = c.PostForm("email")
	ps.Token = c.PostForm("token")

	if password != confirmPass {
		G.Msg.Fail = "Confirm Password Doesn't Match."
		encEmail := base64.URLEncoding.EncodeToString([]byte(ps.Email))
		c.Redirect(http.StatusFound, "/reset-password?email="+encEmail+"&token="+ps.Token)
		return
	}
	if !R.IsResetTokenValid(ps) {
		G.Msg.Fail = "Token Invalid."
		log.Println(G.Msg.Fail)
		c.HTML(http.StatusOK, "404.html", G.Msg.Fail)
		return
	}

	cost := bcrypt.DefaultCost
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), cost)
	user.Password = string(hash)
	user.Email = ps.Email
	if !R.ResetPasswordPost(user, ps, "email=? and token=?", ps.Email, ps.Token) {
		G.Msg.Fail = "Some Internal Server Error Occurred, Please Reload And Try Again Later."
		encEmail := base64.URLEncoding.EncodeToString([]byte(ps.Email))
		c.Redirect(http.StatusFound, "/reset-password?email="+encEmail+"&token="+ps.Token)
		return
	}
	G.Msg.Success = "Your Password Is Reset Successfully."
	c.Redirect(http.StatusFound, "/login")
}

func Logout(c *gin.Context) {
	/*if user := M.GetAuthUser(c, G.FStore); user.ID == 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}*/

	session, _ := G.FStore.Get(c.Request, "login_token")
	session.Options.MaxAge = -1
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		log.Println(err.Error())
	}
	c.Redirect(http.StatusFound, "/login")
}
