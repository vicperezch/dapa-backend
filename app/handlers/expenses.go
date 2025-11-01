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

// @Summary		Get all expenses
// @Description	Returns a list of all expenses
// @Tags		expenses
// @Produce		json
// @Success		200	{object} model.ApiResponse "List of expenses"
// @Failure		500	{object} model.ApiResponse "Error retrieving expenses"
// @Router		/expenses [get]
func GetExpenses(c *gin.Context) {
	var expenses []model.Expense
	if err := database.DB.Find(&expenses).Error; err != nil {
		utils.RespondWithInternalError(c, "Error fetching expenses")
		return
	}
	utils.RespondWithSuccess(c, http.StatusOK, expenses, "Expenses fetched successfully")
}

// @Summary		Get expense by ID
// @Description	Returns a single expense by its ID
// @Tags		expenses
// @Produce		json
// @Param		id path int true "Expense ID"
// @Success		200	{object} model.ApiResponse "Expense"
// @Failure		404	{object} model.ApiResponse "Expense not found"
// @Router		/expenses/{id} [get]
func GetExpense(c *gin.Context) {
	id := c.Param("id")
	var expense model.Expense
	if err := database.DB.First(&expense, id).Error; err != nil {
		utils.RespondWithCustomError(c, http.StatusNotFound, "Expense not found", "No se encontró el egreso")
		return
	}
	utils.RespondWithSuccess(c, http.StatusOK, expense, "Expense fetched successfully")
}

// @Summary		Create expense
// @Description	Creates a new expense
// @Tags		expenses
// @Accept		json
// @Produce		json
// @Param		expense body model.ExpenseDTO true "Expense"
// @Success		201	{object} model.ApiResponse "Expense created"
// @Failure		400	{object} model.ApiResponse "Invalid request body"
// @Failure		500	{object} model.ApiResponse "Error creating expense"
// @Router		/expenses [post]
func CreateExpense(c *gin.Context) {
	var dto model.ExpenseDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	expense := model.Expense{
		Date:             dto.Date,
		TypeID:           dto.TypeID,
		TemporalEmployee: dto.TemporalEmployee,
		Description:      dto.Description,
		Amount:           dto.Amount,
	}

	if err := database.DB.Create(&expense).Error; err != nil {
		utils.RespondWithInternalError(c, "Error creating expense")
		return
	}
	utils.RespondWithSuccess(c, http.StatusCreated, expense, "Expense created successfully")
}

// @Summary		Update expense
// @Description	Updates an existing expense
// @Tags		expenses
// @Accept		json
// @Produce		json
// @Param		id path int true "Expense ID"
// @Param		expense body model.ExpenseDTO true "Expense"
// @Success		200	{object} model.ApiResponse "Expense updated"
// @Failure		400	{object} model.ApiResponse "Invalid request body"
// @Failure		404	{object} model.ApiResponse "Expense not found"
// @Failure		500	{object} model.ApiResponse "Error updating expense"
// @Router		/expenses/{id} [put]
func UpdateExpense(c *gin.Context) {
	id := c.Param("id")
	var expense model.Expense
	if err := database.DB.First(&expense, id).Error; err != nil {
		utils.RespondWithCustomError(c, http.StatusNotFound, "Expense not found", "No se encontró el egreso")
		return
	}

	var dto model.ExpenseDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	expense.Date = dto.Date
	expense.TypeID = dto.TypeID
	expense.TemporalEmployee = dto.TemporalEmployee
	expense.Description = dto.Description
	expense.Amount = dto.Amount

	if err := database.DB.Save(&expense).Error; err != nil {
		utils.RespondWithInternalError(c, "Error updating expense")
		return
	}
	utils.RespondWithSuccess(c, http.StatusOK, expense, "Expense updated successfully")
}

// @Summary		Delete expense
// @Description	Deletes an expense by its ID
// @Tags		expenses
// @Produce		json
// @Param		id path int true "Expense ID"
// @Success		200	{object} model.ApiResponse "Expense deleted"
// @Failure		404	{object} model.ApiResponse "Expense not found"
// @Failure		500	{object} model.ApiResponse "Error deleting expense"
// @Router		/expenses/{id} [delete]
func DeleteExpense(c *gin.Context) {
	id := c.Param("id")
	var expense model.Expense
	if err := database.DB.First(&expense, id).Error; err != nil {
		utils.RespondWithCustomError(c, http.StatusNotFound, "Expense not found", "No se encontró el egreso")
		return
	}

	if err := database.DB.Delete(&expense).Error; err != nil {
		utils.RespondWithInternalError(c, "Error deleting expense")
		return
	}
	utils.RespondWithSuccess(c, http.StatusOK, nil, "Expense deleted successfully")
}
