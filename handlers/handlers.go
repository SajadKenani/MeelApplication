package handlers

type Car struct {
	ID           int      `json:"id,omitempty" db:"id,omitempty"`
	Name         string   `json:"name" db:"name"`
	Price        float64  `json:"price" db:"price"`
	Image        string   `json:"image" db:"image"`
	Brand        string   `json:"brand" db:"brand"`
	Description  string   `json:"description" db:"description"`
	MoreImages   []string `json:"more_images_id" db:"more_images_id"`
	MoreImagesID string   `json:"_more_images_id" db:"_more_images_id"`
	Property     []string `json:"property" db:"property"`
	Info         []Info   `json:"info" db:"info"`
	InfoID       string   `json:"info_id" db:"info_id"`
	IsFavorite   bool   `json:"is_favorite" db:"is_favorite"`
}

type Image struct {
	ID    int    `json:"id,omitempty" db:"id,omitempty"`
	Image string `json:"image" db:"image"`
}

type Brands struct {
	ID    int    `json:"id,omitempty" db:"id,omitempty"`
	Name  string `json:"name" db:"name"`
	Image string `json:"image" db:"image"`
}

type Property struct {
	ID          int    `json:"id,omitempty" db:"id,omitempty"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type Info struct {
	ID          int    `json:"id,omitempty" db:"id,omitempty"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type SubProperty struct {
	ID         int    `json:"id,omitempty" db:"id,omitempty"`
	Name       string `json:"name" db:"name"`
	PropertyID string `json:"property_id" db:"property_id"`
}

// type httpResponse struct {
// 	Data  any   `json:"data"`
// 	Error error `json:"error,omitempty"`
// 	Message any `json:"message"`
// }
