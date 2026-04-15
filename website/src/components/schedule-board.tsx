"use client";

import type { Scenario, Shift, ComplianceReport } from "@/lib/types";
import { Badge } from "@/components/ui/badge";
import { useMemo } from "react";

const DAYS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

// Pure string operations for time display - no Date parsing, no timezone issues
function timeOf(iso: string): string {
  return iso.slice(11, 16); // "2025-03-10T08:00:00" -> "08:00"
}

function dateOf(iso: string): string {
  return iso.slice(0, 10); // "2025-03-10T08:00:00" -> "2025-03-10"
}

function hoursBetween(startIso: string, endIso: string): number {
  return (new Date(endIso).getTime() - new Date(startIso).getTime()) / 3600000;
}

function localDateStr(d: Date): string {
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
}

interface Props {
  scenario: Scenario;
  shifts: Shift[];
  report: ComplianceReport;
  onCellClick: (workerId: string, date: string) => void;
  onShiftClick: (uid: string) => void;
}

export function ScheduleBoard({ scenario, shifts, report, onCellClick, onShiftClick }: Props) {
  const violatedStaff = useMemo(
    () => new Set((report.violations || []).map(v => v.staff_id)),
    [report]
  );

  const violationCounts = useMemo(() => {
    const m = new Map<string, number>();
    for (const v of report.violations || []) m.set(v.staff_id, (m.get(v.staff_id) || 0) + 1);
    return m;
  }, [report]);

  // Calculate visible date range
  const { weekStart, numDays } = useMemo(() => {
    if (!shifts.length) return { weekStart: new Date(), numDays: 7 };
    const times = shifts.flatMap(s => [new Date(s.start).getTime(), new Date(s.end).getTime()]);
    const min = new Date(Math.min(...times));
    min.setHours(0, 0, 0, 0);
    const dow = min.getDay();
    min.setDate(min.getDate() - ((dow + 6) % 7)); // back to Monday
    let n = Math.max(7, Math.ceil((Math.max(...times) - min.getTime()) / 86400000));
    if (n % 7) n = Math.ceil(n / 7) * 7;
    return { weekStart: min, numDays: n };
  }, [shifts]);

  // Build day headers
  const days = useMemo(() => {
    return Array.from({ length: numDays }, (_, i) => {
      const d = new Date(weekStart);
      d.setDate(d.getDate() + i);
      return {
        label: DAYS[i % 7],
        date: d.getDate(),
        month: d.toLocaleString("en", { month: "short" }),
        dateStr: localDateStr(d),
        isWeekend: d.getDay() === 0 || d.getDay() === 6,
      };
    });
  }, [weekStart, numDays]);

  return (
    <div className="border border-neutral-200 rounded-xl overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full border-collapse" style={{ minWidth: Math.max(720, numDays * 100 + 176) }}>
          <thead>
            <tr>
              <th className="bg-neutral-50 text-left text-[11px] font-medium text-neutral-500 uppercase tracking-wide p-3 pl-4 border-b border-neutral-200 w-44 min-w-[176px] sticky left-0 z-10">
                Staff
              </th>
              {days.map((day, i) => (
                <th key={i} className={`text-center text-[11px] font-medium p-2 border-b border-l border-neutral-200 ${day.isWeekend ? "bg-neutral-100/50 text-neutral-400" : "bg-neutral-50 text-neutral-500"}`}>
                  <div className="font-semibold">{day.label}</div>
                  <div className="font-normal text-neutral-400 text-[10px]">{day.date} {day.month}</div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {scenario.workers.map(worker => {
              const isFlagged = violatedStaff.has(worker.id);
              const vCount = violationCounts.get(worker.id) || 0;
              const workerShifts = shifts
                .filter(s => s.staff_id === worker.id)
                .sort((a, b) => new Date(a.start).getTime() - new Date(b.start).getTime());

              return (
                <tr key={worker.id} className={isFlagged ? "bg-red-50/20" : ""}>
                  {/* Name cell */}
                  <td className={`p-3 pl-4 border-b border-neutral-100 sticky left-0 z-10 bg-white ${isFlagged ? "border-l-[3px] border-l-red-400 !bg-red-50/60" : ""}`}>
                    <div className="flex items-center gap-2">
                      <div className="w-2.5 h-2.5 rounded-full shrink-0" style={{ background: worker.color }} />
                      <div className="min-w-0 flex-1">
                        <div className={`text-[13px] font-semibold leading-tight truncate ${isFlagged ? "text-red-700" : ""}`}>
                          {worker.name}
                        </div>
                        <div className="text-[10px] text-neutral-400">{worker.role}</div>
                      </div>
                      {isFlagged && (
                        <Badge variant="destructive" className="text-[9px] px-1.5 py-0 h-4 shrink-0">{vCount}</Badge>
                      )}
                    </div>
                  </td>

                  {/* Day cells */}
                  {days.map((day, di) => {
                    const dayStart = new Date(weekStart);
                    dayStart.setDate(dayStart.getDate() + di);
                    dayStart.setHours(0, 0, 0, 0);
                    const dayEnd = new Date(dayStart);
                    dayEnd.setDate(dayEnd.getDate() + 1);

                    // Shifts visible on this day
                    const dayShifts = workerShifts.filter(sh => {
                      const s = new Date(sh.start), e = new Date(sh.end);
                      return s < dayEnd && e > dayStart;
                    });

                    return (
                      <td key={di} className={`border-b border-l border-neutral-100 p-1 align-top ${day.isWeekend ? "bg-neutral-50/30" : ""}`}>
                        <div
                          className="min-h-[52px] rounded cursor-pointer group transition-colors hover:bg-blue-50/40 p-0.5"
                          onClick={() => onCellClick(worker.id, day.dateStr)}
                        >
                          {dayShifts.length === 0 && (
                            <div className="h-[48px] flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                              <div className="w-5 h-5 rounded bg-neutral-200/80 flex items-center justify-center text-neutral-400 text-[11px]">+</div>
                            </div>
                          )}

                          {dayShifts.map((sh, si) => {
                            const dur = hoursBetween(sh.start, sh.end);
                            // Visible portion times
                            const visStart = new Date(Math.max(new Date(sh.start).getTime(), dayStart.getTime()));
                            const visEnd = new Date(Math.min(new Date(sh.end).getTime(), dayEnd.getTime()));
                            const visDur = (visEnd.getTime() - visStart.getTime()) / 3600000;

                            // Rest gap from previous shift
                            let restGap: number | null = null;
                            if (isFlagged && si === 0) {
                              const wi = workerShifts.indexOf(sh);
                              if (wi > 0) {
                                const prevEnd = new Date(workerShifts[wi - 1].end);
                                const thisStart = new Date(sh.start);
                                if (prevEnd >= dayStart && prevEnd <= thisStart) {
                                  const gap = (thisStart.getTime() - prevEnd.getTime()) / 3600000;
                                  if (gap > 0 && gap < 12) restGap = gap;
                                }
                              }
                            }

                            return (
                              <div key={sh._uid || si}>
                                {restGap !== null && (
                                  <div className="text-[9px] font-bold text-red-500 text-center py-0.5 mb-1 rounded bg-red-50 border border-dashed border-red-200">
                                    {Math.round(restGap)}h rest
                                  </div>
                                )}
                                <button
                                  type="button"
                                  onClick={e => { e.stopPropagation(); if (sh._uid) onShiftClick(sh._uid); }}
                                  className={`
                                    w-full text-left rounded-lg p-2 mb-1 last:mb-0 text-white text-[11px]
                                    cursor-pointer transition-all hover:brightness-110 hover:shadow-md active:scale-[0.98]
                                    ${isFlagged ? "ring-2 ring-red-400 ring-offset-1 ring-offset-white" : ""}
                                    ${sh.on_call ? "bg-stripes" : ""}
                                  `}
                                  style={{ backgroundColor: worker.color }}
                                >
                                  {sh.on_call && (
                                    <div className="text-[8px] font-bold uppercase tracking-wider opacity-60 mb-0.5">On-call</div>
                                  )}
                                  <div className="font-semibold text-[12px] tabular-nums tracking-tight">
                                    {timeOf(visStart.getFullYear() === dayStart.getFullYear() ? sh.start : `${localDateStr(visStart)}T${String(visStart.getHours()).padStart(2,"0")}:${String(visStart.getMinutes()).padStart(2,"0")}:00`)}
                                    <span className="opacity-40 mx-0.5">{"\u2013"}</span>
                                    {timeOf(visEnd.getFullYear() === dayEnd.getFullYear() ? sh.end : `${localDateStr(visEnd)}T${String(visEnd.getHours()).padStart(2,"0")}:${String(visEnd.getMinutes()).padStart(2,"0")}:00`)}
                                  </div>
                                  <div className="opacity-50 text-[10px]">
                                    {dur === visDur ? `${Math.round(dur)}h` : `${Math.round(visDur)}h of ${Math.round(dur)}h`}
                                  </div>
                                </button>
                              </div>
                            );
                          })}

                          {dayShifts.length > 0 && (
                            <div className="opacity-0 group-hover:opacity-100 transition-opacity flex justify-center pt-0.5">
                              <div className="w-4 h-4 rounded bg-neutral-100 flex items-center justify-center text-neutral-300 text-[9px]">+</div>
                            </div>
                          )}
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
