"use client";

import type { Scenario, Shift, ComplianceReport, Violation } from "@/lib/types";
import { Badge } from "@/components/ui/badge";
import { useMemo } from "react";

function fmt(iso: string) {
  const d = new Date(iso);
  return `${String(d.getHours()).padStart(2, "0")}:${String(d.getMinutes()).padStart(2, "0")}`;
}

function dur(s: Shift) {
  return (new Date(s.end).getTime() - new Date(s.start).getTime()) / 3600000;
}

const DAYS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

interface Props {
  scenario: Scenario;
  shifts: Shift[];
  report: ComplianceReport;
  onCellClick: (workerId: string, date: string) => void;
  onShiftClick: (index: number) => void;
}

export function ScheduleBoard({ scenario, shifts, report, onCellClick, onShiftClick }: Props) {
  const violatedStaff = useMemo(
    () => new Set((report.violations || []).map(v => v.staff_id)),
    [report]
  );

  // Build a map of violations per staff+shift for precise highlighting
  const shiftViolations = useMemo(() => {
    const map = new Map<string, Violation[]>();
    for (const v of report.violations || []) {
      const existing = map.get(v.staff_id) || [];
      existing.push(v);
      map.set(v.staff_id, existing);
    }
    return map;
  }, [report]);

  // Calculate date range
  const { weekStart, numDays } = useMemo(() => {
    const times = shifts.flatMap(s => [new Date(s.start).getTime(), new Date(s.end).getTime()]);
    if (!times.length) return { weekStart: new Date(), numDays: 7 };
    const min = new Date(Math.min(...times));
    min.setHours(0, 0, 0, 0);
    const dow = min.getDay();
    min.setDate(min.getDate() - ((dow + 6) % 7));
    let n = Math.max(7, Math.ceil((Math.max(...times) - min.getTime()) / 86400000));
    if (n % 7) n = Math.ceil(n / 7) * 7;
    return { weekStart: min, numDays: n };
  }, [shifts]);

  // Count violations per worker
  const violationCounts = useMemo(() => {
    const counts = new Map<string, number>();
    for (const v of report.violations || []) {
      counts.set(v.staff_id, (counts.get(v.staff_id) || 0) + 1);
    }
    return counts;
  }, [report]);

  return (
    <div className="border border-neutral-200 rounded-xl overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full border-collapse min-w-[720px]">
          <thead>
            <tr>
              <th className="bg-neutral-50/80 text-left text-[11px] font-medium text-neutral-500 uppercase tracking-wide p-3 pl-4 border-b border-neutral-200 w-44 min-w-[176px]">
                Staff member
              </th>
              {Array.from({ length: numDays }, (_, d) => {
                const dt = new Date(weekStart);
                dt.setDate(dt.getDate() + d);
                const isWeekend = dt.getDay() === 0 || dt.getDay() === 6;
                return (
                  <th key={d} className={`text-center text-[11px] font-medium p-2.5 border-b border-l border-neutral-200 ${isWeekend ? "bg-neutral-100/60 text-neutral-400" : "bg-neutral-50/80 text-neutral-500"}`}>
                    <div className="font-semibold">{DAYS[d % 7]}</div>
                    <div className="font-normal text-neutral-400 text-[10px] mt-0.5">
                      {dt.getDate()} {dt.toLocaleString("en", { month: "short" })}
                    </div>
                  </th>
                );
              })}
            </tr>
          </thead>
          <tbody>
            {scenario.workers.map(w => {
              const isFlagged = violatedStaff.has(w.id);
              const vCount = violationCounts.get(w.id) || 0;
              const workerShifts = shifts
                .filter(s => s.staff_id === w.id)
                .sort((a, b) => new Date(a.start).getTime() - new Date(b.start).getTime());

              return (
                <tr key={w.id} className={`transition-colors ${isFlagged ? "bg-red-50/30" : "hover:bg-neutral-50/30"}`}>
                  {/* Worker name cell */}
                  <td className={`p-3 pl-4 border-b border-neutral-100 ${isFlagged ? "border-l-[3px] border-l-red-400 bg-red-50/50" : ""}`}>
                    <div className="flex items-center gap-2.5">
                      <div className="w-2.5 h-2.5 rounded-full shrink-0 ring-2 ring-white" style={{ background: w.color }} />
                      <div className="min-w-0">
                        <div className={`text-[13px] font-semibold leading-tight truncate ${isFlagged ? "text-red-700" : "text-neutral-900"}`}>
                          {w.name}
                        </div>
                        <div className="text-[10px] text-neutral-400 mt-0.5">{w.role}</div>
                      </div>
                      {isFlagged && (
                        <Badge variant="destructive" className="text-[9px] px-1.5 py-0 h-4 shrink-0">
                          {vCount}
                        </Badge>
                      )}
                    </div>
                  </td>

                  {/* Day cells */}
                  {Array.from({ length: numDays }, (_, d) => {
                    const dS = new Date(weekStart);
                    dS.setDate(dS.getDate() + d);
                    dS.setHours(0, 0, 0, 0);
                    const dE = new Date(dS);
                    dE.setDate(dE.getDate() + 1);
                    const dateStr = dS.toISOString().slice(0, 10);
                    const isWeekend = dS.getDay() === 0 || dS.getDay() === 6;
                    const dayShifts = workerShifts.filter(sh => new Date(sh.start) < dE && new Date(sh.end) > dS);

                    return (
                      <td key={d} className={`border-b border-l border-neutral-100 p-1 align-top ${isWeekend ? "bg-neutral-50/40" : ""}`}>
                        <div
                          className="min-h-[56px] rounded-md cursor-pointer group relative transition-colors hover:bg-neutral-100/50 p-0.5"
                          onClick={() => onCellClick(w.id, dateStr)}
                        >
                          {/* Add shift hint on hover */}
                          {dayShifts.length === 0 && (
                            <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                              <div className="w-6 h-6 rounded-full bg-neutral-200 flex items-center justify-center text-neutral-400 text-xs font-bold">+</div>
                            </div>
                          )}

                          {dayShifts.map((sh, si) => {
                            const idx = shifts.findIndex(s => s.staff_id === sh.staff_id && s.start === sh.start && s.end === sh.end);
                            const st = new Date(sh.start);
                            const en = new Date(sh.end);
                            const vS = new Date(Math.max(st.getTime(), dS.getTime()));
                            const vE = new Date(Math.min(en.getTime(), dE.getTime()));
                            const totalDur = dur(sh);
                            const visDur = (vE.getTime() - vS.getTime()) / 3600000;

                            // Rest gap check
                            let restGap: number | null = null;
                            if (si === 0 && isFlagged) {
                              const wi = workerShifts.indexOf(sh);
                              if (wi > 0) {
                                const pE = new Date(workerShifts[wi - 1].end);
                                if (pE >= dS && pE <= st) {
                                  const g = (st.getTime() - pE.getTime()) / 3600000;
                                  if (g > 0 && g < 12) restGap = g;
                                }
                              }
                            }

                            return (
                              <div key={si}>
                                {restGap !== null && (
                                  <div className="text-[9px] font-bold text-red-500 text-center py-0.5 mb-1 rounded bg-red-50 border border-dashed border-red-200/80">
                                    {restGap.toFixed(0)}h rest
                                  </div>
                                )}
                                <button
                                  type="button"
                                  onClick={(e) => { e.stopPropagation(); onShiftClick(idx); }}
                                  className={`
                                    w-full text-left rounded-lg p-2 mb-1 last:mb-0 text-white text-[11px] leading-tight
                                    cursor-pointer transition-all duration-150 hover:brightness-110 hover:shadow-md
                                    ${isFlagged ? "ring-2 ring-red-400/80 ring-offset-1 ring-offset-white shadow-sm shadow-red-100" : "hover:shadow-sm"}
                                    ${sh.on_call ? "bg-stripes" : ""}
                                  `}
                                  style={{ backgroundColor: w.color }}
                                >
                                  {sh.on_call && (
                                    <div className="text-[8px] font-bold uppercase tracking-wider opacity-60 mb-0.5">On-call guard</div>
                                  )}
                                  <div className="font-semibold text-[12px]">
                                    {fmt(vS.toISOString())} - {fmt(vE.toISOString())}
                                  </div>
                                  <div className="opacity-60 text-[10px] mt-0.5">
                                    {totalDur === visDur ? `${totalDur.toFixed(0)}h` : `${visDur.toFixed(0)}h / ${totalDur.toFixed(0)}h total`}
                                  </div>
                                </button>
                              </div>
                            );
                          })}
                        </div>
                      </td>
                    );
                  })}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
