// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';
import api from '@/lib/api';
import { maskPhone, formatDate } from '@/lib/utils';
import type { Customer, CustomerWithHistory, Chart, ServiceMenu, Staff } from '@/types';

export default function CustomersPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [selectedCustomer, setSelectedCustomer] = useState<CustomerWithHistory | null>(null);
  const [serviceList, setServiceList] = useState<ServiceMenu[]>([]);
  const [staffList, setStaffList] = useState<Staff[]>([]);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchCustomers = useCallback(async () => {
    setIsLoading(true);
    try {
      if (searchQuery.length >= 2) {
        const data = await api.get<Customer[]>(`/api/customers/search?q=${encodeURIComponent(searchQuery)}&limit=30`);
        setCustomers(data || []);
      } else {
        const res = await api.get<{ data: Customer[]; total: number }>('/api/customers?limit=30');
        setCustomers(res?.data || []);
      }
    } catch { setCustomers([]); }
    setIsLoading(false);
  }, [searchQuery]);

  const fetchServices = useCallback(async () => {
    try {
      const data = await api.get<ServiceMenu[]>('/api/services');
      setServiceList(data || []);
    } catch { setServiceList([]); }
  }, []);

  const fetchStaff = useCallback(async () => {
    try {
      const data = await api.get<Staff[]>('/api/staffs');
      setStaffList(data || []);
    } catch { setStaffList([]); }
  }, []);

  useEffect(() => {
    const timer = setTimeout(fetchCustomers, 300);
    return () => clearTimeout(timer);
  }, [fetchCustomers]);

  useEffect(() => { fetchServices(); fetchStaff(); }, [fetchServices, fetchStaff]);

  const handleSelectCustomer = async (id: string) => {
    try {
      const data = await api.get<CustomerWithHistory>(`/api/customers/${id}`);
      setSelectedCustomer(data);
    } catch { /* error handled */ }
  };

  const handleCreateCustomer = async (form: Record<string, string>) => {
    // Errors propagate to the modal so the user can see the failure reason
    await api.post('/api/customers', form);
    setShowCreateModal(false);
    fetchCustomers();
  };

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Search Header */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <h1 className="text-2xl font-bold text-white">👤 고객 및 차트 관리</h1>
        <button id="create-customer-btn" onClick={() => setShowCreateModal(true)} className="btn-primary text-sm py-2">
          + 새 고객 등록
        </button>
      </div>

      {/* Sticky Search Bar */}
      <div className="sticky top-0 z-20 bg-dark-bg/80 backdrop-blur-xl py-3">
        <input
          id="customer-search"
          type="search"
          className="glass-input w-full text-lg py-4"
          placeholder="🔍 고객명 또는 전화번호 뒷자리로 검색..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          autoComplete="off"
        />
      </div>

      <div className="grid lg:grid-cols-5 gap-6">
        {/* Customer List */}
        <div className="lg:col-span-2 space-y-2 max-h-[70vh] overflow-y-auto">
          {isLoading ? (
            <div className="flex justify-center py-12"><div className="loading-spinner w-8 h-8" /></div>
          ) : customers.length === 0 ? (
            <div className="text-center py-12 text-dark-muted">
              <p className="text-4xl mb-3">🔍</p>
              <p>검색 결과가 없습니다</p>
            </div>
          ) : customers.map((c) => (
            <button key={c.id} id={`customer-${c.id}`}
              onClick={() => handleSelectCustomer(c.id)}
              className={`w-full glass-card p-4 text-left hover:border-salon-500/30 transition-all ${
                selectedCustomer?.id === c.id ? 'border-salon-500/50 bg-salon-500/5' : ''
              }`}>
              <div className="flex items-center justify-between">
                <div>
                  <span className="text-base font-semibold text-white">{c.name}</span>
                  <span className="ml-2 text-sm text-dark-muted">{maskPhone(c.phone)}</span>
                </div>
                <span className="text-xs text-dark-muted">방문 {c.visit_count}회</span>
              </div>
              {c.last_visited_at && (
                <p className="text-xs text-dark-muted mt-1">최근 방문: {formatDate(c.last_visited_at)}</p>
              )}
              {c.tags && c.tags.length > 0 && (
                <div className="flex gap-1 mt-2 flex-wrap">
                  {c.tags.map((tag, i) => (
                    <span key={i} className="badge bg-dark-surface text-dark-muted">{tag}</span>
                  ))}
                </div>
              )}
            </button>
          ))}
        </div>

        {/* Customer Detail / History Book */}
        <div className="lg:col-span-3">
          {selectedCustomer ? (
            <CustomerDetail customer={selectedCustomer} serviceList={serviceList} staffList={staffList} onRefresh={() => handleSelectCustomer(selectedCustomer.id)} />
          ) : (
            <div className="glass-card p-12 text-center">
              <p className="text-4xl mb-3">📋</p>
              <p className="text-dark-muted">고객을 선택하면 상세 정보와 시술 히스토리를 확인할 수 있습니다</p>
            </div>
          )}
        </div>
      </div>

      {/* Create Customer Modal */}
      {showCreateModal && (
        <CreateCustomerModal onClose={() => setShowCreateModal(false)} onSubmit={handleCreateCustomer} />
      )}
    </div>
  );
}

function CustomerDetail({ customer, serviceList, staffList, onRefresh }: { customer: CustomerWithHistory; serviceList: ServiceMenu[]; staffList: Staff[]; onRefresh: () => void }) {
  const [showChartForm, setShowChartForm] = useState(false);
  const [showMembershipForm, setShowMembershipForm] = useState(false);

  const handleCreateChart = async (form: Record<string, string>) => {
    // Errors propagate to the form so the user can see the failure reason
    await api.post('/api/charts', {
      ...form,
      customer_id: customer.id,
    });
    setShowChartForm(false);
    onRefresh();
  };

  const handleCreateMembership = async (form: Record<string, unknown>) => {
    await api.post('/api/memberships', {
      ...form,
      customer_id: customer.id,
    });
    setShowMembershipForm(false);
    onRefresh();
  };

  return (
    <div className="space-y-4">
      {/* Customer Info Card */}
      <div className="glass-card p-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-4">
            <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-salon-400 to-purple-600 flex items-center justify-center text-white text-xl font-bold">
              {customer.name.charAt(0)}
            </div>
            <div>
              <h2 className="text-xl font-bold text-white">{customer.name}</h2>
              <p className="text-sm text-dark-muted">{maskPhone(customer.phone)}</p>
            </div>
          </div>
          <div className="text-right">
            <p className="text-sm text-dark-muted">총 방문</p>
            <p className="text-2xl font-bold text-salon-400">{customer.visit_count}회</p>
          </div>
        </div>
        {customer.memo && (
          <div className="bg-dark-surface/50 rounded-xl p-3 mt-3">
            <p className="text-sm text-dark-muted">{customer.memo}</p>
          </div>
        )}
      </div>

      {/* Memberships */}
      <div className="glass-card p-4">
        <div className="flex items-center justify-between mb-3">
          <h3 className="font-semibold text-white">💳 멤버십</h3>
          <button id="add-membership-btn" onClick={() => setShowMembershipForm(true)} className="btn-primary text-xs py-1.5 px-3">
            + 회원권 등록
          </button>
        </div>

        {showMembershipForm && (
          <MembershipForm onSubmit={handleCreateMembership} onCancel={() => setShowMembershipForm(false)} />
        )}

        {(!customer.memberships || customer.memberships.length === 0) ? (
          <p className="text-dark-muted text-sm text-center py-4">등록된 회원권이 없습니다</p>
        ) : (
          <div className="space-y-2">
            {customer.memberships.map((m) => (
              <div key={m.id} className="flex items-center justify-between bg-dark-surface/50 rounded-xl p-3">
                <div>
                  <span className="text-sm font-medium text-white">{m.name}</span>
                  <span className="ml-2 badge bg-salon-500/20 text-salon-400">
                    {m.type === 'money' ? '금액형' : '횟수형'}
                  </span>
                  {!m.is_active && <span className="ml-2 badge bg-red-500/20 text-red-400">만료</span>}
                </div>
                <div className="text-right">
                  {m.type === 'money' ? (
                    <p className="text-sm font-bold text-white">{m.balance.toLocaleString()}원</p>
                  ) : (
                    <p className="text-sm font-bold text-white">{m.remaining_count}회 남음</p>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Chart History */}
      <div className="glass-card p-4">
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-semibold text-white">📖 시술 히스토리북</h3>
          <button id="add-chart-btn" onClick={() => setShowChartForm(true)} className="btn-primary text-xs py-1.5 px-3">
            + 시술 기록 추가
          </button>
        </div>

        {showChartForm && (
          <ChartForm serviceList={serviceList} staffList={staffList} onSubmit={handleCreateChart} onCancel={() => setShowChartForm(false)} />
        )}

        <div className="space-y-3">
          {(!customer.charts || customer.charts.length === 0) ? (
            <p className="text-dark-muted text-sm text-center py-6">시술 기록이 없습니다</p>
          ) : customer.charts.map((chart) => (
            <ChartEntry key={chart.id} chart={chart} />
          ))}
        </div>
      </div>
    </div>
  );
}

function ChartEntry({ chart }: { chart: Chart }) {
  return (
    <div className="bg-dark-surface/50 rounded-xl p-4 border-l-4 border-salon-500/30">
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          {chart.treatment_name && <span className="text-sm font-semibold text-salon-400">{chart.treatment_name}</span>}
          {chart.staff_name && <span className="text-xs text-dark-muted">by {chart.staff_name}</span>}
        </div>
        <span className="text-xs text-dark-muted">{formatDate(chart.created_at)}</span>
      </div>
      {chart.recipe && <p className="text-sm text-dark-text">{chart.recipe}</p>}
      {chart.notes && <p className="text-xs text-dark-muted mt-1">{chart.notes}</p>}
      {(chart.before_img_url || chart.after_img_url) && (
        <div className="flex gap-3 mt-3">
          {chart.before_img_url && (
            <div className="flex-1">
              <p className="text-xs text-dark-muted mb-1">Before</p>
              <div className="aspect-video bg-dark-surface rounded-lg overflow-hidden">
                <img src={chart.before_img_url} alt="Before" className="w-full h-full object-cover" />
              </div>
            </div>
          )}
          {chart.after_img_url && (
            <div className="flex-1">
              <p className="text-xs text-dark-muted mb-1">After</p>
              <div className="aspect-video bg-dark-surface rounded-lg overflow-hidden">
                <img src={chart.after_img_url} alt="After" className="w-full h-full object-cover" />
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

function ChartForm({ serviceList, staffList, onSubmit, onCancel }: { serviceList: ServiceMenu[]; staffList: Staff[]; onSubmit: (form: Record<string, string>) => void; onCancel: () => void }) {
  const [form, setForm] = useState({ staff_id: '', service_id: '', recipe: '', treatment_name: '', notes: '' });
  const [error, setError] = useState('');

  const handleServiceChange = (id: string) => {
    const service = serviceList.find(s => s.id === id);
    if (service) {
      setForm({ ...form, service_id: id, treatment_name: service.name });
    } else {
      setForm({ ...form, service_id: '', treatment_name: '' });
    }
  };

  const handleSubmit = async () => {
    setError('');
    try {
      if (!form.staff_id) throw new Error('담당 직원을 선택해주세요.');
      if (!form.service_id && !form.treatment_name) throw new Error('시술 메뉴를 선택하거나 시술명을 입력해주세요.');
      await onSubmit(form);
    } catch (err: any) {
      setError(err.message || '시술 기록 저장에 실패했습니다.');
    }
  };

  return (
    <div className="glass-card p-4 mb-4 border-salon-500/30">
      {error && (
        <div className="mb-3 p-2 bg-red-500/20 border border-red-500/50 rounded-lg">
          <p className="text-xs text-red-400">{error}</p>
        </div>
      )}
      <div className="space-y-3">
        <select id="chart-staff" className="glass-input w-full text-sm" value={form.staff_id}
          onChange={(e) => setForm({ ...form, staff_id: e.target.value })}>
          <option value="">담당 직원 선택 *</option>
          {staffList.map(s => (
            <option key={s.id} value={s.id}>{s.name}</option>
          ))}
        </select>
        <select className="glass-input w-full text-sm" value={form.service_id}
          onChange={(e) => handleServiceChange(e.target.value)}>
          <option value="">시술 메뉴 선택 (또는 직접 입력)</option>
          {serviceList.filter(s => s.is_active).map(s => (
            <option key={s.id} value={s.id}>{s.category} - {s.name}</option>
          ))}
        </select>
        {!form.service_id && (
          <input id="chart-treatment" className="glass-input w-full text-sm" placeholder="시술명 직접 입력" value={form.treatment_name}
            onChange={(e) => setForm({ ...form, treatment_name: e.target.value })} />
        )}
        <textarea id="chart-recipe" className="glass-input w-full text-sm" rows={3} placeholder="레시피 (약제, 시간 등)" value={form.recipe}
          onChange={(e) => setForm({ ...form, recipe: e.target.value })} />
        <textarea id="chart-notes" className="glass-input w-full text-sm" rows={2} placeholder="메모" value={form.notes}
          onChange={(e) => setForm({ ...form, notes: e.target.value })} />
        <div className="flex gap-2">
          <button onClick={onCancel} className="btn-secondary text-xs flex-1">취소</button>
          <button id="chart-submit" onClick={handleSubmit} className="btn-primary text-xs flex-1">저장</button>
        </div>
      </div>
    </div>
  );
}

function MembershipForm({ onSubmit, onCancel }: { onSubmit: (form: Record<string, unknown>) => Promise<void>; onCancel: () => void }) {
  const [form, setForm] = useState({
    type: 'money' as 'money' | 'count',
    name: '',
    initial_balance: '',
    initial_count: '',
    target_treatment: '',
    expired_at: '',
  });
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    setError('');
    try {
      if (!form.name) throw new Error('회원권 이름을 입력해주세요.');
      if (form.type === 'money' && !(parseFloat(form.initial_balance) > 0)) throw new Error('충전 금액을 입력해주세요.');
      if (form.type === 'count' && !(parseInt(form.initial_count) > 0)) throw new Error('사용 횟수를 입력해주세요.');
      await onSubmit({
        type: form.type,
        name: form.name,
        initial_balance: form.type === 'money' ? parseFloat(form.initial_balance) || 0 : 0,
        initial_count: form.type === 'count' ? parseInt(form.initial_count) || 0 : 0,
        target_treatment: form.target_treatment,
        expired_at: form.expired_at ? new Date(`${form.expired_at}T23:59:59`).toISOString() : '',
      });
    } catch (err: any) {
      setError(err.message || '회원권 등록에 실패했습니다.');
    }
  };

  return (
    <div className="glass-card p-4 mb-4 border-salon-500/30">
      {error && (
        <div className="mb-3 p-2 bg-red-500/20 border border-red-500/50 rounded-lg">
          <p className="text-xs text-red-400">{error}</p>
        </div>
      )}
      <div className="space-y-3">
        <div className="flex gap-2">
          <button onClick={() => setForm({ ...form, type: 'money' })}
            className={`flex-1 py-2 rounded-xl text-sm font-medium transition-all ${
              form.type === 'money' ? 'bg-salon-500/20 text-salon-400 border border-salon-500/30' : 'bg-dark-surface text-dark-muted'}`}>
            금액형 (정액권)
          </button>
          <button onClick={() => setForm({ ...form, type: 'count' })}
            className={`flex-1 py-2 rounded-xl text-sm font-medium transition-all ${
              form.type === 'count' ? 'bg-salon-500/20 text-salon-400 border border-salon-500/30' : 'bg-dark-surface text-dark-muted'}`}>
            횟수형 (회원권)
          </button>
        </div>
        <input id="membership-name" className="glass-input w-full text-sm" placeholder="회원권 이름 (예: 50만원 정액권)" value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })} />
        {form.type === 'money' ? (
          <input id="membership-balance" type="number" className="glass-input w-full text-sm" placeholder="충전 금액 (원)" value={form.initial_balance}
            onChange={(e) => setForm({ ...form, initial_balance: e.target.value })} />
        ) : (
          <>
            <input id="membership-count" type="number" className="glass-input w-full text-sm" placeholder="사용 횟수 (회)" value={form.initial_count}
              onChange={(e) => setForm({ ...form, initial_count: e.target.value })} />
            <input id="membership-target" className="glass-input w-full text-sm" placeholder="대상 시술 (예: 두피 클리닉)" value={form.target_treatment}
              onChange={(e) => setForm({ ...form, target_treatment: e.target.value })} />
          </>
        )}
        <div>
          <label htmlFor="membership-expiry" className="block text-xs text-dark-muted mb-1">유효기간 (선택)</label>
          <input id="membership-expiry" type="date" className="glass-input w-full text-sm" value={form.expired_at}
            onChange={(e) => setForm({ ...form, expired_at: e.target.value })} />
        </div>
        <div className="flex gap-2">
          <button onClick={onCancel} className="btn-secondary text-xs flex-1">취소</button>
          <button id="membership-submit" onClick={handleSubmit} className="btn-primary text-xs flex-1">등록</button>
        </div>
      </div>
    </div>
  );
}

function CreateCustomerModal({ onClose, onSubmit }: { onClose: () => void; onSubmit: (data: Record<string, string>) => void }) {
  const [form, setForm] = useState({ name: '', phone: '', email: '', birth_date: '', memo: '' });
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    setError('');
    try {
      if (!form.name) throw new Error('이름을 입력해주세요.');
      if (!form.phone) throw new Error('전화번호를 입력해주세요.');
      await onSubmit(form);
    } catch (err: any) {
      setError(err.message || '고객 등록에 실패했습니다.');
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-xl font-bold text-white mb-6">👤 새 고객 등록</h2>
        {error && (
          <div className="mb-4 p-3 bg-red-500/20 border border-red-500/50 rounded-lg">
            <p className="text-sm text-red-400">{error}</p>
          </div>
        )}
        <div className="space-y-4">
          <div>
            <label htmlFor="cust-name" className="block text-sm text-dark-muted mb-1">이름 *</label>
            <input id="cust-name" className="glass-input w-full" value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          </div>
          <div>
            <label htmlFor="cust-phone" className="block text-sm text-dark-muted mb-1">전화번호 *</label>
            <input id="cust-phone" className="glass-input w-full" value={form.phone}
              onChange={(e) => setForm({ ...form, phone: e.target.value })} required />
          </div>
          <div>
            <label htmlFor="cust-email" className="block text-sm text-dark-muted mb-1">이메일</label>
            <input id="cust-email" type="email" className="glass-input w-full" value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })} />
          </div>
          <div>
            <label htmlFor="cust-birth" className="block text-sm text-dark-muted mb-1">생년월일</label>
            <input id="cust-birth" type="date" className="glass-input w-full" value={form.birth_date}
              onChange={(e) => setForm({ ...form, birth_date: e.target.value })} />
          </div>
          <div>
            <label htmlFor="cust-memo" className="block text-sm text-dark-muted mb-1">메모</label>
            <textarea id="cust-memo" className="glass-input w-full" rows={2} value={form.memo}
              onChange={(e) => setForm({ ...form, memo: e.target.value })} />
          </div>
        </div>
        <div className="flex gap-3 mt-6">
          <button onClick={onClose} className="btn-secondary flex-1">취소</button>
          <button id="cust-submit" onClick={handleSubmit} className="btn-primary flex-1">등록</button>
        </div>
      </div>
    </div>
  );
}
