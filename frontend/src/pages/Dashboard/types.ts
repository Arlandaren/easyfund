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
  isPositive?: boolean;
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
  // Account summary
  accountBalance: number;
  totalDebt: number;
  
  // Credit products stats
  creditCount: number;
  creditCardCount: number;
  
  // Progress
  progress: ProgressData;
  
  // Credit rating
  creditRating: CreditRatingData;
  
  // Payments
  payments: PaymentItem[];
  
  // Transactions
  transactions: Transaction[];
  
  // Debts by bank
  debtsByBank: DebtItem[];
}

// API Response types based on OpenAPI schema
export interface ApiUser {
  user_id: string;
  full_name: string;
  email: string;
  phone: string;
  created_at: string;
  updated_at: string;
}

export interface BalanceSummary {
  user_id: string;
  total_balance: string;
  currency: string;
  by_bank: Array<{
    bank_id: number;
    balance: string;
  }>;
}

export interface ApiTransaction {
  transaction_id: number;
  user_id: string;
  bank_id: number;
  occurred_at: string;
  amount: string;
  category: string;
  description: string;
}

export interface ApiLoan {
  loan_id: string;
  user_id: string;
  amount: string;
  rate: string;
  months: number;
  status: string;
  created_at: string;
}

export interface UserDebt {
  user_id: string;
  total_debt: string;
  by_loan: Array<{
    loan_id: string;
    outstanding: string;
  }>;
}

export interface ApiApplication {
  application_id: string;
  user_id: string;
  amount: string;
  term_months: number;
  purpose: string;
  status: string;
  created_at: string;
}

// Export User type for use in other components
export type { User } from '../../context/AuthContext';