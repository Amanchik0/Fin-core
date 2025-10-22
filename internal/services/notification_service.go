package services

import (
	"database/sql"
	"errors"
	"fmt"
	"justTest/internal/interfaces"
	"justTest/internal/models"
	"justTest/internal/models/events"
	"log"
	"time"
)

type NotificationService struct {
	notificationRepo interfaces.NotificationRepository
	settingsRepo     interfaces.UserNotificationSettingsRepository

	publisher interface{}
}

func NewNotificationService(
	notificationRepo interfaces.NotificationRepository,
	settingsRepo interfaces.UserNotificationSettingsRepository,
	publisher interface{},
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		settingsRepo:     settingsRepo,
		publisher:        publisher,
	}
}

func (s *NotificationService) GetUserNotifications(userID string, limit, offset int) ([]*models.Notification, error) {
	return s.notificationRepo.GetUserNotifications(userID, limit, offset)
}

func (s *NotificationService) MarkAsRead(id int64) error {
	return s.notificationRepo.MarkAsRead(id)
}

func (s *NotificationService) MarkAllAsRead(userID string) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}
func (s *NotificationService) GetSettings(userID string) (*models.UserNotificationSettings, error) {
	settings, err := s.settingsRepo.GetSettings(userID)
	if err == sql.ErrNoRows {
		defaultSettings := &models.UserNotificationSettings{
			UserID:               userID,
			BudgetAlertsEnabled:  true,
			BalanceAlertsEnabled: true,
			BudgetWarningPercent: 80,
			LowBalanceThreshold:  1000,
			PreferredChannel:     "email",
		}
		if err := s.settingsRepo.SaveSettings(defaultSettings); err != nil {
			return nil, err
		}
		return defaultSettings, nil
	}
	if err != nil {
		return nil, err
	}
	return settings, nil
}
func (s *NotificationService) SaveSettings(settings *models.UserNotificationSettings) (*models.UserNotificationSettings, error) {
	if settings == nil || settings.UserID == "" {
		return nil, errors.New("invalid settings")
	}
	existing, err := s.settingsRepo.GetSettings(settings.UserID)
	if err == sql.ErrNoRows {
		settings.CreatedAt = time.Now()
		settings.UpdatedAt = time.Now()
		if err := s.settingsRepo.SaveSettings(settings); err != nil {
			return nil, err
		}
		return settings, nil
	}
	if err != nil {
		return nil, err
	}

	existing.BudgetAlertsEnabled = settings.BudgetAlertsEnabled
	existing.BalanceAlertsEnabled = settings.BalanceAlertsEnabled
	existing.BudgetWarningPercent = settings.BudgetWarningPercent
	existing.LowBalanceThreshold = settings.LowBalanceThreshold
	existing.PreferredChannel = settings.PreferredChannel

	return s.settingsRepo.UpdateSettings(existing)
}

func (s *NotificationService) HandleLowBalance(event events.LowBalanceEvent) error {
	settings, err := s.GetSettings(event.UserID)
	if err != nil {
		return err
	}
	if !settings.BalanceAlertsEnabled {
		return nil
	}
	if event.CurrentBalance >= settings.LowBalanceThreshold {
		return nil
	}

	n := &models.Notification{
		UserID:    event.UserID,
		Type:      "low_balance",
		Title:     "Низкий баланс",
		Message:   "Баланс по счету " + event.AccountName + " ниже порога",
		Data:      map[string]interface{}{"bank_account_id": event.BankAccountID, "current_balance": event.CurrentBalance, "threshold": settings.LowBalanceThreshold},
		IsRead:    false,
		Priority:  "high",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.notificationRepo.SaveNotification(n); err != nil {
		return err
	}

	// Опционально — отправить в общую очередь "notification"
	if s.publisher != nil {
		if publisher, ok := s.publisher.(interface {
			PublishNotification(events.NotificationEvent) error
		}); ok {
			_ = publisher.PublishNotification(events.NotificationEvent{
				UserID:    n.UserID,
				Type:      n.Type,
				Title:     n.Title,
				Message:   n.Message,
				Data:      n.Data,
				Timestamp: time.Now(),
				Priority:  n.Priority,
			})
		}
	}
	return nil
}

// TODO: дедупликация, как нибудь потом надо будет добавить
func (s *NotificationService) HandleBudgetExceeded(event events.BudgetExceededEvent) error {
	settings, err := s.GetSettings(event.UserID)
	if err != nil {
		return err
	}
	if !settings.BudgetAlertsEnabled {
		log.Printf("[NotificationService] Budget alerts disabled for user %s", event.UserID)
		return nil
	}

	message := fmt.Sprintf(
		"Бюджет '%s' превышен на %.2f. Потрачено: %.2f из %.2f",
		event.BudgetName,
		event.ExcessAmount,
		event.SpentAmount,
		event.BudgetAmount,
	)
	notification := &models.Notification{
		UserID:  event.UserID,
		Type:    "budget_exceeded",
		Title:   "Бюджет превышен",
		Message: message,
		Data: map[string]interface{}{
			"budget_id":     event.BudgetID,
			"budget_name":   event.BudgetName,
			"budget_amount": event.BudgetAmount,
			"spent_amount":  event.SpentAmount,
			"excess_amount": event.ExcessAmount,
			"category_id":   event.CategoryID,
		},
		IsRead:    false,
		Priority:  "high",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.notificationRepo.SaveNotification(notification); err != nil {
		log.Printf("[NotificationService] Error saving notification: %v", err)
		return fmt.Errorf("save notification: %w", err)
	}
	return nil

}
func (s *NotificationService) HandleBudgetWarning(event events.BudgetWarningEvent) error {

	log.Printf("[NotificationService] Handling BudgetWarning: UserID=%s, BudgetID=%d, Used=%.0f%%",
		event.UserID, event.BudgetID, event.WarningPercent)

	settings, err := s.GetSettings(event.UserID)
	if err != nil {
		return err

	}
	if !settings.BudgetAlertsEnabled {
		log.Printf("[NotificationService] Budget alerts disabled")
		return nil
	}
	if event.WarningPercent < float64(settings.BudgetWarningPercent) {
		log.Printf("[NotificationService] Warning percent %.0f%% below user threshold %d%%",
			event.WarningPercent, settings.BudgetWarningPercent)
		return nil
	}
	message := fmt.Sprintf(
		"Бюджет '%s' использован на %.0f%%. Потрачено: %.2f из %.2f",
		event.BudgetName,
		event.WarningPercent,
		event.SpentAmount,
		event.BudgetAmount,
	)
	notification := &models.Notification{
		UserID:  event.UserID,
		Type:    "budget_warning",
		Title:   "Предупреждение о бюджете",
		Message: message,
		Data: map[string]interface{}{
			"budget_id":       event.BudgetID,
			"budget_name":     event.BudgetName,
			"budget_amount":   event.BudgetAmount,
			"spent_amount":    event.SpentAmount,
			"warning_percent": event.WarningPercent,
			"category_id":     event.CategoryID,
		},
		IsRead:    false,
		Priority:  "medium",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.notificationRepo.SaveNotification(notification); err != nil {
		return fmt.Errorf("save notification: %w", err)
	}
	return nil
}
