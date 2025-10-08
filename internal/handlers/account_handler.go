package handlers

import (
	"justTest/internal/models"
	"justTest/internal/services"
	"justTest/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// CreateAccount POST /api/v1/accounts
// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new financial account for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body models.CreateAccountRequest true "Account creation request"
// @Success 201 {object} models.AccountResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /account [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {

	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	var req models.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false,
			"error": "invalid", "details": err.Error()})
		return
	}

	account, err := h.accountService.CreateAccount(userID, req.DisplayName)
	if err != nil {
		if err.Error() == "Account already exists" {
			c.JSON(http.StatusConflict, gin.H{"success": false, "error": "Account already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false,
			"error":   "failed to create account",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Account created successfully",
		"data": models.AccountResponse{
			ID:          account.ID,
			UserID:      account.UserID,
			DisplayName: account.DisplayName,
			Name:        account.Name,
			Timezone:    account.Timezone,
			IsActive:    account.IsActive,
			CreatedAt:   account.CreatedAt,
			UpdatedAt:   account.UpdatedAt,
		},
	})
}

// GetAccount godoc
// @Summary Get user account
// @Description Get the authenticated user's account information
// @Tags accounts
// @Produce json
// @Success 200 {object} models.AccountResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /account [get]
func (h *AccountHandler) GetAccount(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	account, err := h.accountService.GetUserAccount(userID)
	if err != nil {
		if err.Error() == "Account not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Account not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false,
			"error":   "failed to get account",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": models.AccountResponse{
			ID:          account.ID,
			UserID:      account.UserID,
			DisplayName: account.DisplayName,
			Name:        account.Name,
			Timezone:    account.Timezone,
			IsActive:    account.IsActive,
			CreatedAt:   account.CreatedAt,
			UpdatedAt:   account.UpdatedAt,
		},
	})
}

//func (h *AccountHandler) GetAccountSummary(c *gin.Context) {
//	userID, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{
//
//			"success": false,
//			"error":   "unauthorized",
//		})
//		return
//	}
//	userIDInt, ok := userID.(int64)
//	if !ok {
//		c.JSON(http.StatusUnauthorized, gin.H{
//			"success": false,
//			"error":   "invalid user id format",
//		})
//		return
//	}
//	account, err := h.accountService.GetUserAccount(userIDInt)
//	if err != nil {
//		if err.Error() == "Account not found" {
//			c.JSON(http.StatusNotFound, gin.H{
//				"success": false,
//				"error":   "Account not found",
//			})
//			return
//		}
//		c.JSON(http.StatusInternalServerError, gin.H{"success": false,
//			"error": "failed to get account",
//		})
//	}
//}

// UpdateAccount PUT /api/v1/accounts

// GetAccountByID GET /api/v1/accounts/:id (для админа или внутренних нужд)
