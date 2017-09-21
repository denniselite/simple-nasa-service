package structs

import "time"

type NEOResponse struct {
	Reference   string `json:"reference"`
	Name        string `json:"name"`
	IsHazardous bool `json:"is_hazardous"`
	Date        time.Time `json:"date"`
	Speed       float64 `json:"speed"`
}
