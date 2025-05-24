package main

// @title Hello Professor API
// @version 1.0
// @description This is the API documentation for the Hello Professor parking management system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

import (
	"hello-professor_backend/database"
	"hello-professor_backend/routers"
	"log"
	"os"

	"github.com/joho/godotenv"
	// _ "hello-professor_backend/docs"
)

func main() {
	loadEnv()

	if err := database.InitDB(); err != nil {
		log.Fatalf("無法初始化資料庫: %v", err)
	}

	// 設定並啟動路由
	router := routers.SetupRouter()
	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("無法啟動伺服器: %v", err)
	}
}

func loadEnv() {
	// if .env file not found, skip it
	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		log.Println("無法載入 .env 檔案，跳過")
		return
	}
	err := godotenv.Overload(envFile)
	if err != nil {
		log.Println("無法載入 .env 檔案，跳過")
	}
}
