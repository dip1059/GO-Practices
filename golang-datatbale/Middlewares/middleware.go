package Middlewares

import (
	G "gold-store/Globals"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
)

func GetAuthUser(c *gin.Context, store *sessions.FilesystemStore) Mod.User {
	var user Mod.User
	var success bool

	session, _ := store.Get(c.Request, "login_token")
	email := session.Values["userEmail"]

	if email != nil {
		user.Email = session.Values["userEmail"].(string)
		user, success = R.ReadUserWithEmail(user)
		if !success {
			c.Redirect(http.StatusFound, "/logout")
			return user
		}
	}
	return user
}

func IsGuest(c *gin.Context, store *sessions.FilesystemStore) (Mod.User, bool) {
	var user Mod.User

	session, _ := store.Get(c.Request, "login_token")
	email := session.Values["userEmail"]

	var success bool
	if email != nil {
		user.Email = session.Values["userEmail"].(string)
		user, success = R.ReadUserWithEmail(user)

		if !success {
			c.Redirect(http.StatusFound, "/logout")
			return user, true
		}

		if user.ActiveStatus == 2 {
			c.Redirect(http.StatusFound, "/logout")
			return user, false
		} else if user.ActiveStatus == 1 && user.RoleID == 1 {
			c.Redirect(http.StatusFound, "/dashboard")
			return user, false
		} else if user.ActiveStatus == 1 && user.RoleID == 2 {
			return user, true
		}/* else if user.ActiveStatus == 1 && user.RoleID == 3 {
			return user, true
		}*/

	}
	return user, true
}


func IsAuthUser(c *gin.Context, store *sessions.FilesystemStore) (Mod.User, bool) {
	var user Mod.User

	session, _ := store.Get(c.Request, "login_token")
	email := session.Values["userEmail"]

	var success bool
	if email != nil {
		user.Email = session.Values["userEmail"].(string)
		user, success = R.ReadUserWithEmail(user)

		if !success {
			G.Msg.Fail = "User Doesn't Exist Anymore."
			c.Redirect(http.StatusFound, "/logout")
			return user, false
		}
		if user.ActiveStatus == 2 {
			G.Msg.Fail = "You Are Suspended. Contact With The Authority Quickly."
			c.Redirect(http.StatusFound, "/logout")
		} else if user.ActiveStatus == 1 && user.RoleID == 1 {
			c.Redirect(http.StatusFound, "/dashboard")
		} else if user.ActiveStatus == 1 && user.RoleID == 2 {
			return user, true
		}/* else if user.ActiveStatus == 1 && user.RoleID == 3 {
			c.Redirect(http.StatusFound, "/vendor/dashboard")
		}*/
		return user, false
	}
	return user, false

}


/*func IsAuthVendorUser(c *gin.Context, store *sessions.FilesystemStore) (Mod.User, bool) {
	var user Mod.User

	session, _ := store.Get(c.Request, "login_token")
	email := session.Values["userEmail"]

	var success bool
	if email != nil {
		user.Email = session.Values["userEmail"].(string)
		user, success = R.ReadUserWithEmail(user)

		if !success {
			G.Msg.Fail = "User Doesn't Exist Anymore."
			c.Redirect(http.StatusFound, "/logout")
			return user, false
		}
		if user.ActiveStatus == 2 {
			G.Msg.Fail = "You Are Suspended. Contact With The Authority Quickly."
			c.Redirect(http.StatusFound, "/logout")
		} else if user.ActiveStatus == 1 && user.RoleID == 1 {
			c.Redirect(http.StatusFound, "/dashboard")
		} else if user.ActiveStatus == 1 && user.RoleID == 2 {
			c.Redirect(http.StatusFound, "/user/dashboard")
		} else if user.ActiveStatus == 1 && user.RoleID == 3 {
			return user, true
		}
		return user, false
	}
	return user, false

}
*/

func IsAuthAdminUser(c *gin.Context, store *sessions.FilesystemStore) (Mod.User, bool) {
	var user Mod.User
	session, _ := store.Get(c.Request, "login_token")
	email := session.Values["userEmail"]

	var success bool

	if email != nil {
		user.Email = session.Values["userEmail"].(string)
		user, success = R.ReadUserWithEmail(user)

		if !success {
			G.Msg.Fail = "User Doesn't Exist Anymore."
			c.Redirect(http.StatusFound, "/logout")
			return user, false
		}
		if user.ActiveStatus == 2 {
			G.Msg.Fail = "You Are Suspended. Contact With The Authority Quickly."
			c.Redirect(http.StatusFound, "/logout")
		} else if user.ActiveStatus == 1 && user.RoleID == 2 {
			c.Redirect(http.StatusFound, "/")
		/*} else if user.ActiveStatus == 1 && user.RoleID == 3 {
			c.Redirect(http.StatusFound, "/vendor/dashboard")*/
		} else if user.ActiveStatus == 1 && user.RoleID == 1 {
			return user, true
		}
		return user, false
	}
	c.Redirect(http.StatusFound, "/login")
	return user, false
}


/*func GenerateJWT(id uint) (string, int64, bool) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = id
	claims["iat"] = time.Now().Unix()
	//claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	key := os.Getenv("JWT_KEY")
	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Something Went Wrong: %s", err.Error())
		return "", 0, false
	}

	return tokenString, claims["iat"].(int64), true
}*/


/*func IsTokenValid(c *gin.Context) (Mod.User, bool) {
	var user Mod.User

	if c.Request.Header["Token"] != nil {

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(c.Request.Header["Token"][0], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing mehod error")
			}
			return []byte(os.Getenv("JWT_KEY")), nil
		})
		if err != nil {
			log.Println(err.Error())
			return user, false
		}

		if token.Valid {
			var apiSes Mod.ApiSession
			apiSes.UserID = uint(claims["client"].(float64))
			apiSes.IssuedAt = int64(claims["iat"].(float64))
			apiSes.Token = token.Raw
			apiSes.User.ID = apiSes.UserID

			apiSes, success := R.GetApiSession(apiSes, "token=? and user_id=? and issued_at=?", apiSes.Token, apiSes.UserID, apiSes.IssuedAt)
			if success {
				if apiSes.User.ActiveStatus == 1 {
					return apiSes.User, true
				} else {
					R.DeleteApiSession(apiSes, "user_id=?", apiSes.UserID)
					return user, false
				}
			} else {
				return user, false
			}
		}
	} else {
		return user, false
	}
	return user, false
}*/