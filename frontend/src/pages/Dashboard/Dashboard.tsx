import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Layout, Card, Button } from '../../components';
import { useAuth } from '../../context/AuthContext';
import { loanApplicationsAPI } from '../../utils/api';
import './Dashboard.css';

interface LoanApplication {
  id: string;
  amount: number;
  purpose: string;
  status: string;
  created_at: string;
}

export const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const [applications, setApplications] = useState<LoanApplication[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (user) {
      fetchApplications();
    }
  }, [user]);

  const fetchApplications = async () => {
    try {
      const response = await loanApplicationsAPI.getAll();
      setApplications(response.data);
    } catch (error) {
      console.error('Error fetching applications:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    const colors: { [key: string]: string } = {
      draft: 'var(--color-neutral-500)',
      submitted: 'var(--color-info-500)',
      under_review: 'var(--color-warning-500)',
      approved: 'var(--color-success-500)',
      rejected: 'var(--color-error-500)',
      partially_approved: 'var(--color-secondary-500)',
    };
    return colors[status] || 'var(--color-neutral-500)';
  };

  if (loading) {
    return (
      <Layout user={user}>
        <div className="dashboard__loading">Loading...</div>
      </Layout>
    );
  }

  return (
    <Layout user={user}>
      <div className="dashboard">
        <div className="dashboard__header">
          <h1 className="dashboard__title">Dashboard</h1>
          <Link to="/applications/new">
            <Button>New Application</Button>
          </Link>
        </div>

        <div className="dashboard__stats">
          <Card variant="elevated" className="dashboard__stat-card">
            <h3 className="dashboard__stat-label">Total Applications</h3>
            <p className="dashboard__stat-value">{applications.length}</p>
          </Card>

          <Card variant="elevated" className="dashboard__stat-card">
            <h3 className="dashboard__stat-label">Under Review</h3>
            <p className="dashboard__stat-value">
              {applications.filter((app) => app.status === 'under_review').length}
            </p>
          </Card>

          <Card variant="elevated" className="dashboard__stat-card">
            <h3 className="dashboard__stat-label">Approved</h3>
            <p className="dashboard__stat-value">
              {applications.filter((app) => app.status === 'approved').length}
            </p>
          </Card>
        </div>

        <div className="dashboard__applications">
          <h2 className="dashboard__section-title">Recent Applications</h2>
          {applications.length === 0 ? (
            <Card variant="outlined" className="dashboard__empty">
              <p>No applications yet. Create your first application!</p>
              <Link to="/applications/new">
                <Button style={{ marginTop: 'var(--spacing-md)' }}>
                  Create Application
                </Button>
              </Link>
            </Card>
          ) : (
            <div className="dashboard__app-list">
              {applications.map((app) => (
                <Card
                  key={app.id}
                  variant="outlined"
                  className="dashboard__app-card"
                  onClick={() => window.location.href = `/applications/${app.id}`}
                >
                  <div className="dashboard__app-header">
                    <h3 className="dashboard__app-amount">
                      ${app.amount.toLocaleString()}
                    </h3>
                    <span
                      className="dashboard__app-status"
                      style={{ color: getStatusColor(app.status) }}
                    >
                      {app.status.replace('_', ' ').toUpperCase()}
                    </span>
                  </div>
                  <p className="dashboard__app-purpose">{app.purpose}</p>
                  <p className="dashboard__app-date">
                    Created: {new Date(app.created_at).toLocaleDateString()}
                  </p>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

