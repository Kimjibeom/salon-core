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

interface NaverMapping {
  id: string;
  internal_type: string;
  internal_id: string;
  naver_id: string;
  created_at: string;
}

interface StaffItem {
  id: string;
  name: string;
}

interface ServiceItem {
  id: string;
  name: string;
  category: string;
}

export default function IntegrationSettingsPage() {
  const [settings, setSettings] = useState<Setting[]>([]);
  const [mappings, setMappings] = useState<NaverMapping[]>([]);
  const [staffList, setStaffList] = useState<StaffItem[]>([]);
  const [serviceList, setServiceList] = useState<ServiceItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState('');

  // New mapping form
  const [newMappingType, setNewMappingType] = useState('STAFF');
  const [newMappingInternalId, setNewMappingInternalId] = useState('');
  const [newMappingNaverId, setNewMappingNaverId] = useState('');

  const fetchAll = useCallback(async () => {
    try {
      const [settingsData, mappingsData, staffData, servicesData] = await Promise.all([
        api.get<Setting[]>('/api/settings'),
        api.get<NaverMapping[]>('/api/naver-mappings'),
        api.get<StaffItem[]>('/api/staffs'),
        api.get<ServiceItem[]>('/api/services'),
      ]);
      setSettings(settingsData || []);
      setMappings(mappingsData || []);
      setStaffList(staffData || []);
      setServiceList(servicesData || []);
    } catch {
      setMessage('데이터를 불러오는데 실패했습니다.');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchAll();
  }, [fetchAll]);

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
    } catch {
      setMessage('설정 저장에 실패했습니다.');
    } finally {
      setIsSaving(false);
    }
  };

  const handleManualSync = async () => {
    try {
      await api.post('/api/sync/naver/manual', {});
      setMessage('네이버 예약 수동 동기화가 요청되었습니다.');
      setTimeout(() => setMessage(''), 3000);
    } catch {
      setMessage('동기화 요청에 실패했습니다.');
    }
  };

  const handleAddMapping = async () => {
    if (!newMappingInternalId || !newMappingNaverId) {
      setMessage('내부 항목과 네이버 ID를 모두 입력해주세요.');
      return;
    }
    try {
      await api.post('/api/naver-mappings', {
        internal_type: newMappingType,
        internal_id: newMappingInternalId,
        naver_id: newMappingNaverId,
      });
      setNewMappingInternalId('');
      setNewMappingNaverId('');
      setMessage('매핑이 추가되었습니다.');
      setTimeout(() => setMessage(''), 3000);
      // Refresh mappings
      const updated = await api.get<NaverMapping[]>('/api/naver-mappings');
      setMappings(updated || []);
    } catch {
      setMessage('매핑 추가에 실패했습니다.');
    }
  };

  const handleDeleteMapping = async (id: string) => {
    try {
      await api.delete(`/api/naver-mappings/${id}`);
      setMappings(prev => prev.filter(m => m.id !== id));
      setMessage('매핑이 삭제되었습니다.');
      setTimeout(() => setMessage(''), 3000);
    } catch {
      setMessage('매핑 삭제에 실패했습니다.');
    }
  };

  const getInternalName = (type: string, id: string) => {
    if (type === 'STAFF') {
      return staffList.find(s => s.id === id)?.name || id.substring(0, 8);
    }
    const svc = serviceList.find(s => s.id === id);
    return svc ? `${svc.category} - ${svc.name}` : id.substring(0, 8);
  };

  const backendUrl = typeof window !== 'undefined' ? (process.env.NEXT_PUBLIC_BACKEND_URL || window.location.origin) : '';

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

      {/* Shop Basic Info (shown on the public booking page) */}
      <div className="glass-card p-6 space-y-6">
        <div>
          <h2 className="text-xl font-bold text-white mb-2">🏠 매장 기본 정보</h2>
          <p className="text-sm text-dark-muted">
            매장명, 연락처, 영업시간 등의 정보입니다. 자체 예약 사이트 상단에 노출되며,<br/>
            영업시간과 예약 슬롯 간격은 예약 가능 시간 계산에 사용됩니다.
          </p>
        </div>
        <div className="space-y-4 max-w-2xl">
          {settings.filter(s => s.key.startsWith('shop_') || s.key === 'booking_slot_interval').map(setting => (
            <div key={setting.key}>
              <label className="block text-sm font-medium text-dark-muted mb-1">
                {setting.description || setting.key}
              </label>
              <input
                type={setting.key.endsWith('_time') ? 'time' : setting.key === 'booking_slot_interval' ? 'number' : 'text'}
                className="glass-input w-full"
                value={setting.value}
                onChange={(e) => handleChange(setting.key, e.target.value)}
                placeholder={`${setting.description || setting.key} 입력`}
              />
            </div>
          ))}
        </div>
      </div>

      {/* Webhook URL Info */}
      <div className="glass-card p-6">
        <h2 className="text-xl font-bold text-white mb-2">📡 웹훅 URL</h2>
        <p className="text-sm text-dark-muted mb-3">
          아래 URL을 네이버 스마트플레이스 웹훅 설정에 입력하세요.
        </p>
        <div className="bg-dark-bg/50 rounded-lg p-3 font-mono text-sm text-green-400 break-all">
          {backendUrl}/api/sync/naver/webhook
        </div>
      </div>

      {/* Naver Settings */}
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

      {/* Naver ID Mapping */}
      <div className="glass-card p-6 space-y-6">
        <div>
          <h2 className="text-xl font-bold text-white mb-2">🔗 네이버 ID 매핑</h2>
          <p className="text-sm text-dark-muted mb-4">
            내부 시스템의 디자이너/시술과 네이버 시스템의 ID를 1:1 매핑합니다.<br/>
            이 매핑이 있어야 웹훅으로 수신된 예약이 올바르게 연결됩니다.
          </p>
        </div>

        {/* Add Mapping Form */}
        <div className="bg-dark-bg/30 rounded-lg p-4 space-y-3">
          <h4 className="text-sm font-semibold text-white">새 매핑 추가</h4>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-3 items-end">
            <div>
              <label className="block text-xs text-dark-muted mb-1">유형</label>
              <select
                className="glass-input w-full"
                value={newMappingType}
                onChange={(e) => { setNewMappingType(e.target.value); setNewMappingInternalId(''); }}
              >
                <option value="STAFF">디자이너</option>
                <option value="SERVICE">시술</option>
              </select>
            </div>
            <div>
              <label className="block text-xs text-dark-muted mb-1">내부 항목</label>
              <select
                className="glass-input w-full"
                value={newMappingInternalId}
                onChange={(e) => setNewMappingInternalId(e.target.value)}
              >
                <option value="">선택...</option>
                {newMappingType === 'STAFF'
                  ? staffList.map(s => <option key={s.id} value={s.id}>{s.name}</option>)
                  : serviceList.map(s => <option key={s.id} value={s.id}>{s.category} - {s.name}</option>)
                }
              </select>
            </div>
            <div>
              <label className="block text-xs text-dark-muted mb-1">네이버 ID</label>
              <input
                type="text"
                className="glass-input w-full"
                value={newMappingNaverId}
                onChange={(e) => setNewMappingNaverId(e.target.value)}
                placeholder="네이버 측 ID"
              />
            </div>
            <button onClick={handleAddMapping} className="btn-primary h-10">
              추가
            </button>
          </div>
        </div>

        {/* Existing Mappings */}
        {mappings.length > 0 ? (
          <div className="space-y-2">
            {mappings.map(m => (
              <div key={m.id} className="flex items-center justify-between bg-dark-bg/20 rounded-lg p-3">
                <div className="flex items-center gap-3">
                  <span className={`badge ${m.internal_type === 'STAFF' ? 'bg-blue-500/20 text-blue-400' : 'bg-purple-500/20 text-purple-400'}`}>
                    {m.internal_type === 'STAFF' ? '디자이너' : '시술'}
                  </span>
                  <span className="text-white font-medium">{getInternalName(m.internal_type, m.internal_id)}</span>
                  <span className="text-dark-muted">→</span>
                  <span className="text-green-400 font-mono text-sm">{m.naver_id}</span>
                </div>
                <button
                  onClick={() => handleDeleteMapping(m.id)}
                  className="text-red-400 hover:text-red-300 text-sm"
                >
                  삭제
                </button>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-sm text-dark-muted text-center py-4">
            등록된 매핑이 없습니다. 위 폼에서 매핑을 추가해주세요.
          </p>
        )}
      </div>

      {/* Public Booking Page Link */}
      <div className="glass-card p-6">
        <h2 className="text-xl font-bold text-white mb-2">🌐 자체 예약 페이지</h2>
        <p className="text-sm text-dark-muted mb-3">
          고객이 직접 예약할 수 있는 공개 페이지입니다.
        </p>
        <div className="bg-dark-bg/50 rounded-lg p-3 font-mono text-sm text-blue-400 break-all">
          {typeof window !== 'undefined' ? window.location.origin : ''}/booking
        </div>
      </div>
    </div>
  );
}
