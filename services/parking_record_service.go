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
	GetParkingRecordsByLicensePlate(licensePlate string) ([]models.ParkingRecord, error)
	UpdateParkingRecord(parkingRecord *models.ParkingRecord) error
	DeleteParkingRecord(id uint) error
	GetAllParkingRecords(limit int, offset int) ([]models.ParkingRecord, error)
	GetLatestParkingRecordByLicensePlate(licensePlate string) (*models.ParkingRecord, error)
	RecordVehicleEntry(licensePlate string, sensorEntryID string) (*models.ParkingRecord, error)
	RecordVehicleExit(licensePlate string, sensorExitID string) (*models.ParkingRecord, error)
	UpdateUserVerifiedLicensePlate(recordID uint, verifiedLicensePlate string) (*models.ParkingRecord, error)
}

// parkingRecordService 是 ParkingRecordService 的實作
type parkingRecordService struct {
	parkingRecordRepo repositories.ParkingRecordRepository
}

// NewParkingRecordService 建立一個新的 ParkingRecordService 實例
func NewParkingRecordService(prRepo repositories.ParkingRecordRepository) ParkingRecordService {
	return &parkingRecordService{
		parkingRecordRepo: prRepo,
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

// GetParkingRecordsByLicensePlate 呼叫 repository 透過 LicensePlate 取得相關的所有停車記錄
func (s *parkingRecordService) GetParkingRecordsByLicensePlate(licensePlate string) ([]models.ParkingRecord, error) {
	return s.parkingRecordRepo.GetParkingRecordsByLicensePlate(licensePlate)
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

// GetLatestParkingRecordByLicensePlate 呼叫 repository 透過 LicensePlate 取得最新的停車記錄
func (s *parkingRecordService) GetLatestParkingRecordByLicensePlate(licensePlate string) (*models.ParkingRecord, error) {
	return s.parkingRecordRepo.GetLatestParkingRecordByLicensePlate(licensePlate)
}

// RecordVehicleEntry 記錄車輛進場
func (s *parkingRecordService) RecordVehicleEntry(licensePlate string, sensorEntryID string) (*models.ParkingRecord, error) {
	// 1. 檢查車輛是否已經在場內（是否有未出場的停車記錄）
	latestRecord, err := s.parkingRecordRepo.GetLatestParkingRecordByLicensePlate(licensePlate)
	if err != nil {
		return nil, err
	}
	if latestRecord != nil && latestRecord.ExitTime == nil {
		return nil, errors.New("vehicle already in parking lot")
	}

	// 2. 建立新的停車記錄
	now := time.Now()
	newRecord := &models.ParkingRecord{
		LicensePlate:  licensePlate,
		EntryTime:     now,
		SensorEntryID: sensorEntryID,
		PaymentStatus: "Pending", // 預設為待支付
	}

	err = s.parkingRecordRepo.CreateParkingRecord(newRecord)
	if err != nil {
		return nil, err
	}

	return newRecord, nil
}

// RecordVehicleExit 記錄車輛出場並計算費用
func (s *parkingRecordService) RecordVehicleExit(licensePlate string, sensorExitID string) (*models.ParkingRecord, error) {
	// 1. 獲取該車輛最新的（未出場的）停車記錄
	latestRecord, err := s.parkingRecordRepo.GetLatestParkingRecordByLicensePlate(licensePlate)
	if err != nil {
		return nil, err
	}
	if latestRecord == nil || latestRecord.ExitTime != nil {
		return nil, errors.New("no active parking record found for this vehicle")
	}

	// 2. 更新出場時間、計算停車時長和費用
	now := time.Now()
	latestRecord.ExitTime = &now
	latestRecord.SensorExitID = sensorExitID

	duration := now.Sub(latestRecord.EntryTime)
	latestRecord.ActualDurationMinutes = int(duration.Minutes())

	// TODO: 實作費率計算邏輯 (CalculateParkingFee)
	// 暫時設定一個固定費用或簡單計算
	latestRecord.CalculatedAmount = float64(latestRecord.ActualDurationMinutes) * 0.5 // 例如：每分鐘 0.5 元

	// 3. 更新停車記錄
	err = s.parkingRecordRepo.UpdateParkingRecord(latestRecord)
	if err != nil {
		return nil, err
	}

	return latestRecord, nil
}

// UpdateUserVerifiedLicensePlate 更新使用者驗證的車牌號碼
func (s *parkingRecordService) UpdateUserVerifiedLicensePlate(recordID uint, verifiedLicensePlate string) (*models.ParkingRecord, error) {
	record, err := s.parkingRecordRepo.GetParkingRecordByID(recordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("parking record not found")
	}

	// 更新使用者驗證的車牌號碼
	record.UserVerifiedLicensePlate = &verifiedLicensePlate

	if err := s.parkingRecordRepo.UpdateParkingRecord(record); err != nil {
		return nil, err
	}
	return record, nil
}

// TODO: 需要一個費率計算函式
// func CalculateParkingFee(entryTime, exitTime time.Time, rateStructure interface{}) float64 {
// 	 // 根據費率結構計算費用
// 	 return 0.0
// }
