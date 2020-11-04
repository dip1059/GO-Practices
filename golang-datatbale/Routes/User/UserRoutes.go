package User

import (
	"github.com/gin-gonic/gin"
	"gold-store/Controllers/User"
	S "gold-store/Services"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/", User.Home)
	r.POST("/get-custom-product-ajax", User.GetCustomProductAjax)
	r.GET("/product-details/:id", User.ProductDetails)
	r.GET("/add-cart/:id/:quantity/:data", User.AddCart)
	r.GET("/add-all-to-cart", User.AddAllToCart)
	r.POST("/add-cart", User.AddCartCustom)
	r.GET("/show-cart", User.ShowCart)
	r.GET("/delete-from-cart/:key", User.DeleteFromCart)
	r.GET("/checkout", User.CheckoutGet)
	r.POST("/checkout", User.CheckoutPost)
	r.GET("/about-us", User.AboutPage)
	r.GET("/contact-us", User.ContactPage)
	r.GET("/alerts", User.AlertasPage)
	r.GET("/privacy-policy", User.PrivacidadePage)
	r.GET("/cookies-policies", User.CookiesPage)
	r.GET("/terms-of-use", User.TermosPage)
	r.GET("/page/:url", User.DynamicPage)
	r.GET("/account/orders", User.Orders)
	r.GET("/account/wallet", User.MyWallet)
	r.GET("/account/address", User.MyAddress)
	r.GET("/account/add-address", User.AddAddress)
	r.GET("/account/edit-address", User.EditAddress)
	r.POST("/account/update-address", User.UpdateShippingAddress)
	r.GET("/account/set-as-default-address", User.SetAsDefaultShippingAddress)
	r.GET("/account/unset-as-default-address", User.UnsetAsDefaultShippingAddress)
	r.GET("/account/delete-address", User.DeleteShippingAddress)
	r.GET("/account/wishlist", User.MyWishlist)
	//r.GET("/user-order-details/:id", User.UserOrderDetails)
	r.GET("/add-wish/:proId", User.AddWish)
	r.GET("/remove-from-wishlist/:id/:proId", User.RemoveFromWishlist)
	r.POST("/update-profile-pic", User.UpdateProfilePic)
	r.POST("/update-user/:data", User.UpdateUser)
	r.POST("/update-social-user", User.UpdateSocialUser)
	r.POST("/add-order-doc", User.AddOrderDocPost)
	r.POST("/subs-newsletter", User.SubsNewsletter)
	r.POST("/contact-us", User.ContactMessage)
	r.POST("/change-lang", User.ChangeLang)
	//r.GET("/download-checkout-pdf/:bankId/:ref", User.DownloadCheckoutPDF)
	r.GET("/download-order-pdf/:id", User.DownloadOrderPDF)
	r.GET("/download-order-invoice/:id", User.DownloadOrderInvoice)

	/*r.GET("/redsys-process", User.RedsysCheckoutProcess)
	r.GET("/redsys-success-notification", User.RedsysCheckoutSuccess)
	r.GET("/redsys-error-notification", User.RedsysCheckoutFail)

	r.GET("/lusopay-process", User.LusopayCheckoutProcess)*/

	r.GET("/buy-now/:id", User.BuyNow)

	r.POST("/check-cash-point", User.CheckCashPoint)

	r.POST("/get-discount-coupon", User.GetDiscountCoupon)

	/*r.GET("/free-bids-checkout", User.FreeBidsCheckoutGet)
	r.POST("/free-bids-checkout", User.FreeBidsCheckoutPost)*/

	r.GET("/mollie/redirect", S.PaymentRedirect)
	r.POST("/mollie/webhook", S.PaymentWebhook)

	r.GET("/transfer-details/:id", User.TransfersDetails)
	r.POST("/send-gold", User.SendGold)
	r.GET("/send-gold-email-verify", User.ConfirmGoldSendingEmail)
	r.GET("/update-gold-send", User.UpdateGoldSend)
}
