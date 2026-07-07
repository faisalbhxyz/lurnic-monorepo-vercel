"use client";

import React, { createContext, useContext, useMemo } from "react";
import {
  CurrencyCode,
  DEFAULT_CURRENCY,
  formatCurrency,
  getCurrencySymbol,
  resolveCurrency,
} from "@/lib/currency";

type CurrencyContextValue = {
  currency: CurrencyCode;
  symbol: string;
  formatAmount: (amount: number) => string;
};

const CurrencyContext = createContext<CurrencyContextValue>({
  currency: DEFAULT_CURRENCY,
  symbol: getCurrencySymbol(DEFAULT_CURRENCY),
  formatAmount: (amount) => formatCurrency(amount, DEFAULT_CURRENCY),
});

export function CurrencyProvider({
  currency,
  children,
}: {
  currency?: string | null;
  children: React.ReactNode;
}) {
  const value = useMemo(() => {
    const resolved = resolveCurrency(currency);
    return {
      currency: resolved.code,
      symbol: resolved.symbol,
      formatAmount: (amount: number) => formatCurrency(amount, resolved.code),
    };
  }, [currency]);

  return (
    <CurrencyContext.Provider value={value}>{children}</CurrencyContext.Provider>
  );
}

export function useCurrency() {
  return useContext(CurrencyContext);
}
