package services

import (
	"fmt"
	"justTest/internal/interfaces"
	"justTest/internal/models"
	"time"
)

type BankAccService struct {
	BankAccountRepository interfaces.BankAccountRepository
	accountRepo           interfaces.AccountRepository
}

func NewBankAccService(
	BankAccountRepository interfaces.BankAccountRepository,
	accountRepo interfaces.AccountRepository,
) *BankAccService {
	return &BankAccService{
		BankAccountRepository: BankAccountRepository,
		accountRepo:           accountRepo,
	}
}
func (s *BankAccService) CreateBankAccount(userID string, name, currency, accountType, bankName string) (*models.BankAccount, error) {
	if userID == "" {
		return nil, fmt.Errorf("invalid user id")
	}
	if name == "" {
		return nil, fmt.Errorf("empty name")
	}
	if currency == "" {
		return nil, fmt.Errorf("empty currency")
	}
	if accountType == "" {
		return nil, fmt.Errorf("empty accountType")
	}
	if bankName == "" {
		return nil, fmt.Errorf("empty bankName")
	}
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	exists, err := s.BankAccountRepository.ExsitsAccountIDAndName(account.ID, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("duplicate bank account")
	}
	newBankAccount := &models.BankAccount{
		AccountID:   account.ID,
		Name:        name,
		Currency:    currency,
		AccountType: accountType,
		BankName:    bankName,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	createdBankAccount, err := s.BankAccountRepository.Create(newBankAccount)
	if err != nil {
		return nil, err
	}
	return createdBankAccount, nil
}
func (s *BankAccService) GetBankAccount(userID string, bankAccountID int64) (*models.BankAccount, error) {
	if userID == "" {
		return nil, fmt.Errorf("invalid user id")
	}
	if bankAccountID <= 0 {
		return nil, fmt.Errorf("invalid bank account id")
	}
	bankAccount, err := s.BankAccountRepository.GetByBankAccountID(bankAccountID)
	if err != nil {
		return nil, err
	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if bankAccount.AccountID != userAccount.ID {
		return nil, fmt.Errorf("invalid user account")
	}

	return bankAccount, nil
}
func (s *BankAccService) GetBankAccountsByAccountID(userID string) ([]*models.BankAccount, error) {
	if userID == "" {
		return nil, fmt.Errorf("invalid user id")
	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user account")
	}
	return s.BankAccountRepository.GetByAccountID(userAccount.ID)

}
func (s *BankAccService) DeActiveBankAccount(userID string, bankAccountID int64) error {
	if userID == "" {
		return fmt.Errorf("invalid user id")
	}
	if bankAccountID <= 0 {
		return fmt.Errorf("invalid bank account id")
	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("invalid user account")
	}
	bankAccount, err := s.BankAccountRepository.GetByBankAccountID(bankAccountID)
	if err != nil {
		return fmt.Errorf("invalid user account")
	}
	if bankAccount.AccountID != userAccount.ID {
		return fmt.Errorf("invalid user account")

	}
	return s.BankAccountRepository.DeActiveBankAccount(bankAccountID)
}

func (s *BankAccService) ActivateBankAccount(userID string, bankAccountID int64) error {
	if userID == "" {
		return fmt.Errorf("invalid user id")
	}
	if bankAccountID <= 0 {
		return fmt.Errorf("invalid bank account id")
	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("invalid user account")
	}
	bankAccount, err := s.BankAccountRepository.GetByBankAccountID(bankAccountID)
	if err != nil {
		return fmt.Errorf("invalid user account")
	}
	if bankAccount.AccountID != userAccount.ID {
		return fmt.Errorf("invalid user account")

	}
	return s.BankAccountRepository.ActivateBankAccount(bankAccountID)
}

func (s *BankAccService) DeleteBankAccount(userID string, bankAccountID int64) error {
	if userID == "" {
		return fmt.Errorf("invalid user id")
	}
	if bankAccountID <= 0 {
		return fmt.Errorf("invalid bank account id")
	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("invalid user account")
	}
	bankAccount, err := s.BankAccountRepository.GetByBankAccountID(bankAccountID)
	if err != nil {
		return fmt.Errorf("invalid user account")
	}
	if bankAccount.AccountID != userAccount.ID {
		return fmt.Errorf("invalid user account")

	}
	return s.BankAccountRepository.DeleteBankAccount(bankAccountID)
}
