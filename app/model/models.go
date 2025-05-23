package model

import (
	"time"
)

type User struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	Name           string     `json:"name" gorm:"size:50;not null"`
	LastName       string     `json:"lastName" gorm:"column:last_name;size:50;not null"`
	Phone          string     `json:"phone" gorm:"size:20;not null"`
	Email          string     `json:"email" gorm:"default:unknown"`
	LastModifiedAt time.Time  `json:"lastModifiedAt" gorm:"column:last_modified_at;autoUpdateTime"`
	DeletedAt      *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	IsActive       bool       `json:"isActive" gorm:"column:is_active;default:true"`
}

type Employee struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	UserID   uint
	User     User   `gorm:"constraint:OnDelete:CASCADE"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"size:20;not null"`
}

type UserWithRole struct {
	User
	Role string `json:"role"`
}

type Vehicle struct {
	ID                     uint       `json:"id" gorm:"primaryKey"`
	Brand                  string     `json:"brand" gorm:"size:50;not null"`
	Model                  string     `json:"model" gorm:"size:50;not null"`
	LicensePlate           string     `json:"licensePlate" gorm:"size:15;not null;unique"`
	CapacityKg             float64    `json:"capacityKg" gorm:"column:capacity_kg;check:capacity_kg > 0"`
	Available              *bool      `json:"available"`
	CurrentMileage         float64    `json:"currentMileage" gorm:"column:current_mileage;not null;check:current_mileage > 0"`
	NextMaintenanceMileage float64    `json:"nextMaintenanceMileage" gorm:"column:next_maintenance_mileage;not null;check:next_maintenance_mileage > 0"`
	LastModifiedAt         time.Time  `json:"lastModifiedAt" gorm:"column:last_modified_at;autoUpdateTime"`
	DeletedAt              *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	IsActive               bool       `json:"isActive" gorm:"column:is_active;default:true"`
}
