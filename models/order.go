package models

import "time"

type Order struct {
	ID        uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint        `json:"user_id" gorm:"not null"`
	Status    string      `json:"status" gorm:"not null;default:placed"`
	Total     float64     `json:"total" gorm:"not null"`
	CreatedAt time.Time   `json:"created_at"`
	User      User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items     []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   uint    `json:"order_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}
