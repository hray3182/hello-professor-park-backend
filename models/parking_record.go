package models

import "time"

// ParkingRecord 停車紀錄
// 對應 PostgreSQL 的 'parking_records' 表
type ParkingRecord struct {
	// RecordID 作為主鍵
	RecordID uint `gorm:"primaryKey"`
	// LicensePlate 車牌號碼 (通常來自 OCR)
	LicensePlate string `gorm:"type:varchar(20);not null;index"`
	// UserVerifiedLicensePlate 使用者驗證/修正後的車牌號碼，可以為 NULL
	UserVerifiedLicensePlate *string `gorm:"type:varchar(20)"`
	// EntryTime 進場時間
	EntryTime time.Time `gorm:"not null"`
	// ExitTime 出場時間，如果尚未出場則為 NULL
	ExitTime *time.Time
	// ActualDurationMinutes 實際停車時長（分鐘）
	ActualDurationMinutes int `gorm:"default:0"` // 預設值為 0
	// CalculatedAmount 應付停車費用
	CalculatedAmount float64 `gorm:"type:decimal(10,2);default:0.00"`
	// PaymentStatus 支付狀態：Pending, Paid, Refunded
	PaymentStatus string `gorm:"type:varchar(20);not null;default:'Pending'"`
	// TransactionID 關聯到 Transactions 表的外鍵，如果尚未支付或無交易則為 NULL
	TransactionID *uint // 使用指針表示可為 NULL
	// SensorEntryID 入場感應器記錄ID
	SensorEntryID string `gorm:"type:varchar(100)"`
	// SensorExitID 出場感應器記錄ID
	SensorExitID string `gorm:"type:varchar(100)"`

	// GORM 模型關聯定義
	// Vehicle     Vehicle     `gorm:"foreignKey:VehicleID"` // 移除 Vehicle 關聯
	Transaction Transaction `gorm:"foreignKey:TransactionID"`

	// New fields
	Image *string `json:"image,omitempty" gorm:"type:text"` // 新增欄位: 圖片 URL 或 Base64
}
