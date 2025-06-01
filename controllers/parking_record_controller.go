package controllers

import (
	"encoding/base64"
	"hello-professor_backend/dtos"
	"hello-professor_backend/models"
	"hello-professor_backend/services"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ParkingRecordController 定義停車記錄控制器
type ParkingRecordController struct {
	parkingRecordService services.ParkingRecordService
	// reportService services.ReportService // 未來可以考慮引入專門的報表服務
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
// @Success 201 {object} dtos.SuccessResponseWithData{data=models.ParkingRecord}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records [post]
func (prc *ParkingRecordController) CreateParkingRecordHandler(c *gin.Context) {
	var record models.ParkingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := prc.parkingRecordService.CreateParkingRecord(&record); err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create parking record: "+err.Error())
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusCreated, "Parking record created successfully.", record)
}

// GetParkingRecordByIDHandler godoc
// @Summary Get a parking record by ID
// @Description Get details of a parking record by its ID
// @Tags parking_records
// @Produce  json
// @Param   id path int true "Parking Record ID"
// @Success 200 {object} dtos.SuccessResponseWithData{data=models.ParkingRecord}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/{id} [get]
func (prc *ParkingRecordController) GetParkingRecordByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid parking record ID format")
		return
	}

	record, err := prc.parkingRecordService.GetParkingRecordByID(uint(id))
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get parking record: "+err.Error())
		return
	}
	if record == nil {
		dtos.SendErrorResponse(c, http.StatusNotFound, "Parking record not found")
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Parking record retrieved successfully.", record)
}

// GetParkingRecordsByLicensePlateHandler godoc
// @Summary Get parking records by License Plate
// @Description Get all parking records associated with a specific License Plate
// @Tags parking_records
// @Produce json
// @Param licensePlate path string true "License Plate"
// @Success 200 {object} dtos.SuccessResponseWithData{data=[]models.ParkingRecord}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/license/{licensePlate} [get]
func (prc *ParkingRecordController) GetParkingRecordsByLicensePlateHandler(c *gin.Context) {
	licensePlate := c.Param("licensePlate")
	if licensePlate == "" {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "License plate cannot be empty")
		return
	}

	records, err := prc.parkingRecordService.GetParkingRecordsByLicensePlate(licensePlate)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get parking records by license plate: "+err.Error())
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Parking records retrieved successfully.", records)
}

// SearchParkingRecordsByLicensePlateHandler godoc
// @Summary Search parking records by License Plate (fuzzy search)
// @Description Search all parking records by a partial or full License Plate (case-insensitive)
// @Tags parking_records
// @Produce json
// @Param q query string true "License Plate Query"
// @Success 200 {object} dtos.SuccessResponseWithData{data=[]models.ParkingRecord}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/search/license [get]
func (prc *ParkingRecordController) SearchParkingRecordsByLicensePlateHandler(c *gin.Context) {
	licensePlateQuery := c.Query("q")
	if licensePlateQuery == "" {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "License plate query parameter 'q' cannot be empty")
		return
	}

	records, err := prc.parkingRecordService.SearchParkingRecordsByLicensePlate(licensePlateQuery)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to search parking records by license plate: "+err.Error())
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Parking records searched successfully.", records)
}

// GetLatestParkingRecordByLicensePlateHandler godoc
// @Summary Get the latest parking record by License Plate
// @Description Get the most recent parking record for a specific License Plate
// @Tags parking_records
// @Produce json
// @Param licensePlate path string true "License Plate"
// @Success 200 {object} dtos.SuccessResponseWithData{data=models.ParkingRecord}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/license/{licensePlate}/latest [get]
func (prc *ParkingRecordController) GetLatestParkingRecordByLicensePlateHandler(c *gin.Context) {
	licensePlate := c.Param("licensePlate")
	if licensePlate == "" {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "License plate cannot be empty")
		return
	}

	record, err := prc.parkingRecordService.GetLatestParkingRecordByLicensePlate(licensePlate)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get latest parking record by license plate: "+err.Error())
		return
	}
	if record == nil {
		dtos.SendErrorResponse(c, http.StatusNotFound, "No parking record found for this license plate")
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Latest parking record retrieved successfully.", record)
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
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid parking record ID format")
		return
	}

	var recordUpdates models.ParkingRecord
	if err := c.ShouldBindJSON(&recordUpdates); err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	recordUpdates.RecordID = uint(id) // 確保更新的是正確的 ID

	if err := prc.parkingRecordService.UpdateParkingRecord(nil, &recordUpdates); err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update parking record: "+err.Error())
		return
	}
	dtos.SendSuccessResponse(c, http.StatusOK, "Parking record updated successfully")
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
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid parking record ID format")
		return
	}

	if err := prc.parkingRecordService.DeleteParkingRecord(uint(id)); err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete parking record: "+err.Error())
		return
	}
	dtos.SendSuccessResponse(c, http.StatusOK, "Parking record deleted successfully")
}

// GetAllParkingRecordsHandler godoc
// @Summary Get all parking records
// @Description Get a list of all parking records, with pagination
// @Tags parking_records
// @Produce  json
// @Param limit query int false "Limit number of parking records returned" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} dtos.SuccessResponseWithData{data=[]models.ParkingRecord}
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
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get all parking records: "+err.Error())
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "All parking records retrieved successfully.", records)
}

// RecordVehicleEntryHandler godoc
// @Summary Record a vehicle entry event
// @Description Records when a vehicle enters the parking lot, accepting license plate and an optional image file.
// @Tags parking_records
// @Accept multipart/form-data
// @Produce json
// @Param licensePlate formData string true "Vehicle License Plate" example:"ABC-1234"
// @Param image formData file false "Optional image of the vehicle/license plate"
// @Success 201 {object} dtos.SuccessResponseWithData{data=models.ParkingRecord}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 409 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/entry [post]
func (prc *ParkingRecordController) RecordVehicleEntryHandler(c *gin.Context) {
	var payload dtos.SimpleEntryPayload
	if err := c.ShouldBind(&payload); err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if payload.LicensePlate == "" {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "License plate cannot be empty")
		return
	}

	var imageBase64 *string
	if payload.Image != nil {
		file, err := payload.Image.Open()
		if err != nil {
			dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to open image file: "+err.Error())
			return
		}
		defer file.Close()

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to read image file: "+err.Error())
			return
		}

		var mimeType string
		if len(payload.Image.Header["Content-Type"]) > 0 {
			mimeType = payload.Image.Header["Content-Type"][0]
		} else {
			mimeType = "application/octet-stream"
		}

		base64Str := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(bytes)
		imageBase64 = &base64Str
	}

	record, err := prc.parkingRecordService.RecordSimpleVehicleEntry(payload.LicensePlate, imageBase64)
	if err != nil {
		if strings.Contains(err.Error(), "vehicle already in parking lot") {
			dtos.SendErrorResponse(c, http.StatusConflict, err.Error())
		} else {
			dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to record vehicle entry: "+err.Error())
		}
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusCreated, "Vehicle entry recorded successfully.", record)
}

// UpdateUserVerifiedLicensePlateHandler godoc
// @Summary Update user-verified license plate for a parking record
// @Description Allows a user to correct or verify the license plate for an existing parking record.
// @Tags parking_records
// @Accept  json
// @Produce  json
// @Param   id path int true "Parking Record ID"
// @Param   license_plate_info body dtos.VerifyLicensePlatePayload true "Verified License Plate Information"
// @Success 200 {object} dtos.SuccessResponseWithData{data=models.ParkingRecord}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/{id}/verify-license-plate [patch]
func (prc *ParkingRecordController) UpdateUserVerifiedLicensePlateHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid parking record ID format")
		return
	}

	var payload dtos.VerifyLicensePlatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if payload.LicensePlate == "" {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "License plate cannot be empty")
		return
	}

	record, err := prc.parkingRecordService.UpdateUserVerifiedLicensePlate(uint(id), payload.LicensePlate)
	if err != nil {
		if err.Error() == "parking record not found" {
			dtos.SendErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update verified license plate: "+err.Error())
		}
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "User verified license plate updated successfully.", record)
}

// RecordVehicleExitHandler godoc
// @Summary Record a vehicle exit event
// @Description Records when a vehicle exits the parking lot. Checks for payment status.
// @Tags parking_records
// @Accept  json
// @Produce  json
// @Param   exit_info body dtos.SimpleEntryPayload true "Vehicle Exit Information (License Plate Only)"
// @Success 200 {object} dtos.SuccessResponseWithData{data=models.ParkingRecord}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 402 {object} dtos.ErrorResponseWithRecord
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /parking-records/exit [post]
func (prc *ParkingRecordController) RecordVehicleExitHandler(c *gin.Context) {
	var payload dtos.SimpleEntryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if payload.LicensePlate == "" {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "License plate cannot be empty")
		return
	}

	record, err := prc.parkingRecordService.RecordVehicleExit(payload.LicensePlate)
	if err != nil {
		if strings.HasPrefix(err.Error(), "payment_required:") {
			response := dtos.ErrorResponseWithRecord{
				Error: err.Error(),
			}
			if record != nil {
				response.ParkingRecordID = record.RecordID
				response.LicensePlate = record.LicensePlate
				response.CalculatedAmount = record.CalculatedAmount
				response.PaymentStatus = record.PaymentStatus
				response.EntryTime = record.EntryTime
			}
			c.JSON(http.StatusPaymentRequired, response)
		} else if strings.Contains(err.Error(), "no active parking record found") {
			dtos.SendErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to record vehicle exit: "+err.Error())
		}
		return
	}

	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Vehicle exit recorded successfully.", record)
}

// PrepareParkingRecordForPaymentHandler godoc
// @Summary Prepare a parking record for payment by calculating/retrieving its fee
// @Description Calculates and stores the parking fee if not already calculated for an active parking record. Returns the record with payment details.
// @Tags parking_records
// @Produce  json
// @Param   id path int true "Parking Record ID"
// @Success 200 {object} dtos.SuccessResponseWithData{data=models.ParkingRecord} "Successfully calculated/retrieved fee, record ready for payment"
// @Failure 400 {object} dtos.ErrorResponse "Invalid Record ID or record already exited/paid"
// @Failure 404 {object} dtos.ErrorResponse "Parking Record not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /parking-records/{id}/prepare-payment [post]
func (prc *ParkingRecordController) PrepareParkingRecordForPaymentHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid parking record ID format")
		return
	}

	record, err := prc.parkingRecordService.PrepareParkingRecordForPayment(uint(id))
	if err != nil {
		if strings.HasPrefix(err.Error(), "vehicle_exited:") || strings.HasPrefix(err.Error(), "already_paid:") {
			dtos.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		} else if strings.Contains(err.Error(), "not found") {
			dtos.SendErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			dtos.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Parking fee prepared successfully.", record)
}

// PayForParkingRecordHandler handles the request to pay for a parking record.
// @Summary Pay for a parking record
// @Description Marks a parking record as paid and ideally creates a transaction record.
// @Tags Parking Records
// @Accept json
// @Produce json
// @Param id path uint true "Parking Record ID"
// @Param paymentPayload body dtos.ParkingPaymentPayload true "Payment Details"
// @Success 200 {object} dtos.SuccessResponseWithData{data=dtos.ParkingRecordWithTransactionResponse} "Payment successful"
// @Failure 400 {object} dtos.ErrorResponse "Invalid request (e.g., validation error, amount mismatch)"
// @Failure 402 {object} dtos.ErrorResponse "Payment required conditions not met (e.g., fee not calculated, already paid, vehicle exited)"
// @Failure 404 {object} dtos.ErrorResponse "Parking record not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /parking-records/{id}/pay [post]
func (prc *ParkingRecordController) PayForParkingRecordHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid parking record ID format")
		return
	}

	var payload dtos.ParkingPaymentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	parkingRecord, transaction, err := prc.parkingRecordService.PayForParkingRecord(uint(id), payload)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			dtos.SendErrorResponse(c, http.StatusNotFound, errMsg)
		} else if strings.Contains(errMsg, "amount_mismatch") || strings.Contains(errMsg, "fee_not_calculated") || strings.Contains(errMsg, "already_paid") || strings.Contains(errMsg, "vehicle_exited") {
			dtos.SendErrorResponse(c, http.StatusPaymentRequired, errMsg)
		} else {
			dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to process payment: "+errMsg)
		}
		return
	}

	response := dtos.ParkingRecordWithTransactionResponse{
		ParkingRecord: *parkingRecord,
		Transaction:   *transaction,
	}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Payment processed successfully.", response)
}

// --- 新增報表 API ---

// parseTimeRangeParameters 輔助函數，用於解析時間範圍參數
func parseTimeRangeParameters(c *gin.Context) (startTime, endTime *time.Time, err error) {
	startTimeStr := c.Query("startTime") // e.g., "2023-01-01T00:00:00Z"
	endTimeStr := c.Query("endTime")     // e.g., "2023-01-31T23:59:59Z"

	if startTimeStr != "" {
		st, errP := time.Parse(time.RFC3339, startTimeStr)
		if errP != nil {
			return nil, nil, errP
		}
		startTime = &st
	}
	if endTimeStr != "" {
		et, errP := time.Parse(time.RFC3339, endTimeStr)
		if errP != nil {
			return nil, nil, errP
		}
		endTime = &et
	}
	return startTime, endTime, nil
}

// GetTotalParkingCountHandler godoc
// @Summary Get total parking count within a time range
// @Description Retrieves the total number of parking events (vehicle entries).
// @Tags reports
// @Produce json
// @Param startTime query string false "Start time for the report (RFC3339 format, e.g., 2023-01-01T00:00:00Z)"
// @Param endTime query string false "End time for the report (RFC3339 format, e.g., 2023-01-31T23:59:59Z)"
// @Success 200 {object} dtos.SuccessResponseWithData{data=dtos.TotalParkingCountResponse}
// @Failure 400 {object} dtos.ErrorResponse "Invalid time format"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /reports/traffic/total-count [get]
func (prc *ParkingRecordController) GetTotalParkingCountHandler(c *gin.Context) {
	startTime, endTime, err := parseTimeRangeParameters(c)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid time format: "+err.Error())
		return
	}

	count, err := prc.parkingRecordService.GetTotalParkingCount(startTime, endTime)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get total parking count: "+err.Error())
		return
	}
	responseData := dtos.TotalParkingCountResponse{TotalCount: count}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Total parking count retrieved successfully.", responseData)
}

// GetTotalRevenueHandler godoc
// @Summary Get total revenue from parking fees within a time range
// @Description Retrieves the total revenue collected from parking fees.
// @Tags reports
// @Produce json
// @Param startTime query string false "Start time for the report (RFC3339 format, e.g., 2023-01-01T00:00:00Z)"
// @Param endTime query string false "End time for the report (RFC3339 format, e.g., 2023-01-31T23:59:59Z)"
// @Success 200 {object} dtos.SuccessResponseWithData{data=dtos.TotalRevenueResponse}
// @Failure 400 {object} dtos.ErrorResponse "Invalid time format"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /reports/revenue/total [get]
func (prc *ParkingRecordController) GetTotalRevenueHandler(c *gin.Context) {
	startTime, endTime, err := parseTimeRangeParameters(c)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid time format: "+err.Error())
		return
	}

	revenue, err := prc.parkingRecordService.GetTotalRevenue(startTime, endTime)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get total revenue: "+err.Error())
		return
	}
	responseData := dtos.TotalRevenueResponse{TotalRevenue: revenue, Currency: "TWD"} // 假設幣別為 TWD
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Total revenue retrieved successfully.", responseData)
}

// GetImageAttachmentRateHandler godoc
// @Summary Get the rate of parking entries with images
// @Description Calculates the percentage of vehicle entries that have an associated image.
// @Tags reports
// @Produce json
// @Param startTime query string false "Start time for the report (RFC3339 format, e.g., 2023-01-01T00:00:00Z)"
// @Param endTime query string false "End time for the report (RFC3339 format, e.g., 2023-01-31T23:59:59Z)"
// @Success 200 {object} dtos.SuccessResponseWithData{data=dtos.ImageAttachmentRateResponse}
// @Failure 400 {object} dtos.ErrorResponse "Invalid time format"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /reports/operations/image-attachment-rate [get]
func (prc *ParkingRecordController) GetImageAttachmentRateHandler(c *gin.Context) {
	startTime, endTime, err := parseTimeRangeParameters(c)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusBadRequest, "Invalid time format: "+err.Error())
		return
	}

	rateResponse, err := prc.parkingRecordService.GetImageAttachmentRate(startTime, endTime)
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get image attachment rate: "+err.Error())
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Image attachment rate retrieved successfully.", rateResponse)
}

// GetAvailableParkingSpotsHandler godoc
// @Summary Get available parking spots
// @Description Retrieves the total capacity, occupied spots, and available spots in the parking lot.
// @Tags reports
// @Produce json
// @Success 200 {object} dtos.SuccessResponseWithData{data=dtos.AvailableSpotsResponse}
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /reports/parking-lot/available-spots [get]
func (prc *ParkingRecordController) GetAvailableParkingSpotsHandler(c *gin.Context) {
	spotsResponse, err := prc.parkingRecordService.GetAvailableParkingSpots()
	if err != nil {
		dtos.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get available parking spots: "+err.Error())
		return
	}
	dtos.SendSuccessResponseWithData(c, http.StatusOK, "Available parking spots retrieved successfully.", spotsResponse)
}
