package Services

import (
	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
	G "gold-store/Globals"
	"gopkg.in/danilopolani/gocialite.v1"
	"log"
)

var (
	gocial = gocialite.NewDispatcher()
)

func SocialLogin(driver string, clientID  string, clientSecret string, callBackUrl string) (string, bool) {

	authURL, err := gocial.New().
		Driver(driver). // Set provider
		//Scopes([]string{"public_repo"}). // Set optional scope(s)
		Redirect( //
			clientID, // Client ID
			clientSecret, // Client Secret
			G.AppEnv.Url+callBackUrl, // Redirect URL
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		log.Println(err.Error())
		return authURL, false
	}
	return authURL, true
}

func SocialLoginCallback(state string, code string) (*structs.User, *oauth2.Token, bool) {


	// Handle callback and check for errors
	user, token, err := gocial.Handle(state, code)
	if err != nil {
		log.Println(err.Error())
		return user, token, false
	}
	return user, token, true
}