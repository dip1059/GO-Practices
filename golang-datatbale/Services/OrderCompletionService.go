package Services

import (
	G "gold-store/Globals"
	H "gold-store/Helpers"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"strconv"
	"time"
)


func SendOrderEmail(order Mod.Order, c *gin.Context, filename string) bool{

	orderID := 10000+order.ID
	order.PayMethod.ID = order.PayMethodID
	order.PayMethod = R.PayMethod(order.PayMethod)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	htmlString, err := H.ParseTemplate("Views/Email/order-email.html", map[string]interface{}{
		"User":order.User, "Wbsts":wbsts, "Order":order, "Adm":G.Adm,"OrderID":orderID, "Lang":lang, "AppEnv":G.AppEnv})

	if err != nil {
		log.Println(err.Error())
		return false
	}

	To := []string{order.User.Email}
	Subject := "Order Email"
	HtmlString := htmlString
	if !SendEmail(To, Subject, HtmlString, filename) {
		return false
	}
	return true
}


func GenerateInvoice(order Mod.Order) (string, bool){
	//order = R.Order(order)
	var invoice Mod.Invoice
	invoice.OrderID = order.ID
	invoice = R.AddInvoice(invoice)

	//order.User = R.GetUserAddress(order.User)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	var lang G.Lang
	lang.LangValue = G.LangVal

	address := template.HTML(G.Adm["Invoice_Company_Address"].Value.String)

	pdf := NewRequestPdf("")
	err := pdf.ParseTemplate("Views/PdfTemp/invoice.html", map[string]interface{}{
		"Wbsts":wbsts, "Adm":G.Adm, "Lang":lang, "Address":address, "AppEnv":G.AppEnv,"Order":order})
	if err != nil {
		log.Println(err.Error())
		return "", false
	}

	t := time.Now().Unix()
	fileName := strconv.FormatInt(int64(t), 10)
	success := pdf.GeneratePDF("Storage/Temp/"+fileName+".pdf")
	if success {
		log.Println("PDF Generated Successfully.")
		//c.Redirect(http.StatusFound, "/assets/Storage/Temp/"+fileName+".pdf")
	}
	return "Storage/Temp/"+fileName+".pdf", true
}