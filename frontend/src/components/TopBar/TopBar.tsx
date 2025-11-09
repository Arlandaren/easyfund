import React, { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import easyfundLogoSvg from '../../utils/img/easyfund-logo.svg';
import profileImage from '../../utils/img/profile.png';
import './TopBar.css';

type TopBarVariant = 'static' | 'overlay';

interface TopBarProps {
  variant?: TopBarVariant;
  className?: string;
}

export const TopBar: React.FC<TopBarProps> = ({ variant = 'static', className = '' }) => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!dropdownOpen) {
      return;
    }

    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setDropdownOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [dropdownOpen]);

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  return (
    <header
      className={`topbar topbar--${variant} ${className}`.trim()}
    >
      <div className="topbar__container">
        <button
          className="topbar__logo-button"
          onClick={() => navigate('/dashboard')}
          type="button"
          aria-label="Перейти на дашборд"
        >
          <img className="topbar__logo" alt="EasyFund Logo" src={easyfundLogoSvg} />
        </button>

        <div className="topbar__actions">
          <button className="topbar__icon" type="button" aria-label="Поиск">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path
                d="M21 21L15 15M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>
          <button className="topbar__icon" type="button" aria-label="Уведомления">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path
                d="M18 8C18 6.4087 17.3679 4.88258 16.2426 3.75736C15.1174 2.63214 13.5913 2 12 2C10.4087 2 8.88258 2.63214 7.75736 3.75736C6.63214 4.88258 6 6.4087 6 8C6 15 3 17 3 17H21C21 17 18 15 18 8Z"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
              <path
                d="M13.73 21C13.5542 21.3031 13.3019 21.5547 12.9982 21.7295C12.6946 21.9044 12.3504 21.9965 12 21.9965C11.6496 21.9965 11.3054 21.9044 11.0018 21.7295C10.6982 21.5547 10.4458 21.3031 10.27 21"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>
          <div className="topbar__avatar" ref={dropdownRef}>
            <button
              className="topbar__avatar-button"
              onClick={() => setDropdownOpen(!dropdownOpen)}
              type="button"
              aria-label="Меню пользователя"
            >
              <img
                className="topbar__avatar-image"
                src={profileImage}
                alt={user?.full_name || user?.email || 'Пользователь'}
              />
            </button>
            {dropdownOpen && (
              <div className="topbar__dropdown">
                <button
                  className="topbar__dropdown-item"
                  type="button"
                  onClick={() => {
                    setDropdownOpen(false);
                    navigate('/dashboard');
                  }}
                >
                  Профиль
                </button>
                <button
                  className="topbar__dropdown-item"
                  type="button"
                  onClick={() => {
                    setDropdownOpen(false);
                    handleLogout();
                  }}
                >
                  Выйти
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};


