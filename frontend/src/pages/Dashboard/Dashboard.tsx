import React, { useEffect, useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import {
  AccountSummarySection,
  CreditScoreSection,
  DebtOverviewSection,
  FinancialGoalsSection,
  PaymentHistorySection,
  ProgressSection,
  CreditRatingSection,
} from './components';
import {
  DashboardData,
  BalanceSummary,
  UserDebt,
  ApiLoan,
  ApiTransaction,
  ApiApplication,
} from './types';
import { dashboardAPI } from '../../utils/api';
import { TopBar } from '../../components';
import './Dashboard.css';

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [loading, setLoading] = useState(true);
  const [dashboardData, setDashboardData] = useState<DashboardData | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Fallback data
  const defaultData: DashboardData = {
    accountBalance: 0,
    totalDebt: 0,
    creditCount: 0,
    creditCardCount: 0,
    progress: {
      currentDebt: 0,
      initialDebt: 0,
      targetDebt: 0,
      percentage: 0,
    },
    creditRating: {
      score: 645,
      min: 300,
      max: 850,
      labels: ['Низкий', 'Неплохой', 'Хороший', 'Отличный'],
    },
    payments: [],
    transactions: [],
    debtsByBank: [],
  };

  useEffect(() => {
    fetchDashboardData();
  }, [user?.user_id]);

  // Close dropdown when clicking outside
  // Helper functions
  const safeToString = (value: any): string => {
    if (value === null || value === undefined) return '';
    if (typeof value === 'string') return value;
    if (typeof value === 'number') return value.toString();
    return String(value);
  };

  const safeParseFloat = (value: any): number => {
    if (value === null || value === undefined) return 0;
    const num = parseFloat(safeToString(value));
    return isNaN(num) ? 0 : num;
  };

  const normalizeArray = <T,>(x: T[] | null | undefined): T[] => (Array.isArray(x) ? x : []);

  const fetchDashboardData = async () => {
    if (!user?.user_id) {
      console.log('No user ID available, using mock data');
      setTimeout(() => {
        setDashboardData(defaultData);
        setLoading(false);
      }, 300);
      return;
    }

    try {
      setLoading(true);
      setError(null);
      console.log('🔄 Fetching dashboard data for user:', user.user_id);

      // Получаем все данные без ограничения по транзакциям
      const [balanceResponse, debtResponse, loansResponse, transactionsResponse, applicationsResponse] =
        await dashboardAPI.getFullDashboardData(user.user_id);

      // Normalize responses
      const balanceData: BalanceSummary = balanceResponse?.data || {
        user_id: user.user_id,
        total_balance: '0',
        currency: 'RUB',
        by_bank: [],
      };

      const debtData: UserDebt = debtResponse?.data || {
        user_id: user.user_id,
        total_debt: '0',
        by_loan: [],
      };

      const loansData: ApiLoan[] = normalizeArray<ApiLoan>(loansResponse?.data);

      // Правильное извлечение массива транзакций из объекта { transactions, total_spent }
      const txArraySource = transactionsResponse?.data?.transactions;
      const txArray: ApiTransaction[] = Array.isArray(txArraySource) ? txArraySource : [];

      const applicationsData: ApiApplication[] = normalizeArray<ApiApplication>(
        applicationsResponse?.data
      );

      console.log('✅ API raw transactions length:', Array.isArray(txArraySource) ? txArraySource.length : 0);

      console.log('✅ API data received:', {
        balance: balanceData,
        debt: debtData,
        loans: loansData.length,
        transactions: txArray.length,
        applications: applicationsData.length,
        transactionsSample: txArray.slice(0, 3) // Показываем первые 3 транзакции для отладки
      });

      // Transform API data to frontend format
      let totalDebtAmount = safeParseFloat(debtData.total_debt);
      if (totalDebtAmount === 0 && loansData.length > 0) {
        totalDebtAmount = loansData.reduce((sum, loan) => {
          const loanData = loan as any;
          const remaining = safeParseFloat(
            loanData?.remaining_balance ?? loanData?.remaining_amount ?? loanData?.outstanding
          );
          const fallbackAmount = safeParseFloat(
            loan.amount ?? loanData?.original_amount ?? loanData?.requested_amount ?? loanData?.principal
          );
          return sum + (remaining > 0 ? remaining : fallbackAmount);
        }, 0);
      }

      const transformedData: DashboardData = {
        accountBalance: safeParseFloat(balanceData.total_balance),
        totalDebt: totalDebtAmount,
        creditCount: loansData.length,
        creditCardCount: applicationsData.filter((app) => app.status === 'active').length,
        progress: {
          currentDebt: totalDebtAmount * 0.6,
          initialDebt: totalDebtAmount,
          targetDebt: 0,
          percentage: totalDebtAmount > 0 ? 60 : 0,
        },
        creditRating: {
          score: 645,
          min: 300,
          max: 850,
          labels: ['Низкий', 'Неплохой', 'Хороший', 'Отличный'],
        },
        payments: loansData.map((loan, index) => ({
          id: index + 1,
          title: `Кредит #${loan.loan_id ?? index + 1}`,
          dueDate: 'Ближайший платеж скоро',
          amount: safeToString(loan.amount || '0'),
        })),
        // ✅ ВСЕ транзакции пользователя
        transactions: txArray.map((transaction) => {
          const amountStr = safeToString(transaction.amount);
          const amountNum = safeParseFloat(amountStr);
          return {
            id: transaction.transaction_id,
            title: transaction.description || 'Транзакция',
            amount: amountStr,
            isPositive: amountNum > 0,
            company: transaction.category || 'Unknown',
            occurredAt: transaction.occurred_at,
            bankId: transaction.bank_id,
          };
        }),
        debtsByBank: [
          { id: 1, bankName: 'ВТБ', amount: 213123, color: '#5218f4' },
          { id: 2, bankName: 'Сбербанк', amount: 650000, color: '#d081e4' },
          { id: 3, bankName: 'Альфа-Банк', amount: 180000, color: '#189CF4' },
        ],
      };

      console.log('📊 Transformed data:', {
        transactionsCount: transformedData.transactions.length,
        paymentsCount: transformedData.payments.length,
        firstTransaction: transformedData.transactions[0],
        lastTransaction: transformedData.transactions[transformedData.transactions.length - 1]
      });

      setDashboardData(transformedData);
      setLoading(false);
    } catch (err) {
      console.error('❌ Error fetching dashboard data:', err);
      setError('Не удалось загрузить данные с сервера');
      setTimeout(() => {
        setDashboardData(defaultData);
        setLoading(false);
      }, 300);
    }
  };

  const userName = useMemo(() => {
    return user?.full_name || user?.email?.split('@')[0] || 'Пользователь';
  }, [user]);

  if (error) {
    return (
      <div className="dashboard">
        <div className="dashboard__error">
          <h2>Ошибка загрузки</h2>
          <p>{error}</p>
          <button onClick={fetchDashboardData} className="dashboard__retry-btn">
            Попробовать снова
          </button>
        </div>
      </div>
    );
  }

  if (loading || !dashboardData) {
    return (
      <div className="dashboard dashboard--loading">
        <div className="dashboard__loading-spinner">Загрузка данных...</div>
      </div>
    );
  }

  return (
    <div className="dashboard">
      <div className="dashboard__container">
        {/* Background */}
        <div className="dashboard__background" />

        <TopBar variant="overlay" />

        {/* Greeting */}
        <h1 className="dashboard__greeting">С возвращением, {userName}</h1>

        {/* Main Sections */}
        <div className="dashboard__sections">
          <CreditScoreSection
            accountBalance={dashboardData.accountBalance}
            onTransfer={() => console.log('Transfer clicked')}
            onTopUp={() => console.log('Top up clicked')}
          />

          <PaymentHistorySection
            totalDebt={dashboardData.totalDebt}
            creditCount={dashboardData.creditCount}
            creditCardCount={dashboardData.creditCardCount}
            onViewAllProducts={() => navigate('/applications')}
          />

          {/* История трат и Вы почти у цели в одной строке */}
          <div className="dashboard__row">
            <DebtOverviewSection
              transactions={dashboardData.transactions}
              onFilterChange={(filter: string) => console.log('Filter changed:', filter)}
            />

            <ProgressSection progress={dashboardData.progress} totalDebt={dashboardData.totalDebt} />
          </div>

          <AccountSummarySection
            payments={dashboardData.payments}
            onViewAll={() => console.log('View all payments')}
          />

          <FinancialGoalsSection debtsByBank={dashboardData.debtsByBank} />

          <CreditRatingSection creditRating={dashboardData.creditRating} />
        </div>

      </div>
    </div>
  );
};