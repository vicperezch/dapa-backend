package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Get financial report
// @Description	Returns a financial report of all delivered orders
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Financial report"
// @Failure		500	{object} model.ApiResponse "Error retrieving financial report"
// @Router		/reports/financial [get]
func FinancialReport(c *gin.Context) {
	var orders []model.Order
	err := database.DB.Where("status = ?", "delivered").Find(&orders).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching orders")
		return
	}

	var report []model.FinancialReportDTO
	for _, order := range orders {
		var user model.User
		if order.UserID != nil {
			err := database.DB.Where("id = ?", order.UserID).First(&user).Error
			if err != nil {
				utils.RespondWithInternalError(c, "Error fetching user")
				return
			}
		}

		report = append(report, model.FinancialReportDTO{
			Date:        order.Date,
			Type:        order.Type,
			TotalAmount: order.TotalAmount,
			User:        user.Name + " " + user.LastName,
		})
	}

	utils.RespondWithSuccess(c, http.StatusOK, report, "Financial report fetched successfully")
}
