// Copyright 2026. Kimjibeom. All rights reserved.

import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('ko-KR', { style: 'currency', currency: 'KRW' }).format(amount);
}

export function formatNumber(num: number): string {
  return new Intl.NumberFormat('ko-KR').format(num);
}

/** Masks a phone number for PII protection: 010-****-1234 */
export function maskPhone(phone: string): string {
  if (phone.length >= 8) {
    return phone.slice(0, 3) + '-****-' + phone.slice(-4);
  }
  return phone;
}

export function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString('ko-KR', { year: 'numeric', month: '2-digit', day: '2-digit' });
}

export function formatTime(timeStr: string): string {
  return timeStr.slice(0, 5); // HH:MM
}

export function getStatusColor(status: string): string {
  switch (status) {
    case 'reserved': return 'bg-blue-500/20 text-blue-400';
    case 'waiting': return 'bg-amber-500/20 text-amber-400';
    case 'in_progress': return 'bg-green-500/20 text-green-400';
    case 'completed': return 'bg-gray-500/20 text-gray-400';
    case 'canceled': return 'bg-red-500/20 text-red-400';
    case 'no_show': return 'bg-red-700/20 text-red-500';
    default: return 'bg-gray-500/20 text-gray-400';
  }
}

export function getStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    reserved: '예약됨',
    waiting: '대기중',
    in_progress: '시술중',
    completed: '완료',
    canceled: '취소됨',
    no_show: '노쇼',
  };
  return labels[status] || status;
}

export function getSourceLabel(source: string): string {
  const labels: Record<string, string> = {
    online: '온라인',
    offline: '방문',
    naver: '네이버',
  };
  return labels[source] || source;
}

export function getCategoryLabel(category: string): string {
  return category === 'service' ? '시술' : '점판';
}

export function getPaymentLabel(method: string): string {
  const labels: Record<string, string> = {
    card: '카드',
    cash: '현금',
    membership: '정액권',
    mixed: '복합',
  };
  return labels[method] || method;
}
