import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Input } from '../../components';
import { useAuth } from '../../context/AuthContext';
import dicesImg from '../../utils/img/dices.png';
import gosuslugiLogo from '../../utils/img/gosuslugi-logo.png';
import maxLogo from '../../utils/img/Max_logo_2025.png';
import pochtaLogo from '../../utils/img/pochta-logo.png';
import adminIcon from '../../utils/img/profile.png';
import easyfundLogo from '../../utils/img/easyfund-logo.png';
import './Login.css';

type SocialLoginOption = {
  id: string;
  label: string;
  icon: string;
  alt: string;
  onClick: () => void | Promise<void>;
  isRandom?: boolean;
  isAdmin?: boolean;
};

const RANDOM_USERS = [
  { email: 'ivan@example.com', password: 'password123' },
  { email: 'anna@example.com', password: 'password123' },
  { email: 'pavel@example.com', password: 'password123' },
  { email: 'olga@example.com', password: 'password123' },
  { email: 'sergey@example.com', password: 'password123' },
];

const ADMIN_USER = { email: 'ivan@example.com', password: 'password123' };

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

  const handleRandomLogin = async () => {
    if (loading) {
      return;
    }

    setError('');
    setLoading(true);

    const randomUser = RANDOM_USERS[Math.floor(Math.random() * RANDOM_USERS.length)];

    setFormData((prev) => ({
      ...prev,
      email: randomUser.email,
      password: randomUser.password,
    }));

    try {
      await login(randomUser.email, randomUser.password);
      navigate('/welcome');
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Не удалось войти как случайный пользователь');
    } finally {
      setLoading(false);
    }
  };

  const handleAdminLogin = async () => {
    if (loading) {
      return;
    }

    setError('');
    setLoading(true);

    const { email, password } = ADMIN_USER;

    setFormData((prev) => ({
      ...prev,
      email,
      password,
    }));

    try {
      await login(email, password);
      navigate('/welcome');
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Не удалось войти как администратор');
    } finally {
      setLoading(false);
    }
  };

  const socialLogins: SocialLoginOption[] =
    formData.role === 'client'
      ? [
          {
            id: 'random-user',
            label: 'Случайный пользователь',
            icon: dicesImg,
            alt: 'Случайный пользователь',
            onClick: handleRandomLogin,
            isRandom: true,
          },
          {
            id: 'gosuslugi',
            label: 'Госуслуги',
            icon: gosuslugiLogo,
            alt: 'Госуслуги',
            onClick: () => {
              console.log('Gosuslugi login clicked');
            },
          },
          {
            id: 'max',
            label: 'Макс',
            icon: maxLogo,
            alt: 'Банк Макс',
            onClick: () => {
              console.log('Max login clicked');
            },
          },
          {
            id: 'pochta',
            label: 'Почта',
            icon: pochtaLogo,
            alt: 'Почта банк',
            onClick: () => {
              console.log('Pochta login clicked');
            },
          },
        ]
      : [
          {
            id: 'admin',
            label: 'Войти как админ',
            icon: adminIcon,
            alt: 'Войти как администратор',
            onClick: handleAdminLogin,
            isAdmin: true,
          },
        ];

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
              <h3 className="login-page__role-title">Я клиент</h3>
              <p className="login-page__role-description">
                Возможность управлять своими финансами в одном сервисе. Использование кредитных продуктов и многое другое.
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
              Полный контроль над входящими заявками. Проверяйте статус, запрашивайте дополнительные документы и выносите решение где угодно — в офисе, в дороге или дома.
              </p>
            </div>
          </div>
        </div>

        {/* Login form */}
        <form onSubmit={handleSubmit} className="login-page__form">
          <Input
            type="email"
            name="email"
            label={formData.role === 'bank_employee' ? 'Корпоративная почта' : 'Почта'}
            placeholder={formData.role === 'bank_employee' ? 'nastyasber@sber.ru' : 'yanavtb@ya.ru'}
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
          <div
            className={`login-page__social-buttons${formData.role === 'bank_employee' ? ' login-page__social-buttons--single' : ''}`}
          >
            {socialLogins.map(({ id, label, icon, alt, onClick, isRandom, isAdmin }) => (
              <button
                key={id}
                type="button"
                className={[
                  'login-page__social-btn',
                  isRandom ? 'login-page__social-btn--accent login-page__social-btn--random' : '',
                  isAdmin ? 'login-page__social-btn--accent login-page__social-btn--admin' : '',
                ]
                  .filter(Boolean)
                  .join(' ')}
                onClick={onClick}
                disabled={loading}
              >
                <img src={icon} alt={alt} className="login-page__social-logo" />
                <span>{label}</span>
              </button>
            ))}
          </div>
        </div>

        {formData.role === 'bank_employee' && (
          <div className="login-page__role-note">
            Чтобы <span>узнать свои данные</span>, обратитесь к официальному <span>представителю</span> вашего банка
          </div>
        )}
      </div>
      <div className="login-page__footer">
        <img src={easyfundLogo} alt="EasyFund" className="login-page__footer-logo" />
      </div>
    </div>
  );
};