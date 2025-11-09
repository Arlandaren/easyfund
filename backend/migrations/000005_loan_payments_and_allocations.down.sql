BEGIN;

-- allocation зависит от payments, поэтому сначала allocations
DROP TABLE IF EXISTS payment_allocations;

DROP TABLE IF EXISTS loan_payments;

COMMIT;
