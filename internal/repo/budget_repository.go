package repo

import (
	"database/sql"
	"fmt"
	"justTest/internal/models"
)

type BudgetRepository struct {
	db *sql.DB
}

func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{
		db: db,
	}
}

func (r *BudgetRepository) CreateBudget(budget *models.Budget) (*models.Budget, error) {
	query := ` insert into budgets (
                     account_id , budget_limit_name, category_id, amount, 
                     period, start_date, end_date,is_active, created_at, updated_at
 ) 
 values ($1, $2,$3,$4,$5,$6,$7,$8,$9, $10)
 returning id `

	err := r.db.QueryRow(query,
		budget.AccountID,
		budget.BudgetLimitName,
		budget.CategoryID,
		budget.Amount,
		budget.Period,
		budget.StartDate,
		budget.EndDate,
		budget.IsActive,
		budget.CreatedAt,
		budget.UpdatedAt,
	).Scan(&budget.ID)
	if err != nil {
		return nil, err

	}
	return budget, nil

}
func (r *BudgetRepository) GetBudget(budgetID int64) (*models.Budget, error) {
	query := ` select id , account_id , budget_limit_name, category_id, amount, 
                     period, start_date, end_date,is_active, created_at, updated_at from budgets where id = $1`
	row := r.db.QueryRow(query, budgetID)
	budget := &models.Budget{}
	err := row.Scan(&budget.ID,
		&budget.AccountID,
		&budget.BudgetLimitName,
		&budget.CategoryID,
		&budget.Amount,
		&budget.Period,
		&budget.StartDate,
		&budget.EndDate,
		&budget.IsActive,
		&budget.CreatedAt,
		&budget.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(`no budget found with id %d`, budgetID)
		}
		return nil, err
	}
	return budget, nil
}
func (r *BudgetRepository) GetBudgetByCategoryID(categoryID int64) (*models.Budget, error) {
	query := ` select id , account_id , budget_limit_name, category_id, amount, 
                     period, start_date, end_date,is_active, created_at, updated_at from budgets where category_id = $1`
	row := r.db.QueryRow(query, categoryID)
	budget := &models.Budget{}
	err := row.Scan(&budget.ID,
		&budget.AccountID,
		&budget.BudgetLimitName,
		&budget.CategoryID,
		&budget.Amount,
		&budget.Period,
		&budget.StartDate,
		&budget.EndDate,
		&budget.IsActive,
		&budget.CreatedAt,
		&budget.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(`no budget found with id %d`, categoryID)
		}
		return nil, err
	}
	return budget, nil
}

func (r *BudgetRepository) GetBudgetByCategoryAndMonth(categoryID int64, year, month int) (*models.Budget, error) {
	query := `
		SELECT id, account_id, budget_limit_name, category_id, amount, 
		       period, start_date, end_date, is_active, created_at, updated_at 
		FROM budgets 
		WHERE category_id = $1 
		AND EXTRACT(YEAR FROM start_date) = $2 
		AND EXTRACT(MONTH FROM start_date) = $3
		AND is_active = true
	`
	row := r.db.QueryRow(query, categoryID, year, month)
	budget := &models.Budget{}
	err := row.Scan(
		&budget.ID,
		&budget.AccountID,
		&budget.BudgetLimitName,
		&budget.CategoryID,
		&budget.Amount,
		&budget.Period,
		&budget.StartDate,
		&budget.EndDate,
		&budget.IsActive,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no budget found for category %d in %d-%02d", categoryID, year, month)
		}
		return nil, err
	}
	return budget, nil
}

func (r *BudgetRepository) GetBudgetsByAccountAndMonth(accountID int64, year, month int) ([]*models.Budget, error) {
	query := `
		SELECT id, account_id, budget_limit_name, category_id, amount, 
		       period, start_date, end_date, is_active, created_at, updated_at 
		FROM budgets 
		WHERE account_id = $1 
		AND EXTRACT(YEAR FROM start_date) = $2 
		AND EXTRACT(MONTH FROM start_date) = $3
		AND is_active = true
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, accountID, year, month)
	if err != nil {
		return nil, fmt.Errorf("error getting budgets by account and month: %v", err)
	}
	defer rows.Close()

	var budgets []*models.Budget
	for rows.Next() {
		budget := &models.Budget{}
		err := rows.Scan(
			&budget.ID,
			&budget.AccountID,
			&budget.BudgetLimitName,
			&budget.CategoryID,
			&budget.Amount,
			&budget.Period,
			&budget.StartDate,
			&budget.EndDate,
			&budget.IsActive,
			&budget.CreatedAt,
			&budget.UpdatedAt,
		)
		if err != nil {
			return budgets, fmt.Errorf("error scanning budget: %v", err)
		}
		budgets = append(budgets, budget)
	}
	return budgets, nil
}

//func (r *BudgetRepository) UpdateBudget()
