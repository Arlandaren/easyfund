-- 000010_seed_demo_data_down.sql

DELETE FROM transactions
WHERE user_id IN (SELECT user_id FROM users WHERE email IN
  ('ivan@example.com','anna@example.com','pavel@example.com','olga@example.com','sergey@example.com'));

DELETE FROM loan_payments
WHERE loan_id IN (
  SELECT l.loan_id
  FROM loans l
  JOIN users u ON u.user_id = l.user_id
  WHERE u.email IN ('ivan@example.com','anna@example.com','pavel@example.com','olga@example.com','sergey@example.com')
);

DELETE FROM loans
WHERE user_id IN (SELECT user_id FROM users WHERE email IN
  ('ivan@example.com','anna@example.com','pavel@example.com','olga@example.com','sergey@example.com'));

DELETE FROM user_bank_accounts
WHERE user_id IN (SELECT user_id FROM users WHERE email IN
  ('ivan@example.com','anna@example.com','pavel@example.com','olga@example.com','sergey@example.com'));

DELETE FROM users
WHERE email IN ('ivan@example.com','anna@example.com','pavel@example.com','olga@example.com','sergey@example.com');