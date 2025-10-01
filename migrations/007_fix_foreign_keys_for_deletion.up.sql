-- Удаляем старые ограничения внешних ключей
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_bank_account_id_fkey;
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_to_account_id_fkey;

-- Добавляем новые ограничения с CASCADE для удаления
ALTER TABLE transactions ADD CONSTRAINT transactions_bank_account_id_fkey 
    FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE CASCADE;

ALTER TABLE transactions ADD CONSTRAINT transactions_to_account_id_fkey 
    FOREIGN KEY (to_account_id) REFERENCES bank_accounts(id) ON DELETE SET NULL;
