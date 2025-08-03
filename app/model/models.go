package model

import (
	"time"
)

type FormStatus string

const (
	FormStatusPending   FormStatus = "pending"
	FormStatusCancelled FormStatus = "cancelled"
	FormStatusApproved  FormStatus = "approved"
)

type User struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	Name           string     `json:"name" gorm:"size:50;not null"`
	LastName       string     `json:"lastName" gorm:"column:last_name;size:50;not null"`
	Phone          string     `json:"phone" gorm:"size:20;not null"`
	CreatedAt	   time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
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
	CapacityKg             float64    `json:"capacityKg" gorm:"column:capacity_kg;check:capacity_kg > 0"`
	Available              *bool      `json:"available"`
	CurrentMileage         float64    `json:"currentMileage" gorm:"column:current_mileage;not null;check:current_mileage > 0"`
	NextMaintenanceMileage float64    `json:"nextMaintenanceMileage" gorm:"column:next_maintenance_mileage;not null;check:next_maintenance_mileage > 0"`
	CreatedAt 			   time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	LastModifiedAt         time.Time  `json:"lastModifiedAt" gorm:"column:last_modified_at;autoUpdateTime"`
	DeletedAt              *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	IsActive               bool       `json:"isActive" gorm:"column:is_active;default:true"`
}

// ******************** FORMULARIO ********************
// Tipos de preguntas disponibles
type QuestionType struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Type string `json:"type" gorm:"size:50;not null"`
}

// Pregunta del formulario
type Question struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Question    string         `json:"question" gorm:"size:50;not null"`
	Description *string        `json:"description,omitempty" gorm:"size:255"`
	TypeID      uint           `json:"typeId" gorm:"not null"`
	IsActive    bool           `json:"isActive" gorm:"not null;default:true"`
	Type        QuestionType   `json:"type" gorm:"foreignKey:TypeID"`
	Options     []QuestionOption `json:"options,omitempty" gorm:"foreignKey:QuestionID"`
}

// Opciones de una pregunta
type QuestionOption struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	QuestionID uint   `json:"questionId" gorm:"not null"`
	Option     string `json:"option" gorm:"size:50;not null"`
}

// Envío de formulario
type Submission struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UserID      uint       `json:"userId" gorm:"not null"`
	SubmittedAt time.Time  `json:"submittedAt" gorm:"default:CURRENT_TIMESTAMP"`
	Status      FormStatus `json:"status" gorm:"type:form_status;not null;default:'pending'"`
	User        User       `json:"user" gorm:"foreignKey:UserID"`
	Answers     []Answer   `json:"answers,omitempty" gorm:"foreignKey:SubmissionID"`
}

// Respuesta a una pregunta
type Answer struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	SubmissionID uint            `json:"submissionId" gorm:"not null"`
	QuestionID   *uint           `json:"questionId,omitempty"`
	Answer       *string         `json:"answer,omitempty"`
	OptionID     *uint           `json:"optionId,omitempty"`
	Question     *Question       `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	Option       *QuestionOption `json:"option,omitempty" gorm:"foreignKey:OptionID"`
}

