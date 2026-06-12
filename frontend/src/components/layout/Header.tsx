// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useAuth } from '@/hooks/useAuth';

interface HeaderProps {
  isConnected: boolean;
}

export default function Header({ isConnected }: HeaderProps) {
  const { staff } = useAuth();
  const now = new Date();

  return (
    <header className="sticky top-0 z-30 bg-dark-bg/80 backdrop-blur-xl border-b border-dark-border px-4 lg:px-6 py-4">
      <div className="flex items-center justify-between">
        {/* Mobile Logo */}
        <div className="lg:hidden flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-salon-400 to-salon-600 flex items-center justify-center text-white font-bold text-sm">
            S
          </div>
          <span className="font-bold text-white">Salon Core</span>
        </div>

        {/* Date & Greeting */}
        <div className="hidden lg:block">
          <p className="text-sm text-dark-muted">
            {now.toLocaleDateString('ko-KR', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' })}
          </p>
          <h2 className="text-lg font-semibold text-white">
            안녕하세요, {staff?.name}님 👋
          </h2>
        </div>

        {/* Status indicators */}
        <div className="flex items-center gap-4">
          {/* WS Connection Status */}
          <div className="flex items-center gap-2">
            <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-400 animate-pulse-soft' : 'bg-red-400'}`} />
            <span className="text-xs text-dark-muted hidden sm:inline">
              {isConnected ? '실시간 연결됨' : '연결 끊김'}
            </span>
          </div>
        </div>
      </div>
    </header>
  );
}
