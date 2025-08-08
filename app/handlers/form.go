package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ---------- TIPOS DE PREGUNTA ----------

// Crear tipo de pregunta
func CreateQuestionType(c *gin.Context) {
	var req model.CreateQuestionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Formato inválido", http.StatusBadRequest)
		return
	}
	qt := model.QuestionType{Type: req.Type}
	if err := database.DB.Create(&qt).Error; err != nil {
		utils.RespondWithError(c, "Error al crear tipo", http.StatusInternalServerError)
		return
	}
	utils.RespondWithJSON(c, model.ApiResponse{Success: true, Data: qt})
}

// Listar tipos de pregunta
func GetQuestionTypes(c *gin.Context) {
	var types []model.QuestionType
	if err := database.DB.Find(&types).Error; err != nil {
		utils.RespondWithError(c, "Error al obtener tipos", http.StatusInternalServerError)
		return
	}
	utils.RespondWithJSON(c, types)
}

// ---------- PREGUNTAS ----------

// Crear pregunta
func CreateQuestion(c *gin.Context) {
	var req model.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Formato inválido", http.StatusBadRequest)
		return
	}
	q := model.Question{
		Question:    req.Question,
		Description: req.Description,
		TypeID:      req.TypeID,
		IsActive:    req.IsActive == nil || *req.IsActive,
	}
	if err := database.DB.Create(&q).Error; err != nil {
		utils.RespondWithError(c, "Error al crear pregunta", http.StatusInternalServerError)
		return
	}
	// Crear opciones si existen
	for _, opt := range req.Options {
		option := model.QuestionOption{QuestionID: q.ID, Option: opt.Option}
		database.DB.Create(&option)
	}
	utils.RespondWithJSON(c, model.ApiResponse{Success: true, Data: q})
}

// Listar preguntas
func GetQuestions(c *gin.Context) {
	var questions []model.Question
	if err := database.DB.Preload("Options").Preload("Type").Find(&questions).Error; err != nil {
		utils.RespondWithError(c, "Error al obtener preguntas", http.StatusInternalServerError)
		return
	}
	utils.RespondWithJSON(c, questions)
}

func GetActiveQuestions(c *gin.Context) {
	var activeQuestions []model.Question
	if err := database.DB.Preload("Options").Preload("Type").Find(&activeQuestions).Where("is_active = false").Error; err != nil {
		utils.RespondWithError(c, "Error al obtener preguntas", http.StatusInternalServerError)
		return
	}
	utils.RespondWithJSON(c, activeQuestions)
}

// Obtener pregunta por ID
func GetQuestionByID(c *gin.Context) {
	id := c.Param("id")
	var question model.Question
	if err := database.DB.Preload("Options").Preload("Type").First(&question, id).Error; err != nil {
		utils.RespondWithError(c, "Pregunta no encontrada", http.StatusNotFound)
		return
	}
	utils.RespondWithJSON(c, question)
}

// Actualizar pregunta
func UpdateQuestion(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateQuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Formato inválido", http.StatusBadRequest)
		return
	}

	var question model.Question
	if err := database.DB.Preload("Type").Preload("Options").First(&question, id).Error; err != nil {
		utils.RespondWithError(c, "Pregunta no encontrada", http.StatusNotFound)
		return
	}

	// Usar una transacción para garantizar consistencia
	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.RespondWithError(c, "Error iniciando transacción", http.StatusInternalServerError)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Actualizar campos básicos
	if req.Question != nil {
		question.Question = *req.Question
	}
	if req.Description != nil {
		question.Description = req.Description
	}
	if req.TypeID != nil {
		question.TypeID = *req.TypeID
		var newType model.QuestionType
		if err := tx.First(&newType, question.TypeID).Error; err != nil {
			tx.Rollback()
			utils.RespondWithError(c, "Error recargando tipo de pregunta", http.StatusInternalServerError)
			return
		}
		question.Type = newType
	}
	if req.IsActive != nil {
		question.IsActive = *req.IsActive
	}

	// Guardar cambios básicos de la pregunta
	if err := tx.Save(&question).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(c, "Error actualizando pregunta", http.StatusInternalServerError)
		return
	}

	// Obtener el tipo actual de la pregunta para validación
	var questionType model.QuestionType
	if err := tx.First(&questionType, question.TypeID).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(c, "Error obteniendo tipo de pregunta", http.StatusInternalServerError)
		return
	}

	// Verificar si el tipo requiere opciones
	requiresOptions := questionType.Type == "multiple" || questionType.Type == "dropdown" || questionType.Type == "unique"

	// Manejar opciones si vienen en el request O si el tipo cambió
	if req.Options != nil {
		// Validación: preguntas con opciones deben tener al menos una
		if requiresOptions && len(req.Options) == 0 {
			tx.Rollback()
			utils.RespondWithError(c, "Las preguntas de tipo múltiple, lista o única deben tener al menos una opción", http.StatusBadRequest)
			return
		}

		// Si el tipo NO requiere opciones pero se enviaron, ignorarlas
		if !requiresOptions {
			// Eliminar opciones existentes ya que el tipo no las necesita
			if err := tx.Where("question_id = ?", question.ID).Delete(&model.QuestionOption{}).Error; err != nil {
				tx.Rollback()
				utils.RespondWithError(c, "Error eliminando opciones existentes", http.StatusInternalServerError)
				return
			}
		} else {
			// El tipo SÍ requiere opciones, procesarlas
			// Eliminar TODAS las opciones existentes de esta pregunta
			if err := tx.Where("question_id = ?", question.ID).Delete(&model.QuestionOption{}).Error; err != nil {
				tx.Rollback()
				utils.RespondWithError(c, "Error eliminando opciones existentes", http.StatusInternalServerError)
				return
			}

			// Crear las nuevas opciones
			if len(req.Options) > 0 {
				var newOptions []model.QuestionOption
				for _, optReq := range req.Options {
					newOption := model.QuestionOption{
						QuestionID: question.ID,
						Option:     optReq.Option,
					}
					newOptions = append(newOptions, newOption)
				}

				// Insertar todas las opciones de una vez
				if err := tx.Create(&newOptions).Error; err != nil {
					tx.Rollback()
					utils.RespondWithError(c, "Error creando nuevas opciones", http.StatusInternalServerError)
					return
				}
			}
		}
	} else {
		// req.Options es nil - no se enviaron opciones
		// Si el tipo cambió a uno que no requiere opciones, eliminarlas
		if !requiresOptions {
			if err := tx.Where("question_id = ?", question.ID).Delete(&model.QuestionOption{}).Error; err != nil {
				tx.Rollback()
				utils.RespondWithError(c, "Error eliminando opciones existentes", http.StatusInternalServerError)
				return
			}
		}
		// Si el tipo requiere opciones pero no se enviaron, no hacer nada (mantener existentes)
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		utils.RespondWithError(c, "Error confirmando cambios", http.StatusInternalServerError)
		return
	}

	// Recargar la pregunta con sus relaciones actualizadas para la respuesta
	if err := database.DB.Preload("Type").Preload("Options").First(&question, question.ID).Error; err != nil {
		utils.RespondWithError(c, "Error recargando pregunta actualizada", http.StatusInternalServerError)
		return
	}

	// Convertir a response DTO
	response := model.QuestionResponse{
		ID:          question.ID,
		Question:    question.Question,
		Description: question.Description,
		TypeID:      question.TypeID,
		Type:        question.Type.Type,
		IsActive:    question.IsActive,
		Options:     make([]model.QuestionOptionResponse, 0),
	}

	// Mapear opciones a response
	for _, opt := range question.Options {
		optResponse := model.QuestionOptionResponse{
			ID:     opt.ID,
			Option: opt.Option,
		}
		response.Options = append(response.Options, optResponse)
	}

	utils.RespondWithJSON(c, response)
}

// Eliminar pregunta (soft delete)
func DeleteQuestion(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Model(&model.Question{}).
		Delete("id = ?", id).Error; err != nil {
		utils.RespondWithError(c, "Error al eliminar pregunta", http.StatusInternalServerError)
		return
	}
	utils.RespondWithJSON(c, model.ApiResponse{Success: true, Message: "Pregunta eliminada"})
}

// ---------- OPCIONES DE PREGUNTA ----------

// Crear opción
func CreateQuestionOption(c *gin.Context) {
	var req model.QuestionOptionRequest
	questionID, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		utils.RespondWithError(c, "ID inválido", http.StatusBadRequest)
		return
	}
	option := model.QuestionOption{QuestionID: uint(questionID), Option: req.Option}
	if err := database.DB.Create(&option).Error; err != nil {
		utils.RespondWithError(c, "Error al crear opción", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{Success: true, Data: option})
}

// ---------- ENVÍOS DE FORMULARIO ----------

// Crear envío
func CreateSubmission(c *gin.Context) {
	var req model.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Formato inválido", http.StatusBadRequest)
		return
	}
	sub := model.Submission{
		UserID:      req.UserID,
		SubmittedAt: time.Now(),
		Status:      model.FormStatusPending,
	}
	if err := database.DB.Create(&sub).Error; err != nil {
		utils.RespondWithError(c, "Error al crear envío", http.StatusInternalServerError)
		return
	}
	// Guardar respuestas
	for _, ans := range req.Answers {
		answer := model.Answer{
			SubmissionID: sub.ID,
			QuestionID:   ans.QuestionID,
			Answer:       ans.Answer,
			OptionID:     ans.OptionID,
		}
		database.DB.Create(&answer)
	}
	utils.RespondWithJSON(c, model.ApiResponse{Success: true, Data: sub})
}

// Listar envíos
func GetSubmissions(c *gin.Context) {
	var submissions []model.Submission
	if err := database.DB.Preload("Answers").Preload("User").Find(&submissions).Error; err != nil {
		utils.RespondWithError(c, "Error al obtener envíos", http.StatusInternalServerError)
		return
	}
	utils.RespondWithJSON(c, submissions)
}

// Actualizar estado de envío
func UpdateSubmissionStatus(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateSubmissionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Formato inválido", http.StatusBadRequest)
		return
	}
	if err := database.DB.Model(&model.Submission{}).Where("id = ?", id).
		Update("status", req.Status).Error; err != nil {
		utils.RespondWithError(c, "Error al actualizar estado", http.StatusInternalServerError)
		return
	}
	utils.RespondWithJSON(c, model.ApiResponse{Success: true, Message: "Estado actualizado"})
}

// @Summary		Get all quotes
// @Description	Returns a list of all quotes in the system.
// @Tags		quotes
// @Produce		json
// @Success		200	{array} model.Quote "List of quotes"
// @Failure		403	{object} model.ApiResponse "Insufficient permissions"
// @Failure		500	{object} model.ApiResponse "Error fetching quotes"
// @Router		/quotes/ [get]
func GetQuotes(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)

	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var quotes []model.QuoteView

	if err := database.DB.Find(&quotes).Error; err != nil {
		log.Println("Error fetching quotes:", err)
		utils.RespondWithError(c, "Error getting all quotes", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, quotes)
}

// @Summary		Assign a driver and vehicle to an order
// @Description	Adds the driver and vehicle id to the quotes table.
// @Tags		quotes
// @Accept		json
// @Produce		json
// @Param		user body model.AssignQuote true "Information to add"
// @Success		200	{object} model.ApiResponse "Successfully assigned"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error assigning driver and vehicle"
// @Router		/quotes/{id} [post]
func AssignQuoteInfo(c *gin.Context) {
	id := c.Param("id")
	var req model.AssignQuote

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Binding error: ", err)
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	var quote model.Quote

	if err := database.DB.First(&quote, id).Error; err != nil {
		log.Println("Error fetching quote: ", err)
		utils.RespondWithError(c, "Error assigning information", http.StatusInternalServerError)
	}

	quote.DriverID = req.Driver
	quote.VehicleID = req.Vehicle
	quote.Details = req.Details

	if err := database.DB.Save(&quote).Error; err != nil {
		log.Println("Error updating quote: ", err)
		utils.RespondWithError(c, "Error assigning information", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Quote updated successfully",
	})
}
