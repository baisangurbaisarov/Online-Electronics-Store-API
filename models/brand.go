package models

type Brand struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name" gorm:"not null;unique"`
	Country string `json:"country"`
}
