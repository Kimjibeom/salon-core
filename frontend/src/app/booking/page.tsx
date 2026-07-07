// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { useState, useEffect, useCallback } from 'react';

const API_BASE = process.env.NEXT_PUBLIC_BACKEND_URL || '';

interface ShopInfo {
  shop_name: string;
  shop_phone: string;
  shop_address: string;
  shop_open_time: string;
  shop_close_time: string;
}

interface Service {
  id: string;
  category: string;
  name: string;
  price: number;
  duration: number;
}

interface Designer {
  id: string;
  name: string;
  day_off: number[];
}

interface TimeSlot {
  start_time: string;
  end_time: string;
  available: boolean;
}

interface AvailabilityResponse {
  date: string;
  staff_id: string;
  staff_name: string;
  slots: TimeSlot[];
}

type Step = 1 | 2 | 3 | 4;

export default function BookingPage() {
  const [step, setStep] = useState<Step>(1);
  const [shopInfo, setShopInfo] = useState<ShopInfo | null>(null);
  const [services, setServices] = useState<Service[]>([]);
  const [designers, setDesigners] = useState<Designer[]>([]);
  const [availability, setAvailability] = useState<AvailabilityResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isComplete, setIsComplete] = useState(false);
  const [error, setError] = useState('');

  // Selections
  const [selectedService, setSelectedService] = useState<Service | null>(null);
  const [selectedDesigner, setSelectedDesigner] = useState<Designer | null>(null);
  const [selectedDate, setSelectedDate] = useState('');
  const [selectedSlot, setSelectedSlot] = useState<TimeSlot | null>(null);
  const [customerName, setCustomerName] = useState('');
  const [customerPhone, setCustomerPhone] = useState('');
  const [memo, setMemo] = useState('');

  // Fetch shop info + services + designers on mount
  useEffect(() => {
    const fetchInitialData = async () => {
      try {
        const [shopRes, servicesRes, designersRes] = await Promise.all([
          fetch(`${API_BASE}/api/booking/shop`),
          fetch(`${API_BASE}/api/booking/services`),
          fetch(`${API_BASE}/api/booking/staff`),
        ]);
        if (shopRes.ok) setShopInfo(await shopRes.json());
        if (servicesRes.ok) setServices(await servicesRes.json() || []);
        if (designersRes.ok) setDesigners(await designersRes.json() || []);
      } catch {
        console.error('Failed to load booking data');
      }
    };
    fetchInitialData();
  }, []);

  // Fetch availability when date + designer + service are selected
  const fetchAvailability = useCallback(async () => {
    if (!selectedDate || !selectedDesigner || !selectedService) return;
    setIsLoading(true);
    setAvailability(null);
    try {
      const res = await fetch(
        `${API_BASE}/api/booking/availability?date=${selectedDate}&staff_id=${selectedDesigner.id}&service_id=${selectedService.id}`
      );
      if (res.ok) {
        const data: AvailabilityResponse = await res.json();
        setAvailability(data);
      }
    } catch {
      setError('시간 정보를 불러오는데 실패했습니다.');
    }
    setIsLoading(false);
  }, [selectedDate, selectedDesigner, selectedService]);

  useEffect(() => {
    fetchAvailability();
  }, [fetchAvailability]);

  // Generate next 14 days
  const getDateOptions = () => {
    const dates: { value: string; label: string; dayOfWeek: number }[] = [];
    const today = new Date();
    for (let i = 1; i <= 14; i++) {
      const d = new Date(today);
      d.setDate(today.getDate() + i);
      const value = d.toISOString().split('T')[0];
      const dayNames = ['일', '월', '화', '수', '목', '금', '토'];
      const label = `${d.getMonth() + 1}/${d.getDate()} (${dayNames[d.getDay()]})`;
      dates.push({ value, label, dayOfWeek: d.getDay() });
    }
    return dates;
  };

  const handleSubmit = async () => {
    if (!selectedService || !selectedDesigner || !selectedSlot || !customerName || !customerPhone) {
      setError('모든 필수 항목을 입력해주세요.');
      return;
    }

    setIsSubmitting(true);
    setError('');
    try {
      const res = await fetch(`${API_BASE}/api/booking/reservations`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          customer_name: customerName,
          customer_phone: customerPhone,
          service_id: selectedService.id,
          staff_id: selectedDesigner.id,
          date: selectedDate,
          start_time: selectedSlot.start_time,
          memo,
        }),
      });

      if (res.ok) {
        setIsComplete(true);
      } else {
        const data = await res.json();
        setError(data.error || '예약에 실패했습니다. 다시 시도해주세요.');
      }
    } catch {
      setError('네트워크 오류가 발생했습니다.');
    }
    setIsSubmitting(false);
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('ko-KR').format(price) + '원';
  };

  // Group services by category
  const servicesByCategory = services.reduce<Record<string, Service[]>>((acc, svc) => {
    if (!acc[svc.category]) acc[svc.category] = [];
    acc[svc.category].push(svc);
    return acc;
  }, {});

  const categoryLabels: Record<string, string> = {
    cut: '✂️ 커트',
    perm: '🌀 펌',
    color: '🎨 염색',
    clinic: '💊 클리닉',
    product: '🧴 제품',
  };

  // Completion screen
  if (isComplete) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4">
        <div className="max-w-md w-full text-center space-y-6">
          <div className="w-20 h-20 mx-auto bg-green-100 rounded-full flex items-center justify-center">
            <svg className="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-gray-900">예약이 완료되었습니다!</h1>
          <div className="bg-white rounded-2xl shadow-lg p-6 text-left space-y-3">
            <div className="flex justify-between">
              <span className="text-gray-500">시술</span>
              <span className="font-medium">{selectedService?.name}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">담당</span>
              <span className="font-medium">{selectedDesigner?.name} 디자이너</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">날짜</span>
              <span className="font-medium">{selectedDate}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">시간</span>
              <span className="font-medium">{selectedSlot?.start_time} ~ {selectedSlot?.end_time}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">예약자</span>
              <span className="font-medium">{customerName}</span>
            </div>
          </div>
          <p className="text-sm text-gray-500">
            예약 확인 및 변경은 매장에 직접 문의해주세요.
            {shopInfo?.shop_phone && <><br/>📞 {shopInfo.shop_phone}</>}
          </p>
          <button
            onClick={() => window.location.reload()}
            className="w-full py-3 bg-purple-600 text-white rounded-xl font-semibold hover:bg-purple-700 transition-colors"
          >
            새 예약하기
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen" style={{ fontFamily: "'Inter', sans-serif" }}>
      {/* Header */}
      <header className="sticky top-0 z-50 bg-white/80 backdrop-blur-xl border-b border-gray-100">
        <div className="max-w-lg mx-auto px-4 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-lg font-bold text-gray-900">{shopInfo?.shop_name || '미용실'}</h1>
            <p className="text-xs text-gray-400">온라인 예약</p>
          </div>
          {shopInfo?.shop_phone && (
            <a href={`tel:${shopInfo.shop_phone}`} className="text-purple-600 text-sm font-medium">
              📞 전화
            </a>
          )}
        </div>
      </header>

      {/* Progress Bar */}
      <div className="max-w-lg mx-auto px-4 pt-6">
        <div className="flex items-center gap-1 mb-2">
          {[1, 2, 3, 4].map((s) => (
            <div
              key={s}
              className={`flex-1 h-1.5 rounded-full transition-all duration-300 ${
                s <= step ? 'bg-purple-600' : 'bg-gray-200'
              }`}
            />
          ))}
        </div>
        <div className="flex justify-between text-xs text-gray-400 mb-6">
          <span className={step >= 1 ? 'text-purple-600 font-medium' : ''}>시술 선택</span>
          <span className={step >= 2 ? 'text-purple-600 font-medium' : ''}>디자이너</span>
          <span className={step >= 3 ? 'text-purple-600 font-medium' : ''}>날짜/시간</span>
          <span className={step >= 4 ? 'text-purple-600 font-medium' : ''}>정보 입력</span>
        </div>
      </div>

      {/* Content */}
      <main className="max-w-lg mx-auto px-4 pb-32">
        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-xl text-sm text-red-600">
            {error}
          </div>
        )}

        {/* Step 1: Service Selection */}
        {step === 1 && (
          <div className="space-y-6 animate-fadeIn">
            <h2 className="text-xl font-bold text-gray-900">어떤 시술을 원하시나요?</h2>
            {Object.entries(servicesByCategory).map(([category, svcs]) => (
              <div key={category}>
                <h3 className="text-sm font-semibold text-gray-500 uppercase mb-3">
                  {categoryLabels[category] || category}
                </h3>
                <div className="space-y-2">
                  {svcs.map((svc) => (
                    <button
                      key={svc.id}
                      onClick={() => {
                        setSelectedService(svc);
                        setSelectedSlot(null);
                        setStep(2);
                        setError('');
                      }}
                      className={`w-full text-left p-4 rounded-xl border-2 transition-all duration-200 hover:shadow-md ${
                        selectedService?.id === svc.id
                          ? 'border-purple-500 bg-purple-50 shadow-md'
                          : 'border-gray-100 bg-white hover:border-purple-200'
                      }`}
                    >
                      <div className="flex justify-between items-start">
                        <div>
                          <p className="font-semibold text-gray-900">{svc.name}</p>
                          <p className="text-sm text-gray-400 mt-0.5">약 {svc.duration}분 소요</p>
                        </div>
                        <span className="text-purple-600 font-bold">{formatPrice(svc.price)}</span>
                      </div>
                    </button>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Step 2: Designer Selection */}
        {step === 2 && (
          <div className="space-y-6 animate-fadeIn">
            <h2 className="text-xl font-bold text-gray-900">담당 디자이너를 선택해주세요</h2>
            <div className="grid grid-cols-2 gap-3">
              {designers.map((designer) => (
                <button
                  key={designer.id}
                  onClick={() => {
                    setSelectedDesigner(designer);
                    setSelectedSlot(null);
                    setSelectedDate('');
                    setStep(3);
                    setError('');
                  }}
                  className={`p-5 rounded-xl border-2 transition-all duration-200 hover:shadow-md text-center ${
                    selectedDesigner?.id === designer.id
                      ? 'border-purple-500 bg-purple-50 shadow-md'
                      : 'border-gray-100 bg-white hover:border-purple-200'
                  }`}
                >
                  <div className="w-14 h-14 mx-auto bg-gradient-to-br from-purple-400 to-indigo-500 rounded-full flex items-center justify-center text-white text-xl font-bold mb-3">
                    {designer.name.charAt(0)}
                  </div>
                  <p className="font-semibold text-gray-900">{designer.name}</p>
                  <p className="text-xs text-gray-400 mt-1">디자이너</p>
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Step 3: Date & Time */}
        {step === 3 && (
          <div className="space-y-6 animate-fadeIn">
            <h2 className="text-xl font-bold text-gray-900">날짜와 시간을 선택해주세요</h2>

            {/* Date Selection */}
            <div>
              <h3 className="text-sm font-semibold text-gray-500 mb-3">📅 날짜 선택</h3>
              <div className="flex gap-2 overflow-x-auto pb-2 -mx-1 px-1">
                {getDateOptions().map((d) => {
                  const isDayOff = selectedDesigner?.day_off?.includes(d.dayOfWeek);
                  return (
                    <button
                      key={d.value}
                      disabled={isDayOff}
                      onClick={() => {
                        setSelectedDate(d.value);
                        setSelectedSlot(null);
                        setError('');
                      }}
                      className={`flex-shrink-0 px-4 py-3 rounded-xl text-sm font-medium transition-all duration-200 ${
                        isDayOff
                          ? 'bg-gray-100 text-gray-300 cursor-not-allowed'
                          : selectedDate === d.value
                          ? 'bg-purple-600 text-white shadow-lg shadow-purple-200'
                          : 'bg-white border border-gray-200 text-gray-700 hover:border-purple-300'
                      }`}
                    >
                      {d.label}
                      {isDayOff && <span className="block text-xs">휴무</span>}
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Time Slots */}
            {selectedDate && (
              <div>
                <h3 className="text-sm font-semibold text-gray-500 mb-3">⏰ 시간 선택</h3>
                {isLoading ? (
                  <div className="flex items-center justify-center py-12">
                    <div className="w-8 h-8 border-3 border-purple-200 border-t-purple-600 rounded-full animate-spin" />
                  </div>
                ) : availability && availability.slots.length > 0 ? (
                  <div className="grid grid-cols-3 gap-2">
                    {availability.slots.map((slot) => (
                      <button
                        key={slot.start_time}
                        disabled={!slot.available}
                        onClick={() => {
                          setSelectedSlot(slot);
                          setStep(4);
                          setError('');
                        }}
                        className={`py-3 px-2 rounded-xl text-sm font-medium transition-all duration-200 ${
                          !slot.available
                            ? 'bg-gray-100 text-gray-300 cursor-not-allowed line-through'
                            : selectedSlot?.start_time === slot.start_time
                            ? 'bg-purple-600 text-white shadow-lg shadow-purple-200'
                            : 'bg-white border border-gray-200 text-gray-700 hover:border-purple-300 hover:bg-purple-50'
                        }`}
                      >
                        {slot.start_time}
                      </button>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-400">
                    <p className="text-lg mb-1">😢</p>
                    <p className="text-sm">해당 날짜에 예약 가능한 시간이 없습니다.</p>
                  </div>
                )}
              </div>
            )}
          </div>
        )}

        {/* Step 4: Customer Info */}
        {step === 4 && (
          <div className="space-y-6 animate-fadeIn">
            <h2 className="text-xl font-bold text-gray-900">예약 정보를 입력해주세요</h2>

            {/* Summary */}
            <div className="bg-purple-50 border border-purple-100 rounded-xl p-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-purple-600">시술</span>
                <span className="font-medium text-gray-900">{selectedService?.name} ({selectedService?.duration}분)</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-purple-600">담당</span>
                <span className="font-medium text-gray-900">{selectedDesigner?.name} 디자이너</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-purple-600">날짜</span>
                <span className="font-medium text-gray-900">{selectedDate}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-purple-600">시간</span>
                <span className="font-medium text-gray-900">{selectedSlot?.start_time} ~ {selectedSlot?.end_time}</span>
              </div>
              <div className="flex justify-between text-sm pt-2 border-t border-purple-200">
                <span className="text-purple-600 font-semibold">예상 금액</span>
                <span className="font-bold text-purple-700">{selectedService ? formatPrice(selectedService.price) : ''}</span>
              </div>
            </div>

            {/* Form */}
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">이름 *</label>
                <input
                  type="text"
                  value={customerName}
                  onChange={(e) => setCustomerName(e.target.value)}
                  placeholder="홍길동"
                  className="w-full px-4 py-3 rounded-xl border border-gray-200 bg-white focus:border-purple-500 focus:ring-2 focus:ring-purple-200 outline-none transition-all"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">연락처 *</label>
                <input
                  type="tel"
                  value={customerPhone}
                  onChange={(e) => setCustomerPhone(e.target.value)}
                  placeholder="010-1234-5678"
                  className="w-full px-4 py-3 rounded-xl border border-gray-200 bg-white focus:border-purple-500 focus:ring-2 focus:ring-purple-200 outline-none transition-all"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">요청 사항 (선택)</label>
                <textarea
                  value={memo}
                  onChange={(e) => setMemo(e.target.value)}
                  placeholder="예: 앞머리 자르기 원합니다"
                  rows={3}
                  className="w-full px-4 py-3 rounded-xl border border-gray-200 bg-white focus:border-purple-500 focus:ring-2 focus:ring-purple-200 outline-none transition-all resize-none"
                />
              </div>
            </div>

            <button
              onClick={handleSubmit}
              disabled={isSubmitting || !customerName || !customerPhone}
              className="w-full py-4 bg-gradient-to-r from-purple-600 to-indigo-600 text-white rounded-xl font-bold text-lg shadow-lg shadow-purple-200 hover:shadow-xl hover:from-purple-700 hover:to-indigo-700 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting ? (
                <span className="flex items-center justify-center gap-2">
                  <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  예약 중...
                </span>
              ) : (
                '예약 확정하기'
              )}
            </button>
          </div>
        )}
      </main>

      {/* Bottom Navigation (Back Button) */}
      {step > 1 && (
        <div className="fixed bottom-0 left-0 right-0 bg-white/90 backdrop-blur-xl border-t border-gray-100 p-4">
          <div className="max-w-lg mx-auto">
            <button
              onClick={() => {
                setStep((s) => Math.max(1, s - 1) as Step);
                setError('');
              }}
              className="w-full py-3 border-2 border-gray-200 rounded-xl font-semibold text-gray-600 hover:border-purple-300 hover:text-purple-600 transition-all"
            >
              ← 이전 단계
            </button>
          </div>
        </div>
      )}

      <style jsx>{`
        @keyframes fadeIn {
          from { opacity: 0; transform: translateY(12px); }
          to { opacity: 1; transform: translateY(0); }
        }
        .animate-fadeIn {
          animation: fadeIn 0.35s ease-out;
        }
      `}</style>
    </div>
  );
}
