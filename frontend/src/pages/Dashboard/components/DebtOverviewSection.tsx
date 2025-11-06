import React from 'react';
import { Transaction } from '../types';
import './DebtOverviewSection.css';

interface DebtOverviewSectionProps {
  transactions: Transaction[];
  onFilterChange?: (filter: string) => void;
}

export const DebtOverviewSection: React.FC<DebtOverviewSectionProps> = ({
  transactions,
  onFilterChange,
}) => {
  return (
    <section className="debt-overview-section">
      <div className="debt-overview-section__content">
        <header className="debt-overview-section__header">
          <h2 className="debt-overview-section__title">История трат</h2>
          <button
            type="button"
            className="debt-overview-section__filter"
            onClick={() => onFilterChange?.('all')}
          >
            Все траты
          </button>
        </header>
        
        <ul className="debt-overview-section__list">
          {transactions.map((transaction) => (
            <li key={transaction.id} className="debt-overview-section__item">
              <div className="debt-overview-section__item-content">
                {transaction.image && (
                  <img
                    className="debt-overview-section__icon"
                    alt={transaction.title}
                    src={transaction.image}
                  />
                )}
                <div className="debt-overview-section__item-info">
                  {transaction.company && (
                    <div className="debt-overview-section__item-company">{transaction.company}</div>
                  )}
                  <div className="debt-overview-section__item-title">{transaction.title}</div>
                </div>
              </div>
              <div
                className={`debt-overview-section__item-amount ${
                  transaction.isPositive
                    ? 'debt-overview-section__item-amount--positive'
                    : ''
                }`}
              >
                {transaction.amount}
              </div>
            </li>
          ))}
        </ul>
      </div>
    </section>
  );
};

