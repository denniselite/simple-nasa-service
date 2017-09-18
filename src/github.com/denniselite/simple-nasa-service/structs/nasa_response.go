package structs

type NasaResponse struct {
	Links            Links `json:"links"`
	Page             Page `json:"page"`
	NearEarthObjects []NearEarthObject `json:"near_earth_objects"`
}

type Links struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
	Self string `json:"self"`
}

type Page struct {
	Size          int `json:"size"`
	TotalElements int `json:"total_elements"`
	TotalPages    int `json:"total_pages"`
	Number        int `json:"number"`
}

type NearEarthObject struct {
	NeoReferenceId                 string `json:"neo_reference_id"`
	Name                           string `json:"name"`
	IsPotentiallyHazardousAsteroid bool `json:"is_potentially_hazardous_asteroid"`
	NearEarthObjectData            []NearEarthObjectDataItem `json:"close_approach_data"`
}

type NearEarthObjectDataItem struct {
	Date             string `json:"close_approach_date"`
	RelativeVelocity RelativeVelocity `json:"relative_velocity"`
}

type RelativeVelocity struct {
	KilometersPerHour string `json:"kilometers_per_hour"`
}
