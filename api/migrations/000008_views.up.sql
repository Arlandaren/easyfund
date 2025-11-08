BEGIN;

-- Общий баланс по всем банкам пользователя
CREATE MATERIALIZED VIEW total_balance_view AS
 SELECT user_id, SUM(balance) as total_balance
 FROM user_bank_accounts
 GROUP BY user_id;

-- Общая задолженность
CREATE MATERIALIZED VIEW total_debt_view AS
 SELECT l.user_id, SUM(ls.remaining_principal) as total_debt
 FROM loans l
 JOIN loan_splits ls ON l.loan_id = ls.loan_id
 WHERE l.status = 'ACTIVE'
 GROUP BY l.user_id;

-- Процент выплаченного по каждому кредиту
CREATE MATERIALIZED VIEW credit_progress_view AS
 SELECT l.loan_id, l.user_id,
        SUM(ls.split_amount - ls.remaining_principal) / SUM(ls.split_amount) * 100 as percent_paid
 FROM loans l
 JOIN loan_splits ls ON l.loan_id = ls.loan_id
 GROUP BY l.loan_id, l.user_id;

COMMIT;
