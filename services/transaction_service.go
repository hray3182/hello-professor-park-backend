package services

import (
	"hello-professor_backend/models"
	"hello-professor_backend/repositories"

	"gorm.io/gorm"
)

// TransactionService 定義交易服務的介面
type TransactionService interface {
	CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetTransactionsByParkingRecordID(parkingRecordID uint) ([]models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id uint) error
	GetAllTransactions(limit int, offset int) ([]models.Transaction, error)
}

// transactionService 是 TransactionService 的實作
type transactionService struct {
	transactionRepo repositories.TransactionRepository
}

// NewTransactionService 建立一個新的 TransactionService 實例
func NewTransactionService(repo repositories.TransactionRepository) TransactionService {
	return &transactionService{transactionRepo: repo}
}

// CreateTransaction 呼叫 repository 來新增交易記錄
func (s *transactionService) CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error {
	// 在此處可以加入業務邏輯，例如：
	// - 檢查交易金額是否大於0
	// - 根據 ParkingRecordID 檢查停車記錄是否存在等
	return s.transactionRepo.CreateTransaction(tx, transaction)
}

// GetTransactionByID 呼叫 repository 透過 ID 取得交易記錄
func (s *transactionService) GetTransactionByID(id uint) (*models.Transaction, error) {
	return s.transactionRepo.GetTransactionByID(id)
}

// GetTransactionsByParkingRecordID 呼叫 repository 透過 ParkingRecordID 取得相關的所有交易記錄
func (s *transactionService) GetTransactionsByParkingRecordID(parkingRecordID uint) ([]models.Transaction, error) {
	return s.transactionRepo.GetTransactionsByParkingRecordID(parkingRecordID)
}

// UpdateTransaction 呼叫 repository 來更新交易記錄
func (s *transactionService) UpdateTransaction(transaction *models.Transaction) error {
	// 在此處可以加入業務邏輯
	return s.transactionRepo.UpdateTransaction(transaction)
}

// DeleteTransaction 呼叫 repository 透過 ID 刪除交易記錄
func (s *transactionService) DeleteTransaction(id uint) error {
	return s.transactionRepo.DeleteTransaction(id)
}

// GetAllTransactions 呼叫 repository 取得所有交易記錄，支援分頁
func (s *transactionService) GetAllTransactions(limit int, offset int) ([]models.Transaction, error) {
	return s.transactionRepo.GetAllTransactions(limit, offset)
}
