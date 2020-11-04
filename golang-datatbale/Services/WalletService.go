package Services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	Mod "gold-store/Models"
	"log"
	"strconv"

	//Cfg "gold-store/Config"
	//G "gold-store/Globals"
	R "gold-store/Repositories"
)


func UpdateWalletGoldAmount(db *gorm.DB,amountToAdd float64, userID uint, note string) bool {
	var walletHistory Mod.WalletHistory
	var wallet Mod.Wallet
	wallet = R.OnlyWallet(wallet, "user_id=?", userID)

	if wallet.ID > 0 {
		walletHistory.GoldAmountBefore = wallet.GoldAmount
		walletHistory.AddedGoldAmount,_ = strconv.ParseFloat(fmt.Sprintf("%.3f",amountToAdd), 64)
		walletHistory.GoldAmountAfter,_ = strconv.ParseFloat(fmt.Sprintf("%.3f",wallet.GoldAmount+amountToAdd), 64)
		//walletHistory.GoldAmountAfter = math.Floor(walletHistory.GoldAmountAfter * 1000) / 1000
		walletHistory.Note = note

		wallet.GoldAmount = walletHistory.GoldAmountAfter
		if R.SaveWallet(db, wallet) {
			log.Println("Wallet:", wallet.ID, "of User:",userID ,"updated successfully with amount:", amountToAdd)

			walletHistory.WalletID = wallet.ID
			ok := R.AddWalletHistory(db, walletHistory)
			log.Println("Wallet History Created:", ok, note)
			return true
		} else {
			log.Println("Wallet:", wallet.ID, "of User:",userID ,"update failed with amount:", amountToAdd)
			return false
		}

	} else if wallet.ID == 0 {
		wallet.UserID = userID
		wallet.GoldAmount,_ = strconv.ParseFloat(fmt.Sprintf("%.3f",amountToAdd), 64)

		wallet = R.AddWallet(db, wallet)
		if wallet.ID > 0 {
			log.Println("New wallet created for user:", userID)
			log.Println("Wallet:", wallet.ID, "of User:",userID ,"updated successfully with amount:", amountToAdd)

			walletHistory.WalletID = wallet.ID
			walletHistory.GoldAmountBefore = 0
			walletHistory.AddedGoldAmount,_ = strconv.ParseFloat(fmt.Sprintf("%.3f",amountToAdd), 64)
			walletHistory.GoldAmountAfter = wallet.GoldAmount
			walletHistory.Note = note

			ok := R.AddWalletHistory(db, walletHistory)
			log.Println("Wallet History Created:", ok, note)

			return true
		} else {
			log.Println("Wallet:", wallet.ID, "of User:",userID ,"update failed with amount:", amountToAdd)
			return false
		}
	}
	log.Println("Weird Wallet id:", wallet.ID)
	return false
}
