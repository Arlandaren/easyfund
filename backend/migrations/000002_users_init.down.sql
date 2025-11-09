BEGIN;

DROP TRIGGER IF EXISTS trg_users_updated ON users;
DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS users;

COMMIT;
