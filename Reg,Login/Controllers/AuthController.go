package Controllers

import (
	G "PracticeGoland/Globals"
	//H "PracticeGoland/Helpers"
	M "PracticeGoland/Middlewares"
	Mod "PracticeGoland/Models"
	R "PracticeGoland/Repositories"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

var (
	sc = securecookie.New([]byte("secret"),nil)
	Msg G.Message
	user Mod.User
	store = sessions.NewCookieStore([]byte("secret"))
)

func Welcome(c *gin.Context) {
	if M.Guest(c, store) {
		c.HTML(http.StatusOK, "welcome.html", nil)
	}
	return
}

func RegisterGet(c *gin.Context)  {
	if M.Guest(c, store) {
		c.HTML(http.StatusOK, "register.html", Msg)
		G.Msg.Success = ""
		G.Msg.Fail = ""
	}
	return
}

func RegisterPost(c *gin.Context) {
	user.Name = c.PostForm("full_name")
	user.Email = c.PostForm("email")
	password := c.PostForm("password")
	cost := bcrypt.DefaultCost
	hash,_ := bcrypt.GenerateFromPassword([]byte(password), cost)
	user.Password = string(hash)

	success := R.Register(user)
	if success {
		Msg.Success = "Successfully Added."
		c.Redirect(http.StatusFound, "/register")
	} else {
		Msg.Fail = "Some error occurred, please try again."
		c.Redirect(http.StatusFound, "/register")
	}
}

func LoginGet(c *gin.Context) {
	if M.Guest(c, store) {
		c.HTML(http.StatusOK, "login.html", Msg)
		Msg.Success = ""
		Msg.Fail = ""
	}
	return
}

func LoginPost(c *gin.Context) {
	user.Email = c.PostForm("email")
	password := c.PostForm("password")
	rememberToken, _ := strconv.Atoi(c.PostForm("remember_token"))
	var success bool
	user, success = R.Login(user)
	if success {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			Msg.Fail = "Wrong Credentials."
			c.Redirect(http.StatusFound, "/login")
		} else {
			session, _ := store.Get(c.Request, "login_token")
			session.Values["data"] = user.Name+"  ID: "+strconv.Itoa(user.ID)+" Email: "+user.Email
			session.Save(c.Request, c.Writer)
			if rememberToken == 1 {
				SetRememberToken(c)
			}
			c.Redirect(http.StatusFound, "/home")
		}
	} else {
		Msg.Fail = "User not found or Some Internal Server error occured."
		c.Redirect(http.StatusFound, "/login")
	}

}

func Home( c *gin.Context) {
	if M.AuthUser(c, store) {
		session, _ := store.Get(c.Request, "login_token")
		data := session.Values["data"];
		if  data == nil{
			val,_ := c.Cookie("remember_token")
			if err := sc.Decode("remember_token", val, &val); err == nil {
				data = val
			}
		}
		c.String(http.StatusOK, "Welcome Home "+data.(string))
	} else {
		Msg.Fail = "You are not Logged in."
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

func SetRememberToken(c *gin.Context) {
	val := user.Name+"  ID: "+strconv.Itoa(user.ID)+" Email: "+user.Email
	encoded, _ := sc.Encode("remember_token", val)

	cookie1 := http.Cookie{
		Name:     "remember_token",
		Value:    encoded,
		MaxAge:   60 * 60 * 24 * 365,
	}

	http.SetCookie(c.Writer, &cookie1)
}

