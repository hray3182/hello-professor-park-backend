package repositories

import (
	"hello-professor_backend/database"
	"hello-professor_backend/models"

	"strings"

	"gorm.io/gorm"
)

// ParkingRecordRepository 定義停車記錄資料庫操作的介面
type ParkingRecordRepository interface {
	CreateParkingRecord(parkingRecord *models.ParkingRecord) error
	GetParkingRecordByID(id uint) (*models.ParkingRecord, error)
	GetParkingRecordsByLicensePlate(licensePlate string) ([]models.ParkingRecord, error)
	SearchParkingRecordsByLicensePlate(licensePlateQuery string) ([]models.ParkingRecord, error)
	UpdateParkingRecord(parkingRecord *models.ParkingRecord) error
	DeleteParkingRecord(id uint) error
	GetAllParkingRecords(limit int, offset int) ([]models.ParkingRecord, error)
	GetLatestParkingRecordByLicensePlate(licensePlate string) (*models.ParkingRecord, error)
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
	// Preload Transaction to get associated data
	result := r.db.Preload("Transaction").First(&record, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 或者回傳一個特定的 not found 錯誤
		}
		return nil, result.Error
	}
	return &record, nil
}

// GetParkingRecordsByLicensePlate 透過 LicensePlate 取得相關的所有停車記錄
// 會同時比對 LicensePlate 和 UserVerifiedLicensePlate 欄位
func (r *parkingRecordRepository) GetParkingRecordsByLicensePlate(licensePlate string) ([]models.ParkingRecord, error) {
	var records []models.ParkingRecord
	result := r.db.Preload("Transaction").
		Where("license_plate = ? OR user_verified_license_plate = ?", licensePlate, licensePlate).
		Order("entry_time DESC").
		Find(&records)
	return records, result.Error
}

// SearchParkingRecordsByLicensePlate 透過 LicensePlate 模糊搜尋相關的所有停車記錄
// 會同時比對 LicensePlate 和 UserVerifiedLicensePlate 欄位，不區分大小寫
func (r *parkingRecordRepository) SearchParkingRecordsByLicensePlate(licensePlateQuery string) ([]models.ParkingRecord, error) {
	var records []models.ParkingRecord
	// 定義相似度閾值，您可以根據需求調整這個值（0.0 到 1.0 之間）
	const similarityThreshold = 0.3 // 例如，0.3 表示 30% 的相似度

	// 將查詢字串轉換為小寫，以進行不區分大小寫的相似度比較
	lowerLicensePlateQuery := strings.ToLower(licensePlateQuery)

	result := r.db.Preload("Transaction").
		Where("similarity(LOWER(license_plate), ?) > ? OR similarity(LOWER(user_verified_license_plate), ?) > ?", lowerLicensePlateQuery, similarityThreshold, lowerLicensePlateQuery, similarityThreshold).
		Order("entry_time DESC").
		Find(&records)
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
	dbQuery := r.db.Preload("Transaction")
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}
	result := dbQuery.Find(&records)
	return records, result.Error
}

// GetLatestParkingRecordByLicensePlate 透過 LicensePlate 取得最新的停車記錄（基於 EntryTime 降序）
// 會同時比對 LicensePlate 和 UserVerifiedLicensePlate 欄位
func (r *parkingRecordRepository) GetLatestParkingRecordByLicensePlate(licensePlate string) (*models.ParkingRecord, error) {
	var record models.ParkingRecord
	// 查詢條件修改為同時檢查 LicensePlate 或 UserVerifiedLicensePlate，並且 ExitTime 為 NULL (表示仍在場內)
	result := r.db.Preload("Transaction").
		Where("(license_plate = ? OR user_verified_license_plate = ?) AND exit_time IS NULL", licensePlate, licensePlate).
		Order("entry_time DESC").First(&record)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 或者回傳一個特定的 not found 錯誤
		}
		return nil, result.Error
	}
	return &record, nil
}
