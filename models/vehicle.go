package models

import (
	"time"
)

// Vehicles 車輛資訊
// 對應 PostgreSQL 的 'vehicles' 表
type Vehicle struct {
	// VehicleID 作為主鍵，GORM 會自動設為 primary_key 和 auto_increment
	VehicleID uint `gorm:"primaryKey"`
	// LicensePlate 車牌號碼，唯一索引，確保每輛車有獨特的車牌
	LicensePlate string `gorm:"type:varchar(20);uniqueIndex;not null"`
	// LastSeen 最後一次感應到車輛的時間
	LastSeen time.Time

	// GORM 模型關聯定義
	ParkingRecords []ParkingRecord `gorm:"foreignKey:VehicleID"`
}
