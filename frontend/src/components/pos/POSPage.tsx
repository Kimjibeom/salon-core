// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';
import api from '@/lib/api';
import { formatCurrency, getCategoryLabel, getPaymentLabel } from '@/lib/utils';
import type { Staff, Membership, Sale, ServiceMenu } from '@/types';

export default function POSPage() {
  const [staffList, setStaffList] = useState<Staff[]>([]);
  const [recentSales, setRecentSales] = useState<Sale[]>([]);
  const [serviceList, setServiceList] = useState<ServiceMenu[]>([]);
  const [showSaleModal, setShowSaleModal] = useState(false);

  const fetchStaff = useCallback(async () => {
    try {
      const data = await api.get<Staff[]>('/api/staffs');
      setStaffList(data || []);
    } catch { setStaffList([]); }
  }, []);

  const fetchRecentSales = useCallback(async () => {
    const today = new Date().toISOString().split('T')[0];
    try {
      const data = await api.get<Sale[]>(`/api/sales?start_date=${today}&end_date=${today}`);
      setRecentSales(data || []);
    } catch { setRecentSales([]); }
  }, []);

  const fetchServices = useCallback(async () => {
    try {
      const data = await api.get<ServiceMenu[]>('/api/services');
      setServiceList(data || []);
    } catch { setServiceList([]); }
  }, []);

  useEffect(() => { fetchStaff(); fetchRecentSales(); fetchServices(); }, [fetchStaff, fetchRecentSales, fetchServices]);

  const todayTotal = recentSales.reduce((sum, s) => sum + s.total_amount, 0);
  const serviceTotal = recentSales.filter((s) => s.category === 'service').reduce((sum, s) => sum + s.total_amount, 0);
  const productTotal = recentSales.filter((s) => s.category === 'product').reduce((sum, s) => sum + s.total_amount, 0);

  const handleSubmitSale = async (saleData: Record<string, unknown>) => {
    try {
      await api.post('/api/sales', saleData);
      setShowSaleModal(false);
      fetchRecentSales();
    } catch { /* error handled */ }
  };

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">💳 간편 POS</h1>
        <button id="new-sale-btn" onClick={() => setShowSaleModal(true)} className="btn-primary text-sm py-2">
          + 영업 입력
        </button>
      </div>

      {/* Today Summary */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="stat-card">
          <span className="stat-label">오늘 총 매출</span>
          <span className="stat-value">{formatCurrency(todayTotal)}</span>
          <span className="text-xs text-dark-muted">{recentSales.length}건</span>
        </div>
        <div className="stat-card">
          <span className="stat-label">시술 매출</span>
          <span className="text-xl font-bold text-green-400">{formatCurrency(serviceTotal)}</span>
        </div>
        <div className="stat-card">
          <span className="stat-label">점판 매출</span>
          <span className="text-xl font-bold text-blue-400">{formatCurrency(productTotal)}</span>
        </div>
      </div>

      {/* Recent Sales Table */}
      <div className="glass-card overflow-hidden">
        <div className="p-4 border-b border-dark-border">
          <h3 className="font-semibold text-white">📝 오늘 매출 내역</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-dark-surface/50">
              <tr>
                <th className="table-header text-left">시간</th>
                <th className="table-header text-left">항목</th>
                <th className="table-header text-left">분류</th>
                <th className="table-header text-left">담당</th>
                <th className="table-header text-left">결제</th>
                <th className="table-header text-right">금액</th>
              </tr>
            </thead>
            <tbody>
              {recentSales.length === 0 ? (
                <tr>
                  <td colSpan={6} className="table-cell text-center text-dark-muted py-8">오늘 매출 내역이 없습니다</td>
                </tr>
              ) : recentSales.map((sale) => (
                <tr key={sale.id} className="table-row">
                  <td className="table-cell text-dark-muted">{new Date(sale.created_at).toLocaleTimeString('ko-KR', { hour: '2-digit', minute: '2-digit' })}</td>
                  <td className="table-cell font-medium text-white">{sale.item_name}</td>
                  <td className="table-cell">
                    <span className={`badge ${sale.category === 'service' ? 'bg-green-500/20 text-green-400' : 'bg-blue-500/20 text-blue-400'}`}>
                      {getCategoryLabel(sale.category)}
                    </span>
                  </td>
                  <td className="table-cell text-dark-muted">{sale.staff_name}</td>
                  <td className="table-cell">
                    <span className="badge bg-dark-surface text-dark-muted">{getPaymentLabel(sale.payment_method)}</span>
                  </td>
                  <td className="table-cell text-right font-semibold text-white">{formatCurrency(sale.total_amount)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Sale Entry Modal */}
      {showSaleModal && (
        <SaleEntryModal staffList={staffList} serviceList={serviceList} onClose={() => setShowSaleModal(false)} onSubmit={handleSubmitSale} />
      )}
    </div>
  );
}

function SaleEntryModal({ staffList, serviceList, onClose, onSubmit }: {
  staffList: Staff[];
  serviceList: ServiceMenu[];
  onClose: () => void;
  onSubmit: (data: Record<string, unknown>) => void;
}) {
  const [form, setForm] = useState({
    staff_id: '',
    service_id: '',
    item_name: '',
    total_amount: '',
    category: 'service' as 'service' | 'product',
    payment_method: 'card' as 'card' | 'cash' | 'membership' | 'mixed',
    card_amount: '',
    cash_amount: '',
    membership_amount: '',
    card_commission_rate: '2.5',
    membership_id: '',
    memo: '',
  });

  const totalAmount = parseFloat(form.total_amount) || 0;

  const handlePaymentMethodChange = (method: string) => {
    setForm((prev) => ({
      ...prev,
      payment_method: method as typeof prev.payment_method,
      card_amount: method === 'card' ? String(totalAmount) : method === 'mixed' ? prev.card_amount : '0',
      cash_amount: method === 'cash' ? String(totalAmount) : method === 'mixed' ? prev.cash_amount : '0',
      membership_amount: method === 'membership' ? String(totalAmount) : method === 'mixed' ? prev.membership_amount : '0',
    }));
  };

  const handleServiceChange = (id: string) => {
    const service = serviceList.find(s => s.id === id);
    if (!service) {
      setForm(prev => ({ ...prev, service_id: '', item_name: '', total_amount: '0' }));
      return;
    }
    
    setForm(prev => ({
      ...prev,
      service_id: id,
      item_name: service.name,
      category: service.category === '점판' ? 'product' : 'service',
      total_amount: service.price.toString()
    }));
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content max-w-xl" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-xl font-bold text-white mb-6">💳 영업 입력</h2>
        <div className="space-y-4">
          {/* Staff Selection */}
          <div>
            <label htmlFor="sale-staff" className="block text-sm text-dark-muted mb-1">담당 직원 *</label>
            <select id="sale-staff" className="glass-input w-full" value={form.staff_id}
              onChange={(e) => setForm({ ...form, staff_id: e.target.value })} required>
              <option value="">선택</option>
              {staffList.map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
            </select>
          </div>

          {/* Item & Category */}
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label htmlFor="sale-service" className="block text-sm text-dark-muted mb-1">시술 메뉴</label>
              <select id="sale-service" className="glass-input w-full" value={form.service_id}
                onChange={(e) => handleServiceChange(e.target.value)}>
                <option value="">메뉴 선택 (또는 직접 입력)</option>
                {serviceList.filter(s => s.is_active).map((s) => (
                  <option key={s.id} value={s.id}>{s.category} - {s.name} ({formatCurrency(s.price)})</option>
                ))}
              </select>
            </div>
            {!form.service_id && (
              <div>
                <label htmlFor="sale-item" className="block text-sm text-dark-muted mb-1">항목명 *</label>
                <input id="sale-item" className="glass-input w-full" value={form.item_name}
                  onChange={(e) => setForm({ ...form, item_name: e.target.value })} placeholder="직접 입력" required />
              </div>
            )}
          </div>
          
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-sm text-dark-muted mb-1">분류 *</label>
              <div className="flex gap-2">
                <button onClick={() => setForm({ ...form, category: 'service' })}
                  className={`flex-1 py-3 rounded-xl text-sm font-medium transition-all ${
                    form.category === 'service' ? 'bg-green-500/20 text-green-400 border border-green-500/30' : 'bg-dark-surface text-dark-muted'}`}>
                  시술
                </button>
                <button onClick={() => setForm({ ...form, category: 'product' })}
                  className={`flex-1 py-3 rounded-xl text-sm font-medium transition-all ${
                    form.category === 'product' ? 'bg-blue-500/20 text-blue-400 border border-blue-500/30' : 'bg-dark-surface text-dark-muted'}`}>
                  점판
                </button>
              </div>
            </div>
          </div>

          {/* Amount */}
          <div>
            <label htmlFor="sale-amount" className="block text-sm text-dark-muted mb-1">총 금액 *</label>
            <input id="sale-amount" type="number" className="glass-input w-full text-lg font-bold" value={form.total_amount}
              onChange={(e) => setForm({ ...form, total_amount: e.target.value })} placeholder="0" required />
          </div>

          {/* Payment Method */}
          <div>
            <label className="block text-sm text-dark-muted mb-2">결제 수단 *</label>
            <div className="grid grid-cols-4 gap-2">
              {(['card', 'cash', 'membership', 'mixed'] as const).map((method) => (
                <button key={method}
                  onClick={() => handlePaymentMethodChange(method)}
                  className={`py-3 rounded-xl text-sm font-medium transition-all ${
                    form.payment_method === method ? 'bg-salon-500/20 text-salon-400 border border-salon-500/30' : 'bg-dark-surface text-dark-muted'}`}>
                  {getPaymentLabel(method)}
                </button>
              ))}
            </div>
          </div>

          {/* Mixed payment breakdown */}
          {form.payment_method === 'mixed' && (
            <div className="grid grid-cols-3 gap-3 bg-dark-surface/50 rounded-xl p-3">
              <div>
                <label htmlFor="sale-card-amt" className="text-xs text-dark-muted">카드</label>
                <input id="sale-card-amt" type="number" className="glass-input w-full text-sm mt-1" value={form.card_amount}
                  onChange={(e) => setForm({ ...form, card_amount: e.target.value })} />
              </div>
              <div>
                <label htmlFor="sale-cash-amt" className="text-xs text-dark-muted">현금</label>
                <input id="sale-cash-amt" type="number" className="glass-input w-full text-sm mt-1" value={form.cash_amount}
                  onChange={(e) => setForm({ ...form, cash_amount: e.target.value })} />
              </div>
              <div>
                <label htmlFor="sale-member-amt" className="text-xs text-dark-muted">정액권</label>
                <input id="sale-member-amt" type="number" className="glass-input w-full text-sm mt-1" value={form.membership_amount}
                  onChange={(e) => setForm({ ...form, membership_amount: e.target.value })} />
              </div>
            </div>
          )}

          {/* Card commission */}
          {(form.payment_method === 'card' || form.payment_method === 'mixed') && (
            <div>
              <label htmlFor="sale-commission" className="block text-sm text-dark-muted mb-1">카드 수수료율 (%)</label>
              <input id="sale-commission" type="number" step="0.1" className="glass-input w-full text-sm" value={form.card_commission_rate}
                onChange={(e) => setForm({ ...form, card_commission_rate: e.target.value })} />
            </div>
          )}

          {/* Memo */}
          <div>
            <label htmlFor="sale-memo" className="block text-sm text-dark-muted mb-1">메모</label>
            <input id="sale-memo" className="glass-input w-full text-sm" value={form.memo}
              onChange={(e) => setForm({ ...form, memo: e.target.value })} />
          </div>
        </div>

        <div className="flex gap-3 mt-6">
          <button onClick={onClose} className="btn-secondary flex-1">취소</button>
          <button id="sale-submit" onClick={() => {
            const amount = parseFloat(form.total_amount) || 0;
            onSubmit({
              staff_id: form.staff_id,
              service_id: form.service_id || undefined,
              item_name: form.item_name,
              total_amount: amount,
              category: form.category,
              payment_method: form.payment_method,
              card_amount: form.payment_method === 'card' ? amount : parseFloat(form.card_amount) || 0,
              cash_amount: form.payment_method === 'cash' ? amount : parseFloat(form.cash_amount) || 0,
              membership_amount: form.payment_method === 'membership' ? amount : parseFloat(form.membership_amount) || 0,
              card_commission_rate: parseFloat(form.card_commission_rate) || 0,
              membership_id: form.membership_id || undefined,
              memo: form.memo,
            });
          }} className="btn-primary flex-1">
            💰 매출 등록
          </button>
        </div>
      </div>
    </div>
  );
}
