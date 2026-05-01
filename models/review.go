package models

type Review struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	UserID    uint    `json:"user_id" gorm:"not null"`
	Rating    int     `json:"rating" gorm:"not null"` // 1–5
	Comment   string  `json:"comment"`
	Sentiment string  `json:"sentiment"` // filled by external API
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	User      User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}