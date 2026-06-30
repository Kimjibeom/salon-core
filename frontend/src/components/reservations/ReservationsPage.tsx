// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';
import api from '@/lib/api';
import { formatTime, getStatusColor, getStatusLabel, getSourceLabel } from '@/lib/utils';
import type { Reservation, WaitingQueueEntry, Staff, ServiceMenu, Customer } from '@/types';
import SaleEntryModal from '../pos/SaleEntryModal';

type ViewMode = 'day' | 'week' | 'month';

export default function ReservationsPage() {
  const [viewMode, setViewMode] = useState<ViewMode>('day');
  const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0]);
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [waitingQueue, setWaitingQueue] = useState<WaitingQueueEntry[]>([]);
  const [staffList, setStaffList] = useState<Staff[]>([]);
  const [serviceList, setServiceList] = useState<ServiceMenu[]>([]);
  const [customerList, setCustomerList] = useState<Customer[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [showWalkinModal, setShowWalkinModal] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedReservationForSale, setSelectedReservationForSale] = useState<Reservation | null>(null);

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

  const fetchCustomers = useCallback(async () => {
    try {
      const data = await api.get<Customer[]>('/api/customers');
      setCustomerList(data || []);
    } catch { setCustomerList([]); }
  }, []);

  useEffect(() => { fetchReservations(); fetchWaitingQueue(); fetchStaff(); fetchServices(); fetchCustomers(); }, [fetchReservations, fetchWaitingQueue, fetchStaff, fetchServices, fetchCustomers]);

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

  const handleCompleteClick = (r: Reservation) => {
    setSelectedReservationForSale(r);
  };

  const handleSubmitSale = async (saleData: Record<string, unknown>) => {
    try {
      await api.post('/api/sales', saleData);
      if (selectedReservationForSale) {
        await handleStatusChange(selectedReservationForSale.id, 'completed');
      }
      setSelectedReservationForSale(null);
    } catch (err: any) {
      alert(err.message || '매출 등록에 실패했습니다.');
    }
  };

  const handleAddWalkin = async (formData: Record<string, string>) => {
    try {
      await api.post('/api/reservations/waiting', {
        ...formData,
        source: 'offline',
      });
      setShowWalkinModal(false);
      fetchWaitingQueue();
    } catch (err: any) {
      alert(err.message || '워크인 등록에 실패했습니다.');
    }
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
                                <button onClick={() => handleCompleteClick(r)}
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
            onClick={() => setShowWalkinModal(true)}
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
          customerList={customerList}
          onClose={() => setShowModal(false)}
          onSubmit={handleCreateReservation}
          fetchCustomers={fetchCustomers}
        />
      )}

      {/* Walkin Modal */}
      {showWalkinModal && (
        <WalkinModal
          serviceList={serviceList}
          customerList={customerList}
          onClose={() => setShowWalkinModal(false)}
          onSubmit={handleAddWalkin}
          fetchCustomers={fetchCustomers}
        />
      )}

      {/* Sale Entry Modal for Completion */}
      {selectedReservationForSale && (
        <SaleEntryModal
          staffList={staffList}
          serviceList={serviceList}
          onClose={() => setSelectedReservationForSale(null)}
          onSubmit={handleSubmitSale}
          initialData={{
            staff_id: selectedReservationForSale.staff_id || '',
            service_id: selectedReservationForSale.service_id || '',
            item_name: selectedReservationForSale.treatment_name || '',
            memo: selectedReservationForSale.memo || '',
          }}
        />
      )}
    </div>
  );
}

function ReservationModal({ staffList, serviceList, customerList, onClose, onSubmit, fetchCustomers }: {
  staffList: Staff[];
  serviceList: ServiceMenu[];
  customerList: Customer[];
  onClose: () => void;
  onSubmit: (data: Record<string, string>) => void;
  fetchCustomers: () => Promise<void>;
}) {
  const [isNewCustomer, setIsNewCustomer] = useState(false);
  const [form, setForm] = useState({
    customer_id: '',
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

  const [error, setError] = useState('');

  const handleCustomerSelect = (id: string) => {
    const cust = customerList.find(c => c.id === id);
    if (cust) {
      setForm(prev => ({ ...prev, customer_id: id, customer_name: cust.name, customer_phone: cust.phone || '' }));
    } else {
      setForm(prev => ({ ...prev, customer_id: '', customer_name: '', customer_phone: '' }));
    }
  };

  const handleSubmit = async () => {
    setError('');
    try {
      if (isNewCustomer) {
        if (!form.customer_name) throw new Error('고객명을 입력해주세요.');
        if (!form.customer_phone) throw new Error('연락처를 입력해주세요.');
        // Create customer first
        const custRes = await api.post<{id: string}>('/api/customers', {
          name: form.customer_name,
          phone: form.customer_phone,
        });
        if (custRes && custRes.id) {
          form.customer_id = custRes.id;
          await fetchCustomers();
        }
      } else {
        if (!form.customer_id) throw new Error('고객을 선택해주세요.');
      }
      
      if (!form.date || !form.start_time || !form.end_time) throw new Error('날짜와 시간을 확인해주세요.');
      
      await onSubmit(form);
    } catch (err: any) {
      setError(err.message || '예약 등록에 실패했습니다.');
    }
  };

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
        
        {error && (
          <div className="mb-4 p-3 bg-red-500/20 border border-red-500/50 rounded-lg">
            <p className="text-sm text-red-400">{error}</p>
          </div>
        )}

        <div className="flex gap-4 mb-4">
          <label className="flex items-center gap-2 cursor-pointer text-white">
            <input type="radio" name="customer_type" checked={!isNewCustomer} onChange={() => setIsNewCustomer(false)} />
            기존 고객
          </label>
          <label className="flex items-center gap-2 cursor-pointer text-white">
            <input type="radio" name="customer_type" checked={isNewCustomer} onChange={() => setIsNewCustomer(true)} />
            신규 고객
          </label>
        </div>

        <div className="space-y-4">
          {!isNewCustomer ? (
            <div>
              <label htmlFor="res-customer-select" className="block text-sm text-dark-muted mb-1">고객 선택 *</label>
              <select id="res-customer-select" className="glass-input w-full" value={form.customer_id}
                onChange={(e) => handleCustomerSelect(e.target.value)}>
                <option value="">고객을 선택하세요</option>
                {customerList.map(c => (
                  <option key={c.id} value={c.id}>{c.name} ({c.phone})</option>
                ))}
              </select>
            </div>
          ) : (
            <>
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
            </>
          )}

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
          <button id="res-submit" onClick={handleSubmit} className="btn-primary flex-1">예약 등록</button>
        </div>
      </div>
    </div>
  );
}

function WalkinModal({ serviceList, customerList, onClose, onSubmit, fetchCustomers }: {
  serviceList: ServiceMenu[];
  customerList: Customer[];
  onClose: () => void;
  onSubmit: (data: Record<string, string>) => void;
  fetchCustomers: () => Promise<void>;
}) {
  const [isNewCustomer, setIsNewCustomer] = useState(false);
  const [form, setForm] = useState({
    customer_id: '',
    customer_name: '',
    customer_phone: '',
    service_id: '',
    treatment_name: '',
    start_time: new Date().toLocaleTimeString('en-US', { hour12: false, hour: '2-digit', minute: '2-digit' }),
    end_time: new Date().toLocaleTimeString('en-US', { hour12: false, hour: '2-digit', minute: '2-digit' }),
  });

  const [error, setError] = useState('');

  const handleCustomerSelect = (id: string) => {
    const cust = customerList.find(c => c.id === id);
    if (cust) {
      setForm(prev => ({ ...prev, customer_id: id, customer_name: cust.name, customer_phone: cust.phone || '' }));
    } else {
      setForm(prev => ({ ...prev, customer_id: '', customer_name: '', customer_phone: '' }));
    }
  };

  const handleServiceChange = (id: string) => {
    const service = serviceList.find(s => s.id === id);
    setForm(prev => ({
      ...prev,
      service_id: id,
      treatment_name: service ? service.name : ''
    }));
  };

  const handleSubmit = async () => {
    setError('');
    try {
      if (isNewCustomer) {
        if (!form.customer_name) throw new Error('고객명을 입력해주세요.');
        const custRes = await api.post<{id: string}>('/api/customers', {
          name: form.customer_name,
          phone: form.customer_phone,
        });
        if (custRes && custRes.id) {
          form.customer_id = custRes.id;
          await fetchCustomers();
        }
      } else {
        if (!form.customer_id && !form.customer_name) {
          // Allow anonymous walkin if name is given
          if (!form.customer_name) throw new Error('고객을 선택하거나 이름을 입력해주세요.');
        }
      }
      
      await onSubmit(form);
    } catch (err: any) {
      setError(err.message || '워크인 등록에 실패했습니다.');
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-xl font-bold text-white mb-6">🚶 워크인 대기열 추가</h2>
        
        {error && (
          <div className="mb-4 p-3 bg-red-500/20 border border-red-500/50 rounded-lg">
            <p className="text-sm text-red-400">{error}</p>
          </div>
        )}

        <div className="flex gap-4 mb-4">
          <label className="flex items-center gap-2 cursor-pointer text-white">
            <input type="radio" name="walkin_customer_type" checked={!isNewCustomer} onChange={() => setIsNewCustomer(false)} />
            기존 고객
          </label>
          <label className="flex items-center gap-2 cursor-pointer text-white">
            <input type="radio" name="walkin_customer_type" checked={isNewCustomer} onChange={() => setIsNewCustomer(true)} />
            신규 고객
          </label>
        </div>

        <div className="space-y-4">
          {!isNewCustomer ? (
            <div>
              <label htmlFor="walkin-customer-select" className="block text-sm text-dark-muted mb-1">고객 선택</label>
              <select id="walkin-customer-select" className="glass-input w-full" value={form.customer_id}
                onChange={(e) => handleCustomerSelect(e.target.value)}>
                <option value="">익명 (비회원)</option>
                {customerList.map(c => (
                  <option key={c.id} value={c.id}>{c.name} ({c.phone})</option>
                ))}
              </select>
              {!form.customer_id && (
                <div className="mt-2">
                  <label htmlFor="walkin-anon-name" className="block text-sm text-dark-muted mb-1">비회원 이름</label>
                  <input id="walkin-anon-name" className="glass-input w-full" value={form.customer_name}
                    onChange={(e) => setForm({ ...form, customer_name: e.target.value })} placeholder="Walk-in" />
                </div>
              )}
            </div>
          ) : (
            <>
              <div>
                <label htmlFor="walkin-customer-name" className="block text-sm text-dark-muted mb-1">고객명 *</label>
                <input id="walkin-customer-name" className="glass-input w-full" value={form.customer_name}
                  onChange={(e) => setForm({ ...form, customer_name: e.target.value })} required />
              </div>
              <div>
                <label htmlFor="walkin-customer-phone" className="block text-sm text-dark-muted mb-1">연락처</label>
                <input id="walkin-customer-phone" className="glass-input w-full" value={form.customer_phone}
                  onChange={(e) => setForm({ ...form, customer_phone: e.target.value })} />
              </div>
            </>
          )}

          <div>
            <label htmlFor="walkin-treatment" className="block text-sm text-dark-muted mb-1">시술 메뉴</label>
            <select id="walkin-treatment" className="glass-input w-full" value={form.service_id}
              onChange={(e) => handleServiceChange(e.target.value)}>
              <option value="">선택 안함</option>
              {serviceList.filter(s => s.is_active).map((s) => (
                <option key={s.id} value={s.id}>{s.category} - {s.name}</option>
              ))}
            </select>
          </div>
          {!form.service_id && (
            <div>
              <label htmlFor="walkin-treatment-custom" className="block text-sm text-dark-muted mb-1">직접 입력 (시술명)</label>
              <input id="walkin-treatment-custom" className="glass-input w-full" value={form.treatment_name}
                onChange={(e) => setForm({ ...form, treatment_name: e.target.value })} placeholder="예: 커트, 펌" />
            </div>
          )}
        </div>
        <div className="flex gap-3 mt-6">
          <button onClick={onClose} className="btn-secondary flex-1">취소</button>
          <button onClick={handleSubmit} className="btn-primary flex-1">대기열 추가</button>
        </div>
      </div>
    </div>
  );
}
