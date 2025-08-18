package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary		Get all vehicles
// @Description	Returns a list of all vehicles in the system.
// @Tags		vehicles
// @Produce		json
// @Success		200	{array} model.Vehicle "List of vehicles"
// @Failure		403	{object} model.ApiResponse "Insufficient permissions"
// @Failure		500	{object} model.ApiResponse "Error fetching vehicles"
// @Router		/vehicles/ [get]
func GetVehicles(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var vehicles []model.Vehicle
	err := database.DB.
		Where("is_active = ?", true).
		Find(&vehicles).Error

	if err != nil {
		utils.RespondWithError(c, "Error retrieving vehicles", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, vehicles)
}

// @Summary		Get vehicle by ID
// @Description	Returns the vehicle information based on the given ID.
// @Tags		vehicles
// @Produce		json
// @Param		id path int true "Vehicle ID"
// @Success		200	{object} model.Vehicle "Vehicle found"
// @Failure		403	{object} model.ApiResponse "Insufficient permissions"
// @Failure		500	{object} model.ApiResponse "Error fetching vehicle"
// @Router		/vehicles/{id} [get]
func GetVehicleById(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	id := c.Param("id")
	var vehicle model.Vehicle

	err := database.DB.
		Where("id = ? AND is_active = ?", id, true).
		First(&vehicle).Error

	if err != nil {
		utils.RespondWithError(c, "Vehicle not found", http.StatusNotFound)
		return
	}

	utils.RespondWithJSON(c, vehicle)
}

// @Summary		Update vehicle by ID
// @Description	Updates the vehicle's information based on the given ID.
// @Tags		vehicles
// @Accept		json
// @Produce		json
// @Param		id path int true "Vehicle ID"
// @Param		vehicle body model.UpdateVehicleRequest true "Updated vehicle information"
// @Success		200	{object} model.ApiResponse "Successfully updated vehicle"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error updating vehicle"
// @Router		/vehicles/{id} [put]
func UpdateVehicle(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var req model.VehicleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	id := c.Param("id")
	var vehicle model.Vehicle
	if err := database.DB.Where("id = ? AND is_active = ?", id, true).First(&vehicle).Error; err != nil {
		utils.RespondWithError(c, "Vehicle not found", http.StatusNotFound)
		return
	}

	updated := model.Vehicle{
		ID:             vehicle.ID,
		Brand:          req.Brand,
		Model:          req.Model,
		LicensePlate:   req.LicensePlate,
		CapacityKg:     req.CapacityKg,
		IsAvailable:    req.IsAvailable,
		InsuranceDate:  req.InsuranceDate,
		IsActive:       true,
		CreatedAt:      vehicle.CreatedAt,
		LastModifiedAt: time.Now(),
	}

	if err := database.DB.Save(&updated).Error; err != nil {
		utils.RespondWithError(c, "Error updating vehicle", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Vehicle updated successfully",
	})
}

// @Summary		Create a new vehicle
// @Description	Creates a new vehicle entry in the database.
// @Tags		vehicles
// @Accept		json
// @Produce		json
// @Param		vehicle body model.CreateVehicleRequest true "Vehicle information to create"
// @Success		200	{object} model.ApiResponse "Successfully created vehile"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error creating new vehicle"
// @Router		/vehicles/ [post]
func CreateVehicle(c *gin.Context) {
	var req model.VehicleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	vehicle := model.Vehicle{
		Brand:         req.Brand,
		Model:         req.Model,
		LicensePlate:  req.LicensePlate,
		CapacityKg:    req.CapacityKg,
		IsAvailable:   req.IsAvailable,
		InsuranceDate: req.InsuranceDate,
	}

	err := database.DB.Create(&vehicle).Error
	if err != nil {
		utils.RespondWithError(c, "Error creating vehicle", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Vehicle created successfully",
	})
}

// @Summary		Mark vehicle as inactive
// @Description	Marks the vehicle as inactive instead of permanently deleting.
// @Tags		vehicles
// @Produce		json
// @Param		id path int true "Vehicle ID"
// @Success		200	{object} model.ApiResponse "Successfully marked vehicle as inactive"
// @Failure		403	{object} model.ApiResponse "Insufficient permissions"
// @Failure		500	{object} model.ApiResponse "Error updating vehicle status"
// @Router		/vehicles/{id} [delete]
func DeleteVehicle(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	id := c.Param("id")

	err := database.DB.Model(&model.Vehicle{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"is_active":  false,
		}).Error

	if err != nil {
		utils.RespondWithError(c, "Error deleting vehicle", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Vehicle successfully deleted",
	})
}
