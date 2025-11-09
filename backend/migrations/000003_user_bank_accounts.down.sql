BEGIN;

DROP INDEX IF EXISTS idx_user_bank_accounts_user;

DROP TABLE IF EXISTS user_bank_accounts;

COMMIT;
