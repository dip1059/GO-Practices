package Config

import(
	G "gold-store/Globals"
	"github.com/joho/godotenv"
	"os"
)

func AppConfig() {
	godotenv.Load()
	G.AppEnv.Name = os.Getenv("APP_NAME")
	G.AppEnv.Url = os.Getenv("APP_URL")
	G.AppEnv.Port = os.Getenv("APP_PORT")
	G.AppEnv.Debug = os.Getenv("APP_DEBUG")
}
