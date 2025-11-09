// Types for Dashboard data

export interface PaymentItem {
  id: number;
  icon?: string;
  title: string;
  dueDate: string;
  amount: string;
  bankName?: string;
}

export interface Transaction {
  id: number;
  image?: string;
  company?: string;
  title: string;
  amount: string;
  isPositive: boolean;
  occurredAt?: string;
  bankId?: number;
}

export interface DebtItem {
  id: number;
  bankName: string;
  amount: number;
  color: string;
}

export interface ProgressData {
  currentDebt: number;
  initialDebt: number;
  targetDebt: number;
  percentage: number;
}

export interface CreditRatingData {
  score: number;
  min: number;
  max: number;
  labels: string[];
}

export interface DashboardData {
  accountBalance: number;
  totalDebt: number;
  creditCount: number;
  creditCardCount: number;
  progress: ProgressData;
  creditRating: CreditRatingData;
  payments: PaymentItem[];
  transactions: Transaction[];
  debtsByBank: DebtItem[];
}

// API Response types (согласно документации)
export interface ApiUser {
  user_id: number;
  full_name: string;
  email: string;
  phone: string;
  created_at: string;
  updated_at: string;
}

export interface BalanceSummary {
  user_id: number;
  total_balance: string;
  currency: string;
  by_bank: Array<{
    bank_id: number;
    balance: string;
  }>;
}

// ✅ Правильная структура транзакции согласно dock.yaml
export interface ApiTransaction {
  transaction_id: number;
  user_id: number;
  bank_id: number;
  occurred_at: string;
  amount: string; // может быть положительное или отрицательное
  category: string;
  description: string;
}

export interface ApiLoan {
  loan_id: number;
  user_id: number;
  amount?: string;
  rate?: string;
  months?: number;
  status?: string;
  created_at?: string;
}

export interface UserDebt {
  user_id: number;
  total_debt: string;
  by_loan?: Array<{
    loan_id: number;
    outstanding?: string;
  }>;
}

export interface ApiApplication {
  application_id: number;
  user_id: number;
  amount?: string;
  term_months?: number;
  purpose?: string;
  status?: string;
  created_at?: string;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
  status?: number;
}

export interface TransactionFilters {
  bank_id?: number;
  from?: string;
  to?: string;
  limit?: number;
  offset?: number;
  category?: string;
}