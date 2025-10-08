package interfaces

import (
	"justTest/internal/models"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) (*models.Transaction, error)
	GetByBankAccountID(BankAccountID int64, limit, offset int) ([]*models.Transaction, error)
	GetByCategoryID(CategoryID int64) ([]*models.Transaction, error)
	GetByTransactionID(TransactionID int64) (*models.Transaction, error)
	GetByAccountID(AccountID int64, limit, offset int) ([]*models.Transaction, error)
	GetTransfersByAccountID(AccountID int64) ([]*models.Transaction, error)
	CreateTransaction(AccountID, FromBankAccountID, categoryID *int64, toAccountID int64, amount float64, description string, transferRate *float64) error
	GetTotalAmountByBankAccountID(BankAccountID int64) (float64, error)
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
