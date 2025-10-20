package services

import (
	"justTest/internal/events"
	"justTest/internal/interfaces"
	"justTest/internal/models"
)

type NotificationService struct {
	notificationRepo interfaces.NotificationRepository
	settingsRepo     UserNotificationSettingsRepository

	budgetService  *BudgetService
	bankAccService *BankAccountService
	publisher      *events.Publisher
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// CheckBudgetAlerts - проверить превышение бюджетов
func (s *NotificationService) CheckBudgetAlerts(userID string) ([]*models.BudgetAlert, error) {
	// TODO: Реализовать проверку превышения бюджетов
	return nil, nil
}

// CheckLowBalance - проверить низкий баланс
func (s *NotificationService) CheckLowBalance(userID string) ([]*models.BalanceAlert, error) {
	// TODO: Реализовать проверку низкого баланса
	return nil, nil
}
