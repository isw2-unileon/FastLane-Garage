// Package dto defines request and response data transfer objects for API endpoints.
package dto

// TopPartResponse representa los datos de una pieza y cuántas veces ha sido pedida.
type TopPartResponse struct {
	PartID       uint   `json:"part_id"`
	PartName     string `json:"part_name"`
	CarZone      string `json:"car_zone"`
	TotalOrdered int    `json:"total_ordered"`
}
