package dtos

// ErrorResponse represents a generic error message.
type ErrorResponse struct {
	Error   string `json:"error" example:"Detailed error message"`
	Details string `json:"details,omitempty" example:"Optional additional details"`
}
