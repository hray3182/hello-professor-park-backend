package services

import (
	"hello-professor_backend/dtos"
	"hello-professor_backend/models"
	"hello-professor_backend/repositories"
)

// VehicleService 定義車輛服務的介面
type VehicleService interface {
	CreateVehicle(vehicle *models.Vehicle) error
	GetVehicleByID(id uint) (*models.Vehicle, error)
	GetVehicleByLicensePlate(licensePlate string) (*models.Vehicle, error)
	UpdateVehicle(vehicle *models.Vehicle) error
	DeleteVehicle(id uint) error
	GetAllVehicles(limit int, offset int) ([]models.Vehicle, error)
	SearchAndGetVehiclesWithActiveParking(plateQuery string, limit int) ([]dtos.VehicleSearchResult, error)
}

// vehicleService 是 VehicleService 的實作
type vehicleService struct {
	vehicleRepo       repositories.VehicleRepository
	parkingRecordRepo repositories.ParkingRecordRepository
}

// NewVehicleService 建立一個新的 VehicleService 實例
// 此處使用 repositories.NewVehicleRepository() 作為預設的 repository
// 在測試或需要不同 repository 實作時，可以傳入不同的 VehicleRepository
func NewVehicleService(vRepo repositories.VehicleRepository, prRepo repositories.ParkingRecordRepository) VehicleService {
	return &vehicleService{
		vehicleRepo:       vRepo,
		parkingRecordRepo: prRepo,
	}
}

// CreateVehicle 呼叫 repository 來新增車輛記錄
func (s *vehicleService) CreateVehicle(vehicle *models.Vehicle) error {
	// 在此處可以加入業務邏輯，例如：
	// - 檢查車牌號碼格式是否正確
	// - 檢查車輛是否已存在等
	return s.vehicleRepo.CreateVehicle(vehicle)
}

// GetVehicleByID 呼叫 repository 透過 ID 取得車輛記錄
func (s *vehicleService) GetVehicleByID(id uint) (*models.Vehicle, error) {
	return s.vehicleRepo.GetVehicleByID(id)
}

// GetVehicleByLicensePlate 呼叫 repository 透過車牌號碼取得車輛記錄
func (s *vehicleService) GetVehicleByLicensePlate(licensePlate string) (*models.Vehicle, error) {
	return s.vehicleRepo.GetVehicleByLicensePlate(licensePlate)
}

// UpdateVehicle 呼叫 repository 來更新車輛記錄
func (s *vehicleService) UpdateVehicle(vehicle *models.Vehicle) error {
	// 在此處可以加入業務邏輯，例如：
	// - 檢查更新的資料是否合法
	return s.vehicleRepo.UpdateVehicle(vehicle)
}

// DeleteVehicle 呼叫 repository 透過 ID 刪除車輛記錄
func (s *vehicleService) DeleteVehicle(id uint) error {
	return s.vehicleRepo.DeleteVehicle(id)
}

// GetAllVehicles 呼叫 repository 取得所有車輛記錄，支援分頁
func (s *vehicleService) GetAllVehicles(limit int, offset int) ([]models.Vehicle, error) {
	return s.vehicleRepo.GetAllVehicles(limit, offset)
}

// SearchAndGetVehiclesWithActiveParking 根據車牌模糊搜尋車輛，並附帶其活躍的停車記錄
func (s *vehicleService) SearchAndGetVehiclesWithActiveParking(plateQuery string, limit int) ([]dtos.VehicleSearchResult, error) {
	vehicles, err := s.vehicleRepo.SearchVehiclesByPlateFuzzy(plateQuery, limit)
	if err != nil {
		return nil, err
	}

	var searchResults []dtos.VehicleSearchResult
	for _, v := range vehicles {
		searchResult := dtos.VehicleSearchResult{Vehicle: v}

		// 獲取最新的停車記錄
		latestParkingRecord, err := s.parkingRecordRepo.GetLatestParkingRecordByVehicleID(v.VehicleID)
		if err != nil {
			// 如果只是找不到停車記錄，不應該中斷整個搜尋，可以選擇記錄錯誤或忽略
			// log.Printf("Error fetching latest parking record for vehicle %d: %v", v.VehicleID, err)
		} else if latestParkingRecord != nil && latestParkingRecord.ExitTime == nil {
			// 檢查是否為活躍的停車記錄 (尚未出場)
			// 未來可能還需要檢查 PaymentStatus 是否為 "Pending" 等
			searchResult.ActiveParkingRecord = latestParkingRecord
		}
		searchResults = append(searchResults, searchResult)
	}

	return searchResults, nil
}
