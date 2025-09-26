package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"
	"time"

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

// @Summary		Get financial report by date range
// @Description	Returns a financial report of all delivered orders within a specific date range
// @Tags		reports
// @Produce		json
// @Param		startDate query string true "Start date for the report"
// @Param		endDate query string true "End date for the report"
// @Success		200	{object} model.ApiResponse "Financial report"
// @Failure		400	{object} model.ApiResponse "Invalid date format"
// @Failure		500	{object} model.ApiResponse "Error retrieving financial report"
// @Router		/reports/financial/date [get]
func FinancialReportByDate(c *gin.Context) {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid start date format")
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid end date format")
		return
	}

	var orders []model.Order
	err = database.DB.Where("status = ? AND date BETWEEN ? AND ?", "delivered", startDate, endDate).Find(&orders).Error
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

// @Summary		Get drivers report
// @Description	Returns a report of drivers with delivered orders
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Drivers report"
// @Failure		500	{object} model.ApiResponse "Error retrieving drivers report"
// @Router		/reports/drivers [get]
func DriversReport(c *gin.Context) {
	var orders []model.Order
	err := database.DB.Where("status = ?", "delivered").Find(&orders).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching orders")
		return
	}

	driverStats := make(map[uint]struct {
		totalOrders int
		firstOrder  time.Time
		lastOrder   time.Time
	})

	for _, order := range orders {
		if order.UserID != nil {
			stats := driverStats[*order.UserID]
			stats.totalOrders++
			if stats.firstOrder.IsZero() || order.Date.Before(stats.firstOrder) {
				stats.firstOrder = order.Date
			}
			if stats.lastOrder.IsZero() || order.Date.After(stats.lastOrder) {
				stats.lastOrder = order.Date
			}
			driverStats[*order.UserID] = stats
		}
	}

	var report []model.DriverReportDTO
	for driverID, stats := range driverStats {
		var driver model.User
		err := database.DB.Where("id = ?", driverID).First(&driver).Error
		if err != nil {
			utils.RespondWithInternalError(c, "Error fetching driver")
			return
		}

		weeks := stats.lastOrder.Sub(stats.firstOrder).Hours() / (24 * 7)
		if weeks < 1 {
			weeks = 1
		}
		ordersPerWeek := float64(stats.totalOrders) / weeks

		report = append(report, model.DriverReportDTO{
			DriverName:      driver.Name + " " + driver.LastName,
			TotalOrders:     stats.totalOrders,
			OrdersPerWeek: ordersPerWeek,
		})
	}

	utils.RespondWithSuccess(c, http.StatusOK, report, "Drivers report fetched successfully")
}

// @Summary		Get total income report
// @Description	Returns the total income from all delivered orders
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Total income report"
// @Failure		500	{object} model.ApiResponse "Error retrieving total income report"
// @Router		/reports/income [get]
func TotalIncomeReport(c *gin.Context) {
	var totalIncome float64
	err := database.DB.Model(&model.Order{}).Where("status = ?", "delivered").Select("sum(total_amount)").Row().Scan(&totalIncome)
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching total income")
		return
	}

	report := model.TotalIncomeReportDTO{
		TotalIncome: totalIncome,
	}

	utils.RespondWithSuccess(c, http.StatusOK, report, "Total income report fetched successfully")
}
