package repo

import (
	"database/sql"
	"fmt"
	"justTest/internal/models"
)

type BankAccountRepository struct {
	db *sql.DB
}

func NewBankAccountRepository(db *sql.DB) *BankAccountRepository {
	return &BankAccountRepository{
		db: db,
	}
}
func (r *BankAccountRepository) Create(bankAccount *models.BankAccount) (*models.BankAccount, error) {
	query := `
		insert into bank_accounts  ( account_id, name, currency, account_type, bank_name, is_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		returning id;
		
`
	err := r.db.QueryRow(query,
		bankAccount.AccountID,
		bankAccount.Name,
		bankAccount.Currency,
		bankAccount.AccountType,
		bankAccount.BankName,
		bankAccount.IsActive,
		bankAccount.CreatedAt,
		bankAccount.UpdatedAt,
	).Scan(&bankAccount.ID)
	if err != nil {
		return nil, fmt.Errorf("Error creating bank account: %v", err)

	}
	return bankAccount, nil
}

func (r *BankAccountRepository) GetByAccountID(accountID int64) ([]*models.BankAccount, error) {
	query := `
	select id, account_id, name, currency, account_type, bank_name, is_active, created_at, updated_at
	from bank_accounts 
	where account_id = $1;
`
	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("Error getting bank account: %v", err)
	}
	defer rows.Close()
	var accounts []*models.BankAccount = make([]*models.BankAccount, 0)
	for rows.Next() {
		account := &models.BankAccount{}
		err := rows.Scan(
			&account.ID,
			&account.AccountID,
			&account.Name,
			&account.Currency,
			&account.AccountType,
			&account.IsActive,
			&account.BankName,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		accounts = append(accounts, account)
		if err != nil {
			return nil, fmt.Errorf("Error getting bank account: %v", err)
		}
	}
	return accounts, nil

}
func (r *BankAccountRepository) GetActiveBankAccounts(accountID int64) ([]*models.BankAccount, error) {
	query := `
	select id, account_id, name, currency, account_type, bank_name, is_active, created_at, updated_at
	from bank_accounts
	where account_id = $1 and is_active = true;
`
	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("Error getting bank account: %v", err)
	}
	defer rows.Close()
	var accounts []*models.BankAccount = make([]*models.BankAccount, 0)
	for rows.Next() {
		account := &models.BankAccount{}
		err := rows.Scan(
			&account.ID,
			&account.AccountID,
			&account.Name,
			&account.Currency,
			&account.AccountType,
			&account.IsActive,
			&account.BankName,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		accounts = append(accounts, account)

		if err != nil {
			return nil, fmt.Errorf("Error getting bank account: %v", err)
		}
	}
	return accounts, nil

}
func (r *BankAccountRepository) GetBankAccountByCurrency(accountID int64, currency string) ([]*models.BankAccount, error) {
	query := `
select id, account_id, name, currency, account_type, bank_name, is_active, created_at, updated_at
from bank_accounts
where account_id = $1 and currency = $2;
`
	rows, err := r.db.Query(query, accountID, currency)
	if err != nil {
		return nil, fmt.Errorf("Error getting bank account: %v", err)
	}
	defer rows.Close()
	var accounts []*models.BankAccount = make([]*models.BankAccount, 0)
	for rows.Next() {
		account := &models.BankAccount{}
		err := rows.Scan(
			&account.ID,
			&account.AccountID,
			&account.Name,
			&account.Currency,
			&account.AccountType,
			&account.IsActive,
			&account.BankName,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		accounts = append(accounts, account)

		if err != nil {
			return nil, fmt.Errorf("Error getting bank account: %v", err)

		}

	}
	return accounts, nil
}
func (r *BankAccountRepository) GetByBankAccountID(BankAccountID int64) (*models.BankAccount, error) {
	query := `
		select id, account_id, name, currency, account_type, bank_name, is_active, created_at, updated_at
from bank_accounts
where id = $1; `
	bankAccount := &models.BankAccount{}
	err := r.db.QueryRow(query, BankAccountID).Scan(
		&bankAccount.ID,
		&bankAccount.AccountID,
		&bankAccount.Name,
		&bankAccount.Currency,
		&bankAccount.AccountType,
		&bankAccount.IsActive,
		&bankAccount.BankName,
		&bankAccount.CreatedAt,
		&bankAccount.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(`No bank account found with id %d`, BankAccountID)

		}
		return nil, fmt.Errorf("Error getting bank account: %v", err)

	}
	return bankAccount, nil
}

// middle checks

func (r *BankAccountRepository) ExsitsAccountIDAndName(accountID int64, name string) (bool, error) {
	query := `
select count(*) from bank_accounts where account_id = $1 and name = $2;
`
	var count int

	err := r.db.QueryRow(query, accountID, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("Error getting bank account: %v", err)
	}
	return count > 0, nil
}
func (r *BankAccountRepository) DeActiveBankAccount(bankAccountID int64) error {
	query := `
update bank_accounts set is_active = false , updated_at = now() where id = $1`
	res, err := r.db.Exec(query, bankAccountID)
	if err != nil {
		return fmt.Errorf("Error deleting bank account: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error deleting bank account: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf(`No bank account found with id %d`, bankAccountID)
	}
	return nil

}
func (r *BankAccountRepository) DeleteBankAccount(bankAccountID int64) error {
	query := `
delete from bank_accounts where id = $1`
	res, err := r.db.Exec(query, bankAccountID)
	if err != nil {
		return fmt.Errorf("Error deleting bank account: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error deleting bank account: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf(`No bank account found with id %d`, bankAccountID)

	}
	return nil

}
