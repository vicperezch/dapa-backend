package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary		Creates a form submission
// @Description	Creates a new set of responses for the client form
// @Tags		form
// @Produce		json
// @Success		200	{object} model.ApiResponse "Submission successfully created"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error creating submission"
// @Router		/form/submissions [post]
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

// @Summary		Gets all form submissions
// @Description	Fetches all the created submissions
// @Tags		form
// @Produce		json
// @Success		200	{object} model.ApiResponse "Submissions fetched successfully"
// @Failure		500	{object} model.ApiResponse "Error fetching submissions"
// @Router		/form/submissions [get]
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

// @Summary		Gets a form submission
// @Description	Fetches the submission with the specified ID
// @Tags		form
// @Produce		json
// @Param		id path int true "Submission ID"
// @Success		200	{object} model.ApiResponse "Submission fetched successfully"
// @Failure		500	{object} model.ApiResponse "Error fetching submission"
// @Router		/form/submissions/{id} [get]
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

// @Summary		Updates a submission status
// @Description	Changes the submission status to the specifiec one
// @Tags		form
// @Produce		json
// @Param		id path int true "Submission ID"
// @Success		200	{object} model.ApiResponse "Submission updated successfully"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error updating submission"
// @Router		/form/submissions/{id}/status [patch]
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
