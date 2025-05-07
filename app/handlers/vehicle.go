package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"log"
	"net/http"

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
// @Router		/users/{id} [get]
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
