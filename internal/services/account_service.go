package services

import (
	"fmt"
	"justTest/internal/interfaces"
	"justTest/internal/models"
	"time"
)

type AccountService struct {
	accountRepo     interfaces.AccountRepository
	BankAccountRepo interfaces.BankAccountRepository
	transactionRepo interfaces.TransactionRepository
	authService     models.AuthService
}

func NewAccountService(
	accountRepo interfaces.AccountRepository,
	bankAccountRepo interfaces.BankAccountRepository,
	transactionRepo interfaces.TransactionRepository,
	authService models.AuthService,
) *AccountService {
	return &AccountService{
		accountRepo:     accountRepo,
		BankAccountRepo: bankAccountRepo,
		transactionRepo: transactionRepo,
		authService:     authService,
	}
}
func (s *AccountService) GetUserAccount(userID string) (*models.Account, error) {
	if userID == "" {
		return nil, fmt.Errorf("invalid user id: %s", userID)
	}
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: empty")
	}

	return account, nil

}
func (s *AccountService) CreateAccount(userID string, displayName string) (*models.Account, error) {
	if userID == "" {
		return nil, fmt.Errorf("invalid user id: empty")
	}
	if displayName == "" {
		return nil, fmt.Errorf("display name is required")

	}
	existingAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get account by user id failed, err:%v", err)
	}
	if existingAccount != nil {
		return nil, fmt.Errorf("account with id %s already exists", existingAccount.UserID)
	}
	//authUser, err := s.authService.GetUserByID(fmt.Sprintf("%d", userID))
	//if err != nil {
	//	return nil, fmt.Errorf("get user by id failed, err:%v", err)
	//}
	newAccount := &models.Account{
		UserID:      userID,
		DisplayName: displayName,
		Name:        displayName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		// потом добавлю еще чтот по сути
	}
	createdAccount, err := s.accountRepo.Create(newAccount)
	if err != nil {
		return nil, fmt.Errorf("create account failed, err:%v", err)
	}
	return createdAccount, nil

}

// summary надо написать потом как нибудь
