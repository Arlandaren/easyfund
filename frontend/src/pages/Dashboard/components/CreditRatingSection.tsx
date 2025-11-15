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

  const indicatorColor = useMemo(() => {
    const colorStops = [
      { stop: 0, color: '#DB1414' },
      { stop: 0.33, color: '#f59e0b' },
      { stop: 0.66, color: '#22c55e' },
      { stop: 1, color: '#0ea5e9' },
    ];

    const progress = Math.max(
      0,
      Math.min(1, (creditRating.score - creditRating.min) / (creditRating.max - creditRating.min))
    );

    for (let i = 0; i < colorStops.length - 1; i++) {
      const start = colorStops[i];
      const end = colorStops[i + 1];
      if (progress >= start.stop && progress <= end.stop) {
        const localProgress = (progress - start.stop) / (end.stop - start.stop);
        const startColor = start.color;
        const endColor = end.color;

        const interpolateChannel = (channel: number) => {
          const startValue = parseInt(startColor.slice(channel, channel + 2), 16);
          const endValue = parseInt(endColor.slice(channel, channel + 2), 16);
          return Math.round(startValue + (endValue - startValue) * localProgress)
            .toString(16)
            .padStart(2, '0');
        };

        const r = interpolateChannel(1);
        const g = interpolateChannel(3);
        const b = interpolateChannel(5);
        return `#${r}${g}${b}`;
      }
    }

    return colorStops[colorStops.length - 1].color;
  }, [creditRating]);

  return (
    <section className="credit-rating-section">
      <h2 className="credit-rating-section__title">Ваш рейтинг</h2>
      <div
        className="credit-rating-section__score"
        style={{ color: indicatorColor, left: `${indicatorPosition}px` }}
      >
        {creditRating.score}
      </div>
      
      <div className="credit-rating-section__slider-container">
        <div className="credit-rating-section__slider">
          <div className="credit-rating-section__slider-track" />
          <div
            className="credit-rating-section__slider-progress"
            style={{ width: `${indicatorPosition}px`, background: `linear-gradient(to right, #DB1414 0%, #f59e0b 33%, #22c55e 66%, #0ea5e9 100%)` }}
          />
          <div
            className="credit-rating-section__slider-indicator"
            style={{ left: `${indicatorPosition}px`, backgroundColor: indicatorColor }}
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

