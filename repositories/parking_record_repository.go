package repositories

import (
	"hello-professor_backend/database"
	"hello-professor_backend/models"

	"gorm.io/gorm"
)

// ParkingRecordRepository 定義停車記錄資料庫操作的介面
type ParkingRecordRepository interface {
	CreateParkingRecord(parkingRecord *models.ParkingRecord) error
	GetParkingRecordByID(id uint) (*models.ParkingRecord, error)
	GetParkingRecordsByVehicleID(vehicleID uint) ([]models.ParkingRecord, error)
	UpdateParkingRecord(parkingRecord *models.ParkingRecord) error
	DeleteParkingRecord(id uint) error
	GetAllParkingRecords(limit int, offset int) ([]models.ParkingRecord, error)
	GetLatestParkingRecordByVehicleID(vehicleID uint) (*models.ParkingRecord, error)
}

// parkingRecordRepository 是 ParkingRecordRepository 的 GORM 實作
type parkingRecordRepository struct {
	db *gorm.DB
}

// NewParkingRecordRepository 建立一個新的 ParkingRecordRepository 實例
func NewParkingRecordRepository() ParkingRecordRepository {
	return &parkingRecordRepository{db: database.GetDB()}
}

// CreateParkingRecord 新增停車記錄
func (r *parkingRecordRepository) CreateParkingRecord(parkingRecord *models.ParkingRecord) error {
	result := r.db.Create(parkingRecord)
	return result.Error
}

// GetParkingRecordByID 透過 ID 取得停車記錄
func (r *parkingRecordRepository) GetParkingRecordByID(id uint) (*models.ParkingRecord, error) {
	var record models.ParkingRecord
	// Preload Vehicle and Transaction to get associated data
	result := r.db.Preload("Vehicle").Preload("Transaction").First(&record, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 或者回傳一個特定的 not found 錯誤
		}
		return nil, result.Error
	}
	return &record, nil
}

// GetParkingRecordsByVehicleID 透過 VehicleID 取得相關的所有停車記錄
func (r *parkingRecordRepository) GetParkingRecordsByVehicleID(vehicleID uint) ([]models.ParkingRecord, error) {
	var records []models.ParkingRecord
	result := r.db.Preload("Vehicle").Preload("Transaction").Where("vehicle_id = ?", vehicleID).Find(&records)
	return records, result.Error
}

// UpdateParkingRecord 更新停車記錄
func (r *parkingRecordRepository) UpdateParkingRecord(parkingRecord *models.ParkingRecord) error {
	result := r.db.Save(parkingRecord)
	return result.Error
}

// DeleteParkingRecord 透過 ID 刪除停車記錄
func (r *parkingRecordRepository) DeleteParkingRecord(id uint) error {
	result := r.db.Delete(&models.ParkingRecord{}, id)
	return result.Error
}

// GetAllParkingRecords 取得所有停車記錄，支援分頁
func (r *parkingRecordRepository) GetAllParkingRecords(limit int, offset int) ([]models.ParkingRecord, error) {
	var records []models.ParkingRecord
	dbQuery := r.db.Preload("Vehicle").Preload("Transaction")
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}
	result := dbQuery.Find(&records)
	return records, result.Error
}

// GetLatestParkingRecordByVehicleID 透過 VehicleID 取得最新的停車記錄（基於 EntryTime 降序）
func (r *parkingRecordRepository) GetLatestParkingRecordByVehicleID(vehicleID uint) (*models.ParkingRecord, error) {
	var record models.ParkingRecord
	result := r.db.Preload("Vehicle").Preload("Transaction").Where("vehicle_id = ?", vehicleID).Order("entry_time DESC").First(&record)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 或者回傳一個特定的 not found 錯誤
		}
		return nil, result.Error
	}
	return &record, nil
}
