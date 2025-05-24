package models

import "time"

// Transaction 交易紀錄
// 對應 PostgreSQL 的 'transactions' 表
type Transaction struct {
	// TransactionID 作為主鍵
	TransactionID uint `gorm:"primaryKey"`
	// ParkingRecordID 關聯到 ParkingRecords 表的外鍵
	ParkingRecordID uint `gorm:"not null"`
	// Amount 交易金額
	Amount float64 `gorm:"type:decimal(10,2);not null"`
	// TransactionTime 交易時間
	TransactionTime time.Time `gorm:"not null"`
	// PaymentMethod 付款方式，例如 "CreditCard", "MobilePay", "Cash"
	PaymentMethod string `gorm:"type:varchar(50);not null"`
	// Status 交易狀態：Success, Failed, Refunded
	Status string `gorm:"type:varchar(20);not null;default:'Success'"`
	// PaymentGatewayResponse 支付閘道回傳的詳細資訊 (JSON或TEXT)
	PaymentGatewayResponse string `gorm:"type:text"`
}
