// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import type { WebSocketEvent } from '@/types';

interface NotificationPopupProps {
  event: WebSocketEvent;
  onDismiss: () => void;
}

const eventLabels: Record<string, { title: string; icon: string; color: string }> = {
  NEW_RESERVATION: { title: '새 예약', icon: '📅', color: 'border-blue-500/50' },
  RESERVATION_UPDATE: { title: '예약 변경', icon: '🔄', color: 'border-amber-500/50' },
  WAITING_QUEUE_UPDATE: { title: '대기열 변경', icon: '⏳', color: 'border-green-500/50' },
  CID_INCOMING: { title: '전화 수신', icon: '📞', color: 'border-salon-500/50' },
  ONLINE_BOOKING: { title: '온라인 예약', icon: '🌐', color: 'border-purple-500/50' },
  NOTIFICATION: { title: '알림', icon: '🔔', color: 'border-dark-border' },
};

export default function NotificationPopup({ event, onDismiss }: NotificationPopupProps) {
  const config = eventLabels[event.type] || eventLabels.NOTIFICATION;

  return (
    <div className={`glass-card p-4 border-l-4 ${config.color} animate-slide-in-right shadow-2xl cursor-pointer`}
         onClick={onDismiss} role="alert">
      <div className="flex items-start gap-3">
        <span className="text-2xl">{config.icon}</span>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-semibold text-white">{config.title}</p>
          <p className="text-xs text-dark-muted mt-0.5 truncate">
            {new Date(event.time).toLocaleTimeString('ko-KR')}
          </p>
        </div>
        <button onClick={(e) => { e.stopPropagation(); onDismiss(); }}
                className="text-dark-muted hover:text-white text-lg leading-none">×</button>
      </div>
    </div>
  );
}
