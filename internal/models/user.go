package models

import (
	"time"
)

// Account - единственный финансовый аккаунт пользователя
type Account struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`             // ID из auth-сервиса (Node.js)
	Name         string    `json:"name" db:"name"`                   // "Мой аккаунт", или имя пользователя
	DisplayName  string    `json:"display_name" db:"display_name"`   // имя для отображения в UI
	Timezone     string    `json:"timezone" db:"timezone"`           // для корректного отображения времени
	BaseCurrency string    `json:"base_currency" db:"base_currency"` // основная валюта для расчетов (KZT, USD, EUR)
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// BankAccount - банковский счет внутри аккаунта
type BankAccount struct {
	ID          int64     `json:"id" db:"id"`
	AccountID   int64     `json:"account_id" db:"account_id"`     // ссылка на главный аккаунт
	Name        string    `json:"name" db:"name"`                 // "Каспи рубли", "Халык доллары", "Наличные USD"
	Currency    string    `json:"currency" db:"currency"`         // "KZT", "USD", "EUR"
	AccountType string    `json:"account_type" db:"account_type"` // "cash", "debit", "credit", "savings"
	BankName    string    `json:"bank_name" db:"bank_name"`       // "Kaspi", "Halyk", "Cash"
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Transaction - операции по банковским счетам
type Transaction struct {
	ID              int64     `json:"id" db:"id"`
	BankAccountID   int64     `json:"bank_account_id" db:"bank_account_id"`
	CategoryID      *int64    `json:"category_id" db:"category_id"` // может быть null для переводов
	Amount          float64   `json:"amount" db:"amount"`           // положительное для доходов, отрицательное для расходов
	Description     string    `json:"description" db:"description"`
	TransactionType string    `json:"transaction_type" db:"transaction_type"` // "income", "expense", "transfer"
	Date            time.Time `json:"date" db:"date"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	// Для переводов между банковскими счетами
	ToAccountID  *int64   `json:"to_account_id" db:"to_account_id"` // ID другого банковского счета
	TransferRate *float64 `json:"transfer_rate" db:"transfer_rate"` // курс валют если перевод между валютами
}

// Category - категории транзакций
type Category struct {
	ID        int64     `json:"id" db:"id"`
	AccountID int64     `json:"account_id" db:"account_id"` // привязаны к аккаунту
	Name      string    `json:"name" db:"name"`             // "Продукты", "Транспорт", "Зарплата"
	Type      string    `json:"type" db:"type"`             // "income", "expense"
	Color     string    `json:"color" db:"color"`           // для UI
	Icon      string    `json:"icon" db:"icon"`             // для UI
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Budget - бюджеты на категории
type Budget struct {
	ID         int64     `json:"id" db:"id"`
	AccountID  int64     `json:"account_id" db:"account_id"`
	CategoryID int64     `json:"category_id" db:"category_id"`
	Amount     float64   `json:"amount" db:"amount"` // лимит на период в базовой валюте аккаунта
	Period     string    `json:"period" db:"period"` // "monthly", "weekly", "yearly"
	StartDate  time.Time `json:"start_date" db:"start_date"`
	EndDate    time.Time `json:"end_date" db:"end_date"`
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// BankAccountBalance - кэшированные балансы банковских счетов
type BankAccountBalance struct {
	BankAccountID int64     `json:"bank_account_id" db:"bank_account_id"`
	Balance       float64   `json:"balance" db:"balance"`
	Currency      string    `json:"currency" db:"currency"`
	LastUpdated   time.Time `json:"last_updated" db:"last_updated"`
}

// CurrencyRate - курсы валют для расчета общего баланса
type CurrencyRate struct {
	ID           int64     `json:"id" db:"id"`
	FromCurrency string    `json:"from_currency" db:"from_currency"`
	ToCurrency   string    `json:"to_currency" db:"to_currency"`
	Rate         float64   `json:"rate" db:"rate"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// DTO для взаимодействия с auth-сервисом
type AuthUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Value objects для бизнес-логики
type CurrencyBalance struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type BudgetStatus struct {
	Budget     *Budget `json:"budget"`
	Spent      float64 `json:"spent"`       // потрачено в текущем периоде в валюте бюджета
	Remaining  float64 `json:"remaining"`   // осталось
	Progress   float64 `json:"progress"`    // процент использования (0-100)
	IsExceeded bool    `json:"is_exceeded"` // превышен ли бюджет
}

type BankAccountSummary struct {
	BankAccount *BankAccount `json:"bank_account"`
	Balance     float64      `json:"balance"`
}

type AccountSummary struct {
	Account         *Account              `json:"account"`
	Balances        []*CurrencyBalance    `json:"balances"` // балансы по валютам: [{currency: "KZT", amount: 150000}, ...]
	BankAccounts    []*BankAccountSummary `json:"bank_accounts"`
	MonthlyIncome   []*CurrencyBalance    `json:"monthly_income"`   // доходы за месяц по валютам
	MonthlyExpenses []*CurrencyBalance    `json:"monthly_expenses"` // расходы за месяц по валютам
	Budgets         []*BudgetStatus       `json:"budgets"`
}

type CreateAccountRequest struct {
	DisplayName string `json:"display_name" binding:"required,min=2,max=40"`
}
type UpdateAccountRequest struct {
	DisplayName string `json:"display_name" binding:"required,min=2,max=40"`
	Timezone    string `json:"timezone" binding:"required"`
}
type AccountResponse struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Timezone    string    `json:"timezone"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateBankAccountRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=40"`
	Currency    string `json:"currency" binding:"required,currency"`
	AccountType string `json:"account_type" binding:"required,currency"`
	BankName    string `json:"bank_name" binding:"required,min=2,max=40"`
}
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50"`
	Type  string `json:"type" binding:"required,oneof=income expense"`
	Color string `json:"color" binding:"required"`
	Icon  string `json:"icon" binding:"required"`
}
