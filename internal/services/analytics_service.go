package services

import (
	"justTest/internal/models"
	"justTest/internal/repo"
	"time"
)

type AnalyticsService struct {
	transactionRepo repo.TransactionRepository
	accountRepo     repo.AccountRepository
}

func NewAnalyticsService(transactionRepo repo.TransactionRepository, accountRepo repo.AccountRepository) *AnalyticsService {
	return &AnalyticsService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

// GetMonthlyReport - отчет за месяц
func (s *AnalyticsService) GetMonthlyReport(userID string, year int, month int) (*models.MonthlyReport, error) {
	// TODO: Реализовать получение отчета за месяц
	return nil, nil
}

// GetCategorySpending - траты по категориям
func (s *AnalyticsService) GetCategorySpending(userID string, startDate, endDate time.Time) ([]*models.CategorySpending, error) {
	// TODO: Реализовать анализ трат по категориям
	return nil, nil
}

// GetIncomeVsExpenses - доходы vs расходы
func (s *AnalyticsService) GetIncomeVsExpenses(userID string, startDate, endDate time.Time) (*models.IncomeExpenseReport, error) {
	// TODO: Реализовать сравнение доходов и расходов
	return nil, nil
}
