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

interface BankDistributionSegment {
  bankId: number;
  bankName: string;
  share: number;
  amount: number;
  color: string;
}

interface ApplicationHistoryItem {
  id: number;
  bankId: number | null;
  typeCode: string;
  requestedAmount: number;
  status: string;
  submittedAt?: string;
  distribution?: BankDistributionSegment[];
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

const BANK_COLORS: Record<string, string> = {
  TBANK: '#FFDE34',
  VTB: '#002782',
  SBER: '#046A38',
  ALFA: '#EF3125',
  OPT: '#C2FF05',
};

const buildDistributionGradient = (segments: BankDistributionSegment[]): string => {
  if (!segments.length) {
    return 'conic-gradient(#e2e8f0 0 100%)';
  }

  let offset = 0;
  const stops = segments
    .map((segment, index) => {
      const start = offset;
      offset += segment.share;
      if (index === segments.length - 1) {
        offset = 100;
      }
      return `${segment.color} ${start}% ${offset}%`;
    })
    .join(', ');

  return `conic-gradient(${stops})`;
};

const FALLBACK_HISTORY: ApplicationHistoryItem[] = [];

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

const normalizeShares = (
  bankIds: number[],
  baseShares: Record<number, number>,
): Record<number, number> => {
  if (!bankIds.length) {
    return {};
  }

  const sanitized = bankIds.reduce<Record<number, number>>((acc, id) => {
    acc[id] = Math.max(0, baseShares[id] ?? 0);
    return acc;
  }, {});

  const total = bankIds.reduce((sum, id) => sum + sanitized[id], 0);

  if (total <= 0) {
    const evenShare = Math.floor(100 / bankIds.length);
    let remainder = 100 - evenShare * bankIds.length;
    const evenResult: Record<number, number> = {};
    bankIds.forEach((id) => {
      const bonus = remainder > 0 ? 1 : 0;
      evenResult[id] = evenShare + bonus;
      remainder -= bonus;
    });
    return evenResult;
  }

  const normalized: Record<number, number> = {};
  let remainder = 100;
  bankIds.forEach((id, index) => {
    let share: number;
    if (index === bankIds.length - 1) {
      share = remainder;
    } else {
      share = Math.round((sanitized[id] / total) * 100);
      share = Math.max(0, Math.min(share, remainder));
    }
    normalized[id] = share;
    remainder -= share;
  });

  return normalized;
};

export const NewApplication: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuth();

  const [banks, setBanks] = useState<Bank[]>([]);
  const [history, setHistory] = useState<ApplicationHistoryItem[]>([]);
  const [selectedTemplate, setSelectedTemplate] = useState<string>(CREDIT_TEMPLATES[0].code);
  const [selectedBankIds, setSelectedBankIds] = useState<number[]>([]);
  const [loanAmount, setLoanAmount] = useState<number>(150_000);
  const [bankShares, setBankShares] = useState<Record<number, number>>({});

  const [initialLoading, setInitialLoading] = useState(false);
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
      const banksResponse = await banksAPI.getAll().catch(() => null);
      const historyResponse = await (user?.user_id
        ? applicationsAPI.getUserApplications(String(user.user_id)).catch(() => null)
        : Promise.resolve(null));

      if (banksResponse?.data && Array.isArray(banksResponse.data) && banksResponse.data.length > 0) {
        banksData = banksResponse.data.map((bank: Bank) => {
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

      if (historyResponse?.data && Array.isArray(historyResponse.data)) {
        historyData = historyResponse.data
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
    } catch (fetchError) {
      console.warn('Failed to load application dependencies', fetchError);
    }

      if (!isMounted) return;

      setBanks(banksData);
      setSelectedBankIds((prev) => {
        const validPrev = prev.filter((id) => banksData.some((bank) => bank.bank_id === id));
        if (validPrev.length) {
          return validPrev;
        }
        return banksData.slice(0, Math.min(2, banksData.length)).map((bank) => bank.bank_id);
      });
      setHistory(historyData);
      setInitialLoading(false);
    };

    loadData();

    return () => {
      isMounted = false;
    };
  }, [user?.user_id]);

  useEffect(() => {
    setBankShares((prev) => normalizeShares(selectedBankIds, prev));
  }, [selectedBankIds]);

  const decoratedHistory = useMemo(() => {
    return history.map((item) => {
      const template = CREDIT_TEMPLATES.find((t) => t.code === item.typeCode);
      const bank = banks.find((b) => b.bank_id === item.bankId);
      const statusMeta = STATUS_MAP[item.status] ?? STATUS_MAP.PENDING;
      const bankNameOverride =
        item.distribution && item.distribution.length > 1 ? 'Несколько банков' : bank?.name ?? 'Банк не выбран';

      return {
        ...item,
        title: template?.title ?? 'Заявка на кредит',
        bankName: bankNameOverride,
        statusLabel: statusMeta.label,
        statusTone: statusMeta.tone,
        formattedAmount: formatCurrency(item.requestedAmount || 0),
        formattedDate: item.submittedAt ? new Date(item.submittedAt).toLocaleDateString('ru-RU') : '—',
        logo: bank?.logo,
        distribution: item.distribution,
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

  const getBankAmount = (bankId: number) => {
    const share = bankShares[bankId] ?? 0;
    return Math.round((share / 100) * loanAmount);
  };

  const handleBankAmountChange = (bankId: number, amount: number) => {
    if (!selectedBankIds.includes(bankId)) {
      return;
    }

    const clamped = Math.max(0, Math.min(amount, loanAmount));
    const share = loanAmount === 0 ? 0 : Math.round((clamped / loanAmount) * 100);
    setBankShares((prev) => normalizeShares(selectedBankIds, { ...prev, [bankId]: share }));
  };

  const handleBankInputChange = (bankId: number, value: string) => {
    const sanitized = Number(sanitizeAmount(value)) || 0;
    handleBankAmountChange(bankId, sanitized);
  };

  const distributionData = useMemo<BankDistributionSegment[]>(() => {
    if (!selectedBankIds.length || loanAmount <= 0) {
      return [];
    }

    const entries = selectedBankIds
      .map((bankId) => {
        const bank = banks.find((item) => item.bank_id === bankId);
        if (!bank) {
          return null;
        }
        const share = Math.max(0, Math.min(100, bankShares[bankId] ?? 0));
        const amount = Math.round((share / 100) * loanAmount);
        const color = BANK_COLORS[bank.code] ?? '#189cf4';
        return {
          bankId,
          bankName: bank.name,
          share,
          amount,
          color,
        };
      })
      .filter((entry): entry is BankDistributionSegment => Boolean(entry));

    return entries.filter((entry) => entry.share > 0);
  }, [selectedBankIds, bankShares, loanAmount, banks]);

  const distributionGradient = useMemo(() => {
    if (!distributionData.length) {
      return 'conic-gradient(#e2e8f0 0 100%)';
    }

    let offset = 0;
    const segments = distributionData
      .map((entry, index) => {
        const start = offset;
        offset += entry.share;
        if (index === distributionData.length - 1) {
          offset = 100;
        }
        return `${entry.color} ${start}% ${offset}%`;
      })
      .join(', ');

    return `conic-gradient(${segments})`;
  }, [distributionData]);

  const toggleBankSelection = (bankId: number) => {
    setSelectedBankIds((prev) => {
      const exists = prev.includes(bankId);
      const updated = exists ? prev.filter((id) => id !== bankId) : [...prev, bankId];
      return updated;
    });
    setSuccess(null);
    setError(null);
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSuccess(null);

    if (!selectedBankIds.length) {
      setError('Выберите хотя бы один банк для отправки заявки.');
      return;
    }

    setSubmitting(true);

    try {
      const primaryBankId = selectedBankIds[0];
      const primaryAmount = getBankAmount(primaryBankId) || loanAmount;

      await applicationsAPI.create({
        bank_id: primaryBankId,
        type_code: selectedTemplate,
        requested_amount: primaryAmount,
      });

      setSuccess(
        selectedBankIds.length > 1
          ? `Заявка отправлена. Банков в списке: ${selectedBankIds.length}.`
          : 'Заявка успешно отправлена. Мы уведомим вас о статусе.',
      );

      const historyResponse = await applicationsAPI
        .getUserApplications(String(user?.user_id))
        .catch(() => null);

      const fallbackEntry: ApplicationHistoryItem = {
        id: Date.now(),
        bankId: primaryBankId,
        typeCode: selectedTemplate,
        requestedAmount: primaryAmount,
        status: 'PENDING',
        submittedAt: new Date().toISOString(),
      };

      const freshHistory: ApplicationHistoryItem[] =
        historyResponse?.data && Array.isArray(historyResponse.data)
          ? historyResponse.data
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
              .slice(0, 8)
          : [fallbackEntry, ...history].slice(0, 8);

      if (distributionData.length > 1) {
        const multiEntry: ApplicationHistoryItem = {
          id: Date.now() + 1,
          bankId: null,
          typeCode: selectedTemplate,
          requestedAmount: loanAmount,
          status: 'PENDING',
          submittedAt: new Date().toISOString(),
          distribution: distributionData,
        };
        setHistory([multiEntry, ...freshHistory].slice(0, 8));
      } else {
        setHistory(freshHistory);
      }
    } catch (err) {
      console.error('Failed to submit application, storing locally', err);

      const fallbackEntry: ApplicationHistoryItem = {
        id: Date.now(),
        bankId: distributionData.length > 1 ? null : selectedBankIds[0],
        typeCode: selectedTemplate,
        requestedAmount: loanAmount,
        status: 'PENDING',
        submittedAt: new Date().toISOString(),
        distribution: distributionData.length > 1 ? distributionData : undefined,
      };

      setHistory((prev) => [fallbackEntry, ...prev].slice(0, 8));
      setSuccess('Заявка сохранена. Мы уведомим вас, как только получим ответ от банков.');
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
            const isHighlighted = template.highlight && template.code === 'PERSONAL';
            const cardClasses = [
              'new-application__template-card',
              isSelected ? 'new-application__template-card--selected' : '',
              isHighlighted ? 'new-application__template-card--highlighted' : '',
            ]
              .filter(Boolean)
              .join(' ');

            return (
              <Card key={template.code} className={cardClasses} variant={isSelected ? 'elevated' : 'default'}>
                {template.highlight && (
                  <span
                    className={`new-application__template-badge ${
                      isHighlighted ? 'new-application__template-badge--compact' : ''
                    }`}
                  >
                    {template.highlight}
                  </span>
                )}
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
                    const isActive = selectedBankIds.includes(bank.bank_id);
                    return (
                      <button
                        key={`${bank.bank_id}-chip`}
                        type="button"
                        className={`new-application__chip ${isActive ? 'is-active' : ''}`}
                        onClick={() => toggleBankSelection(bank.bank_id)}
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
                    const isSelected = selectedBankIds.includes(bank.bank_id);
                    const bankInfo = BANK_META[bank.code] ?? { rate: '—' };
                    const allocatedAmount = isSelected ? getBankAmount(bank.bank_id) : 0;
                    const sliderMax = loanAmount;
                    const sliderStep = sliderMax > 0 ? Math.max(1, Math.round(sliderMax / 100)) : 1;
                    const progressPercent =
                      !isSelected || sliderMax === 0 ? 0 : Math.min(100, (allocatedAmount / sliderMax) * 100);

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
                            onClick={() => toggleBankSelection(bank.bank_id)}
                          >
                            {isSelected ? 'Убрать' : 'Выбрать'}
                          </Button>
                        </div>

                        <div className="new-application__bank-fields">
                          <div className="new-application__bank-progress">
                            <div className="new-application__bank-progress-bar">
                              <div
                                className="new-application__bank-progress-fill"
                                style={{ width: `${progressPercent}%` }}
                              />
                            </div>
                            <div className="new-application__bank-progress-scale">
                              <span>0 ₽</span>
                              <span>{formatCurrency(sliderMax)}</span>
                            </div>
                          </div>
                          <div className="new-application__bank-inputs">
                            <Input
                              type="number"
                              label="Сумма для банка, ₽"
                              value={allocatedAmount}
                              onChange={(event) => handleBankInputChange(bank.bank_id, event.target.value)}
                              placeholder={loanAmount.toString()}
                              fullWidth
                              disabled={!isSelected || sliderMax === 0}
                            />
                            <div className="new-application__bank-slider">
                              <input
                                type="range"
                                min={0}
                                max={sliderMax}
                                step={sliderStep}
                                value={allocatedAmount}
                                onChange={(event) =>
                                  handleBankAmountChange(bank.bank_id, Number(event.target.value) || 0)
                                }
                                disabled={!isSelected || sliderMax === 0}
                              />
                              <div className="new-application__bank-slider-label">
                                {isSelected ? formatCurrency(allocatedAmount) : '0 ₽'}
                              </div>
                            </div>
                          </div>
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
                      {item.distribution && item.distribution.length > 1 ? (
                        <div
                          className="new-application__history-distribution"
                          style={{ backgroundImage: buildDistributionGradient(item.distribution) }}
                          aria-label="Распределение по банкам"
                        />
                      ) : item.logo ? (
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


