package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Get all expense types
// @Description	Returns a list of all expense types
// @Tags		expense-types
// @Produce		json
// @Success		200	{object} model.ApiResponse "List of expense types"
// @Failure		500	{object} model.ApiResponse "Error retrieving expense types"
// @Router		/expense-types [get]
func GetExpenseTypes(c *gin.Context) {
	var expenseTypes []model.ExpenseType
	if err := database.DB.Find(&expenseTypes).Error; err != nil {
		utils.RespondWithInternalError(c, "Error fetching expense types")
		return
	}
	utils.RespondWithSuccess(c, http.StatusOK, expenseTypes, "Expense types fetched successfully")
}

// @Summary		Get expense type by ID
// @Description	Returns a single expense type by its ID
// @Tags		expense-types
// @Produce		json
// @Param		id path int true "Expense Type ID"
// @Success		200	{object} model.ApiResponse "Expense type"
// @Failure		404	{object} model.ApiResponse "Expense type not found"
// @Router		/expense-types/{id} [get]
func GetExpenseType(c *gin.Context) {
	id := c.Param("id")
	var expenseType model.ExpenseType
	if err := database.DB.First(&expenseType, id).Error; err != nil {
		utils.RespondWithCustomError(c, http.StatusNotFound, "Expense type not found", "No se encontró el tipo de egreso")
		return
	}
	utils.RespondWithSuccess(c, http.StatusOK, expenseType, "Expense type fetched successfully")
}

// @Summary		Create expense type
// @Description	Creates a new expense type
// @Tags		expense-types
// @Accept		json
// @Produce		json
// @Param		expenseType body model.ExpenseTypeDTO true "Expense Type"
// @Success		201	{object} model.ApiResponse "Expense type created"
// @Failure		400	{object} model.ApiResponse "Invalid request body"
// @Failure		500	{object} model.ApiResponse "Error creating expense type"
// @Router		/expense-types [post]
func CreateExpenseType(c *gin.Context) {
	var dto model.ExpenseTypeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	expenseType := model.ExpenseType{
		Type: dto.Type,
	}

	if err := database.DB.Create(&expenseType).Error; err != nil {
		utils.RespondWithInternalError(c, "Error creating expense type")
		return
	}
	utils.RespondWithSuccess(c, http.StatusCreated, expenseType, "Expense type created successfully")
}

// @Summary		Update expense type
// @Description	Updates an existing expense type
// @Tags		expense-types
// @Accept		json
// @Produce		json
// @Param		id path int true "Expense Type ID"
// @Param		expenseType body model.ExpenseTypeDTO true "Expense Type"
// @Success		200	{object} model.ApiResponse "Expense type updated"
// @Failure		400	{object} model.ApiResponse "Invalid request body"
// @Failure		404	{object} model.ApiResponse "Expense type not found"
// @Failure		500	{object} model.ApiResponse "Error updating expense type"
// @Router		/expense-types/{id} [put]
func UpdateExpenseType(c *gin.Context) {
	id := c.Param("id")
	var expenseType model.ExpenseType
	if err := database.DB.First(&expenseType, id).Error; err != nil {
		utils.RespondWithCustomError(c, http.StatusNotFound, "Expense type not found", "No se encontró el tipo de egreso")
		return
	}

	var dto model.ExpenseTypeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	expenseType.Type = dto.Type
	if err := database.DB.Save(&expenseType).Error; err != nil {
		utils.RespondWithInternalError(c, "Error updating expense type")
		return
	}
	utils.RespondWithSuccess(c, http.StatusOK, expenseType, "Expense type updated successfully")
}

// @Summary		Delete expense type
// @Description	Deletes an expense type by its ID
// @Tags		expense-types
// @Produce		json
// @Param		id path int true "Expense Type ID"
// @Success		200	{object} model.ApiResponse "Expense type deleted"
// @Failure		404	{object} model.ApiResponse "Expense type not found"
// @Failure		500	{object} model.ApiResponse "Error deleting expense type"
// @Router		/expense-types/{id} [delete]
func DeleteExpenseType(c *gin.Context) {
	id := c.Param("id")
	var expenseType model.ExpenseType
	if err := database.DB.First(&expenseType, id).Error; err != nil {
		utils.RespondWithCustomError(c, http.StatusNotFound, "Expense type not found", "No se encontró el tipo de egreso")
		return
	}

	if err := database.DB.Delete(&expenseType).Error; err != nil {
		utils.RespondWithInternalError(c, "Error deleting expense type")
		return
	}
	utils.RespondWithSuccess(c, http.StatusOK, nil, "Expense type deleted successfully")
}
