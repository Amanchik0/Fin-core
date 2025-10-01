package handlers

import (
	"github.com/gin-gonic/gin"
	"justTest/internal/models"
	"justTest/internal/services"
	"justTest/internal/utils"
	"net/http"
	"strconv"
)

// ================================
// TRANSACTION HANDLER
// ================================

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
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

	response := h.transactionToResponse(transaction)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Transaction created",
	})
}

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
