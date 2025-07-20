package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Get all quotes
// @Description	Returns a list of all quotes in the system.
// @Tags		quotes
// @Produce		json
// @Success		200	{array} model.Quote "List of quotes"
// @Failure		403	{object} model.ApiResponse "Insufficient permissions"
// @Failure		500	{object} model.ApiResponse "Error fetching quotes"
// @Router		/quotes/ [get]
func GetQuotes(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)

	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var quotes []model.QuoteView

	if err := database.DB.Find(&quotes).Error; err != nil {
		log.Println("Error fetching quotes:", err)
		utils.RespondWithError(c, "Error getting all quotes", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, quotes)
}

// @Summary		Assign a driver and vehicle to an order
// @Description	Adds the driver and vehicle id to the quotes table.
// @Tags		quotes
// @Accept		json
// @Produce		json
// @Param		user body model.AssignQuote true "Information to add"
// @Success		200	{object} model.ApiResponse "Successfully assigned"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error assigning driver and vehicle"
// @Router		/quotes/{id} [post]
func AssignQuoteInfo(c *gin.Context) {
	id := c.Param("id")
	var req model.AssignQuote

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Binding error: ", err)
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	var quote model.Quote

	if err := database.DB.First(&quote, id).Error; err != nil {
		log.Println("Error fetching quote: ", err)
		utils.RespondWithError(c, "Error assigning information", http.StatusInternalServerError)
	}

	quote.DriverID = req.Driver
	quote.VehicleID = req.Vehicle
	quote.Details = req.Details

	if err := database.DB.Save(&quote).Error; err != nil {
		log.Println("Error updating quote: ", err)
		utils.RespondWithError(c, "Error assigning information", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Quote updated successfully",
	})
}
