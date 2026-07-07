export type CurrencyCode = "BDT" | "USD" | "EUR" | "GBP" | "INR";

export type CurrencyOption = {
  code: CurrencyCode;
  label: string;
  symbol: string;
  locale: string;
};

export const CURRENCIES: CurrencyOption[] = [
  { code: "BDT", label: "Bangladeshi Taka (BDT)", symbol: "৳", locale: "en-BD" },
  { code: "USD", label: "US Dollar (USD)", symbol: "$", locale: "en-US" },
  { code: "EUR", label: "Euro (EUR)", symbol: "€", locale: "de-DE" },
  { code: "GBP", label: "British Pound (GBP)", symbol: "£", locale: "en-GB" },
  { code: "INR", label: "Indian Rupee (INR)", symbol: "₹", locale: "en-IN" },
];

export const DEFAULT_CURRENCY: CurrencyCode = "BDT";

const currencyMap = new Map(CURRENCIES.map((item) => [item.code, item]));

export function resolveCurrency(code?: string | null): CurrencyOption {
  const match = currencyMap.get((code ?? DEFAULT_CURRENCY) as CurrencyCode);
  return match ?? currencyMap.get(DEFAULT_CURRENCY)!;
}

export function getCurrencySymbol(code?: string | null) {
  return resolveCurrency(code).symbol;
}

export function formatCurrency(amount: number, code?: string | null) {
  const currency = resolveCurrency(code);
  return `${currency.symbol} ${amount.toLocaleString(currency.locale)}`;
}
