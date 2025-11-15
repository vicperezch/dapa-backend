package model

import (
	"gorm.io/gorm"
	"time"
)

type FormStatus string

const (
	FormStatusPending   FormStatus = "pending"
	FormStatusCancelled FormStatus = "cancelled"
	FormStatusApproved  FormStatus = "approved"
)

type VerificationResponse struct {
	AcceptAll bool   `json:"accept_all"`
	State     string `json:"state"`
}

// ******************** ENTIDADES ********************
type User struct {
	ID                    uint       `json:"id" gorm:"primaryKey"`
	Name                  string     `json:"name" gorm:"size:50;not null"`
	LastName              string     `json:"lastName" gorm:"column:last_name;size:50;not null"`
	Phone                 string     `json:"phone" gorm:"size:20;not null"`
	Email                 string     `json:"email" gorm:"size:50;unique;not null"`
	PasswordHash          string     `json:"password" gorm:"column:password_hash;size:255;not null"`
	Role                  string     `json:"role" gorm:"not null"`
	LicenseExpirationDate time.Time  `json:"licenseExpirationDate" gorm:"column:license_expiration_date"`
	CreatedAt             time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	LastModifiedAt        time.Time  `json:"lastModifiedAt" gorm:"column:last_modified_at;autoUpdateTime"`
	DeletedAt             *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	IsActive              bool       `json:"isActive" gorm:"column:is_active;default:true"`
}

type Vehicle struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	Brand          string     `json:"brand" gorm:"size:50;not null"`
	Model          string     `json:"model" gorm:"size:50;not null"`
	LicensePlate   string     `json:"licensePlate" gorm:"size:15;not null"`
	CapacityKg     float64    `json:"capacityKg" gorm:"column:capacity_kg"`
	IsAvailable    bool       `json:"isAvailable" gorm:"column:is_available;not null;default:true"`
	InsuranceDate  time.Time  `json:"insuranceDate" gorm:"column:insurance"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	LastModifiedAt time.Time  `json:"lastModifiedAt" gorm:"column:last_modified_at;autoUpdateTime"`
	DeletedAt      *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	IsActive       bool       `json:"isActive" gorm:"column:is_active;default:true"`
}

type ResetToken struct {
	ID     uint      `gorm:"primaryKey"`
	Token  string    `gorm:"size:255;not null"`
	Expiry time.Time `gorm:"not null"`
	IsUsed bool      `gorm:"not null"`
	UserID uint      `gorm:"not null"`
}

type Order struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SubmissionID uint      `json:"submissionId"`
	UserID       *uint     `json:"userId,omitempty" gorm:"default:null"`
	VehicleID    *uint     `json:"vehicleId,omitempty" gorm:"default:null"`
	HelperID     *uint     `json:"helperId,omitempty" gorm:"default:null"`
	ClientName   string    `json:"clientName"`
	ClientPhone  string    `json:"clientPhone"`
	Origin       string    `json:"origin" gorm:"size:100;not null"`
	Destination  string    `json:"destination" gorm:"size:100;not null"`
	TotalAmount  float64   `json:"totalAmount" gorm:"not null"`
	Details      string    `json:"details"`
	Status       string    `json:"status" gorm:"default:pending"`
	Type         string    `json:"type" gorm:"not null"`
	Date         time.Time `json:"date" gorm:"type:date"`
	MeetingDate  time.Time `json:"meetingDate" gorm:"type:date;not null"`
}

type OrderToken struct {
	ID      uint       `json:"id" gorm:"primaryKey"`
	OrderID uint       `json:"orderId" gorm:"unique;not null;column:order_id"`
	Token   string     `json:"token" gorm:"not null;unique;size:255"`
	Expiry  *time.Time `json:"expiry"`
}

// ******************** FORMULARIO ********************
// Tipos de preguntas disponibles
type QuestionType struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Type string `json:"type" gorm:"size:50;not null" validate:"required,question_type"`
}

// Pregunta del formulario
type Question struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Question    string           `json:"question" gorm:"size:50;not null" validate:"required,question_text"`
	Description *string          `json:"description,omitempty" gorm:"size:255" validate:"omitempty,question_desc"`
	TypeID      uint             `json:"typeId" gorm:"not null" validate:"required,gt=0"`
	IsActive    bool             `json:"isActive" gorm:"not null;default:true"`
	Position    int              `json:"position" gorm:"not null;default:1"`
	IsRequired  bool             `json:"isRequired" gorm:"not null;default:true"`
	IsMutable   bool             `json:"isMutable" gorm:"not null;default:true"`
	Type        QuestionType     `json:"type" gorm:"foreignKey:TypeID"`
	Options     []QuestionOption `json:"options,omitempty" gorm:"foreignKey:QuestionID"`
	DeletedAt   gorm.DeletedAt   `json:"deletedAt" gorm:"index"`
}

// Opciones de una pregunta
type QuestionOption struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	QuestionID uint           `json:"questionId" gorm:"not null" validate:"required,gt=0"`
	Option     string         `json:"option" gorm:"size:50;not null" validate:"required,question_option"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// Envío de formulario
type Submission struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	SubmittedAt time.Time  `json:"submittedAt" gorm:"default:CURRENT_TIMESTAMP"`
	Status      FormStatus `json:"status" gorm:"type:form_status;not null;default:'pending'" validate:"required,submission_status"`
	Answers     []Answer   `json:"answers,omitempty" gorm:"foreignKey:SubmissionID"`
}

// Respuesta a una pregunta
type Answer struct {
	ID           uint             `json:"id" gorm:"primaryKey"`
	SubmissionID uint             `json:"submissionId" gorm:"not null" validate:"required,gt=0"`
	QuestionID   uint             `json:"questionId,omitempty" validate:"omitempty,gt=0"`
	Answer       *string          `json:"answer,omitempty" validate:"omitempty,max=255"`
	Question     Question         `json:"question" gorm:"foreignKey:QuestionID"`
	Options      []QuestionOption `json:"options,omitempty" gorm:"many2many:answer_options;"`
}

type ExpenseType struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Type string `json:"type" gorm:"not null;size:25"`
}

type Expense struct {
	ID               uint        `json:"id" gorm:"primaryKey"`
	Date             time.Time   `json:"date" gorm:"not null"`
	TypeID           uint        `json:"typeId" gorm:"not null"`
	Type             ExpenseType `json:"type" gorm:"foreignKey:TypeID"`
	TemporalEmployee bool        `json:"temporalEmployee" gorm:"not null;column:temporal_employee"`
	Description      string      `json:"description" gorm:"size:255"`
	Amount           float64     `json:"amount" gorm:"not null" validate:"gt=0"`
}

type PerformanceGoal struct {
	ID                  uint    `json:"id" gorm:"primaryKey"`
	OrderGoal           int     `json:"orderGoal" gorm:"not null"`
	UtilityGoal         float64 `json:"utilityGoal" gorm:"not null"`
	AveragePerOrderGoal float64 `json:"averagePerOrderGoal" gorm:"not null"`
	TravelGoal          int     `json:"travelGoal" gorm:"not null"`
	DeliveryGoal        float64 `json:"deliveryGoal" gorm:"not null"`
	AchievementRateGoal float64 `json:"achievementRateGoal" gorm:"not null"`
}
