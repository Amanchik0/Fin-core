package events

import "time"

type TransactionCreatedEvent struct {
	TransactionID int64     `json:"transaction_id"`
	UserID        string    `json:"user_id"`
	CategoryID    int64     `json:"category_id"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	Timestamp     time.Time `json:"timestamp"`
}

type BudgetExceededEvent struct {
	UserID       string    `json:"user_id"`
	BudgetID     int64     `json:"budget_id"`
	BudgetName   string    `json:"budget_name"`
	BudgetAmount float64   `json:"budget_amount"`
	SpentAmount  float64   `json:"spent_amount"`
	ExcessAmount float64   `json:"excess_amount"`
	CategoryID   int64     `json:"category_id"`
	Timestamp    time.Time `json:"timestamp"`
}

type LowBalanceEvent struct {
	UserID         string    `json:"user_id"`
	BankAccountID  int64     `json:"bank_account_id"`
	AccountName    string    `json:"account_name"`
	CurrentBalance float64   `json:"current_balance"`
	AlertThreshold float64   `json:"alert_threshold"`
	Timestamp      time.Time `json:"timestamp"`
}

type BudgetWarningEvent struct {
	UserID         string    `json:"user_id"`
	BudgetID       int64     `json:"budget_id"`
	BudgetName     string    `json:"budget_name"`
	BudgetAmount   float64   `json:"budget_amount"`
	SpentAmount    float64   `json:"spent_amount"`
	WarningPercent float64   `json:"warning_percent"`
	ExcessAmount   float64   `json:"excess_amount"`
	CategoryID     int64     `json:"category_id"`
	Timestamp      time.Time `json:"timestamp"`
}

type NotificationEvent struct {
	UserID    string                 `json:"user_id"`
	Type      string                 `json:"type"` // "budget_exceeded", "low_balance", "budget_warning"
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Priority  string                 `json:"priority"` // "low", "medium", "high"
}
