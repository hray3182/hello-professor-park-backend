package routers

import (
	"hello-professor_backend/controllers"
	"hello-professor_backend/docs" // 匯入 swag 產生的 docs
	"hello-professor_backend/repositories"
	"hello-professor_backend/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 初始化 Repositories
	// vehicleRepo := repositories.NewVehicleRepository() // 移除
	transactionRepo := repositories.NewTransactionRepository()
	parkingRecordRepo := repositories.NewParkingRecordRepository()

	// 初始化 Services
	// vehicleService := services.NewVehicleService(vehicleRepo, parkingRecordRepo) // 移除
	transactionService := services.NewTransactionService(transactionRepo)
	// ParkingRecordService 不再需要 vehicleRepo
	parkingRecordService := services.NewParkingRecordService(parkingRecordRepo) // 修改此處

	// 初始化 Controllers
	// vehicleController := controllers.NewVehicleController(vehicleService) // 移除
	transactionController := controllers.NewTransactionController(transactionService)
	parkingRecordController := controllers.NewParkingRecordController(parkingRecordService)

	// Swagger 文件的基本路徑，對應 main.go 中的 @BasePath
	docs.SwaggerInfo.BasePath = "/api/v1"

	// API v1 路由群組
	apiV1 := router.Group("/api/v1")
	{
		// 車輛路由 (移除)
		// vehicleRoutes := apiV1.Group("/vehicles")
		// {
		// 	vehicleRoutes.POST("", vehicleController.CreateVehicleHandler)
		// 	vehicleRoutes.GET("/search", vehicleController.SearchVehiclesHandler)
		// 	vehicleRoutes.GET("/:id", vehicleController.GetVehicleByIDHandler)
		// 	vehicleRoutes.GET("/license/:plate", vehicleController.GetVehicleByLicensePlateHandler)
		// 	vehicleRoutes.PUT("/:id", vehicleController.UpdateVehicleHandler)
		// 	vehicleRoutes.DELETE("/:id", vehicleController.DeleteVehicleHandler)
		// 	vehicleRoutes.GET("", vehicleController.GetAllVehiclesHandler)
		// }

		// 交易路由
		transactionRoutes := apiV1.Group("/transactions")
		{
			transactionRoutes.POST("", transactionController.CreateTransactionHandler)
			transactionRoutes.GET("/:id", transactionController.GetTransactionByIDHandler)
			transactionRoutes.GET("/parking/:parkingRecordID", transactionController.GetTransactionsByParkingRecordIDHandler)
			transactionRoutes.PUT("/:id", transactionController.UpdateTransactionHandler)
			transactionRoutes.DELETE("/:id", transactionController.DeleteTransactionHandler)
			transactionRoutes.GET("", transactionController.GetAllTransactionsHandler)
		}

		// 停車記錄路由
		parkingRecordRoutes := apiV1.Group("/parking-records")
		{
			parkingRecordRoutes.POST("/entry", parkingRecordController.RecordVehicleEntryHandler)
			parkingRecordRoutes.POST("/exit", parkingRecordController.RecordVehicleExitHandler)
			parkingRecordRoutes.POST("", parkingRecordController.CreateParkingRecordHandler) // 通用建立
			parkingRecordRoutes.GET("/search/license", parkingRecordController.SearchParkingRecordsByLicensePlateHandler)
			parkingRecordRoutes.GET("/:id", parkingRecordController.GetParkingRecordByIDHandler)
			// 修改路由以使用 licensePlate 而非 vehicleID
			parkingRecordRoutes.GET("/license/:licensePlate", parkingRecordController.GetParkingRecordsByLicensePlateHandler)
			parkingRecordRoutes.GET("/license/:licensePlate/latest", parkingRecordController.GetLatestParkingRecordByLicensePlateHandler)
			parkingRecordRoutes.PATCH("/:id/verify-license-plate", parkingRecordController.UpdateUserVerifiedLicensePlateHandler)
			parkingRecordRoutes.POST("/:id/prepare-payment", parkingRecordController.PrepareParkingRecordForPaymentHandler)
			parkingRecordRoutes.POST("/:id/pay", parkingRecordController.PayForParkingRecordHandler)
			parkingRecordRoutes.PUT("/:id", parkingRecordController.UpdateParkingRecordHandler)
			parkingRecordRoutes.DELETE("/:id", parkingRecordController.DeleteParkingRecordHandler)
			parkingRecordRoutes.GET("", parkingRecordController.GetAllParkingRecordsHandler)
		}
	}

	// 設定 Swagger UI 路由
	// 存取 /swagger/index.html 可以看到 API 文件
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
