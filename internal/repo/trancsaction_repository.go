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
insert into transactions ( bank_account_id, category_id, amount, description, transaction_type, date, 
                          created_at, updated_at, to_account_id, transfer_rate)
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

func (r *TransactionRepository) GetByBankAccountID(BankAccountID int64, limit, offset int) ([]*models.Transaction, error) {
	query := ` 
	select id, bank_account_id, category_id, amount, description, transaction_type, 
               date, created_at, updated_at, to_account_id, transfer_rate from transactions where bank_account_id = $1
	order by created_at desc limit $2 offset $3
`
	rows, err := r.db.Query(query, BankAccountID, limit, offset)
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
		if err != nil {
			return transactions, fmt.Errorf("error getting transaction: %v", err)
		}
		transactions = append(transactions, transaction)

	}
	return transactions, nil

}

// TODO надо везде дату добавить а то есть лимиты но нет даты я хз на каком уровне его добавить но надо
func (r *TransactionRepository) GetByCategoryID(CategoryID int64, limit, offset int) ([]*models.Transaction, error) {
	query := ` 
	select id, bank_account_id, category_id, amount, description, transaction_type, 
               date, created_at, updated_at, to_account_id, transfer_rate 
	from transactions where category_id = $1
		order by created_at desc limit $2 offset $3

`
	rows, err := r.db.Query(query, CategoryID, limit, offset)
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

func (r *TransactionRepository) GetByAccountID(AccountID int64, limit, offset int) ([]*models.Transaction, error) {
	query := ` 
	select t.id , t.bank_account_id, t.category_id, t.amount, t.description, t.transaction_type,
	t.date, t.created_at, t.updated_at, to_account_id, t.transfer_rate 
from transactions t 
	join bank_accounts ba on t.bank_account_id = ba.id
	where ba.account_id = $1
	order by t.date desc
	limit $2 offset $3
	
`
	rows, err := r.db.Query(query, AccountID, limit, offset)
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

func (r *TransactionRepository) GetSpentAmountByCategoryAndMonth(categoryID int64, year, month int) (float64, error) {
	query := ` 
	select COALESCE(SUM(ABS(amount)), 0) 
-- 	    as total // можно тотал убрать и после amount умножить все в тг 
from transactions where category_id =$1 
	and transaction_type ='expense' 
	and extract(year from date) =$2 
	and extract(month from date) =$3
-- 	group by currency; // хз вот убрать или нет 
`
	row := r.db.QueryRow(query, categoryID, year, month)
	var amount float64
	err := row.Scan(&amount)
	if err != nil {
		return amount, fmt.Errorf("error getting total amount by bank account id: %v", err)
	}
	return amount, nil
}

func (r *TransactionRepository) GetTransactionsByCategoryAndMonth(categoryID int64, year, month int, limit, offset int) ([]*models.Transaction, error) {
	query := `
        SELECT id, bank_account_id, category_id, amount, description, transaction_type, 
               date, created_at, updated_at, to_account_id, transfer_rate 
        FROM transactions 
        WHERE category_id = $1 
        AND EXTRACT(YEAR FROM date) = $2 
        AND EXTRACT(MONTH FROM date) = $3
        ORDER BY date DESC 
        LIMIT $4 OFFSET $5
    `
	rows, err := r.db.Query(query, categoryID, year, month, limit, offset)
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

func (r *TransactionRepository) GetTransactionsByDateRangeWithCategory(categoryID int64, startDate, endDate time.Time, limit, offset int) ([]*models.Transaction, error) {
	query := `
     SELECT id, bank_account_id, category_id, amount, description, transaction_type, 
               date, created_at, updated_at, to_account_id, transfer_rate 
        FROM transactions
        where category_id = $1
    AND date >= $2 
        AND date <= $3
        ORDER BY date DESC 
        LIMIT $4 OFFSET $5        `
	rows, err := r.db.Query(query, categoryID, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions by date range: %v", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
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
			return transactions, fmt.Errorf("error scanning transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *TransactionRepository) GetTransactionsByDateRange(startDate, endDate time.Time, limit, offset int) ([]*models.Transaction, error) {
	query := `
     SELECT id, bank_account_id, category_id, amount, description, transaction_type, 
               date, created_at, updated_at, to_account_id, transfer_rate 
        FROM transactions
        where date >= $1
        AND date <= $2
        ORDER BY date DESC 
        LIMIT $3 OFFSET $4        `
	rows, err := r.db.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions by date range: %v", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
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
			return transactions, fmt.Errorf("error scanning transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
