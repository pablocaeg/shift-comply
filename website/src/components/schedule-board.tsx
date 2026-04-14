"use client";

import type { Scenario, Shift, ComplianceReport, Violation } from "@/lib/types";
import { Badge } from "@/components/ui/badge";
import { useMemo } from "react";
import { DndContext, DragOverlay, useDraggable, useDroppable, type DragEndEvent, type DragStartEvent } from "@dnd-kit/core";
import { useState } from "react";

function fmtTime(d: Date): string {
  return `${String(d.getHours()).padStart(2, "0")}:${String(d.getMinutes()).padStart(2, "0")}`;
}

function shiftDur(s: Shift): number {
  return (new Date(s.end).getTime() - new Date(s.start).getTime()) / 3600000;
}

function localDateStr(d: Date): string {
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
}

const DAYS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

interface Props {
  scenario: Scenario;
  shifts: Shift[];
  report: ComplianceReport;
  onCellClick: (workerId: string, date: string) => void;
  onShiftClick: (uid: string) => void;
  onMoveShift: (uid: string, newWorkerId: string) => void;
}

function DraggableShift({ shift, worker, isFlagged, onShiftClick }: {
  shift: Shift; worker: { color: string }; isFlagged: boolean;
  onShiftClick: (uid: string) => void;
}) {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({
    id: shift._uid || "",
    data: { shift },
  });

  const st = new Date(shift.start);
  const en = new Date(shift.end);
  const dur = shiftDur(shift);

  return (
    <div
      ref={setNodeRef}
      {...listeners}
      {...attributes}
      onClick={(e) => { e.stopPropagation(); onShiftClick(shift._uid || ""); }}
      className={`
        w-full text-left rounded-lg p-2 mb-1 last:mb-0 text-white text-[11px] leading-tight
        cursor-grab active:cursor-grabbing transition-all duration-150
        ${isDragging ? "opacity-30 scale-95" : "hover:brightness-110 hover:shadow-md"}
        ${isFlagged ? "ring-2 ring-red-400/80 ring-offset-1 ring-offset-white shadow-sm shadow-red-100" : "hover:shadow-sm"}
        ${shift.on_call ? "bg-stripes" : ""}
      `}
      style={{ backgroundColor: worker.color }}
    >
      {shift.on_call && (
        <div className="text-[8px] font-bold uppercase tracking-wider opacity-60 mb-0.5">On-call guard</div>
      )}
      <div className="font-semibold text-[12px] tabular-nums">
        {fmtTime(st)}<span className="opacity-50 mx-0.5">to</span>{fmtTime(en)}
      </div>
      <div className="opacity-60 text-[10px] mt-0.5">{dur.toFixed(0)}h</div>
    </div>
  );
}

function DroppableCell({ workerId, dateStr, children }: {
  workerId: string; dateStr: string; children: React.ReactNode;
}) {
  const { setNodeRef, isOver } = useDroppable({
    id: `${workerId}|${dateStr}`,
    data: { workerId, dateStr },
  });

  return (
    <div
      ref={setNodeRef}
      className={`min-h-[56px] rounded-md transition-colors p-0.5 ${isOver ? "bg-blue-50 ring-2 ring-blue-300 ring-inset" : "group"}`}
    >
      {children}
    </div>
  );
}

export function ScheduleBoard({ scenario, shifts, report, onCellClick, onShiftClick, onMoveShift }: Props) {
  const [activeShift, setActiveShift] = useState<Shift | null>(null);

  const violatedStaff = useMemo(
    () => new Set((report.violations || []).map(v => v.staff_id)),
    [report]
  );

  const violationCounts = useMemo(() => {
    const counts = new Map<string, number>();
    for (const v of report.violations || []) {
      counts.set(v.staff_id, (counts.get(v.staff_id) || 0) + 1);
    }
    return counts;
  }, [report]);

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

  function handleDragStart(event: DragStartEvent) {
    const shift = event.active.data.current?.shift as Shift | undefined;
    setActiveShift(shift || null);
  }

  function handleDragEnd(event: DragEndEvent) {
    setActiveShift(null);
    if (!event.over) return;
    const uid = event.active.id as string;
    const targetWorkerId = (event.over.data.current as { workerId?: string })?.workerId;
    if (targetWorkerId && uid) {
      onMoveShift(uid, targetWorkerId);
    }
  }

  const activeWorker = activeShift ? scenario.workers.find(w => w.id === activeShift.staff_id) : null;

  return (
    <DndContext onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
      <div className="border border-neutral-200 rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full border-collapse min-w-[720px]">
            <thead>
              <tr>
                <th className="bg-neutral-50/80 text-left text-[11px] font-medium text-neutral-500 uppercase tracking-wide p-3 pl-4 border-b border-neutral-200 w-44 min-w-[176px]">
                  Staff member
                </th>
                {Array.from({ length: numDays }, (_, d) => {
                  const dt = new Date(weekStart); dt.setDate(dt.getDate() + d);
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

                    {Array.from({ length: numDays }, (_, d) => {
                      const dS = new Date(weekStart); dS.setDate(dS.getDate() + d); dS.setHours(0, 0, 0, 0);
                      const dE = new Date(dS); dE.setDate(dE.getDate() + 1);
                      const dateStr = localDateStr(dS);
                      const isWeekend = dS.getDay() === 0 || dS.getDay() === 6;
                      const dayShifts = workerShifts.filter(sh => new Date(sh.start) < dE && new Date(sh.end) > dS);

                      return (
                        <td key={d} className={`border-b border-l border-neutral-100 p-1 align-top ${isWeekend ? "bg-neutral-50/40" : ""}`}>
                          <DroppableCell workerId={w.id} dateStr={dateStr}>
                            {dayShifts.length === 0 && (
                              <div
                                className="h-full min-h-[48px] flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer"
                                onClick={() => onCellClick(w.id, dateStr)}
                              >
                                <div className="w-6 h-6 rounded-full bg-neutral-200 flex items-center justify-center text-neutral-400 text-xs font-bold">+</div>
                              </div>
                            )}

                            {dayShifts.map((sh, si) => {
                              // Rest gap
                              let restGap: number | null = null;
                              if (si === 0 && isFlagged) {
                                const wi = workerShifts.indexOf(sh);
                                if (wi > 0) {
                                  const pE = new Date(workerShifts[wi - 1].end);
                                  const cS = new Date(sh.start);
                                  if (pE >= dS && pE <= cS) {
                                    const g = (cS.getTime() - pE.getTime()) / 3600000;
                                    if (g > 0 && g < 12) restGap = g;
                                  }
                                }
                              }

                              return (
                                <div key={sh._uid || si}>
                                  {restGap !== null && (
                                    <div className="text-[9px] font-bold text-red-500 text-center py-0.5 mb-1 rounded bg-red-50 border border-dashed border-red-200/80">
                                      {restGap.toFixed(0)}h rest
                                    </div>
                                  )}
                                  <DraggableShift
                                    shift={sh}
                                    worker={w}
                                    isFlagged={isFlagged}
                                    onShiftClick={onShiftClick}
                                  />
                                </div>
                              );
                            })}

                            {dayShifts.length > 0 && (
                              <div
                                className="h-4 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer"
                                onClick={() => onCellClick(w.id, dateStr)}
                              >
                                <div className="w-4 h-4 rounded-full bg-neutral-100 flex items-center justify-center text-neutral-300 text-[10px] font-bold">+</div>
                              </div>
                            )}
                          </DroppableCell>
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

      <DragOverlay>
        {activeShift && activeWorker && (
          <div
            className="rounded-lg p-2 text-white text-[11px] leading-tight shadow-xl opacity-90 w-[100px]"
            style={{ backgroundColor: activeWorker.color }}
          >
            <div className="font-semibold text-[12px] tabular-nums">
              {fmtTime(new Date(activeShift.start))}<span className="opacity-50 mx-0.5">to</span>{fmtTime(new Date(activeShift.end))}
            </div>
            <div className="opacity-60 text-[10px] mt-0.5">{shiftDur(activeShift).toFixed(0)}h</div>
          </div>
        )}
      </DragOverlay>
    </DndContext>
  );
}
