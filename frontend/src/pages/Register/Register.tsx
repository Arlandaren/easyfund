import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Layout, Button, Input, Card } from '../../components';
import { useAuth } from '../../context/AuthContext';
import './Register.css';

export const Register: React.FC = () => {
  const navigate = useNavigate();
  const { login, user } = useAuth();
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
    phone: '',
    role: 'borrower',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    setLoading(true);

    try {
      const response = await fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: formData.name,
          email: formData.email,
          password: formData.password,
          phone: formData.phone,
          role: formData.role,
        }),
      });

      if (response.ok) {
        // Автоматически логиним пользователя после регистрации
        try {
          await login(formData.email, formData.password);
          // Небольшая задержка для гарантии, что состояние обновилось
          setTimeout(() => {
            navigate('/welcome');
          }, 100);
        } catch (loginErr) {
          // Если логин не удался, всё равно показываем welcome (для случая, когда бэкенд недоступен)
          // Создаём mock пользователя для welcome страницы
          const mockUser = {
            id: '1',
            email: formData.email,
            name: formData.name,
            role: formData.role,
          };
          localStorage.setItem('token', 'mock_token_after_registration');
          localStorage.setItem('user', JSON.stringify(mockUser));
          setTimeout(() => {
            navigate('/welcome');
          }, 100);
        }
      } else {
        const data = await response.json();
        setError(data.message || 'Registration failed');
      }
    } catch (err: any) {
      // Если бэкенд недоступен, создаём mock пользователя и показываем welcome
      const mockUser = {
        id: '1',
        email: formData.email,
        name: formData.name,
        role: formData.role,
      };
      localStorage.setItem('token', 'mock_token_after_registration');
      localStorage.setItem('user', JSON.stringify(mockUser));
          // Переходим на welcome, AuthContext обновит состояние при следующей проверке
          navigate('/welcome', { replace: true });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="register">
        <Card variant="elevated" className="register__card">
          <h2 className="register__title">Create Account</h2>
          <p className="register__subtitle">Sign up to start your loan application journey.</p>

          <form onSubmit={handleSubmit} className="register__form">
            <Input
              type="text"
              name="name"
              label="Full Name"
              placeholder="Enter your full name"
              value={formData.name}
              onChange={handleChange}
              required
              fullWidth
            />

            <Input
              type="email"
              name="email"
              label="Email"
              placeholder="Enter your email"
              value={formData.email}
              onChange={handleChange}
              required
              fullWidth
            />

            <Input
              type="tel"
              name="phone"
              label="Phone Number"
              placeholder="Enter your phone number"
              value={formData.phone}
              onChange={handleChange}
              required
              fullWidth
            />

            <div className="register__form-group">
              <label htmlFor="role" className="input-label">
                Role
              </label>
              <select
                id="role"
                name="role"
                value={formData.role}
                onChange={handleChange}
                className="input"
                required
              >
                <option value="borrower">Borrower</option>
                <option value="bank_risk_manager">Bank Risk Manager</option>
                <option value="bank_analyst">Bank Analyst</option>
              </select>
            </div>

            <Input
              type="password"
              name="password"
              label="Password"
              placeholder="Enter your password"
              value={formData.password}
              onChange={handleChange}
              required
              fullWidth
              autoComplete="new-password"
            />

            <Input
              type="password"
              name="confirmPassword"
              label="Confirm Password"
              placeholder="Confirm your password"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
              fullWidth
              autoComplete="new-password"
            />

            {error && <div className="register__error">{error}</div>}

            <Button
              type="submit"
              fullWidth
              isLoading={loading}
              size="lg"
              style={{ marginTop: 'var(--spacing-lg)' }}
            >
              Create Account
            </Button>
          </form>

          <p className="register__footer">
            Already have an account?{' '}
            <Link to="/login" className="register__link">
              Login here
            </Link>
          </p>
        </Card>
      </div>
    </Layout>
  );
};

