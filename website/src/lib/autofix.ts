import type { Shift, Violation, Worker } from "./types";
import { hoursBetween, formatDateTime } from "./dates";

function shiftHours(s: Shift): number {
  return hoursBetween(s.start, s.end);
}

export interface FixResult {
  shifts: Shift[];
  description: string;
}

export function applyFix(shifts: Shift[], violation: Violation, workers: Worker[]): FixResult {
  const sid = violation.staff_id;
  const k = violation.rule_key;
  let next = [...shifts];

  const staffShifts = (id: string) =>
    next.filter(s => s.staff_id === id).sort((a, b) => a.start.localeCompare(b.start));

  const patch = (uid: string, updates: Partial<Shift>) => {
    next = next.map(s => s._uid === uid ? { ...s, ...updates } : s);
  };

  // Weekly hours exceeded: shorten the longest shifts to reduce total hours
  if (k.includes("max-weekly") || k.includes("max-combined") || k.includes("max-ordinary")) {
    // violation.actual is the averaged weekly hours. We need to figure out total excess.
    const ss = staffShifts(sid);
    const totalHours = ss.reduce((sum, s) => sum + shiftHours(s), 0);
    // Estimate weeks from the date range
    const firstMs = new Date(ss[0].start).getTime();
    const lastMs = new Date(ss[ss.length - 1].end).getTime();
    const spanWeeks = Math.max(1, (lastMs - firstMs) / (7 * 86400000));
    const targetTotal = violation.limit * spanWeeks;
    let toRemove = Math.ceil(totalHours - targetTotal);
    if (toRemove <= 0) toRemove = Math.ceil(violation.actual - violation.limit); // fallback

    // Trim longest shifts first, keeping at least 6h each
    const sorted = [...ss].sort((a, b) => shiftHours(b) - shiftHours(a));
    let removed = 0;
    for (const s of sorted) {
      if (removed >= toRemove) break;
      const dur = shiftHours(s);
      const cut = Math.min(dur - 6, toRemove - removed);
      if (cut > 0) {
        const trimEnd = new Date(new Date(s.start).getTime() + (dur - cut) * 3600000);
        patch(s._uid!, { end: formatDateTime(trimEnd) });
        removed += cut;
      }
    }
    return { shifts: next, description: `Reduced ${Math.round(removed)} total hours across shifts to meet ${violation.limit}h/week average` };
  }

  // Days off: remove the shortest shift
  if (k.includes("days-off") || k.includes("day-of-rest")) {
    const ss = staffShifts(sid).sort((a, b) => shiftHours(a) - shiftHours(b));
    if (ss.length) {
      next = next.filter(s => s._uid !== ss[0]._uid);
      return { shifts: next, description: `Removed ${Math.round(shiftHours(ss[0]))}h shift to create a day off` };
    }
  }

  // Rest gap: push shifts forward
  if (k.includes("rest-between") || k.includes("min-rest")) {
    const ss = staffShifts(sid);
    for (let i = 1; i < ss.length; i++) {
      const prev = next.find(s => s._uid === ss[i - 1]._uid)!;
      const cur = next.find(s => s._uid === ss[i]._uid)!;
      const gap = (new Date(cur.start).getTime() - new Date(prev.end).getTime()) / 3600000;
      if (gap >= 0 && gap < violation.limit) {
        const dur = shiftHours(cur);
        const newStart = new Date(new Date(prev.end).getTime() + violation.limit * 3600000);
        const newEnd = new Date(newStart.getTime() + dur * 3600000);
        patch(cur._uid!, { start: formatDateTime(newStart), end: formatDateTime(newEnd) });
      }
    }
    return { shifts: next, description: `Adjusted shift times to ensure ${violation.limit}h rest between shifts` };
  }

  // Max shift hours: trim to limit
  if (k.includes("max-shift")) {
    for (const s of staffShifts(sid)) {
      if (shiftHours(s) > violation.limit) {
        const trimEnd = new Date(new Date(s.start).getTime() + violation.limit * 3600000);
        patch(s._uid!, { end: formatDateTime(trimEnd) });
      }
    }
    return { shifts: next, description: `Trimmed shifts to ${violation.limit}h maximum` };
  }

  // Guard duty exceeded: reassign excess to colleague with fewest guards
  if (k.includes("guards") || k.includes("on-call")) {
    const myGuards = staffShifts(sid).filter(s => s.on_call);
    const excess = myGuards.length - Math.floor(violation.limit);
    if (excess > 0) {
      const guardCounts = new Map<string, number>();
      for (const w of workers) {
        guardCounts.set(w.id, next.filter(s => s.staff_id === w.id && s.on_call).length);
      }
      const toReassign = myGuards.slice(-excess);
      let targetName = "";
      for (const guard of toReassign) {
        let minWorker = workers.find(w => w.id !== sid) || workers[0];
        let minCount = Infinity;
        for (const w of workers) {
          if (w.id === sid) continue;
          const count = guardCounts.get(w.id) || 0;
          if (count < minCount) { minCount = count; minWorker = w; }
        }
        patch(guard._uid!, { staff_id: minWorker.id, staff_type: minWorker.type });
        guardCounts.set(minWorker.id, (guardCounts.get(minWorker.id) || 0) + 1);
        guardCounts.set(sid, (guardCounts.get(sid) || 0) - 1);
        targetName = minWorker.name;
      }
      return { shifts: next, description: `Reassigned ${excess} guard${excess > 1 ? "s" : ""} to ${targetName}` };
    }
  }

  // Fallback: shorten the last shift
  const ss = staffShifts(sid);
  if (ss.length) {
    const last = ss[ss.length - 1];
    const trimEnd = new Date(new Date(last.start).getTime() + Math.min(shiftHours(last), 8) * 3600000);
    patch(last._uid!, { end: formatDateTime(trimEnd) });
  }
  return { shifts: next, description: "Shortened last shift" };
}
