package Services

import (
	G "gold-store/Globals"
	"github.com/braintree-go/braintree-go"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func CreateTransaction(c *gin.Context) (*braintree.Transaction, bool){
	//ctx := context.Background()
	unscaled, _ := strconv.Atoi(c.PostForm("amount"))
	bt := braintree.New(
		braintree.Sandbox,
		G.Adm["Brain_Tree_Merchant_ID"].Value.String,
		G.Adm["Brain_Tree_Public_Key"].Value.String,
		G.Adm["Brain_Tree_Private_Key"].Value.String)
	//bt.ClientToken().Generate(c)
	tx := &braintree.TransactionRequest{
		Type:   "sale",
		Amount: braintree.NewDecimal(int64(unscaled), 2),
		PaymentMethodNonce:c.PostForm("payment_method_nonce"),
		Options:&braintree.TransactionOptions{
			SubmitForSettlement:true,
		},
	}

	trx, err := bt.Transaction().Create(c, tx)

	if err == nil {
		return trx, true
	} else {
		log.Println("BrainTreeService.go log1: ", err.Error())
		return trx, false
	}
}
