"use client";

import { useMemo } from "react";
import type { Scenario, Shift, ComplianceReport } from "@/lib/types";
import { dateOf, hoursBetween, mondayOf, addDays, dayName, dayNum, dayOfWeek, shiftOverlapsDate } from "@/lib/dates";
import { Badge } from "@/components/ui/badge";

interface Props {
  scenario: Scenario;
  shifts: Shift[];
  report: ComplianceReport;
  onCellClick: (workerId: string, date: string) => void;
  onShiftClick: (uid: string) => void;
}

export function ScheduleBoard({ scenario, shifts, report, onCellClick, onShiftClick }: Props) {
  const violatedStaff = useMemo(() => new Set((report.violations || []).map(v => v.staff_id)), [report]);
  const vCounts = useMemo(() => {
    const m = new Map<string, number>();
    for (const v of report.violations || []) m.set(v.staff_id, (m.get(v.staff_id) || 0) + 1);
    return m;
  }, [report]);

  // Grid dates
  const gridDates = useMemo(() => {
    if (!shifts.length) return [];
    const allDates = [...shifts.map(s => dateOf(s.start)), ...shifts.map(s => dateOf(s.end))].sort();
    const monday = mondayOf(allDates[0]);
    const last = allDates[allDates.length - 1];
    const dates: string[] = [];
    let cur = monday;
    while (cur <= last || dates.length % 7 !== 0) { dates.push(cur); cur = addDays(cur, 1); }
    if (dates.length < 7) while (dates.length < 7) { dates.push(cur); cur = addDays(cur, 1); }
    return dates;
  }, [shifts]);

  // Aggregate: for each worker+date, total hours and whether any shift is on-call
  const cellData = useMemo(() => {
    const map = new Map<string, { hours: number; onCall: boolean; uids: string[] }>();
    for (const sh of shifts) {
      const dur = hoursBetween(sh.start, sh.end);
      // A shift might span multiple days
      for (const d of gridDates) {
        if (shiftOverlapsDate(sh.start, sh.end, d)) {
          // Calculate hours on this specific day
          const dayStartMs = new Date(+d.slice(0,4), +d.slice(5,7)-1, +d.slice(8,10), 0).getTime();
          const dayEndMs = dayStartMs + 86400000;
          const sMs = new Date(+sh.start.slice(0,4), +sh.start.slice(5,7)-1, +sh.start.slice(8,10), +sh.start.slice(11,13)||0, +sh.start.slice(14,16)||0).getTime();
          const eMs = new Date(+sh.end.slice(0,4), +sh.end.slice(5,7)-1, +sh.end.slice(8,10), +sh.end.slice(11,13)||0, +sh.end.slice(14,16)||0).getTime();
          const overlapH = (Math.min(eMs, dayEndMs) - Math.max(sMs, dayStartMs)) / 3600000;
          if (overlapH <= 0) continue;
          const key = `${sh.staff_id}|${d}`;
          const existing = map.get(key) || { hours: 0, onCall: false, uids: [] };
          existing.hours += overlapH;
          if (sh.on_call) existing.onCall = true;
          existing.uids.push(sh._uid || "");
          map.set(key, existing);
        }
      }
    }
    return map;
  }, [shifts, gridDates]);

  // Week separators
  const weekStarts = useMemo(() => {
    const s = new Set<number>();
    gridDates.forEach((d, i) => { if (i > 0 && dayOfWeek(d) === 1) s.add(i); });
    return s;
  }, [gridDates]);

  return (
    <div className="border border-neutral-200 rounded-xl overflow-hidden bg-white">
      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr>
              <th className="bg-neutral-50 text-left text-[10px] font-semibold text-neutral-400 uppercase tracking-wider p-2.5 pl-4 border-b border-neutral-200 sticky left-0 z-10 min-w-[160px]">
                Staff
              </th>
              {gridDates.map((d, i) => {
                const isWe = dayOfWeek(d) === 0 || dayOfWeek(d) === 6;
                const isWeekBorder = weekStarts.has(i);
                return (
                  <th key={d} className={`text-center text-[10px] font-medium p-1.5 border-b border-neutral-200 min-w-[44px] ${isWe ? "bg-neutral-100/60 text-neutral-300" : "bg-neutral-50 text-neutral-400"} ${isWeekBorder ? "border-l-2 border-l-neutral-300" : "border-l border-l-neutral-100"}`}>
                    <div className="font-semibold text-[9px]">{dayName(d).slice(0, 2)}</div>
                    <div>{dayNum(d)}</div>
                  </th>
                );
              })}
              <th className="bg-neutral-50 text-center text-[10px] font-semibold text-neutral-400 uppercase tracking-wider p-2.5 border-b border-l-2 border-neutral-200 min-w-[50px]">
                Total
              </th>
            </tr>
          </thead>
          <tbody>
            {scenario.workers.map(worker => {
              const flagged = violatedStaff.has(worker.id);
              const vc = vCounts.get(worker.id) || 0;
              const totalHours = shifts.filter(s => s.staff_id === worker.id).reduce((sum, s) => sum + hoursBetween(s.start, s.end), 0);

              return (
                <tr key={worker.id} className={flagged ? "bg-red-50/30" : "hover:bg-neutral-50/50"}>
                  {/* Name */}
                  <td className={`p-0 border-b border-neutral-100 sticky left-0 z-10 ${flagged ? "bg-red-50" : "bg-white"}`}>
                    <div className="flex items-center gap-0 h-full min-h-[44px]">
                      <div className="w-1 self-stretch shrink-0" style={{ backgroundColor: flagged ? "#ef4444" : worker.color }} />
                      <div className="flex items-center gap-2 px-3 py-2 flex-1 min-w-0">
                        <div className="min-w-0">
                          <div className={`text-[12px] font-semibold leading-tight truncate ${flagged ? "text-red-700" : ""}`}>{worker.name}</div>
                          <div className="text-[9px] text-neutral-400">{worker.role}</div>
                        </div>
                        {flagged && <Badge variant="destructive" className="text-[8px] px-1 py-0 h-3.5 shrink-0 ml-auto">{vc}</Badge>}
                      </div>
                    </div>
                  </td>

                  {/* Day cells */}
                  {gridDates.map((d, i) => {
                    const key = `${worker.id}|${d}`;
                    const cell = cellData.get(key);
                    const isWe = dayOfWeek(d) === 0 || dayOfWeek(d) === 6;
                    const isWeekBorder = weekStarts.has(i);
                    const hasShift = cell && cell.hours > 0;

                    return (
                      <td key={d} className={`border-b border-neutral-100 text-center align-middle p-0 ${isWeekBorder ? "border-l-2 border-l-neutral-200" : "border-l border-l-neutral-100"}`}>
                        <button
                          type="button"
                          onClick={() => {
                            if (hasShift && cell.uids[0]) onShiftClick(cell.uids[0]);
                            else onCellClick(worker.id, d);
                          }}
                          className={`
                            w-full h-full min-h-[44px] text-[11px] font-semibold tabular-nums transition-all cursor-pointer
                            ${!hasShift && isWe ? "bg-neutral-50" : ""}
                            ${!hasShift && !isWe ? "bg-white hover:bg-neutral-50" : ""}
                            ${hasShift ? "text-white hover:brightness-110" : "text-neutral-300 hover:text-neutral-400"}
                            ${hasShift && flagged ? "ring-2 ring-inset ring-red-500" : ""}
                          `}
                          style={hasShift ? {
                            backgroundColor: cell.onCall ? "#f59e0b" : worker.color,
                          } : undefined}
                        >
                          {hasShift ? (
                            <>
                              <div className="text-[12px]">{Math.round(cell.hours)}h</div>
                              {cell.onCall && <div className="text-[7px] font-bold uppercase opacity-60 leading-none mt-0.5">guard</div>}
                            </>
                          ) : (
                            <div className="text-[10px] opacity-0 hover:opacity-100 transition-opacity">+</div>
                          )}
                        </button>
                      </td>
                    );
                  })}

                  {/* Total */}
                  <td className="border-b border-l-2 border-neutral-200 p-2 text-center">
                    <div className={`text-[12px] font-bold tabular-nums ${flagged ? "text-red-600" : "text-neutral-700"}`}>
                      {Math.round(totalHours)}h
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
