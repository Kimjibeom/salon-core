// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';
import { AreaChart, Area, BarChart, Bar, PieChart, Pie, Cell, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from 'recharts';
import api from '@/lib/api';
import { formatCurrency } from '@/lib/utils';
import type { DailySummary, AnalyticsSummary, PaymentBreakdown, ChurnAnalysis } from '@/types';

const CHART_COLORS = ['#ec4899', '#8b5cf6', '#06b6d4', '#10b981', '#f59e0b', '#ef4444'];
const PIE_COLORS = ['#ec4899', '#8b5cf6', '#06b6d4'];

export default function AnalyticsPage() {
  const [period, setPeriod] = useState<'daily' | 'monthly' | 'trend'>('daily');
  const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0]);
  const [dailySummary, setDailySummary] = useState<DailySummary | null>(null);
  const [monthlyTrend, setMonthlyTrend] = useState<DailySummary[]>([]);
  const [analyticsSummary, setAnalyticsSummary] = useState<AnalyticsSummary[]>([]);
  const [paymentBreakdown, setPaymentBreakdown] = useState<PaymentBreakdown[]>([]);
  const [customerStats, setCustomerStats] = useState<{ new_customers: number; returning_customers: number; total: number } | null>(null);
  const [churnData, setChurnData] = useState<ChurnAnalysis[]>([]);

  const now = new Date();
  const monthStart = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`;
  const monthEnd = selectedDate;

  const fetchData = useCallback(async () => {
    try {
      const [daily, payment, customers, churn] = await Promise.all([
        api.get<DailySummary>(`/api/analytics/daily?date=${selectedDate}`),
        api.get<PaymentBreakdown[]>(`/api/analytics/payments?start_date=${monthStart}&end_date=${monthEnd}`),
        api.get<{ new_customers: number; returning_customers: number; total: number }>(`/api/analytics/customers?start_date=${monthStart}&end_date=${monthEnd}`),
        api.get<ChurnAnalysis[]>('/api/analytics/churn?inactive_days=90'),
      ]);
      setDailySummary(daily);
      setPaymentBreakdown(payment || []);
      setCustomerStats(customers);
      setChurnData(churn || []);
    } catch { /* error handled */ }

    try {
      const monthly = await api.get<DailySummary[]>(`/api/analytics/monthly?year=${now.getFullYear()}&month=${now.getMonth() + 1}`);
      setMonthlyTrend(monthly || []);
    } catch { setMonthlyTrend([]); }

    try {
      const summary = await api.get<AnalyticsSummary[]>('/api/analytics/summary?months=6');
      setAnalyticsSummary(summary || []);
    } catch { setAnalyticsSummary([]); }
  }, [selectedDate, monthStart, monthEnd]);

  useEffect(() => { fetchData(); }, [fetchData]);

  const paymentMethodLabels: Record<string, string> = { card: '카드', cash: '현금', membership: '정액권' };

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <h1 className="text-2xl font-bold text-white">📊 매출 분석 대시보드</h1>
        <input type="date" id="analytics-date" value={selectedDate}
          onChange={(e) => setSelectedDate(e.target.value)} className="glass-input text-sm py-2" />
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="stat-card">
          <span className="stat-label">일 매출</span>
          <span className="stat-value">{formatCurrency(dailySummary?.total_revenue || 0)}</span>
          <span className="text-xs text-dark-muted">{dailySummary?.transaction_count || 0}건</span>
        </div>
        <div className="stat-card">
          <span className="stat-label">시술 매출</span>
          <span className="text-xl font-bold text-green-400">{formatCurrency(dailySummary?.service_revenue || 0)}</span>
          <span className="text-xs text-dark-muted">
            객단가: {formatCurrency(dailySummary?.transaction_count ? (dailySummary.service_revenue / Math.max(dailySummary.customer_count, 1)) : 0)}
          </span>
        </div>
        <div className="stat-card">
          <span className="stat-label">점판 매출</span>
          <span className="text-xl font-bold text-blue-400">{formatCurrency(dailySummary?.product_revenue || 0)}</span>
          <span className="text-xs text-dark-muted">
            객단가: {formatCurrency(dailySummary?.transaction_count ? (dailySummary.product_revenue / Math.max(dailySummary.customer_count, 1)) : 0)}
          </span>
        </div>
        <div className="stat-card">
          <span className="stat-label">고객 수</span>
          <span className="text-xl font-bold text-purple-400">{dailySummary?.customer_count || 0}명</span>
        </div>
      </div>

      <div className="grid lg:grid-cols-2 gap-6">
        {/* Monthly Revenue Trend */}
        <div className="glass-card p-6">
          <h3 className="text-lg font-semibold text-white mb-4">📈 월간 매출 트렌드</h3>
          <ResponsiveContainer width="100%" height={280}>
            <AreaChart data={monthlyTrend}>
              <defs>
                <linearGradient id="colorRevenue" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#ec4899" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#ec4899" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#2d2d50" />
              <XAxis dataKey="date" tick={{ fill: '#8888a8', fontSize: 11 }} tickFormatter={(v) => v.slice(8)} />
              <YAxis tick={{ fill: '#8888a8', fontSize: 11 }} tickFormatter={(v) => `${(v / 10000).toFixed(0)}만`} />
              <Tooltip
                contentStyle={{ background: '#1a1a2e', border: '1px solid #2d2d50', borderRadius: '12px', color: '#e4e4f0' }}
                formatter={(value: number) => [formatCurrency(value), '매출']}
              />
              <Area type="monotone" dataKey="total_revenue" stroke="#ec4899" strokeWidth={2} fill="url(#colorRevenue)" />
            </AreaChart>
          </ResponsiveContainer>
        </div>

        {/* 6-Month Trend */}
        <div className="glass-card p-6">
          <h3 className="text-lg font-semibold text-white mb-4">📊 6개월 매출 변화</h3>
          <ResponsiveContainer width="100%" height={280}>
            <BarChart data={analyticsSummary}>
              <CartesianGrid strokeDasharray="3 3" stroke="#2d2d50" />
              <XAxis dataKey="period" tick={{ fill: '#8888a8', fontSize: 11 }} tickFormatter={(v) => v.slice(5)} />
              <YAxis tick={{ fill: '#8888a8', fontSize: 11 }} tickFormatter={(v) => `${(v / 10000).toFixed(0)}만`} />
              <Tooltip
                contentStyle={{ background: '#1a1a2e', border: '1px solid #2d2d50', borderRadius: '12px', color: '#e4e4f0' }}
                formatter={(value: number) => [formatCurrency(value)]}
              />
              <Legend />
              <Bar dataKey="service_revenue" name="시술" fill="#10b981" radius={[4, 4, 0, 0]} />
              <Bar dataKey="product_revenue" name="점판" fill="#06b6d4" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* Payment Method Pie Chart */}
        <div className="glass-card p-6">
          <h3 className="text-lg font-semibold text-white mb-4">💳 결제 수단별 비중</h3>
          <ResponsiveContainer width="100%" height={280}>
            <PieChart>
              <Pie
                data={paymentBreakdown}
                cx="50%" cy="50%"
                innerRadius={60}
                outerRadius={100}
                dataKey="amount"
                nameKey="method"
                label={({ name, ratio }) => `${paymentMethodLabels[name] || name} ${ratio}%`}
                labelLine={false}
              >
                {paymentBreakdown.map((_, index) => (
                  <Cell key={`cell-${index}`} fill={PIE_COLORS[index % PIE_COLORS.length]} />
                ))}
              </Pie>
              <Tooltip
                contentStyle={{ background: '#1a1a2e', border: '1px solid #2d2d50', borderRadius: '12px', color: '#e4e4f0' }}
                formatter={(value: number) => [formatCurrency(value)]}
              />
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* New vs Returning */}
        <div className="glass-card p-6">
          <h3 className="text-lg font-semibold text-white mb-4">🔄 신규 vs 재방문</h3>
          {customerStats && (
            <div className="space-y-6">
              <div className="flex items-center justify-center gap-8">
                <div className="text-center">
                  <p className="text-3xl font-bold text-salon-400">{customerStats.new_customers}</p>
                  <p className="text-sm text-dark-muted mt-1">신규 고객</p>
                </div>
                <div className="text-4xl text-dark-muted">vs</div>
                <div className="text-center">
                  <p className="text-3xl font-bold text-purple-400">{customerStats.returning_customers}</p>
                  <p className="text-sm text-dark-muted mt-1">재방문 고객</p>
                </div>
              </div>
              {customerStats.total > 0 && (
                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span className="text-dark-muted">재방문율</span>
                    <span className="text-white font-semibold">
                      {((customerStats.returning_customers / customerStats.total) * 100).toFixed(1)}%
                    </span>
                  </div>
                  <div className="gauge-track h-3">
                    <div className="gauge-fill" style={{ width: `${(customerStats.returning_customers / customerStats.total) * 100}%` }} />
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Churn Analysis Table */}
      <div className="glass-card overflow-hidden">
        <div className="p-4 border-b border-dark-border">
          <h3 className="font-semibold text-white">📉 이탈 고객 분석 (방문 회차별)</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-dark-surface/50">
              <tr>
                <th className="table-header text-left">방문 회차</th>
                <th className="table-header text-right">전체 고객</th>
                <th className="table-header text-right">이탈 고객</th>
                <th className="table-header text-right">이탈률</th>
                <th className="table-header text-left">그래프</th>
              </tr>
            </thead>
            <tbody>
              {churnData.length === 0 ? (
                <tr><td colSpan={5} className="table-cell text-center text-dark-muted py-8">데이터가 없습니다</td></tr>
              ) : churnData.map((row) => (
                <tr key={row.visit_number} className="table-row">
                  <td className="table-cell font-medium text-white">{row.visit_number}회차</td>
                  <td className="table-cell text-right text-dark-muted">{row.total_customers}명</td>
                  <td className="table-cell text-right text-red-400">{row.churned_customers}명</td>
                  <td className="table-cell text-right font-semibold text-white">{row.churn_rate}%</td>
                  <td className="table-cell">
                    <div className="gauge-track h-2 w-32">
                      <div className="h-full rounded-full bg-gradient-to-r from-green-500 to-red-500 transition-all" style={{ width: `${Math.min(row.churn_rate, 100)}%` }} />
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
