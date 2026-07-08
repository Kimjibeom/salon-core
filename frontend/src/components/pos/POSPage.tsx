// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';
import api from '@/lib/api';
import { formatCurrency, getCategoryLabel, getPaymentLabel } from '@/lib/utils';
import type { Staff, Membership, Sale, ServiceMenu } from '@/types';
import SaleEntryModal from './SaleEntryModal';

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
    } catch (err: any) {
      alert(err.message || '매출 등록에 실패했습니다.');
    }
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
