package handlers

import (
	"justTest/internal/models"
	"justTest/internal/services"
	"justTest/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}

}

// create category psot
// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new expense or income category for the authenticated user
// @Tags categories
// @Accept json
// @Produce json
// @Param request body models.CreateCategoryRequest true "Category creation request"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	var req models.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	category := &models.Category{
		Name:  req.Name,
		Type:  req.Type,
		Color: req.Color,
		Icon:  req.Icon,
	}
	newCategory, err := h.categoryService.CreateCategory(userID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    newCategory,
	})

}

// GetCategoryByID godoc
// @Summary Get a specific category
// @Description Get a specific category by ID for the authenticated user
// @Tags categories
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /categories/{category_id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}
	category, err := h.categoryService.GetByCategoryID(userID, categoryID)
	if err != nil {
		if err.Error() == "ibvalid user account" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    category,
	})

}

// GetByAccountID godoc
// @Summary Get all categories
// @Description Get all categories for the authenticated user's account
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /categories [get]
func (h *CategoryHandler) GetByAccountID(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	// Получаем account_id через user_id
	account, err := h.categoryService.GetAccountByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to get user account",
			"details": err.Error(),
		})
		return
	}

	accountCategories, err := h.categoryService.GetByAccountID(account.ID)
	if err != nil {
		if err.Error() == "ibvalid user account" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    accountCategories,
	})
}

// DeleteCategoryByID godoc
// @Summary Delete a category
// @Description Delete a category by ID for the authenticated user
// @Tags categories
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /categories/{category_id} [delete]
func (h *CategoryHandler) DeleteCategoryByID(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}
	err = h.categoryService.DeleteCategory(userID, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    nil,
		"message": "category deleted successfully",
	})
}

// нужно потом update прописать
