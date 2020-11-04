package User

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	S "gold-store/Services"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func TransfersDetails(c *gin.Context) {
	user, success := M.IsAuthUser(c, G.FStore)
	if !success {
		return
	}
	id := c.Param("id")
	var gt Mod.GoldTransfer
	gt = R.GoldTransfer(gt, "sender_wallet_id=? and transfer_id=?", user.Wallet.ID, id)
	if gt.SenderWalletID != user.Wallet.ID {
		return
	}
	user.GoldTransfers = append(user.GoldTransfers, gt)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var menus []Mod.Menu
	menus = R.Menus(menus, "status=?", 1)

	session, _ := G.Store.Get(c.Request, "cart")
	lenCart := len(session.Values)

	G.LangVal = H.GetCookie("secret", nil, "lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	c.HTML(http.StatusOK, "transfer-details.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Banner": "Account", "Carts": lenCart, "Lang": lang,
		"Msg": G.Msg, "Adm": G.Adm, "Title": "Account", "Active": "Account",
		"Wbsts": wbsts, "Menus": menus, "NavActive": "Wallet",
	})
	G.Msg.Success = ""
	G.Msg.Fail = ""
	return
}

func SendGold(c *gin.Context) {
	user := M.GetAuthUser(c, G.FStore)
	if user.RoleID < 2 {
		return
	}
	var gt Mod.GoldTransfer
	var success bool
	email := c.PostForm("email")
	var receiver Mod.User
	var amount, senderFeesPercent, senderFeesFixed, totalSenderFees, receiverFeesPercent, receiverFeesFixed, totalReceiverFees, totalFees float64

	amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Invalid amount."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}
	amount = math.Round(amount*1000) / 1000
	if amount < 0.1 {
		G.Msg.Fail = "Amount can't be less than 0.1g."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	delivery, err := strconv.Atoi(c.PostForm("delivery_type"))
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Invalid sending type."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return

	} else if delivery == 2 {
		success, gt.ShipAddress = S.AddressValidation(c)
		if !success {
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}
		receiver = user

		if amount > user.Wallet.GoldAmount {
			G.Msg.Fail = "You don't have enough balance including fees."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}

	} else if delivery == 1 {
		receiver.Email = email
		receiver, success = R.ReadUserWithEmail(receiver)
		if !success {
			G.Msg.Fail = "No user found with the email."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}

		if receiver.ID == user.ID {
			G.Msg.Fail = "You can't send to your own wallet."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}

		senderFeesPercent, err = strconv.ParseFloat(G.Adm["Transfer_Sender_Fees_In_Percent"].Value.String, 64)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Invalid fees, tell the admin to set valid fees."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}
		senderFeesFixed, err = strconv.ParseFloat(G.Adm["Transfer_Sender_Fees_Fixed_In_Milligram"].Value.String, 64)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Invalid fees, tell the admin to set valid fees."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}

		totalSenderFees = (amount*senderFeesPercent)/100.0
		totalSenderFees += senderFeesFixed / 1000
		totalSenderFees = math.Ceil(totalSenderFees * 1000) / 1000

		if amount+totalSenderFees > user.Wallet.GoldAmount {
			G.Msg.Fail = "You don't have enough balance including fees."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}

		receiverFeesPercent, err = strconv.ParseFloat(G.Adm["Transfer_Receiver_Fees_In_Percent"].Value.String, 64)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Invalid fees, tell the admin to set valid fees."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}
		receiverFeesFixed, err = strconv.ParseFloat(G.Adm["Transfer_Receiver_Fees_Fixed_In_Milligram"].Value.String, 64)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Invalid fees, tell the authority to set valid fees."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}

		totalReceiverFees = (amount*receiverFeesPercent)/100.0
		totalReceiverFees += receiverFeesFixed / 1000
		totalReceiverFees = math.Ceil(totalReceiverFees * 1000) / 1000

		totalFees = totalSenderFees+totalReceiverFees
	} else {
		G.Msg.Fail = "Invalid sending type."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	trnsfrId, err := uuid.NewV4()
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Something went wrong."
		return
	}

	gt.TransferID = strings.ReplaceAll(trnsfrId.String(), "-", "")
	gt.ReceiverUser.Email = receiver.Email
	gt.SenderWalletID = user.Wallet.ID
	gt.ReceiverWalletID = receiver.Wallet.ID
	gt.GoldAmount = amount
	gt.SenderFeesPercent = senderFeesPercent
	gt.SenderFeesFixed = senderFeesFixed
	gt.SenderTotalFees = totalSenderFees
	gt.AmountToDeduct,_ = strconv.ParseFloat(fmt.Sprintf("%.3f",totalSenderFees + gt.GoldAmount), 64)
	gt.ReceiverFeesPercent = receiverFeesPercent
	gt.ReceiverFeesFixed = receiverFeesFixed
	gt.ReceiverTotalFees = totalReceiverFees
	gt.ReceiverAmount, _ = strconv.ParseFloat(fmt.Sprintf("%.3f",gt.GoldAmount - gt.ReceiverTotalFees), 64)
	gt.TotalFees, _ = strconv.ParseFloat(fmt.Sprintf("%.3f",totalFees), 64)
	gt.DeliveryType = delivery
	gt.EmailVerification = H.RandomString(60)
	gt.EmailVerifyCode = H.RandomNumber(10)

	if !R.AddGoldTransfer(gt) {
		G.Msg.Fail = "Something went wrong."
		return
	}

	success = S.SendGoldTransferEmail(gt, user, c, "")
	log.Println("Gold Transfer Email Confirmation Email Sending:", success)

	G.Msg.Success = "Please go to your email to confirm."
	c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
	return
}

func ConfirmGoldSendingEmail(c *gin.Context) {
	user := M.GetAuthUser(c, G.FStore)
	if user.RoleID > 1 {
		encEmail := c.Query("encEmail")
		emailVf := c.Query("emailVf")
		transferId := c.Query("id")
		var err error
		var decoded []byte
		decoded, err = base64.URLEncoding.DecodeString(encEmail)
		if err != nil {
			log.Println(err.Error())
			c.HTML(http.StatusOK, "404.html", nil)
			return
		}
		if user.Email != string(decoded) {
			c.HTML(http.StatusOK, "404.html", nil)
			return
		}

		var gt Mod.GoldTransfer
		gt.EmailVerification = emailVf
		gt.TransferID = transferId

		gt = R.GoldTransfer(gt, "status=0 and transfer_id=? and email_verification=? and sender_wallet_id=?", gt.TransferID, gt.EmailVerification, user.Wallet.ID)
		if gt.ID > 0 {
			user.GoldTransfers = append(user.GoldTransfers, gt)

			var wbsts []Mod.WebsiteSetting
			wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

			var menus []Mod.Menu
			menus = R.Menus(menus, "status=?", 1)

			session, _ := G.Store.Get(c.Request, "cart")
			lenCart := len(session.Values)

			G.LangVal = H.GetCookie("secret", nil, "lang", c)
			if G.LangVal == "" {
				G.LangVal = "en"
				H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
			}
			var lang G.Lang
			lang.LangValue = G.LangVal

			c.HTML(http.StatusOK, "confirm-gold-transfer.html", map[string]interface{}{
				"AppEnv": G.AppEnv, "User": user, "Banner": "Account", "Carts": lenCart, "Lang": lang,
				"Msg": G.Msg, "Adm": G.Adm, "Title": "Account", "Active": "Account",
				"Wbsts": wbsts, "Menus": menus, "NavActive": "Wallet", "EncEmail": encEmail, "EmailVf": emailVf,
			})
			G.Msg.Success = ""
			G.Msg.Fail = ""
			return

		} else {
			c.HTML(http.StatusOK, "404.html", nil)
		}
	} else {
		c.Redirect(http.StatusFound, "/login?val=account/wallet")
	}

}

func UpdateGoldSend(c *gin.Context) {
	user := M.GetAuthUser(c, G.FStore)
	if user.RoleID > 1 {
		encEmail := c.Query("encEmail")
		emailVf := c.Query("emailVf")
		transferId := c.Query("id")
		status := c.Query("status")
		var err error
		var decoded []byte
		decoded, err = base64.URLEncoding.DecodeString(encEmail)
		if err != nil {
			log.Println(err.Error())
			return
		}
		if user.Email != string(decoded) {
			return
		}

		var gt Mod.GoldTransfer
		gt.EmailVerification = emailVf
		gt.TransferID = transferId

		gt = R.GoldTransfer(gt, "status=0 and transfer_id=? and email_verification=? and sender_wallet_id=?", gt.TransferID, gt.EmailVerification, user.Wallet.ID)
		if gt.AmountToDeduct > user.Wallet.GoldAmount {
			return
		}
		if gt.ID > 0 {
			if status == "cancel" {
				gt.Status = 3
				gt.EmailVerification = ""
				gt.EmailVerifyCode = ""
				if !R.SaveGoldTransfer(gt) {
					G.Msg.Fail = "Something went wrong."
				} else {
					G.Msg.Success = "Cancelled successfully."
				}

			} else if status == "confirm" {
				if !S.ProcessGoldTransfer(gt, user, c) {
					if G.Msg.Fail == "" {
						G.Msg.Fail = "Something went wrong."
					}
				} else {
					if gt.DeliveryType == 2 {
						gt.Status = 2
						gt.EmailVerification = ""
						gt.EmailVerifyCode = ""
						if !R.SaveGoldTransfer(gt) {
							G.Msg.Fail = "Something went wrong."
						} else {
							G.Msg.Success = "Your transfer is under process."
						}
					} else if gt.DeliveryType == 1 {
						G.Msg.Success = "Your transfer is under process. Please reload after a few minutes."
					}
				}
			}

			c.Redirect(http.StatusFound, "account/wallet")
		} else {
			return
		}
	} else {
		return
	}
}
