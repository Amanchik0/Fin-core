package repo

import (
	"database/sql"
	"fmt"
	"justTest/internal/models"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}
func (r *AccountRepository) Create(account *models.Account) (*models.Account, error) {
	query := ` 
	insert into accounts (user_id, name, display_name, timezone, is_active ,created_at, updated_at )
	values ($1, $2, $3, $4, $5, $6, $7)
	returning id;`
	err := r.db.QueryRow(query,
		account.UserID,
		account.Name,
		account.DisplayName,
		account.Timezone,
		account.IsActive,
		account.CreatedAt,
		account.UpdatedAt).Scan(&account.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %v", err)
	}
	return account, nil

}

func (r *AccountRepository) GetByUserID(userID string) (*models.Account, error) {
	query := `
    SELECT id, user_id, name, display_name, timezone, is_active, created_at, updated_at 
    FROM accounts
    WHERE user_id = $1 AND is_active = true`

	var account models.Account
	err := r.db.QueryRow(query, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.DisplayName,
		&account.Timezone,
		&account.IsActive,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Account not found") // ✅ Возвращаем ошибку
		}
		return nil, fmt.Errorf("error getting account: %v", err)
	}

	return &account, nil
}
func (r *AccountRepository) GetByID(id int64) (*models.Account, error) {
	query := `
	select id, user_id, name, display_name, timezone, is_active, created_at, updated_at 
	from accounts 
	where id = $1`
	account := &models.Account{}

	err := r.db.QueryRow(query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.DisplayName,
		&account.Timezone,
		&account.IsActive,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil

}
