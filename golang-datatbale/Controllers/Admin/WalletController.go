package Admin

import (
	"github.com/gin-gonic/gin"
	G "gold-store/Globals"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	S "gold-store/Services"
	"log"
	"net/http"
	"strconv"
)

func Wallets(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var wallets []Mod.Wallet
	wallets = R.OnlyWallets(wallets)
	c.HTML(http.StatusOK, "wallets.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user,  "Nav":"wallets", "Title": "Wallets", "Msg": G.Msg, "Wallets": wallets})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func WalletHistories(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var wallet Mod.Wallet
	id, _ := strconv.Atoi(c.Param("id"))

	wallet.ID = uint(id)
	wallet = R.WalletWithOthers(wallet)

	c.HTML(http.StatusOK, "wallet-histories.html",map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user,  "Nav":"wallets", "Title":"Wallet-History", "Wallet":wallet, "Msg":G.Msg,
		})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func GoldTransfers(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var gts []Mod.GoldTransfer
	gts = R.GoldTransfers(gts)
	c.HTML(http.StatusOK, "gold-transfers.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user,  "Nav":"wallets", "Title": "Gold-Transfers", "Msg": G.Msg, "Transfers": gts})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func TransferShippingDetails(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var gt Mod.GoldTransfer
	gt.ID = uint(id)
	gt = R.GoldTransfer(gt)
	c.HTML(http.StatusOK, "transfer-shipping-details.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user,  "Nav":"wallets", "Title": "Gold-Transfers", "Msg": G.Msg, "Transfer": gt})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func MakeTransferCancel(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var gt Mod.GoldTransfer
	gt.ID = uint(id)
	if R.UpdateGoldTransfer(gt, map[string]interface{}{"status": 3}) {
		G.Msg.Success = "Cancelled Successfully"
	} else {
		G.Msg.Fail = "Some Error Occurred,Order Status Update Failed. Please Try Again Later."
	}
	c.Redirect(http.StatusFound, "/gold-transfers")
}

func AddTransferTrackCode(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.PostForm("id"))
	code := c.PostForm("track_code")
	var gt Mod.GoldTransfer
	gt.ID = uint(id)
	if R.UpdateGoldTransfer(gt, map[string]interface{}{"track_code": code, "status": 4}) {
		G.Msg.Success = "Track Code Added Successfully"
		gt = R.GoldTransfer(gt)
		success := S.SendGoldTransferEmail(gt, gt.SenderUser, c, "")
		log.Println("Transfer Email Sending Success:", success)
	} else {
		G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
	}
	c.Redirect(http.StatusFound, "/gold-transfers")
}
