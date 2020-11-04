package Routes

import (
	G "gold-store/Globals"
	"gold-store/Routes/Admin"
	"gold-store/Routes/Auth"
	"gold-store/Routes/User"
	"github.com/bykovme/gotrans"
	"github.com/gin-gonic/gin"
	"log"
)

func Routes() {

	err := gotrans.InitLocales("Langs")
	if err != nil {
		log.Println(err.Error())
	}

	if G.AppEnv.Debug == "false" || G.AppEnv.Debug == "" {
		gin.SetMode("release")
	}
	r := gin.Default()

	r.Static("/assets", "./")
	r.LoadHTMLGlob("Views/**/*.html")

	Auth.AuthRoutes(r)
	Admin.AdminRoutes(r)
	User.UserRoutes(r)

	if G.AppEnv.Port == "" {
		G.AppEnv.Port = "5500"
	}
	err = r.Run(":"+G.AppEnv.Port)
	if err != nil {
		log.Println(err.Error())
	}

}
