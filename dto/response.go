package dto

// Response is a data transfer object that maps
// the json response to specific data transfer objects
type Response struct {
	Players Players `json:"players"`
}
