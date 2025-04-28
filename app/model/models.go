package model

import (
	"time"
)
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"size:50;not null"`
	LastName string `json:"lastName" gorm:"column:last_name;size:50;not null"`
	Phone    string `json:"phone" gorm:"size:20;not null"`
	Email    string `json:"email" gorm:"default:unknown"`
	LastModifiedAt time.Time `json:"lastModifiedAt" gorm:"column:last_modified_at;autoUpdateTime"`
	DeletedAt  *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	IsActive   bool       `json:"isActive" gorm:"column:is_active;default:true"`
}

type Employee struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	UserID   uint
	User     User   `gorm:"constraint:OnDelete:CASCADE"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"size:20;not null"`
}
