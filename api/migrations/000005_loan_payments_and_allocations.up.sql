BEGIN;

CREATE TABLE loan_payments (
  payment_id bigserial PRIMARY KEY,
  loan_id bigint NOT NULL REFERENCES loans(loan_id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(user_id),
  paid_at timestamptz NOT NULL DEFAULT now(),
  total_amount numeric(18,2) NOT NULL,
  comment text
);

CREATE TABLE payment_allocations (
  allocation_id bigserial PRIMARY KEY,
  payment_id bigint NOT NULL REFERENCES loan_payments(payment_id) ON DELETE CASCADE,
  split_id bigint NOT NULL REFERENCES loan_splits(split_id) ON DELETE CASCADE,
  principal_paid numeric(18,2) NOT NULL,
  interest_paid numeric(18,2) NOT NULL
);

COMMIT;
