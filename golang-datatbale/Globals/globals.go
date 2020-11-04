package Globals

import (
	Mod "gold-store/Models"
	"github.com/bykovme/gotrans"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"html/template"
)

type DB_ENV struct {
	Host, Port, Dialect, Username, Password, DBname string
}

type App_env struct {
	Name, Url, Port, Debug string
}

type RedsysEnv struct {
	Url string
}

type SessionCooikeEnv struct {
	Name string
}

type WalletAppEnv struct {
	Url, Key, SocketUrl string
}

type RabbitMQEnv struct {
	Url, ReqEx, ReqQ, ResEx, ResQ string
}

type Message struct {
	Success template.HTML
	Fail template.HTML
}

type EmailGenerals struct {
	From, To, Subject, HtmlString string
}

type UserDataForEmail struct {
	EncEmail string
	User Mod.User
	PS Mod.PasswordReset
	AppEnv App_env
}

type SocialAuthEnv struct {
	FacebookClientID, FacebookClientSecret, GoogleClientID, GoogleClientSecret string
}


type Cart struct {
	Key string
	Product Mod.Product
	Quantity float64
	Total float64
	TotalDiscount float64
	TotalWithDiscount float64
}

type FinalCart struct {
	Carts []Cart
	SubTotal float64
	TotalGrmAmount float64
	Fees1Percent float64
	Fees1Fixed float64
	Fees2Percent float64
	Fees2Fixed float64
	Fees3Percent float64
	Fees3Fixed float64
	TotalFees float64
	TotalDiscount float64
	GrandTotal float64
}

type Lang struct {
	LangValue string
}

func (l Lang) ConvertString(str string) string{
	str = gotrans.Tr(LangVal, str)
	return str
}

func (l Lang) ConvertHtml(html template.HTML) template.HTML{
	html = template.HTML(gotrans.Tr(LangVal, string(html)))
	return html
}

var(
	LangVal string
	Store = sessions.NewCookieStore([]byte("secret"))
	FStore = sessions.NewFilesystemStore("./Storage/Session",[]byte("secret"))
	DBEnv DB_ENV
	DB *gorm.DB
	Adm = make(map[string]Mod.AdminSetting)
	Msg Message
	AppEnv App_env
	SocialEnv SocialAuthEnv
	Order = make(map[uint]Mod.Order)
	Country = make(map[string]string)
)


