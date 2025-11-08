BEGIN;

CREATE TABLE user_bank_accounts (
  account_id bigserial PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  bank_id smallint NOT NULL REFERENCES banks(bank_id),
  balance numeric(18,2) NOT NULL DEFAULT 0.00,
  currency char(3) NOT NULL DEFAULT 'RUB',
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE(user_id, bank_id)
);

CREATE INDEX idx_user_bank_accounts_user ON user_bank_accounts(user_id);

COMMIT;
