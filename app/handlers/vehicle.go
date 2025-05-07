package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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
		utils.RespondWithError(c,"Insufficient permissions",http.StatusForbidden )
		return
	}

	var vehicle model.Vehicle

	id := c.Param("id")
	if err := database.DB.First(&vehicle, id).Error; err != nil {
		log.Println("Error fetching vehicle:", err)
		utils.RespondWithError(c, "Error getting vehicle", http.StatusInternalServerError)
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
// @Param		user body model.UpdateVehicleRequest true "Updated vehicle information"
// @Success		200	{object} model.ApiResponse "Successfully updated vehicle"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error updating vehicle"
// @Router		/vehicles/{id} [put]
func UpdateVehicle(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)

	if claims.Role != "admin" {
		utils.RespondWithError(c,"Insufficient permissions",http.StatusForbidden )
		return
	}

	var req model.UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request:", err)
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	id := c.Param("id")

	var vehicle model.Vehicle
	if err := database.DB.First(&vehicle, id).Error; err != nil {
		log.Println("Error finding vehicle:", err)
		utils.RespondWithError(c, "Vehicle not found", http.StatusInternalServerError)
		return
	}

	vehicle.Brand = req.Brand
	vehicle.Model = req.Model
	vehicle.LicensePlate = req.LicensePlate
	vehicle.CapacityKg = req.CapacityKg
	vehicle.Available = req.Available
	vehicle.CurrentMileage = req.CurrentMileage
	vehicle.NextMaintenanceMileage = req.NextMaintenanceMileage
	vehicle.LastModifiedAt = time.Now()

	if err := database.DB.Save(&vehicle).Error; err != nil {
		log.Println("Error updating vehicle:", err)
		utils.RespondWithError(c, "Error updating vehicle", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Successfully updated user",
	})
}

// @Summary		Create a new vehicle
// @Description	Creates a new vehicle entry in the database.
// @Tags		vehicles
// @Accept		json
// @Produce		json
// @Param		user body model.CreateVehicleRequest true "Vehicle information to create"
// @Success		200	{object} model.ApiResponse "Successfully created vehile"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error creating new vehicle"
// @Router		/vehicles/ [post]
func CreateVehicle(c *gin.Context) {
	var req model.CreateVehicleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request:", err)
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	vehicle := model.Vehicle{
		Brand:                  req.Brand,
		Model:                  req.Model,
		LicensePlate:           req.LicensePlate,
		CapacityKg:             req.CapacityKg,
		Available:              req.Available,
		CurrentMileage:         req.CurrentMileage,
		NextMaintenanceMileage: req.NextMaintenanceMileage,
	}

	if err := database.DB.Create(&vehicle).Error; err != nil {
		log.Println("Error creating new vehicle:", err)
		utils.RespondWithError(c, "Error creating new vehicle", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Successfully created vehicle",
	})
}
