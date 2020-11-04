package Models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	FirstName             string `gorm:"not null; default:''" form:"first_name" binding:"required"`
	LastName              string `gorm:"not null; default:''" form:"last_name" binding:"required"`
	RegType               int    `gorm:"type:tinyint; not null; default:1" json:"-"`
	Email                 string `gorm:"not null; unique" form:"email" binding:"email"`
	Phone                 string `gorm:"not null; unique" form:"phone" binding:"required"`
	CountryCode           string `form:"country" binding:"required"`
	City                  string `form:"city" binding:"required"`
	Address               string `form:"address" binding:"required"`
	ZipCode               int    `form:"zip_code" binding:"required"`
	PhoneVerification     string `json:"-"`
	Password              string `gorm:"not null; default:''" json:"-" form:"password" binding:"required"`
	ActiveStatus          int    `gorm:"type:tinyint(4); not null; default:0" json:"-"`
	RoleID                uint   `gorm:"index; not null"`
	SocialAuthID          string `json:"-"`
	EmailVerification     string `json:"-"`
	EmailVerifyCode       string `json:"-"`
	ProfilePic            string
	Role                  Role                  `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	Orders                []Order               `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	Wishlist              []Wishlist            `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	Wallet                Wallet                `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	GoldTransfers         []GoldTransfer        `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	EmailChangeHistories  []EmailChangeHistory  `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	DefaultAddress        ShippingAddress       `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	ShippingAddresses     []ShippingAddress     `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	OldEmail              string                `gorm:"-"`
}

/*type User struct {
	ID int	`gorm:"primary_key"`
	FirstName         string	`gorm:"not null" form:"first_name"`
	LastName          string	`gorm:"not null" form:"last_name"`
	Email             string	`gorm:"not null" form:"email"`
	Password          string	`gorm:"not null"`
	CountryCode		  sql.NullString
	Phone             sql.NullString
	PhoneNumber       string
	Photo sql.NullString
	Country string
	ActiveStatus      int	`gorm:"not null; default:0"`
	ActivityStatus    int	`gorm:"not null; default:9"`
	Role            int	`gorm:"not null"`
	Level sql.NullInt64
	Notification sql.NullString
	EmailVerified string	`gorm:"not null; default:0"`
	PhoneVerified string	`gorm:"not null; default:0"`
	ResetToken     sql.NullString
	GoogleCode sql.NullString
	LoggedInWith2fa int		`gorm:"column:logged_in_with_2fa; not null; default:0"`
	LoggedInWith2faFromApp int	`gorm:"column:logged_in_with_2fa_from_app; not null; default:0"`
	RememberToken sql.NullString
	StripeID sql.NullString
	CardBrand sql.NullString
	CardLastFour sql.NullString
	TrialEndsAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	FullAddress	sql.NullString	`gorm:"save_associations:false; association_save_reference:false"`
	Orders  []Order	`gorm:"save_associations:false; association_save_reference:false"`
}*/
