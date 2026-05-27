// Defines the data structures for the application
package models

// Part represents an automobile part in the catalog
// It maps to the "parts" table in the SQLite database
type Part struct {
	// ID is the unique identifier for the part (primary key)
	ID uint `gorm:"primaryKey" json:"id"`

	// Name is the descriptive name of the part
	Name string `json:"name"`

	// CarZone specifies which zone of the car this part belongs to
	CaroZone string `json:"car_zone"`

	// ImageURL is the URL to the image of the part for frontend rendering
	ImageURL string `json:"image_url"`

	// Price is the cost of the part in currency
	Price float64 `json:"price"`
}

// TabelName specifies the database table name for the Part model
// GORM uses this method to determinate which table to query
func (Part) TableName() string {
	return "parts"
}
