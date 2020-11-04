package Services

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	//Cfg "gold-store/Config"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	//M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"log"
	"strconv"
)

func AddressValidation(c *gin.Context) (bool, Mod.TransferShippingAddress) {
	var shipAddress Mod.TransferShippingAddress
	err := c.ShouldBind(&shipAddress)

	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill all the required fields with valid values."
		return false, shipAddress
	}

	_, err = strconv.Atoi(shipAddress.Phone)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Not a valid phone number."
		return false, shipAddress
	}

	_, err = strconv.Atoi(shipAddress.ZipCode)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Not a valid zip code."
		return false, shipAddress
	}
	return true, shipAddress
}

func SendGoldTransferEmail(transfer Mod.GoldTransfer, user Mod.User, c *gin.Context, filename string) bool {
	encEmail := base64.URLEncoding.EncodeToString([]byte(user.Email))

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	G.LangVal = H.GetCookie("secret", nil, "lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	var view, Subject string
	if transfer.Status == 0 {
		view = "Views/Email/gold-transfer-confirm.html"
		Subject = "Gold Transfer Verification"
	} else if transfer.Status > 0 {
		view = "Views/Email/gold-transfer-email.html"
		Subject = "Gold Transfer"
	}

	htmlString, err := H.ParseTemplate(view, map[string]interface{}{"EncEmail": encEmail,
		"User": user, "Wbsts": wbsts, "Transfer": transfer, "Adm": G.Adm, "Lang": lang, "AppEnv": G.AppEnv})

	if err != nil {
		log.Println(err.Error())
		return false
	}

	To := []string{user.Email}
	HtmlString := htmlString
	if !SendEmail(To, Subject, HtmlString, filename) {
		return false
	}
	return true
}

func ProcessGoldTransfer(gt Mod.GoldTransfer, user Mod.User, c *gin.Context) bool {
	count := R.CountWalletCertificate(user.ID)
	log.Println(count, int(user.Wallet.GoldAmount*1000), count - int(user.Wallet.GoldAmount*1000))

	if (count - int(user.Wallet.GoldAmount*1000)) > 1 ||  (count - int(user.Wallet.GoldAmount*1000)) < -1 {
		G.Msg.Fail = "You have a balance mismatch."
		return false
	} else if count - int(user.Wallet.GoldAmount*1000) == 1 {
		UpdateWalletGoldAmount(nil, 0.001, user.ID, "fraction mismatch solution")
	} else if count - int(user.Wallet.GoldAmount*1000) == -1 {
		UpdateWalletGoldAmount(nil, -0.001, user.ID, "fraction mismatch solution")
	}

	count = R.CountWalletCertificate(gt.ReceiverUser.ID)
	gt.ReceiverUser,_ = R.ReadUserWithEmail(gt.ReceiverUser)
	log.Println(count, int(gt.ReceiverUser.Wallet.GoldAmount*1000), count - int(gt.ReceiverUser.Wallet.GoldAmount*1000))

	if (count - int(gt.ReceiverUser.Wallet.GoldAmount*1000)) > 1 || (count - int(gt.ReceiverUser.Wallet.GoldAmount*1000)) < -1 {
		G.Msg.Fail = "Receiver has a balance mismatch."
		return false
	} else if count - int(gt.ReceiverUser.Wallet.GoldAmount*1000) == 1 {
		UpdateWalletGoldAmount(nil, 0.001, gt.ReceiverUser.ID, "from fraction mismatch solution")
	} else if count - int(gt.ReceiverUser.Wallet.GoldAmount*1000) == -1 {
		UpdateWalletGoldAmount(nil, -0.001, gt.ReceiverUser.ID, "from fraction mismatch solution")
	}

	go func() {
		TransferCertificate(gt, user, c)
	}()

	return true
}
