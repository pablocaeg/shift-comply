"use client";

import { useMemo } from "react";
import type { Scenario, Shift, ComplianceReport } from "@/lib/types";
import { dateOf, hoursBetween, mondayOf, addDays, dayName, dayNum, monthShort, dayOfWeek, shiftOverlapsDate } from "@/lib/dates";

const COLS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

interface Props {
  scenario: Scenario;
  shifts: Shift[];
  report: ComplianceReport;
  fixedStaff: Set<string>;
  onCellClick: (workerId: string, date: string) => void;
  onShiftClick: (uid: string) => void;
}

export function ScheduleBoard({ scenario, shifts, report, fixedStaff, onCellClick, onShiftClick }: Props) {
  const violatedStaff = useMemo(() => new Set((report.violations || []).map(v => v.staff_id)), [report]);

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
        ex.hours += h; if (sh.on_call) ex.onCall = true; ex.uids.push(sh._uid || "");
        map.set(key, ex);
      }
    }
    return map;
  }, [shifts, weeks]);

  return (
    <div>
      {/* Day headers */}
      <div className="grid grid-cols-7 gap-1 mb-1">
        {COLS.map(c => <div key={c} className="text-center text-[10px] font-semibold text-neutral-400 uppercase tracking-wider py-1">{c}</div>)}
      </div>

      {/* Weeks */}
      <div className="space-y-1">
        {weeks.map((week, wi) => (
          <div key={wi} className="grid grid-cols-7 gap-1">
            {week.map(d => {
              const isWe = dayOfWeek(d) === 0 || dayOfWeek(d) === 6;
              const hasViolation = scenario.workers.some(w => violatedStaff.has(w.id) && cellData.has(`${w.id}|${d}`));
              const hasFixed = scenario.workers.some(w => fixedStaff.has(w.id) && !violatedStaff.has(w.id) && cellData.has(`${w.id}|${d}`));

              return (
                <div key={d} className={`rounded-lg border min-h-[90px] flex flex-col transition-colors ${
                  hasViolation ? "border-red-300 bg-red-50/30" :
                  hasFixed ? "border-emerald-200 bg-emerald-50/20" :
                  isWe ? "border-neutral-100 bg-neutral-50/40" :
                  "border-neutral-200 bg-white"
                }`}>
                  <div className="px-2 pt-1.5 flex items-baseline justify-between">
                    <span className={`text-[11px] font-semibold tabular-nums ${isWe ? "text-neutral-300" : "text-neutral-500"}`}>{dayNum(d)}</span>
                    {dayNum(d) <= 7 && wi === 0 && <span className="text-[8px] font-medium text-neutral-400 uppercase">{monthShort(d)}</span>}
                  </div>

                  <div className="px-1 pb-1 flex flex-col gap-[3px] flex-1 mt-0.5">
                    {scenario.workers.map(w => {
                      const cell = cellData.get(`${w.id}|${d}`);
                      if (!cell || cell.hours <= 0) return null;

                      const flagged = violatedStaff.has(w.id);
                      const fixed = fixedStaff.has(w.id) && !flagged;
                      const lastName = w.name.split(" ").pop() || w.name;
                      const sh = shifts.find(s => s._uid === cell.uids[0]);
                      const timeStr = sh ? `${sh.start.slice(11,16)}\u2013${sh.end.slice(11,16)}` : "";

                      let bg = w.color;
                      let extraClass = "";
                      if (flagged) { bg = "#ef4444"; extraClass = "ring-2 ring-red-600 ring-offset-1 shadow-sm shadow-red-200 animate-pulse"; }
                      else if (fixed) { bg = "#10b981"; extraClass = "ring-2 ring-emerald-400 ring-offset-1"; }
                      else if (cell.onCall) { bg = "#f59e0b"; extraClass = "bg-stripes border border-amber-600/30"; }

                      return (
                        <button key={w.id} type="button"
                          onClick={() => { if (cell.uids[0]) onShiftClick(cell.uids[0]); }}
                          className={`w-full rounded-[5px] px-1.5 py-1 text-left text-white cursor-pointer select-none transition-all hover:brightness-110 active:scale-[0.97] ${extraClass}`}
                          style={{ backgroundColor: bg }}
                        >
                          <div className="flex items-center gap-1 text-[10px] font-semibold">
                            {flagged && <span className="text-[9px] shrink-0">&#9888;</span>}
                            {fixed && <span className="text-[9px] shrink-0">&#10003;</span>}
                            <span className="truncate">{lastName}</span>
                            <span className="tabular-nums opacity-70 ml-auto shrink-0">{Math.round(cell.hours)}h</span>
                            {cell.onCall && !flagged && !fixed && <span className="text-[7px] font-bold opacity-60 shrink-0">G</span>}
                          </div>
                          <div className="text-[8px] opacity-50 tabular-nums">{timeStr}</div>
                        </button>
                      );
                    })}

                    <button type="button" onClick={() => onCellClick(scenario.workers[0].id, d)}
                      className="min-h-[20px] flex items-center justify-center text-neutral-300 text-[10px] opacity-0 hover:opacity-100 transition-opacity rounded hover:bg-neutral-100 mt-auto">+</button>
                  </div>
                </div>
              );
            })}
          </div>
        ))}
      </div>

      {/* Legend */}
      <div className="flex items-center gap-4 flex-wrap pt-4 px-1">
        {scenario.workers.map(w => {
          const f = violatedStaff.has(w.id);
          const fx = fixedStaff.has(w.id) && !f;
          return (
            <div key={w.id} className="flex items-center gap-1.5">
              <div className={`w-3 h-3 rounded ${f ? "bg-red-500" : fx ? "bg-emerald-500" : ""}`} style={f || fx ? undefined : { backgroundColor: w.color }} />
              <span className={`text-[11px] ${f ? "text-red-600 font-semibold" : fx ? "text-emerald-600 font-semibold" : "text-neutral-500"}`}>
                {w.name}{f ? " - violations" : ""}{fx ? " - fixed" : ""}
              </span>
            </div>
          );
        })}
        <div className="flex items-center gap-1.5">
          <div className="w-3 h-3 rounded bg-amber-500 bg-stripes" />
          <span className="text-[11px] text-neutral-500">On-call guard</span>
        </div>
      </div>
    </div>
  );
}
