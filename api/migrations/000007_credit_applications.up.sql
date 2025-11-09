BEGIN;

CREATE TABLE credit_applications (
  application_id bigserial PRIMARY KEY,
  user_id bigint NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  bank_id smallint NOT NULL REFERENCES banks(bank_id),
  type_code text NOT NULL REFERENCES credit_application_types(type_code),
  status_code text NOT NULL REFERENCES credit_application_statuses(status_code),
  requested_amount numeric(18,2) NOT NULL,
  loan_id bigint REFERENCES loans(loan_id) ON DELETE SET NULL,
  submitted_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION set_ca_updated_at() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN NEW.updated_at := now(); RETURN NEW; END $$;

CREATE TRIGGER trg_credit_applications_updated
BEFORE UPDATE ON credit_applications FOR EACH ROW
EXECUTE FUNCTION set_ca_updated_at();

COMMIT;