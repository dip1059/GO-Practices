package Services

import (
	"github.com/gin-gonic/gin"
)

type LusopayData struct {
	Entity string		`json:"entity, omitempty"`
	Value float64	`json:"value, omitempty"`
	Reference string	`json:"reference, omitempty"`
	Table string	`json:"table, omitempty"`
	Status bool	`json:"status"`
	Message string	`json:"message, omitempty"`
}


func LusopayProcess(c *gin.Context) {
	/*_, success :=M.IsAuthUser(c,G.FStore)
	if !success {
		c.Redirect(http.StatusFound, "/checkout")
		return
	}
	ord := H.RandomNumber(12)
	finalCart := ProcessCart(c)*/

	/*sc := securecookie.New([]byte(G.WalletApp.Key), nil)
	encOrd, err := sc.Encode("encryptedOrder", ord)
	if err != nil {
		log.Println(err.Error())
	}*/

	/*cookie, err := c.Cookie(G.SessionCookie.Name)
	if err != nil {
		log.Println(err.Error())
	}
	sessionId, err := golaravelsession.GetSessionID(cookie, G.WalletApp.Key)
	if err != nil {
		log.Println(err.Error())
	}

	req2, _ := json.Marshal(map[string]interface{}{
		"order_id": ord,
		"order_value": finalCart.GrandTotal,
		"secret_key": sessionId,
	})
	log.Println(string(req2))

	url := G.WalletApp.Url+"/api/get-goldstore-lusopay-payment-data"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req2))
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(url)

	Bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(Bytes))

	var d LusopayData
	err = json.Unmarshal(Bytes, &d)
	if err != nil {
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, d)

	defer resp.Body.Close()*/
}
