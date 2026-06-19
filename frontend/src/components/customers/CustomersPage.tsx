// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';
import api from '@/lib/api';
import { maskPhone, formatDate } from '@/lib/utils';
import type { Customer, CustomerWithHistory, Chart, ServiceMenu } from '@/types';

export default function CustomersPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [selectedCustomer, setSelectedCustomer] = useState<CustomerWithHistory | null>(null);
  const [serviceList, setServiceList] = useState<ServiceMenu[]>([]);
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

  useEffect(() => {
    const timer = setTimeout(fetchCustomers, 300);
    return () => clearTimeout(timer);
  }, [fetchCustomers]);

  useEffect(() => { fetchServices(); }, [fetchServices]);

  const handleSelectCustomer = async (id: string) => {
    try {
      const data = await api.get<CustomerWithHistory>(`/api/customers/${id}`);
      setSelectedCustomer(data);
    } catch { /* error handled */ }
  };

  const handleCreateCustomer = async (form: Record<string, string>) => {
    try {
      await api.post('/api/customers', form);
      setShowCreateModal(false);
      fetchCustomers();
    } catch { /* error handled */ }
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
            <CustomerDetail customer={selectedCustomer} serviceList={serviceList} onRefresh={() => handleSelectCustomer(selectedCustomer.id)} />
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

function CustomerDetail({ customer, serviceList, onRefresh }: { customer: CustomerWithHistory; serviceList: ServiceMenu[]; onRefresh: () => void }) {
  const [showChartForm, setShowChartForm] = useState(false);

  const handleCreateChart = async (form: Record<string, string>) => {
    try {
      await api.post('/api/charts', {
        ...form,
        customer_id: customer.id,
      });
      setShowChartForm(false);
      onRefresh();
    } catch { /* error handled */ }
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
      {customer.memberships && customer.memberships.length > 0 && (
        <div className="glass-card p-4">
          <h3 className="font-semibold text-white mb-3">💳 멤버십</h3>
          <div className="space-y-2">
            {customer.memberships.map((m) => (
              <div key={m.id} className="flex items-center justify-between bg-dark-surface/50 rounded-xl p-3">
                <div>
                  <span className="text-sm font-medium text-white">{m.name}</span>
                  <span className="ml-2 badge bg-salon-500/20 text-salon-400">
                    {m.type === 'money' ? '금액형' : '횟수형'}
                  </span>
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
        </div>
      )}

      {/* Chart History */}
      <div className="glass-card p-4">
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-semibold text-white">📖 시술 히스토리북</h3>
          <button id="add-chart-btn" onClick={() => setShowChartForm(true)} className="btn-primary text-xs py-1.5 px-3">
            + 시술 기록 추가
          </button>
        </div>

        {showChartForm && (
          <ChartForm serviceList={serviceList} onSubmit={handleCreateChart} onCancel={() => setShowChartForm(false)} />
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

function ChartForm({ serviceList, onSubmit, onCancel }: { serviceList: ServiceMenu[]; onSubmit: (form: Record<string, string>) => void; onCancel: () => void }) {
  const [form, setForm] = useState({ staff_id: '', service_id: '', recipe: '', treatment_name: '', notes: '' });

  const handleServiceChange = (id: string) => {
    const service = serviceList.find(s => s.id === id);
    if (service) {
      setForm({ ...form, service_id: id, treatment_name: service.name });
    } else {
      setForm({ ...form, service_id: '', treatment_name: '' });
    }
  };

  return (
    <div className="glass-card p-4 mb-4 border-salon-500/30">
      <div className="space-y-3">
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
          <button id="chart-submit" onClick={() => onSubmit(form)} className="btn-primary text-xs flex-1">저장</button>
        </div>
      </div>
    </div>
  );
}

function CreateCustomerModal({ onClose, onSubmit }: { onClose: () => void; onSubmit: (data: Record<string, string>) => void }) {
  const [form, setForm] = useState({ name: '', phone: '', email: '', birth_date: '', memo: '' });
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-xl font-bold text-white mb-6">👤 새 고객 등록</h2>
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
          <button id="cust-submit" onClick={() => onSubmit(form)} className="btn-primary flex-1">등록</button>
        </div>
      </div>
    </div>
  );
}
