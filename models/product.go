package models

type Product struct {
	ID         uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string   `json:"name" gorm:"not null"`
	Price      float64  `json:"price" gorm:"not null"`
	Stock      int      `json:"stock" gorm:"default:0"`
	BrandID    uint     `json:"brand_id"`
	CategoryID uint     `json:"category_id"`
	Brand      Brand    `json:"brand,omitempty" gorm:"foreignKey:BrandID"`
	Category   Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}
