import React from 'react';
import { Header } from './Header';
import { Footer } from './Footer';
import './Layout.css';

interface LayoutProps {
  children: React.ReactNode;
  user?: {
    name: string;
    email: string;
    role: string;
  };
  onLogout?: () => void;
}

export const Layout: React.FC<LayoutProps> = ({ children, user, onLogout }) => {
  return (
    <div className="layout">
      <Header user={user} onLogout={onLogout} />
      <main className="layout__main">
        {children}
      </main>
      <Footer />
    </div>
  );
};

