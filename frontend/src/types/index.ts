// Copyright 2026. Kimjibeom. All rights reserved.

export interface Staff {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'designer' | 'staff';
  phone?: string;
  service_incentive_rate: number;
  product_incentive_rate: number;
  base_salary: number;
  monthly_target: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ServiceMenu {
  id: string;
  category: string;
  name: string;
  price: number;
  duration: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Customer {
  id: string;
  name: string;
  phone: string;
  email?: string;
  birth_date?: string;
  memo: string;
  tags: string[];
  created_at: string;
  last_visited_at?: string;
  visit_count: number;
}

export interface CustomerWithHistory extends Customer {
  charts: Chart[];
  memberships: Membership[];
  recent_sales?: Sale[];
}

export interface Chart {
  id: string;
  customer_id: string;
  staff_id: string;
  staff_name?: string;
  service_id?: string;
  recipe: string;
  treatment_name?: string;
  treatment_duration?: number;
  notes: string;
  before_img_url?: string;
  after_img_url?: string;
  consent_doc_url?: string;
  created_at: string;
}

export interface Membership {
  id: string;
  customer_id: string;
  type: 'money' | 'count';
  name: string;
  initial_balance: number;
  balance: number;
  initial_count: number;
  remaining_count: number;
  target_treatment?: string;
  expired_at?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface MembershipTransaction {
  id: string;
  membership_id: string;
  sale_id?: string;
  amount: number;
  count_change: number;
  balance_before: number;
  balance_after: number;
  count_before: number;
  count_after: number;
  description: string;
  created_at: string;
}

export interface Reservation {
  id: string;
  customer_id?: string;
  staff_id?: string;
  customer_name: string;
  customer_phone: string;
  staff_name?: string;
  service_id?: string;
  treatment_name?: string;
  date: string;
  start_time: string;
  end_time: string;
  status: 'reserved' | 'waiting' | 'in_progress' | 'completed' | 'canceled' | 'no_show';
  source: 'online' | 'offline' | 'naver';
  waiting_number?: number;
  waiting_started_at?: string;
  memo: string;
  created_at: string;
  updated_at: string;
}

export interface WaitingQueueEntry extends Reservation {
  wait_time_minutes: number;
  position: number;
}

export interface Sale {
  id: string;
  reservation_id?: string;
  customer_id?: string;
  staff_id: string;
  staff_name?: string;
  customer_name?: string;
  service_id?: string;
  item_name: string;
  total_amount: number;
  category: 'service' | 'product';
  payment_method: 'card' | 'cash' | 'membership' | 'mixed';
  card_amount: number;
  cash_amount: number;
  membership_amount: number;
  card_commission_rate: number;
  membership_id?: string;
  memo: string;
  created_at: string;
}

export interface DailySummary {
  date: string;
  total_revenue: number;
  service_revenue: number;
  product_revenue: number;
  card_revenue: number;
  cash_revenue: number;
  membership_revenue: number;
  transaction_count: number;
  customer_count: number;
}

export interface AnalyticsSummary {
  period: string;
  total_revenue: number;
  service_revenue: number;
  product_revenue: number;
  avg_service_price: number;
  avg_product_price: number;
  card_revenue: number;
  cash_revenue: number;
  membership_revenue: number;
  new_customers: number;
  returning_customers: number;
  total_customers: number;
}

export interface ChurnAnalysis {
  visit_number: number;
  total_customers: number;
  churned_customers: number;
  churn_rate: number;
  last_staff_name?: string;
}

export interface PaymentBreakdown {
  method: string;
  amount: number;
  ratio: number;
}

export interface StaffPerformance {
  staff_id: string;
  staff_name: string;
  monthly_target: number;
  service_revenue: number;
  product_revenue: number;
  total_revenue: number;
  achievement_rate: number;
  service_incentive: number;
  product_incentive: number;
  total_incentive: number;
  net_service_revenue: number;
  net_product_revenue: number;
  base_salary: number;
  total_payroll: number;
}

export interface WebSocketEvent {
  type: 'NEW_RESERVATION' | 'RESERVATION_UPDATE' | 'WAITING_QUEUE_UPDATE' | 'CID_INCOMING' | 'ONLINE_BOOKING' | 'NOTIFICATION';
  payload: unknown;
  time: string;
}
