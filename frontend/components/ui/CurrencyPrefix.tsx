"use client";

import { useCurrency } from "@/context/CurrencyContext";

export default function CurrencyPrefix({
  className = "w-10 h-8 flex items-center justify-center border-r border-gray-300",
}: {
  className?: string;
}) {
  const { symbol } = useCurrency();
  return <span className={className}>{symbol}</span>;
}
