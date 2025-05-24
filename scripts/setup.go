package main

import (
	"fmt"
	"log"

	"hello-professor_backend/database"
	"hello-professor_backend/models"

	"github.com/joho/godotenv"
)

func main() {
	// load env
	godotenv.Overload(".env")
	// 初始化資料庫連接
	if err := database.InitDB(); err != nil {
		log.Fatalf("資料庫連接失敗: %v", err)
	}

	// 創建資料表（AutoMigrate 內部會處理順序和依賴）
	log.Println("開始進行資料庫遷移...")
	if err := database.AutoMigrate(
		&models.Vehicle{},
		&models.ParkingRecord{},
		&models.Transaction{},
	); err != nil {
		log.Fatalf("資料庫遷移失敗: %v", err)
	}

	fmt.Println("資料庫設置完成！")
}
