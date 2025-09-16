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
func GetVehiclesHandler(c *gin.Context) {
	var vehicles []model.Vehicle
	err := database.DB.
		Where("is_active = ?", true).
		Find(&vehicles).Error

	if err != nil {
		utils.RespondWithInternalError(c, "Error fetching vehicles")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, vehicles, "Vehicles fetched successfully")
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
func GetVehicleHandler(c *gin.Context) {
	id := c.Param("id")
	var vehicle model.Vehicle

	err := database.DB.
		Where("id = ? AND is_active = ?", id, true).
		First(&vehicle).Error

	if err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusNotFound,
			"Vehicle not found",
			"Something went wrong",
		)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, vehicle, "Vehicle fetched successfully")
}

// @Summary		Update vehicle by ID
// @Description	Updates the vehicle's information based on the given ID.
// @Tags		vehicles
// @Accept		json
// @Produce		json
// @Param		id path int true "Vehicle ID"
// @Param		vehicle body model.VehicleDTO true "Updated vehicle information"
// @Success		200	{object} model.ApiResponse "Successfully updated vehicle"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error updating vehicle"
// @Router		/vehicles/{id} [put]
func UpdateVehicleHandler(c *gin.Context) {
	var req model.VehicleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	id := c.Param("id")
	var vehicle model.Vehicle
	if err := database.DB.Where("id = ? AND is_active = ?", id, true).First(&vehicle).Error; err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusNotFound,
			"Vehicle not found",
			"Something went wrong",
		)
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
		utils.RespondWithInternalError(c, "Error updating vehicle")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Vehicle updated successfully")
}

// @Summary		Create a new vehicle
// @Description	Creates a new vehicle entry in the database.
// @Tags		vehicles
// @Accept		json
// @Produce		json
// @Param		vehicle body model.VehicleDTO true "Vehicle information to create"
// @Success		200	{object} model.ApiResponse "Successfully created vehile"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error creating new vehicle"
// @Router		/vehicles/ [post]
func CreateVehicleHandler(c *gin.Context) {
	var req model.VehicleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
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
		utils.RespondWithInternalError(c, "Error creating vehicle")
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, nil, "Vehicle created successfully")
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
func DeleteVehicleHandler(c *gin.Context) {
	id := c.Param("id")

	err := database.DB.Model(&model.Vehicle{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"is_active":  false,
		}).Error

	if err != nil {
		utils.RespondWithInternalError(c, "Error deleting vehicle")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Vehicle deleted successfully")
}
