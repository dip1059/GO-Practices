package Middlewares

import (
	G "PracticeGoland/Globals"
	R "PracticeGoland/Repositories"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

func IsGuest(c *gin.Context, store *sessions.CookieStore, sc *securecookie.SecureCookie) bool {
	session, _ := store.Get(c.Request, "login_token")
	data := session.Values["userId"]
	cookie, err := c.Cookie("remember_token")
	if data != nil || err == nil {
		if data != nil {
			G.User.Email = session.Values["userEmail"].(string)
			G.User, _ = R.ReadWithEmail(G.User)
		} else {
			if err := sc.Decode("remember_token", cookie, &cookie); err == nil {
				G.User.Email = cookie
				G.User, _ = R.ReadWithEmail(G.User)
			} else {
				log.Println("middleware.go log1", err.Error())
				G.Msg.Fail = "Internal Server Error Occurred. Please Try Again Later."
				return false
			}
		}
		if G.User.ActiveStauts == 2 {
			G.Msg.Fail = "You Are Suspended. Contact With The Authority Quickly."
			c.Redirect(http.StatusFound, "/logout")
		} else if G.User.ActiveStauts == 1 && G.User.Role == 1 {
			c.Redirect(http.StatusFound, "/dashboard")
		} else if G.User.ActiveStauts == 1 && G.User.Role == 2 {
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
			G.User.Email = session.Values["userEmail"].(string)
			G.User, _ = R.ReadWithEmail(G.User)
		} else {
			if err := sc.Decode("remember_token", cookie, &cookie); err == nil {
				G.User.Email = cookie
				G.User, _ = R.ReadWithEmail(G.User)
			} else {
				log.Println("middleware.go log1", err.Error())
				G.Msg.Fail = "Internal Server Error Occurred. Please Try Again Later."
				return false
			}
		}
		if G.User.ActiveStauts == 1 && G.User.Role == 2{
			return true
		} else if G.User.ActiveStauts == 2 {
			G.Msg.Fail = "You Are Suspended. Contact With The Authority Quickly."
			c.Redirect(http.StatusFound, "/logout")
			return false
		} else {
			c.Redirect(http.StatusFound, "/login")
			return false
		}
	} else {
		c.Redirect(http.StatusFound, "/login")
		return false
	}
}


func IsAuthAdminUser(c *gin.Context, store *sessions.CookieStore, sc *securecookie.SecureCookie) bool {
	session, _ := store.Get(c.Request, "login_token")
	data := session.Values["userId"]

	cookie, err := c.Cookie("remember_token")

	if data != nil || err == nil {
		if data != nil {
			G.User.Email = session.Values["userEmail"].(string)
			G.User, _ = R.ReadWithEmail(G.User)
		} else {
			if err := sc.Decode("remember_token", cookie, &cookie); err == nil {
				G.User.Email = cookie
				G.User, _ = R.ReadWithEmail(G.User)
			} else {
				log.Println("middleware.go log1", err.Error())
				G.Msg.Fail = "Internal Server Error Occurred. Please Try Again Later."
				return false
			}
		}
		if G.User.ActiveStauts == 1 && G.User.Role == 1 {
			return true
		} else if G.User.ActiveStauts == 2 {
			G.Msg.Fail = "You Are Suspended. Contact With The Authority Quickly."
			c.Redirect(http.StatusFound, "/logout")
			return false
		} else {
			c.Redirect(http.StatusFound, "/login")
			return false
		}
	} else {
		c.Redirect(http.StatusFound, "/login")
		return false
	}
}
