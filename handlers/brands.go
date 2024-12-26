package handlers

import (
	"encoding/base64"
	"backend/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListAllBrands(ctx *gin.Context) {
	var brands []Brands

	err := db.DB.Select(&brands, `SELECT name, image from brands`)
	if err != nil {
		log.Print("Error getting brands")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data."})
		return
	}

	for i, brand := range brands {
		brandImage := string(base64.StdEncoding.EncodeToString([]byte(brand.Image)))
		brands[i].Image = brandImage
	}

	ctx.JSON(http.StatusOK, gin.H{"data": brands})
}

func ListBrandsWithoutImage(ctx *gin.Context) {
	var Brands []Brands

	err := db.DB.Select(&Brands, `SELECT name from brands`)
	if err != nil {
		log.Print("Error getting brands")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Brands": Brands})
}

// brand.Image = string(base64.StdEncoding.EncodeToString([]byte(brand.Image)))

func RemoveBrand(ctx *gin.Context){
	var brand Brands

	err := ctx.ShouldBindJSON(&brand)
	if err != nil {
		log.Print("Failed, error with binding process.")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed, error with binding process."})
		return
	}

	_, err = db.DB.Exec("delete from brands where id = $1", brand.ID)
	if err != nil {
		log.Print("Error deleting the brand")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the brand."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "The brand was deleted successfully!"})

}

func InsertBrand(ctx *gin.Context) {
	var brand Brands

	err := ctx.ShouldBindJSON(&brand)
	if err != nil {
		log.Print("Failed, error with binding process.")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed, error with binding process."})
		return
	}

	image, err := SettingImage(brand.Image, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to insert additional images"})
		log.Print("Error inserting additional images: ", err)
		return
	}

	_, err = db.DB.Exec("INSERT INTO brands (name, image) VALUES ($1, $2)", brand.Name, image)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert brand data into database"})
		log.Printf("Database insertion error: %v", err)
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "brand was inserted!"})
}

func UpdateBrand(ctx *gin.Context) {
	var brand Brands
	_, err := db.DB.Exec(`Update brand (name) = $1 WHERE id = $2`, brand.Name, brand.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert brand data into database"})
		log.Printf("Database insertion error: %v", err)
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Brand was updated!"})
}

