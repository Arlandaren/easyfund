BEGIN;

-- Справочники удаляются последними, чтобы не ломать внешние ключи
DROP TABLE IF EXISTS credit_application_types;
DROP TABLE IF EXISTS credit_application_statuses;
DROP TABLE IF EXISTS banks;

COMMIT;
