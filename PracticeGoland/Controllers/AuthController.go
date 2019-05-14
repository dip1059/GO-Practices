package Controllers

import (
	G "PracticeGoland/Globals"
	H "PracticeGoland/Helpers"
	M "PracticeGoland/Middlewares"
	R "PracticeGoland/Repositories"
	S "PracticeGoland/Services"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var (
	sc    = securecookie.New([]byte("secret"), nil)
	store = sessions.NewCookieStore([]byte("secret"))
)

func Welcome(c *gin.Context) {
	if M.IsGuest(c, store, sc) {
		c.HTML(http.StatusOK, "welcome.html", nil)
	}
	return
}

func RegisterGet(c *gin.Context) {
	if M.IsGuest(c, store, sc) {
		c.HTML(http.StatusOK, "register.html", G.Msg)
		G.Msg.Success = ""
		G.Msg.Fail = ""
	}
	return
}

func RegisterPost(c *gin.Context) {
	G.User.Name = c.PostForm("full_name")
	G.User.Email = c.PostForm("email")
	password := c.PostForm("password")
	cost := bcrypt.DefaultCost
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), cost)
	G.User.Password = string(hash)
	G.User.EmailVf.String = H.RandomString(60)
	G.User.EmailVf.Valid = true
	var success bool
	G.User, success = R.Register(G.User)
	if success {
		if S.SendVerificationEmail() {
			var link template.HTML
			link = "<a href='http://localhost:2000/resend-email-verification'>Click Here To Resend</a>"
			G.Msg.Success = "Successfully Registered. Please Check Your Verification Email. If You Don't Get it " + link + "."
		}
		c.Redirect(http.StatusFound, "/register")
	} else {
		if G.Msg.Fail == "" {
			G.Msg.Fail = "Some Internal Server Error Occurred, Please Try Again."
		}
		c.Redirect(http.StatusFound, "/register")
	}
}

func ResendEmailVf(c *gin.Context) {
	if G.User.Email != "" {
		if G.User.ActiveStauts == 0 {
			if S.SendVerificationEmail() {
				G.Msg.Success = "Email Has Been Sent Successfully."
			}
		} else {
			G.Msg.Success = "Already Activated."
		}
	}
	c.Redirect(http.StatusFound, "/login")
}


func ActivateAccount(c * gin.Context) {
	encEmail := c.Param("encEmail")
	emailVf := c.Param("emailVf")
	var err error
	var decoded []byte
	decoded, err = base64.URLEncoding.DecodeString(encEmail)
	if err != nil {
		log.Println("AuthController.go Log1", err.Error())
		c.HTML(http.StatusOK, "404.html",nil)
		return
	}

	G.User.Email = string(decoded)
	G.User.EmailVf.String = emailVf
	var success bool

	G.User, success = R.ActivateAccount(G.User)
	if success {
		G.Msg.Success = "Congratulations, Your Account Is Activated."
		c.Redirect(http.StatusFound, "/login")
	} else {
		c.HTML(http.StatusOK, "404.html",nil)
	}
}


func LoginGet(c *gin.Context) {
	if M.IsGuest(c, store, sc) {
		c.HTML(http.StatusOK, "login.html", G.Msg)
		G.Msg.Success = ""
		G.Msg.Fail = ""
	}
	return
}

func LoginPost(c *gin.Context) {
	G.User.Email = c.PostForm("email")
	password := c.PostForm("password")
	rememberToken, _ := strconv.Atoi(c.PostForm("remember_token"))
	var success bool
	G.User, success = R.Login(G.User)
	if success {

		err := bcrypt.CompareHashAndPassword([]byte(G.User.Password), []byte(password))
		if err != nil {
			G.Msg.Fail = "Wrong Credentials."
			c.Redirect(http.StatusFound, "/login")
		} else {
			if G.User.ActiveStauts == 1 {
				session, _ := store.Get(c.Request, "login_token")
				session.Values["userEmail"] = G.User.Email
				session.Save(c.Request, c.Writer)
				if rememberToken == 1 {
					S.SetRememberToken(c, sc)
				}
				if G.User.Role == 1 {
					c.Redirect(http.StatusFound, "/dashboard")
				} else if G.User.Role == 2 {
					c.Redirect(http.StatusFound, "/home")
				}
			} else if G.User.ActiveStauts == 2 {
				if G.Msg.Fail == "" {
					G.Msg.Fail = "You Are Suspended. Contact With The Authority Quickly."
				}
				c.Redirect(http.StatusFound, "/login")
			} else {
				var link template.HTML
				link = "<a href='http://localhost:2000/resend-email-verification'>Click Here To Send Verification Email</a>"
				if G.Msg.Fail == "" {
					G.Msg.Fail = "Please Activate Your Account, " + link + "."
				}
				c.Redirect(http.StatusFound, "/login")
			}
		}

	} else {
		if G.Msg.Fail == "" {
			G.Msg.Fail = "User Not Found."
		}
		c.Redirect(http.StatusFound, "/login")
	}
}

func Home(c *gin.Context) {
	if M.IsAuthUser(c, store, sc) {
		c.HTML(http.StatusOK, "home.html", G.User)
	}
	return
}

func Dashboard(c *gin.Context) {
	if M.IsAuthAdminUser(c, store, sc) {
		c.HTML(http.StatusOK, "dashboard.html", G.User)
	}
	return
}

func Logout(c *gin.Context) {
	session, _ := store.Get(c.Request, "login_token")
	session.Save(c.Request, c.Writer)
	cookie := http.Cookie{
		Name:   "login_token",
		MaxAge: -1,
	}
	http.SetCookie(c.Writer, &cookie)

	cookie2 := http.Cookie{
		Name:   "remember_token",
		MaxAge: -1,
	}
	http.SetCookie(c.Writer, &cookie2)
	c.Redirect(http.StatusFound, "/login")
}
