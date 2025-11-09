BEGIN;

CREATE TABLE banks (
  bank_id smallint PRIMARY KEY,
  code text NOT NULL UNIQUE,
  name text NOT NULL
);

INSERT INTO banks (bank_id, code, name) VALUES
  (1, 'ALFA', 'Альфа'),
  (2, 'VTB', 'ВТБ'),
  (3, 'SBER', 'Сбер'),
  (4, 'TBANK', 'ТБанк'),
  (5, 'OPT', 'ОПТ-Банк');

CREATE TABLE credit_application_statuses (
  status_code text PRIMARY KEY,
  display_name text NOT NULL
);

INSERT INTO credit_application_statuses(status_code, display_name) VALUES
  ('PENDING','На рассмотрении'),
  ('APPROVED','Одобрено'),
  ('REJECTED','Отклонено'),
  ('CANCELLED','Отменено');

CREATE TABLE credit_application_types (
  type_code text PRIMARY KEY,
  display_name text NOT NULL
);

INSERT INTO credit_application_types(type_code, display_name) VALUES
  ('PERSONAL','Потребительский'),
  ('AUTO','Авто'),
  ('MORTGAGE','Ипотека'),
  ('OTHER','Другое');

COMMIT;
