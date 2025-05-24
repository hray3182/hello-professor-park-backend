package dtos

// VehicleEntryExitPayload defines the expected JSON structure for vehicle entry/exit requests.
type VehicleEntryExitPayload struct {
	VehicleID uint   `json:"vehicleID" binding:"required" example:"1"`
	SensorID  string `json:"sensorID" example:"EntrySensor001"`
}
