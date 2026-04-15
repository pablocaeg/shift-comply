"use client";

import { useMemo } from "react";
import type { Scenario, Shift, ComplianceReport } from "@/lib/types";
import { dateOf, hoursBetween, mondayOf, addDays, dayName, dayNum, monthShort, dayOfWeek, shiftOverlapsDate } from "@/lib/dates";

const COLS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

interface Props {
  scenario: Scenario;
  shifts: Shift[];
  report: ComplianceReport;
  onCellClick: (workerId: string, date: string) => void;
  onShiftClick: (uid: string) => void;
}

export function ScheduleBoard({ scenario, shifts, report, onCellClick, onShiftClick }: Props) {
  const violatedStaff = useMemo(() => new Set((report.violations || []).map(v => v.staff_id)), [report]);

  // Build weeks (each week is an array of 7 date strings)
  const weeks = useMemo(() => {
    if (!shifts.length) return [];
    const allDates = [...shifts.map(s => dateOf(s.start)), ...shifts.map(s => dateOf(s.end))].sort();
    const monday = mondayOf(allDates[0]);
    const last = allDates[allDates.length - 1];
    const dates: string[] = [];
    let cur = monday;
    while (cur <= last || dates.length % 7 !== 0) { dates.push(cur); cur = addDays(cur, 1); }
    const w: string[][] = [];
    for (let i = 0; i < dates.length; i += 7) w.push(dates.slice(i, i + 7));
    return w;
  }, [shifts]);

  // Cell data: worker+date -> { hours, onCall, uids }
  const cellData = useMemo(() => {
    const map = new Map<string, { hours: number; onCall: boolean; uids: string[] }>();
    const allDates = weeks.flat();
    for (const sh of shifts) {
      for (const d of allDates) {
        if (!shiftOverlapsDate(sh.start, sh.end, d)) continue;
        const dayStartMs = new Date(+d.slice(0,4), +d.slice(5,7)-1, +d.slice(8,10), 0).getTime();
        const dayEndMs = dayStartMs + 86400000;
        const sMs = new Date(+sh.start.slice(0,4), +sh.start.slice(5,7)-1, +sh.start.slice(8,10), +(sh.start.slice(11,13)||0), +(sh.start.slice(14,16)||0)).getTime();
        const eMs = new Date(+sh.end.slice(0,4), +sh.end.slice(5,7)-1, +sh.end.slice(8,10), +(sh.end.slice(11,13)||0), +(sh.end.slice(14,16)||0)).getTime();
        const h = (Math.min(eMs, dayEndMs) - Math.max(sMs, dayStartMs)) / 3600000;
        if (h <= 0) continue;
        const key = `${sh.staff_id}|${d}`;
        const ex = map.get(key) || { hours: 0, onCall: false, uids: [] };
        ex.hours += h;
        if (sh.on_call) ex.onCall = true;
        ex.uids.push(sh._uid || "");
        map.set(key, ex);
      }
    }
    return map;
  }, [shifts, weeks]);

  return (
    <div className="space-y-1">
      {/* Day headers */}
      <div className="grid grid-cols-7 gap-1">
        {COLS.map(d => (
          <div key={d} className="text-center text-[10px] font-semibold text-neutral-400 uppercase tracking-wider py-1">{d}</div>
        ))}
      </div>

      {/* Weeks */}
      {weeks.map((week, wi) => (
        <div key={wi} className="grid grid-cols-7 gap-1">
          {week.map(d => {
            const isWe = dayOfWeek(d) === 0 || dayOfWeek(d) === 6;

            return (
              <div key={d} className={`rounded-lg border min-h-[100px] flex flex-col ${isWe ? "border-neutral-100 bg-neutral-50/50" : "border-neutral-200 bg-white"}`}>
                {/* Date label */}
                <div className="px-2 pt-1.5 pb-1 flex items-baseline justify-between">
                  <span className={`text-[12px] font-semibold ${isWe ? "text-neutral-300" : "text-neutral-500"}`}>{dayNum(d)}</span>
                  {dayNum(d) === 1 && <span className="text-[9px] font-medium text-neutral-400 uppercase">{monthShort(d)}</span>}
                </div>

                {/* Worker pills */}
                <div className="px-1 pb-1 flex flex-col gap-0.5 flex-1">
                  {scenario.workers.map(w => {
                    const key = `${w.id}|${d}`;
                    const cell = cellData.get(key);
                    if (!cell || cell.hours <= 0) return null;
                    const flagged = violatedStaff.has(w.id);

                    return (
                      <button
                        key={w.id}
                        type="button"
                        onClick={() => { if (cell.uids[0]) onShiftClick(cell.uids[0]); }}
                        className={`
                          w-full rounded px-1.5 py-0.5 text-white text-[10px] font-semibold
                          flex items-center justify-between gap-1
                          cursor-pointer transition-all hover:brightness-110
                          ${flagged ? "ring-1 ring-red-500 ring-offset-1" : ""}
                        `}
                        style={{ backgroundColor: cell.onCall ? "#f59e0b" : w.color }}
                      >
                        <span className="truncate">{w.name.split(" ").pop()}</span>
                        <span className="tabular-nums opacity-70 shrink-0">
                          {Math.round(cell.hours)}h
                          {cell.onCall && <span className="text-[7px] ml-0.5">G</span>}
                        </span>
                      </button>
                    );
                  })}

                  {/* Add button if no shifts for any worker */}
                  {!scenario.workers.some(w => cellData.has(`${w.id}|${d}`)) && (
                    <button
                      type="button"
                      onClick={() => onCellClick(scenario.workers[0].id, d)}
                      className="flex-1 flex items-center justify-center text-neutral-300 text-[11px] opacity-0 hover:opacity-100 transition-opacity"
                    >
                      +
                    </button>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      ))}

      {/* Legend */}
      <div className="flex items-center gap-4 pt-3 px-1">
        {scenario.workers.map(w => (
          <div key={w.id} className="flex items-center gap-1.5">
            <div className="w-2.5 h-2.5 rounded-sm" style={{ backgroundColor: w.color }} />
            <span className={`text-[11px] ${violatedStaff.has(w.id) ? "text-red-600 font-semibold" : "text-neutral-500"}`}>
              {w.name}
              {violatedStaff.has(w.id) && <span className="text-red-400 ml-1 text-[9px]">violations</span>}
            </span>
          </div>
        ))}
        <div className="flex items-center gap-1.5">
          <div className="w-2.5 h-2.5 rounded-sm bg-amber-500" />
          <span className="text-[11px] text-neutral-500">On-call guard</span>
        </div>
      </div>
    </div>
  );
}
