package controllers

import (
	"hello-professor_backend/models"
	"hello-professor_backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TransactionController 定義交易控制器
type TransactionController struct {
	transactionService services.TransactionService
}

// NewTransactionController 建立一個新的 TransactionController 實例
func NewTransactionController(ts services.TransactionService) *TransactionController {
	return &TransactionController{transactionService: ts}
}

// CreateTransactionHandler godoc
// @Summary Create a new transaction
// @Description Add a new transaction to the system
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param   transaction_info body models.Transaction true "Transaction Information"
// @Success 201 {object} models.Transaction
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /transactions [post]
func (tc *TransactionController) CreateTransactionHandler(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := tc.transactionService.CreateTransaction(nil, &transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, transaction)
}

// GetTransactionByIDHandler godoc
// @Summary Get a transaction by ID
// @Description Get details of a transaction by its ID
// @Tags transactions
// @Produce  json
// @Param   id path int true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /transactions/{id} [get]
func (tc *TransactionController) GetTransactionByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}

	transaction, err := tc.transactionService.GetTransactionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transaction: " + err.Error()})
		return
	}
	if transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

// GetTransactionsByParkingRecordIDHandler godoc
// @Summary Get transactions by ParkingRecord ID
// @Description Get all transactions associated with a specific ParkingRecord ID
// @Tags transactions
// @Produce json
// @Param parkingRecordID path int true "Parking Record ID"
// @Success 200 {array} models.Transaction
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /transactions/parking/{parkingRecordID} [get]
func (tc *TransactionController) GetTransactionsByParkingRecordIDHandler(c *gin.Context) {
	parkingRecordIDStr := c.Param("parkingRecordID")
	parkingRecordID, err := strconv.ParseUint(parkingRecordIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parking record ID format"})
		return
	}

	transactions, err := tc.transactionService.GetTransactionsByParkingRecordID(uint(parkingRecordID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions by parking record ID: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// UpdateTransactionHandler godoc
// @Summary Update an existing transaction
// @Description Update details of an existing transaction by its ID
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param   id path int true "Transaction ID"
// @Param   transaction_update body models.Transaction true "Transaction Update Information"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /transactions/{id} [put]
func (tc *TransactionController) UpdateTransactionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}

	var transactionUpdates models.Transaction
	if err := c.ShouldBindJSON(&transactionUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	transactionUpdates.TransactionID = uint(id) // 確保更新的是正確的 ID

	if err := tc.transactionService.UpdateTransaction(&transactionUpdates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}

// DeleteTransactionHandler godoc
// @Summary Delete a transaction by ID
// @Description Remove a transaction from the system by its ID
// @Tags transactions
// @Produce  json
// @Param   id path int true "Transaction ID"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /transactions/{id} [delete]
func (tc *TransactionController) DeleteTransactionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}

	if err := tc.transactionService.DeleteTransaction(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

// GetAllTransactionsHandler godoc
// @Summary Get all transactions
// @Description Get a list of all transactions, with pagination
// @Tags transactions
// @Produce  json
// @Param limit query int false "Limit number of transactions returned" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.Transaction
// @Failure 500 {object} dtos.ErrorResponse
// @Router /transactions [get]
func (tc *TransactionController) GetAllTransactionsHandler(c *gin.Context) {
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

	transactions, err := tc.transactionService.GetAllTransactions(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all transactions: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
