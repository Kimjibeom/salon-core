// Copyright 2026. Kimjibeom. All rights reserved.
import type { Metadata } from 'next';
import './globals.css';
import { AuthProvider } from '@/hooks/useAuth';

export const metadata: Metadata = {
  title: 'Salon Core CRM - 미용실 통합 관리 시스템',
  description: '미용실 및 뷰티숍을 위한 통합 CRM 시스템. 예약, 고객, 매출, 직원 관리를 한 곳에서.',
  keywords: ['미용실', 'CRM', '예약관리', '매출분석', '고객관리', '뷰티숍'],
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ko" className="dark">
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta httpEquiv="X-Content-Type-Options" content="nosniff" />
        <meta httpEquiv="X-Frame-Options" content="DENY" />
      </head>
      <body className="min-h-screen bg-dark-bg">
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}
