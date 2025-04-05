package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	LastName string `json:"lastName" gorm:"column:lastName"`
	Phone    string `json:"phone"`
	Email    string `json:"email" gorm:"default:unknown"`
}