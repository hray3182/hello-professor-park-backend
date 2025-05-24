package repositories

import (
	"hello-professor_backend/database"
	"hello-professor_backend/models"
	"strings"

	"gorm.io/gorm"
)

// VehicleRepository 定義車輛資料庫操作的介面
type VehicleRepository interface {
	CreateVehicle(vehicle *models.Vehicle) error
	GetVehicleByID(id uint) (*models.Vehicle, error)
	GetVehicleByLicensePlate(licensePlate string) (*models.Vehicle, error)
	UpdateVehicle(vehicle *models.Vehicle) error
	DeleteVehicle(id uint) error
	GetAllVehicles(limit int, offset int) ([]models.Vehicle, error)
	SearchVehiclesByPlateFuzzy(plateQuery string, limit int) ([]models.Vehicle, error)
}

// vehicleRepository 是 VehicleRepository 的 GORM 實作
type vehicleRepository struct {
	db *gorm.DB
}

// NewVehicleRepository 建立一個新的 VehicleRepository 實例
func NewVehicleRepository() VehicleRepository {
	return &vehicleRepository{db: database.GetDB()}
}

// CreateVehicle 新增車輛記錄
func (r *vehicleRepository) CreateVehicle(vehicle *models.Vehicle) error {
	result := r.db.Create(vehicle)
	return result.Error
}

// GetVehicleByID 透過 ID 取得車輛記錄
func (r *vehicleRepository) GetVehicleByID(id uint) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	result := r.db.First(&vehicle, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 或者回傳一個特定的 not found 錯誤
		}
		return nil, result.Error
	}
	return &vehicle, nil
}

// GetVehicleByLicensePlate 透過車牌號碼取得車輛記錄
func (r *vehicleRepository) GetVehicleByLicensePlate(licensePlate string) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	result := r.db.Where("license_plate = ?", licensePlate).First(&vehicle)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 或者回傳一個特定的 not found 錯誤
		}
		return nil, result.Error
	}
	return &vehicle, nil
}

// UpdateVehicle 更新車輛記錄
func (r *vehicleRepository) UpdateVehicle(vehicle *models.Vehicle) error {
	result := r.db.Save(vehicle)
	return result.Error
}

// DeleteVehicle 透過 ID 刪除車輛記錄
func (r *vehicleRepository) DeleteVehicle(id uint) error {
	result := r.db.Delete(&models.Vehicle{}, id)
	return result.Error
}

// GetAllVehicles 取得所有車輛記錄，支援分頁
func (r *vehicleRepository) GetAllVehicles(limit int, offset int) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	dbQuery := r.db
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}
	result := dbQuery.Find(&vehicles)
	return vehicles, result.Error
}

// normalizePlate 標準化車牌字串 (轉小寫，移除空格和破折號)
func normalizePlate(plate string) string {
	plate = strings.ToLower(plate)
	plate = strings.ReplaceAll(plate, " ", "")
	plate = strings.ReplaceAll(plate, "-", "")
	return plate
}

// SearchVehiclesByPlateFuzzy 根據車牌號碼進行模糊搜尋
func (r *vehicleRepository) SearchVehiclesByPlateFuzzy(plateQuery string, limit int) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle

	// 標準化查詢字串
	normalizedQuery := normalizePlate(plateQuery)

	// 使用 ILIKE 進行模糊查詢，並在查詢前對資料庫中的 license_plate 進行相同的標準化處理
	// 注意：在 WHERE 子句中使用函數 (如 LOWER, REPLACE) 可能會影響索引的使用效率。
	// 對於非常大的資料表，可能需要考慮其他優化方式，例如預先儲存標準化後的車牌欄位。
	dbQuery := r.db.Where("REPLACE(REPLACE(LOWER(license_plate), ' ', ''), '-', '') ILIKE ?", "%"+normalizedQuery+"%")

	if limit <= 0 {
		limit = 5 // 預設最多返回 5 筆
	}
	dbQuery = dbQuery.Limit(limit)

	result := dbQuery.Find(&vehicles)
	return vehicles, result.Error
}
