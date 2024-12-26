package handlers

import (
	"encoding/base64"

	"strings"

	// "io"
	// "strings"

	// "fmt"
	"backend/db"
	"net/http"

	// "strings"

	"log"

	"github.com/gin-gonic/gin"
)

func ListAllCars(ctx *gin.Context) {
	var cars []Car
	err := db.DB.Select(&cars, `SELECT id, name, price, description, brand, image, _more_images_id FROM cars`)
	if err != nil {
		log.Print("Error getting data: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cars"})
		return
	}

	ImageProcessing(cars)

	ctx.JSON(http.StatusOK, gin.H{"data": cars})
}
func ListSpecifiedCar(ctx *gin.Context) {
	var car Car

	id := ctx.Param("id")
	err := db.DB.QueryRow(`
		SELECT id, name, price, description, brand, image, _more_images_id, info 
		FROM cars WHERE id = $1`, id).
		Scan(&car.ID, &car.Name, &car.Price, &car.Description, &car.Brand, &car.Image, &car.MoreImagesID, &car.InfoID)

	if err != nil {
		log.Printf("Error getting car data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the car"})
		return
	}

	// Decode the image data to base64
	carImage := base64.StdEncoding.EncodeToString([]byte(car.Image))
	car.Image = carImage

	// Populate MoreImages
	PopulateMoreImages(&car)
	PopulateInfo(&car)


	// Send the response
	ctx.JSON(http.StatusOK, gin.H{"Data": car})
}


func InsertCar(ctx *gin.Context) {
	var car Car

	if err := ctx.ShouldBindJSON(&car); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind data"})
		log.Print("Binding error: ", err)
		return
	}

	log.Print(car.Property)

	if car.Property == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "The property is null!"})
		log.Print(car.Property)
		return
	}

	imageData, err := SettingImage(car.Image, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process image"})
		log.Print("Image processing error: ", err)
		return
	}
	additionalImageIDs, err := InsertingMultipleImages(car)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to insert additional images"})
		log.Print("Error inserting additional images: ", err)
		return
	}

	moreImagesID := strings.Join(additionalImageIDs, ",")

	addInfo, err := InsertInfo(car)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to insert the info"})
		log.Print("Error inserting info: ", err)
		return
	}

	allInfoID := strings.Join(addInfo, ",")

	_, err = db.DB.Exec(`INSERT INTO cars (name, price, image, brand, description, _more_images_id, property, info) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		car.Name, car.Price, imageData, car.Brand, car.Description, "["+moreImagesID+"]", car.Property, "["+allInfoID+"]")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert car data into database"})
		log.Printf("Database insertion error: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Car data and images saved successfully"})
}

func UpdateCar(ctx *gin.Context) {
	var car Car
	err := ctx.ShouldBindJSON(&car)
	if err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err = db.DB.Exec(`UPDATE cars SET name = $1, price = $2, description = $3, brand = $4 WHERE id = $6`,
		car.Name, car.Price, car.Description, car.Brand, car.ID)
	if err != nil {
		log.Printf("Error updating data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the car"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Car was updated successfully"})
}

func SetToFavorite(ctx *gin.Context){
	var car Car

	err := ctx.ShouldBindJSON(&car)
	if err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	myBoolean := car.IsFavorite

	_, err = db.DB.Exec(`UPDATE cars SET is_favorite = $1 WHERE id = $2`, !myBoolean, car.ID)
	if err != nil {
		log.Printf("Error updating data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the car"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Car was added/removed to/from favorites"})

}

func DeleteCar(ctx *gin.Context) {
	var car Car
	err := ctx.ShouldBindJSON(&car)
	if err != nil {
		log.Printf("Error deleting data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the car"})
		return
	}

	_, err = db.DB.Exec(`delete from cars where id = $1`, car.ID)

	if err != nil {
		log.Printf("Error deleting data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the car"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Car was deleted successfully"})

}
