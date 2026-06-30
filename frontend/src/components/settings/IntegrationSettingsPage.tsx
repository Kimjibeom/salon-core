// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';
import api from '@/lib/api';

interface Setting {
  key: string;
  value: string;
  description: string;
  updated_at: string;
}

export default function IntegrationSettingsPage() {
  const [settings, setSettings] = useState<Setting[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState('');

  const fetchSettings = useCallback(async () => {
    try {
      const data = await api.get<Setting[]>('/api/settings');
      setSettings(data || []);
    } catch (err: any) {
      setMessage('설정을 불러오는데 실패했습니다.');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchSettings();
  }, [fetchSettings]);

  const handleChange = (key: string, value: string) => {
    setSettings(prev => prev.map(s => s.key === key ? { ...s, value } : s));
  };

  const handleSave = async () => {
    setIsSaving(true);
    setMessage('');
    try {
      await api.put('/api/settings', settings);
      setMessage('설정이 성공적으로 저장되었습니다.');
      setTimeout(() => setMessage(''), 3000);
    } catch (err: any) {
      setMessage('설정 저장에 실패했습니다.');
    } finally {
      setIsSaving(false);
    }
  };

  const handleManualSync = async () => {
    try {
      await api.post('/api/sync/naver/manual');
      setMessage('네이버 예약 수동 동기화가 요청되었습니다.');
      setTimeout(() => setMessage(''), 3000);
    } catch (err: any) {
      setMessage('동기화 요청에 실패했습니다.');
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <div className="loading-spinner w-8 h-8" />
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-fade-in max-w-4xl mx-auto">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">⚙️ 외부 서비스 연동 설정</h1>
        <button onClick={handleSave} disabled={isSaving} className="btn-primary">
          {isSaving ? '저장 중...' : '설정 저장'}
        </button>
      </div>

      {message && (
        <div className={`p-4 rounded-lg border ${message.includes('실패') ? 'bg-red-500/20 border-red-500/50 text-red-400' : 'bg-green-500/20 border-green-500/50 text-green-400'}`}>
          {message}
        </div>
      )}

      <div className="glass-card p-6 space-y-6">
        <div>
          <h2 className="text-xl font-bold text-white mb-2">네이버 예약 연동</h2>
          <p className="text-sm text-dark-muted mb-6">
            네이버 스마트플레이스 예약과 연동하기 위한 설정입니다.<br/>
            웹훅을 통해 실시간으로 예약이 동기화됩니다.
          </p>
        </div>

        <div className="space-y-4 max-w-2xl">
          {settings.filter(s => s.key.startsWith('naver_')).map(setting => (
            <div key={setting.key}>
              <label className="block text-sm font-medium text-dark-muted mb-1">
                {setting.description || setting.key}
              </label>
              <input
                type="text"
                className="glass-input w-full"
                value={setting.value}
                onChange={(e) => handleChange(setting.key, e.target.value)}
                placeholder={`${setting.description} 입력`}
              />
            </div>
          ))}
        </div>

        <div className="pt-6 border-t border-dark-border/50">
          <h3 className="text-lg font-bold text-white mb-2">수동 동기화</h3>
          <p className="text-sm text-dark-muted mb-4">
            웹훅 유실 등 예약 누락이 의심될 경우, 수동으로 최근 예약을 동기화할 수 있습니다.
          </p>
          <button onClick={handleManualSync} className="btn-secondary">
            네이버 예약 수동 가져오기
          </button>
        </div>
      </div>
    </div>
  );
}
