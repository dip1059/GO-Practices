package Middlewares

import (
	G "PracticeGoland/Globals"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	R "PracticeGoland/Repositories"
)

func IsGuest(c *gin.Context, store *sessions.CookieStore, sc *securecookie.SecureCookie) bool {
	session, _ := store.Get(c.Request, "login_token")
	data := session.Values["userId"]
	cookie, err := c.Cookie("remember_token")
	if data != nil || err == nil {
		if data != nil {
			G.User.ID = session.Values["userId"].(int)
			G.User.Email = session.Values["userEmail"].(string)
			G.User.ActiveStauts = session.Values["userActiveStatus"].(int)
			G.User.Role = session.Values["userRole"].(int)
		} else {
			if err := sc.Decode("remember_token", cookie, &cookie); err == nil {
				G.User.Email = cookie
				G.User, _= R.Read(G.User)
			}
		}
		if G.User.ActiveStauts == 1 && G.User.Role == 1 {
			c.Redirect(http.StatusFound, "/dashboard")
		} else if G.User.ActiveStauts == 1 && G.User.Role == 2  {
			c.Redirect(http.StatusFound, "/home")
		}
		return false
	}
	return true

}

func IsAuthUser(c *gin.Context, store *sessions.CookieStore, sc *securecookie.SecureCookie) bool {
	session, _ := store.Get(c.Request, "login_token")
	data := session.Values["userId"]

	cookie, err := c.Cookie("remember_token")

	if data != nil || err == nil {
		if data != nil {
			G.User.ID = session.Values["userId"].(int)
			G.User.Email = session.Values["userEmail"].(string)
			G.User.ActiveStauts = session.Values["userActiveStatus"].(int)
			G.User.Role = session.Values["userRole"].(int)
		} else {
			if err := sc.Decode("remember_token", cookie, &cookie); err == nil {
				G.User.Email = cookie
				G.User, _= R.Read(G.User)
			}
		}
		if G.User.ActiveStauts == 1 {
			return true
		}
		var link template.HTML
		link = "<a href='http://localhost:2000/resend-email-verification'>Click Here To Send Verification Email</a>"
		G.Msg.Fail = "Please Activate Your Account"+link+"."
		return false
	} else {
		c.Redirect(http.StatusFound, "/login")
		return false
	}
	return false
}
