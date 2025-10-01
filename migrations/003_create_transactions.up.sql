CREATE TABLE IF NOT EXISTS transactions (
                                            id BIGSERIAL PRIMARY KEY,
                                            bank_account_id BIGINT NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
    category_id BIGINT,
    amount DECIMAL(15,2) NOT NULL,
    description TEXT NOT NULL,
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('income', 'expense', 'transfer')),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
                                                                                                             to_account_id BIGINT REFERENCES bank_accounts(id),
    transfer_rate DECIMAL(15,6),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at timestamp with time zone default now()
    );

CREATE INDEX idx_transactions_bank_account_id ON transactions(bank_account_id);
CREATE INDEX idx_transactions_date ON transactions(date);
CREATE INDEX idx_transactions_type ON transactions(transaction_type);
CREATE INDEX idx_transactions_category_id ON transactions(category_id);