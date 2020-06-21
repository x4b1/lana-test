package server

type errorResponse struct {
	Message string `json:"message"`
}

type basketCreated struct {
	ID string `json:"id"`
}

type basketTotal struct {
	Total string `json:"total"`
}
