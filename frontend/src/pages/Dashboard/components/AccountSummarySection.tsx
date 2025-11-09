import React from 'react';
import { PaymentItem } from '../types';
import './AccountSummarySection.css';

interface AccountSummarySectionProps {
  payments: PaymentItem[];
  onViewAll?: () => void;
}

export const AccountSummarySection: React.FC<AccountSummarySectionProps> = ({
  payments,
  onViewAll,
}) => {
  return (
    <section className="account-summary-section">
      <div className="account-summary-section__content">
        <div className="account-summary-section__header">
          <h2 className="account-summary-section__title">Платежи по кредитам</h2>
        </div>
        
        <ul className="account-summary-section__list">
          {payments.map((payment) => (
            <li key={payment.id} className="account-summary-section__item">
              <div className="account-summary-section__item-content">
                {payment.icon && (
                  <img
                    className="account-summary-section__icon"
                    alt={payment.title}
                    src={payment.icon}
                  />
                )}
                <div className="account-summary-section__item-info">
                  <div className="account-summary-section__item-title">{payment.title}</div>
                  <div className="account-summary-section__item-date">{payment.dueDate}</div>
                </div>
              </div>
              <div className="account-summary-section__item-amount">{payment.amount}</div>
            </li>
          ))}
        </ul>
        
        <button
          className="account-summary-section__view-all"
          onClick={onViewAll}
          type="button"
        >
          Все платежи
        </button>
      </div>
    </section>
  );
};