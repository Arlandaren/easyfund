import React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import homeIcon from '../../utils/img/home.svg';
import applicationsIcon from '../../utils/img/applications.svg';
import './BottomNav.css';

const HIDDEN_PATHS = ['/login', '/welcome'];

const getMaskedIconStyle = (icon: string): React.CSSProperties => ({
  WebkitMaskImage: `url(${icon})`,
  maskImage: `url(${icon})`,
  WebkitMaskRepeat: 'no-repeat',
  maskRepeat: 'no-repeat',
  WebkitMaskPosition: 'center',
  maskPosition: 'center',
  WebkitMaskSize: 'contain',
  maskSize: 'contain',
  backgroundColor: 'currentColor',
});

export const BottomNav: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();

  if (HIDDEN_PATHS.includes(location.pathname)) {
    return null;
  }

  const isDashboard = location.pathname.startsWith('/dashboard');
  const isApplications = location.pathname.startsWith('/applications');

  const activeId = isDashboard ? 'dashboard' : isApplications ? 'applications' : null;
  const indicatorLeft = activeId === 'applications' ? 'calc(100% - 72px)' : '22px';

  return (
    <nav className="bottom-nav" aria-label="Main navigation">
      {activeId && (
        <div
          className="bottom-nav__indicator"
          style={{ left: indicatorLeft }}
          aria-hidden="true"
        />
      )}

      <button
        type="button"
        className={`bottom-nav__btn${
          activeId === 'dashboard'
            ? ' bottom-nav__btn--active'
            : activeId === 'applications'
              ? ' bottom-nav__btn--secondary-active'
              : ''
        }`}
        aria-label="Главная"
        onClick={() => navigate('/dashboard')}
      >
        <span
          className="bottom-nav__icon bottom-nav__icon--home"
          style={getMaskedIconStyle(homeIcon)}
          aria-hidden="true"
        />
      </button>

      <button
        type="button"
        className={`bottom-nav__btn${
          activeId === 'applications'
            ? ' bottom-nav__btn--active'
            : activeId === 'dashboard'
              ? ' bottom-nav__btn--secondary-active'
              : ''
        }`}
        aria-label="Заявки"
        onClick={() => navigate('/applications')}
      >
        <span
          className="bottom-nav__icon bottom-nav__icon--applications"
          style={getMaskedIconStyle(applicationsIcon)}
          aria-hidden="true"
        />
      </button>
    </nav>
  );
};
