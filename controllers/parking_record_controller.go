package controllers

import (
	"hello-professor_backend/dtos"
	"hello-professor_backend/models"
	"hello-professor_backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParkingRecordController 定義停車記錄控制器
type ParkingRecordController struct {
	parkingRecordService services.ParkingRecordService
}

// NewParkingRecordController 建立一個新的 ParkingRecordController 實例
func NewParkingRecordController(prs services.ParkingRecordService) *ParkingRecordController {
	return &ParkingRecordController{parkingRecordService: prs}
}

// CreateParkingRecordHandler godoc
// @Summary Create a new parking record
// @Description Add a new parking record to the system. This is a general endpoint, for specific entry/exit events, use /parking-records/entry and /parking-records/exit.
// @Tags parking_records
// @Accept  json
// @Produce  json
// @Param   parking_record_info body models.ParkingRecord true "Parking Record Information"
// @Success 201 {object} models.ParkingRecord
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records [post]
func (prc *ParkingRecordController) CreateParkingRecordHandler(c *gin.Context) {
	var record models.ParkingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := prc.parkingRecordService.CreateParkingRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create parking record: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, record)
}

// GetParkingRecordByIDHandler godoc
// @Summary Get a parking record by ID
// @Description Get details of a parking record by its ID
// @Tags parking_records
// @Produce  json
// @Param   id path int true "Parking Record ID"
// @Success 200 {object} models.ParkingRecord
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/{id} [get]
func (prc *ParkingRecordController) GetParkingRecordByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parking record ID format"})
		return
	}

	record, err := prc.parkingRecordService.GetParkingRecordByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parking record: " + err.Error()})
		return
	}
	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parking record not found"})
		return
	}
	c.JSON(http.StatusOK, record)
}

// GetParkingRecordsByVehicleIDHandler godoc
// @Summary Get parking records by Vehicle ID
// @Description Get all parking records associated with a specific Vehicle ID
// @Tags parking_records
// @Produce json
// @Param vehicleID path int true "Vehicle ID"
// @Success 200 {array} models.ParkingRecord
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/vehicle/{vehicleID} [get]
func (prc *ParkingRecordController) GetParkingRecordsByVehicleIDHandler(c *gin.Context) {
	vehicleIDStr := c.Param("vehicleID")
	vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID format"})
		return
	}

	records, err := prc.parkingRecordService.GetParkingRecordsByVehicleID(uint(vehicleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parking records by vehicle ID: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

// GetLatestParkingRecordByVehicleIDHandler godoc
// @Summary Get the latest parking record by Vehicle ID
// @Description Get the most recent parking record for a specific Vehicle ID
// @Tags parking_records
// @Produce json
// @Param vehicleID path int true "Vehicle ID"
// @Success 200 {object} models.ParkingRecord
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/vehicle/{vehicleID}/latest [get]
func (prc *ParkingRecordController) GetLatestParkingRecordByVehicleIDHandler(c *gin.Context) {
	vehicleIDStr := c.Param("vehicleID")
	vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID format"})
		return
	}

	record, err := prc.parkingRecordService.GetLatestParkingRecordByVehicleID(uint(vehicleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest parking record by vehicle ID: " + err.Error()})
		return
	}
	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No parking record found for this vehicle"})
		return
	}
	c.JSON(http.StatusOK, record)
}

// UpdateParkingRecordHandler godoc
// @Summary Update an existing parking record
// @Description Update details of an existing parking record by its ID. Can be used for manual adjustments.
// @Tags parking_records
// @Accept  json
// @Produce  json
// @Param   id path int true "Parking Record ID"
// @Param   parking_record_update body models.ParkingRecord true "Parking Record Update Information"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/{id} [put]
func (prc *ParkingRecordController) UpdateParkingRecordHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parking record ID format"})
		return
	}

	var recordUpdates models.ParkingRecord
	if err := c.ShouldBindJSON(&recordUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	recordUpdates.RecordID = uint(id) // 確保更新的是正確的 ID

	if err := prc.parkingRecordService.UpdateParkingRecord(&recordUpdates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parking record: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Parking record updated successfully"})
}

// DeleteParkingRecordHandler godoc
// @Summary Delete a parking record by ID
// @Description Remove a parking record from the system by its ID
// @Tags parking_records
// @Produce  json
// @Param   id path int true "Parking Record ID"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/{id} [delete]
func (prc *ParkingRecordController) DeleteParkingRecordHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parking record ID format"})
		return
	}

	if err := prc.parkingRecordService.DeleteParkingRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete parking record: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Parking record deleted successfully"})
}

// GetAllParkingRecordsHandler godoc
// @Summary Get all parking records
// @Description Get a list of all parking records, with pagination
// @Tags parking_records
// @Produce  json
// @Param limit query int false "Limit number of parking records returned" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.ParkingRecord
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records [get]
func (prc *ParkingRecordController) GetAllParkingRecordsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	records, err := prc.parkingRecordService.GetAllParkingRecords(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all parking records: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

// RecordVehicleEntryHandler godoc
// @Summary Record a vehicle entry event
// @Description Records when a vehicle enters the parking lot.
// @Tags parking_records
// @Accept  json
// @Produce  json
// @Param   entry_info body dtos.VehicleEntryExitPayload true "Vehicle Entry Information"
// @Success 201 {object} models.ParkingRecord
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/entry [post]
func (prc *ParkingRecordController) RecordVehicleEntryHandler(c *gin.Context) {
	var payload dtos.VehicleEntryExitPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	record, err := prc.parkingRecordService.RecordVehicleEntry(payload.VehicleID, payload.SensorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vehicle entry: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, record)
}

// RecordVehicleExitHandler godoc
// @Summary Record a vehicle exit event
// @Description Records when a vehicle exits the parking lot and calculates fees.
// @Tags parking_records
// @Accept  json
// @Produce  json
// @Param   exit_info body dtos.VehicleEntryExitPayload true "Vehicle Exit Information"
// @Success 200 {object} models.ParkingRecord
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/exit [post]
func (prc *ParkingRecordController) RecordVehicleExitHandler(c *gin.Context) {
	var payload dtos.VehicleEntryExitPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	record, err := prc.parkingRecordService.RecordVehicleExit(payload.VehicleID, payload.SensorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vehicle exit: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, record)
}
