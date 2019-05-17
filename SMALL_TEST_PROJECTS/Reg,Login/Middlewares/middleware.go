package Middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
)

func Guest(c *gin.Context, store *sessions.CookieStore) bool {
	session, _ := store.Get(c.Request, "login_token")
	data := session.Values["data"]
	_, err := c.Cookie("remember_token")
	if data != nil || err == nil {
		c.Redirect(http.StatusFound, "/home")
		return false
	}
	return true
}

func AuthUser(c *gin.Context, store *sessions.CookieStore) bool {
	session, _ := store.Get(c.Request, "login_token")
	data := session.Values["data"]
	_, err := c.Cookie("remember_token")
	if data != nil || err == nil {
		return true
	} else {
		c.Redirect(http.StatusFound, "/login")
		return false
	}
}
