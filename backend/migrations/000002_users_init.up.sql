BEGIN;

-- users: PK как BIGSERIAL
CREATE TABLE users (
  user_id BIGSERIAL PRIMARY KEY,
  full_name text NOT NULL,
  email text UNIQUE NOT NULL,
  phone text,
  password_hash text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);


CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN NEW.updated_at := now(); RETURN NEW; END $$;

CREATE TRIGGER trg_users_updated
BEFORE UPDATE ON users FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

COMMIT;