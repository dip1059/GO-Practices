package Admin

import (
	"gold-store/Controllers/Admin"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	r.GET("/dashboard", Admin.Dashboard)

	r.GET("/add-karat", Admin.AddKaratGet)
	r.POST("/add-karat", Admin.AddKaratPost)
	r.GET("/all-karat", Admin.AllKarat)
	r.GET("/update-status-karat/:id/:status", Admin.UpdateStatusKarat)
	r.GET("/edit-karat/:id", Admin.EditKarat)
	r.POST("/update-karat", Admin.UpdateKarat)
	r.GET("/delete-karat/:id", Admin.DeleteKarat)

	r.GET("/add-product", Admin.AddProductGet)
	r.POST("/add-product", Admin.AddProductPost)
	r.GET("/all-product", Admin.AllProduct)
	r.GET("/make-product-inactive/:id", Admin.MakeProductInactive)
	r.GET("/make-product-active/:id", Admin.MakeProductActive)
	r.GET("/edit-product/:id", Admin.EditProduct)
	r.POST("/update-product", Admin.UpdateProduct)
	//r.GET("/delete-product/:id", Admin.DeleteProduct)

	r.GET("/orders", Admin.Orders)
	//r.GET("/make-payment-pending/:id/:userId", Admin.MakePaymentPending)
	//r.GET("/make-payment-done/:id/:userId", Admin.MakePaymentDone)
	r.GET("/make-order-cancel/:id/:userId", Admin.MakeOrderCancel)
	r.GET("/make-order-complete/:id/:userId", Admin.MakeOrderComplete)
	r.POST("/add-track-code/:userId", Admin.AddTrackCode)
	r.GET("/order-details/:id/:userId", Admin.OrderDetails)
	r.GET("/delete-order/:id/:userId", Admin.DeleteOrder)

	r.GET("/users", Admin.Users)
	r.GET("/make-user-suspend/:id", Admin.MakeUserSuspend)
	r.GET("/make-user-active/:id", Admin.MakeUserActive)
	r.GET("/user-orders/:id", Admin.UserOrders)
	//r.GET("/delete-user/:id", Admin.DeleteUser)

	r.GET("/admin-settings/:data", Admin.AdminSettings)

	r.GET("/edit-admin-setting/:slug", Admin.EditPrimarySetting)
	r.POST("/update-admin-setting", Admin.UpdatePrimarySetting)
	r.GET("/delete-admin-setting/:id", Admin.DeletePrimarySetting)

	r.POST("/add-bank", Admin.AddBank)
	r.GET("/make-bank-inactive/:id", Admin.MakeBankInactive)
	r.GET("/make-bank-active/:id", Admin.MakeBankActive)
	r.GET("/edit-bank/:id", Admin.EditBank)
	r.POST("/update-bank", Admin.UpdateBank)
	//r.GET("/delete-bank/:id", Admin.DeleteBank)

	r.GET("/make-payMethod-inactive/:id", Admin.MakePayMethodInactive)
	r.GET("/make-payMethod-active/:id", Admin.MakePayMethodActive)
	r.GET("/edit-payMethod/:id", Admin.EditPayMethod)
	r.POST("/update-payMethod", Admin.UpdatePayMethod)
	//r.GET("/delete-payMethod/:id", Admin.DeletePayMethod)

	r.GET("/website-settings/:data", Admin.WebsiteSettings)

	r.POST("/update-website-settings", Admin.UpdateWebsiteSettings)

	r.POST("/add-page", Admin.AddPage)
	r.GET("/make-page-inactive/:id", Admin.MakePageInactive)
	r.GET("/make-page-active/:id", Admin.MakePageActive)
	r.GET("/edit-page/:id", Admin.EditPage)
	r.POST("/update-page", Admin.UpdatePage)
	r.GET("/delete-page/:id", Admin.DeletePage)

	r.GET("/edit-home-content/:id", Admin.EditHomeContent)
	r.POST("/update-home-content", Admin.UpdateHomeContent)

	r.POST("/add-menu", Admin.AddMenu)
	r.GET("/make-menu-inactive/:id", Admin.MakeMenuInactive)
	r.GET("/make-menu-active/:id", Admin.MakeMenuActive)
	r.GET("/edit-menu/:id", Admin.EditMenu)
	r.POST("/update-menu", Admin.UpdateMenu)
	r.GET("/delete-menu/:id", Admin.DeleteMenu)

	r.GET("/coupon-settings", Admin.CouponSettings)
	r.GET("/add-coupon", Admin.AddCouponGet)
	r.POST("/add-coupon", Admin.AddCouponPost)
	r.GET("/make-coupon-inactive/:id", Admin.MakeCouponInactive)
	r.GET("/make-coupon-active/:id", Admin.MakeCouponActive)
	r.GET("/edit-coupon/:id", Admin.EditCoupon)
	r.POST("/update-coupon", Admin.UpdateCoupon)
	//r.GET("/delete-coupon/:id", Admin.DeleteCoupon)

	//r.POST("/check", Admin.Check)

	r.GET("/wallets", Admin.Wallets)
	r.GET("/wallet-history/:id", Admin.WalletHistories)

	r.GET("/gold-transfers", Admin.GoldTransfers)
	r.GET("/make-transfer-cancel/:id", Admin.MakeTransferCancel)
	r.POST("/add-transfer-track-code", Admin.AddTransferTrackCode)
	r.GET("/transfer-shipping-details/:id", Admin.TransferShippingDetails)


	//datatable-test
	r.GET("/order-datatable", Admin.OrderDatatable)
	r.POST("/order-datatable", Admin.OrderDatatable)
}
