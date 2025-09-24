package handlers

import (
	"context"
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary		Creates and order from a submission
// @Description	Changes a submission status and creates a new order with the information provided
// @Tags		orders
// @Produce		json
// @Param		order body model.AcceptSubmissionDTO true "Order information"
// @Success		200	{object} model.ApiResponse "Order successfully created"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error creating order"
// @Router		/orders/ [post]
func CreateOrderHandler(c *gin.Context) {
	var req model.AcceptSubmissionDTO
	var err error

	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	ctx := context.Background()
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		_, txErr := gorm.G[model.Submission](tx).Where("id = ?", req.SubmissionID).Update(ctx, "status", "approved")
		if txErr != nil {
			return txErr
		}

		currentDate := time.Now().Truncate(24 * time.Hour)

		order := model.Order{
			SubmissionID: req.SubmissionID,
			UserID:       nil,
			VehicleID:    nil,
			ClientName:   req.ClientName,
			ClientPhone:  req.ClientPhone,
			Origin:       req.Origin,
			Destination:  req.Destination,
			TotalAmount:  req.TotalAmount,
			Details:      req.Details,
			Status:       "pending",
			Type:         req.Type,
			Date:         currentDate,
		}
		txErr = gorm.G[model.Order](tx).Create(ctx, &order)
		if txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		utils.RespondWithInternalError(c, "Error creating order")
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, nil, "Order successfully created")
}

// @Summary		Get all orders in the system
// @Description	Returns a list of all orders, if the user is a driver it returns only their associated orders
// @Tags		orders
// @Produce		json
// @Param       status query string false "Order status"
// @Success		200	{object} model.ApiResponse "List of orders"
// @Failure		500	{object} model.ApiResponse "Error retrieving orders"
// @Router		/orders/ [get]
func GetOrdersHandler(c *gin.Context) {
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

	err = database.DB.Where("user_id = ?", claims.UserID).Find(&orders).Error
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
// @Success		200	{object} model.ApiResponse "Order"
// @Failure		403	{object} model.ApiResponse "Insufficient permissions"
// @Failure		500	{object} model.ApiResponse "Error retrieving order"
// @Router		/orders/{id} [get]
func GetOrderHandler(c *gin.Context) {
	var order model.Order
	claims := c.MustGet("claims").(*model.EmployeeClaims)

	id := c.Param("id")
	err := database.DB.Where("id = ?", id).First(&order).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching order")
		return
	}

	if claims.Role == "driver" && (order.UserID == nil || *order.UserID != claims.UserID) {
		utils.RespondWithUnathorizedError(c)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, order, "Order fetched successfully")
}

// @Summary		Update one order
// @Description	Updates one order's data
// @Tags		orders
// @Produce		json
// @Param       id path int true "Order ID"
// @Param		order body model.OrderDTO true "Updated order information"
// @Success		200	{object} model.ApiResponse "Order updated successfully"
// @Failure		400 {object} model.ApiResponse "Invalid request format"
// @Failure		403 {object} model.ApiResponse "Insufficient permissions"
// @Failure		404 {object} model.ApiResponse "Order not found"
// @Failure		500 {object} model.ApiResponse "Error updating order"
// @Router		/orders/{id} [put]
func UpdateOrderHandler(c *gin.Context) {
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

	order.ClientName = req.ClientName
	order.ClientPhone = req.ClientPhone
	order.Origin = req.Origin
	order.Destination = req.Destination
	order.TotalAmount = req.TotalAmount
	order.Type = req.Type

	if req.UserID != nil {
		order.UserID = req.UserID
	}

	if req.VehicleID != nil {
		order.VehicleID = req.VehicleID
	}

	if req.Details != nil {
		order.Details = *req.Details
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
// @Router		/orders/{id}/assign [patch]
func AssignOrderHandler(c *gin.Context) {
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

	order.UserID = &req.UserID
	order.VehicleID = &req.VehicleID
	order.Status = "assigned"

	err = database.DB.Save(&order).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error assigning order")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Order assigned successfully")
}

// @Summary		Changes and order's status
// @Description	Updates the order's data to the provided status
// @Tags		orders
// @Produce		json
// @Param       id path int true "Order ID"
// @Param		status body model.OrderStatusDTO true "Updated order status"
// @Success		200	{object} model.ApiResponse "Order status updated successfully"
// @Failure		400 {object} model.ApiResponse "Invalid request format"
// @Failure		500 {object} model.ApiResponse "Error updating order status"
// @Router		/orders/{id}/status [patch]
func ChangeOrderStatusHandler(c *gin.Context) {
	var req model.OrderStatusDTO
	var err error

	id := c.Param("id")
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	ctx := context.Background()
	_, err = gorm.G[model.Order](database.DB).Where("id = ?", id).Update(ctx, "status", req.Status)
	if err != nil {
		utils.RespondWithInternalError(c, "Error updating order status")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Order status updated successfully")
}