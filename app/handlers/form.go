package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ---------- TIPOS DE PREGUNTA ----------
// Listar tipos de pregunta
func GetQuestionTypesHandler(c *gin.Context) {
	var types []model.QuestionType
	if err := database.DB.Find(&types).Error; err != nil {
		utils.RespondWithInternalError(c, "Error fetching question types")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, types, "Fetched question types successfully")
}

// ---------- PREGUNTAS ----------

// Crear pregunta
func CreateQuestionHandler(c *gin.Context) {
	var req model.QuestionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	// Obtener la posición máxima actual
	var maxPos int
	database.DB.Model(&model.Question{}).Select("COALESCE(MAX(position), 0)").Scan(&maxPos)

	q := model.Question{
		Question:    req.Question,
		Description: req.Description,
		TypeID:      req.TypeID,
		IsActive:    req.IsActive == nil || *req.IsActive,
		Position:    maxPos + 1, // nueva pregunta al final
	}

	if err := database.DB.Create(&q).Error; err != nil {
		utils.RespondWithInternalError(c, "Error creating question")
		return
	}

	// Crear opciones si existen
	if len(req.Options) > 0 {
		var options []model.QuestionOption
		for _, opt := range req.Options {
			options = append(options, model.QuestionOption{
				QuestionID: q.ID,
				Option:     opt.Option,
			})
		}
		if err := database.DB.Create(&options).Error; err != nil {
			utils.RespondWithInternalError(c, "Error creating question options")
			return
		}
	}

	utils.RespondWithSuccess(c, http.StatusCreated, q, "Question created successfully")
}

// Listar preguntas
func GetQuestionsHandler(c *gin.Context) {
	var questions []model.Question
	err := database.DB.
		Preload("Options").
		Preload("Type").
		Order("position ASC").
		Find(&questions).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching questions")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, questions, "Questions fetched successfully")
}

func GetActiveQuestionsHandler(c *gin.Context) {
	var activeQuestions []model.Question
	if err := database.DB.
		Preload("Options").
		Preload("Type").
		Where("is_active = ?", true).
		Order("position ASC").
		Find(&activeQuestions).Error; err != nil {
		utils.RespondWithInternalError(c, "Error fetching active questions")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, activeQuestions, "Active questions fetched successfully")
}

// Obtener pregunta por ID
func GetQuestionHandler(c *gin.Context) {
	id := c.Param("id")
	var question model.Question
	if err := database.DB.
		Preload("Options").
		Preload("Type").
		First(&question, id).Error; err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusNotFound,
			"Question not found",
			"Something went wrong",
		)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, question, "Question fetched successfully")
}

// Actualizar pregunta
func UpdateQuestionHandler(c *gin.Context) {
	id := c.Param("id")
	var req model.QuestionDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	var question model.Question
	if err := database.DB.Preload("Type").Preload("Options").First(&question, id).Error; err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusNotFound,
			"Question not found",
			"Something went wrong",
		)
		return
	}

	// Usar una transacción para garantizar consistencia
	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.RespondWithInternalError(c, "Error updating question")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Actualizar campos básicos
	question.Question = req.Question
	if req.Description != nil {
		question.Description = req.Description
	}

	question.TypeID = req.TypeID
	var newType model.QuestionType
	if err := tx.First(&newType, question.TypeID).Error; err != nil {
		tx.Rollback()
		utils.RespondWithInternalError(c, "Error updating question")
		return
	}

	question.Type = newType
	if req.IsActive != nil {
		question.IsActive = *req.IsActive
	}

	// Guardar cambios básicos de la pregunta
	if err := tx.Save(&question).Error; err != nil {
		tx.Rollback()
		utils.RespondWithInternalError(c, "Error updating question")
		return
	}

	// Obtener el tipo actual de la pregunta para validación
	var questionType model.QuestionType
	if err := tx.First(&questionType, question.TypeID).Error; err != nil {
		tx.Rollback()
		utils.RespondWithInternalError(c, "Error updating question")
		return
	}

	// Verificar si el tipo requiere opciones
	requiresOptions := questionType.Type == "multiple" || questionType.Type == "dropdown" || questionType.Type == "unique"

	// Manejar opciones si vienen en el request O si el tipo cambió
	if req.Options != nil {
		// Validación: preguntas con opciones deben tener al menos una
		if requiresOptions && len(req.Options) == 0 {
			tx.Rollback()
			utils.RespondWithCustomError(
				c,
				http.StatusBadRequest,
				"Question type must have at least one option",
				"Invalid request",
			)
			return
		}

		// Si el tipo NO requiere opciones pero se enviaron, ignorarlas
		if !requiresOptions {
			// Eliminar opciones existentes ya que el tipo no las necesita
			if err := tx.Where("question_id = ?", question.ID).Delete(&model.QuestionOption{}).Error; err != nil {
				tx.Rollback()
				utils.RespondWithInternalError(c, "Error updating question")
				return
			}

		} else {
			// El tipo SÍ requiere opciones, procesarlas
			// Eliminar TODAS las opciones existentes de esta pregunta
			if err := tx.Where("question_id = ?", question.ID).Delete(&model.QuestionOption{}).Error; err != nil {
				tx.Rollback()
				utils.RespondWithInternalError(c, "Error updating question")
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
					utils.RespondWithInternalError(c, "Error updating question")
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
				utils.RespondWithInternalError(c, "Error updating question")
				return
			}
		}
		// Si el tipo requiere opciones pero no se enviaron, no hacer nada (mantener existentes)
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		utils.RespondWithInternalError(c, "Error updating question")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Question updated")
}

// Reordenar preguntas (intercambiar posiciones)
func ReorderQuestionsHandler(c *gin.Context) {
	var req model.ReorderQuestionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	var source, target model.Question

	// Buscar las preguntas involucradas
	if err := database.DB.First(&source, req.SourceID).Error; err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusNotFound,
			"Origin question not found",
			"Error reordering questions",
		)
		return
	}
	if err := database.DB.First(&target, req.TargetID).Error; err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusNotFound,
			"Target question not found",
			"Error reordering questions",
		)
		return
	}

	// Intercambiar posiciones dentro de una transacción
	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.RespondWithInternalError(c, "Error reordering questions")
		return
	}

	sourcePos := source.Position
	targetPos := target.Position

	if err := tx.Model(&source).Update("position", targetPos).Error; err != nil {
		tx.Rollback()
		utils.RespondWithInternalError(c, "Error reordering questions")
		return
	}

	if err := tx.Model(&target).Update("position", sourcePos).Error; err != nil {
		tx.Rollback()
		utils.RespondWithInternalError(c, "Error reordering questions")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.RespondWithInternalError(c, "Error reordering questions")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, map[string]any{
		"source": source.ID,
		"newPos": targetPos,
		"target": target.ID,
		"oldPos": targetPos,
	},
		"Questions reordered",
	)
}

func ToggleQuestionActiveHandler(c *gin.Context) {
	id := c.Param("id")

	var question model.Question
	if err := database.DB.First(&question, id).Error; err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusNotFound,
			"Question not found",
			"Error updating question",
		)
		return
	}

	question.IsActive = !question.IsActive

	if err := database.DB.Save(&question).Error; err != nil {
		utils.RespondWithInternalError(c, "Error updating question")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Question updated successfully")
}

// Eliminar pregunta
func DeleteQuestionHandler(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&model.Question{}, id).Error; err != nil {
		utils.RespondWithInternalError(c, "Error deleting question")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Question deleted")
}

// ---------- OPCIONES DE PREGUNTA ----------

// Crear opción
func CreateQuestionOptionHandler(c *gin.Context) {
	var req model.QuestionOptionDTO
	questionID, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusBadRequest,
			"Invalid ID",
			"Invalid request format",
		)
		return
	}
	option := model.QuestionOption{QuestionID: uint(questionID), Option: req.Option}
	if err := database.DB.Create(&option).Error; err != nil {
		utils.RespondWithInternalError(c, "Error creating question option")
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, option, "Question option created successfully")
}

// ---------- ENVÍOS DE FORMULARIO ----------

// Crear envio
func CreateSubmissionHandler(c *gin.Context) {
	var req model.CreateSubmissionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	submission := model.Submission{
		SubmittedAt: time.Now(),
		Status:      "pending",
	}

	// Preparar las respuestas
	var answers []model.Answer
	for _, ans := range req.Answers {
		answer := model.Answer{
			QuestionID: ans.QuestionID,
		}

		if ans.Answer != nil {
			answer.Answer = ans.Answer
		}

		if len(ans.OptionsID) > 0 {
			var options []model.QuestionOption
			if err := database.DB.Where("id IN ?", ans.OptionsID).Find(&options).Error; err != nil {
				utils.RespondWithInternalError(c, "Error creating submission")
				return
			}
			answer.Options = options
		}

		answers = append(answers, answer)
	}

	submission.Answers = answers

	if err := database.DB.Create(&submission).Error; err != nil {
		utils.RespondWithInternalError(c, "Error creating submission")
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, submission, "Submission created")
}

// Listar envíos
func GetSubmissionsHandler(c *gin.Context) {
	var submissions []model.Submission

	if err := database.DB.
		Preload("Answers").
		Preload("Answers.Question").
		Preload("Answers.Question.Type").
		Preload("Answers.Question.Options").
		Preload("Answers.Options").
		Find(&submissions).Error; err != nil {
		utils.RespondWithInternalError(c, "Error fetching submissions")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, submissions, "Submissions fetched successfully")
}

// GetSubmissionByID obtiene una submission con sus respuestas, preguntas y opciones asociadas
func GetSubmissionHandler(c *gin.Context) {
	id := c.Param("id")

	var submission model.Submission
	if err := database.DB.
		Preload("Answers").
		Preload("Answers.Question").
		Preload("Answers.Question.Type").
		Preload("Answers.Question.Options").
		Preload("Answers.Options").
		First(&submission, "id = ?", id).Error; err != nil {
		utils.RespondWithInternalError(c, "Error fetching submission")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, submission, "Submission fetched successfully")
}

// Actualizar estado de envío
func UpdateSubmissionStatusHandler(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateSubmissionStatusDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	if err := database.DB.Model(&model.Submission{}).Where("id = ?", id).
		Update("status", req.Status).Error; err != nil {
		utils.RespondWithInternalError(c, "Error updating submission")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Submission updated successfully")
}
