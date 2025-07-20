package model

import (
	"time"
)

type User struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	Name           string     `json:"name" gorm:"size:50;not null"`
	LastName       string     `json:"lastName" gorm:"column:last_name;size:50;not null"`
	Phone          string     `json:"phone" gorm:"size:20;not null"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	Email          string     `json:"email" gorm:"default:unknown;unique"`
	LastModifiedAt time.Time  `json:"lastModifiedAt" gorm:"column:last_modified_at;autoUpdateTime"`
	DeletedAt      *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	IsActive       bool       `json:"isActive" gorm:"column:is_active;default:true"`
}

type Employee struct {
	ID                uint `json:"id" gorm:"primaryKey"`
	UserID            uint
	User              User      `gorm:"constraint:OnDelete:CASCADE"`
	LicenseExpiration time.Time `json:"licenseExpirationDate" gorm:"column:license_expiration_date"`
	Password          string    `json:"password"`
	Role              string    `json:"role" gorm:"size:20;not null"`
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
	CapacityKg             float64    `json:"capacityKg" gorm:"column:capacity_kg"`
	Available              *bool      `json:"available"`
	CurrentMileage         float64    `json:"currentMileage" gorm:"column:current_mileage;not null;"`
	NextMaintenanceMileage float64    `json:"nextMaintenanceMileage" gorm:"column:next_maintenance_mileage;not null;"`
	CreatedAt              time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	LastModifiedAt         time.Time  `json:"lastModifiedAt" gorm:"column:last_modified_at;autoUpdateTime"`
	DeletedAt              *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	IsActive               bool       `json:"isActive" gorm:"column:is_active;default:true"`
}

type Quote struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	SubmissionID uint    `gorm:"column:submission_id"`
	DriverID     uint    `gorm:"column:driver_id"`
	VehicleID    uint    `gorm:"column:vehicle_id"`
	StateID      uint    `gorm:"column:state_id"`
	ServiceType  uint    `gorm:"column:service_type"`
	TotalAmount  float64 `gorm:"column:total_amount"`
	Date         time.Time
	Details      string
}

type QuoteView struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Date         time.Time `json:"date"`
	Client       string    `json:"clientName"`
	VehicleBrand string    `json:"vehicleBrand" gorm:"column:vehicle_brand"`
	VehicleModel string    `json:"vehicleModel" gorm:"column:vehicle_model"`
	Driver       string    `json:"driver"`
}

func (QuoteView) TableName() string {
	return "pending_orders"
}
