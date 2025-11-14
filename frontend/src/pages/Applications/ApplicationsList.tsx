import React, { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Layout, Card, Button, Input } from '../../components';
import { useAuth } from '../../context/AuthContext';
import { loansAPI, applicationsAPI } from '../../utils/api';
import './ApplicationsList.css';

interface CreditApplication {
  id: number;
  name: string;
  amount: string;
  interest_rate: string;
  months?: number;
  status?: string;
  purpose?: string;
  created_at?: string;
  submitted_at?: string;
  updated_at?: string;
  bank_id?: number;
  term_months?: number;
}

export const ApplicationsList: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [applications, setApplications] = useState<CreditApplication[]>([]);
  const [filteredApplications, setFilteredApplications] = useState<CreditApplication[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState('');
  const [statusMenuOpen, setStatusMenuOpen] = useState(false);
  const statusDropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (user?.user_id) {
      fetchApplications();
    }
  }, [user?.user_id]);

  useEffect(() => {
    filterApplications();
  }, [applications, searchTerm, filterStatus]);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        statusDropdownRef.current &&
        !statusDropdownRef.current.contains(event.target as Node)
      ) {
        setStatusMenuOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const statusOptions = [
    { value: '', label: 'Выберите статус', disabled: true },
    { value: 'all', label: 'Все статусы' },
    { value: 'active', label: 'Активный' },
    { value: 'pending', label: 'На рассмотрении' },
    { value: 'approved', label: 'Одобрено' },
    { value: 'rejected', label: 'Отклонено' },
    { value: 'completed', label: 'Завершено' },
  ];

  const fetchApplications = async () => {
    if (!user?.user_id) return;

    try {
      setLoading(true);
      setError(null);

      // Получаем оба типа предложений: пользовательские приложения и кредиты
      const [loansResponse, applicationsResponse] = await Promise.all([
        loansAPI.getUserLoans(String(user.user_id)).catch(() => ({ data: [] })),
        applicationsAPI.getUserApplications(String(user.user_id)).catch(() => ({ data: [] })),
      ]);

      const loansData = Array.isArray(loansResponse?.data) ? loansResponse.data : [];
      const appsData = Array.isArray(applicationsResponse?.data) ? applicationsResponse.data : [];

      // Преобразуем кредиты в единый формат
      const formattedLoans: CreditApplication[] = loansData.map((loan: any) => ({
        id: loan.loan_id ?? loan.id ?? 0,
        name: `Кредит #${loan.loan_id ?? loan.id ?? 0}`,
        amount: loan.original_amount || loan.amount || '0',
        interest_rate: loan.interest_rate || loan.rate || '0',
        months: loan.months || 0,
        status: loan.status || 'active',
        purpose: loan.purpose || 'Кредит',
        created_at: loan.taken_at || loan.created_at,
      }));

      // Преобразуем приложения
      const formattedApplications: CreditApplication[] = appsData.map((app: any) => ({
        id: app.application_id ?? app.id ?? 0,
        name: `Приложение #${app.application_id ?? app.id ?? 0}`,
        amount: app.requested_amount || app.amount || '0',
        interest_rate: '0',
        status: app.status || app.status_code || 'pending',
        purpose: app.purpose || 'Заявка на кредит',
        term_months: app.term_months,
        bank_id: app.bank_id,
        submitted_at: app.submitted_at,
        updated_at: app.updated_at,
      }));

      // Объединяем и сортируем по дате
      const combined = [...formattedLoans, ...formattedApplications];
      setApplications(combined);
    } catch (err) {
      console.error('❌ Error fetching applications:', err);
      setError('Не удалось загрузить кредитные предложения');
    } finally {
      setLoading(false);
    }
  };

  const filterApplications = () => {
    let filtered = applications;

    if (searchTerm) {
      filtered = filtered.filter((app) =>
        app.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        (app.purpose ?? '').toLowerCase().includes(searchTerm.toLowerCase()) ||
        app.amount.includes(searchTerm)
      );
    }

    if (filterStatus && filterStatus !== 'all') {
      filtered = filtered.filter((app) => app.status === filterStatus);
    }

    setFilteredApplications(filtered);
  };

  const getStatusBadge = (status?: string) => {
    if (!status) return { label: 'Неизвестно', color: '#9ca3af' };

    const statusMap: Record<string, { label: string; color: string }> = {
      active: { label: 'Активный', color: '#10b981' },
      pending: { label: 'На рассмотрении', color: '#f59e0b' },
      approved: { label: 'Одобрено', color: '#3b82f6' },
      rejected: { label: 'Отклонено', color: '#ef4444' },
      completed: { label: 'Завершено', color: '#6366f1' },
      under_review: { label: 'На проверке', color: '#f59e0b' },
      partially_approved: { label: 'Частично одобрено', color: '#8b5cf6' },
    };

    return statusMap[status] || { label: status.toUpperCase(), color: '#9ca3af' };
  };

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return '';
    try {
      return new Date(dateStr).toLocaleDateString('ru-RU');
    } catch {
      return dateStr;
    }
  };

  if (loading) {
    return (
      <Layout>
        <div className="applications-container loading">
          <div className="applications-spinner">Загрузка кредитных предложений...</div>
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout>
        <div className="applications-container error">
          <div className="applications-error">{error}</div>
          <Button onClick={() => fetchApplications()}>Попробовать снова</Button>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="applications-container">
        <div className="applications-header">
          <h1 className="applications-title">Кредитные предложения</h1>
          <Button onClick={() => navigate('/applications/new')} className="applications-btn-new">
            + Новая заявка
          </Button>
        </div>

        <div className="applications-filters">
          <Input
            type="text"
            placeholder="Поиск по названию или целям..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="applications-search"
          />

          <div
            className={`applications-select ${filterStatus ? '' : 'applications-select--placeholder'}`}
            onClick={() => setStatusMenuOpen((prev) => !prev)}
            ref={statusDropdownRef}
            role="button"
            tabIndex={0}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                setStatusMenuOpen((prev) => !prev);
              }
            }}
            aria-haspopup="listbox"
            aria-expanded={statusMenuOpen}
          >
            <span>
              {statusOptions.find((opt) => opt.value === filterStatus)?.label || 'Выберите статус'}
            </span>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path
                d="M6 9L12 15L18 9"
                stroke="#082131"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
            {statusMenuOpen && (
              <ul className="applications-select__menu" role="listbox">
                {statusOptions.map((option) => (
                  <li
                    key={option.value}
                    className={`applications-select__option ${
                      filterStatus === option.value ? 'applications-select__option--selected' : ''
                    } ${option.disabled ? 'applications-select__option--disabled' : ''}`}
                    onClick={(e) => {
                      e.stopPropagation();
                      if (option.disabled) return;
                      setFilterStatus(option.value);
                      setStatusMenuOpen(false);
                    }}
                    role="option"
                    aria-selected={filterStatus === option.value}
                  >
                    {option.label}
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>

        {filteredApplications.length === 0 ? (
          <Card className="applications-empty">
            <p>Кредитные предложения не найдены</p>
          </Card>
        ) : (
          <div className="applications-grid">
            {filteredApplications.map((app) => {
              const statusBadge = getStatusBadge(app.status);
              return (
                <Card
                  key={app.id}
                  onClick={() => navigate(`/applications/${app.id}`)}
                  className="applications-card"
                >
                  <div className="applications-card-header">
                    <h3 className="applications-card-title">{app.name}</h3>
                    <span
                      className="applications-badge"
                      style={{ backgroundColor: statusBadge.color }}
                    >
                      {statusBadge.label}
                    </span>
                  </div>

                  <div className="applications-card-content">
                    <div className="applications-field">
                      <span className="applications-label">Сумма:</span>
                      <span className="applications-value">{app.amount} ₽</span>
                    </div>

                    {app.interest_rate && app.interest_rate !== '0' && (
                      <div className="applications-field">
                        <span className="applications-label">Процентная ставка:</span>
                        <span className="applications-value">{app.interest_rate}%</span>
                      </div>
                    )}

                    {app.months && app.months > 0 && (
                      <div className="applications-field">
                        <span className="applications-label">Срок:</span>
                        <span className="applications-value">{app.months} месяцев</span>
                      </div>
                    )}

                    {app.purpose && (
                      <div className="applications-field">
                        <span className="applications-label">Назначение:</span>
                        <span className="applications-value">{app.purpose}</span>
                      </div>
                    )}

                    {app.created_at && (
                      <div className="applications-field">
                        <span className="applications-label">Создано:</span>
                        <span className="applications-value">{formatDate(app.created_at)}</span>
                      </div>
                    )}
                  </div>

                  <div className="applications-card-footer">
                    <Button className="applications-card-btn">Просмотреть</Button>
                  </div>
                </Card>
              );
            })}
          </div>
        )}
      </div>
    </Layout>
  );
};
