CREATE TABLE user_notification_settings (
                                            id SERIAL PRIMARY KEY,
                                            user_id VARCHAR(255) UNIQUE NOT NULL,
                                            budget_alerts_enabled BOOLEAN DEFAULT TRUE,
                                            balance_alerts_enabled BOOLEAN DEFAULT TRUE,
                                            budget_warning_percent INTEGER DEFAULT 80,  -- процент для предупреждения
                                            low_balance_threshold DECIMAL(10,2) DEFAULT 100.00,
                                            preferred_channel VARCHAR(20) DEFAULT 'email', -- 'email', 'push', 'sms'
                                            created_at TIMESTAMP DEFAULT NOW(),
                                            updated_at TIMESTAMP DEFAULT NOW()
);