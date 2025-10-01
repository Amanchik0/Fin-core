package services

import (
	"fmt"
	"justTest/internal/models"
	"justTest/internal/repo"
	"time"
)

type CategoryService struct {
	accountRepo  *repo.AccountRepository
	categoryRepo *repo.CategoryRepository
	authService  models.AuthService
}

func NewCategoryService(
	accountRepo *repo.AccountRepository,
	categoryRepo *repo.CategoryRepository,
	authService models.AuthService,
) *CategoryService {
	return &CategoryService{
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
		authService:  authService,
	}

}

func (s *CategoryService) GetAccountByUserID(userID string) (*models.Account, error) {
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	return account, nil
}

func (s *CategoryService) CreateCategory(userID string, category *models.Category) (*models.Category, error) {
	if category == nil {
		return nil, fmt.Errorf("category is nil")
	}
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	category.AccountID = account.ID
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	category.IsActive = true
	newCategory, err := s.categoryRepo.CreateCategory(category)
	if err != nil {

		return nil, fmt.Errorf("create category: %w", err)
	}
	return newCategory, nil
}

func (s *CategoryService) UpdateCategory(userID string, category *models.Category) (*models.Category, error) {
	if category == nil {
		return nil, fmt.Errorf("category is nil")
	}
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	if account.ID != category.AccountID {
		return nil, fmt.Errorf("account is not owned by another user")
	}
	category.UpdatedAt = time.Now()
	updatedCategory, err := s.categoryRepo.UpdateCategory(category)
	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}
	return updatedCategory, nil
}

func (s *CategoryService) DeleteCategory(userID string, categoryID int64) error {
	if categoryID == 0 {
		return fmt.Errorf("category is nil")
	}
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}
	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return fmt.Errorf("get category: %w", err)
	}
	if account.ID != category.AccountID {
		return fmt.Errorf("account is not owned by another user")
	}
	return s.categoryRepo.DeleteCategory(categoryID)
}

func (s *CategoryService) GetByAccountID(accountID int64) ([]*models.Category, error) {
	if accountID == 0 {
		return nil, fmt.Errorf("account is nil")
	}
	category, err := s.categoryRepo.GetByAccountID(accountID)
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category is nil")
	}
	return category, nil
}

func (s *CategoryService) GetByCategoryID(userID string, categoryID int64) (*models.Category, error) {
	if categoryID == 0 {
		return nil, fmt.Errorf("category is nil")
	}

	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category is nil")
	}
	return category, nil

}
