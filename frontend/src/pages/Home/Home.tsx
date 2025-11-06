import React from 'react';
import { Link } from 'react-router-dom';
import { Layout, Button, Card } from '../../components';
import './Home.css';

export const Home: React.FC = () => {
  return (
    <Layout>
      <div className="home">
        <div className="home__hero">
          <h1 className="home__title">Welcome to EasyFund</h1>
          <p className="home__subtitle">
            Your trusted platform for managing loan applications and connecting with banks
          </p>
          <div className="home__actions">
            <Link to="/login">
              <Button size="lg">Get Started</Button>
            </Link>
            <Link to="/register">
              <Button variant="outline" size="lg">
                Create Account
              </Button>
            </Link>
          </div>
        </div>

        <div className="home__features">
          <Card variant="elevated" className="home__feature-card">
            <h3 className="home__feature-title">🏦 Bank Integration</h3>
            <p className="home__feature-text">
              Connect with multiple banks and get the best loan offers tailored to your needs.
            </p>
          </Card>

          <Card variant="elevated" className="home__feature-card">
            <h3 className="home__feature-title">📝 Easy Applications</h3>
            <p className="home__feature-text">
              Submit loan applications with ease and track their status in real-time.
            </p>
          </Card>

          <Card variant="elevated" className="home__feature-card">
            <h3 className="home__feature-title">🔒 Secure & Reliable</h3>
            <p className="home__feature-text">
              Your data is protected with bank-level encryption and security measures.
            </p>
          </Card>
        </div>
      </div>
    </Layout>
  );
};

