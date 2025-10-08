package repo

import (
	"database/sql"
	"fmt"
	"justTest/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}

}
func (r *CategoryRepository) CreateCategory(category *models.Category) (*models.Category, error) {
	query := `
	insert into categories (account_id, name, type , color, icon, is_active, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8)
	returning id;`
	err := r.db.QueryRow(query,
		category.AccountID,
		category.Name,
		category.Type,
		category.Color,
		category.Icon,
		category.IsActive,
		category.CreatedAt,
		category.UpdatedAt,
	).Scan(&category.ID)
	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	return category, nil
}

func (r *CategoryRepository) UpdateCategory(category *models.Category) (*models.Category, error) {
	query := ` 

update categories 
set name = $1, type = $2, color = $3, icon = $4, is_active = $5, updated_at = $6 
	where id = $7
returning id;`

	row := r.db.QueryRow(query,
		category.Name,
		category.Type,
		category.Color,
		category.Icon,
		category.IsActive,
		category.UpdatedAt,
		category.ID)
	var updatedAt *models.Category
	if err := row.Scan(&updatedAt); err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}
	return updatedAt, nil
}
func (r *CategoryRepository) DeleteCategory(categoryID int64) error {
	query := `
	delete from categories where id = $1;
`
	_, err := r.db.Exec(query,
		categoryID,
	)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}
func (r *CategoryRepository) GetByAccountID(accountID int64) ([]*models.Category, error) {
	query := ` 
select id, name, type, color, icon, is_active, created_at, updated_at
from categories
where account_id = $1`
	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("get categories by id: %w", err)
	}
	defer rows.Close()
	var categories []*models.Category = make([]*models.Category, 0)
	for rows.Next() {
		category := &models.Category{}
		err := rows.Scan(&category.ID, &category.Name, &category.Type, &category.Color, &category.Icon, &category.IsActive, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("get categories by id: %w", err)
		}
		categories = append(categories, category)

	}
	return categories, nil
}
func (r *CategoryRepository) GetByID(categoryID int64) (*models.Category, error) {
	query := `
	select id, account_id, name, type, color, icon, is_active, created_at, updated_at 
 from categories where id = $1`
	category := &models.Category{}
	err := r.db.QueryRow(query, categoryID).Scan(
		&category.ID,
		&category.AccountID,
		&category.Name,
		&category.Type,
		&category.Color,
		&category.Icon,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(`no category with id %d `, categoryID)
		}
		return nil, fmt.Errorf(`smg get wrong  %d`, err)

	}
	return category, nil
}

//func (r *CategoryRepository) GetByType(name string) (*models.Category, error) {}
