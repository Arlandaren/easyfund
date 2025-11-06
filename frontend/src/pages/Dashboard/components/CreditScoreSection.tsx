import React from 'react';
import './CreditScoreSection.css';

interface CreditScoreSectionProps {
  accountBalance: number;
  onTransfer?: () => void;
  onTopUp?: () => void;
}

export const CreditScoreSection: React.FC<CreditScoreSectionProps> = ({
  accountBalance,
  onTransfer,
  onTopUp,
}) => {
  const formatBalance = (amount: number): string => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  return (
    <section className="credit-score-section">
      <div className="credit-score-section__content">
        <h2 className="credit-score-section__title">Общий счет</h2>
        <p className="credit-score-section__balance">{formatBalance(accountBalance)}</p>
        <div className="credit-score-section__actions">
          <button
            type="button"
            onClick={onTransfer}
            className="credit-score-section__btn credit-score-section__btn--primary"
          >
            Перевести
          </button>
          <button
            type="button"
            onClick={onTopUp}
            className="credit-score-section__btn credit-score-section__btn--outline"
          >
            Пополнить
          </button>
        </div>
      </div>
    </section>
  );
};

