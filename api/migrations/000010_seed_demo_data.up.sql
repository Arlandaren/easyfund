-- 000010_seed_demo_data_up.sql
-- Согласовано с:
-- users(user_id, full_name, email, phone, password_hash)
-- banks(bank_id, ...)
-- user_bank_accounts(user_id, bank_id, balance, currency, created_at)
-- loans(loan_id, user_id, amount, rate, months, status, created_at)   -- Имена по 000004
-- loan_payments(loan_id, amount, paid_at, method, status)              -- Имена по 000005
-- transactions(transaction_id, user_id, bank_id, occurred_at, amount, category, description) -- По вашему DDL

-- 1) Пользователи
INSERT INTO users (full_name, email, phone, password_hash)
VALUES
  ('Ivan Petrov','ivan@example.com','+79998887761','$2y$10$Kjh5eRBHUIhi38Dc/9za9ORrAwDzn8qE0.ZgR1W5NsJDndjzk/mJ2'),
  ('Anna Sidorova','anna@example.com','+79998887762','$2y$10$Kjh5eRBHUIhi38Dc/9za9ORrAwDzn8qE0.ZgR1W5NsJDndjzk/mJ2'),
  ('Pavel Smirnov','pavel@example.com','+79998887763','$2y$10$Kjh5eRBHUIhi38Dc/9za9ORrAwDzn8qE0.ZgR1W5NsJDndjzk/mJ2'),
  ('Olga Ivanova','olga@example.com','+79998887764','$2y$10$Kjh5eRBHUIhi38Dc/9za9ORrAwDzn8qE0.ZgR1W5NsJDndjzk/mJ2'),
  ('Sergey Volkov','sergey@example.com','+79998887765','$2y$10$Kjh5eRBHUIhi38Dc/9za9ORrAwDzn8qE0.ZgR1W5NsJDndjzk/mJ2');

-- 2) Счета (user_bank_accounts без account_number)
WITH u AS (
  SELECT user_id, email
  FROM users
  WHERE email IN ('ivan@example.com','anna@example.com','pavel@example.com','olga@example.com','sergey@example.com')
),
ub AS (
  SELECT u.user_id, b.bank_id
  FROM u CROSS JOIN banks b
)
INSERT INTO user_bank_accounts (user_id, bank_id, balance, currency, created_at)
SELECT ub.user_id, ub.bank_id,
       round((random()*90000 + 10000)::numeric, 2),
       'RUB',
       now()
FROM ub;

-- 3) Кредиты (по одному на пользователя)
-- Используем original_amount и interest_rate, поля term_months нет в схеме
INSERT INTO loans (user_id, original_amount, interest_rate, status, purpose, created_at)
SELECT u.user_id, v.original_amount, v.interest_rate, 'ACTIVE', 'Demo loan purpose', now()
FROM (VALUES
  ('ivan@example.com'  , 200000::numeric, 14.9::numeric),
  ('anna@example.com'  , 350000::numeric, 16.5::numeric),
  ('pavel@example.com' , 150000::numeric, 12.5::numeric),
  ('olga@example.com'  , 500000::numeric, 17.9::numeric),
  ('sergey@example.com', 120000::numeric, 13.9::numeric)
) AS v(email, original_amount, interest_rate)
JOIN users u ON u.email = v.email;

-- 4) Платежи по кредитам - исправлено: используем total_amount вместо amount, убраны method и status
INSERT INTO loan_payments (loan_id, user_id, total_amount, paid_at, comment)
SELECT l.loan_id,
       l.user_id,
       round((l.original_amount / 12)::numeric, 2),
       now() - interval '7 days',
       'Demo payment'
FROM loans l
JOIN users u ON u.user_id = l.user_id
WHERE u.email IN ('ivan@example.com','anna@example.com','pavel@example.com','olga@example.com','sergey@example.com');

-- 5) Транзакции (по 3 на пользователя в каждом банке)
WITH ua AS (
  SELECT uba.user_id, uba.bank_id
  FROM user_bank_accounts uba
  JOIN users u ON u.user_id = uba.user_id
  WHERE u.email IN ('ivan@example.com','anna@example.com','pavel@example.com','olga@example.com','sergey@example.com')
),
rows AS (
  SELECT user_id, bank_id, generate_series(1,3) AS n FROM ua
)
INSERT INTO transactions (user_id, bank_id, amount, category, description, occurred_at)
SELECT r.user_id,
       r.bank_id,
       round((random()*-5000 - 1000)::numeric, 2),
       (ARRAY['groceries','transport','entertainment','utilities'])[1 + floor(random()*4)::int],
       'Demo expense',
       now() - ((random()*20)::int || ' days')::interval
FROM rows r;