package Services

import (
	"github.com/gin-gonic/gin"
	"github.com/rikvdh/go-mollie"
	G "gold-store/Globals"
	"log"
	"strconv"
)


func GetMethods() {
	m := mollie.Get(G.Adm["Mollie_Api_Key"].Value.String)
	methods, err := m.Methods().List()
	if err != nil {
		panic(err)
	}

	for _, method := range methods {
		log.Printf("method %s: %s\n", method.ID, method.Description)
	}
}

func MolliePayment(amount float64, userEmail string, OrderID uint, c *gin.Context) (*mollie.PaymentReply, bool){

	m := mollie.Get(G.Adm["Mollie_Api_Key"].Value.String)

	payment := mollie.PaymentData {
		Amount:      amount,
		Description: "Payment by user "+userEmail,
		RedirectURL: G.AppEnv.Url+"/mollie/redirect?order_id="+strconv.Itoa(int(OrderID)),
		WebhookURL:  G.AppEnv.Url+"/mollie/webhook",
		/*Method:      "ideal",*/
		Metadata: map[string]string {
			"order_id": strconv.Itoa(10000 + int(OrderID)),
		},
	}

	response, err := m.Payments().New(payment)
	if err != nil {
		log.Println(err.Error())
		return nil, false
	}

	//c.Redirect(http.StatusFound, response.Links["paymentUrl"])
	return response, true
}


func GetPayment(id string) *mollie.PaymentReply{
	m := mollie.Get(G.Adm["Mollie_Api_Key"].Value.String)
	response, err := m.Payments().Get(id)
	if err != nil {
		log.Println(err.Error())
	}
	return response
}
