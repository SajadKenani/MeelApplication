package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"backend/db"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)


func SettingImage(image string, ctx *gin.Context) ([]byte, error) {
	if image == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No image provided"})
		return nil, fmt.Errorf("no image provided")
	}

	supportedFormats := []string{"image/png", "image/jpeg", "image/gif", "image/webp"}


	const prefix = "data:"
	if len(image) > len(prefix) && image[:len(prefix)] == prefix {
		formatEnd := strings.Index(image, ";base64,")
		if formatEnd == -1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 image format"})
			return nil, fmt.Errorf("invalid base64 image format")
		}

		imageFormat := image[len(prefix):formatEnd]

		isValidFormat := false
		for _, format := range supportedFormats {
			if imageFormat == format {
				isValidFormat = true
				break
			}
		}

		if !isValidFormat {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unsupported image format: %s", imageFormat)})
			return nil, fmt.Errorf("unsupported image format: %s", imageFormat)
		}

		image = image[formatEnd+8:] 
	}

	imageData, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding image"})
		log.Printf("Error decoding image: %v", err)
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	if len(imageData) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Decoded image is empty"})
		return nil, fmt.Errorf("decoded image is empty")
	}

	return imageData, nil
}

func InsertingMultipleImages(car Car) ([]string, error) {
	var additionalImageIDs []string
	for _, image := range car.MoreImages {
		var id int
		// Insert each additional image and get its ID
		err := db.DB.QueryRow(`INSERT INTO cars_images (image) VALUES ($1) RETURNING id`, image).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error inserting image: %v", err)
		}

		log.Print("Inserted additional image ID: ", strconv.Itoa(id))
		additionalImageIDs = append(additionalImageIDs, strconv.Itoa(id))
	}

	return additionalImageIDs, nil
}

func FetchImageByID(imgID int) (string, error) {
	var image string
	err := db.DB.QueryRow(`SELECT image FROM cars_images WHERE id = $1`, imgID).Scan(&image)
	if err != nil {
		log.Printf("Error querying database for imgID %d: %s", imgID, err)
		return "", err
	}
	return image, nil
}

func PopulateMoreImages(car *Car) {
	var moreImagesID interface{}

	// Attempt to unmarshal the MoreImagesID JSON string
	err := json.Unmarshal([]byte(car.MoreImagesID), &moreImagesID)
	if err != nil {
		log.Printf("Error unmarshalling MoreImagesID: %s", err)
		return
	}

	// Check if moreImagesID is an array or a single value
	switch v := moreImagesID.(type) {
	case []interface{}:
		for _, imgID := range v {
			if id, ok := imgID.(float64); ok {
				image, err := FetchImageByID(int(id))
				if err == nil {
					car.MoreImages = append(car.MoreImages, image)
				}
			} else {
				log.Printf("Invalid type for image ID in array: %v", imgID)
			}
		}
	case float64:
		image, err := FetchImageByID(int(v))
		if err == nil {
			car.MoreImages = append(car.MoreImages, image)
		}
	default:
		log.Printf("Unexpected type for MoreImagesID: %T", v)
	}
}

func ImageProcessing(cars []Car) {
	for i, car := range cars {
		if len(car.MoreImages) != 0 {
			var moreImages []string
			err := json.Unmarshal([]byte(car.MoreImages[i]), &moreImages)
			if err != nil {
				log.Printf("Error unmarshaling more_images_id for car %d: %v", car.ID, err)
				continue
			}
			cars[i].MoreImages = moreImages
		}
		carImage := base64.StdEncoding.EncodeToString([]byte(car.Image))
		cars[i].Image = string(carImage)
	}
}
