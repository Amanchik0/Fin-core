package services

import (
	"fmt"
	"justTest/internal/events"
	"justTest/internal/interfaces"
	"justTest/internal/models"
	"log"
	"time"
)

type BudgetService struct {
	budgetRepo      interfaces.BudgetRepository
	transactionRepo interfaces.TransactionRepository
	accountRepo     interfaces.AccountRepository
	categoryRepo    interfaces.CategoryRepository
	publisher       *events.Publisher
}

func NewBudgetService(
	budgetRepo interfaces.BudgetRepository,
	transactionRepo interfaces.TransactionRepository,
	accountRepo interfaces.AccountRepository,
	categoryRepo interfaces.CategoryRepository,
	publisher *events.Publisher,
) *BudgetService {
	return &BudgetService{
		budgetRepo:      budgetRepo,
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		publisher:       publisher,
	}
}
func (s *BudgetService) CheckBudgetAfterTransaction(event events.TransactionCreatedEvent) error {
	log.Printf("[BudgetService] Checking budget for transaction: ID=%d, CategoryID=%d, Amount=%.2f",
		event.TransactionID, event.CategoryID, event.Amount)
	if event.Amount >= 0 {
		log.Printf("[BudgetService] Skipping: not an expense (amount=%.2f)", event.Amount)
		return nil
	}
	//category, err := s.categoryRepo.GetByID(event.CategoryID)
	//if err != nil {
	//	return fmt.Errorf("get category: %w", err)
	//}
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	budget, err := s.budgetRepo.GetBudgetByCategoryAndMonth(event.CategoryID, year, month)
	if err != nil {
		log.Printf("[BudgetService] No budget found for category %d in %d-%02d", event.CategoryID, year, month)
		return nil
	}
	account, err := s.accountRepo.GetByUserID(event.UserID)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}
	if budget.AccountID != account.ID {
		log.Printf("[BudgetService] Budget does not belong to user")
		return nil
	}
	spentAmount, err := s.transactionRepo.GetSpentAmountByCategoryAndMonth(event.CategoryID, year, month)
	if err != nil {
		return fmt.Errorf("get spent amount: %w", err)
	}
	log.Printf("[BudgetService] Budget check: Spent=%.2f / Limit=%.2f (%.0f%%)",
		spentAmount, budget.Amount, (spentAmount/budget.Amount)*100)
	if spentAmount > budget.Amount {
		excessAmount := spentAmount - budget.Amount
		log.Printf("[BudgetService] ⚠️ Budget EXCEEDED by %.2f", excessAmount)

		// Публикуем событие превышения бюджета
		if s.publisher != nil {
			err := s.publisher.PublishBudgetExceeded(events.BudgetExceededEvent{
				UserID:       event.UserID,
				BudgetID:     budget.ID,
				BudgetName:   budget.BudgetLimitName,
				BudgetAmount: budget.Amount,
				SpentAmount:  spentAmount,
				ExcessAmount: excessAmount,
				CategoryID:   event.CategoryID,
				Timestamp:    time.Now(),
			})
			if err != nil {
				log.Printf("[BudgetService] Error publishing BudgetExceeded event: %v", err)
			} else {
				log.Printf("[BudgetService] ✅ Published BudgetExceeded event")
			}
		}
		return nil
	}
	warningThreshold := budget.Amount * 0.80
	if spentAmount >= warningThreshold {
		percentUsed := (spentAmount / budget.Amount) * 100
		log.Printf("[BudgetService] ⚠ Budget WARNING: %.0f%% used", percentUsed)

		// Публикуем предупреждение
		if s.publisher != nil {
			err := s.publisher.PublishBudgetWarning(events.BudgetWarningEvent{
				UserID:         event.UserID,
				BudgetID:       budget.ID,
				BudgetName:     budget.BudgetLimitName,
				BudgetAmount:   budget.Amount,
				SpentAmount:    spentAmount,
				WarningPercent: percentUsed,
				CategoryID:     event.CategoryID,
				Timestamp:      time.Now(),
			})
			if err != nil {
				log.Printf("[BudgetService] Error publishing BudgetWarning event: %v", err)
			} else {
				log.Printf("[BudgetService]  Published BudgetWarning event")
			}
		}
	}

	log.Printf("[BudgetService]  Budget check completed")
	return nil
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
