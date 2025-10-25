package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary		Get all question types available in the system
// @Description	Returns a list of all question types
// @Tags		form
// @Produce		json
// @Success		200	{object} model.ApiResponse "List of question types"
// @Failure		500	{object} model.ApiResponse "Error fetching question types"
// @Router		/form/question-types [get]
func GetQuestionTypesHandler(c *gin.Context) {
	var types []model.QuestionType
	if err := database.DB.Find(&types).Error; err != nil {
		utils.RespondWithInternalError(c, "Error fetching question types")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, types, "Fetched question types successfully")
}

// @Summary		Create a new question for the form
// @Description	Creates a new question to display on the client form
// @Tags		form
// @Produce		json
// @Param       question body model.QuestionDTO true "Question information"
// @Success		200	{object} model.ApiResponse "Returns the data of the newly created question"
// @Failure     400 {object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error creating question"
// @Router		/form/questions [post]
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
		IsRequiered: req.IsRequiered == nil || *req.IsRequiered,
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

// @Summary		Get all created questions
// @Description	Fetches all questions in the form
// @Tags		form
// @Produce		json
// @Param       active query boolean false "question status (is active)"
// @Success		200	{object} model.ApiResponse "Returns all questions with the filter applied"
// @Failure     400 {object} model.ApiResponse "The active filter must be of type boolean"
// @Failure		500	{object} model.ApiResponse "Error fetching questions"
// @Router		/form/questions [get]
func GetQuestionsHandler(c *gin.Context) {
	var questions []model.Question
	var err error

	statusStr := c.Query("active")
	if statusStr == "" {
		err = database.DB.
			Preload("Options").
			Preload("Type").
			Order("position ASC").
			Find(&questions).Error

	} else {
		status, err := strconv.ParseBool(statusStr)
		if err != nil {
			utils.RespondWithCustomError(
				c,
				http.StatusBadRequest,
				"Active filter must be of type boolean",
				"Invalid request format",
			)
		}

		err = database.DB.
			Preload("Options").
			Preload("Type").
			Where("is_active = ?", status).
			Order("position ASC").
			Find(&questions).Error
	}

	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching questions")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, questions, "Questions fetched successfully")
}

// @Summary		Get one form question
// @Description	Fetches the question with the specified ID
// @Tags		form
// @Produce		json
// @Param		id path int true "Question ID"
// @Success		200	{object} model.ApiResponse "Returns the required question"
// @Failure		500	{object} model.ApiResponse "Error fetching question"
// @Router		/form/questions/{id} [get]
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

// @Summary		Updates one question in the system
// @Description	Updates the specified question with the information provided
// @Tags		form
// @Produce		json
// @Param		id path int true "Question ID"
// @Param		question body model.QuestionDTO true "Updated question information"
// @Success		200	{object} model.ApiResponse "Question successfully updated"
// @Failure     400 {object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error updating question"
// @Router		/form/questions/{id} [put]
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

	if req.IsRequiered != nil {
    	question.IsRequiered = *req.IsRequiered
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

// @Summary		Changes the position of two questions
// @Description	Updates the order in which the questions appear in the form
// @Tags		form
// @Produce		json
// @Param		question body model.ReorderQuestionDTO true "Questions information"
// @Success		200	{object} model.ApiResponse "Questions successfully reordered"
// @Failure     400 {object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error reordering questions"
// @Router		/form/questions/reorder [patch]
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

// @Summary		Changes the status of one question
// @Description	Sets the question active or inactive
// @Tags		form
// @Produce		json
// @Param		id path int true "Question ID"
// @Success		200	{object} model.ApiResponse "Question status changed"
// @Failure     400 {object} model.ApiResponse "Invalid request format"
// @Failure     404 {object} model.ApiResponse "Question not found"
// @Failure		500	{object} model.ApiResponse "Error updating question"
// @Router		/form/questions/{id}/active [patch]
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

// @Summary		Changes if a question is requiered o not
// @Description	Sets the question requiered true or false
// @Tags		form
// @Produce		json
// @Param		id path int true "Question ID"
// @Success		200	{object} model.ApiResponse "Question status changed"
// @Failure     400 {object} model.ApiResponse "Invalid request format"
// @Failure     404 {object} model.ApiResponse "Question not found"
// @Failure		500	{object} model.ApiResponse "Error updating question"
// @Router		/form/questions/{id}/requiered [patch]
func ToggleQuestionRequieredHandler(c *gin.Context) {
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

	question.IsRequiered = !question.IsRequiered

	if err := database.DB.Save(&question).Error; err != nil {
		utils.RespondWithInternalError(c, "Error updating question")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Question updated successfully")
}

// @Summary		Deletes one question
// @Description	Soft deletes the specifiec question
// @Tags		form
// @Produce		json
// @Param		id path int true "Question ID"
// @Success		200	{object} model.ApiResponse "Question successfully deleted"
// @Failure		500	{object} model.ApiResponse "Error deleting question"
// @Router		/form/questions/{id} [delete]
func DeleteQuestionHandler(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&model.Question{}, id).Error; err != nil {
		utils.RespondWithInternalError(c, "Error deleting question")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Question deleted")
}

// @Summary		Creates an option for one form question
// @Description	Sets a new possible option to be selected in the client form for one question
// @Tags		form
// @Produce		json
// @Param		id path int true "Question ID"
// @Success		200	{object} model.ApiResponse "Option successfully created"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error creating question option"
// @Router		/form/questions/{id}/options [post]
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
