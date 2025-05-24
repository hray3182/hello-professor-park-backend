package repositories

import (
	"hello-professor_backend/database"
	"hello-professor_backend/models"

	"gorm.io/gorm"
)

// TransactionRepository 定義交易資料庫操作的介面
type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetTransactionsByParkingRecordID(parkingRecordID uint) ([]models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id uint) error
	GetAllTransactions(limit int, offset int) ([]models.Transaction, error)
}

// transactionRepository 是 TransactionRepository 的 GORM 實作
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository 建立一個新的 TransactionRepository 實例
func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{db: database.GetDB()}
}

// CreateTransaction 新增交易記錄
func (r *transactionRepository) CreateTransaction(transaction *models.Transaction) error {
	result := r.db.Create(transaction)
	return result.Error
}

// GetTransactionByID 透過 ID 取得交易記錄
func (r *transactionRepository) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	result := r.db.First(&transaction, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 或者回傳一個特定的 not found 錯誤
		}
		return nil, result.Error
	}
	return &transaction, nil
}

// GetTransactionsByParkingRecordID 透過 ParkingRecordID 取得相關的所有交易記錄
func (r *transactionRepository) GetTransactionsByParkingRecordID(parkingRecordID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := r.db.Where("parking_record_id = ?", parkingRecordID).Find(&transactions)
	return transactions, result.Error
}

// UpdateTransaction 更新交易記錄
func (r *transactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	result := r.db.Save(transaction)
	return result.Error
}

// DeleteTransaction 透過 ID 刪除交易記錄
func (r *transactionRepository) DeleteTransaction(id uint) error {
	result := r.db.Delete(&models.Transaction{}, id)
	return result.Error
}

// GetAllTransactions 取得所有交易記錄，支援分頁
func (r *transactionRepository) GetAllTransactions(limit int, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	dbQuery := r.db
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}
	result := dbQuery.Find(&transactions)
	return transactions, result.Error
}
