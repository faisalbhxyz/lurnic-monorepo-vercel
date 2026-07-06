export const formatDate = (date: string): string => {
  const parsedDate = new Date(date);
  const day = String(parsedDate.getDate()).padStart(2, "0");
  const month = parsedDate.toLocaleString("en-US", { month: "long" }); // Full month name
  const year = parsedDate.getFullYear();

  return `${day} ${month}, ${year}`;
};

export const formatDateTime = (date: string): string => {
  const parsedDate = new Date(date);
  const day = String(parsedDate.getDate()).padStart(2, "0");
  const month = parsedDate.toLocaleString("en-US", { month: "short" });
  const year = parsedDate.getFullYear();
  const hours = parsedDate.getHours();
  const minutes = String(parsedDate.getMinutes()).padStart(2, "0");
  const period = hours < 12 ? "AM" : "PM";
  const hour12 = hours % 12 === 0 ? 12 : hours % 12;

  return `${day} ${month} ${year}, ${hour12}:${minutes} ${period}`;
};

/** DB TIME / ISO datetime → backend UI format "10:05 AM". */
export const dbTimeToPickerFormat = (timeStr: string): string => {
  const hhmm = scheduleTimeToHHMM(timeStr);
  if (!hhmm) return "";
  return hhmmToScheduleTime(hhmm);
};

export function getFirstFormError(
  errors: Record<string, unknown>
): { path: string; message: string } | null {
  const walk = (
    obj: Record<string, unknown>,
    path = ""
  ): { path: string; message: string } | null => {
    for (const [key, value] of Object.entries(obj)) {
      const nextPath = path ? `${path}.${key}` : key;
      if (value && typeof value === "object") {
        if (
          "message" in value &&
          typeof (value as { message?: unknown }).message === "string"
        ) {
          const message = (value as { message: string }).message;
          if (message.trim()) return { path: nextPath, message };
        }
        const nested = walk(value as Record<string, unknown>, nextPath);
        if (nested) return nested;
      }
    }
    return null;
  };
  return walk(errors);
}

/** 24h "HH:MM" for native `<input type="time" />` from DB "HH:MM:SS" or UI "hh:mm AM/PM". */
export function scheduleTimeToHHMM(
  value: string | null | undefined
): string {
  if (value == null || String(value).trim() === "") return "";
  const t = String(value).trim();
  // Go / JSON may serialize TIME as full ISO datetime
  if (/^\d{4}-\d{2}-\d{2}T/.test(t)) {
    const d = new Date(t);
    if (!Number.isNaN(d.getTime())) {
      return `${String(d.getHours()).padStart(2, "0")}:${String(d.getMinutes()).padStart(2, "0")}`;
    }
  }
  const twelve = t.match(/^(\d{1,2}):(\d{2})\s*(AM|PM)$/i);
  if (twelve) {
    let h = parseInt(twelve[1], 10);
    const m = parseInt(twelve[2], 10);
    const ap = twelve[3].toUpperCase();
    if (ap === "PM" && h !== 12) h += 12;
    if (ap === "AM" && h === 12) h = 0;
    return `${String(h).padStart(2, "0")}:${String(m).padStart(2, "0")}`;
  }
  const twentyFour = t.match(/^(\d{1,2}):(\d{2})(?::(\d{2}))?$/);
  if (twentyFour) {
    const h = parseInt(twentyFour[1], 10);
    const m = parseInt(twentyFour[2], 10);
    return `${String(h).padStart(2, "0")}:${String(m).padStart(2, "0")}`;
  }
  return "";
}

/** Convert native time input "HH:MM" to backend format "03:04 PM" (Go: time.Parse("03:04 PM", …)). */
export function hhmmToScheduleTime(hhmm: string): string {
  const [hs, ms] = hhmm.split(":");
  const h24 = parseInt(hs, 10);
  const m = parseInt(ms ?? "0", 10);
  if (Number.isNaN(h24) || Number.isNaN(m)) return "";
  const period = h24 < 12 ? "AM" : "PM";
  const h12 = h24 % 12 === 0 ? 12 : h24 % 12;
  return `${String(h12).padStart(2, "0")}:${String(m).padStart(2, "0")} ${period}`;
}
