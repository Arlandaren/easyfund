BEGIN;

DROP TRIGGER IF EXISTS trg_credit_applications_updated ON credit_applications;
DROP FUNCTION IF EXISTS set_ca_updated_at();

DROP TABLE IF EXISTS credit_applications;

COMMIT;
