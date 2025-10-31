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

	if len(orders) == 0 {
		utils.RespondWithSuccess(c, http.StatusOK, []model.FinancialReportDTO{}, "No financial data available for the given date range")
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
	err := database.DB.
			Model(&model.Order{}).
			Where("status = ?", "delivered").
			Select("COALESCE(SUM(total_amount), 0)").
			Row().
			Scan(&totalIncome)
	if err != nil {
			utils.RespondWithInternalError(c, "Error fetching total income")
			return
	}

	report := model.TotalIncomeReportDTO{
		TotalIncome: totalIncome,
	}

	utils.RespondWithSuccess(c, http.StatusOK, report, "Total income report fetched successfully")
}

// @Summary		Get completed quotations chart
// @Description	Returns data for the completed quotations chart
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Completed quotations chart data"
// @Failure		500	{object} model.ApiResponse "Error retrieving completed quotations chart data"
// @Router		/reports/completed-quotations [get]
func CompletedQuotationsChart(c *gin.Context) {
	var results []struct {
		Month string
		Count int
	}
	err := database.DB.Model(&model.Submission{}).
		Select("TO_CHAR(submitted_at, 'Month YYYY') as month, COUNT(*) as count").
		Where("status = ?", model.FormStatusApproved).
		Group("month").
		Order("month").
		Scan(&results).Error

	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching completed quotations")
		return
	}

	var categories []string
	var data []int
	for _, result := range results {
		categories = append(categories, result.Month)
		data = append(data, result.Count)
	}

	chartData := model.CompletedQuotationsDTO{
		Series: []struct {
			Data []int  `json:"data"`
		}{
			{Data: data},
		},
		Categories: categories,
	}

	utils.RespondWithSuccess(c, http.StatusOK, chartData, "Completed quotations chart data fetched successfully")
}

// @Summary		Get quotations status chart
// @Description	Returns data for the quotations status chart
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Quotations status chart data"
// @Failure		500	{object} model.ApiResponse "Error retrieving quotations status chart data"
// @Router		/reports/quotations-status [get]
func QuotationsStatusChart(c *gin.Context) {
	var results []struct {
		Status string
		Count  int
	}
	err := database.DB.Model(&model.Submission{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching quotations status")
		return
	}

	var series []float64
	var labels []string
	for _, result := range results {
		series = append(series, float64(result.Count))
		labels = append(labels, result.Status)
	}

	chartData := model.QuotationsStatusDTO{
		Series: series,
		Labels: labels,
	}

	utils.RespondWithSuccess(c, http.StatusOK, chartData, "Quotations status chart data fetched successfully")
}

// @Summary		Get drivers performance chart
// @Description	Returns data for the drivers performance chart
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Drivers performance chart data"
// @Failure		500	{object} model.ApiResponse "Error retrieving drivers performance chart data"
// @Router		/reports/drivers-performance [get]
func DriversPerformanceChart(c *gin.Context) {
	var drivers []model.User
	err := database.DB.Where("role = ?", "driver").Find(&drivers).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching drivers")
		return
	}

	var categories []string
	var deliveredData []int
	var pendingData []int

	for _, driver := range drivers {
		categories = append(categories, driver.Name+" "+driver.LastName)

		var deliveredCount int64
		database.DB.Model(&model.Order{}).Where("user_id = ? AND status = ?", driver.ID, "delivered").Count(&deliveredCount)
		deliveredData = append(deliveredData, int(deliveredCount))

		var pendingCount int64
		database.DB.Model(&model.Order{}).Where("user_id = ? AND status = ?", driver.ID, "pending").Count(&pendingCount)
		pendingData = append(pendingData, int(pendingCount))
	}

	chartData := model.DriversBarDataDTO{
		Series: []struct {
			Name string `json:"name"`
			Data []int  `json:"data"`
		}{
			{Name: "delivered", Data: deliveredData},
			{Name: "pending", Data: pendingData},
		},
		Categories: categories,
	}

	utils.RespondWithSuccess(c, http.StatusOK, chartData, "Drivers performance chart data fetched successfully")
}

// @Summary		Get drivers trip participation chart
// @Description	Returns data for the drivers trip participation chart
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Drivers trip participation chart data"
// @Failure		500	{object} model.ApiResponse "Error retrieving drivers trip participation chart data"
// @Router		/reports/drivers-trip-participation [get]
func DriversTripParticipationChart(c *gin.Context) {
	var results []struct {
		Driver string
		Count  int
	}
	err := database.DB.Model(&model.Order{}).
		Select("users.name || ' ' || users.last_name as driver, COUNT(*) as count").
		Joins("join users on users.id = orders.user_id").
		Where("orders.status = ?", "delivered").
		Group("driver").
		Scan(&results).Error

	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching drivers trip participation")
		return
	}

	var series []int
	var labels []string
	for _, result := range results {
		series = append(series, result.Count)
		labels = append(labels, result.Driver)
	}

	chartData := model.TripParticipationDTO{
		Series: series,
		Labels: labels,
	}

	utils.RespondWithSuccess(c, http.StatusOK, chartData, "Drivers trip participation chart data fetched successfully")
}

func FinancialControlIncome(c *gin.Context) {
	// Opcional: filtrar por rango de fechas usando query params startDate y endDate (YYYY-MM-DD)
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate time.Time
	var err error
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid start date format")
			return
		}
	}
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid end date format")
			return
		}
	}

	// Resultado intermedio
	var results []struct {
		Date          time.Time `json:"date"`
		InType        string    `json:"in_type"`
		Amount        float64   `json:"amount"`
		PaymentMethod string    `json:"payment_method"`
		Assigned      string    `json:"assigned"`
		Description   string    `json:"description"`
	}

	q := database.DB.Model(&model.Order{}).
		Select("orders.date, orders.type as in_type, orders.total_amount as amount, orders.details as description, COALESCE(users.name || ' ' || users.last_name, 'Sin responsable') as assigned, users.id as user_id").
		Joins("left join users on users.id = orders.user_id").
		Where("orders.status = ?", "delivered")

	if !startDate.IsZero() && !endDate.IsZero() {
		q = q.Where("orders.date BETWEEN ? AND ?", startDate, endDate)
	} else if !startDate.IsZero() {
		q = q.Where("orders.date >= ?", startDate)
	} else if !endDate.IsZero() {
		q = q.Where("orders.date <= ?", endDate)
	}

	if err := q.Order("orders.date desc").Scan(&results).Error; err != nil {
		utils.RespondWithInternalError(c, "Error fetching income records")
		return
	}

	// PaymentMethod no está en la entidad Order actualmente, dejar como vacío o N/A
	for i := range results {
		if results[i].PaymentMethod == "" {
			results[i].PaymentMethod = "N/A"
		}
	}

	utils.RespondWithSuccess(c, http.StatusOK, results, "Financial income control fetched successfully")
}

func FinancialControlSpending(c *gin.Context) {
	// Opcional: filtrar por rango de fechas usando query params startDate y endDate (YYYY-MM-DD)
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate time.Time
	var err error
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid start date format")
			return
		}
	}
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid end date format")
			return
		}
	}

	var results []struct {
		Date          time.Time `json:"date"`
		ExType        string    `json:"ex_type"`
		Assigned      string    `json:"assigned"`
		Description   string    `json:"description"`
		PaymentMethod string    `json:"payment_method"`
		Amount        float64   `json:"amount"`
	}

	q := database.DB.Model(&model.Expense{}).
		Select("expenses.date, expense_types.type as ex_type, expenses.temporal_employee, expenses.description, expenses.amount").
		Joins("join expense_types on expense_types.id = expenses.type_id")

	if !startDate.IsZero() && !endDate.IsZero() {
		q = q.Where("expenses.date BETWEEN ? AND ?", startDate, endDate)
	} else if !startDate.IsZero() {
		q = q.Where("expenses.date >= ?", startDate)
	} else if !endDate.IsZero() {
		q = q.Where("expenses.date <= ?", endDate)
	}

	// Escanear en una estructura temporal que incluya temporal_employee para mapear Assigned
	var tmpResults []struct {
		Date            time.Time `json:"date"`
		ExType          string    `json:"ex_type"`
		TemporalEmployee bool     `json:"temporal_employee"`
		Description     string    `json:"description"`
		Amount          float64   `json:"amount"`
	}

	if err := q.Order("expenses.date desc").Scan(&tmpResults).Error; err != nil {
		utils.RespondWithInternalError(c, "Error fetching spending records")
		return
	}

	// Mapear tmpResults a results y convertir temporal_employee a Assigned
	for _, r := range tmpResults {
		assigned := ""
		if r.TemporalEmployee {
			assigned = "temporal"
		}
		results = append(results, struct {
			Date          time.Time `json:"date"`
			ExType        string    `json:"ex_type"`
			Assigned      string    `json:"assigned"`
			Description   string    `json:"description"`
			PaymentMethod string    `json:"payment_method"`
			Amount        float64   `json:"amount"`
		}{
			Date:          r.Date,
			ExType:        r.ExType,
			Assigned:      assigned,
			Description:   r.Description,
			PaymentMethod: "N/A",
			Amount:        r.Amount,
		})
	}

	utils.RespondWithSuccess(c, http.StatusOK, results, "Financial spending control fetched successfully")
}

// @Summary		Get income per month
// @Description	Returns the total income for each month
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Income per month"
// @Failure		500	{object} model.ApiResponse "Error retrieving income per month"
// @Router		/reports/income-per-month [get]
func IncomePerMonth(c *gin.Context) {
	var results []struct {
		Month  string
		Amount float64
	}
	err := database.DB.Model(&model.Order{}).
		Select("TO_CHAR(date, 'Month YYYY') as month, SUM(total_amount) as amount").
		Where("status = ?", "delivered").
		Group("month").
		Order("month").
		Scan(&results).Error

	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching income per month")
		return
	}

	if len(results) == 0 {
		utils.RespondWithSuccess(c, http.StatusOK, model.IncomePerMonthDTO{}, "No income data available")
		return
	}

	var categories []string
	var data []float64
	for _, result := range results {
		categories = append(categories, result.Month)
		data = append(data, result.Amount)
	}

	chartData := model.IncomePerMonthDTO{
		Series: []struct {
			Data []float64 `json:"data"`
		}{
			{Data: data},
		},
		Categories: categories,
	}

	utils.RespondWithSuccess(c, http.StatusOK, chartData, "Income per month fetched successfully")
}

// @Summary		Get expenses per type
// @Description	Returns the total expenses for each expense type
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Expenses per type"
// @Failure		500	{object} model.ApiResponse "Error retrieving expenses per type"
// @Router		/reports/expenses-per-type [get]
func ExpensesPerType(c *gin.Context) {
	var results []struct {
		Type   string
		Amount float64
	}
	err := database.DB.Model(&model.Expense{}).
		Select("expense_types.type, SUM(expenses.amount) as amount").
		Joins("join expense_types on expense_types.id = expenses.type_id").
		Group("expense_types.type").
		Scan(&results).Error

	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching expenses per type")
		return
	}

	if len(results) == 0 {
		utils.RespondWithSuccess(c, http.StatusOK, model.ExpensesPerTypeDTO{}, "No expense data available")
		return
	}

	var series []float64
	var labels []string
	for _, result := range results {
		series = append(series, result.Amount)
		labels = append(labels, result.Type)
	}

	chartData := model.ExpensesPerTypeDTO{
		Series: series,
		Labels: labels,
	}

	utils.RespondWithSuccess(c, http.StatusOK, chartData, "Expenses per type fetched successfully")
}

// @Summary		Get expenses per month
// @Description	Returns the total expenses for each month
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Expenses per month"
// @Failure		500	{object} model.ApiResponse "Error retrieving expenses per month"
// @Router		/reports/expenses-per-month [get]
func ExpensesPerMonth(c *gin.Context) {
	var results []struct {
		Month  string
		Amount float64
	}
	err := database.DB.Model(&model.Expense{}).
		Select("TO_CHAR(date, 'Month YYYY') as month, SUM(amount) as amount").
		Group("month").
		Order("month").
		Scan(&results).Error

	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching expenses per month")
		return
	}

	if len(results) == 0 {
		utils.RespondWithSuccess(c, http.StatusOK, model.ExpensesPerMonthDTO{}, "No expense data available")
		return
	}

	var categories []string
	var data []float64
	for _, result := range results {
		categories = append(categories, result.Month)
		data = append(data, result.Amount)
	}

	chartData := model.ExpensesPerMonthDTO{
		Series: []struct {
			Data []float64 `json:"data"`
		}{
			{Data: data},
		},
		Categories: categories,
	}

	utils.RespondWithSuccess(c, http.StatusOK, chartData, "Expenses per month fetched successfully")
}

// @Summary		Get order type distribution
// @Description	Returns the distribution of order types
// @Tags		reports
// @Produce		json
// @Success		200	{object} model.ApiResponse "Order type distribution"
// @Failure		500	{object} model.ApiResponse "Error retrieving order type distribution"
// @Router		/reports/order-type-distribution [get]
func OrderTypeDistribution(c *gin.Context) {
	var results []struct {
		OrderType string `json:"order_type"`
		Count     int    `json:"count"`
	}
	err := database.DB.Model(&model.Order{}).
		Select("type as order_type, COUNT(*) as count").
		Where("status = ?", "delivered").
		Group("type").
		Scan(&results).Error

	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching order type distribution")
		return
	}

	if len(results) == 0 {
		utils.RespondWithSuccess(c, http.StatusOK, model.OrderTypeDistributionDTO{}, "No order type data available")
		return
	}

	var series []int
	var categories []string
	for _, result := range results {
		series = append(series, result.Count)
		categories = append(categories, result.OrderType)
	}

	chartData := model.OrderTypeDistributionDTO{
		Series: series,
		Categories: categories,
	}

	utils.RespondWithSuccess(c, http.StatusOK, chartData, "Order type distribution fetched successfully")
}
