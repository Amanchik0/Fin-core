package handlers

import (
	"justTest/internal/models"
	"justTest/internal/services"
	"justTest/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BudgetHandler struct {
	budgetService *services.BudgetService
}

func NewBudgetHandler(budgetService *services.BudgetService) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
	}
}

// CreateBudget godoc
// @Summary Create a new budget
// @Description Create a new monthly budget for a category
// @Tags budgets
// @Accept json
// @Produce json
// @Param request body models.CreateBudgetRequest true "Budget creation request"
// @Success 201 {object} models.Budget
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /budgets [post]
func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	var req models.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error()})
		return
	}

	budget, err := h.budgetService.CreateBudget(userID, &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    budget,
		"message": "budget created",
	})
}

// GetBudgets godoc
// @Summary Get all budgets for a month
// @Description Get all budgets with status for a specific month
// @Tags budgets
// @Produce json
// @Param year query int true "Year" default(2024)
// @Param month query int true "Month (1-12)" default(10)
// @Success 200 {array} models.BudgetWithStatus
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /budgets [get]
func (h *BudgetHandler) GetBudgets(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid year parameter",
		})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid month parameter",
		})
		return
	}
	if month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "month must be between 1 and 12",
		})
		return
	}

	budgets, err := h.budgetService.GetBudgets(userID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    budgets,
		"message": "budgets found",
	})

}

// GetBudgetStatus godoc
// @Summary Get budget status for a category
// @Description Get budget status (spent/remaining) for a specific category and month
// @Tags budgets
// @Produce json
// @Param category_id path int true "Category ID"
// @Param year query int true "Year" default(2024)
// @Param month query int true "Month (1-12)" default(10)
// @Success 200 {object} models.BudgetStatus
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Budget not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /budgets/{category_id}/status [get]
func (h *BudgetHandler) GetBudgetStatus(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	categoryIdStr := c.Param("category_id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid category id parameter",
		})
		return
	}
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid year parameter",
		})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid month parameter",
		})
		return
	}
	if month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "month must be between 1 and 12",
		})
		return
	}

	budget, err := h.budgetService.GetBudgetStatus(userID, categoryId, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    budget,
		"message": "budget status",
	})

}

// GetBudgetSummary godoc
// @Summary Get budget summary for a month
// @Description Get overall budget summary (total planned vs spent) for a month
// @Tags budgets
// @Produce json
// @Param year query int true "Year" default(2024)
// @Param month query int true "Month (1-12)" default(10)
// @Success 200 {object} models.BudgetSummary
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /budgets/summary [get]
func (h *BudgetHandler) GetBudgetSummary(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return

	}

	yearStr := c.Query("year")
	monthStr := c.Query("month")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid year parameter",
		})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid month parameter",
		})
		return
	}
	if month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "month must be between 1 and 12",
		})
		return
	}

	budget, err := h.budgetService.GetBudgetSummary(userID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "something bad happened",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    budget,
		"message": "budget summary",
	})

}
