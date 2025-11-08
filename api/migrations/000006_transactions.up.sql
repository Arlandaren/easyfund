BEGIN;

CREATE TABLE transactions (
  transaction_id bigserial PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  bank_id smallint NOT NULL REFERENCES banks(bank_id),
  occurred_at timestamptz NOT NULL DEFAULT now(),
  amount numeric(18,2) NOT NULL,
  category text,
  description text
);

CREATE INDEX idx_transactions_user ON transactions(user_id);

COMMIT;
