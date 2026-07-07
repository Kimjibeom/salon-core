// Copyright 2026. Kimjibeom. All rights reserved.
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: '온라인 예약 - Salon Core',
  description: '간편하게 미용실 예약을 진행하세요. 시술 선택, 디자이너 선택, 날짜와 시간을 골라 바로 예약할 수 있습니다.',
  keywords: ['미용실 예약', '헤어살롱', '온라인 예약', '뷰티'],
};

export default function BookingLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ko">
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-purple-50">
        {children}
      </body>
    </html>
  );
}
