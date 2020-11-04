package Services

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	"html/template"
	"log"
	"net"
	"strconv"
	"time"
)

func StripeCreateCharge(c *gin.Context, amount int, orderID uint) (*stripe.Charge, bool) {
	stripe.Key = G.Adm["Stripe_Secret_Key"].Value.String
	stripeToken := c.PostForm("stripeToken")
	stripeEmail := c.PostForm("stripeEmail")
	randomString := H.RandomString(60)
	/*formAmount, _ := strconv.Atoi(c.PostForm("amount"))
	log.Println("Stripe Amount from checkout form:", formAmount, "----- Cart Current Total Amount:", amount)
	if formAmount != amount {
		return nil, false
	}*/

	for {
		chargeParams := &stripe.ChargeParams {
			Amount:      stripe.Int64(int64(amount)),
			Currency:    stripe.String(string(stripe.CurrencyEUR)),
			Description: stripe.String("Payment by " + stripeEmail+ " on Order: #"+strconv.Itoa(int(orderID))),
		}
		err := chargeParams.SetSource(stripeToken)
		if err != nil {
			log.Println(err.Error())
		}
		chargeParams.SetIdempotencyKey(randomString)
		ch, err := charge.New(chargeParams)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				time.Sleep(time.Second)
				continue
			}
			var stripeErr stripe.Error
			decErr := json.Unmarshal([]byte(err.Error()), &stripeErr)
			if decErr != nil {
				log.Println(decErr.Error())
			}
			G.Msg.Fail = template.HTML("Stripe error message: "+stripeErr.Msg)
			log.Println(err.Error())
			return ch, false
		} else {
			return ch, true
		}
	}
}

func StripeGetCharge(chID string) (*stripe.Charge, bool){
	stripe.Key = G.Adm["Stripe_Secret_Key"].Value.String
	ch, err := charge.Get(chID, nil)
	if err !=nil {
		log.Println("StripeService log2", err.Error())
		return ch, false
	} else {
		return ch, true
	}
}

