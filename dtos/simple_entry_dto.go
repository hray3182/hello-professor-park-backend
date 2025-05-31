package dtos

import "mime/multipart"

// SimpleEntryPayload defines the JSON structure for simple vehicle entry requests using multipart/form-data.
type SimpleEntryPayload struct {
	LicensePlate string                `form:"licensePlate" binding:"required" example:"ABC-1234"`
	Image        *multipart.FileHeader `form:"image" swaggerignore:"true"` // Image file upload
}
