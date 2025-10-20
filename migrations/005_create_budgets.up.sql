CREATE TABLE IF NOT EXISTS budgets (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    budget_limit_name VARCHAR(255) NOT NULL,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL,
    period VARCHAR(20) NOT NULL CHECK (period IN ('monthly', 'weekly', 'yearly')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_budget_per_category_period UNIQUE (category_id, period, start_date)
);

CREATE INDEX idx_budgets_account_id ON budgets(account_id);
CREATE INDEX idx_budgets_category_id ON budgets(category_id);
CREATE INDEX idx_budgets_period ON budgets(period);
CREATE INDEX idx_budgets_is_active ON budgets(is_active);
