package repo

import (
	"database/sql"
	"justTest/internal/models"
)

type UserNotificationSettingsRepository struct {
	db *sql.DB
}

func NewUserNotificationSettingsRepository(db *sql.DB) *UserNotificationSettingsRepository {
	return &UserNotificationSettingsRepository{db: db}

}
func (r *UserNotificationSettingsRepository) GetSettings(userID string) (*models.UserNotificationSettings, error) {
	query := ` 
	SELECT id, user_id, budget_alerts_enabled, balance_alerts_enabled, budget_warning_percent,
       low_balance_threshold, preferred_channel, created_at, updated_at
FROM user_notification_settings
WHERE user_id = $1

`
	row := r.db.QueryRow(query, userID)
	var settings models.UserNotificationSettings
	err := row.Scan(
		&settings.ID,
		&settings.UserID,
		&settings.BudgetAlertsEnabled,
		&settings.BalanceAlertsEnabled,
		&settings.BudgetWarningPercent,
		&settings.LowBalanceThreshold,
		&settings.PreferredChannel,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}
func (r *UserNotificationSettingsRepository) SaveSettings(settings *models.UserNotificationSettings) error {
	query := ` 
insert into user_notification_settings (
                                        user_id, budget_alerts_enabled, balance_alerts_enabled,
                                        budget_warning_percent,low_balance_threshold, preferred_channel, created_at, updated_at
                                        
) VALUES ($1, $2, $3, $4, $5, $6, $7 , $8) 
returning id; `
	return r.db.QueryRow(
		query,
		settings.UserID,
		settings.BudgetAlertsEnabled,
		settings.BalanceAlertsEnabled,
		settings.BudgetWarningPercent,
		settings.LowBalanceThreshold,
		settings.PreferredChannel,
		settings.CreatedAt,
		settings.UpdatedAt,
	).Scan(&settings.ID)

}
func (r *UserNotificationSettingsRepository) UpdateSettings(settings *models.UserNotificationSettings) (*models.UserNotificationSettings, error) {
	query := ` 
	update user_notification_settings
SET budget_alerts_enabled = $2,
    balance_alerts_enabled = $3,
    budget_warning_percent = $4,
    low_balance_threshold = $5,
    preferred_channel = $6,
    updated_at = NOW()
WHERE user_id = $1
	RETURNING id, user_id, budget_alerts_enabled, balance_alerts_enabled, budget_warning_percent,
	          low_balance_threshold, preferred_channel, created_at, updated_at;
	
`
	var updated models.UserNotificationSettings
	err := r.db.QueryRow(
		query,
		settings.UserID,
		settings.BudgetAlertsEnabled,
		settings.BalanceAlertsEnabled,
		settings.BudgetWarningPercent,
		settings.LowBalanceThreshold,
		settings.PreferredChannel,
	).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.BudgetAlertsEnabled,
		&updated.BalanceAlertsEnabled,
		&updated.BudgetWarningPercent,
		&updated.LowBalanceThreshold,
		&updated.PreferredChannel,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}
