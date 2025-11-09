BEGIN;

-- splits зависит от loans, поэтому сначала splits
DROP TABLE IF EXISTS loan_splits;

DROP TABLE IF EXISTS loans;

COMMIT;
