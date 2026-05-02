package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"-" gorm:"not null"`
}