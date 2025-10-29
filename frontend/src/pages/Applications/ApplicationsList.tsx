import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Layout, Card, Button, Input } from '../../components';
import { useAuth } from '../../context/AuthContext';
import { loanApplicationsAPI } from '../../utils/api';
import './ApplicationsList.css';

interface LoanApplication {
  id: string;
  amount: number;
  purpose: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export const ApplicationsList: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [applications, setApplications] = useState<LoanApplication[]>([]);
  const [filteredApplications, setFilteredApplications] = useState<LoanApplication[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');

  useEffect(() => {
    fetchApplications();
  }, []);

  useEffect(() => {
    filterApplications();
  }, [applications, searchTerm, filterStatus]);

  const fetchApplications = async () => {
    try {
      const response = await loanApplicationsAPI.getAll();
      setApplications(response.data);
      setFilteredApplications(response.data);
    } catch (error) {
      console.error('Error fetching applications:', error);
    } finally {
      setLoading(false);
    }
  };

  const filterApplications = () => {
    let filtered = applications;

    if (searchTerm) {
      filtered = filtered.filter(
        (app) =>
          app.purpose.toLowerCase().includes(searchTerm.toLowerCase()) ||
          app.amount.toString().includes(searchTerm)
      );
    }

    if (filterStatus !== 'all') {
      filtered = filtered.filter((app) => app.status === filterStatus);
    }

    setFilteredApplications(filtered);
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
        <div className="applications__loading">Loading...</div>
      </Layout>
    );
  }

  return (
    <Layout user={user}>
      <div className="applications">
        <div className="applications__header">
          <h1 className="applications__title">Loan Applications</h1>
          <Button onClick={() => navigate('/applications/new')}>
            New Application
          </Button>
        </div>

        <div className="applications__filters">
          <Input
            type="search"
            placeholder="Search applications..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="applications__search"
          />
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
            className="input"
            style={{ maxWidth: '200px' }}
          >
            <option value="all">All Status</option>
            <option value="draft">Draft</option>
            <option value="submitted">Submitted</option>
            <option value="under_review">Under Review</option>
            <option value="approved">Approved</option>
            <option value="rejected">Rejected</option>
          </select>
        </div>

        <div className="applications__list">
          {filteredApplications.length === 0 ? (
            <Card variant="outlined" className="applications__empty">
              <p>No applications found.</p>
            </Card>
          ) : (
            filteredApplications.map((app) => (
              <Card
                key={app.id}
                variant="outlined"
                className="applications__item"
                onClick={() => navigate(`/applications/${app.id}`)}
              >
                <div className="applications__item-header">
                  <div>
                    <h3 className="applications__item-amount">
                      ${app.amount.toLocaleString()}
                    </h3>
                    <p className="applications__item-purpose">{app.purpose}</p>
                  </div>
                  <span
                    className="applications__item-status"
                    style={{ color: getStatusColor(app.status) }}
                  >
                    {app.status.replace('_', ' ').toUpperCase()}
                  </span>
                </div>
                <div className="applications__item-footer">
                  <span className="applications__item-date">
                    Created: {new Date(app.created_at).toLocaleDateString()}
                  </span>
                  {app.updated_at !== app.created_at && (
                    <span className="applications__item-date">
                      Updated: {new Date(app.updated_at).toLocaleDateString()}
                    </span>
                  )}
                </div>
              </Card>
            ))
          )}
        </div>
      </div>
    </Layout>
  );
};

