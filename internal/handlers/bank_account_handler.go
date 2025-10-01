package handlers

import (
	"github.com/gin-gonic/gin"
	"justTest/internal/models"
	"justTest/internal/services"
	"justTest/internal/utils"
	"net/http"
	"strconv"
)

type BankAccountHandler struct {
	BankAccService *services.BankAccService
}

func NewBankAccountHandler(bankAccService *services.BankAccService) *BankAccountHandler {
	return &BankAccountHandler{
		BankAccService: bankAccService,
	}
}

// create bank account  api/v1/bankaccount

func (h *BankAccountHandler) CreateBankAccount(c *gin.Context) {

	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	var req models.CreateBankAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	bankAccount, err := h.BankAccService.CreateBankAccount(
		userID,
		req.Name,
		req.Currency,
		req.AccountType,
		req.BankName,
	)
	if err != nil {
		if err.Error() == "bank account already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "bank account already exists",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "failed to create bank account",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "bank account created",
		"data":    bankAccount,
	})

}

// get bank account

func (h *BankAccountHandler) GetBankAccounts(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	bankAccounts, err := h.BankAccService.GetBankAccountsByAccountID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"error":   "failed to get bank accounts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bankAccounts,
	})
}

func (h *BankAccountHandler) GetBankAccount(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	bankAccountIDStr := c.Param("bank_account_id")
	bankAccountID, err := strconv.ParseInt(bankAccountIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "invalid bank account id",
		})
		return
	}

	bankAccount, err := h.BankAccService.GetBankAccount(userID, bankAccountID)
	if err != nil {
		if err.Error() == "invalid user account" {
			c.JSON(http.StatusNotFound, gin.H{
				"status": false,
				"error":  "account not found",
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "failed to get bank account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bankAccount,
	})

}

func (h *BankAccountHandler) DeactivateBankAccount(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	bankAccountIDStr := c.Param("bank_account_id")
	bankAccountID, err := strconv.ParseInt(bankAccountIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "invalid bank account id",
		})
		return
	}
	err = h.BankAccService.DeActiveBankAccount(userID, bankAccountID)
	if err != nil {
		if err.Error() == "invlaid user acc" {
			c.JSON(http.StatusForbidden, gin.H{
				"status": false,
				"error":  "access denied",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"error":   "failed to deactivate bank account",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "bank account deactivated",
	})
}

func (h *BankAccountHandler) DeleteBankAccount(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	bankAccountIDStr := c.Param("bank_account_id")
	bankAccountID, err := strconv.ParseInt(bankAccountIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "invalid bank account id",
		})
		return
	}
	err = h.BankAccService.DeleteBankAccount(userID, bankAccountID)
	if err != nil {
		if err.Error() == "invlaid user acc" {
			c.JSON(http.StatusForbidden, gin.H{
				"status": false,
				"error":  "access denied",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"error":   "failed to deactivate bank account",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "bank account deactivated",
	})
}
