import React from 'react';
import './PaymentHistorySection.css';

interface PaymentHistorySectionProps {
  totalDebt: number;
  creditCount: number;
  creditCardCount: number;
  onViewAllProducts?: () => void;
}

export const PaymentHistorySection: React.FC<PaymentHistorySectionProps> = ({
  totalDebt,
  creditCount,
  creditCardCount,
  onViewAllProducts,
}) => {
  const formatAmount = (amount: number): string => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  return (
    <section className="payment-history-section">
      <div className="payment-history-section__content">
        <h2 className="payment-history-section__title">Задолженность</h2>
        <p className="payment-history-section__amount">{formatAmount(totalDebt)}</p>
        <div className="payment-history-section__actions">
          <button
            type="button"
            className="payment-history-section__btn"
            onClick={onViewAllProducts}
          >
            Все продукты
          </button>
          <div className="payment-history-section__stats">
            <div className="payment-history-section__stat">
              <div className="payment-history-section__stat-badge">{creditCount}</div>
              <span className="payment-history-section__stat-label">Кредита</span>
            </div>
            <div className="payment-history-section__stat">
              <div className="payment-history-section__stat-badge">{creditCardCount}</div>
              <span className="payment-history-section__stat-label">Кредитных карт</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

