import React, { useEffect, useState, useMemo, useRef } from 'react';
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
import easyfundLogoSvg from '../../utils/img/easyfund-logo.svg';
import profileImage from '../../utils/img/profile.png';
import './Dashboard.css';

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [loading, setLoading] = useState(true);
  const [dashboardData, setDashboardData] = useState<DashboardData | null>(null);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
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
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setDropdownOpen(false);
      }
    };

    if (dropdownOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [dropdownOpen]);

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

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
      const totalDebtAmount = safeParseFloat(debtData.total_debt);

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

        {/* Header */}
        <header className="dashboard__header">
          <button
            className="dashboard__logo-link"
            onClick={() => navigate('/dashboard')}
            type="button"
            aria-label="Go to dashboard"
          >
            <img className="dashboard__logo" alt="EasyFund Logo" src={easyfundLogoSvg} />
          </button>
          <div className="dashboard__header-actions">
            <button className="dashboard__header-icon" type="button" aria-label="Search">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path
                  d="M21 21L15 15M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </button>
            <button className="dashboard__header-icon" type="button" aria-label="Notifications">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path
                  d="M18 8C18 6.4087 17.3679 4.88258 16.2426 3.75736C15.1174 2.63214 13.5913 2 12 2C10.4087 2 8.88258 2.63214 7.75736 3.75736C6.63214 4.88258 6 6.4087 6 8C6 15 3 17 3 17H21C21 17 18 15 18 8Z"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
                <path
                  d="M13.73 21C13.5542 21.3031 13.3019 21.5547 12.9982 21.7295C12.6946 21.9044 12.3504 21.9965 12 21.9965C11.6496 21.9965 11.3054 21.9044 11.0018 21.7295C10.6982 21.5547 10.4458 21.3031 10.27 21"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </button>
            <div className="dashboard__avatar" ref={dropdownRef}>
              <button
                className="dashboard__avatar-button"
                onClick={() => setDropdownOpen(!dropdownOpen)}
                type="button"
                aria-label="User menu"
              >
                <img
                  className="dashboard__avatar-image"
                  src={profileImage}
                  alt={user?.full_name || 'User'}
                />
              </button>
              {dropdownOpen && (
                <div className="dashboard__dropdown">
                  <button
                    className="dashboard__dropdown-item"
                    onClick={() => {
                      setDropdownOpen(false);
                    }}
                    type="button"
                  >
                    Перейти в профиль
                  </button>
                  <button
                    className="dashboard__dropdown-item"
                    onClick={() => {
                      setDropdownOpen(false);
                      handleLogout();
                    }}
                    type="button"
                  >
                    Выйти
                  </button>
                </div>
              )}
            </div>
          </div>
        </header>

        {/* Greeting */}
        <h1 className="dashboard__greeting">Добрый день, {userName}!</h1>

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

            <ProgressSection progress={dashboardData.progress} />
          </div>

          <AccountSummarySection
            payments={dashboardData.payments}
            onViewAll={() => console.log('View all payments')}
          />

          <FinancialGoalsSection debtsByBank={dashboardData.debtsByBank} />

          <CreditRatingSection creditRating={dashboardData.creditRating} />
        </div>

        {/* Bottom Navigation */}
        <nav className="dashboard__nav" aria-label="Main navigation">
          <div className="dashboard__nav-indicator" />
          <button
            className="dashboard__nav-btn dashboard__nav-btn--active"
            aria-label="Home"
            onClick={() => navigate('/dashboard')}
          >
            <svg width="30" height="30" viewBox="0 0 24 24" fill="none">
              <path
                d="M3 12L5 10M5 10L12 3L19 10M5 10V20C5 20.5523 5.44772 21 6 21H9M19 10L21 12M19 10V20C19 20.5523 18.5523 21 18 21H15M9 21C9.55228 21 10 20.5523 10 20V16C10 15.4477 10.4477 15 11 15H13C13.5523 15 14 15.4477 14 16V20C14 20.5523 14.4477 21 15 21M9 21H15"
                stroke="#FFFFFF"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>
          <button
            className="dashboard__nav-btn"
            aria-label="Applications"
            onClick={() => navigate('/applications')}
          >
            <svg width="30" height="30" viewBox="0 0 24 24" fill="none">
              <path
                d="M9 12H15M9 16H15M17 21H7C5.89543 21 5 20.1046 5 19V5C5 3.89543 5.89543 3 7 3H12.5858C12.851 3 13.1054 3.10536 13.2929 3.29289L18.7071 8.70711C18.8946 8.89464 19 9.149 19 9.41421V19C19 20.1046 18.1046 21 17 21Z"
                stroke="#082131"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>
        </nav>
      </div>
    </div>
  );
};