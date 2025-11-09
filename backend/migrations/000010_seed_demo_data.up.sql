-- 000010_seed_demo_data_up.sql

-- 1) Пользователи
INSERT INTO users (user_id, full_name, email, phone, password_hash)
VALUES
  (1, 'Ivan Petrov','ivan@example.com','+79998887761','$2y$10$Kjh5eRBHUIhi38Dc/9za9ORrAwDzn8qE0.ZgR1W5NsJDndjzk/mJ2'),
  (2, 'Anna Sidorova','anna@example.com','+79998887762','$2y$10$Kjh5eRBHUIhi38Dc/9za9ORrAwDzn8qE0.ZgR1W5NsJDndjzk/mJ2'),
  (3, 'Pavel Smirnov','pavel@example.com','+79998887763','$2y$10$Kjh5eRBHUIhi38Dc/9za9ORrAwDzn8qE0.ZgR1W5NsJDndjzk/mJ2'),
  (4, 'Olga Ivanova','olga@example.com','+79998887764','$2y$10$Kjh5eRBHUIhi38Dc/9za9ORrAwDzn8qE0.ZgR1W5NsJDndjzk/mJ2'),
  (5, 'Sergey Volkov','sergey@example.com','+79998887765','$2y$10$Kjh5eRBHUIhi38Dc/9za9ORrAwDzn8qE0.ZgR1W5NsJDndjzk/mJ2');

-- Устанавливаем последовательность для users
SELECT setval('users_user_id_seq', (SELECT MAX(user_id) FROM users));

-- 2) Счета (user_bank_accounts)
INSERT INTO user_bank_accounts (user_id, bank_id, balance, currency, created_at)
SELECT u.user_id, b.bank_id,
       round((random()*90000 + 10000)::numeric, 2),
       'RUB',
       now()
FROM users u
CROSS JOIN banks b
WHERE u.user_id BETWEEN 1 AND 5;

-- 3) Кредиты (по одному на пользователя)
INSERT INTO loans (user_id, original_amount, interest_rate, status, purpose, created_at)
SELECT u.user_id, v.original_amount, v.interest_rate, 'ACTIVE', 'Demo loan purpose', now()
FROM (VALUES
  (1, 200000::numeric, 14.9::numeric),
  (2, 350000::numeric, 16.5::numeric),
  (3, 150000::numeric, 12.5::numeric),
  (4, 500000::numeric, 17.9::numeric),
  (5, 120000::numeric, 13.9::numeric)
) AS v(user_id, original_amount, interest_rate)
JOIN users u ON u.user_id = v.user_id;

-- 4) Платежи по кредитам
INSERT INTO loan_payments (loan_id, amount, paid_at, method, status)
SELECT l.loan_id,
       round((l.original_amount / 12)::numeric, 2) AS amount,
       now() - interval '7 days'                   AS paid_at,
       'card'                                      AS method,
       'posted'                                    AS status
FROM loans l
WHERE l.user_id BETWEEN 1 AND 5;

-- 5) Транзакции (по 3 на пользователя в каждом банке)
WITH user_banks AS (
  SELECT u.user_id, b.bank_id
  FROM users u
  CROSS JOIN banks b
  WHERE u.user_id BETWEEN 1 AND 5
),
transaction_rows AS (
  SELECT user_id, bank_id, generate_series(1,3) AS transaction_num
  FROM user_banks
)
INSERT INTO transactions (user_id, bank_id, amount, category, description, occurred_at)
SELECT 
  tr.user_id,
  tr.bank_id,
  round((random()*-5000 - 1000)::numeric, 2) as amount,
  (ARRAY['groceries','transport','entertainment','utilities'])[1 + floor(random()*4)::int] as category,
  'Demo expense ' || tr.transaction_num as description,
  now() - ((random()*20)::int || ' days')::interval as occurred_at
FROM transaction_rows tr;