package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"size:50;not null"`
	LastName string `json:"lastName" gorm:"column:last_name;size:50;not null"`
	Phone    string `json:"phone" gorm:"size:20;not null"`
	Email    string `json:"email" gorm:"default:unknown"`
}

type Employee struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	UserID   uint
	User     User   `gorm:"constraint:OnDelete:CASCADE"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"size:20;not null"`
}
