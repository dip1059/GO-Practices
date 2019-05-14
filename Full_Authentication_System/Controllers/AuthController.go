package Controllers

import (
	G "PracticeGoland/Globals"
	H "PracticeGoland/Helpers"
	M "PracticeGoland/Middlewares"
	R "PracticeGoland/Repositories"
	S "PracticeGoland/Services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strconv"
)

var (
	sc = securecookie.New([]byte("secret"),nil)
	store = sessions.NewCookieStore([]byte("secret"))
)

func Welcome(c *gin.Context) {
	if M.IsGuest(c, store, sc) {
		c.HTML(http.StatusOK, "welcome.html", nil)
	}
	return
}

func RegisterGet(c *gin.Context)  {
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
	hash,_ := bcrypt.GenerateFromPassword([]byte(password), cost)
	G.User.Password = string(hash)
	G.User.EmailVf.String = H.RandomString(60)
	G.User.EmailVf.Valid = true
	var success bool
	G.User, success = R.Register(G.User)
	if success {
		if S.SendVerificationEmail() {
			var link template.HTML
			link = "<a href='http://localhost:2000/resend-email-verification'>Click Here To Resend</a>"
			G.Msg.Success = "Successfully Registered. Please Check Your Verification Email. If You Don't Get it "+link+"."
		}
		c.Redirect(http.StatusFound, "/register")
	} else {
		G.Msg.Fail = "Some error occurred, please try again."
		c.Redirect(http.StatusFound, "/register")
	}
}

func ResendEmailVf(c *gin.Context) {
	if S.ResendVerificationEmail() {
		var link template.HTML
		link = "<a href='http://localhost:2000/resend-email-verification'>Click Here To Resend</a>"
		G.Msg.Success = "Email Has Been Sent Successfully. If You Don't Get it "+link+"."
	}
	c.Redirect(http.StatusFound, "/register")
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
	G.User, success = R.Read(G.User)
	if success {
		err := bcrypt.CompareHashAndPassword([]byte(G.User.Password), []byte(password))
		if err != nil {
			G.Msg.Fail = "Wrong Credentials."
			c.Redirect(http.StatusFound, "/login")
		} else {
			session, _ := store.Get(c.Request, "login_token")
			session.Values["userId"] = G.User.ID
			session.Values["userName"] = G.User.Name
			session.Values["userEmail"] = G.User.Email
			session.Values["userRole"] = G.User.Role
			session.Values["userActiveStatus"] = G.User.ActiveStauts
			session.Save(c.Request, c.Writer)
			if rememberToken == 1 {
				S.SetRememberToken(c, sc)
			}
			c.Redirect(http.StatusFound, "/home")
		}
	} else {
		G.Msg.Fail = "User Not Found Or Some Internal Server Error Occured."
		c.Redirect(http.StatusFound, "/login")
	}
}

func Home( c *gin.Context) {
	if M.IsAuthUser(c, store, sc) {
		/*session, _ := store.Get(c.Request, "login_token")
		data := session.Values["userId"]
		if  data == nil{
			val,_ := c.Cookie("remember_token")
			if err := sc.Decode("remember_token", val, &val); err == nil {
				data = val
			}
		}*/
		c.HTML(http.StatusOK, "User/home.html",G.User)
	} else {
		G.Msg.Fail = "You are not Logged in."
	}
}

func Dashboard( c *gin.Context) {
	if M.IsAuthUser(c, store, sc) {
		c.HTML(http.StatusOK, "Admin/dashboard.html",G.User)
	} else {
		G.Msg.Fail = "You are not Logged in."
	}
}

func Logout(c *gin.Context) {
	session, _ := store.Get(c.Request, "login_token")
	session.Save(c.Request, c.Writer)
	cookie := http.Cookie{
		Name:     "login_token",
		MaxAge:   -1,
	}
	http.SetCookie(c.Writer, &cookie)

	cookie2 := http.Cookie{
		Name:     "remember_token",
		MaxAge:   -1,
	}
	http.SetCookie(c.Writer, &cookie2)
	c.Redirect(http.StatusFound, "/login")
}





