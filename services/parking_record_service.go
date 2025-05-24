package services

import (
	"errors"
	"hello-professor_backend/models"
	"hello-professor_backend/repositories"
	"time"
)

// ParkingRecordService 定義停車記錄服務的介面
type ParkingRecordService interface {
	CreateParkingRecord(parkingRecord *models.ParkingRecord) error
	GetParkingRecordByID(id uint) (*models.ParkingRecord, error)
	GetParkingRecordsByVehicleID(vehicleID uint) ([]models.ParkingRecord, error)
	UpdateParkingRecord(parkingRecord *models.ParkingRecord) error
	DeleteParkingRecord(id uint) error
	GetAllParkingRecords(limit int, offset int) ([]models.ParkingRecord, error)
	GetLatestParkingRecordByVehicleID(vehicleID uint) (*models.ParkingRecord, error)
	RecordVehicleEntry(vehicleID uint, sensorEntryID string) (*models.ParkingRecord, error)
	RecordVehicleExit(vehicleID uint, sensorExitID string) (*models.ParkingRecord, error)
}

// parkingRecordService 是 ParkingRecordService 的實作
type parkingRecordService struct {
	parkingRecordRepo repositories.ParkingRecordRepository
	vehicleRepo       repositories.VehicleRepository // 可能需要用來更新車輛最後出現時間等
}

// NewParkingRecordService 建立一個新的 ParkingRecordService 實例
func NewParkingRecordService(prRepo repositories.ParkingRecordRepository, vRepo repositories.VehicleRepository) ParkingRecordService {
	return &parkingRecordService{
		parkingRecordRepo: prRepo,
		vehicleRepo:       vRepo,
	}
}

// CreateParkingRecord 呼叫 repository 來新增停車記錄
func (s *parkingRecordService) CreateParkingRecord(parkingRecord *models.ParkingRecord) error {
	return s.parkingRecordRepo.CreateParkingRecord(parkingRecord)
}

// GetParkingRecordByID 呼叫 repository 透過 ID 取得停車記錄
func (s *parkingRecordService) GetParkingRecordByID(id uint) (*models.ParkingRecord, error) {
	return s.parkingRecordRepo.GetParkingRecordByID(id)
}

// GetParkingRecordsByVehicleID 呼叫 repository 透過 VehicleID 取得相關的所有停車記錄
func (s *parkingRecordService) GetParkingRecordsByVehicleID(vehicleID uint) ([]models.ParkingRecord, error) {
	return s.parkingRecordRepo.GetParkingRecordsByVehicleID(vehicleID)
}

// UpdateParkingRecord 呼叫 repository 來更新停車記錄
func (s *parkingRecordService) UpdateParkingRecord(parkingRecord *models.ParkingRecord) error {
	return s.parkingRecordRepo.UpdateParkingRecord(parkingRecord)
}

// DeleteParkingRecord 呼叫 repository 透過 ID 刪除停車記錄
func (s *parkingRecordService) DeleteParkingRecord(id uint) error {
	return s.parkingRecordRepo.DeleteParkingRecord(id)
}

// GetAllParkingRecords 呼叫 repository 取得所有停車記錄，支援分頁
func (s *parkingRecordService) GetAllParkingRecords(limit int, offset int) ([]models.ParkingRecord, error) {
	return s.parkingRecordRepo.GetAllParkingRecords(limit, offset)
}

// GetLatestParkingRecordByVehicleID 呼叫 repository 透過 VehicleID 取得最新的停車記錄
func (s *parkingRecordService) GetLatestParkingRecordByVehicleID(vehicleID uint) (*models.ParkingRecord, error) {
	return s.parkingRecordRepo.GetLatestParkingRecordByVehicleID(vehicleID)
}

// RecordVehicleEntry 記錄車輛進場
func (s *parkingRecordService) RecordVehicleEntry(vehicleID uint, sensorEntryID string) (*models.ParkingRecord, error) {
	// 1. 檢查車輛是否存在
	vehicle, err := s.vehicleRepo.GetVehicleByID(vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("vehicle not found")
	}

	// 2. 檢查車輛是否已經在場內（是否有未出場的停車記錄）
	latestRecord, err := s.parkingRecordRepo.GetLatestParkingRecordByVehicleID(vehicleID)
	if err != nil {
		return nil, err
	}
	if latestRecord != nil && latestRecord.ExitTime == nil {
		return nil, errors.New("vehicle already in parking lot")
	}

	// 3. 建立新的停車記錄
	now := time.Now()
	newRecord := &models.ParkingRecord{
		VehicleID:     vehicleID,
		EntryTime:     now,
		SensorEntryID: sensorEntryID,
		PaymentStatus: "Pending", // 預設為待支付
	}

	err = s.parkingRecordRepo.CreateParkingRecord(newRecord)
	if err != nil {
		return nil, err
	}

	// 4. (可選) 更新車輛的 LastSeen 時間
	vehicle.LastSeen = now
	if err := s.vehicleRepo.UpdateVehicle(vehicle); err != nil {
		// 此處可以選擇記錄錯誤，但不一定要中斷主要流程
		// log.Printf("Failed to update vehicle last seen time: %v", err)
	}

	return newRecord, nil
}

// RecordVehicleExit 記錄車輛出場並計算費用
func (s *parkingRecordService) RecordVehicleExit(vehicleID uint, sensorExitID string) (*models.ParkingRecord, error) {
	// 1. 檢查車輛是否存在
	vehicle, err := s.vehicleRepo.GetVehicleByID(vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("vehicle not found")
	}

	// 2. 獲取該車輛最新的（未出場的）停車記錄
	latestRecord, err := s.parkingRecordRepo.GetLatestParkingRecordByVehicleID(vehicleID)
	if err != nil {
		return nil, err
	}
	if latestRecord == nil || latestRecord.ExitTime != nil {
		return nil, errors.New("no active parking record found for this vehicle")
	}

	// 3. 更新出場時間、計算停車時長和費用
	now := time.Now()
	latestRecord.ExitTime = &now
	latestRecord.SensorExitID = sensorExitID

	duration := now.Sub(latestRecord.EntryTime)
	latestRecord.ActualDurationMinutes = int(duration.Minutes())

	// TODO: 實作費率計算邏輯 (CalculateParkingFee)
	// 暫時設定一個固定費用或簡單計算
	latestRecord.CalculatedAmount = float64(latestRecord.ActualDurationMinutes) * 0.5 // 例如：每分鐘 0.5 元

	// 4. 更新停車記錄
	err = s.parkingRecordRepo.UpdateParkingRecord(latestRecord)
	if err != nil {
		return nil, err
	}

	// 5. (可選) 更新車輛的 LastSeen 時間
	vehicle.LastSeen = now
	if err := s.vehicleRepo.UpdateVehicle(vehicle); err != nil {
		// log.Printf("Failed to update vehicle last seen time: %v", err)
	}

	return latestRecord, nil
}

// TODO: 需要一個費率計算函式
// func CalculateParkingFee(entryTime, exitTime time.Time, rateStructure interface{}) float64 {
// 	 // 根據費率結構計算費用
// 	 return 0.0
// }
