import React, { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { Layout, Card, Button } from '../../components';
import { useAuth } from '../../context/AuthContext';
import { loanApplicationsAPI, banksAPI } from '../../utils/api';
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
  const [payments, setPayments] = useState<{ date: string; amount: number; bank: string }[]>([]);
  const [debtsByBank, setDebtsByBank] = useState<{ bank: string; amount: number }[]>([]);
  const [overall, setOverall] = useState<{ totalDebt: number; totalPaid: number }>({ totalDebt: 0, totalPaid: 0 });

  useEffect(() => {
    if (user) {
      fetchApplications();
      // NOTE: Backend endpoints for these analytics aren't present yet.
      // Using lightweight mocked data so UI is functional.
      seedDemoData();
      fetchBanksForDemo();
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

  const fetchBanksForDemo = async () => {
    try {
      await banksAPI.getAll();
      // If API exists, you could map real balances here. For now keep demo mix.
    } catch {
      // ignore
    }
  };

  const seedDemoData = () => {
    const now = new Date();
    const months = Array.from({ length: 12 }).map((_, i) => {
      const d = new Date(now.getFullYear(), now.getMonth() - (11 - i), 1);
      return d.toISOString().slice(0, 10);
    });
    const series = months.map((d, idx) => ({
      date: d,
      amount: 120 + Math.round(40 * Math.sin(idx / 2) + Math.random() * 30),
      bank: ['Sberbank', 'VTB', 'Alfa-Bank'][idx % 3],
    }));
    setPayments(series);
    setDebtsByBank([
      { bank: 'Sberbank', amount: 650000 },
      { bank: 'VTB', amount: 420000 },
      { bank: 'T-Bank', amount: 180000 },
      { bank: 'Alfa-Bank', amount: 90000 },
    ]);
    setOverall({ totalDebt: 1314593, totalPaid: 1314593 - 2314593 + 1000000 });
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

  const payoffPercent = useMemo(() => {
    const { totalDebt, totalPaid } = overall;
    if (totalDebt <= 0) return 0;
    return Math.max(0, Math.min(100, Math.round((totalPaid / totalDebt) * 100)));
  }, [overall]);

  const totalDebtAmount = useMemo(() => debtsByBank.reduce((s, d) => s + d.amount, 0), [debtsByBank]);

  const donutArcs = useMemo(() => {
    const r = 56;
    const c = 2 * Math.PI * r;
    let acc = 0;
    return debtsByBank.map((d) => {
      const frac = totalDebtAmount ? d.amount / totalDebtAmount : 0;
      const len = frac * c;
      const dasharray = `${len} ${c - len}`;
      const dashoffset = c - acc * c;
      acc += frac;
      return { bank: d.bank, dasharray, dashoffset };
    });
  }, [debtsByBank, totalDebtAmount]);

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

        {/* Top KPIs */}
        <div className="dashboard__stats">
          <Card variant="elevated" className="dashboard__stat-card">
            <h3 className="dashboard__stat-label">Overall Loan</h3>
            <p className="dashboard__stat-value dashboard__stat-value--debt">
              {totalDebtAmount.toLocaleString('ru-RU')} ₽
            </p>
          </Card>
          <Card variant="elevated" className="dashboard__stat-card">
            <h3 className="dashboard__stat-label">Paid</h3>
            <p className="dashboard__stat-value">
              {(overall.totalPaid || Math.round(totalDebtAmount * payoffPercent / 100)).toLocaleString('ru-RU')} ₽
            </p>
          </Card>
          <Card variant="elevated" className="dashboard__stat-card">
            <h3 className="dashboard__stat-label">Applications</h3>
            <p className="dashboard__stat-value">{applications.length}</p>
          </Card>
        </div>

        {/* Main grid */}
        <div className="dashboard__grid">
          {/* 1) Progress bar */}
          <Card variant="outlined" className="dashboard__panel">
            <h3 className="dashboard__panel-title">You are almost there!</h3>
            <div className="dashboard__progress-wrap">
              <div className="dashboard__progress">
                <div className="dashboard__progress-bar" style={{ width: `${payoffPercent}%` }} />
              </div>
              <div className="dashboard__progress-meta">
                <span>{payoffPercent}%</span>
                <span>
                  {totalDebtAmount.toLocaleString('ru-RU')} ₽ →{' '}
                  {(totalDebtAmount - (overall.totalPaid || totalDebtAmount * payoffPercent / 100)).toLocaleString('ru-RU')} ₽ left
                </span>
              </div>
            </div>
          </Card>

          {/* 2) History list */}
          <Card variant="outlined" className="dashboard__panel">
            <h3 className="dashboard__panel-title">Payment history</h3>
            <ul className="dashboard__history">
              {payments.slice(-7).reverse().map((p, i) => (
                <li key={i} className="dashboard__history-item">
                  <span className="dashboard__history-icon" />
                  <div className="dashboard__history-info">
                    <div className="dashboard__history-title">Payment to {p.bank}</div>
                    <div className="dashboard__history-sub">
                      {new Date(p.date).toLocaleDateString('ru-RU')}
                    </div>
                  </div>
                  <div className="dashboard__history-amount">-{p.amount.toLocaleString('ru-RU')} ₽</div>
                </li>
              ))}
            </ul>
          </Card>

          {/* 3) Donut chart */}
          <Card variant="outlined" className="dashboard__panel">
            <h3 className="dashboard__panel-title">Debt by bank</h3>
            <div className="dashboard__donut">
              <svg width="140" height="140" viewBox="0 0 140 140">
                <g transform="translate(70,70)">
                  <circle r="56" fill="none" stroke="var(--color-background-tertiary)" strokeWidth="16" />
                  {donutArcs.map((arc, idx) => (
                    <circle
                      key={idx}
                      r="56"
                      fill="none"
                      stroke={`var(--chart-${(idx % 6) + 1})`}
                      strokeWidth="16"
                      strokeDasharray={arc.dasharray}
                      strokeDashoffset={arc.dashoffset}
                      transform="rotate(-90)"
                      strokeLinecap="round"
                    />
                  ))}
                </g>
              </svg>
              <div className="dashboard__legend">
                {debtsByBank.map((d, idx) => (
                  <div key={d.bank} className="dashboard__legend-row">
                    <span className="dashboard__legend-dot" style={{ background: `var(--chart-${(idx % 6) + 1})` }} />
                    <span className="dashboard__legend-label">{d.bank}</span>
                    <span className="dashboard__legend-value">{d.amount.toLocaleString('ru-RU')} ₽</span>
                  </div>
                ))}
              </div>
            </div>
          </Card>

          {/* 5) Line chart */}
          <Card variant="outlined" className="dashboard__panel">
            <h3 className="dashboard__panel-title">Monthly payments</h3>
            <svg viewBox="0 0 500 160" className="dashboard__linechart">
              <polyline
                fill="url(#fill)"
                stroke="var(--color-info-500)"
                strokeWidth="2"
                points={payments
                  .map((p, i) => {
                    const x = (i / Math.max(1, payments.length - 1)) * 500;
                    const y = 140 - (p.amount / 200) * 120;
                    return `${x},${Math.max(0, Math.min(140, y))}`;
                  })
                  .join(' ')}
              />
              <defs>
                <linearGradient id="fill" x1="0" x2="0" y1="0" y2="1">
                  <stop offset="0%" stopColor="var(--color-info-200)" stopOpacity="0.6" />
                  <stop offset="100%" stopColor="var(--color-background)" stopOpacity="0" />
                </linearGradient>
              </defs>
            </svg>
          </Card>

          {/* 6) Scoring */}
          <Card variant="outlined" className="dashboard__panel">
            <h3 className="dashboard__panel-title">Predicted scoring</h3>
            <div className="dashboard__score">
              <div className="dashboard__score-value">742</div>
              <div className="dashboard__score-meter">
                <div className="dashboard__score-good" style={{ width: '70%' }} />
              </div>
              <div className="dashboard__score-scale">
                <span>Poor</span><span>Fair</span><span>Good</span><span>Excellent</span>
              </div>
            </div>
          </Card>

          {/* 7) Suggestions */}
          <Card variant="outlined" className="dashboard__panel">
            <h3 className="dashboard__panel-title">New offers</h3>
            <ul className="dashboard__offers">
              {['Credit card', 'Refinance loan', 'Cashback card', 'Low-rate loan'].map((t, i) => (
                <li key={i} className="dashboard__offer-row">
                  <div className="dashboard__offer-icon" />
                  <div className="dashboard__offer-main">
                    <div className="dashboard__offer-title">{t}</div>
                    <div className="dashboard__offer-sub">Limit up to {(1000 + i * 250).toLocaleString('ru-RU')} ₽</div>
                  </div>
                  <Button size="sm">Details</Button>
                </li>
              ))}
            </ul>
          </Card>
        </div>

        {/* Recent applications keep at bottom */}
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
                  onClick={() => (window.location.href = `/applications/${app.id}`)}
                >
                  <div className="dashboard__app-header">
                    <h3 className="dashboard__app-amount">${app.amount.toLocaleString()}</h3>
                    <span className="dashboard__app-status" style={{ color: getStatusColor(app.status) }}>
                      {app.status.replace('_', ' ').toUpperCase()}
                    </span>
                  </div>
                  <p className="dashboard__app-purpose">{app.purpose}</p>
                  <p className="dashboard__app-date">Created: {new Date(app.created_at).toLocaleDateString()}</p>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

