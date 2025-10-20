CREATE TABLE notifications (
                               id SERIAL PRIMARY KEY,
                               user_id VARCHAR(255) NOT NULL,
                               type VARCHAR(50) NOT NULL,           -- 'budget_exceeded', 'low_balance', 'budget_warning'
                               title VARCHAR(255) NOT NULL,
                               message TEXT NOT NULL,
                               data JSONB,                         -- дополнительные данные
                               is_read BOOLEAN DEFAULT FALSE,
                               priority VARCHAR(20) DEFAULT 'medium', -- 'low', 'medium', 'high'
                               created_at TIMESTAMP DEFAULT NOW(),
                               updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);