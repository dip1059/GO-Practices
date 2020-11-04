package Config

import(
	G "gold-store/Globals"
	"github.com/joho/godotenv"
	"os"
)

func SocialAuthConfig() {
	godotenv.Load()
	G.SocialEnv.FacebookClientID = os.Getenv("FACEBOOK_CLIENT_ID")
	G.SocialEnv.FacebookClientSecret = os.Getenv("FACEBOOK_CLIENT_SECRET")
	G.SocialEnv.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	G.SocialEnv.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
}
