package handlers

import (
	"backend/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListAllProperties handles fetching all properties from the database
func ListAllProperties(ctx *gin.Context) {
	var properties []Property

	// Use db.Select to fetch data
	err := db.DB.Select(&properties, `SELECT * FROM property`)
	if err != nil {
		log.Printf("Error getting properties: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve properties"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": properties})
}

// PropertyInsertion handles the insertion of a new property
func PropertyInsertion(ctx *gin.Context) {
	var property Property

	// Bind JSON to the Property struct
	if err := ctx.ShouldBindJSON(&property); err != nil {
		log.Printf("Error binding JSON for property insertion: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Insert property into the database
	_, err := db.DB.Exec(`INSERT INTO property (name, description) VALUES ($1, $2)`, property.Name, property.Description)
	if err != nil {
		log.Printf("Error inserting property into DB: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert property"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Property successfully inserted"})
}

// UpdateProperty handles updating an existing property
func UpdateProperty(ctx *gin.Context) {
	var property Property

	// Bind JSON to the Property struct
	if err := ctx.ShouldBindJSON(&property); err != nil {
		log.Printf("Error binding JSON for property update: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Update the property in the database
	_, err := db.DB.Exec(`UPDATE property SET name = $1, description = $2 WHERE id = $3`, property.Name, property.Description, property.ID)
	if err != nil {
		log.Printf("Error updating property in DB: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update property"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Property successfully updated"})
}

// DeleteProperty handles deleting a property
func DeleteProperty(ctx *gin.Context) {
	var property Property

	// Bind JSON to the Property struct
	if err := ctx.ShouldBindJSON(&property); err != nil {
		log.Printf("Error binding JSON for property deletion: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Delete the property from the database
	_, err := db.DB.Exec(`DELETE FROM property WHERE id = $1`, property.ID)
	if err != nil {
		log.Printf("Error deleting property from DB: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete property"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Property successfully deleted"})
}
