package services

import (
	"fmt"
	"justTest/internal/interfaces"
	"justTest/internal/models"
	"time"
)

type BudgetService struct {
	budgetRepo      interfaces.BudgetRepository
	transactionRepo interfaces.TransactionRepository
	accountRepo     interfaces.AccountRepository
	categoryRepo    interfaces.CategoryRepository
}

func NewBudgetService(
	budgetRepo interfaces.BudgetRepository,
	transactionRepo interfaces.TransactionRepository,
	accountRepo interfaces.AccountRepository,
	categoryRepo interfaces.CategoryRepository,
) *BudgetService {
	return &BudgetService{
		budgetRepo:      budgetRepo,
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
	}
}

func (s *BudgetService) CreateBudget(userID string, req *models.CreateBudgetRequest) (*models.Budget, error) {
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}

	category, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}
	if category.AccountID != account.ID {
		return nil, fmt.Errorf("category does not belong to user")
	}

	startDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // последний день месяца

	existingBudget, err := s.budgetRepo.GetBudgetByCategoryAndMonth(req.CategoryID, req.Year, req.Month)
	if err == nil && existingBudget != nil {
		return nil, fmt.Errorf("budget already exists for category %d in %d-%02d", req.CategoryID, req.Year, req.Month)
	}

	budget := &models.Budget{
		AccountID:       account.ID,
		BudgetLimitName: req.BudgetName,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Period:          "monthly",
		StartDate:       startDate,
		EndDate:         endDate,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	createdBudget, err := s.budgetRepo.CreateBudget(budget)
	if err != nil {
		return nil, fmt.Errorf("create budget: %w", err)
	}

	return createdBudget, nil
}

func (s *BudgetService) GetBudgets(userID string, year, month int) ([]*models.BudgetWithStatus, error) {

	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}

	budgets, err := s.budgetRepo.GetBudgetsByAccountAndMonth(account.ID, year, month)
	if err != nil {
		return nil, fmt.Errorf("get budgets: %w", err)
	}

	var budgetsWithStatus []*models.BudgetWithStatus
	for _, budget := range budgets {
		status, err := s.getBudgetStatus(budget, year, month)
		if err != nil {
			return nil, fmt.Errorf("get budget status: %w", err)
		}

		budgetWithStatus := &models.BudgetWithStatus{
			Budget: budget,
			Status: status,
		}
		budgetsWithStatus = append(budgetsWithStatus, budgetWithStatus)
	}

	return budgetsWithStatus, nil
}

func (s *BudgetService) GetBudgetStatus(userID string, categoryID int64, year, month int) (*models.BudgetStatus, error) {
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}

	budget, err := s.budgetRepo.GetBudgetByCategoryAndMonth(categoryID, year, month)
	if err != nil {
		return nil, fmt.Errorf("get budget: %w", err)
	}

	if budget.AccountID != account.ID {
		return nil, fmt.Errorf("budget does not belong to user")
	}

	status, err := s.getBudgetStatus(budget, year, month)
	if err != nil {
		return nil, fmt.Errorf("get budget status: %w", err)
	}

	return status, nil
}

func (s *BudgetService) getBudgetStatus(budget *models.Budget, year, month int) (*models.BudgetStatus, error) {
	spentAmount, err := s.transactionRepo.GetSpentAmountByCategoryAndMonth(budget.CategoryID, year, month)
	if err != nil {
		return nil, fmt.Errorf("get spent amount: %w", err)
	}

	remainingAmount := budget.Amount - spentAmount

	var progress float64
	if budget.Amount > 0 {
		progress = (spentAmount / budget.Amount) * 100
	}

	isExceeded := spentAmount > budget.Amount

	status := &models.BudgetStatus{
		Budget:     budget,
		Spent:      spentAmount,
		Remaining:  remainingAmount,
		Progress:   progress,
		IsExceeded: isExceeded,
	}

	return status, nil
}

func (s *BudgetService) GetBudgetSummary(userID string, year, month int) (*models.BudgetSummary, error) {
	budgets, err := s.GetBudgets(userID, year, month)
	if err != nil {
		return nil, fmt.Errorf("get monthly budgets: %w", err)
	}

	var totalPlanned float64
	var totalSpent float64
	var totalRemaining float64

	for _, budgetWithStatus := range budgets {
		totalPlanned += budgetWithStatus.Budget.Amount
		totalSpent += budgetWithStatus.Status.Spent
		totalRemaining += budgetWithStatus.Status.Remaining
	}

	summary := &models.BudgetSummary{
		TotalPlanned:   totalPlanned,
		TotalSpent:     totalSpent,
		TotalRemaining: totalRemaining,
		IsOverBudget:   totalSpent > totalPlanned,
		Budgets:        budgets,
	}

	return summary, nil
}
