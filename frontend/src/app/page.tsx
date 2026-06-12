// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useCallback } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { useWebSocket } from '@/hooks/useWebSocket';
import type { WebSocketEvent } from '@/types';
import Sidebar from '@/components/layout/Sidebar';
import Header from '@/components/layout/Header';
import LoginPage from '@/components/auth/LoginPage';
import ReservationsPage from '@/components/reservations/ReservationsPage';
import CustomersPage from '@/components/customers/CustomersPage';
import POSPage from '@/components/pos/POSPage';
import AnalyticsPage from '@/components/analytics/AnalyticsPage';
import StaffPage from '@/components/staff/StaffPage';
import NotificationPopup from '@/components/ui/NotificationPopup';

type TabKey = 'reservations' | 'customers' | 'pos' | 'analytics' | 'staff';

export default function Home() {
  const { staff, isLoading } = useAuth();
  const [activeTab, setActiveTab] = useState<TabKey>('reservations');
  const [notifications, setNotifications] = useState<WebSocketEvent[]>([]);

  const handleWebSocketEvent = useCallback((event: WebSocketEvent) => {
    setNotifications((prev) => [event, ...prev].slice(0, 5));
    // Auto-dismiss after 8 seconds
    setTimeout(() => {
      setNotifications((prev) => prev.filter((n) => n !== event));
    }, 8000);
  }, []);

  const { isConnected } = useWebSocket(staff ? handleWebSocketEvent : undefined);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-dark-bg">
        <div className="flex flex-col items-center gap-4">
          <div className="loading-spinner w-12 h-12" />
          <p className="text-dark-muted text-sm">로딩 중...</p>
        </div>
      </div>
    );
  }

  if (!staff) {
    return <LoginPage />;
  }

  const renderTab = () => {
    switch (activeTab) {
      case 'reservations': return <ReservationsPage />;
      case 'customers': return <CustomersPage />;
      case 'pos': return <POSPage />;
      case 'analytics': return <AnalyticsPage />;
      case 'staff': return <StaffPage />;
    }
  };

  return (
    <div className="min-h-screen flex bg-dark-bg">
      <Sidebar activeTab={activeTab} onTabChange={setActiveTab} />
      <div className="flex-1 flex flex-col min-h-screen lg:ml-64">
        <Header isConnected={isConnected} />
        <main className="flex-1 p-4 lg:p-6 overflow-y-auto">
          {renderTab()}
        </main>
      </div>
      {/* Notification Popups */}
      <div className="fixed top-4 right-4 z-[100] flex flex-col gap-2 max-w-sm">
        {notifications.map((n, i) => (
          <NotificationPopup key={`${n.time}-${i}`} event={n} onDismiss={() => {
            setNotifications((prev) => prev.filter((_, idx) => idx !== i));
          }} />
        ))}
      </div>
    </div>
  );
}
