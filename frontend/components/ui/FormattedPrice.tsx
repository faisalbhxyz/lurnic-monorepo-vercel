"use client";

import { useCurrency } from "@/context/CurrencyContext";

export default function FormattedPrice({
  amount,
  className,
}: {
  amount: number;
  className?: string;
}) {
  const { formatAmount } = useCurrency();
  return <span className={className}>{formatAmount(amount)}</span>;
}
