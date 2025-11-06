import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Input } from '../../components';
import { useAuth } from '../../context/AuthContext';
import gosuslugiLogo from '../../utils/img/gosuslugi-logo.png';
import easyfundLogo from '../../utils/img/easyfund-logo.png';
import './Login.css';

export const Login: React.FC = () => {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [formData, setFormData] = useState({
    email: 'yanavtb@ya.ru',
    password: '',
    role: 'client' as 'client' | 'bank_employee',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleRoleChange = (role: 'client' | 'bank_employee') => {
    setFormData({
      ...formData,
      role,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login(formData.email, formData.password);
      navigate('/welcome');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Неверный email или пароль');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      {/* Background decorative element */}
      <div className="login-page__background"></div>
      
      {/* Main content card */}
      <div className="login-page__card">
        <h1 className="login-page__title">
          Начни <span className="login-page__title-highlight">управлять счетами</span>
        </h1>

        {/* Role selection */}
        <div className="login-page__roles">
          <div
            className={`login-page__role-card ${formData.role === 'client' ? 'login-page__role-card--selected' : ''}`}
            onClick={() => handleRoleChange('client')}
          >
            <div className="login-page__role-radio">
              {formData.role === 'client' && <div className="login-page__role-radio-inner"></div>}
            </div>
            <div className="login-page__role-content">
              <h3 className="login-page__role-title">Я клиент банка</h3>
              <p className="login-page__role-description">
                Возможность управлять своими финансами в нашем сервисе. Использование кредитных продуктов.
              </p>
            </div>
          </div>

          <div
            className={`login-page__role-card ${formData.role === 'bank_employee' ? 'login-page__role-card--selected' : ''}`}
            onClick={() => handleRoleChange('bank_employee')}
          >
            <div className="login-page__role-radio">
              {formData.role === 'bank_employee' && <div className="login-page__role-radio-inner"></div>}
            </div>
            <div className="login-page__role-content">
              <h3 className="login-page__role-title">Я сотрудник банка</h3>
              <p className="login-page__role-description">
                Возможность управлять своими финансами в нашем сервисе. Использование кредитных продуктов.
              </p>
            </div>
          </div>
        </div>

        {/* Login form */}
        <form onSubmit={handleSubmit} className="login-page__form">
          <Input
            type="email"
            name="email"
            label="Почта"
            placeholder="yanavtb@ya.ru"
            value={formData.email}
            onChange={handleChange}
            className="login-page__input"
            fullWidth
          />

          <Input
            type="password"
            name="password"
            label="Пароль"
            placeholder="Введите пароль"
            value={formData.password}
            onChange={handleChange}
            className="login-page__input"
            fullWidth
            autoComplete="current-password"
          />

          {error && <div className="login-page__error">{error}</div>}

          <Button
            type="submit"
            fullWidth
            isLoading={loading}
            size="lg"
            className="login-page__submit-btn"
          >
            Войти
          </Button>
        </form>

        {/* Social login */}
        <div className="login-page__social">
          <p className="login-page__social-text">Или войдите с помощью</p>
          <div className="login-page__social-buttons">
            {[1, 2, 3, 4].map((i) => (
              <button
                key={i}
                type="button"
                className="login-page__gosuslugi-btn"
                onClick={() => {
                  // Handle Gosuslugi login
                  console.log('Gosuslugi login clicked');
                }}
              >
                <img src={gosuslugiLogo} alt="Госуслуги" className="login-page__gosuslugi-logo" />
                <span>Госуслуги</span>
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Footer logo */}
      <div className="login-page__footer">
        <img src={easyfundLogo} alt="EasyFund" className="login-page__footer-logo" />
      </div>
    </div>
  );
};