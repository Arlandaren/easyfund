BEGIN;

DROP INDEX IF EXISTS idx_transactions_user;

DROP TABLE IF EXISTS transactions;

COMMIT;
