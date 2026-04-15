// Pure date utilities using ONLY string operations and local-time Date construction.
// Every function that takes a date string expects "YYYY-MM-DD" or "YYYY-MM-DDThh:mm:ss".
// No toISOString(). No UTC. No timezone surprises.

/** Extract "YYYY-MM-DD" from "YYYY-MM-DDThh:mm:ss" */
export function dateOf(iso: string): string {
  return iso.slice(0, 10);
}

/** Extract "hh:mm" from "YYYY-MM-DDThh:mm:ss" */
export function timeOf(iso: string): string {
  return iso.slice(11, 16);
}

/** Hours between two ISO datetime strings */
export function hoursBetween(a: string, b: string): number {
  return (toMs(b) - toMs(a)) / 3600000;
}

/** Convert "YYYY-MM-DDThh:mm:ss" to epoch ms (LOCAL time) */
function toMs(iso: string): number {
  // Parse manually to guarantee local time interpretation
  const y = +iso.slice(0, 4), m = +iso.slice(5, 7) - 1, d = +iso.slice(8, 10);
  const h = +iso.slice(11, 13) || 0, min = +iso.slice(14, 16) || 0, s = +iso.slice(17, 19) || 0;
  return new Date(y, m, d, h, min, s).getTime();
}

/** Format a local Date to "YYYY-MM-DD" */
export function formatDate(d: Date): string {
  const p = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}`;
}

/** Format a local Date to "YYYY-MM-DDThh:mm:ss" */
export function formatDateTime(d: Date): string {
  const p = (n: number) => String(n).padStart(2, "0");
  return `${formatDate(d)}T${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`;
}

/** Add days to "YYYY-MM-DD", returns "YYYY-MM-DD" */
export function addDays(dateStr: string, n: number): string {
  const d = new Date(+dateStr.slice(0, 4), +dateStr.slice(5, 7) - 1, +dateStr.slice(8, 10), 12); // noon to dodge DST
  d.setDate(d.getDate() + n);
  return formatDate(d);
}

/** Day of week for "YYYY-MM-DD" (0=Sun, 1=Mon, ..., 6=Sat) */
export function dayOfWeek(dateStr: string): number {
  return new Date(+dateStr.slice(0, 4), +dateStr.slice(5, 7) - 1, +dateStr.slice(8, 10), 12).getDay();
}

/** Get Monday of the week containing dateStr */
export function mondayOf(dateStr: string): string {
  const dow = dayOfWeek(dateStr);
  const offset = dow === 0 ? -6 : 1 - dow; // Sunday -> -6, Mon -> 0, Tue -> -1, etc.
  return addDays(dateStr, offset);
}

/** Day name */
const DAYS = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
export function dayName(dateStr: string): string {
  return DAYS[dayOfWeek(dateStr)];
}

/** Month short name */
export function monthShort(dateStr: string): string {
  const m = +dateStr.slice(5, 7);
  return ["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"][m - 1];
}

/** Day number */
export function dayNum(dateStr: string): number {
  return +dateStr.slice(8, 10);
}

/** Does a shift overlap with a given date? Both args are "YYYY-MM-DD" style. */
export function shiftOverlapsDate(shiftStart: string, shiftEnd: string, dateStr: string): boolean {
  const dayStartMs = toMs(dateStr + "T00:00:00");
  const dayEndMs = toMs(addDays(dateStr, 1) + "T00:00:00");
  const sMs = toMs(shiftStart);
  const eMs = toMs(shiftEnd);
  return sMs < dayEndMs && eMs > dayStartMs;
}

/** Clamp a shift's visible start/end to a single day. Returns times as "hh:mm". */
export function clampToDay(shiftStart: string, shiftEnd: string, dateStr: string): { visStart: string; visEnd: string; visDur: number } {
  const dayStartMs = toMs(dateStr + "T00:00:00");
  const dayEndMs = toMs(addDays(dateStr, 1) + "T00:00:00");
  const sMs = Math.max(toMs(shiftStart), dayStartMs);
  const eMs = Math.min(toMs(shiftEnd), dayEndMs);
  const visStartDate = new Date(sMs);
  const visEndDate = new Date(eMs);
  const p = (n: number) => String(n).padStart(2, "0");
  return {
    visStart: `${p(visStartDate.getHours())}:${p(visStartDate.getMinutes())}`,
    visEnd: eMs === dayEndMs ? "24:00" : `${p(visEndDate.getHours())}:${p(visEndDate.getMinutes())}`,
    visDur: (eMs - sMs) / 3600000,
  };
}
