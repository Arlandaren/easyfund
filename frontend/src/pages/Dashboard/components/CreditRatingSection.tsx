import React, { useMemo } from 'react';
import { CreditRatingData } from '../types';
import './CreditRatingSection.css';

interface CreditRatingSectionProps {
  creditRating: CreditRatingData;
}

export const CreditRatingSection: React.FC<CreditRatingSectionProps> = ({
  creditRating,
}) => {
  // Calculate position of the indicator on the slider (0-458px)
  const indicatorPosition = useMemo(() => {
    const range = creditRating.max - creditRating.min;
    const position = ((creditRating.score - creditRating.min) / range) * 458;
    return Math.max(0, Math.min(458, position));
  }, [creditRating]);

  return (
    <section className="credit-rating-section">
      <h2 className="credit-rating-section__title">Ваш рейтинг</h2>
      <div className="credit-rating-section__score">{creditRating.score}</div>
      
      <div className="credit-rating-section__slider-container">
        <div className="credit-rating-section__slider">
          <div className="credit-rating-section__slider-track" />
          <div
            className="credit-rating-section__slider-indicator"
            style={{ left: `${indicatorPosition}px` }}
          />
        </div>
        <div className="credit-rating-section__labels">
          {creditRating.labels.map((label, index) => (
            <span key={index} className="credit-rating-section__label">
              {label}
            </span>
          ))}
        </div>
      </div>
      
      <p className="credit-rating-section__description">
        Достаточно высокий шанс на одобрение кредитного продукта.
      </p>
    </section>
  );
};

