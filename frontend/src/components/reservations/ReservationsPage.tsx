// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';
import api from '@/lib/api';
import { formatTime, getStatusColor, getStatusLabel, getSourceLabel } from '@/lib/utils';
import type { Reservation, WaitingQueueEntry, Staff, ServiceMenu } from '@/types';

type ViewMode = 'day' | 'week' | 'month';

export default function ReservationsPage() {
  const [viewMode, setViewMode] = useState<ViewMode>('day');
  const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0]);
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [waitingQueue, setWaitingQueue] = useState<WaitingQueueEntry[]>([]);
  const [staffList, setStaffList] = useState<Staff[]>([]);
  const [serviceList, setServiceList] = useState<ServiceMenu[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchReservations = useCallback(async () => {
    setIsLoading(true);
    try {
      const data = await api.get<Reservation[]>(`/api/reservations?date=${selectedDate}`);
      setReservations(data || []);
    } catch { setReservations([]); }
    setIsLoading(false);
  }, [selectedDate]);

  const fetchWaitingQueue = useCallback(async () => {
    try {
      const data = await api.get<WaitingQueueEntry[]>('/api/reservations/waiting');
      setWaitingQueue(data || []);
    } catch { setWaitingQueue([]); }
  }, []);

  const fetchStaff = useCallback(async () => {
    try {
      const data = await api.get<Staff[]>('/api/staffs');
      setStaffList(data || []);
    } catch { setStaffList([]); }
  }, []);

  const fetchServices = useCallback(async () => {
    try {
      const data = await api.get<ServiceMenu[]>('/api/services');
      setServiceList(data || []);
    } catch { setServiceList([]); }
  }, []);

  useEffect(() => { fetchReservations(); fetchWaitingQueue(); fetchStaff(); fetchServices(); }, [fetchReservations, fetchWaitingQueue, fetchStaff, fetchServices]);

  const handleCreateReservation = async (formData: Record<string, string>) => {
    try {
      await api.post('/api/reservations', {
        ...formData,
        source: 'offline',
      });
      setShowModal(false);
      fetchReservations();
    } catch { /* error handled by API client */ }
  };

  const handleStatusChange = async (id: string, status: string) => {
    try {
      await api.patch(`/api/reservations/${id}/status`, { status });
      fetchReservations();
      fetchWaitingQueue();
    } catch { /* error handled */ }
  };

  const handleAddWalkin = async (formData: Record<string, string>) => {
    try {
      await api.post('/api/reservations/waiting', {
        ...formData,
        source: 'offline',
      });
      fetchWaitingQueue();
    } catch { /* error handled */ }
  };

  const hours = Array.from({ length: 14 }, (_, i) => i + 9); // 9:00 ~ 22:00

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <h1 className="text-2xl font-bold text-white">📅 예약 및 대기 관리</h1>
        <div className="flex items-center gap-3">
          {/* View mode toggle */}
          <div className="flex bg-dark-surface rounded-xl p-1">
            {(['day', 'week', 'month'] as ViewMode[]).map((mode) => (
              <button key={mode} id={`view-mode-${mode}`}
                onClick={() => setViewMode(mode)}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${viewMode === mode ? 'bg-salon-500/20 text-salon-400' : 'text-dark-muted hover:text-white'}`}>
                {mode === 'day' ? '일' : mode === 'week' ? '주' : '월'}
              </button>
            ))}
          </div>
          <input type="date" id="reservation-date-picker" value={selectedDate}
            onChange={(e) => setSelectedDate(e.target.value)}
            className="glass-input text-sm py-2" />
          <button id="new-reservation-btn" onClick={() => setShowModal(true)} className="btn-primary text-sm py-2">
            + 새 예약
          </button>
        </div>
      </div>

      <div className="grid lg:grid-cols-3 gap-6">
        {/* Calendar / Schedule View */}
        <div className="lg:col-span-2 glass-card p-4 overflow-hidden">
          <h3 className="text-lg font-semibold text-white mb-4">
            {new Date(selectedDate).toLocaleDateString('ko-KR', { month: 'long', day: 'numeric', weekday: 'long' })}
          </h3>

          {isLoading ? (
            <div className="flex items-center justify-center py-20">
              <div className="loading-spinner w-8 h-8" />
            </div>
          ) : (
            <div className="space-y-1 max-h-[600px] overflow-y-auto">
              {hours.map((hour) => {
                const hourReservations = reservations.filter((r) => {
                  const startHour = parseInt(r.start_time.split(':')[0]);
                  return startHour === hour;
                });

                return (
                  <div key={hour} className="flex gap-3 min-h-[60px]">
                    <div className="w-16 text-right text-sm text-dark-muted font-mono pt-2 flex-shrink-0">
                      {String(hour).padStart(2, '0')}:00
                    </div>
                    <div className="flex-1 border-l border-dark-border/30 pl-3 space-y-1 py-1">
                      {hourReservations.length > 0 ? hourReservations.map((r) => (
                        <div key={r.id}
                          className="glass-card p-3 cursor-pointer hover:border-salon-500/30 transition-all group">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                              <span className="text-xs font-mono text-dark-muted">
                                {formatTime(r.start_time)} - {formatTime(r.end_time)}
                              </span>
                              <span className={`badge ${getStatusColor(r.status)}`}>
                                {getStatusLabel(r.status)}
                              </span>
                              {r.source !== 'offline' && (
                                <span className="badge bg-purple-500/20 text-purple-400">
                                  {getSourceLabel(r.source)}
                                </span>
                              )}
                            </div>
                            <div className="hidden group-hover:flex gap-1">
                              {r.status === 'reserved' && (
                                <button onClick={() => handleStatusChange(r.id, 'in_progress')}
                                  className="text-xs btn-ghost text-green-400">시술 시작</button>
                              )}
                              {r.status === 'in_progress' && (
                                <button onClick={() => handleStatusChange(r.id, 'completed')}
                                  className="text-xs btn-ghost text-blue-400">완료</button>
                              )}
                              {r.status !== 'canceled' && r.status !== 'completed' && (
                                <button onClick={() => handleStatusChange(r.id, 'canceled')}
                                  className="text-xs btn-ghost text-red-400">취소</button>
                              )}
                            </div>
                          </div>
                          <div className="mt-1.5 flex items-center gap-3">
                            <span className="text-sm font-medium text-white">{r.customer_name}</span>
                            {r.staff_name && <span className="text-xs text-dark-muted">담당: {r.staff_name}</span>}
                            {r.treatment_name && <span className="text-xs text-salon-400">{r.treatment_name}</span>}
                          </div>
                        </div>
                      )) : (
                        <div className="h-full min-h-[40px] border border-dashed border-dark-border/20 rounded-lg" />
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>

        {/* Waiting Queue */}
        <div className="glass-card p-4">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-white">⏳ 대기열</h3>
            <span className="badge bg-amber-500/20 text-amber-400">{waitingQueue.length}명</span>
          </div>

          <div className="space-y-3 max-h-[500px] overflow-y-auto">
            {waitingQueue.length === 0 ? (
              <p className="text-dark-muted text-sm text-center py-8">대기 중인 고객이 없습니다</p>
            ) : waitingQueue.map((entry) => (
              <div key={entry.id} className="glass-card p-3 border-l-4 border-amber-500/50">
                <div className="flex items-center justify-between">
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="text-lg font-bold text-amber-400">#{entry.position}</span>
                      <span className="text-sm font-medium text-white">{entry.customer_name}</span>
                    </div>
                    <p className="text-xs text-dark-muted mt-1">대기시간: {entry.wait_time_minutes}분</p>
                    {entry.treatment_name && <p className="text-xs text-salon-400">{entry.treatment_name}</p>}
                  </div>
                  <div className="flex flex-col gap-1">
                    <button onClick={() => handleStatusChange(entry.id, 'in_progress')}
                      className="text-xs btn-ghost text-green-400">시술 시작</button>
                    <button onClick={() => handleStatusChange(entry.id, 'canceled')}
                      className="text-xs btn-ghost text-red-400">취소</button>
                  </div>
                </div>
              </div>
            ))}
          </div>

          <button id="add-walkin-btn"
            onClick={() => {
              const name = 'Walk-in';
              handleAddWalkin({ customer_name: name, customer_phone: '', start_time: '00:00', end_time: '00:00' });
            }}
            className="btn-secondary w-full mt-4 text-sm">
            + 워크인 고객 추가
          </button>
        </div>
      </div>

      {/* Reservation Modal */}
      {showModal && (
        <ReservationModal
          staffList={staffList}
          serviceList={serviceList}
          onClose={() => setShowModal(false)}
          onSubmit={handleCreateReservation}
        />
      )}
    </div>
  );
}

function ReservationModal({ staffList, serviceList, onClose, onSubmit }: {
  staffList: Staff[];
  serviceList: ServiceMenu[];
  onClose: () => void;
  onSubmit: (data: Record<string, string>) => void;
}) {
  const [form, setForm] = useState({
    customer_name: '',
    customer_phone: '',
    service_id: '',
    treatment_name: '',
    staff_id: '',
    date: new Date().toISOString().split('T')[0],
    start_time: '10:00',
    end_time: '11:00',
    memo: '',
  });

  const handleServiceChange = (id: string) => {
    const service = serviceList.find(s => s.id === id);
    if (!service) {
      setForm(prev => ({ ...prev, service_id: '', treatment_name: '' }));
      return;
    }
    
    // Calculate end time
    const [hours, minutes] = form.start_time.split(':').map(Number);
    const startDate = new Date();
    startDate.setHours(hours, minutes, 0, 0);
    const endDate = new Date(startDate.getTime() + service.duration * 60000);
    const newEndTime = `${String(endDate.getHours()).padStart(2, '0')}:${String(endDate.getMinutes()).padStart(2, '0')}`;

    setForm(prev => ({
      ...prev,
      service_id: id,
      treatment_name: service.name,
      end_time: newEndTime
    }));
  };

  const handleStartTimeChange = (newStartTime: string) => {
    const service = serviceList.find(s => s.id === form.service_id);
    let newEndTime = form.end_time;
    if (service) {
      const [hours, minutes] = newStartTime.split(':').map(Number);
      const startDate = new Date();
      startDate.setHours(hours, minutes, 0, 0);
      const endDate = new Date(startDate.getTime() + service.duration * 60000);
      newEndTime = `${String(endDate.getHours()).padStart(2, '0')}:${String(endDate.getMinutes()).padStart(2, '0')}`;
    }
    setForm(prev => ({ ...prev, start_time: newStartTime, end_time: newEndTime }));
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-xl font-bold text-white mb-6">📅 새 예약 등록</h2>
        <div className="space-y-4">
          <div>
            <label htmlFor="res-customer-name" className="block text-sm text-dark-muted mb-1">고객명 *</label>
            <input id="res-customer-name" className="glass-input w-full" value={form.customer_name}
              onChange={(e) => setForm({ ...form, customer_name: e.target.value })} required />
          </div>
          <div>
            <label htmlFor="res-customer-phone" className="block text-sm text-dark-muted mb-1">연락처 *</label>
            <input id="res-customer-phone" className="glass-input w-full" value={form.customer_phone}
              onChange={(e) => setForm({ ...form, customer_phone: e.target.value })} required />
          </div>
          <div>
            <label htmlFor="res-treatment" className="block text-sm text-dark-muted mb-1">시술 메뉴</label>
            <select id="res-treatment" className="glass-input w-full" value={form.service_id}
              onChange={(e) => handleServiceChange(e.target.value)}>
              <option value="">메뉴 선택 (또는 직접 입력)</option>
              {serviceList.filter(s => s.is_active).map((s) => (
                <option key={s.id} value={s.id}>{s.category} - {s.name} ({s.duration}분)</option>
              ))}
            </select>
          </div>
          {!form.service_id && (
            <div>
              <label htmlFor="res-treatment-custom" className="block text-sm text-dark-muted mb-1">직접 입력 (시술명)</label>
              <input id="res-treatment-custom" className="glass-input w-full" value={form.treatment_name}
                onChange={(e) => setForm({ ...form, treatment_name: e.target.value })} placeholder="메뉴에 없는 시술인 경우 입력" />
            </div>
          )}
          <div>
            <label htmlFor="res-staff" className="block text-sm text-dark-muted mb-1">담당 직원</label>
            <select id="res-staff" className="glass-input w-full" value={form.staff_id}
              onChange={(e) => setForm({ ...form, staff_id: e.target.value })}>
              <option value="">선택 안함</option>
              {staffList.map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
            </select>
          </div>
          <div className="grid grid-cols-3 gap-3">
            <div>
              <label htmlFor="res-date" className="block text-sm text-dark-muted mb-1">날짜</label>
              <input id="res-date" type="date" className="glass-input w-full text-sm" value={form.date}
                onChange={(e) => setForm({ ...form, date: e.target.value })} />
            </div>
            <div>
              <label htmlFor="res-start" className="block text-sm text-dark-muted mb-1">시작</label>
              <input id="res-start" type="time" className="glass-input w-full text-sm" value={form.start_time}
                onChange={(e) => handleStartTimeChange(e.target.value)} />
            </div>
            <div>
              <label htmlFor="res-end" className="block text-sm text-dark-muted mb-1">종료</label>
              <input id="res-end" type="time" className="glass-input w-full text-sm" value={form.end_time}
                onChange={(e) => setForm({ ...form, end_time: e.target.value })} />
            </div>
          </div>
          <div>
            <label htmlFor="res-memo" className="block text-sm text-dark-muted mb-1">메모</label>
            <textarea id="res-memo" className="glass-input w-full" rows={2} value={form.memo}
              onChange={(e) => setForm({ ...form, memo: e.target.value })} />
          </div>
        </div>
        <div className="flex gap-3 mt-6">
          <button onClick={onClose} className="btn-secondary flex-1">취소</button>
          <button id="res-submit" onClick={() => onSubmit(form)} className="btn-primary flex-1">예약 등록</button>
        </div>
      </div>
    </div>
  );
}
