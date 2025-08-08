package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"
	"strconv"
	"time"
    "log"

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
    if err := database.DB.First(&question, id).Error; err != nil {
        utils.RespondWithError(c, "Pregunta no encontrada", http.StatusNotFound)
        return
    }
    if req.Question != nil {
        question.Question = *req.Question
    }
    if req.Description != nil {
        question.Description = req.Description
    }
    if req.TypeID != nil {
        question.TypeID = *req.TypeID
    }
    if req.IsActive != nil {
        question.IsActive = *req.IsActive
    }
    if err := database.DB.Save(&question).Error; err != nil {
        utils.RespondWithError(c, "Error al actualizar pregunta", http.StatusInternalServerError)
        return
    }
    // Opciones: puedes actualizar aquí si lo necesitas
    utils.RespondWithJSON(c, model.ApiResponse{Success: true, Data: question})
}

// Eliminar pregunta (soft delete)
func DeleteQuestion(c *gin.Context) {
    id := c.Param("id")
    if err := database.DB.Model(&model.Question{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "deleted_at": time.Now(),
            "is_active":  false,
        }).Error; err != nil {
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
