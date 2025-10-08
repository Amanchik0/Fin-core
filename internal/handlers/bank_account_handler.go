package handlers

import (
	"justTest/internal/models"
	"justTest/internal/services"
	"justTest/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// CreateBankAccount godoc
// @Summary Create a new bank account
// @Description Create a new bank account for the authenticated user
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param request body models.CreateBankAccountRequest true "Bank account creation request"
// @Success 201 {object} models.BankAccount
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /bankAccounts [post]
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

// GetBankAccounts godoc
// @Summary Get all bank accounts
// @Description Get all bank accounts for the authenticated user
// @Tags bank-accounts
// @Produce json
// @Success 200 {array} models.BankAccount
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /bankAccounts [get]
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

// GetBankAccount godoc
// @Summary Get a specific bank account
// @Description Get a specific bank account by ID for the authenticated user
// @Tags bank-accounts
// @Produce json
// @Param bank_account_id path int true "Bank Account ID"
// @Success 200 {object} models.BankAccount
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Bank account not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /bankAccounts/{bank_account_id} [get]
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

// ActivateBankAccount godoc
// @Summary Activate a bank account
// @Description Activate a previously deactivated bank account by ID
// @Tags bank-accounts
// @Produce json
// @Param bank_account_id path int true "Bank Account ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /bankAccounts/{bank_account_id}/activate [put]
func (h *BankAccountHandler) ActivateBankAccount(c *gin.Context) {
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
	err = h.BankAccService.ActivateBankAccount(userID, bankAccountID)
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
			"error":   "failed to activate bank account",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "bank account activated",
	})
}

// DeactivateBankAccount godoc
// @Summary Deactivate a bank account
// @Description Deactivate a bank account by ID (keeps transactions but hides account)
// @Tags bank-accounts
// @Produce json
// @Param bank_account_id path int true "Bank Account ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /bankAccounts/{bank_account_id}/deactivate [put]
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

// DeleteBankAccount godoc
// @Summary Delete a bank account
// @Description Delete a bank account by ID (will also delete all related transactions)
// @Tags bank-accounts
// @Produce json
// @Param bank_account_id path int true "Bank Account ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /bankAccounts/{bank_account_id} [delete]
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
			"error":   "failed to delete bank account",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "bank account deleted successfully",
	})
}
