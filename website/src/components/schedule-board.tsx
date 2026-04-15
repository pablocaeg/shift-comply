"use client";

import { useMemo, useState } from "react";
import { DndContext, DragOverlay, useDraggable, useDroppable, type DragEndEvent, type DragStartEvent } from "@dnd-kit/core";
import type { Scenario, Shift, ComplianceReport } from "@/lib/types";
import { dateOf, timeOf, hoursBetween, mondayOf, addDays, dayName, monthShort, dayNum, dayOfWeek, shiftOverlapsDate, clampToDay } from "@/lib/dates";
import { Badge } from "@/components/ui/badge";

interface Props {
  scenario: Scenario;
  shifts: Shift[];
  report: ComplianceReport;
  onCellClick: (workerId: string, date: string) => void;
  onShiftClick: (uid: string) => void;
  onMoveShift: (uid: string, toWorkerId: string) => void;
}

// ---- Draggable shift block ----
function ShiftBlock({ shift, color, flagged, onShiftClick }: {
  shift: Shift; color: string; flagged: boolean; onShiftClick: (uid: string) => void;
  dateStr: string;
}) {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({ id: shift._uid || "", data: { shift } });
  const dur = hoursBetween(shift.start, shift.end);
  const vis = clampToDay(shift.start, shift.end, dateStr);

  return (
    <div
      ref={setNodeRef} {...listeners} {...attributes}
      onClick={e => { e.stopPropagation(); if (shift._uid) onShiftClick(shift._uid); }}
      role="button" tabIndex={0}
      className={`
        w-full text-left rounded-lg p-2 mb-1 last:mb-0 text-white text-[11px]
        cursor-grab active:cursor-grabbing select-none
        transition-all hover:brightness-110 hover:shadow-md active:scale-[0.97]
        ${isDragging ? "opacity-30" : ""}
        ${flagged ? "ring-2 ring-red-400 ring-offset-1 ring-offset-white" : ""}
        ${shift.on_call ? "bg-stripes" : ""}
      `}
      style={{ backgroundColor: color }}
    >
      {shift.on_call && <div className="text-[8px] font-bold uppercase tracking-wider opacity-50 mb-0.5">On-call</div>}
      <div className="font-semibold text-[12px] tabular-nums tracking-tight">
        {vis.visStart}<span className="opacity-30 mx-1">&ndash;</span>{vis.visEnd}
      </div>
      <div className="opacity-50 text-[10px]">
        {dur === vis.visDur ? `${Math.round(dur)}h` : `${Math.round(vis.visDur)}h of ${Math.round(dur)}h`}
      </div>
    </div>
  );
}

// ---- Droppable cell ----
function DropCell({ id, children }: { id: string; children: React.ReactNode }) {
  const { setNodeRef, isOver } = useDroppable({ id });
  return (
    <div ref={setNodeRef} className={`min-h-[52px] rounded p-0.5 transition-colors ${isOver ? "bg-blue-50 ring-2 ring-blue-300 ring-inset" : "group"}`}>
      {children}
    </div>
  );
}

// ---- Need this in scope for ShiftBlock ----
let dateStr = "";

export function ScheduleBoard({ scenario, shifts, report, onCellClick, onShiftClick, onMoveShift }: Props) {
  const [dragShift, setDragShift] = useState<Shift | null>(null);

  const violatedStaff = useMemo(() => new Set((report.violations || []).map(v => v.staff_id)), [report]);
  const vCounts = useMemo(() => {
    const m = new Map<string, number>();
    for (const v of report.violations || []) m.set(v.staff_id, (m.get(v.staff_id) || 0) + 1);
    return m;
  }, [report]);

  // Compute grid dates (pure strings)
  const gridDates = useMemo(() => {
    if (!shifts.length) return [];
    const allDates = shifts.map(s => dateOf(s.start));
    const endDates = shifts.map(s => dateOf(s.end));
    const earliest = [...allDates, ...endDates].sort()[0];
    const latest = [...allDates, ...endDates].sort().pop()!;
    const monday = mondayOf(earliest);
    const dates: string[] = [];
    let cur = monday;
    while (cur <= latest || dates.length % 7 !== 0) {
      dates.push(cur);
      cur = addDays(cur, 1);
    }
    if (dates.length < 7) while (dates.length < 7) { dates.push(cur); cur = addDays(cur, 1); }
    return dates;
  }, [shifts]);

  function onDragStart(e: DragStartEvent) {
    setDragShift((e.active.data.current as { shift: Shift })?.shift || null);
  }

  function onDragEnd(e: DragEndEvent) {
    setDragShift(null);
    if (!e.over) return;
    const uid = e.active.id as string;
    const [targetWorker] = (e.over.id as string).split("|");
    if (targetWorker && uid) onMoveShift(uid, targetWorker);
  }

  const dragWorker = dragShift ? scenario.workers.find(w => w.id === dragShift.staff_id) : null;

  return (
    <DndContext onDragStart={onDragStart} onDragEnd={onDragEnd}>
      <div className="border border-neutral-200 rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full border-collapse" style={{ minWidth: Math.max(720, gridDates.length * 100 + 180) }}>
            <thead>
              <tr>
                <th className="bg-neutral-50 text-left text-[11px] font-medium text-neutral-500 uppercase tracking-wide p-3 pl-4 border-b border-neutral-200 w-44 min-w-[176px] sticky left-0 z-10">Staff</th>
                {gridDates.map(d => {
                  const we = dayOfWeek(d) === 0 || dayOfWeek(d) === 6;
                  return (
                    <th key={d} className={`text-center text-[11px] font-medium p-2 border-b border-l border-neutral-200 min-w-[90px] ${we ? "bg-neutral-100/50 text-neutral-400" : "bg-neutral-50 text-neutral-500"}`}>
                      <div className="font-semibold">{dayName(d)}</div>
                      <div className="font-normal text-neutral-400 text-[10px]">{dayNum(d)} {monthShort(d)}</div>
                    </th>
                  );
                })}
              </tr>
            </thead>
            <tbody>
              {scenario.workers.map(w => {
                const flagged = violatedStaff.has(w.id);
                const vc = vCounts.get(w.id) || 0;
                const ws = shifts.filter(s => s.staff_id === w.id).sort((a, b) => a.start.localeCompare(b.start));

                return (
                  <tr key={w.id} className={flagged ? "bg-red-50/20" : ""}>
                    <td className={`p-3 pl-4 border-b border-neutral-100 sticky left-0 z-10 bg-white ${flagged ? "border-l-[3px] border-l-red-400 !bg-red-50/60" : ""}`}>
                      <div className="flex items-center gap-2">
                        <div className="w-2.5 h-2.5 rounded-full shrink-0" style={{ background: w.color }} />
                        <div className="min-w-0 flex-1">
                          <div className={`text-[13px] font-semibold leading-tight truncate ${flagged ? "text-red-700" : ""}`}>{w.name}</div>
                          <div className="text-[10px] text-neutral-400">{w.role}</div>
                        </div>
                        {flagged && <Badge variant="destructive" className="text-[9px] px-1.5 py-0 h-4 shrink-0">{vc}</Badge>}
                      </div>
                    </td>

                    {gridDates.map(d => {
                      const we = dayOfWeek(d) === 0 || dayOfWeek(d) === 6;
                      const dayShifts = ws.filter(sh => shiftOverlapsDate(sh.start, sh.end, d));
                      dateStr = d; // set for ShiftBlock closure

                      return (
                        <td key={d} className={`border-b border-l border-neutral-100 p-1 align-top ${we ? "bg-neutral-50/30" : ""}`}>
                          <DropCell id={`${w.id}|${d}`}>
                            {dayShifts.length === 0 && (
                              <div className="h-[48px] flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer"
                                onClick={() => onCellClick(w.id, d)}>
                                <div className="w-5 h-5 rounded bg-neutral-200/80 flex items-center justify-center text-neutral-400 text-[11px]">+</div>
                              </div>
                            )}
                            {dayShifts.map((sh, si) => {
                              // Rest gap
                              let gap: number | null = null;
                              if (flagged && si === 0) {
                                const wi = ws.indexOf(sh);
                                if (wi > 0) {
                                  const g = hoursBetween(ws[wi - 1].end, sh.start);
                                  if (g > 0 && g < 12 && dateOf(ws[wi - 1].end) === d) gap = g;
                                }
                              }
                              return (
                                <div key={sh._uid || si}>
                                  {gap !== null && (
                                    <div className="text-[9px] font-bold text-red-500 text-center py-0.5 mb-1 rounded bg-red-50 border border-dashed border-red-200">
                                      {Math.round(gap)}h rest
                                    </div>
                                  )}
                                  <ShiftBlock shift={sh} color={w.color} flagged={flagged} onShiftClick={onShiftClick} dateStr={d} />
                                </div>
                              );
                            })}
                            {dayShifts.length > 0 && (
                              <div className="opacity-0 group-hover:opacity-100 transition-opacity flex justify-center pt-0.5 cursor-pointer" onClick={() => onCellClick(w.id, d)}>
                                <div className="w-4 h-4 rounded bg-neutral-100 flex items-center justify-center text-neutral-300 text-[9px]">+</div>
                              </div>
                            )}
                          </DropCell>
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
        {dragShift && dragWorker && (
          <div className="rounded-lg p-2 text-white text-[11px] shadow-2xl w-24 opacity-90" style={{ backgroundColor: dragWorker.color }}>
            <div className="font-semibold tabular-nums">{timeOf(dragShift.start)}<span className="opacity-30 mx-0.5">&ndash;</span>{timeOf(dragShift.end)}</div>
            <div className="opacity-50 text-[10px]">{Math.round(hoursBetween(dragShift.start, dragShift.end))}h</div>
          </div>
        )}
      </DragOverlay>
    </DndContext>
  );
}
