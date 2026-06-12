// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useAuth } from '@/hooks/useAuth';

type TabKey = 'reservations' | 'customers' | 'pos' | 'analytics' | 'staff';

const tabs: { key: TabKey; label: string; icon: string; roles?: string[] }[] = [
  { key: 'reservations', label: '예약 관리', icon: '📅' },
  { key: 'customers', label: '고객 관리', icon: '👤' },
  { key: 'pos', label: '간편 POS', icon: '💳' },
  { key: 'analytics', label: '매출 분석', icon: '📊', roles: ['admin', 'designer'] },
  { key: 'staff', label: '직원 관리', icon: '👥', roles: ['admin'] },
];

interface SidebarProps {
  activeTab: TabKey;
  onTabChange: (tab: TabKey) => void;
}

export default function Sidebar({ activeTab, onTabChange }: SidebarProps) {
  const { staff, logout } = useAuth();

  return (
    <>
      {/* Desktop Sidebar */}
      <aside className="hidden lg:flex w-64 fixed inset-y-0 left-0 z-40 flex-col bg-dark-card/95 backdrop-blur-xl border-r border-dark-border">
        {/* Logo */}
        <div className="p-6 border-b border-dark-border">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-salon-400 to-salon-600 flex items-center justify-center text-white font-bold text-lg shadow-lg shadow-salon-500/30">
              S
            </div>
            <div>
              <h1 className="text-lg font-bold text-white">Salon Core</h1>
              <p className="text-xs text-dark-muted">통합 CRM 시스템</p>
            </div>
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-4 space-y-1">
          {tabs.map((tab) => {
            if (tab.roles && staff && !tab.roles.includes(staff.role)) return null;
            return (
              <button
                key={tab.key}
                id={`sidebar-tab-${tab.key}`}
                onClick={() => onTabChange(tab.key)}
                className={`w-full sidebar-link ${activeTab === tab.key ? 'sidebar-link-active' : ''}`}
              >
                <span className="text-xl">{tab.icon}</span>
                <span>{tab.label}</span>
              </button>
            );
          })}
        </nav>

        {/* User Info */}
        <div className="p-4 border-t border-dark-border">
          <div className="glass-card p-3">
            <div className="flex items-center gap-3">
              <div className="w-9 h-9 rounded-full bg-gradient-to-br from-salon-400 to-purple-600 flex items-center justify-center text-white text-sm font-bold">
                {staff?.name?.charAt(0) || '?'}
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-white truncate">{staff?.name}</p>
                <p className="text-xs text-dark-muted">{staff?.role === 'admin' ? '관리자' : staff?.role === 'designer' ? '디자이너' : '스탭'}</p>
              </div>
            </div>
            <button
              id="logout-button"
              onClick={logout}
              className="mt-3 w-full text-xs text-dark-muted hover:text-red-400 transition-colors py-1.5 rounded-lg hover:bg-red-500/10"
            >
              로그아웃
            </button>
          </div>
        </div>
      </aside>

      {/* Mobile Bottom Nav */}
      <nav className="lg:hidden fixed bottom-0 inset-x-0 z-40 bg-dark-card/95 backdrop-blur-xl border-t border-dark-border">
        <div className="flex justify-around py-2">
          {tabs.map((tab) => {
            if (tab.roles && staff && !tab.roles.includes(staff.role)) return null;
            return (
              <button
                key={tab.key}
                id={`mobile-tab-${tab.key}`}
                onClick={() => onTabChange(tab.key)}
                className={`flex flex-col items-center gap-0.5 px-3 py-1.5 rounded-lg transition-all ${
                  activeTab === tab.key ? 'text-salon-400' : 'text-dark-muted'
                }`}
              >
                <span className="text-lg">{tab.icon}</span>
                <span className="text-[10px] font-medium">{tab.label}</span>
              </button>
            );
          })}
        </div>
      </nav>
    </>
  );
}
