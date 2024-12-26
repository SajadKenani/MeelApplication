package handlers

import (
	"backend/db"
	"encoding/json"
	"fmt"
	"strconv"
	"log"
)

func InsertInfo(car Car) ([]string, error) {
	var InfoIDs []string
	for _, myInfo := range car.Info {
		var id int
		// Insert only `name` and `description`, let `id` auto-increment
		err := db.DB.QueryRow(
			`INSERT INTO info (name, description) VALUES ($1, $2) RETURNING id`,
			myInfo.Name, myInfo.Description, // Provide values for the name and description columns
		).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error inserting info: %v", err)
		}
		InfoIDs = append(InfoIDs, strconv.Itoa(id)) // Add the generated id to the list
	}
	return InfoIDs, nil
}

func FetchInfoByID(infoID int) (Info, error) {
	var info Info
	err := db.DB.QueryRow(`SELECT * FROM info WHERE id = $1`, infoID).Scan(&info.ID, &info.Name, &info.Description)
	if err != nil {
		log.Printf("Error querying database for infoID %d: %s", infoID, err)
		return info, err
	}
	return info, nil
}

func PopulateInfo(car *Car) {
	var moreInfo interface{}
	err := json.Unmarshal([]byte(car.InfoID), &moreInfo)
	if err != nil {
		log.Printf("Error unmarshalling MoreImagesID: %s", err)
		return
	}

	switch v := moreInfo.(type) {
	case []interface{}:
		for _, infoID := range v {
			if id, ok := infoID.(float64); ok {
				info, err := FetchInfoByID(int(id))
				if err == nil {
					car.Info = append(car.Info, info)
				}
			} else {
				log.Printf("Invalid type for info ID in array: %v", infoID)
			}
		}
	case float64:
		info, err := FetchInfoByID(int(v))
		if err == nil {
			car.Info = append(car.Info, info)
		}
	default:
		log.Printf("Unexpected type for MoreImagesID: %T", v)
	}

}
