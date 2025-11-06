import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import easyfundLogo from '../../utils/img/easyfund-logo.png';
import './Welcome.css';

export const Welcome: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [isVisible, setIsVisible] = useState(false);
  const [userName, setUserName] = useState('Пользователь');

  useEffect(() => {
    // Обновляем имя пользователя из localStorage, если user еще не загружен
    const updateUserFromStorage = () => {
      const storedUser = localStorage.getItem('user');
      if (storedUser) {
        try {
          const parsedUser = JSON.parse(storedUser);
          setUserName(parsedUser.name || parsedUser.email?.split('@')[0] || 'Пользователь');
        } catch (e) {
          console.error('Error parsing user from storage:', e);
        }
      } else if (user) {
        setUserName(user.name || user.email?.split('@')[0] || 'Пользователь');
      }
    };

    updateUserFromStorage();
    
    // Плавное появление текста
    const timer = setTimeout(() => {
      setIsVisible(true);
    }, 100);

    // Автоматический редирект через 4 секунды
    const redirectTimer = setTimeout(() => {
      navigate('/dashboard');
    }, 4000);

    return () => {
      clearTimeout(timer);
      clearTimeout(redirectTimer);
    };
  }, [navigate, user]);

  // Обновляем имя, если user изменился
  useEffect(() => {
    if (user) {
      setUserName(user.name || user.email?.split('@')[0] || 'Пользователь');
    }
  }, [user]);

  return (
    <div className="welcome-page">
      {/* Animated gradient background */}
      <div className="welcome-page__background">
        <div className="welcome-page__gradient-shimmer"></div>
      </div>

      {/* Welcome header */}
      <header className={`welcome-page__header ${isVisible ? 'welcome-page__header--visible' : ''}`}>
        <h1 className="welcome-page__title">
          Добро пожаловать, {userName}!
        </h1>
      </header>

      {/* Footer logo */}
      <div className="welcome-page__footer">
        <img src={easyfundLogo} alt="EasyFund" className="welcome-page__footer-logo" />
      </div>
    </div>
  );
};
