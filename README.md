# Salon Core (미용실 및 뷰티숍 통합 CRM 시스템)

Copyright 2026. Kimjibeom. All rights reserved.

핸드SOS를 완벽 대체하는 무료 티어 빌드 기반의 통합 CRM 시스템입니다. 예약, 고객, 매출, 직원 관리를 한 곳에서 효율적으로 처리할 수 있도록 설계되었습니다.

## 주요 기능

- **예약 및 대기 관리**: 일/주/월 캘린더 뷰, 당일 워크인 대기열 보드, 실시간 상태 업데이트
- **고객 및 차트 관리**: 고객별 시술 히스토리북, Before/After 사진 관리, 실시간 고객 검색
- **간편 POS**: 시술/점판 항목 결제, 복합 결제 지원, 멤버십 잔액 자동 차감
- **매출 분석 대시보드**: 일/월 마감 데이터, 6개월 트렌드 차트, 객단가 분석, 신규/재방문 비율, 이탈 고객 분석
- **직원 관리**: 직원별 목표 대비 실적 관리, 순매출 기반 인센티브 및 급여 정산
- **회원권 관리**: 금액 차감형 정액권 및 횟수 차감형 회원권 통합 관리

## 기술 스택

### Frontend
- Next.js 14 (App Router)
- React 18
- Tailwind CSS (Glassmorphism 디자인 시스템)
- Recharts (매출 분석 차트)
- TypeScript

### Backend
- Go 1.22+
- Gin-Gonic 프레임워크
- Gorilla WebSocket (실시간 예약/대기열 동기화)
- pgx (PostgreSQL 커넥션 풀링 - 무료 티어 제한 최적화)
- robfig/cron (백그라운드 알림 스케줄러)

### Database
- Supabase PostgreSQL
- Supabase Storage (미디어 업로드용)

## 아키텍처 제약 조건 및 최적화
1. **DB 커넥션 풀 제약**: Go 백엔드 DB 초기화 시 `MaxOpenConns`를 최대 20개로 제한하여 무료 티어 DB 과부하 방지
2. **미디어 스토리지 분리**: 사진/서류는 Supabase Storage에 업로드하고 DB에는 URL만 저장하여 용량 최적화
3. **콜드 스타트 대응**: Next.js API 요청 타임아웃을 60초로 설정 및 로딩 스피너 UI 구현

## 설치 및 실행 방법

### 요구사항
- Node.js 18+
- Go 1.22+
- Supabase 프로젝트 (PostgreSQL)

### 환경 변수 설정
최상단 디렉토리에 있는 `.env.example` 파일을 복사하여 `.env` 파일을 생성하고, Supabase 프로젝트 정보를 입력합니다.
```bash
cp .env.example .env
```

### 백엔드 실행 (Go)
```bash
cd backend
go mod tidy
go run cmd/main.go
```
서버는 기본적으로 `http://localhost:8080`에서 실행됩니다.

### 프론트엔드 실행 (Next.js)
```bash
cd frontend
npm install
npm run dev
```
클라이언트는 `http://localhost:3000`에서 실행됩니다.

## 보안 가이드

- **JWT 비밀 키**: `JWT_SECRET_KEY` 환경 변수를 반드시 설정하세요.
- **CORS**: 프론트엔드의 도메인을 환경 변수 `FRONTEND_URL`에 명시하여 보안을 유지합니다.
- 프로덕션 배포 시에는 SSL/HTTPS 환경에서 실행되어야 합니다.

## 라이선스
Copyright 2026. Kimjibeom. All rights reserved.
이 소프트웨어의 무단 복제, 배포, 수정을 금지합니다.