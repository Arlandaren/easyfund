// Types for Dashboard data

export interface PaymentItem {
  id: string;
  icon?: string;
  title: string;
  dueDate: string;
  amount: string;
  bankName?: string;
}

export interface Transaction {
  id: string;
  image?: string;
  company?: string;
  title: string;
  amount: string;
  isPositive?: boolean;
}

export interface DebtItem {
  id: string;
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

