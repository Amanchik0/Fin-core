CREATE TABLE IF NOT EXISTS bank_accounts (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    account_type VARCHAR(50) NOT NULL CHECK (account_type IN ('cash', 'debit', 'credit', 'savings')),
    bank_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT unique_account_name UNIQUE (account_id, name)
    );

CREATE INDEX idx_bank_accounts_account_id ON bank_accounts(account_id);
CREATE INDEX idx_bank_accounts_currency ON bank_accounts(currency);
CREATE INDEX idx_bank_accounts_is_active ON bank_accounts(is_active);