package services

import (
	"fmt"
	"justTest/internal/models"
	"justTest/internal/repo"
	"time"
)

type TransactionService struct {
	transactionRepo repo.TransactionRepository
	bankAccountRepo repo.BankAccountRepository
	categoryRepo    repo.CategoryRepository
	accountRepo     repo.AccountRepository
}

func NewTransactionService(
	transactionRepo repo.TransactionRepository,
	bankAccountRepo repo.BankAccountRepository,
	categoryRepo repo.CategoryRepository,
	accountRepo repo.AccountRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		bankAccountRepo: bankAccountRepo,
		categoryRepo:    categoryRepo,
		accountRepo:     accountRepo,
	}
}

func (s *TransactionService) validateCategoryOwnership(userID int64, categoryID int64) error {
	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return fmt.Errorf("category not found", err)

	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("user not found", err)
	}
	if category.AccountID != userAccount.ID {
		return fmt.Errorf("user is not owned by the category", err)
	}
	return nil
}
func (s *TransactionService) validateBankAccountOwnership(userID int64, bankAccountID int64) error {
	bankAccount, err := s.bankAccountRepo.GetByBankAccountID(bankAccountID)
	if err != nil {
		return fmt.Errorf("bank account not found", err)
	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("user not found", err)
	}
	if bankAccount.AccountID != userAccount.ID {
		return fmt.Errorf("user is not owned by the bank account", err)
	}
	return nil
}

func (s *TransactionService) CreateTransaction(userID int64, bankAccountID int64, amount float64, description string, categoryID *int64, transactionType string) (*models.Transaction, error) {

	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	if amount == 0 {
		return nil, fmt.Errorf("invalid amount")
	}
	if transactionType == "" {
		return nil, fmt.Errorf("invalid transaction type")
	}
	if description == "" {
		return nil, fmt.Errorf("invalid description")
	}
	if bankAccountID <= 0 {
		return nil, fmt.Errorf("invalid bank account id")
	}
	err := s.validateBankAccountOwnership(userID, bankAccountID)
	if err != nil {
		return nil, err
	}
	if categoryID != nil {
		err := s.validateCategoryOwnership(userID, *categoryID)
		if err != nil {
			return nil, err
		}
	}
	if transactionType == "expense" && amount > 0 {
		amount = -amount
	}
	if transactionType == "income" && amount < 0 {
		amount = -amount
	}
	transaction := &models.Transaction{
		BankAccountID:   bankAccountID,
		CategoryID:      categoryID,
		Amount:          amount,
		Description:     description,
		TransactionType: transactionType,
		Date:            time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	createdTransaction, err := s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, err
	}
	return createdTransaction, nil

}

// TransferBetweenAccounts
// GetTransactionHistory
// GetBankAccountBalance
// GetTransaction
func (s *TransactionService) TransferBetweenAccounts(userID int64, fromAccountID int64, toAccountID int64, description string, amount float64) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user id")
	}
	if fromAccountID <= 0 {
		return fmt.Errorf("invalid from account id")
	}
	if toAccountID <= 0 {
		return fmt.Errorf("invalid to account id")
	}
	if amount == 0 {
		return fmt.Errorf("invalid amount")
	}
	if description == "" {
		return fmt.Errorf("invalid description")
	}
	err := s.validateBankAccountOwnership(userID, fromAccountID)
	if err != nil {
		return fmt.Errorf("source account: %w", err)
	}

	err = s.validateBankAccountOwnership(userID, toAccountID)
	if err != nil {
		return fmt.Errorf("destination account: %w", err)
	}

	err = s.transactionRepo.CreateTransaction(nil, &fromAccountID, nil, toAccountID, amount, description, nil)
	if err != nil {
		return err
	}
	return nil

}
func (s *TransactionService) GetTransactionHistory(userID int64, bankAccountID int64) ([]*models.Transaction, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	if bankAccountID <= 0 {
		return nil, fmt.Errorf("invalid from bank account id")
	}
	err := s.validateBankAccountOwnership(userID, bankAccountID)
	if err != nil {
		return nil, err
	}
	transactions, err := s.transactionRepo.GetByBankAccountID(bankAccountID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
func (s *TransactionService) GetAllTransactions(userID int64) ([]*models.Transaction, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	userAccount, err := s.accountRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	transactions, err := s.transactionRepo.GetByAccountID(userAccount.ID)
	if err != nil {
		return nil, err

	}
	return transactions, nil
}
func (s *TransactionService) GetBankAccountBalance(userID int64, bankAccountID int64) (float64, error) {
	if userID <= 0 {
		return 0, fmt.Errorf("invalid user id")
	}
	if bankAccountID <= 0 {
		return 0, fmt.Errorf("invalid bank account id")
	}
	err := s.validateBankAccountOwnership(userID, bankAccountID)
	if err != nil {
		return 0, err
	}
	balance, err := s.transactionRepo.GetTotalAmountByBankAccountID(bankAccountID)
	if err != nil {
		return 0, err

	}
	return balance, nil
}
func (s *TransactionService) GetTransactionByID(userID int64, transactionID int64) (*models.Transaction, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")

	}
	if transactionID <= 0 {
		return nil, fmt.Errorf("invalid transaction id")
	}
	transaction, err := s.transactionRepo.GetByTransactionID(transactionID)
	if err != nil {
		return nil, err
	}
	err = s.validateBankAccountOwnership(userID, transaction.BankAccountID)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
