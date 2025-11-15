import React from 'react';
import { ProgressData } from '../types';
import './ProgressSection.css';

interface ProgressSectionProps {
  progress: ProgressData;
  totalDebt: number;
}

export const ProgressSection: React.FC<ProgressSectionProps> = ({
  progress,
  totalDebt,
}) => {
  const formatAmount = (amount: number): string => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const hasDebt = totalDebt > 0;
  const fallbackPercentage =
    typeof progress?.percentage === 'number' ? progress.percentage : 40;
  const paidPercentage = hasDebt
    ? Math.max(0, Math.min(100, fallbackPercentage))
    : 0;
  const progressWidth = Math.min(268, (paidPercentage / 100) * 268);

  return (
    <section className="progress-section">
      <h2 className="progress-section__title">
        {hasDebt ? 'Вы почти у цели!' : 'У вас нет задолженностей'}
      </h2>
      <p className="progress-section__description">
        Благодаря своим усилиям, вы почти
        <br />
        закрыли свои задолженности
      </p>
      
      <div className="progress-section__progress-info">
        <div className="progress-section__progress-percentage">{paidPercentage}%</div>
      </div>

      <div className="progress-section__progress-container">
        <div className="progress-section__progress-bar-bg">
          <div
            className="progress-section__progress-bar-fill"
            style={{ width: `${progressWidth}px` }}
          />
        </div>
        <div className="progress-section__progress-labels">
          <span className="progress-section__progress-label">{formatAmount(0)}</span>
          <span className="progress-section__progress-label">
            {formatAmount(hasDebt ? totalDebt : 0)}
          </span>
        </div>
      </div>
      
      <div className="progress-section__encouragement">Дальше - больше!</div>
    </section>
  );
};

