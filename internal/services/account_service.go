package services

import (
	"fmt"
	"justTest/internal/models"
	"justTest/internal/repo"
	"time"
)

type AccountService struct {
	accountRepo     *repo.AccountRepository
	BankAccountRepo *repo.BankAccountRepository
	transactionRepo *repo.TransactionRepository
	authService     models.AuthService
}

func NewAccountService(
	accountRepo *repo.AccountRepository,
	bankAccountRepo *repo.BankAccountRepository,
	transactionRepo *repo.TransactionRepository,
	authService models.AuthService,
) *AccountService {
	return &AccountService{
		accountRepo:     accountRepo,
		BankAccountRepo: bankAccountRepo,
		transactionRepo: transactionRepo,
		authService:     authService,
	}
}
func (s *AccountService) GetUserAccount(userID int64) (*models.Account, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", userID)
	}
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get account by user id failed, err:%v", err)
	}
	return account, nil

}
func (s *AccountService) CreateAccount(userID int64, displayName string) (*models.Account, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id", userID)
	}
	if displayName == "" {
		return nil, fmt.Errorf("display name is required")

	}
	existingAccount, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get account by user id failed, err:%v", err)
	}
	if existingAccount != nil {
		return nil, fmt.Errorf("account with id %d already exists", existingAccount.UserID)
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
