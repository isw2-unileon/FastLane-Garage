// Defines request and response data transfer objects for API endpoints.
package dto

// CreatePartRequest represents the request payload for creating a new part.
// Each field has a binding tag to validate incoming JSON data.
type CreatePartRequest struct {
	// Name is required and must not be empty.
	Name string `json:"name" binding:"required,min=1,max=100"`

	// CarZone is required and specifies which zone the part belongs to.
	CarZone string `json:"car_zone" binding:"required,min=1,max=50"`

	// ImageURL is optional but if provided, must be a valid URL format.
	ImageURL string `json:"image_url" binding:"omitempty,url"`

	// Price is required and must be a positive number.
	Price float64 `json:"price" binding:"required,gt=0"`
}

// UpdatePartRequest represents the request payload for updating an existing part.
// All fields are optional (omitempty) to allow partial updates.
type UpdatePartRequest struct {
	// Name is optional but if provided, must not be empty.
	Name string `json:"name" binding:"omitempty,min=1,max=100"`

	// CarZone is optional but if provided, must not be empty.
	CarZone string `json:"car_zone" binding:"omitempty,min=1,max=50"`

	// ImageURL is optional but if provided, must be a valid URL format.
	ImageURL string `json:"image_url" binding:"omitempty,url"`

	// Price is optional but if provided, must be a positive number.
	Price *float64 `json:"price" binding:"omitempty,gt=0"`
}

// PartResponse represents the response payload for a part.
// This is returned when retrieving or creating parts.
type PartResponse struct {
	// ID is the unique identifier of the part.
	ID uint `json:"id"`

	// Name is the descriptive name of the part.
	Name string `json:"name"`

	// CarZone specifies which zone of the car this part belongs to.
	CarZone string `json:"car_zone"`

	// ImageURL is the URL to the image of the part.
	ImageURL string `json:"image_url"`

	// Price is the cost of the part.
	Price float64 `json:"price"`
}
