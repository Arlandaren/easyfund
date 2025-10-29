import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Layout, Button, Input, Card } from '../../components';
import { useAuth } from '../../context/AuthContext';
import './Login.css';

export const Login: React.FC = () => {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [formData, setFormData] = useState({
    email: 'admin@easyfund.com',
    password: 'admin123',
    role: 'borrower' as 'borrower' | 'bank_employee',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleRoleChange = (role: 'borrower' | 'bank_employee') => {
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
      // Mock authentication for testing - bypass backend
      await login(formData.email, formData.password);
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Invalid email or password');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="login">
        <div className="login__container">
          <div className="login__left">
            <h2 className="login__title">
              Начни <span className="login__title-highlight">управлять</span> счетами
            </h2>

            <form onSubmit={handleSubmit} className="login__form">
              <Input
                type="email"
                name="email"
                label="Почта"
                placeholder="yanavtb@ya.ru"
                value={formData.email}
                onChange={handleChange}
                className="login__input"
                fullWidth
              />

              <Input
                type="password"
                name="password"
                label="Пароль"
                placeholder="************"
                value={formData.password}
                onChange={handleChange}
                className="login__input"
                fullWidth
              />

              {error && <div className="login__error">{error}</div>}

              <Button
                type="submit"
                fullWidth
                isLoading={loading}
                size="lg"
                className="login__submit-btn"
              >
                Войти
              </Button>
            </form>

            <button className="login__gosuslugi-btn">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="#E01E5A"/>
                <path d="M2 17L12 22L22 17V12L12 17L2 12V17Z" fill="#06B6D4"/>
                <path d="M22 12L12 17L2 12L12 7L22 12Z" fill="#C2EFF3"/>
              </svg>
              Госуслуги
            </button>

            <p className="login__helper-text">
              Или войдите с помощью
            </p>
          </div>

          <div className="login__right">
            <div
              className={`login__role-card ${formData.role === 'borrower' ? 'login__role-card--selected' : ''}`}
              onClick={() => handleRoleChange('borrower')}
            >
              <div className="login__role-radio">
                {formData.role === 'borrower' && <div className="login__role-radio-inner"></div>}
              </div>
              <div className="login__role-content">
                <h3 className="login__role-title">Я клиент банка</h3>
                <p className="login__role-description">
                  Возможность управлять своими финансами в 1 сервисе. Использование кредитных продуктов и многое другое.
                </p>
              </div>
            </div>

            <div
              className={`login__role-card ${formData.role === 'bank_employee' ? 'login__role-card--selected' : ''}`}
              onClick={() => handleRoleChange('bank_employee')}
            >
              <div className="login__role-radio">
                {formData.role === 'bank_employee' && <div className="login__role-radio-inner"></div>}
              </div>
              <div className="login__role-content">
                <h3 className="login__role-title">Я сотрудник банка</h3>
                <p className="login__role-description">
                  Возможность управлять своими финансами в 1 сервисе. Использование кредитных продуктов и многое другое.
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="login__info">
          <p className="login__test-info">
            Test Account: admin@easyfund.com / admin123
          </p>
        </div>
      </div>
    </Layout>
  );
};

