// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect } from 'react';
import api from '@/lib/api';
import { formatCurrency, getPaymentLabel } from '@/lib/utils';
import type { Staff, ServiceMenu, Customer, Membership } from '@/types';

export default function SaleEntryModal({ staffList, serviceList, onClose, onSubmit, initialData }: {
  staffList: Staff[];
  serviceList: ServiceMenu[];
  onClose: () => void;
  onSubmit: (data: Record<string, unknown>) => void;
  initialData?: {
    staff_id?: string;
    service_id?: string;
    customer_id?: string;
    item_name?: string;
    total_amount?: string;
    memo?: string;
  };
}) {
  const [form, setForm] = useState({
    staff_id: initialData?.staff_id || '',
    service_id: initialData?.service_id || '',
    customer_id: initialData?.customer_id || '',
    item_name: initialData?.item_name || '',
    total_amount: initialData?.total_amount || '',
    category: 'service' as 'service' | 'product',
    payment_method: 'card' as 'card' | 'cash' | 'membership' | 'mixed',
    card_amount: '',
    cash_amount: '',
    membership_amount: '',
    card_commission_rate: '2.5',
    membership_id: '',
    memo: initialData?.memo || '',
  });
  const [customerList, setCustomerList] = useState<Customer[]>([]);
  const [membershipList, setMembershipList] = useState<Membership[]>([]);
  const [error, setError] = useState('');

  const totalAmount = parseFloat(form.total_amount) || 0;
  const usesMembership = form.payment_method === 'membership' || form.payment_method === 'mixed';

  useEffect(() => {
    // If initialData has a service_id, try to figure out the category
    if (initialData?.service_id) {
      const service = serviceList.find(s => s.id === initialData.service_id);
      if (service) {
        setForm(prev => ({ ...prev, category: service.category === '점판' ? 'product' : 'service' }));
      }
    }
  }, [initialData, serviceList]);

  useEffect(() => {
    api.get<{ data: Customer[]; total: number }>('/api/customers?limit=100')
      .then((res) => setCustomerList(res?.data || []))
      .catch(() => setCustomerList([]));
  }, []);

  // Load the selected customer's active memberships (needed for balance deduction)
  useEffect(() => {
    if (!form.customer_id) {
      setMembershipList([]);
      setForm(prev => ({ ...prev, membership_id: '' }));
      return;
    }
    api.get<Membership[]>(`/api/memberships/customer/${form.customer_id}`)
      .then((data) => setMembershipList((data || []).filter(m => m.is_active)))
      .catch(() => setMembershipList([]));
  }, [form.customer_id]);

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
        <h2 className="text-xl font-bold text-white mb-6">💳 영업 (결제) 입력</h2>
        {error && (
          <div className="mb-4 p-3 bg-red-500/20 border border-red-500/50 rounded-lg">
            <p className="text-sm text-red-400">{error}</p>
          </div>
        )}
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

          {/* Customer Selection (required for membership payment) */}
          <div>
            <label htmlFor="sale-customer" className="block text-sm text-dark-muted mb-1">
              고객 {usesMembership ? '* (정액권 결제 시 필수)' : '(선택)'}
            </label>
            <select id="sale-customer" className="glass-input w-full" value={form.customer_id}
              onChange={(e) => setForm({ ...form, customer_id: e.target.value, membership_id: '' })}>
              <option value="">선택 안함 (비회원)</option>
              {customerList.map((c) => <option key={c.id} value={c.id}>{c.name} ({c.phone})</option>)}
            </select>
          </div>

          {/* Item & Category */}
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label htmlFor="sale-service" className="block text-sm text-dark-muted mb-1">시술 메뉴</label>
              <select id="sale-service" className="glass-input w-full" value={form.service_id}
                onChange={(e) => handleServiceChange(e.target.value)}>
                <option value="">메뉴 선택 (또는 직접 입력)</option>
                {serviceList.filter(s => s.is_active || s.id === form.service_id).map((s) => (
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

          {/* Membership selection for deduction */}
          {usesMembership && (
            <div>
              <label htmlFor="sale-membership" className="block text-sm text-dark-muted mb-1">사용할 회원권 *</label>
              {form.customer_id ? (
                membershipList.length > 0 ? (
                  <select id="sale-membership" className="glass-input w-full" value={form.membership_id}
                    onChange={(e) => setForm({ ...form, membership_id: e.target.value })}>
                    <option value="">선택</option>
                    {membershipList.map((m) => (
                      <option key={m.id} value={m.id}>
                        {m.name} — {m.type === 'money' ? `잔액 ${formatCurrency(m.balance)}` : `${m.remaining_count}회 남음`}
                      </option>
                    ))}
                  </select>
                ) : (
                  <p className="text-sm text-amber-400 bg-amber-500/10 rounded-lg p-3">해당 고객에게 사용 가능한 회원권이 없습니다.</p>
                )
              ) : (
                <p className="text-sm text-amber-400 bg-amber-500/10 rounded-lg p-3">먼저 고객을 선택해주세요.</p>
              )}
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
            setError('');
            const amount = parseFloat(form.total_amount) || 0;
            if (!form.staff_id) { setError('담당 직원을 선택해주세요.'); return; }
            if (!form.item_name) { setError('항목명을 입력해주세요.'); return; }
            if (amount <= 0) { setError('총 금액을 입력해주세요.'); return; }

            const cardAmount = form.payment_method === 'card' ? amount : form.payment_method === 'mixed' ? parseFloat(form.card_amount) || 0 : 0;
            const cashAmount = form.payment_method === 'cash' ? amount : form.payment_method === 'mixed' ? parseFloat(form.cash_amount) || 0 : 0;
            const membershipAmount = form.payment_method === 'membership' ? amount : form.payment_method === 'mixed' ? parseFloat(form.membership_amount) || 0 : 0;

            if (form.payment_method === 'mixed' && cardAmount + cashAmount + membershipAmount !== amount) {
              setError('복합 결제 금액의 합이 총 금액과 일치해야 합니다.');
              return;
            }
            if (membershipAmount > 0 && !form.membership_id) {
              setError('정액권 결제 시 고객과 사용할 회원권을 선택해주세요.');
              return;
            }

            onSubmit({
              staff_id: form.staff_id,
              service_id: form.service_id || undefined,
              customer_id: form.customer_id || undefined,
              item_name: form.item_name,
              total_amount: amount,
              category: form.category,
              payment_method: form.payment_method,
              card_amount: cardAmount,
              cash_amount: cashAmount,
              membership_amount: membershipAmount,
              card_commission_rate: parseFloat(form.card_commission_rate) || 0,
              membership_id: form.membership_id || undefined,
              memo: form.memo,
            });
          }} className="btn-primary flex-1">
            💰 결제 / 매출 등록
          </button>
        </div>
      </div>
    </div>
  );
}
