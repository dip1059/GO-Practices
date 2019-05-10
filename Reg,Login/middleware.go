package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Guest(c *gin.Context) {
	session, _ := store.Get(c.Request, "auth-cookie")
	userName := session.Values["userName"]
	if userName != nil {
		c.Redirect(http.StatusFound, "/home")
	}
	return
}

func AuthUser(c *gin.Context) {
	session, _ := store.Get(c.Request, "auth-cookie")
	userName := session.Values["userName"]
	if userName == nil {
		Msg.Fail = "You are not Logged in."
		c.Redirect(http.StatusFound, "/login")
	}
	return
}
