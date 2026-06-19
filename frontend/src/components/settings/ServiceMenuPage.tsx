// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';
import api from '@/lib/api';
import { formatCurrency } from '@/lib/utils';
import type { ServiceMenu } from '@/types';

export default function ServiceMenuPage() {
  const [menus, setMenus] = useState<ServiceMenu[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    try {
      const data = await api.get<ServiceMenu[]>('/api/services');
      setMenus(data || []);
    } catch (err) {
      console.error('Failed to fetch menus:', err);
    }
  }, []);

  useEffect(() => { fetchData(); }, [fetchData]);

  const handleSubmit = async (form: Partial<ServiceMenu>) => {
    try {
      if (editingId) {
        await api.put(`/api/services/${editingId}`, form);
      } else {
        await api.post('/api/services', form);
      }
      setShowModal(false);
      fetchData();
    } catch (err: any) {
      alert(`저장 실패: ${err.message}`);
    }
  };

  const handleDelete = async (id: string, name: string) => {
    if (!window.confirm(`정말 '${name}' 메뉴를 삭제하시겠습니까?`)) return;
    try {
      await api.delete(`/api/services/${id}`);
      fetchData();
    } catch (err: any) {
      alert(`삭제 실패: ${err.message}`);
    }
  };

  const openEdit = (menu: ServiceMenu) => {
    setEditingId(menu.id);
    setShowModal(true);
  };

  const openCreate = () => {
    setEditingId(null);
    setShowModal(true);
  };

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">📋 시술 메뉴판 관리</h1>
        <button onClick={openCreate} className="btn-primary text-sm py-2">+ 메뉴 추가</button>
      </div>

      <div className="glass-card overflow-hidden">
        <div className="p-4 border-b border-dark-border">
          <h3 className="font-semibold text-white">메뉴 목록</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-dark-surface/50">
              <tr>
                <th className="table-header text-left">카테고리</th>
                <th className="table-header text-left">메뉴명</th>
                <th className="table-header text-right">기본 가격</th>
                <th className="table-header text-right">소요 시간(분)</th>
                <th className="table-header text-center">상태</th>
                <th className="table-header text-center">관리</th>
              </tr>
            </thead>
            <tbody>
              {menus.map((m) => (
                <tr key={m.id} className="table-row">
                  <td className="table-cell font-medium text-dark-muted">{m.category}</td>
                  <td className="table-cell font-bold text-white">{m.name}</td>
                  <td className="table-cell text-right text-salon-400">{formatCurrency(m.price)}</td>
                  <td className="table-cell text-right">{m.duration}분</td>
                  <td className="table-cell text-center">
                    <span className={`badge ${m.is_active ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'}`}>
                      {m.is_active ? '사용중' : '숨김'}
                    </span>
                  </td>
                  <td className="table-cell text-center">
                    <div className="flex justify-center gap-2">
                      <button onClick={() => openEdit(m)} className="px-3 py-1 bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 rounded-md text-xs transition-colors">수정</button>
                      <button onClick={() => handleDelete(m.id, m.name)} className="px-3 py-1 bg-red-500/10 hover:bg-red-500/20 text-red-500 rounded-md text-xs transition-colors">삭제</button>
                    </div>
                  </td>
                </tr>
              ))}
              {menus.length === 0 && (
                <tr><td colSpan={6} className="table-cell text-center text-dark-muted py-8">등록된 메뉴가 없습니다</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {showModal && (
        <MenuModal
          onClose={() => setShowModal(false)}
          onSubmit={handleSubmit}
          initialData={editingId ? menus.find(m => m.id === editingId) : undefined}
        />
      )}
    </div>
  );
}

function MenuModal({ onClose, onSubmit, initialData }: { onClose: () => void; onSubmit: (data: Partial<ServiceMenu>) => void; initialData?: ServiceMenu }) {
  const [form, setForm] = useState({
    category: initialData?.category || '커트',
    name: initialData?.name || '',
    price: initialData?.price?.toString() || '0',
    duration: initialData?.duration?.toString() || '30',
    is_active: initialData ? initialData.is_active : true,
  });

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content max-w-md" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-xl font-bold text-white mb-6">📋 {initialData ? '메뉴 수정' : '메뉴 추가'}</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm text-dark-muted mb-1">카테고리 *</label>
            <select className="glass-input w-full" value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value })}>
              <option value="커트">커트</option>
              <option value="펌">펌</option>
              <option value="염색">염색</option>
              <option value="클리닉">클리닉</option>
              <option value="드라이">드라이</option>
              <option value="점판">점판 (제품)</option>
              <option value="기타">기타</option>
            </select>
          </div>
          <div>
            <label className="block text-sm text-dark-muted mb-1">메뉴명 *</label>
            <input className="glass-input w-full" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required placeholder="예: 남성 컷트" />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-sm text-dark-muted mb-1">기본 가격 (원) *</label>
              <input type="number" className="glass-input w-full" value={form.price} onChange={(e) => setForm({ ...form, price: e.target.value })} required />
            </div>
            <div>
              <label className="block text-sm text-dark-muted mb-1">소요 시간 (분) *</label>
              <input type="number" className="glass-input w-full" value={form.duration} onChange={(e) => setForm({ ...form, duration: e.target.value })} required />
            </div>
          </div>
          <div className="flex items-center gap-2 mt-2">
            <input type="checkbox" id="is_active" checked={form.is_active} onChange={(e) => setForm({ ...form, is_active: e.target.checked })} className="rounded bg-dark-bg border-dark-border" />
            <label htmlFor="is_active" className="text-sm text-white">메뉴판에 표시 (사용중)</label>
          </div>
        </div>
        <div className="flex gap-3 mt-6">
          <button onClick={onClose} className="btn-secondary flex-1">취소</button>
          <button onClick={() => onSubmit({
            category: form.category,
            name: form.name,
            price: parseFloat(form.price) || 0,
            duration: parseInt(form.duration) || 0,
            is_active: form.is_active
          })} className="btn-primary flex-1">저장</button>
        </div>
      </div>
    </div>
  );
}
