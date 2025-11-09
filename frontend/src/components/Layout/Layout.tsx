import React from 'react';
import { TopBar } from '../TopBar';
import { Footer } from './Footer';
import './Layout.css';

interface LayoutProps {
  children: React.ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="layout">
      <TopBar variant="static" />
      <main className="layout__main">{children}</main>
      <Footer />
    </div>
  );
};

