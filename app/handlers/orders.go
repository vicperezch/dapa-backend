package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Get all orders in the system
// @Description	Returns a list of all orders, if the user is a driver it returns only their associated orders
// @Tags		orders
// @Produce		json
// @Param       status query string false "Order status"
// @Success		200	{array} model.Order "List of orders"
// @Failure		500	{object} model.ApiResponse "Error retrieving orders"
// @Router		/orders/ [get]
func GetOrders(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	status := c.Query("status")

	var orders []model.Order
	var err error

	if claims.Role == "admin" {
		if status == "" {
			err = database.DB.Find(&orders).Error

		} else {
			err = database.DB.Where("status = ?", status).Find(&orders).Error
		}

		if err != nil {
			utils.RespondWithInternalError(c, "Error fetching orders")
			return
		}

		utils.RespondWithSuccess(c, http.StatusOK, orders, "Orders fetched successfully")
		return
	}

	err = database.DB.Where("user_id = ?", claims.ID).Find(&orders).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching orders")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, orders, "Orders fetched successfully")
}

// @Summary		Get one order by ID
// @Description	Returns the order object
// @Tags		orders
// @Produce		json
// @Param       id path int true "Order ID"
// @Success		200	{object} model.Order "Order"
// @Failure		403	{object} model.ApiResponse "Insufficient permissions"
// @Failure		500	{object} model.ApiResponse "Error retrieving order"
// @Router		/orders/{id} [get]
func GetOrderById(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	if claims.Role != "admin" {
		utils.RespondWithUnathorizedError(c)
		return
	}

	var order model.Order

	id := c.Param("id")
	err := database.DB.Where("id = ?", id).First(&order).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching order")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, order, "Order fetched successfully")
}

// @Summary		Update one order
// @Description	Updates one order's data
// @Tags		orders
// @Produce		json
// @Param       id path int true "Order ID"
// @Success		200	{object} model.ApiResponse "Order updated successfully"
// @Failure		400 {object} model.ApiResponse "Invalid request format"
// @Failure		403 {object} model.ApiResponse "Insufficient permissions"
// @Failure		404 {object} model.ApiResponse "Order not found"
// @Failure		500 {object} model.ApiResponse "Error updating order"
// @Router		/orders/{id} [put]
func UpdateOrder(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	if claims.Role != "admin" {
		utils.RespondWithUnathorizedError(c)
		return
	}

	var req model.OrderDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	id := c.Param("id")
	var order model.Order
	var err error
	err = database.DB.Where("id = ?", id).First(&order).Error

	if err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusNotFound,
			"Order not found",
			"Something went wrong",
		)
		return
	}

	order.Origin = req.Origin
	order.Destination = req.Destination
	order.TotalAmount = req.TotalAmount
	order.Type = req.Type

	if req.UserID != nil {
		order.UserID = *req.UserID
	}

	if req.VehicleID != nil {
		order.VehicleID = *req.VehicleID
	}

	if req.Details != nil {
		order.Details = req.Details
	}

	err = database.DB.Save(&order).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error updating order")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Order updated successfully")
}

// @Summary		Assigns a driver and vehicle
// @Description	Updates the order's data to include the driver and vehicle ID
// @Tags		orders
// @Produce		json
// @Param       id path int true "Order ID"
// @Success		200	{object} model.ApiResponse "Order assigned successfully"
// @Failure		400 {object} model.ApiResponse "Invalid request format"
// @Failure		403 {object} model.ApiResponse "Insufficient permissions"
// @Failure		404 {object} model.ApiResponse "Order not found"
// @Failure		500 {object} model.ApiResponse "Error assigning order"
// @Router		/orders/{id} [patch]
func AssignOrder(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	if claims.Role != "admin" {
		utils.RespondWithUnathorizedError(c)
		return
	}

	var req model.AssignOrderDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	id := c.Param("id")
	var order model.Order
	var err error

	err = database.DB.Where("id = ?", id).First(&order).Error
	if err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusNotFound,
			"Order not found",
			"Something went wrong",
		)
		return
	}

	order.UserID = req.UserID
	order.VehicleID = req.VehicleID

	err = database.DB.Save(&order).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error assigning order")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Order assigned successfully")
}
