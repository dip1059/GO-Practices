package Models

import "time"

type GoldTransfer struct {
	ID                  uint      `gorm:"primary_key"`
	CreatedAt           time.Time `gorm:"not null; unique_index:idx_created_at_sender_wallet_id_status"`
	UpdatedAt           time.Time
	DeletedAt           *time.Time              `sql:"index"`
	TransferID          string                  `gorm:"not null; unique"`
	SenderWalletID      uint                    `gorm:"not null; unique_index:idx_created_at_sender_wallet_id_status"`
	ReceiverWalletID    uint                    `gorm:"not null; unique_index:idx_created_at_sender_wallet_id_status"`
	GoldAmount          float64                 `gorm:"not null; default:0.0; unique_index:idx_created_at_sender_wallet_id_status"`
	SenderFeesPercent   float64                 `gorm:"not null; default:0.00"`
	SenderFeesFixed     float64                 `gorm:"not null; default:0.00"`
	SenderTotalFees     float64                 `gorm:"not null; default:0.00"`
	AmountToDeduct      float64                 `gorm:"not null; default:0.00"`
	ReceiverFeesPercent float64                 `gorm:"not null; default:0.00"`
	ReceiverFeesFixed   float64                 `gorm:"not null; default:0.00"`
	ReceiverTotalFees   float64                 `gorm:"not null; default:0.00"`
	ReceiverAmount      float64                 `gorm:"not null; default:0.00"`
	TotalFees           float64                 `gorm:"not null; default:0.0"`
	DeliveryType        int                     `gorm:"type:tinyint(4); not null"`
	TrackCode           string                  `gorm:"not null; default:'N/A'"`
	EmailVerification   string                  `json:"-"`
	EmailVerifyCode     string                  `json:"-"`
	Status              int                     `gorm:"type:tinyint(4); not null; default:0"`
	SenderUser          User                    `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	ReceiverUser        User                    `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	ShipAddress         TransferShippingAddress `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
}
