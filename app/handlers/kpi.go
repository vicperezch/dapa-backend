package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary		Get the performance goal
// @Description	Returns the current performance goal (single row, default zeros if none exists).
// @Tags		kpi
// @Produce		json
// @Success		200	{object} model.PerformanceGoal "Performance goal"
// @Failure		500	{object} model.ApiResponse "Error fetching performance goal"
// @Router		/kpi/goals [get]
func GetPerformanceGoal(c *gin.Context) {
	var goal model.PerformanceGoal
	result := database.DB.First(&goal)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Return a default all-zero goal instead of 500
			defaultGoal := model.PerformanceGoal{
				OrderGoal:           0,
				UtilityGoal:         0.00,
				AveragePerOrderGoal: 0.00,
				TravelGoal:          0,
				DeliveryGoal:        0.00,
				AchievementRateGoal: 0.00,
			}
			utils.RespondWithSuccess(c, http.StatusOK, defaultGoal, "No performance goal found, returning default values")
			return
		}
		utils.RespondWithInternalError(c, "Error fetching performance goal")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, goal, "Performance goal fetched successfully")
}

// @Summary		Update the performance goal
// @Description	Updates the single performance goal (creates if none exists).
// @Tags		kpi
// @Accept		json
// @Produce		json
// @Param		body body model.PerformanceGoal true "Performance goal data"
// @Success		200 {object} model.PerformanceGoal
// @Router		/kpi/goals [post]
func UpsertPerformanceGoal(c *gin.Context) {
	var input model.PerformanceGoal
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithCustomError(c, http.StatusBadRequest, "Invalid request data", "Error actualizando objetivos")
		return
	}

	var goal model.PerformanceGoal
	result := database.DB.First(&goal)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// No existing record → create one
			if err := database.DB.Create(&input).Error; err != nil {
				utils.RespondWithInternalError(c, "Error creating performance goal")
				return
			}
			utils.RespondWithSuccess(c, http.StatusCreated, input, "Performance goal created successfully")
			return
		}
		utils.RespondWithInternalError(c, "Error fetching existing performance goal")
		return
	}

	// Update the existing one
	goal.OrderGoal = input.OrderGoal
	goal.UtilityGoal = input.UtilityGoal
	goal.AveragePerOrderGoal = input.AveragePerOrderGoal
	goal.TravelGoal = input.TravelGoal
	goal.DeliveryGoal = input.DeliveryGoal
	goal.AchievementRateGoal = input.AchievementRateGoal

	if err := database.DB.Save(&goal).Error; err != nil {
		utils.RespondWithInternalError(c, "Error updating performance goal")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, goal, "Performance goal updated successfully")
}

// @Summary		Get current KPIs
// @Description	Returns the current KPIs for the month.
// @Tags		kpi
// @Produce		json
// @Success		200	{object} model.ApiResponse "Current KPIs"
// @Failure		500	{object} model.ApiResponse "Error fetching current KPIs"
// @Router		/kpi/current [get]
func GetCurrentKPIs(c *gin.Context) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// Total Income
	var totalIncome float64
	err := database.DB.Model(&model.Order{}).Where("status = ? AND date BETWEEN ? AND ?", "delivered", startOfMonth, endOfMonth).Select("COALESCE(SUM(total_amount), 0)").Row().Scan(&totalIncome)
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching total income")
		return
	}

	// Total Expenses
	var totalExpenses float64
	err = database.DB.Model(&model.Expense{}).Where("date BETWEEN ? AND ?", startOfMonth, endOfMonth).Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalExpenses)
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching total expenses")
		return
	}

	// Utility
	utility := totalIncome - totalExpenses

	// Average Amount Per Order
	var averageAmount float64
	err = database.DB.Model(&model.Order{}).Where("status = ? AND date BETWEEN ? AND ?", "delivered", startOfMonth, endOfMonth).Select("COALESCE(AVG(total_amount), 0)").Row().Scan(&averageAmount)
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching average amount per order")
		return
	}

	// Completed Trips
	var completedTrips int64
	err = database.DB.Model(&model.Order{}).Where("status = ? AND date BETWEEN ? AND ?", "delivered", startOfMonth, endOfMonth).Count(&completedTrips).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching completed trips")
		return
	}

	// Average Orders Per Employee
	var employeeCount int64
	err = database.DB.Model(&model.Order{}).Where("status = ? AND date BETWEEN ? AND ? AND user_id IS NOT NULL", "delivered", startOfMonth, endOfMonth).Distinct("user_id").Count(&employeeCount).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching employee count")
		return
	}

	var averageOrdersPerEmployee float64
	if employeeCount > 0 {
		averageOrdersPerEmployee = float64(completedTrips) / float64(employeeCount)
	}

	// Fulfillment Rate
	var pendingOrders int64
	err = database.DB.Model(&model.Order{}).Where("status = ? AND date BETWEEN ? AND ?", "pending", startOfMonth, endOfMonth).Count(&pendingOrders).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching pending orders count")
		return
	}

	var fulfillmentRate float64
	if completedTrips+pendingOrders > 0 {
		fulfillmentRate = (float64(completedTrips) / float64(completedTrips+pendingOrders)) * 100
	}

	kpiData := gin.H{
		"utility":                  utility,
		"averagePerOrder":          averageAmount,
		"completedTrips":           completedTrips,
		"deliveredOrders":          completedTrips,
		"averageOrdersPerEmployee": int(averageOrdersPerEmployee + 0.5),
		"fulfillmentRate":          fulfillmentRate,
	}

	utils.RespondWithSuccess(c, http.StatusOK, kpiData, "Current KPIs fetched successfully")
}

