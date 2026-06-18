// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';
import api from '@/lib/api';
import { formatCurrency } from '@/lib/utils';
import type { Staff, StaffPerformance } from '@/types';

export default function StaffPage() {
  const [staffList, setStaffList] = useState<Staff[]>([]);
  const [performance, setPerformance] = useState<StaffPerformance[]>([]);
  const [showCreateModal, setShowCreateModal] = useState(false);

  const now = new Date();
  const monthStart = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`;
  const today = now.toISOString().split('T')[0];

  const fetchData = useCallback(async () => {
    try {
      const [staffs, perf] = await Promise.all([
        api.get<Staff[]>('/api/staffs'),
        api.get<StaffPerformance[]>(`/api/analytics/staff-performance?start_date=${monthStart}&end_date=${today}`),
      ]);
      setStaffList(staffs || []);
      setPerformance(perf || []);
    } catch (err) {
      console.error('Failed to fetch staff data:', err);
      // Removed setStaffList([]) and setPerformance([]) to avoid wiping data on temporary errors
    }
  }, [monthStart, today]);

  useEffect(() => { fetchData(); }, [fetchData]);

  const handleCreateStaff = async (form: Record<string, string>) => {
    try {
      await api.post('/api/auth/register', {
        ...form,
        service_incentive_rate: parseFloat(form.service_incentive_rate) || 0,
        product_incentive_rate: parseFloat(form.product_incentive_rate) || 0,
        base_salary: parseFloat(form.base_salary) || 0,
        monthly_target: parseFloat(form.monthly_target) || 0,
      });
      setShowCreateModal(false);
      fetchData();
    } catch (err: any) {
      alert(`직원 등록 실패: ${err.message || '알 수 없는 오류가 발생했습니다.'}`);
    }
  };

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">👥 직원 관리 및 급여 정산</h1>
        <button id="create-staff-btn" onClick={() => setShowCreateModal(true)} className="btn-primary text-sm py-2">
          + 직원 등록
        </button>
      </div>

      {/* Performance Gauges */}
      <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {performance.map((p) => (
          <div key={p.staff_id} className="glass-card p-5">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-salon-400 to-purple-600 flex items-center justify-center text-white text-sm font-bold">
                  {p.staff_name.charAt(0)}
                </div>
                <div>
                  <p className="font-semibold text-white">{p.staff_name}</p>
                  <p className="text-xs text-dark-muted">목표: {formatCurrency(p.monthly_target)}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-lg font-bold text-white">{Math.round(p.achievement_rate)}%</p>
                <p className="text-xs text-dark-muted">달성률</p>
              </div>
            </div>

            {/* Achievement Gauge */}
            <div className="gauge-track h-3 mb-4">
              <div className="gauge-fill" style={{ width: `${Math.min(p.achievement_rate, 100)}%` }} />
            </div>

            <div className="grid grid-cols-2 gap-3 text-sm">
              <div className="bg-dark-surface/50 rounded-lg p-2">
                <p className="text-xs text-dark-muted">시술 매출</p>
                <p className="font-semibold text-green-400">{formatCurrency(p.service_revenue)}</p>
              </div>
              <div className="bg-dark-surface/50 rounded-lg p-2">
                <p className="text-xs text-dark-muted">점판 매출</p>
                <p className="font-semibold text-blue-400">{formatCurrency(p.product_revenue)}</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Payroll Table */}
      <div className="glass-card overflow-hidden">
        <div className="p-4 border-b border-dark-border">
          <h3 className="font-semibold text-white">💰 실질 급여 정산 (카드 수수료 공제 후)</h3>
          <p className="text-xs text-dark-muted mt-1">
            {new Date(monthStart).toLocaleDateString('ko-KR', { year: 'numeric', month: 'long' })} 기준
          </p>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-dark-surface/50">
              <tr>
                <th className="table-header text-left">직원명</th>
                <th className="table-header text-right">총 매출</th>
                <th className="table-header text-right">순매출 (시술)</th>
                <th className="table-header text-right">순매출 (점판)</th>
                <th className="table-header text-right">시술 인센티브</th>
                <th className="table-header text-right">점판 인센티브</th>
                <th className="table-header text-right">기본급</th>
                <th className="table-header text-right font-bold">총 급여</th>
              </tr>
            </thead>
            <tbody>
              {performance.length === 0 ? (
                <tr><td colSpan={8} className="table-cell text-center text-dark-muted py-8">데이터가 없습니다</td></tr>
              ) : performance.map((p) => (
                <tr key={p.staff_id} className="table-row">
                  <td className="table-cell font-medium text-white">
                    <div className="flex items-center gap-2">
                      <div className="w-7 h-7 rounded-lg bg-gradient-to-br from-salon-400/30 to-purple-600/30 flex items-center justify-center text-salon-400 text-xs font-bold">
                        {p.staff_name.charAt(0)}
                      </div>
                      {p.staff_name}
                    </div>
                  </td>
                  <td className="table-cell text-right">{formatCurrency(p.total_revenue)}</td>
                  <td className="table-cell text-right text-green-400">{formatCurrency(p.net_service_revenue)}</td>
                  <td className="table-cell text-right text-blue-400">{formatCurrency(p.net_product_revenue)}</td>
                  <td className="table-cell text-right">{formatCurrency(p.service_incentive)}</td>
                  <td className="table-cell text-right">{formatCurrency(p.product_incentive)}</td>
                  <td className="table-cell text-right text-dark-muted">{formatCurrency(p.base_salary)}</td>
                  <td className="table-cell text-right font-bold text-salon-400">{formatCurrency(p.total_payroll)}</td>
                </tr>
              ))}
            </tbody>
            {performance.length > 0 && (
              <tfoot className="bg-dark-surface/30 border-t border-dark-border">
                <tr>
                  <td className="table-cell font-bold text-white">합계</td>
                  <td className="table-cell text-right font-bold">{formatCurrency(performance.reduce((s, p) => s + p.total_revenue, 0))}</td>
                  <td className="table-cell text-right font-bold text-green-400">{formatCurrency(performance.reduce((s, p) => s + p.net_service_revenue, 0))}</td>
                  <td className="table-cell text-right font-bold text-blue-400">{formatCurrency(performance.reduce((s, p) => s + p.net_product_revenue, 0))}</td>
                  <td className="table-cell text-right font-bold">{formatCurrency(performance.reduce((s, p) => s + p.service_incentive, 0))}</td>
                  <td className="table-cell text-right font-bold">{formatCurrency(performance.reduce((s, p) => s + p.product_incentive, 0))}</td>
                  <td className="table-cell text-right font-bold">{formatCurrency(performance.reduce((s, p) => s + p.base_salary, 0))}</td>
                  <td className="table-cell text-right font-bold text-salon-400">{formatCurrency(performance.reduce((s, p) => s + p.total_payroll, 0))}</td>
                </tr>
              </tfoot>
            )}
          </table>
        </div>
      </div>

      {/* Staff List */}
      <div className="glass-card overflow-hidden">
        <div className="p-4 border-b border-dark-border">
          <h3 className="font-semibold text-white">📋 직원 목록</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-dark-surface/50">
              <tr>
                <th className="table-header text-left">이름</th>
                <th className="table-header text-left">직급</th>
                <th className="table-header text-right">시술 수수료율</th>
                <th className="table-header text-right">점판 수수료율</th>
                <th className="table-header text-right">기본급</th>
                <th className="table-header text-right">월 목표</th>
              </tr>
            </thead>
            <tbody>
              {staffList.map((s) => (
                <tr key={s.id} className="table-row">
                  <td className="table-cell font-medium text-white">{s.name}</td>
                  <td className="table-cell">
                    <span className={`badge ${s.role === 'admin' ? 'bg-red-500/20 text-red-400' : s.role === 'designer' ? 'bg-purple-500/20 text-purple-400' : 'bg-dark-surface text-dark-muted'}`}>
                      {s.role === 'admin' ? '관리자' : s.role === 'designer' ? '디자이너' : '스탭'}
                    </span>
                  </td>
                  <td className="table-cell text-right">{s.service_incentive_rate}%</td>
                  <td className="table-cell text-right">{s.product_incentive_rate}%</td>
                  <td className="table-cell text-right">{formatCurrency(s.base_salary)}</td>
                  <td className="table-cell text-right">{formatCurrency(s.monthly_target)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Create Staff Modal */}
      {showCreateModal && (
        <CreateStaffModal onClose={() => setShowCreateModal(false)} onSubmit={handleCreateStaff} />
      )}
    </div>
  );
}

function CreateStaffModal({ onClose, onSubmit }: { onClose: () => void; onSubmit: (data: Record<string, string>) => void }) {
  const [form, setForm] = useState({
    name: '', email: '', password: '', role: 'designer',
    phone: '', service_incentive_rate: '30', product_incentive_rate: '10',
    base_salary: '2000000', monthly_target: '5000000',
  });

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content max-w-lg" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-xl font-bold text-white mb-6">👤 직원 등록</h2>
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label htmlFor="staff-name" className="block text-sm text-dark-muted mb-1">이름 *</label>
              <input id="staff-name" className="glass-input w-full" value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })} required />
            </div>
            <div>
              <label htmlFor="staff-role" className="block text-sm text-dark-muted mb-1">직급 *</label>
              <select id="staff-role" className="glass-input w-full" value={form.role}
                onChange={(e) => setForm({ ...form, role: e.target.value })}>
                <option value="admin">관리자</option>
                <option value="designer">디자이너</option>
                <option value="staff">스탭</option>
              </select>
            </div>
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label htmlFor="staff-email" className="block text-sm text-dark-muted mb-1">이메일 *</label>
              <input id="staff-email" type="email" className="glass-input w-full" value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })} required />
            </div>
            <div>
              <label htmlFor="staff-phone" className="block text-sm text-dark-muted mb-1">연락처</label>
              <input id="staff-phone" type="tel" className="glass-input w-full" value={form.phone}
                onChange={(e) => setForm({ ...form, phone: e.target.value })} placeholder="010-0000-0000" />
            </div>
          </div>
          <div>
            <label htmlFor="staff-password" className="block text-sm text-dark-muted mb-1">비밀번호 * (8자 이상)</label>
            <input id="staff-password" type="password" className="glass-input w-full" value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })} required minLength={8} autoComplete="new-password" />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label htmlFor="staff-service-rate" className="block text-sm text-dark-muted mb-1">시술 수수료율 (%)</label>
              <input id="staff-service-rate" type="number" className="glass-input w-full" value={form.service_incentive_rate}
                onChange={(e) => setForm({ ...form, service_incentive_rate: e.target.value })} />
            </div>
            <div>
              <label htmlFor="staff-product-rate" className="block text-sm text-dark-muted mb-1">점판 수수료율 (%)</label>
              <input id="staff-product-rate" type="number" className="glass-input w-full" value={form.product_incentive_rate}
                onChange={(e) => setForm({ ...form, product_incentive_rate: e.target.value })} />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label htmlFor="staff-salary" className="block text-sm text-dark-muted mb-1">기본급 (원)</label>
              <input id="staff-salary" type="number" className="glass-input w-full" value={form.base_salary}
                onChange={(e) => setForm({ ...form, base_salary: e.target.value })} />
            </div>
            <div>
              <label htmlFor="staff-target" className="block text-sm text-dark-muted mb-1">월 목표 (원)</label>
              <input id="staff-target" type="number" className="glass-input w-full" value={form.monthly_target}
                onChange={(e) => setForm({ ...form, monthly_target: e.target.value })} />
            </div>
          </div>
        </div>
        <div className="flex gap-3 mt-6">
          <button onClick={onClose} className="btn-secondary flex-1">취소</button>
          <button id="staff-submit" onClick={() => onSubmit(form)} className="btn-primary flex-1">등록</button>
        </div>
      </div>
    </div>
  );
}
