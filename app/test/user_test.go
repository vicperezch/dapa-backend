package test

import (
	"dapa/app/handlers"
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	utils.Load()

	code := m.Run()
	os.Exit(code)
}

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{})
	return db
}

func TestGetUsers_AdminAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	db := setupTestDB()
	database.DB = db

	// Datos de prueba
	db.Create(&model.User{Name: "Test", Email: "test@example.com", IsActive: true})

	c.Set("claims", &model.EmployeeClaims{Role: "admin"})

	req, _ := http.NewRequest("GET", "/users", nil)
	c.Request = req

	handlers.GetUsersHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test@example.com")
}

func TestGetUsers_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// No administrador
	c.Set("claims", &model.EmployeeClaims{Role: "user"})

	req, _ := http.NewRequest("GET", "/users", nil)
	c.Request = req

	handlers.GetUsersHandler(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Insufficient permissions")
}
