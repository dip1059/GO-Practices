package Auth

import (
	"gold-store/Controllers/Auth"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.GET("/register", Auth.RegisterGet)
	r.POST("/register", Auth.RegisterPost)
	r.GET("/resend-email-verification", Auth.ResendEmailVf)
	r.GET("/email-verify", Auth.ActivateAccount)
	r.GET("/login", Auth.LoginGet)
	r.POST("/login", Auth.LoginPost)
	r.GET("/forgot-password", Auth.ForgotPassword)
	r.POST("/send-password-reset-link", Auth.SendPasswordResetLink)
	r.GET("/reset-password", Auth.ResetPasswordGet)
	r.POST("/reset-password", Auth.ResetPasswordPost)
	r.GET("/logout", Auth.Logout)
	//social auth
	r.GET("/social-login", Auth.SocialLogin)
	r.GET("/social-login-callback", Auth.SocialLoginCallback)
	r.POST("/social-register", Auth.SocialRegister)
}
