package repo

import (
	"database/sql"
	"fmt"
	"justTest/internal/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db}

}
func (r *TransactionRepository) Create(transaction *models.Transaction) (*models.Transaction, error) {
	query := `
insert into transactions ( bank_account_id, category_id, amount, description, transaction_type, date, created_at, updated_at, to_account_id, transfer_rate)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10)
	returning id;`
	err := r.db.QueryRow(query,
		transaction.BankAccountID,
		transaction.CategoryID,
		transaction.Amount,
		transaction.Description,
		transaction.TransactionType,
		transaction.Date,
		transaction.CreatedAt,
		transaction.UpdatedAt,
		transaction.ToAccountID,
		transaction.TransferRate,
	).Scan(&transaction.ID)

	if err != nil {
		return transaction, fmt.Errorf("error creating transaction: %v", err)
	}
	return transaction, nil
}
func (r *TransactionRepository) GetByBankAccountID(BankAccountID int64) ([]*models.Transaction, error) {
	query := ` 
	select id, bank_account_id, category_id, amount, description, transaction_type, 
               date, created_at, updated_at, to_account_id, transfer_rate from transactions where bank_account_id = $1
`
	rows, err := r.db.Query(query, BankAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []*models.Transaction = make([]*models.Transaction, 0)
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.BankAccountID,
			&transaction.CategoryID,
			&transaction.Amount,
			&transaction.Description,
			&transaction.TransactionType,
			&transaction.Date,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.ToAccountID,
			&transaction.TransferRate,
		)
		transactions = append(transactions, transaction)
		if err != nil {
			return transactions, fmt.Errorf("error getting transaction: %v", err)
		}
	}
	return transactions, nil

}

func (r *TransactionRepository) GetByCategoryID(CategoryID int64) ([]*models.Transaction, error) {
	query := ` 
	select id, bank_account_id, category_id, amount, description, transaction_type, 
               date, created_at, updated_at, to_account_id, transfer_rate from transactions where category_id = $1
`
	rows, err := r.db.Query(query, CategoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []*models.Transaction = make([]*models.Transaction, 0)
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.BankAccountID,
			&transaction.CategoryID,
			&transaction.Amount,
			&transaction.Description,
			&transaction.TransactionType,
			&transaction.Date,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.ToAccountID,
			&transaction.TransferRate,
		)
		transactions = append(transactions, transaction)
		if err != nil {
			return transactions, fmt.Errorf("error getting transaction: %v", err)
		}
	}
	return transactions, nil

}
func (r *TransactionRepository) GetByTransactionID(TransactionID int64) (*models.Transaction, error) {
	query := ` 
	select id, bank_account_id, category_id, amount, description, transaction_type, 
               date, created_at, updated_at, to_account_id, transfer_rate 
	from transactions where id = $1
`
	transaction := &models.Transaction{}
	err := r.db.QueryRow(query, TransactionID).Scan(
		&transaction.ID,
		&transaction.BankAccountID,
		&transaction.CategoryID,
		&transaction.Amount,
		&transaction.Description,
		&transaction.TransactionType,
		&transaction.Date,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.ToAccountID,
		&transaction.TransferRate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with id %d not found", TransactionID)
		}
		return nil, fmt.Errorf("error getting transaction: %v", err)
	}
	return transaction, nil

}

func (r *TransactionRepository) GetByAccountID(AccountID int64) ([]*models.Transaction, error) {
	query := ` 
	select t.id , t.bank_account_id, t.category_id, t.amount, t.description, t.transaction_type,
	t.date, t.created_at, t.updated_at, to_account_id, t.transfer_rate 
from transactions t 
	join bank_accounts ba on t.bank_account_id = ba.id
	where ba.account_id = $1
	order by t.date desc
`
	rows, err := r.db.Query(query, AccountID)
	if err != nil {
		return nil, err

	}
	defer rows.Close()
	transactions := make([]*models.Transaction, 0)
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.BankAccountID,
			&transaction.CategoryID,
			&transaction.Amount,
			&transaction.Description,
			&transaction.TransactionType,
			&transaction.Date,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.ToAccountID,
			&transaction.TransferRate,
		)
		if err != nil {
			return transactions, fmt.Errorf("error getting transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *TransactionRepository) GetTransfersByAccountID(AccountID int64) ([]*models.Transaction, error) {
	query := `
	select t.id, t.bank_account_id, t.category_id, t.amount, t.description, t.transaction_type,
	t.date, t.created_at, t.updated_at, to_account_id, t.transfer_rate
from transactions t
	join bank_accounts ba on t.bank_account_id = ba.id
	where ba.account_id = $1 and t.transaction_type = 'transfer'
	order by t.date desc
`
	rows, err := r.db.Query(query, AccountID)
	if err != nil {
		return nil, err

	}
	defer rows.Close()
	transactions := make([]*models.Transaction, 0)
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.BankAccountID,
			&transaction.CategoryID,
			&transaction.Amount,
			&transaction.Description,
			&transaction.TransactionType,
			&transaction.Date,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.ToAccountID,
			&transaction.TransferRate)
		if err != nil {
			return transactions, fmt.Errorf("error getting transaction: %v", err)
		}
		transactions = append(transactions, transaction)

	}
	return transactions, nil

}

func (r *TransactionRepository) CreateTransaction(AccountID, FromBankAccountID, categoryID *int64, toAccountID int64, amount float64, description string, transferRate *float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	outgoingQuery := ` 
	insert into transactions (bank_account_id, category_id, amount, description, transaction_type,date , 
	                          created_at, updated_at, to_account_id, transfer_rate)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = tx.Exec(outgoingQuery,
		FromBankAccountID,
		categoryID,
		-amount,
		description,
		"transfer",
		now,
		now,
		now,
		toAccountID,
		transferRate)
	if err != nil {
		return fmt.Errorf("error creating transaction: %v", err)
	}
	incomingAmount := amount
	if transferRate != nil {
		incomingAmount = amount * (*transferRate)
	}

	incomingQuery := ` 
	insert into transactions (bank_account_id, category_id, amount, description, 
	                          transaction_type, date , created_at , updated_at, to_account_id, transfer_rate )
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`
	_, err = tx.Exec(incomingQuery,
		toAccountID,
		categoryID,
		incomingAmount,
		description,
		"transfer",
		now,
		now,
		now,
		FromBankAccountID,
		transferRate)
	if err != nil {
		return fmt.Errorf("error creating transaction: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func (r *TransactionRepository) GetTotalAmountByBankAccountID(BankAccountID int64) (float64, error) {

	query := `
	select sum(amount) from transactions where bank_account_id = $1
`
	row := r.db.QueryRow(query, BankAccountID)
	var amount float64
	err := row.Scan(&amount)
	if err != nil {
		return amount, fmt.Errorf("error getting total amount by bank account id: %v", err)
	}
	return amount, nil

}
