package handlers

import (
	"justTest/internal/models"
	"justTest/internal/services"
	"justTest/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService}

}
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	limit, offset := utils.GetPaginationParams(1, 20)
	notifications, err := h.notificationService.GetUserNotifications(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    notifications,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	//userID, ok := utils.GetUserID(c)
	//if !ok {
	//	return
	//}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	if err := h.notificationService.MarkAsRead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notification marked as read",
	})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	if err := h.notificationService.MarkAllAsRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All notifications marked as read",
	})
}

func (h *NotificationHandler) GetSettings(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	settings, err := h.notificationService.GetSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

func (h *NotificationHandler) SaveSettings(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}
	var req models.SaveSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	settings, err := h.notificationService.SaveSettings(&models.UserNotificationSettings{
		UserID:               userID,
		BudgetAlertsEnabled:  req.BudgetAlertsEnabled,
		BalanceAlertsEnabled: req.BalanceAlertsEnabled,
		BudgetWarningPercent: req.BudgetWarningPercent,
		LowBalanceThreshold:  req.LowBalanceThreshold,
		PreferredChannel:     req.PreferredChannel,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
		"message": "Settings updated successfully",
	})
}
