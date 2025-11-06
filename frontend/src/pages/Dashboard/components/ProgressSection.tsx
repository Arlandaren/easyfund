import React from 'react';
import { ProgressData } from '../types';
import './ProgressSection.css';

interface ProgressSectionProps {
  progress: ProgressData;
}

export const ProgressSection: React.FC<ProgressSectionProps> = ({
  progress,
}) => {
  const formatAmount = (amount: number): string => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  // Calculate progress bar width (max 268px)
  const progressWidth = Math.min(268, (progress.percentage / 100) * 268);

  return (
    <section className="progress-section">
      <h2 className="progress-section__title">Вы почти у цели!</h2>
      <p className="progress-section__description">
        Благодаря своим усилиям, вы почти
        <br />
        закрыли свои задолженности
      </p>
      
      <div className="progress-section__progress-container">
        <div className="progress-section__progress-bar-bg">
          <div
            className="progress-section__progress-bar-fill"
            style={{ width: `${progressWidth}px` }}
          />
        </div>
        <div className="progress-section__progress-labels">
          <span className="progress-section__progress-label">{formatAmount(progress.initialDebt)}</span>
          <span className="progress-section__progress-label">{formatAmount(progress.currentDebt)}</span>
          <span className="progress-section__progress-label">{formatAmount(progress.targetDebt)}</span>
        </div>
        <div className="progress-section__progress-percentage">{progress.percentage}%</div>
      </div>
      
      <div className="progress-section__encouragement">Дальше - больше!</div>
    </section>
  );
};

