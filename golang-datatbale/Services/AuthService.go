package Services

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"log"
	"net/http"
)

type Country struct {
	Key   string
	Value string
}

func SendVerificationEmail(user Mod.User) bool {
	encEmail := base64.URLEncoding.EncodeToString([]byte(user.Email))

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	/*data := map[string]interface{}{"Wbsts":wbsts,"Adm":G.Adm,
		"AppEnv": G.AppEnv, "EncEmail":encEmail, "User":user}
	body, success := DynamicEmailBodyProcessing("Account-Verification-Email", data)
	if !success {
		return false
	}*/

	htmlString, err := H.ParseTemplate("Views/Email/email-verify.html", map[string]interface{}{
		"Wbsts":wbsts,"Adm":G.Adm,"AppEnv": G.AppEnv, "EncEmail":encEmail, "User":user/*, "Body":body*/})
	if err != nil {
		log.Println(err.Error())
		return false
	}

	//From := G.AppEnv.Name+" <noreply@bidstore.com>"
	To := []string{user.Email}
	Subject := "Account Verification Email" //G.EmailTemplate["Account-Verification-Email"].Subject
	HtmlString := htmlString
	if !SendEmail(To, Subject, HtmlString, "") {
		/*var link template.HTML
		link = "<a href='"+template.HTML(G.AppEnv.Url)+"'>Click Here To Resend</a>"
		G.Msg.Fail = "Verification Email Not Sent, "+link+"."*/
		return false
	}
	return true
}



/*func SendVerificationEmail(user Mod.User) bool {
	encEmail := base64.URLEncoding.EncodeToString([]byte(user.Email))
	htmlString, err := H.ParseTemplate("Views/Email/email-verify.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "EncEmail":encEmail, "User":user})
	if err != nil {
		log.Println("AuthService.go Log1", err.Error())
		return false
	}

	//From := G.AppEnv.Name+" <noreply@goldstore.com>"
	To := user.Email
	Subject := "Account Verification Email"
	HtmlString := htmlString
	if !SendEmail(To, Subject, HtmlString, "") {
		G.Msg.Fail = "Verification Email Not Sent, <a href='http://localhost:2000/resend-email-verification'>Click Here To Resend</a>."
		return false
	}
	return true
}*/


func SendPasswordResetLinkEmail(user Mod.User, ps Mod.PasswordReset) bool {
	encEmail := base64.URLEncoding.EncodeToString([]byte(user.Email))

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	/*data := map[string]interface{}{"Wbsts":wbsts,"Adm":G.Adm,"AppEnv": G.AppEnv, "EncEmail":encEmail, "User":user, "PS":ps}
	body, success := DynamicEmailBodyProcessing("Reset-Password-Email", data)
	if !success {
		return false
	}*/

	htmlString, err := H.ParseTemplate("Views/Email/reset-password-email.html", map[string]interface{}{
		"Wbsts":wbsts,"Adm":G.Adm,"AppEnv": G.AppEnv, "EncEmail":encEmail, "User":user, "PS":ps, /*"Body":body*/})
	if err != nil {
		log.Println(err.Error())
		return false
	}

	//From := G.AppEnv.Name+" <noreply@bidstore.com>"
	To := []string{user.Email}
	Subject := "Reset Password Email" //G.EmailTemplate["Reset-Password-Email"].Subject
	HtmlString := htmlString
	if !SendEmail(To, Subject, HtmlString, "") {
		G.Msg.Fail = "Failed To Send Link, Please Reload And Try Again Later."
		return false
	}
	return true
}


/*func SendPasswordResetLinkEmail(user Mod.User, ps Mod.PasswordReset) bool {
	encEmail := base64.URLEncoding.EncodeToString([]byte(user.Email))

	htmlString, err := H.ParseTemplate("Views/Email/reset-password-email.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "EncEmail":encEmail, "User":user, "PS":ps})
	if err != nil {
		log.Println("AuthService.go Log2", err.Error())
		return false
	}

	//From := G.AppEnv.Name+" <noreply@goldstore.com>"
	To := user.Email
	Subject := "Reset Password Link"
	HtmlString := htmlString
	if !SendEmail(To, Subject, HtmlString, "") {
		G.Msg.Fail = "Failed To Send Link, Please Try Again Later."
		return false
	}
	return true
}*/

func SendEmailChangeMail(user Mod.User) bool {

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	/*data := map[string]interface{}{"Wbsts":wbsts,"Adm":G.Adm,
		"AppEnv": G.AppEnv, "EncEmail":encEmail, "User":user}
	body, success := DynamicEmailBodyProcessing("Account-Verification-Email", data)
	if !success {
		return false
	}*/

	htmlString, err := H.ParseTemplate("Views/Email/email-change.html", map[string]interface{}{
		"Wbsts":wbsts,"Adm":G.Adm,"AppEnv": G.AppEnv, "User":user/*, "Body":body*/})
	if err != nil {
		log.Println(err.Error())
		return false
	}

	//From := G.AppEnv.Name+" <noreply@bidstore.com>"
	To := []string{user.OldEmail, user.Email}
	Subject := "Email Change" //G.EmailTemplate["Account-Verification-Email"].Subject
	HtmlString := htmlString
	if !SendEmail(To, Subject, HtmlString, "") {
		/*var link template.HTML
		link = "<a href='"+template.HTML(G.AppEnv.Url)+"'>Click Here To Resend</a>"
		G.Msg.Fail = "Verification Email Not Sent, "+link+"."*/
		return false
	}
	return true
}

func RedirectToProfileComplete(c *gin.Context, user Mod.User) {
	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	G.LangVal = H.GetCookie("secret", nil, "lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	var countries []Country
	var country Country
	for key, value := range G.Country {
		country.Key = key
		country.Value = value
		countries = append(countries, country)
	}

	c.HTML(http.StatusOK, "profile-complete.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Lang": lang,"Countries": countries,
		"Msg": G.Msg, "Title": "Profile-Completion", "Wbsts": wbsts,
		"Adm": G.Adm})
	G.Msg.Success = ""
	G.Msg.Fail = ""
	return
}