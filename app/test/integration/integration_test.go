package integration

import (
	"bytes"
	"dapa/app/model"
	"dapa/app/routes"
	"dapa/app/utils"
	"dapa/database"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// =======================
//   SETUP GLOBAL
// =======================

// setupTestDB crea una conexión a PostgreSQL y configura la BD para testing
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=database_test user=test_user password=test_password dbname=dapa_test port=5432 sslmode=disable TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error conectando a la base de datos de test: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Error obteniendo la conexión a la base de datos: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}

	// IMPORTANTE: Configurar la variable global que usan los handlers
	database.DB = db

	t.Log("Base de datos de testing conectada correctamente")
	return db
}

// setupRouter crea el router real con las rutas de la app
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", utils.PasswordValidator)
		v.RegisterValidation("phone", utils.PhoneValidator)
		v.RegisterValidation("question_text", utils.QuestionTextValidator)
		v.RegisterValidation("question_desc", utils.QuestionDescriptionValidator)
		v.RegisterValidation("question_type", utils.QuestionTypeValidator)
		v.RegisterValidation("question_option", utils.QuestionOptionValidator)
		v.RegisterValidation("submission_status", utils.SubmissionStatusValidator)
	}
	routes.SetupRoutes(router)
	return router
}

// loginAndGetToken hace login real y retorna un token válido para los tests
func loginAndGetToken(t *testing.T, router *gin.Engine, email, password string) string {
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "El login debe devolver 200. Response: %s", w.Body.String())

	var resp model.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	token, ok := resp.Data.(string)
	assert.True(t, ok, "El token JWT debe ser un string")
	assert.NotEmpty(t, token, "El token no puede estar vacío")

	return token
}

// createTestAdmin crea un usuario admin para testing
func createTestAdmin(t *testing.T, db *gorm.DB) model.User {
	adminPassword := "supersecret123"
	hashedPassword, err := utils.HashPassword(adminPassword)
	assert.NoError(t, err)

	admin := model.User{
		Name:                  "Admin",
		LastName:              "User",
		Email:                 "admin@test.com",
		PasswordHash:          hashedPassword,
		Role:                  "admin",
		LicenseExpirationDate: time.Now().AddDate(1, 0, 0),
		IsActive:              true,
	}

	result := db.Create(&admin)
	assert.NoError(t, result.Error, "Error creando usuario admin")
	return admin
}

func deleteAdminTest(db *gorm.DB) {
	db.Exec("DELETE FROM users WHERE email = ?", "admin@test.com")
}

// =======================
//   TESTS DE INTEGRACIÓN
// =======================

// ----- CICLO DE VIDA DE USUARIOS -----
func TestUser_FullCycle_WithRealAuth(t *testing.T) {
	db := SetupTestDB(t)
	defer deleteAdminTest(db)
	// Migrar modelos necesarios
	err := db.AutoMigrate(&model.User{})
	assert.NoError(t, err)

	router := setupRouter()

	// 1. Crear un usuario admin dummy directamente en la BD
	admin := createTestAdmin(t, db)

	// 2. Hacer login y obtener token JWT
	token := loginAndGetToken(t, router, admin.Email, "supersecret123")
	t.Logf("Token obtenido: %s", token)

	// 3. Crear un nuevo usuario
	newUser := map[string]interface{}{
		"name":                  "John",
		"lastName":              "Doe",
		"phone":                 "555123456",
		"email":                 "john.doe@test.com",
		"password":              "dapa123456",
		"role":                  "driver",
		"licenseExpirationDate": time.Now().AddDate(1, 0, 0).Format("2006-01-02T15:04:05Z07:00"),
	}
	jsonBody, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Create user response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code, "Crear usuario debe devolver 200")

	// 4. Obtener el usuario recién creado desde la base de datos
	var createdUser model.User
	err = db.Where("email = ?", newUser["email"]).First(&createdUser).Error
	assert.NoError(t, err, "El usuario recién creado debe existir en la BD")
	userID := createdUser.ID

	// 5. Listar usuarios
	req, _ = http.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// t.Logf("List users response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	var users []model.User
	err = json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 2) // admin + nuevo

	// 6. Eliminar usuario (soft delete)
	req, _ = http.NewRequest("DELETE", "/api/users/"+strconv.Itoa(int(userID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Delete user response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. Verificar que está inactivo
	var deletedUser model.User
	result := db.First(&deletedUser, userID)
	assert.NoError(t, result.Error)
	assert.False(t, deletedUser.IsActive, "El usuario debe estar marcado como inactivo")
}

// ----- CICLO DE VIDA DE VEHÍCULOS -----
func TestVehicle_FullCycle_WithRealAuth(t *testing.T) {
	db := SetupTestDB(t)
	defer deleteAdminTest(db)
	// Migrar modelos necesarios
	err := db.AutoMigrate(&model.User{}, &model.Vehicle{})
	assert.NoError(t, err)

	router := setupRouter()

	// Crear admin dummy
	admin := createTestAdmin(t, db)
	token := loginAndGetToken(t, router, admin.Email, "supersecret123")

	// 1. Crear vehículo
	newVehicle := map[string]interface{}{
		"brand":         "Toyota",
		"model":         "Hilux",
		"licensePlate":  "P-123ABC",
		"capacityKg":    2500.0,
		"available":     true,
		"insuranceDate": time.Now().AddDate(1, 0, 0),
	}
	jsonBody, _ := json.Marshal(newVehicle)

	req, _ := http.NewRequest("POST", "/api/vehicles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Create vehicle response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code, "Crear vehículo debe devolver 200")

	// 2. Listar vehículos
	req, _ = http.NewRequest("GET", "/api/vehicles", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// t.Logf("List vehicles response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	var vehicles []model.Vehicle
	err = json.Unmarshal(w.Body.Bytes(), &vehicles)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(vehicles), 1, "Debe haber al menos 1 vehículo")

	// 3. Obtener vehículo por ID
	vehicleID := vehicles[0].ID
	req, _ = http.NewRequest("GET", "/api/vehicles/"+strconv.Itoa(int(vehicleID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Actualizar vehículo
	updateVehicle := map[string]interface{}{
		"brand":         "Toyota",
		"model":         "Hilux Updated",
		"licensePlate":  "P-123SBC",
		"capacityKg":    2000.0,
		"available":     true,
		"insuranceDate": time.Now().AddDate(2, 0, 0),
	}
	jsonBody, _ = json.Marshal(updateVehicle)

	req, _ = http.NewRequest("PUT", "/api/vehicles/"+strconv.Itoa(int(vehicleID)), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Update vehicle response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. Eliminar vehículo (soft delete)
	req, _ = http.NewRequest("DELETE", "/api/vehicles/"+strconv.Itoa(int(vehicleID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Delete vehicle response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. Verificar que está inactivo
	var deletedVehicle model.Vehicle
	result := db.First(&deletedVehicle, vehicleID)
	assert.NoError(t, result.Error)
	assert.False(t, deletedVehicle.IsActive, "El vehículo debe estar marcado como inactivo")
}

// ----- CICLO DE VIDA DE PREGUNTAS -----
func TestQuestion_FullCycle_WithRealAuth(t *testing.T) {
	db := SetupTestDB(t)
	defer deleteAdminTest(db)
	// Migrar modelos necesarios
	err := db.AutoMigrate(&model.User{}, &model.Question{}, &model.QuestionType{}, &model.QuestionOption{})
	assert.NoError(t, err)

	router := setupRouter()

	// Crear admin dummy
	admin := createTestAdmin(t, db)
	token := loginAndGetToken(t, router, admin.Email, "supersecret123")

	// 1. Crear tipo de pregunta
	qType := model.QuestionType{Type: "text"}
	result := db.Create(&qType)
	assert.NoError(t, result.Error)

	// 2. Crear pregunta
	newQuestion := map[string]interface{}{
		"question":    "¿Cuál es tu nombre?",
		"description": "Pregunta sobre el nombre del usuario",
		"typeId":      qType.ID,
		"isActive":    true,
	}
	jsonBody, _ := json.Marshal(newQuestion)

	req, _ := http.NewRequest("POST", "/api/questions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Create question response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code, "Crear pregunta debe devolver 200")

	var createResp model.ApiResponse
	err = json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.NoError(t, err)

	// Get question ID from response
	questionData, ok := createResp.Data.(map[string]interface{})
	assert.True(t, ok, "Data should be a map")
	questionID := uint(questionData["id"].(float64))

	// 3. Listar preguntas
	req, _ = http.NewRequest("GET", "/api/questions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// t.Logf("List questions response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	var questions []model.Question
	err = json.Unmarshal(w.Body.Bytes(), &questions)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(questions), 1, "Debe haber al menos 1 pregunta")

	// 4. Actualizar pregunta
	update := map[string]interface{}{
		"question": "¿Cuál es tu apellido?",
	}
	jsonBody, _ = json.Marshal(update)

	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/questions/%d", questionID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Update question response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. Activar/desactivar pregunta
	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/api/questions/%d/active", questionID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Toggle question response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. Eliminar pregunta (hard delete)
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/questions/%d", questionID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Delete question response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. Verificar que la pregunta fue eliminada
	var deletedQuestion model.Question
	result = db.First(&deletedQuestion, questionID)
	assert.Error(t, result.Error, "La pregunta debe haber sido eliminada")
}

// ----- TEST ADICIONAL: AUTENTICACIÓN -----
func TestAuthentication_Flow(t *testing.T) {
	db := SetupTestDB(t)
	defer deleteAdminTest(db)
	err := db.AutoMigrate(&model.User{})
	assert.NoError(t, err)

	router := setupRouter()

	// 1. Crear usuario admin
	admin := createTestAdmin(t, db)

	// 2. Test login exitoso
	loginData := map[string]string{
		"email":    admin.Email,
		"password": "supersecret123",
	}
	jsonBody, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.ApiResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotEmpty(t, resp.Data)

	// 3. Test login con credenciales incorrectas
	wrongLogin := map[string]string{
		"email":    admin.Email,
		"password": "wrongpassword",
	}
	jsonBody, _ = json.Marshal(wrongLogin)

	req, _ = http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 4. Test acceso sin token
	req, _ = http.NewRequest("GET", "/api/users", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
