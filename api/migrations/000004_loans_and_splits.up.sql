BEGIN;

CREATE TABLE loans (
  loan_id bigserial PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  original_amount numeric(18,2) NOT NULL,
  taken_at timestamptz NOT NULL DEFAULT now(),
  interest_rate numeric(5,2) NOT NULL,
  status text NOT NULL DEFAULT 'ACTIVE',
  purpose text,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE loan_splits (
  split_id bigserial PRIMARY KEY,
  loan_id bigint NOT NULL REFERENCES loans(loan_id) ON DELETE CASCADE,
  bank_id smallint NOT NULL REFERENCES banks(bank_id),
  split_amount numeric(18,2) NOT NULL,
  remaining_principal numeric(18,2) NOT NULL,
  UNIQUE(loan_id, bank_id)
);

COMMIT;
