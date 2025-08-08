package test

import (
	"bytes"
	"dapa/app/handlers"
	"dapa/app/model"
	"dapa/database"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupVehicleTestContext() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.Vehicle{})
	database.DB = db
	return db
}

func addAdminClaims(c *gin.Context) {
	c.Set("claims", &model.EmployeeClaims{Role: "admin"})
}

func TestGetVehicles_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupVehicleTestContext()
	db.Create(&model.Vehicle{Brand: "Toyota", IsActive: true})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	addAdminClaims(c)
	req, _ := http.NewRequest("GET", "/vehicles", nil)
	c.Request = req

	handlers.GetVehicles(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetVehicleById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupVehicleTestContext()

	vehicle := model.Vehicle{
		Brand:         "Toyota",
		IsActive:      true,
		CreatedAt:     time.Now(),
		CapacityKg:    100,
		InsuranceDate: time.Now(),
	}
	db.Create(&vehicle)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	addAdminClaims(c)

	req, _ := http.NewRequest("GET", "/vehicles/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = req

	handlers.GetVehicleById(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Toyota")
}

func TestCreateVehicle_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupVehicleTestContext()

	vehicle := model.CreateVehicleRequest{
		Brand:         "Ford",
		Model:         "Transit",
		LicensePlate:  "XYZ-123",
		CapacityKg:    1000,
		Available:     true,
		InsuranceDate: time.Now(),
	}

	body, _ := json.Marshal(vehicle)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("POST", "/vehicles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	handlers.CreateVehicle(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Vehicle created successfully")
}
