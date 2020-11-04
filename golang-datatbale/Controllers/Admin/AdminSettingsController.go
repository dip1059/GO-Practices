package Admin

import (
	"github.com/bykovme/gotrans"
	"github.com/gin-gonic/gin"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"html/template"
	"net/http"
	"strconv"
)


type Country struct {
	Key   string
	Value string
}

var (
	Bank = make(map[uint]Mod.Bank)
	PayMethod = make(map[uint]Mod.PayMethod)
)


func AdminSettings(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	data := c.Param("data")
	//primary
	var admSets []Mod.AdminSetting
	admSets = R.AllAdminSettings(admSets)
	for i, _ := range admSets {
		G.Adm[admSets[i].Slug] = admSets[i]
	}
	//

	//paymethod
	var payMethods []Mod.PayMethod
	payMethods = R.PayMethods(payMethods, "id <> ?", 7)
	for i, _ := range payMethods {
		PayMethod[payMethods[i].ID] = payMethods[i]
	}
	//

	//general
	user.CountryCode = G.Country[user.CountryCode]
	var countries []Country
	var country Country
	for key, value := range G.Country {
		country.Key = key
		country.Value = value
		countries = append(countries, country)
	}
	//

	//bank
	var banks []Mod.Bank
	banks = R.AllBank(banks)
	for i, _ := range banks {
		Bank[banks[i].ID] = banks[i]
	}
	//

	c.HTML(http.StatusOK, "admin-settings.html", map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user, "Nav":"settings", "Title":"Admin-Settings", "Msg":G.Msg,"Adm":G.Adm,
		"AdminSettings":admSets, "Countries":countries, "Data":data, "Banks":banks, "PayMethods":payMethods})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func EditPrimarySetting(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	slug := c.Param("slug")
	tab := c.Query("tab")
	admSet := G.Adm[slug]
	c.HTML(http.StatusOK, "edit-primary-setting.html", map[string]interface{}{"Tab":tab,
		"AppEnv":G.AppEnv,  "User":user, "Nav":"settings", "Title":"Edit-Admin-Setting", "Msg":G.Msg,"Adm":G.Adm, "AdminSetting":admSet})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func UpdatePrimarySetting(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.PostForm("id"))
	value := c.PostForm("value")
	//slug := c.PostForm("slug")
	tab := c.Query("tab")

	var admSet Mod.AdminSetting
	admSet.ID = uint(id)

	if R.UpdateAdminSetting(admSet, map[string]interface{}{"value": value}) {
		G.Msg.Success = "Updated Successfully"
		c.Redirect(http.StatusFound, "/admin-settings/"+tab)
	} else {
		G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
	}

}


func DeletePrimarySetting(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var admSet Mod.AdminSetting
	admSet.ID = uint(id)

	lang := H.GetCookie("secret", nil,"lang", c)
	if lang == "" {
		lang = "en"
		H.SetCookie("secret", nil, lang, "lang", 60*60*24*365, c)
	}

	if R.DeleteAdminSetting(admSet) {
		G.Msg.Success = template.HTML(gotrans.Tr(lang,"Successfully Deleted"))
		c.Redirect(http.StatusFound, "/admin-settings/prm")
	} else {
		G.Msg.Fail = template.HTML(gotrans.Tr(lang, "Some Error Occurred, Delete Failed. Please Try Again Later."))
		c.Redirect(http.StatusFound, "/admin-settings/prm")
	}
}


func AddBank(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var bank Mod.Bank
	bank.Name = c.PostForm("name")
	bank.Details = template.HTML(c.PostForm("details"))
	bank.Status, _ = strconv.Atoi(c.PostForm("status"))
	if R.AddBank(bank) {
		G.Msg.Success = "Bank Added Successfully."
		c.Redirect(http.StatusFound, "/admin-settings/bank")
		return
	} else {
		G.Msg.Fail = "Some Error Occured. Please Try Again Later."
		c.Redirect(http.StatusFound, "/admin-settings/bank")
		return
	}
}


func MakeBankInactive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var bank Mod.Bank
	id, _ := strconv.Atoi(c.Param("id"))
	bank = Bank[uint(id)]
	bank.Status = 0
	if R.UpdateBank(bank) {
		G.Msg.Success = "Status Updated Successfully"
		c.Redirect(http.StatusFound, "/admin-settings/bank")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/admin-settings/bank")
	}
}


func MakeBankActive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var bank Mod.Bank
	id, _ := strconv.Atoi(c.Param("id"))
	bank = Bank[uint(id)]
	bank.Status = 1
	if R.UpdateBank(bank) {
		G.Msg.Success = "Status Updated Successfully."
		c.Redirect(http.StatusFound, "/admin-settings/bank")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/admin-settings/bank")
	}
}


func EditBank(c *gin.Context) {
	var user Mod.User
	var success bool
	user, success = M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var bank Mod.Bank
	id, _ := strconv.Atoi(c.Param("id"))
	bank = Bank[uint(id)]
	
	c.HTML(http.StatusOK, "edit-bank.html",map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user,  "Nav":"settings", "Title":"Edit-Bank", "Bank":bank, "Msg":G.Msg})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func UpdateBank(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var bank Mod.Bank
	id, _ := strconv.Atoi(c.PostForm("id"))
	bank = Bank[uint(id)]
	bank.Name = c.PostForm("name")
	bank.Details = template.HTML(c.PostForm("details"))
	
	if R.UpdateBank(bank) {
		G.Msg.Success = "Updated Successfully"
		c.Redirect(http.StatusFound, "/admin-settings/bank")
	} else {
		G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/admin-settings/bank")
	}
}


func DeleteBank(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var bank Mod.Bank
	id, _ := strconv.Atoi(c.Param("id"))
	bank = Bank[uint(id)]
	if R.DeleteBank(bank) {
		G.Msg.Success = "Deleted Successfully"
		c.Redirect(http.StatusFound, "/admin-settings/bank")
	} else {
		G.Msg.Fail = "Some Error Occurred, Deletion Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/admin-settings/bank")
	}
}


func MakePayMethodInactive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var payMethod Mod.PayMethod
	id, _ := strconv.Atoi(c.Param("id"))
	payMethod = PayMethod[uint(id)]
	payMethod.Status = 0
	if R.UpdatePayMethod(payMethod) {
		G.Msg.Success = "Status Updated Successfully"
		c.Redirect(http.StatusFound, "/admin-settings/pmt")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/admin-settings/pmt")
	}
}


func MakePayMethodActive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var payMethod Mod.PayMethod
	id, _ := strconv.Atoi(c.Param("id"))
	payMethod = PayMethod[uint(id)]
	payMethod.Status = 1
	if R.UpdatePayMethod(payMethod) {
		G.Msg.Success = "Status Updated Successfully."
		c.Redirect(http.StatusFound, "/admin-settings/pmt")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/admin-settings/pmt")
	}
}


func EditPayMethod(c *gin.Context) {
	var user Mod.User
	var success bool
	user, success = M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var payMethod Mod.PayMethod
	id, _ := strconv.Atoi(c.Param("id"))
	payMethod = PayMethod[uint(id)]

	c.HTML(http.StatusOK, "edit-payMethod.html",map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user,
		"Nav":"settings", "Title":"Edit-PayMethod", "PayMethod":payMethod, "Msg":G.Msg})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func UpdatePayMethod(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var payMethod Mod.PayMethod
	id, _ := strconv.Atoi(c.PostForm("id"))
	payMethod = PayMethod[uint(id)]
	payMethod.Method = c.PostForm("method")
	payMethod.Description.String = c.PostForm("description")
	payMethod.Description = H.NullStringProcess(payMethod.Description)

	if R.UpdatePayMethod(payMethod) {
		G.Msg.Success = "Updated Successfully"
		c.Redirect(http.StatusFound, "/admin-settings/pmt")
	} else {
		G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/admin-settings/pmt")
	}
}


func DeletePayMethod(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var payMethod Mod.PayMethod
	id, _ := strconv.Atoi(c.Param("id"))
	payMethod = PayMethod[uint(id)]
	if R.DeletePayMethod(payMethod) {
		G.Msg.Success = "Deleted Successfully"
		c.Redirect(http.StatusFound, "/admin-settings/pmt")
	} else {
		G.Msg.Fail = "Some Error Occurred, Deletion Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/admin-settings/pmt")
	}
}