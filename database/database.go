package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var originalGormConfig *gorm.Config // Store original GORM config

// InitDB 初始化資料庫連接
func InitDB() error {
	// 從環境變數獲取資料庫 URI
	databaseURI := os.Getenv("DATABASE_URI")
	if databaseURI == "" {
		log.Fatal("DATABASE_URI 環境變數未設置")
	}

	// 配置 GORM 日誌
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	originalGormConfig = &gorm.Config{
		Logger: newLogger,
	}

	// 建立資料庫連接
	var err error
	DB, err = gorm.Open(postgres.Open(databaseURI), originalGormConfig)
	if err != nil {
		return fmt.Errorf("無法連接到資料庫: %v", err)
	}

	// 配置連接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("無法獲取資料庫實例: %v", err)
	}

	// 設置連接池參數
	sqlDB.SetMaxIdleConns(10)           // 最大空閒連接數
	sqlDB.SetMaxOpenConns(100)          // 最大打開連接數
	sqlDB.SetConnMaxLifetime(time.Hour) // 連接最大生命週期

	return nil
}

// AutoMigrate 自動遷移資料庫結構
func AutoMigrate(allModels ...interface{}) error {
	// 獲取底層的 sql.DB 連接，以便重用
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("無法獲取底層 sql.DB 以進行遷移: %v", err)
	}

	// 階段一：創建資料表，禁用外鍵約束
	configForPhase1 := &gorm.Config{
		Logger:                                   originalGormConfig.Logger,
		NamingStrategy:                           originalGormConfig.NamingStrategy,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	dbPhase1, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), configForPhase1)
	if err != nil {
		return fmt.Errorf("無法為階段 1 創建 GORM DB 實例: %v", err)
	}

	log.Println("遷移階段 1：創建資料表結構（禁用外鍵）...")
	for _, model := range allModels {
		log.Printf("階段 1: 遷移模型 %T\n", model)
		if err := dbPhase1.AutoMigrate(model); err != nil {
			return fmt.Errorf("自動遷移失敗 (階段 1: 創建資料表) 模型 %T: %v", model, err)
		}
	}
	log.Println("遷移階段 1 完成。")

	// 階段二：使用原始的 DB 實例（啟用了外鍵約束）來創建外鍵
	// 原始的 DB 實例是使用 originalGormConfig 初始化的，其中 DisableForeignKeyConstraintWhenMigrating 預設為 false
	log.Println("遷移階段 2：創建外鍵和約束...")
	for _, model := range allModels {
		log.Printf("階段 2: 遷移模型 %T\n", model)
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("自動遷移失敗 (階段 2: 創建外鍵) 模型 %T: %v", model, err)
		}
	}
	log.Println("遷移階段 2 完成。")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
