// Copyright 2026. Kimjibeom. All rights reserved.
'use client';

import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react';
import api from '@/lib/api';
import type { Staff } from '@/types';

interface AuthContextType {
  staff: Staff | null;
  token: string | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  isAdmin: boolean;
  isDesigner: boolean;
}

const AUTH_TOKEN_KEY = 'salon_core_token';
const AUTH_STAFF_KEY = 'salon_core_staff';

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [staff, setStaff] = useState<Staff | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Restore session from localStorage on mount
    try {
      const savedToken = localStorage.getItem(AUTH_TOKEN_KEY);
      const savedStaff = localStorage.getItem(AUTH_STAFF_KEY);
      if (savedToken && savedStaff) {
        setToken(savedToken);
        setStaff(JSON.parse(savedStaff));
        api.setToken(savedToken);
      }
    } catch {
      // If localStorage is corrupted, clear it
      localStorage.removeItem(AUTH_TOKEN_KEY);
      localStorage.removeItem(AUTH_STAFF_KEY);
    }
    setIsLoading(false);
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const response = await api.post<{ token: string; staff: Staff }>('/api/auth/login', { email, password });
    setToken(response.token);
    setStaff(response.staff);
    api.setToken(response.token);
    // Persist to localStorage
    localStorage.setItem(AUTH_TOKEN_KEY, response.token);
    localStorage.setItem(AUTH_STAFF_KEY, JSON.stringify(response.staff));
  }, []);

  const logout = useCallback(() => {
    setToken(null);
    setStaff(null);
    api.setToken(null);
    // Clear persisted session
    localStorage.removeItem(AUTH_TOKEN_KEY);
    localStorage.removeItem(AUTH_STAFF_KEY);
    // Full page reload to clear all in-memory state
    window.location.href = '/';
  }, []);

  return (
    <AuthContext.Provider value={{
      staff,
      token,
      isLoading,
      login,
      logout,
      isAdmin: staff?.role === 'admin',
      isDesigner: staff?.role === 'admin' || staff?.role === 'designer',
    }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
}
