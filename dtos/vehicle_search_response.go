package dtos

import "hello-professor_backend/models"

// VehicleSearchResult 用於封裝車輛模糊搜尋結果，包含車輛資訊及其活躍的停車記錄。
type VehicleSearchResult struct {
	models.Vehicle
	ActiveParkingRecord *models.ParkingRecord `json:"activeParkingRecord,omitempty"`
}
