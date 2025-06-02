package services

import (
	"errors"
	"fmt"
	"hello-professor_backend/configs"
	"hello-professor_backend/dtos"
	"hello-professor_backend/models"
	"hello-professor_backend/repositories"
	"time"

	"gorm.io/gorm"
)

// ParkingRecordService 定義停車記錄服務的介面
type ParkingRecordService interface {
	CreateParkingRecord(parkingRecord *models.ParkingRecord) error
	GetParkingRecordByID(id uint) (*models.ParkingRecord, error)
	GetParkingRecordsByLicensePlate(licensePlate string) ([]models.ParkingRecord, error)
	SearchParkingRecordsByLicensePlate(licensePlateQuery string) ([]models.ParkingRecord, error)
	UpdateParkingRecord(tx *gorm.DB, parkingRecord *models.ParkingRecord) error
	DeleteParkingRecord(id uint) error
	GetAllParkingRecords(limit int, offset int) ([]models.ParkingRecord, error)
	GetLatestParkingRecordByLicensePlate(licensePlate string) (*models.ParkingRecord, error)
	RecordVehicleEntry(licensePlate string, sensorEntryID string, image *string) (*models.ParkingRecord, error)
	RecordSimpleVehicleEntry(licensePlate string, image *string) (*models.ParkingRecord, error)
	RecordVehicleExit(licensePlate string) (*models.ParkingRecord, error)
	UpdateUserVerifiedLicensePlate(recordID uint, verifiedLicensePlate string) (*models.ParkingRecord, error)
	PrepareParkingRecordForPayment(recordID uint) (*models.ParkingRecord, error)
	PayForParkingRecord(recordID uint, paymentPayload dtos.ParkingPaymentPayload) (*models.ParkingRecord, *models.Transaction, error)
	GetTotalParkingCount(startTime, endTime *time.Time) (int64, error)
	GetTotalRevenue(startTime, endTime *time.Time) (float64, error)
	GetImageAttachmentRate(startTime, endTime *time.Time) (*dtos.ImageAttachmentRateResponse, error)
	GetAvailableParkingSpots() (*dtos.AvailableSpotsResponse, error)
}

// parkingRecordService 是 ParkingRecordService 的實作
type parkingRecordService struct {
	parkingRecordRepo  repositories.ParkingRecordRepository
	transactionService TransactionService
	db                 *gorm.DB
}

// NewParkingRecordService 建立一個新的 ParkingRecordService 實例
func NewParkingRecordService(prRepo repositories.ParkingRecordRepository, ts TransactionService, db *gorm.DB) ParkingRecordService {
	return &parkingRecordService{
		parkingRecordRepo:  prRepo,
		transactionService: ts,
		db:                 db,
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

// SearchParkingRecordsByLicensePlate 呼叫 repository 透過 LicensePlate 模糊搜尋相關的所有停車記錄
func (s *parkingRecordService) SearchParkingRecordsByLicensePlate(licensePlateQuery string) ([]models.ParkingRecord, error) {
	return s.parkingRecordRepo.SearchParkingRecordsByLicensePlate(licensePlateQuery)
}

// UpdateParkingRecord 呼叫 repository 來更新停車記錄
func (s *parkingRecordService) UpdateParkingRecord(tx *gorm.DB, parkingRecord *models.ParkingRecord) error {
	return s.parkingRecordRepo.UpdateParkingRecord(tx, parkingRecord)
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
func (s *parkingRecordService) RecordVehicleEntry(licensePlate string, sensorEntryID string, image *string) (*models.ParkingRecord, error) {
	latestRecord, err := s.parkingRecordRepo.GetLatestParkingRecordByLicensePlate(licensePlate)
	if err != nil {
		return nil, fmt.Errorf("error checking for existing record: %w", err)
	}
	if latestRecord != nil && latestRecord.ExitTime == nil {
		return latestRecord, errors.New("vehicle already in parking lot")
	}

	now := time.Now()
	newRecord := &models.ParkingRecord{
		LicensePlate:  licensePlate,
		EntryTime:     now,
		SensorEntryID: sensorEntryID,
		PaymentStatus: "Pending",
		Image:         image,
	}

	err = s.parkingRecordRepo.CreateParkingRecord(newRecord)
	if err != nil {
		return nil, fmt.Errorf("error creating parking record: %w", err)
	}
	return newRecord, nil
}

// RecordSimpleVehicleEntry 記錄車輛簡易進場 (自動使用預設 SensorID)
func (s *parkingRecordService) RecordSimpleVehicleEntry(licensePlate string, image *string) (*models.ParkingRecord, error) {
	const simpleEntrySensorID = "SIMPLE_ENTRY_PORTAL"
	return s.RecordVehicleEntry(licensePlate, simpleEntrySensorID, image)
}

// RecordVehicleExit 記錄車輛出場，並檢查付款狀態
func (s *parkingRecordService) RecordVehicleExit(licensePlate string) (*models.ParkingRecord, error) {
	const defaultExitSensorID = "DEFAULT_EXIT_SENSOR"

	latestRecord, err := s.parkingRecordRepo.GetLatestParkingRecordByLicensePlate(licensePlate)
	if err != nil {
		return nil, fmt.Errorf("error finding active parking record for license plate %s: %w", licensePlate, err)
	}
	if latestRecord == nil {
		return nil, fmt.Errorf("no active parking record found for license plate %s", licensePlate)
	}

	if latestRecord.PaymentStatus != "Paid" {
		calculatedAmount := latestRecord.CalculatedAmount
		if calculatedAmount == 0 && latestRecord.ExitTime == nil {
			tempExitTimeForCalc := time.Now()
			duration := tempExitTimeForCalc.Sub(latestRecord.EntryTime)
			minutes := int(duration.Minutes())
			if minutes < 0 {
				minutes = 0
			}
			calculatedAmount = float64(minutes) * 0.5
		}
		return latestRecord, fmt.Errorf("payment_required: Parking record ID %d for license plate %s requires payment. Amount due: %.2f", latestRecord.RecordID, latestRecord.LicensePlate, calculatedAmount)
	}

	if latestRecord.ExitTime == nil {
		now := time.Now()
		latestRecord.ExitTime = &now
		latestRecord.SensorExitID = defaultExitSensorID

		duration := now.Sub(latestRecord.EntryTime)
		actualMinutes := int(duration.Minutes())
		if actualMinutes < 0 {
			actualMinutes = 0
		}
		latestRecord.ActualDurationMinutes = actualMinutes

		err = s.parkingRecordRepo.UpdateParkingRecord(nil, latestRecord)
		if err != nil {
			return nil, fmt.Errorf("error updating parking record ID %d on exit: %w", latestRecord.RecordID, err)
		}
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
	record.UserVerifiedLicensePlate = &verifiedLicensePlate
	if err := s.parkingRecordRepo.UpdateParkingRecord(nil, record); err != nil {
		return nil, err
	}
	return record, nil
}

// PrepareParkingRecordForPayment 準備停車記錄以進行付款
func (s *parkingRecordService) PrepareParkingRecordForPayment(recordID uint) (*models.ParkingRecord, error) {
	record, err := s.parkingRecordRepo.GetParkingRecordByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("error finding parking record ID %d: %w", recordID, err)
	}
	if record == nil {
		return nil, fmt.Errorf("parking record ID %d not found", recordID)
	}

	if record.ExitTime != nil {
		return record, fmt.Errorf("vehicle_exited: Vehicle has already exited on %v. Fee is final at %.2f", *record.ExitTime, record.CalculatedAmount)
	}

	if record.PaymentStatus == "Paid" {
		return record, fmt.Errorf("already_paid: Parking record is already paid. Amount was %.2f", record.CalculatedAmount)
	}

	effectiveCalculationTime := time.Now()
	duration := effectiveCalculationTime.Sub(record.EntryTime)
	actualMinutes := int(duration.Minutes())
	if actualMinutes < 0 {
		actualMinutes = 0
	}

	// TODO: Implement proper rate calculation based on configs.RatePerUnit and configs.UnitDurationHours
	calculatedAmount := float64(actualMinutes) * configs.RatePerUnit

	record.ActualDurationMinutes = actualMinutes
	record.CalculatedAmount = calculatedAmount + 10 //基礎費用

	if err = s.parkingRecordRepo.UpdateParkingRecord(nil, record); err != nil {
		return nil, fmt.Errorf("error updating parking record ID %d with calculated fee: %w", recordID, err)
	}

	return record, nil
}

// PayForParkingRecord 處理特定停車記錄的支付
func (s *parkingRecordService) PayForParkingRecord(recordID uint, paymentPayload dtos.ParkingPaymentPayload) (pr *models.ParkingRecord, tr *models.Transaction, err error) {
	pr, err = s.parkingRecordRepo.GetParkingRecordByID(recordID)
	if err != nil {
		err = fmt.Errorf("error finding parking record ID %d: %w", recordID, err)
		return
	}
	if pr == nil {
		err = fmt.Errorf("parking record ID %d not found", recordID)
		return
	}

	if pr.ExitTime != nil {
		err = fmt.Errorf("vehicle_exited: Cannot pay for an already exited record. Fee was %.2f", pr.CalculatedAmount)
		return
	}
	if pr.PaymentStatus == "Paid" {
		err = fmt.Errorf("already_paid: Parking record ID %d is already paid.", recordID)
		return
	}

	if pr.CalculatedAmount <= 0 {
		err = fmt.Errorf("fee_not_calculated: Fee for parking record ID %d has not been calculated. Please call prepare-payment first.", recordID)
		return
	}

	if paymentPayload.AmountPaid != pr.CalculatedAmount {
		err = fmt.Errorf("amount_mismatch: Amount paid (%.2f) does not match calculated amount (%.2f) for parking record ID %d.", paymentPayload.AmountPaid, pr.CalculatedAmount, recordID)
		return
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		err = fmt.Errorf("failed to begin transaction: %w", tx.Error)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", commitErr)
			}
		}
	}()

	newTransaction := &models.Transaction{
		ParkingRecordID:        recordID,
		Amount:                 paymentPayload.AmountPaid,
		TransactionTime:        time.Now(),
		PaymentMethod:          paymentPayload.PaymentMethod,
		Status:                 "Success",
		PaymentGatewayResponse: paymentPayload.PaymentReference,
	}

	if err = s.transactionService.CreateTransaction(tx, newTransaction); err != nil {
		err = fmt.Errorf("failed to create transaction record: %w", err)
		return
	}

	pr.PaymentStatus = "Paid"
	pr.TransactionID = &newTransaction.TransactionID

	if err = s.parkingRecordRepo.UpdateParkingRecord(tx, pr); err != nil {
		err = fmt.Errorf("failed to update parking record status to Paid: %w", err)
		return
	}

	tr = newTransaction
	return
}

// --- 報表服務方法 ---

// GetTotalParkingCount 獲取指定時間範圍內的總停車次數
func (s *parkingRecordService) GetTotalParkingCount(startTime, endTime *time.Time) (int64, error) {
	return s.parkingRecordRepo.CountParkingRecords(startTime, endTime)
}

// GetTotalRevenue 獲取指定時間範圍內的總收入
func (s *parkingRecordService) GetTotalRevenue(startTime, endTime *time.Time) (float64, error) {
	return s.parkingRecordRepo.SumPaidParkingFees(startTime, endTime)
}

// GetImageAttachmentRate 獲取指定時間範圍內停車記錄的圖片附件率
func (s *parkingRecordService) GetImageAttachmentRate(startTime, endTime *time.Time) (*dtos.ImageAttachmentRateResponse, error) {
	totalEntries, err := s.parkingRecordRepo.CountParkingRecords(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error getting total parking records for image rate: %w", err)
	}

	entriesWithImage, err := s.parkingRecordRepo.CountParkingRecordsWithImage(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error getting parking records with image for image rate: %w", err)
	}

	var attachmentRate float64
	if totalEntries > 0 {
		attachmentRate = float64(entriesWithImage) / float64(totalEntries)
	}

	return &dtos.ImageAttachmentRateResponse{
		TotalEntries:     totalEntries,
		EntriesWithImage: entriesWithImage,
		AttachmentRate:   attachmentRate,
	}, nil
}

// GetAvailableParkingSpots 獲取停車場總容量、已佔用車位和可用車位數量
func (s *parkingRecordService) GetAvailableParkingSpots() (*dtos.AvailableSpotsResponse, error) {
	totalCapacity := configs.ParkingLotCapacity

	occupiedSpots, err := s.parkingRecordRepo.CountActiveParkingRecords()
	if err != nil {
		return nil, fmt.Errorf("error counting active parking records: %w", err)
	}

	availableSpots := int64(totalCapacity) - occupiedSpots
	if availableSpots < 0 {
		availableSpots = 0
	}

	return &dtos.AvailableSpotsResponse{
		TotalCapacity:  totalCapacity,
		OccupiedSpots:  occupiedSpots,
		AvailableSpots: availableSpots,
	}, nil
}

// TODO: 需要一個費率計算函式
// func CalculateParkingFee(entryTime, exitTime time.Time, rateStructure interface{}) float64 {
// 	 // 根據費率結構計算費用
// 	 return 0.0
// }
