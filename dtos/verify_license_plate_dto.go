package dtos

// VerifyLicensePlatePayload defines the expected JSON structure for verifying/updating a license plate.
type VerifyLicensePlatePayload struct {
	LicensePlate string `json:"licensePlate" binding:"required" example:"XYZ-7890"`
}
