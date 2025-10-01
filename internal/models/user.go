package models

import (
	"time"
)

// Account - единственный финансовый аккаунт пользователя
type Account struct {
	ID           int64     `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`             // ID из auth-сервиса (Node.js)
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
	ID       string `json:"id"` // ← изменить на string
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
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Timezone    string    `json:"timezone"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MonthlyReport - месячный отчет
type MonthlyReport struct {
	Month        int                 `json:"month"`
	Year         int                 `json:"year"`
	TotalIncome  float64             `json:"total_income"`
	TotalExpense float64             `json:"total_expense"`
	NetIncome    float64             `json:"net_income"`
	Categories   []*CategorySpending `json:"categories"`
	TopExpenses  []*Transaction      `json:"top_expenses"`
}

// CategorySpending - траты по категории
type CategorySpending struct {
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

// IncomeExpenseReport - отчет доходы vs расходы
type IncomeExpenseReport struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetIncome    float64 `json:"net_income"`
	SavingsRate  float64 `json:"savings_rate"` // процент сбережений
}

// BudgetAlert - уведомление о превышении бюджета
type BudgetAlert struct {
	BudgetID     int64   `json:"budget_id"`
	BudgetName   string  `json:"budget_name"`
	BudgetAmount float64 `json:"budget_amount"`
	SpentAmount  float64 `json:"spent_amount"`
	ExcessAmount float64 `json:"excess_amount"`
}

// BalanceAlert - уведомление о низком балансе
type BalanceAlert struct {
	BankAccountID  int64   `json:"bank_account_id"`
	AccountName    string  `json:"account_name"`
	CurrentBalance float64 `json:"current_balance"`
	AlertThreshold float64 `json:"alert_threshold"`
}

type CreateBankAccountRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=40"`
	Currency    string `json:"currency" binding:"required,oneof=KZT USD EUR RUB"`
	AccountType string `json:"account_type" binding:"required,oneof=cash debit credit savings"`
	BankName    string `json:"bank_name" binding:"required,min=2,max=40"`
}
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50"`
	Type  string `json:"type" binding:"required,oneof=income expense"`
	Color string `json:"color" binding:"required"`
	Icon  string `json:"icon" binding:"required"`
}
type CreateTransactionRequest struct {
	BankAccountID   int64   `json:"bank_account_id" binding:"required"`                       // ID банковского счета
	Amount          float64 `json:"amount" binding:"required"`                                // Сумма транзакции
	Description     string  `json:"description" binding:"required,min=1,max=255"`             // Описание транзакции
	CategoryID      *int64  `json:"category_id"`                                              // ID категории (может быть null)
	TransactionType string  `json:"transaction_type" binding:"required,oneof=income expense"` // Тип: доход или расход
}

type TransferRequest struct {
	FromAccountID int64   `json:"from_account_id" binding:"required"`           // Откуда переводим
	ToAccountID   int64   `json:"to_account_id" binding:"required"`             // Куда переводим
	Amount        float64 `json:"amount" binding:"required"`                    // Сумма перевода
	Description   string  `json:"description" binding:"required,min=1,max=255"` // Описание перевода
}

type TransactionResponse struct {
	ID              int64   `json:"id"`
	BankAccountID   int64   `json:"bank_account_id"`
	CategoryID      *int64  `json:"category_id"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	TransactionType string  `json:"transaction_type"`
	Date            string  `json:"date"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
