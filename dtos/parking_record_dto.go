package dtos

import (
	"hello-professor_backend/models"
	"time"
)

// VehicleEntryExitPayload defines the expected JSON structure for vehicle entry/exit requests.
type VehicleEntryExitPayload struct {
	LicensePlate string `json:"licensePlate" binding:"required" example:"ABC-1234"`
	SensorID     string `json:"sensorID" example:"EntrySensor001"`
}

// ParkingRecordWithTransactionResponse combines a ParkingRecord with its associated Transaction.
// Used for responses where both are relevant, e.g., after a payment.
type ParkingRecordWithTransactionResponse struct {
	models.ParkingRecord
	Transaction models.Transaction `json:"transaction"`
}

// ErrorResponseWithRecord defines the JSON structure for an error response that includes parking record details.
// Typically used for 402 Payment Required errors during vehicle exit.
type ErrorResponseWithRecord struct {
	Error            string    `json:"error"`
	ParkingRecordID  uint      `json:"parkingRecordID,omitempty"`
	LicensePlate     string    `json:"licensePlate,omitempty"`
	CalculatedAmount float64   `json:"calculatedAmount,omitempty"`
	PaymentStatus    string    `json:"paymentStatus,omitempty"`
	EntryTime        time.Time `json:"entryTime,omitempty"` // Assuming models.ParkingRecord.EntryTime is time.Time
}
