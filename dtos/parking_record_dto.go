package dtos

// VehicleEntryExitPayload defines the expected JSON structure for vehicle entry/exit requests.
type VehicleEntryExitPayload struct {
	LicensePlate string `json:"licensePlate" binding:"required" example:"ABC-1234"`
	SensorID     string `json:"sensorID" example:"EntrySensor001"`
}
