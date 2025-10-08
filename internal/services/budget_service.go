package services

import (
	"justTest/internal/models"
	"justTest/internal/repo"
)

type BudgetService struct {
	// budgetRepo  repo.BudgetRepository
	accountRepo repo.AccountRepository
}

// func NewBudgetService(budgetRepo repo.BudgetRepository, accountRepo repo.AccountRepository) *BudgetService {
// 	return &BudgetService{
// 		budgetRepo:  budgetRepo,
// 		accountRepo: accountRepo,
// 	}
// }

// CreateBudget - создать бюджет
func (s *BudgetService) CreateBudget(userID string, budget *models.Budget) (*models.Budget, error) {
	// TODO: Реализовать создание бюджета
	return nil, nil
}

// GetBudgets - получить бюджеты пользователя
func (s *BudgetService) GetBudgets(userID string) ([]*models.Budget, error) {
	// TODO: Реализовать получение бюджетов
	return nil, nil
}

// GetBudgetStatus - статус бюджета (потрачено/осталось)
func (s *BudgetService) GetBudgetStatus(userID string, budgetID int64) (*models.BudgetStatus, error) {
	// TODO: Реализовать получение статуса бюджета
	return nil, nil
}
