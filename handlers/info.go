package handlers

import (
	"backend/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RemoveInfo(ctx *gin.Context) {
	var info Info
	err := ctx.ShouldBindJSON(&info)
	if err != nil {
		log.Printf("Error deleting data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the info"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM info where id = $1", info.ID)
	if err != nil {
		log.Printf("Error deleting data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the info"})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "the info was deleted"})

}
