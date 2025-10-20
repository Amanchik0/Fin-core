package handlers

import (
	"justTest/internal/events"
	"justTest/internal/models"
	"justTest/internal/services"
	"justTest/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
	publisher          *events.Publisher
}

func NewTransactionHandler(transactionService *services.TransactionService, publisher *events.Publisher) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		publisher:          publisher,
	}
}

func (h *TransactionHandler) transactionToResponse(transaction *models.Transaction) models.TransactionResponse {
	return models.TransactionResponse{
		ID:              transaction.ID,
		BankAccountID:   transaction.BankAccountID,
		CategoryID:      transaction.CategoryID,
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		TransactionType: transaction.TransactionType,
		Date:            transaction.Date.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:       transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       transaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new income, expense, or transfer transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body models.CreateTransactionRequest true "Transaction creation request"
// @Success 201 {object} models.TransactionResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.CreateTransaction(
		userID,
		req.BankAccountID,
		req.Amount,
		req.Description,
		req.CategoryID,
		req.TransactionType,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.publisher.PublishTransactionCreated(events.TransactionCreatedEvent{
		TransactionID: transaction.ID,
		UserID:        userID,
		CategoryID:    *transaction.CategoryID,
		Amount:        transaction.Amount,
		Description:   transaction.Description,
		Timestamp:     time.Now(),
	})
	response := h.transactionToResponse(transaction)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Transaction created",
	})
}

// TransferBetweenAccounts godoc
// @Summary Transfer money between bank accounts
// @Description Create a transfer transaction between two bank accounts
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body models.TransferRequest true "Transfer request"
// @Success 201 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /transfer [post]
func (h *TransactionHandler) TransferBetweenAccounts(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.transactionService.TransferBetweenAccounts(
		userID,
		req.FromAccountID,
		req.ToAccountID,
		req.Description,
		req.Amount,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transfer completed successfully",
	})
}

// GetTransactionHistory godoc
// @Summary Get transaction history by bank account
// @Description Get all transactions for a specific bank account
// @Tags transactions
// @Produce json
// @Param account_id path int true "Bank Account ID"
// @Success 200 {array} models.TransactionResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /account/{account_id}/transactions [get]
func (h *TransactionHandler) GetTransactionHistory(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	transactions, err := h.transactionService.GetTransactionHistory(userID, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []models.TransactionResponse
	for _, transaction := range transactions {
		response = append(response, h.transactionToResponse(transaction))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetAllTransactions godoc
// @Summary Get all transactions
// @Description Get all transactions for the authenticated user
// @Tags transactions
// @Produce json
// @Success 200 {array} models.TransactionResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	transactions, err := h.transactionService.GetAllTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []models.TransactionResponse
	for _, transaction := range transactions {
		response = append(response, h.transactionToResponse(transaction))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetBankAccountBalance godoc
// @Summary Get bank account balance
// @Description Get the current balance of a specific bank account
// @Tags transactions
// @Produce json
// @Param account_id path int true "Bank Account ID"
// @Success 200 {object} map[string]interface{} "Balance information"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /bank_accounts/{account_id}/balance [get]
func (h *TransactionHandler) GetBankAccountBalance(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	balance, err := h.transactionService.GetBankAccountBalance(userID, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"account_id": accountID,
			"balance":    balance,
		},
	})
}

// GetTransaction godoc
// @Summary Get a specific transaction
// @Description Get a specific transaction by ID for the authenticated user
// @Tags transactions
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} models.TransactionResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	transactionIDStr := c.Param("id")
	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(userID, transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := h.transactionToResponse(transaction)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *TransactionHandler) GetAllTransactionsByCategoryID(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)

	transactions, err := h.transactionService.GetAllTransactionsByCategoryID(userID, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []models.TransactionResponse
	for _, transaction := range transactions {
		response = append(response, h.transactionToResponse(transaction))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
