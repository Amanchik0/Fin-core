package services

import (
	"fmt"
	"justTest/internal/models"
	"justTest/internal/repo"
	"time"
)

type BankAccService struct {
	BankAccountRepository repo.BankAccountRepository
	accountRepo           repo.AccountRepository
}

func NewBankAccService(
	BankAccountRepository repo.BankAccountRepository,
	accountRepo repo.AccountRepository,
) *BankAccService {
	return &BankAccService{
		BankAccountRepository: BankAccountRepository,
		accountRepo:           accountRepo,
	}
}
func (s *BankAccService) CreateBankAccount(userID int64, name, currency, accountType, bankName string) (*models.BankAccount, error) {
	if userID <= 0 {
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
func (s *BankAccService) GetBankAccount(userID int64, bankAccountID int64) (*models.BankAccount, error) {
	if userID <= 0 {
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
func (s *BankAccService) GetBankAccountsByAccountID(userID int64) ([]*models.BankAccount, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	userAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user account")
	}
	return s.BankAccountRepository.GetByAccountID(userAccount.ID)

}
func (s *BankAccService) DeActiveBankAccount(userID int64, bankAccountID int64) error {
	if userID <= 0 {
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
func (s *BankAccService) DeleteBankAccount(userID int64, bankAccountID int64) error {
	if userID <= 0 {
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
