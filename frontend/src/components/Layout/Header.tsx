import React from 'react';
import { Link } from 'react-router-dom';
import './Header.css';

interface HeaderProps {
  user?: {
    name: string;
    email: string;
    role: string;
  };
  onLogout?: () => void;
}

export const Header: React.FC<HeaderProps> = ({ user, onLogout }) => {
  return (
    <header className="header">
      <div className="header__container">
        <Link to="/" className="header__logo">
          <h1>EasyFund</h1>
        </Link>
        <nav className="header__nav">
          <Link to="/dashboard" className="header__nav-link">
            Dashboard
          </Link>
          <Link to="/applications" className="header__nav-link">
            Applications
          </Link>
          {user?.role === 'bank_risk_manager' || user?.role === 'bank_analyst' ? (
            <Link to="/bank/dashboard" className="header__nav-link">
              Bank Dashboard
            </Link>
          ) : null}
        </nav>
        <div className="header__user">
          {user ? (
            <>
              <span className="header__user-name">{user.name}</span>
              <span className="header__user-role">{user.role}</span>
              {onLogout && (
                <button onClick={onLogout} className="header__logout-btn">
                  Logout
                </button>
              )}
            </>
          ) : (
            <Link to="/login" className="header__nav-link">
              Login
            </Link>
          )}
        </div>
      </div>
    </header>
  );
};

