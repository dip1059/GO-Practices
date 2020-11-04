package Services

import (
	"github.com/gin-gonic/gin"
	//"github.com/nu7hatch/gouuid"
	Cfg "gold-store/Config"
	"strconv"
	"time"
	//G "gold-store/Globals"
	//H "gold-store/Helpers"
	//M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"log"
)

//chunk function
func GenerateCertificate(order Mod.Order, c *gin.Context) {
	start := time.Now()
	db := Cfg.DBConnect()
	db = db.Begin()
	log.Println("Certificate generation started for order_id:", order.ID, "user_id:", order.UserID, "total_grm_amount:", order.TotalGrmAmount)
	log.Println("Total rows will be:", int(order.TotalGrmAmount*1000))

	count := 0
	insertQuery := "INSERT INTO `gold_certificates` (`id`, `user_id`, `order_id`, `delivery_type`) "
	valuesQuery := "VALUES "
	userIdStr := strconv.Itoa(int(order.UserID))
	orderIdStr := strconv.Itoa(int(order.ID))
	deliveryStr := strconv.Itoa(int(order.DeliveryType))
	loopEnd := int(order.TotalGrmAmount * 1000)
	chunkCount := 0

	for i := 1; i <= loopEnd; i++ {
		count++
		//id, err := uuid.NewV4()
		//if err != nil {
		//	log.Println(err.Error())
		//	log.Println("UUID generation failed for order_id:",order.ID,"user_id",order.UserID)
		//	log.Println("Failed on row number:", count)
		//	log.Println(time.Now().Sub(start).Seconds())
		//	return
		//}
		//createdAt := time.Now().Format("2006-01-02 15:04:05")

		valuesQuery += "(uuid(), " + userIdStr + ", " + orderIdStr + ", " + deliveryStr + ")"
		if i%10000 == 0 || i+1 > loopEnd {
			chunkCount++
			valuesQuery += ";"
			finalQuery := insertQuery + valuesQuery

			if !R.ExecuteCertificate(db, finalQuery) {
				log.Println("Certificate generation failed for order_id:", order.ID, "user_id:", order.UserID)
				log.Println("Failed on row number:", count, "chunk number", chunkCount)
				log.Println(time.Now().Sub(start).Seconds())
				return
			}
			valuesQuery = "VALUES "
		} else {
			valuesQuery += ", "
		}
	}

	//log.Println(finalQuery)
	db.Commit()
	defer db.Close()
	log.Println("Certificate generated successfully for order_id:", order.ID, "user_id:", order.UserID)
	log.Println("Total chunk:", chunkCount, "Total rows:", count)
	log.Println(time.Now().Sub(start).Seconds())

	if order.DeliveryType == 1 {
		note := "from order " + orderIdStr
		UpdateWalletGoldAmount(nil, order.TotalGrmAmount, order.UserID, note)
	}

	success := SendOrderEmail(order, c, "")
	log.Println("Order Email Sending Success:", success)
}


func TransferCertificate(gt Mod.GoldTransfer, user Mod.User, c *gin.Context) {
	start := time.Now()
	db := Cfg.DBConnect()
	db = db.Begin()
	log.Println("Certificate transfer started for transfer_id:", gt.TransferID)

	var updateQuery string
	if gt.DeliveryType == 1 {
		gt.Status = 1

		var pt Mod.ProcessedTransfer
		pt.TransferID = gt.TransferID
		if !R.AddProcessedTransfer(db, pt) {
			defer db.Close()
			log.Println("Certificate transfer failed for transfer_id:", gt.TransferID)
			log.Println(time.Now().Sub(start).Seconds())
			return
		}

		adminPart := `UPDATE gold_certificates SET user_id=1 
		where user_id=` + strconv.Itoa(int(gt.SenderUser.ID)) + ` and delivery_type=1 order by created_at asc limit
        `+ strconv.Itoa(int(gt.TotalFees * 1000)) + `;`

		if !R.ExecuteCertificate(db, adminPart) {
			defer db.Close()
			log.Println("Certificate transfer failed for transfer_id:", gt.TransferID)
			log.Println(time.Now().Sub(start).Seconds())
			return
		}
		log.Println("Fees Certificate transferred to admin for transfer_id:", gt.TransferID, "total:", gt.TotalFees * 1000)

		receiverPart := `UPDATE gold_certificates SET user_id=`+strconv.Itoa(int(gt.ReceiverUser.ID))+
		` where user_id=` + strconv.Itoa(int(gt.SenderUser.ID)) + ` and delivery_type=1 order by created_at asc limit
		`+ strconv.Itoa(int(gt.ReceiverAmount * 1000)) + `;`

		if !R.ExecuteCertificate(db, receiverPart) {
			db.Rollback()
			defer db.Close()
			log.Println("Certificate transfer failed for transfer_id:", gt.TransferID)
			log.Println(time.Now().Sub(start).Seconds())
			return
		}
		log.Println("Certificate transferred to user:", gt.ReceiverUser.ID,"for transfer_id:", gt.TransferID, "total:", gt.ReceiverAmount * 1000)

	} else if gt.DeliveryType == 2 {
		updateQuery = `UPDATE gold_certificates SET delivery_type=2 where user_id=` + strconv.Itoa(int(gt.SenderUser.ID)) +
			` and delivery_type=1 order by created_at asc limit `+ strconv.Itoa(int(gt.ReceiverAmount * 1000)) + `;`

		if !R.ExecuteCertificate(db, updateQuery) {
			db.Rollback()
			defer db.Close()
			log.Println("Certificate transfer failed for transfer_id:", gt.TransferID)
			log.Println(time.Now().Sub(start).Seconds())
			return
		}
		log.Println("Certificate delivery type updated for user:", gt.SenderUser.ID,"for transfer_id:", gt.TransferID, "total:", gt.AmountToDeduct * 1000)
	}

	note := "from transfer " + gt.TransferID
	if gt.DeliveryType == 1 {
		ok := UpdateWalletGoldAmount(db, -gt.AmountToDeduct, gt.SenderUser.ID, note)
		if !ok {
			db.Rollback()
			defer db.Close()
			log.Println("Certificate transfer failed for transfer_id:", gt.TransferID)
			return
		}
		ok = UpdateWalletGoldAmount(db, gt.TotalFees, 1, note)
		if !ok {
			db.Rollback()
			defer db.Close()
			log.Println("Certificate transfer failed for transfer_id:", gt.TransferID)
			return
		}

		ok = UpdateWalletGoldAmount(db, gt.ReceiverAmount, gt.ReceiverUser.ID, note)
		if !ok {
			db.Rollback()
			defer db.Close()
			log.Println("Certificate transfer failed for transfer_id:", gt.TransferID)
			return
		}
	} else if gt.DeliveryType == 2 {
		ok := UpdateWalletGoldAmount(nil, -gt.AmountToDeduct,gt.SenderUser.ID, note)
		if !ok {
			db.Rollback()
			defer db.Close()
			log.Println("Certificate transfer failed for transfer_id:", gt.TransferID)
			return
		}
	}

	if gt.DeliveryType == 1 {
		gt.Status = 1
		gt.EmailVerification = ""
		gt.EmailVerifyCode = ""
		if !R.SaveGoldTransfer(gt) {
			db.Rollback()
			defer db.Close()
			log.Println("Certificate transfer failed for transfer_id:", gt.TransferID)
			return
		}
	}

	db.Commit()
	defer db.Close()

	gt = R.GoldTransfer(gt)
	success := SendGoldTransferEmail(gt, user, c, "")
	log.Println("Gold Transfer Email Sending:", success)

	log.Println("Certificate transfer succeeded for transfer_id:", gt.TransferID)
	log.Println(time.Now().Sub(start).Seconds())
	return
}

//all in once
/*func GenerateCertificate(order Mod.Order) {
	start := time.Now()
	log.Println("Certificate generation started for order_id:",order.ID,"user_id:",order.UserID, "total_grm_amount:", order.TotalGrmAmount)

	count := 0
	insertQuery := "INSERT INTO `gold_certificates` (`id`, `user_id`, `order_id`, `delivery_type`) "
	valuesQuery := "VALUES "
	userIdStr := strconv.Itoa(int(order.UserID))
	orderIdStr := strconv.Itoa(int(order.ID))

	for i:=1.0; i<=order.TotalGrmAmount*1000; i++ {
		count++
		//id, err := uuid.NewV4()
		//if err != nil {
		//	log.Println(err.Error())
		//	log.Println("UUID generation failed for order_id:",order.ID,"user_id",order.UserID)
		//	log.Println("Failed on row number:", count)
		//	log.Println(time.Now().Sub(start).Seconds())
		//	return
		//}

		//createdAt := time.Now().Format("2006-01-02 15:04:05")
		//valuesQuery += "('"+ createdAt+"', '"+createdAt+"', '"+id.String()+"', "+userIdStr+", "+orderIdStr+")"
		valuesQuery += "(uuid(), "+ userIdStr+", "+orderIdStr+", "+deliveryStr+")"
		if i+1 > order.TotalGrmAmount*10 {
			valuesQuery += ";"
		} else {
			valuesQuery += ", "
		}
	}
	finalQuery := insertQuery+valuesQuery
	//log.Println(finalQuery)
	//log.Println("Total rows:", count)
	//log.Println(time.Now().Sub(start).Seconds())
	//return

	if !R.GenerateCertificate(finalQuery) {
		log.Println("Certificate generation failed for order_id:",order.ID,"user_id:",order.UserID)
		//log.Println("Failed on row number:", count)
		log.Println(time.Now().Sub(start).Seconds())
		return
	}

	log.Println("Certificate generated successfully for order_id:",order.ID,"user_id:",order.UserID)
	log.Println("Total rows:", count)
	log.Println(time.Now().Sub(start).Seconds())
}*/
