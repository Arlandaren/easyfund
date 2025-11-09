import React, { useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Layout, Card, Button, Input } from '../../../components';
import { useAuth } from '../../../context/AuthContext';
import { applicationsAPI, banksAPI } from '../../../utils/api';
import './NewApplication.css';

import alfaLogo from '../../../utils/img/alfa.png';
import vtbLogo from '../../../utils/img/vtb.png';
import sberLogo from '../../../utils/img/sber.png';
import tbankLogo from '../../../utils/img/tbank.png';
import otpLogo from '../../../utils/img/otp.png';

interface Bank {
  bank_id: number;
  code: string;
  name: string;
  logo?: string;
}

interface ApplicationHistoryItem {
  id: number;
  bankId: number | null;
  typeCode: string;
  requestedAmount: number;
  status: string;
  submittedAt?: string;
}

interface CreditTemplate {
  code: string;
  title: string;
  description: string;
  highlight?: string;
}

const CREDIT_TEMPLATES: CreditTemplate[] = [
  {
    code: 'PERSONAL',
    title: 'Потребительский кредит',
    description: 'Ускорьте оформление личного кредита и получите решение в течение одного рабочего дня.',
    highlight: 'Популярно',
  },
  {
    code: 'MORTGAGE',
    title: 'Ипотечное кредитование',
    description: 'Подберите ипотечный продукт и передайте заявку сразу в несколько банков.',
  },
  {
    code: 'AUTO',
    title: 'Автокредит',
    description: 'Получите финансирование на покупку автомобиля с гибкими условиями погашения.',
  },
  {
    code: 'OTHER',
    title: 'Индивидуальная заявка',
    description: 'Сформируйте заявку с пользовательскими параметрами для редких сценариев.',
  },
];

const FALLBACK_BANKS: Bank[] = [
  { bank_id: 1, code: 'ALFA', name: 'Альфа-Банк', logo: alfaLogo },
  { bank_id: 2, code: 'VTB', name: 'ВТБ', logo: vtbLogo },
  { bank_id: 3, code: 'SBER', name: 'Сбербанк', logo: sberLogo },
  { bank_id: 4, code: 'TBANK', name: 'Т-Банк', logo: tbankLogo },
  { bank_id: 5, code: 'OPT', name: 'ОТП Банк', logo: otpLogo },
];

const BANK_META: Record<string, { rate: string; label?: string }> = {
  ALFA: { rate: '25% годовых' },
  VTB: { rate: '13% годовых', label: '20 дней без %' },
  SBER: { rate: '18% годовых', label: '10 дней без %' },
  TBANK: { rate: '10% годовых', label: 'Мгновенное решение' },
  OPT: { rate: '21% годовых' },
};

const FALLBACK_HISTORY: ApplicationHistoryItem[] = [
  {
    id: 101,
    bankId: 2,
    typeCode: 'PERSONAL',
    requestedAmount: 180_000,
    status: 'PENDING',
    submittedAt: new Date().toISOString(),
  },
  {
    id: 102,
    bankId: 3,
    typeCode: 'MORTGAGE',
    requestedAmount: 4_500_000,
    status: 'APPROVED',
    submittedAt: new Date(Date.now() - 86400000 * 5).toISOString(),
  },
  {
    id: 103,
    bankId: 1,
    typeCode: 'AUTO',
    requestedAmount: 950_000,
    status: 'REJECTED',
    submittedAt: new Date(Date.now() - 86400000 * 12).toISOString(),
  },
];

const STATUS_MAP: Record<
  string,
  {
    label: string;
    tone: 'neutral' | 'info' | 'success' | 'warning' | 'danger';
  }
> = {
  PENDING: { label: 'На рассмотрении', tone: 'warning' },
  APPROVED: { label: 'Одобрено', tone: 'success' },
  REJECTED: { label: 'Отклонено', tone: 'danger' },
  CANCELLED: { label: 'Отменено', tone: 'neutral' },
  ACTIVE: { label: 'Активно', tone: 'info' },
};

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(amount);

const sanitizeAmount = (value: string | number): string => {
  if (typeof value === 'number') {
    return Math.max(value, 0).toFixed(0);
  }
  return value.replace(/\s+/g, '').replace(/,/g, '.');
};

export const NewApplication: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuth();

  const [banks, setBanks] = useState<Bank[]>([]);
  const [history, setHistory] = useState<ApplicationHistoryItem[]>([]);
  const [selectedTemplate, setSelectedTemplate] = useState<string>(CREDIT_TEMPLATES[0].code);
  const [selectedBankId, setSelectedBankId] = useState<number | null>(null);
  const [loanAmount, setLoanAmount] = useState<number>(150_000);
  const [bankAmounts, setBankAmounts] = useState<Record<number, string>>({});

  const [initialLoading, setInitialLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadData = async () => {
      setInitialLoading(true);
      let banksData: Bank[] = FALLBACK_BANKS;
      let historyData: ApplicationHistoryItem[] = FALLBACK_HISTORY;

      try {
        const [banksResponse, historyResponse] = await Promise.allSettled([
          banksAPI.getAll(),
          user?.user_id ? applicationsAPI.getUserApplications(String(user.user_id)) : Promise.resolve({ data: [] }),
        ]);

        if (banksResponse.status === 'fulfilled') {
          const data = banksResponse.value?.data;
          if (Array.isArray(data) && data.length > 0) {
            banksData = data.map((bank: Bank) => {
              const code = (bank.code ?? '').toString().toUpperCase();
              const fallback = FALLBACK_BANKS.find(
                (item) => item.bank_id === bank.bank_id || item.code === code || item.code === bank.code,
              );
              return {
                ...bank,
                code: code || fallback?.code || bank.code,
                logo: fallback?.logo,
              };
            });
          }
        }

        if (historyResponse.status === 'fulfilled') {
          const data = historyResponse.value?.data;
          if (Array.isArray(data)) {
            historyData = data
              .map((item: any, index: number) => ({
                id: item.application_id ?? item.id ?? index,
                bankId: item.bank_id ?? null,
                typeCode: (item.type_code ?? item.type ?? 'PERSONAL').toString().toUpperCase(),
                requestedAmount: Number(sanitizeAmount(item.requested_amount ?? item.amount ?? '0')) || 0,
                status: (item.status_code ?? item.status ?? 'PENDING').toString().toUpperCase(),
                submittedAt: item.submitted_at ?? item.created_at,
              }))
              .sort((a, b) => {
                const dateA = a.submittedAt ? new Date(a.submittedAt).getTime() : 0;
                const dateB = b.submittedAt ? new Date(b.submittedAt).getTime() : 0;
                return dateB - dateA;
              })
              .slice(0, 8);
          }
        }
      } catch (fetchError) {
        console.warn('Failed to load application dependencies', fetchError);
      }

      if (!isMounted) return;

      setBanks(banksData);
      setSelectedBankId((prev) => (prev === null && banksData.length ? banksData[0].bank_id : prev));
      setHistory(historyData);
      setInitialLoading(false);
    };

    loadData();

    return () => {
      isMounted = false;
    };
  }, [user?.user_id]);

  const decoratedHistory = useMemo(() => {
    return history.map((item) => {
      const template = CREDIT_TEMPLATES.find((t) => t.code === item.typeCode);
      const bank = banks.find((b) => b.bank_id === item.bankId);
      const statusMeta = STATUS_MAP[item.status] ?? STATUS_MAP.PENDING;

      return {
        ...item,
        title: template?.title ?? 'Заявка на кредит',
        bankName: bank?.name ?? 'Банк не выбран',
        statusLabel: statusMeta.label,
        statusTone: statusMeta.tone,
        formattedAmount: formatCurrency(item.requestedAmount || 0),
        formattedDate: item.submittedAt ? new Date(item.submittedAt).toLocaleDateString('ru-RU') : '—',
        logo: bank?.logo,
      };
    });
  }, [history, banks]);

  const handleLoanAmountChange = (value: number) => {
    const normalized = Math.min(Math.max(value, 0), 1_000_000);
    setLoanAmount(normalized);
  };

  const handleSliderChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    handleLoanAmountChange(parseInt(event.target.value, 10) || 0);
  };

  const handleAmountInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value.replace(/\s+/g, '');
    handleLoanAmountChange(Number(value) || 0);
  };

  const handleBankAmountChange = (bankId: number, value: string) => {
    setBankAmounts((prev) => ({
      ...prev,
      [bankId]: value,
    }));
  };

  const handleSelectBank = (bankId: number) => {
    setSelectedBankId(bankId);
    setSuccess(null);
    setError(null);
    setBankAmounts((prev) => {
      if (prev[bankId]) {
        return prev;
      }
      return {
        ...prev,
        [bankId]: sanitizeAmount(loanAmount),
      };
    });
  };

  const refreshHistory = async () => {
    if (!user?.user_id) return;
    try {
      const response = await applicationsAPI.getUserApplications(String(user.user_id));
      const data = Array.isArray(response?.data) ? response.data : [];
      setHistory(
        data
          .map((item: any, index: number) => ({
            id: item.application_id ?? item.id ?? index,
            bankId: item.bank_id ?? null,
            typeCode: (item.type_code ?? item.type ?? 'PERSONAL').toString().toUpperCase(),
            requestedAmount: Number(sanitizeAmount(item.requested_amount ?? item.amount ?? '0')) || 0,
            status: (item.status_code ?? item.status ?? 'PENDING').toString().toUpperCase(),
            submittedAt: item.submitted_at ?? item.created_at,
          }))
          .sort((a, b) => {
            const dateA = a.submittedAt ? new Date(a.submittedAt).getTime() : 0;
            const dateB = b.submittedAt ? new Date(b.submittedAt).getTime() : 0;
            return dateB - dateA;
          })
          .slice(0, 8),
      );
    } catch (err) {
      console.warn('Failed to refresh applications history', err);
    }
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSuccess(null);

    if (!selectedBankId) {
      setError('Выберите банк для отправки заявки.');
      return;
    }

    const rawAmount =
      bankAmounts[selectedBankId] && bankAmounts[selectedBankId].trim().length > 0
        ? bankAmounts[selectedBankId]
        : loanAmount.toString();
    const sanitized = sanitizeAmount(rawAmount);

    if (!sanitized || Number(sanitized) <= 0) {
      setError('Введите корректную сумму заявки.');
      return;
    }

    setSubmitting(true);

    try {
      await applicationsAPI.create({
        bank_id: selectedBankId,
        type_code: selectedTemplate,
        requested_amount: sanitized,
      });

      setSuccess('Заявка успешно отправлена. Мы уведомим вас о статусе.');
      setBankAmounts((prev) => ({
        ...prev,
        [selectedBankId]: sanitizeAmount(loanAmount),
      }));
      await refreshHistory();
    } catch (err: any) {
      console.error('Failed to submit application', err);
      const apiMessage = err?.response?.data?.error || err?.response?.data?.message;
      setError(apiMessage || 'Не удалось отправить заявку. Попробуйте снова позже.');
    } finally {
      setSubmitting(false);
    }
  };

  const handleCancel = () => navigate('/applications');

  return (
    <Layout>
      <div className="new-application">
        <div className="new-application__header">
          <div>
            <h1 className="new-application__title">Создание заявки на кредитный продукт</h1>
            <p className="new-application__subtitle">
              Выберите нужный шаблон, настройте параметры и отправьте заявку на рассмотрение в банк.
            </p>
          </div>
          <Button variant="outline" onClick={handleCancel}>
            Вернуться к списку
          </Button>
        </div>

        {error && <div className="new-application__alert new-application__alert--error">{error}</div>}
        {success && <div className="new-application__alert new-application__alert--success">{success}</div>}

        <div className="new-application__templates">
          {CREDIT_TEMPLATES.map((template) => {
            const isSelected = template.code === selectedTemplate;
            return (
              <Card
                key={template.code}
                className={`new-application__template-card ${isSelected ? 'new-application__template-card--selected' : ''}`}
                variant={isSelected ? 'elevated' : 'default'}
              >
                {template.highlight && <span className="new-application__template-badge">{template.highlight}</span>}
                <h2 className="new-application__template-title">{template.title}</h2>
                <p className="new-application__template-description">{template.description}</p>
                <button
                  type="button"
                  className={`new-application__template-action ${isSelected ? 'is-active' : ''}`}
                  onClick={() => setSelectedTemplate(template.code)}
                >
                  {isSelected ? 'Активен' : 'Выбрать шаблон'}
                </button>
              </Card>
            );
          })}
        </div>

        <div className="new-application__body">
          <div className="new-application__column">
            <Card className="new-application__form-card" variant="elevated">
              <form className="new-application__form" onSubmit={handleSubmit}>
                <div className="new-application__form-group">
                  <span className="new-application__chips-label">Выберите кредитный продукт</span>
                  <div className="new-application__chips new-application__chips--products">
                    {CREDIT_TEMPLATES.map((template) => {
                      const isActive = selectedTemplate === template.code;
                      return (
                        <button
                          key={`product-${template.code}`}
                          type="button"
                          className={`new-application__chip ${isActive ? 'is-active' : ''}`}
                          onClick={() => setSelectedTemplate(template.code)}
                        >
                          {template.title}
                        </button>
                      );
                    })}
                  </div>
                </div>

                <div className="new-application__form-group">
                  <Input
                    type="number"
                    label="Сумма заявки, ₽"
                    value={loanAmount}
                    onChange={handleAmountInputChange}
                    min={0}
                    max={1_000_000}
                    step={1_000}
                    fullWidth
                    helperText="Укажите ориентировочную сумму. Максимум 1 000 000 ₽."
                  />
                  <div className="new-application__slider">
                    <input
                      type="range"
                      min={0}
                      max={1_000_000}
                      step={5_000}
                      value={loanAmount}
                      onChange={handleSliderChange}
                      aria-valuemin={0}
                      aria-valuemax={1_000_000}
                      aria-valuenow={loanAmount}
                    />
                    <div className="new-application__slider-scale">
                      <span>0 ₽</span>
                      <span>1 000 000 ₽</span>
                    </div>
                  </div>
                </div>

                <div className="new-application__summary">
                  <div>
                    <span className="new-application__summary-label">Выбранный шаблон</span>
                    <strong className="new-application__summary-value">
                      {CREDIT_TEMPLATES.find((template) => template.code === selectedTemplate)?.title}
                    </strong>
                  </div>
                  <div>
                    <span className="new-application__summary-label">Сумма заявки</span>
                    <strong className="new-application__summary-value">{formatCurrency(loanAmount)}</strong>
                  </div>
                </div>

                <div className="new-application__actions">
                  <Button
                    type="submit"
                    variant="primary"
                    size="lg"
                    isLoading={submitting}
                    fullWidth
                  >
                    Отправить заявку
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    onClick={handleCancel}
                    fullWidth
                  >
                    Отменить
                  </Button>
                </div>
              </form>
            </Card>

            <Card className="new-application__banks-card" variant="default">
              <div className="new-application__banks-header">
                <div>
                  <h2 className="new-application__section-title">Выбор банка</h2>
                  <p className="new-application__section-subtitle">
                    Укажите сумму для каждого банка и подтвердите выбор.
                  </p>
                </div>
              </div>

              <div className="new-application__banks-filter">
                <div className="new-application__chips">
                  {banks.map((bank) => {
                    const isActive = bank.bank_id === selectedBankId;
                    return (
                      <button
                        key={`${bank.bank_id}-chip`}
                        type="button"
                        className={`new-application__chip ${isActive ? 'is-active' : ''}`}
                        onClick={() => handleSelectBank(bank.bank_id)}
                      >
                        {bank.name}
                      </button>
                    );
                  })}
                </div>
              </div>

              <div className="new-application__banks-list">
                {initialLoading ? (
                  <div className="new-application__banks-placeholder">Загружаем список банков...</div>
                ) : (
                  banks.map((bank) => {
                    const isSelected = bank.bank_id === selectedBankId;
                    const amountValue = bankAmounts[bank.bank_id] ?? '';
                    const sanitizedAmount = Number(sanitizeAmount(amountValue || loanAmount)) || 0;
                    const progress = Math.min(100, Math.max(0, (sanitizedAmount / 1_000_000) * 100));
                    const bankInfo = BANK_META[bank.code] ?? { rate: '—' };

                    return (
                      <div
                        key={bank.bank_id}
                        className={`new-application__bank ${isSelected ? 'new-application__bank--selected' : ''}`}
                      >
                        <div className="new-application__bank-header">
                          <div>
                            <div className="new-application__bank-chip">
                              {bank.logo && <img src={bank.logo} alt={bank.name} className="new-application__bank-logo" />}
                              <div>
                                <div className="new-application__bank-name">{bank.name}</div>
                                <div className="new-application__bank-code">{bank.code}</div>
                              </div>
                            </div>
                            <div className="new-application__bank-rate">
                              {bankInfo.label && <span className="new-application__bank-pill">{bankInfo.label}</span>}
                              <span>{bankInfo.rate}</span>
                            </div>
                          </div>
                          <Button
                            type="button"
                            variant={isSelected ? 'primary' : 'outline'}
                            size="sm"
                            onClick={() => handleSelectBank(bank.bank_id)}
                          >
                            {isSelected ? 'Выбрано' : 'Выбрать'}
                          </Button>
                        </div>

                        <div className="new-application__bank-fields">
                          <div className="new-application__bank-progress">
                            <div className="new-application__bank-progress-bar">
                              <div
                                className="new-application__bank-progress-fill"
                                style={{ width: `${progress}%` }}
                              />
                            </div>
                            <div className="new-application__bank-progress-scale">
                              <span>0 ₽</span>
                              <span>1 млн ₽</span>
                            </div>
                          </div>
                          <Input
                            type="text"
                            label="Сумма для банка, ₽"
                            value={amountValue}
                            onChange={(event) => handleBankAmountChange(bank.bank_id, event.target.value)}
                            placeholder={loanAmount.toString()}
                            fullWidth
                          />
                        </div>
                      </div>
                    );
                  })
                )}
              </div>
            </Card>
          </div>

          <Card className="new-application__history-card" variant="default">
            <div className="new-application__history-header">
              <div>
                <h2 className="new-application__section-title">История заявок</h2>
                <p className="new-application__section-subtitle">
                  Последние заявки и их статусы. Отслеживайте прогресс в реальном времени.
                </p>
              </div>
              <Button variant="ghost" size="sm" onClick={() => navigate('/applications')}>
                Вся история
              </Button>
            </div>
            <div className="new-application__history-list">
              {initialLoading ? (
                <div className="new-application__history-placeholder">Загружаем историю заявок...</div>
              ) : decoratedHistory.length === 0 ? (
                <div className="new-application__history-empty">
                  У вас пока нет созданных заявок. После отправки они появятся здесь.
                </div>
              ) : (
                decoratedHistory.map((item) => (
                  <div key={item.id} className="new-application__history-item">
                    <div className="new-application__history-left">
                      {item.logo ? (
                        <img src={item.logo} alt={item.bankName} className="new-application__history-logo" />
                      ) : (
                        <div className="new-application__history-avatar">{item.bankName.charAt(0)}</div>
                      )}
                      <div className="new-application__history-main">
                        <div className="new-application__history-title">{item.title}</div>
                        <div className="new-application__history-bank">{item.bankName}</div>
                      </div>
                    </div>
                    <div className="new-application__history-meta">
                      <span className="new-application__history-amount">{item.formattedAmount}</span>
                      <span className={`new-application__history-status new-application__history-status--${item.statusTone}`}>
                        {item.statusLabel}
                      </span>
                      <span className="new-application__history-date">{item.formattedDate}</span>
                    </div>
                  </div>
                ))
              )}
            </div>
          </Card>
        </div>
      </div>
    </Layout>
  );
};


