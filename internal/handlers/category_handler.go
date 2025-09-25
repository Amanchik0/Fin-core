package handlers

import (
	"github.com/gin-gonic/gin"
	"justTest/internal/models"
	"justTest/internal/services"
	"justTest/internal/utils"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}

}

// create category psot
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
