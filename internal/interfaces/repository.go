package interfaces

import (
	"justTest/internal/models"
	"time"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) (*models.Transaction, error)
	GetByBankAccountID(BankAccountID int64, limit, offset int) ([]*models.Transaction, error)
	GetByCategoryID(CategoryID int64, limit, offset int) ([]*models.Transaction, error)
	GetByTransactionID(TransactionID int64) (*models.Transaction, error)
	GetByAccountID(AccountID int64, limit, offset int) ([]*models.Transaction, error)
	GetTransfersByAccountID(AccountID int64) ([]*models.Transaction, error)
	CreateTransaction(AccountID, FromBankAccountID, categoryID *int64, toAccountID int64, amount float64, description string, transferRate *float64) error
	GetTotalAmountByBankAccountID(BankAccountID int64) (float64, error)
	GetSpentAmountByCategoryAndMonth(categoryID int64, year, month int) (float64, error)
	GetTransactionsByCategoryAndMonth(categoryID int64, year, month int, limit, offset int) ([]*models.Transaction, error)
	GetTransactionsByDateRangeWithCategory(categoryID int64, startDate, endDate time.Time, limit, offset int) ([]*models.Transaction, error)
	GetTransactionsByDateRange(startDate, endDate time.Time, limit, offset int) ([]*models.Transaction, error)
}

type AccountRepository interface {
	Create(account *models.Account) (*models.Account, error)
	GetByUserID(userID string) (*models.Account, error)
	GetByID(id int64) (*models.Account, error)
}
type BankAccountRepository interface {
	Create(bankAccount *models.BankAccount) (*models.BankAccount, error)
	GetByAccountID(accountID int64) ([]*models.BankAccount, error)
	GetActiveBankAccounts(accountID int64) ([]*models.BankAccount, error)
	GetBankAccountByCurrency(accountID int64, currency string) ([]*models.BankAccount, error)
	GetByBankAccountID(BankAccountID int64) (*models.BankAccount, error)
	ExsitsAccountIDAndName(accountID int64, name string) (bool, error)
	DeActiveBankAccount(bankAccountID int64) error
	ActivateBankAccount(bankAccountID int64) error
	DeleteBankAccount(bankAccountID int64) error
}

type CategoryRepository interface {
	CreateCategory(category *models.Category) (*models.Category, error)
	UpdateCategory(category *models.Category) (*models.Category, error)
	DeleteCategory(categoryID int64) error
	GetByAccountID(accountID int64) ([]*models.Category, error)
	GetByID(categoryID int64) (*models.Category, error)
}
type BudgetRepository interface {
	CreateBudget(budget *models.Budget) (*models.Budget, error)
	GetBudget(budgetID int64) (*models.Budget, error)
	GetBudgetByCategoryID(categoryID int64) (*models.Budget, error)
	GetBudgetByCategoryAndMonth(categoryID int64, year, month int) (*models.Budget, error)
	GetBudgetsByAccountAndMonth(accountID int64, year, month int) ([]*models.Budget, error)
}

type NotificationRepository interface {
	SaveNotification(notification *models.Notification) error
	GetNotificationByID(id int64) (*models.Notification, error)
	GetUserNotifications(userID string, limit, offset int) ([]*models.Notification, error)
	GetUnreadNotifications(userID string) ([]*models.Notification, error)

	MarkAsRead(id int64) error
	MarkAllAsRead(userID string) error

	DeleteNotification(id int64) error
	DeleteOldNotifications(days int) error
}
