package services

import (
	"fmt"
	"justTest/internal/interfaces"
	"justTest/internal/models"
	"justTest/internal/utils"
	"time"
)

type TransactionService struct {
	transactionRepo interfaces.TransactionRepository
	bankAccountRepo interfaces.BankAccountRepository
	categoryRepo    interfaces.CategoryRepository
	accountRepo     interfaces.AccountRepository
}

func NewTransactionService(
	transactionRepo interfaces.TransactionRepository,
	bankAccountRepo interfaces.BankAccountRepository,
	categoryRepo interfaces.CategoryRepository,
	accountRepo interfaces.AccountRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		bankAccountRepo: bankAccountRepo,
		categoryRepo:    categoryRepo,
		accountRepo:     accountRepo,
	}
}

func (s *TransactionService) validateCategoryOwnership(userID string, categoryID int64) error {
	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if category.AccountID != userAccount.ID {
		return fmt.Errorf("user is not owned by the category: %w", err)
	}
	return nil
}
func (s *TransactionService) validateBankAccountOwnership(userID string, bankAccountID int64) error {
	bankAccount, err := s.bankAccountRepo.GetByBankAccountID(bankAccountID)
	if err != nil {
		return fmt.Errorf("bank account not found: %w", err)
	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if bankAccount.AccountID != userAccount.ID {
		return fmt.Errorf("user is not owned by the bank account: %w", err)
	}
	return nil
}

func (s *TransactionService) CreateTransaction(userID string, bankAccountID int64, amount float64, description string, categoryID *int64, transactionType string) (*models.Transaction, error) {

	if userID == "" {
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
func (s *TransactionService) TransferBetweenAccounts(userID string, fromAccountID int64, toAccountID int64, description string, amount float64) error {
	if userID == "" {
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
func (s *TransactionService) GetTransactionHistory(userID string, bankAccountID int64) ([]*models.Transaction, error) {
	if userID == "" {
		return nil, fmt.Errorf("invalid user id")
	}
	if bankAccountID <= 0 {
		return nil, fmt.Errorf("invalid from bank account id")
	}
	err := s.validateBankAccountOwnership(userID, bankAccountID)
	if err != nil {
		return nil, err
	}
	limit, offset := utils.GetPaginationParams(1, 20) // Для первой страницы

	transactions, err := s.transactionRepo.GetByBankAccountID(bankAccountID, limit, offset)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
func (s *TransactionService) GetAllTransactions(userID string) ([]*models.Transaction, error) {
	if userID == "" {
		return nil, fmt.Errorf("invalid user id")
	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	limit, offset := utils.GetPaginationParams(1, 20) // Для первой страницы

	transactions, err := s.transactionRepo.GetByAccountID(userAccount.ID, limit, offset)
	if err != nil {
		return nil, err

	}
	return transactions, nil
}
func (s *TransactionService) GetBankAccountBalance(userID string, bankAccountID int64) (float64, error) {
	if userID == "" {
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
func (s *TransactionService) GetTransactionByID(userID string, transactionID int64) (*models.Transaction, error) {
	if userID == "" {
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

func (s *TransactionService) GetAllTransactionsByCategoryID(userID string, categoryID int64) ([]*models.Transaction, error) {
	if categoryID <= 0 {
		return nil, fmt.Errorf("invalid category id")
	}
	if userID == "" {
		return nil, fmt.Errorf("invalid user id")

	}
	limit, offset := utils.GetPaginationParams(1, 20)
	transactions, err := s.transactionRepo.GetByCategoryID(categoryID, limit, offset)
	if err != nil {
		return nil, err

	}
	return transactions, nil
}
