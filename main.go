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
	err := godotenv.Overload(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
