BEGIN;


CREATE TABLE loan_payments (
  payment_id BIGSERIAL PRIMARY KEY,
  loan_id BIGINT NOT NULL REFERENCES loans(loan_id) ON DELETE CASCADE,
  amount numeric(18,2) NOT NULL,
  paid_at timestamptz NOT NULL,
  method text NOT NULL,
  status text NOT NULL DEFAULT 'posted'
);

CREATE TABLE payment_allocations (
  allocation_id bigserial PRIMARY KEY,
  payment_id bigint NOT NULL REFERENCES loan_payments(payment_id) ON DELETE CASCADE,
  split_id bigint NOT NULL REFERENCES loan_splits(split_id) ON DELETE CASCADE,
  principal_paid numeric(18,2) NOT NULL,
  interest_paid numeric(18,2) NOT NULL
);

COMMIT;