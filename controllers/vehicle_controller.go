package controllers

import (
	"hello-professor_backend/models"
	"hello-professor_backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// VehicleController 定義車輛控制器
type VehicleController struct {
	vehicleService services.VehicleService
}

// NewVehicleController 建立一個新的 VehicleController 實例
func NewVehicleController(vs services.VehicleService) *VehicleController {
	return &VehicleController{vehicleService: vs}
}

// CreateVehicleHandler godoc
// @Summary Create a new vehicle
// @Description Add a new vehicle to the system
// @Tags vehicles
// @Accept  json
// @Produce  json
// @Param   vehicle_info body models.Vehicle true "Vehicle Information"
// @Success 201 {object} models.Vehicle
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /vehicles [post]
func (vc *VehicleController) CreateVehicleHandler(c *gin.Context) {
	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := vc.vehicleService.CreateVehicle(&vehicle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vehicle: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, vehicle)
}

// GetVehicleByIDHandler godoc
// @Summary Get a vehicle by ID
// @Description Get details of a vehicle by its ID
// @Tags vehicles
// @Produce  json
// @Param   id path int true "Vehicle ID"
// @Success 200 {object} models.Vehicle
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /vehicles/{id} [get]
func (vc *VehicleController) GetVehicleByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID format"})
		return
	}

	vehicle, err := vc.vehicleService.GetVehicleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get vehicle: " + err.Error()})
		return
	}
	if vehicle == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}
	c.JSON(http.StatusOK, vehicle)
}

// GetVehicleByLicensePlateHandler godoc
// @Summary Get a vehicle by license plate
// @Description Get details of a vehicle by its license plate
// @Tags vehicles
// @Produce json
// @Param plate path string true "License Plate"
// @Success 200 {object} models.Vehicle
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /vehicles/license/{plate} [get]
func (vc *VehicleController) GetVehicleByLicensePlateHandler(c *gin.Context) {
	plate := c.Param("plate")
	vehicle, err := vc.vehicleService.GetVehicleByLicensePlate(plate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get vehicle by license plate: " + err.Error()})
		return
	}
	if vehicle == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found with license plate: " + plate})
		return
	}
	c.JSON(http.StatusOK, vehicle)
}

// UpdateVehicleHandler godoc
// @Summary Update an existing vehicle
// @Description Update details of an existing vehicle by its ID
// @Tags vehicles
// @Accept  json
// @Produce  json
// @Param   id path int true "Vehicle ID"
// @Param   vehicle_update body models.Vehicle true "Vehicle Update Information"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /vehicles/{id} [put]
func (vc *VehicleController) UpdateVehicleHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID format"})
		return
	}

	var vehicleUpdates models.Vehicle
	if err := c.ShouldBindJSON(&vehicleUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 確保更新的是正確的 ID
	vehicleUpdates.VehicleID = uint(id)

	if err := vc.vehicleService.UpdateVehicle(&vehicleUpdates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vehicle updated successfully"})
}

// DeleteVehicleHandler godoc
// @Summary Delete a vehicle by ID
// @Description Remove a vehicle from the system by its ID
// @Tags vehicles
// @Produce  json
// @Param   id path int true "Vehicle ID"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /vehicles/{id} [delete]
func (vc *VehicleController) DeleteVehicleHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID format"})
		return
	}

	if err := vc.vehicleService.DeleteVehicle(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vehicle deleted successfully"})
}

// GetAllVehiclesHandler godoc
// @Summary Get all vehicles
// @Description Get a list of all vehicles, with pagination
// @Tags vehicles
// @Produce  json
// @Param limit query int false "Limit number of vehicles returned" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.Vehicle
// @Failure 500 {object} dtos.ErrorResponse
// @Router /vehicles [get]
func (vc *VehicleController) GetAllVehiclesHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")  // 預設每頁 10 筆
	offsetStr := c.DefaultQuery("offset", "0") // 預設從第 0 筆開始

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // 如果轉換失敗或不合法，使用預設值
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // 如果轉換失敗或不合法，使用預設值
	}

	vehicles, err := vc.vehicleService.GetAllVehicles(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all vehicles: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, vehicles)
}

// SearchVehiclesHandler godoc
// @Summary Search vehicles by license plate (fuzzy search)
// @Description Search for vehicles using a partial or potentially incorrect license plate string. Returns a list of matching vehicles, each with their active parking record if any.
// @Tags vehicles
// @Produce  json
// @Param query query string true "License plate query string"
// @Param limit query int false "Limit number of results" default(5)
// @Success 200 {array} dtos.VehicleSearchResult
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /vehicles/search [get]
func (vc *VehicleController) SearchVehiclesHandler(c *gin.Context) {
	plateQuery := c.Query("query")
	if plateQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'query' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5 // Default to 5 if invalid
	}

	results, err := vc.vehicleService.SearchAndGetVehiclesWithActiveParking(plateQuery, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search vehicles: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
