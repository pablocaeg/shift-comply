"use client";

import { useMemo, useState } from "react";
import { DndContext, DragOverlay, useDraggable, useDroppable, type DragEndEvent, type DragStartEvent } from "@dnd-kit/core";
import type { Scenario, Shift, ComplianceReport } from "@/lib/types";
import { dateOf, hoursBetween, mondayOf, addDays, dayName, dayNum, dayOfWeek, shiftOverlapsDate } from "@/lib/dates";
import { Badge } from "@/components/ui/badge";

interface Props {
  scenario: Scenario;
  shifts: Shift[];
  report: ComplianceReport;
  onCellClick: (workerId: string, date: string) => void;
  onShiftClick: (uid: string) => void;
  onMoveShift: (uid: string, toWorkerId: string, toDate: string) => void;
}

function DropZone({ id, children, className }: { id: string; children: React.ReactNode; className?: string }) {
  const { setNodeRef, isOver } = useDroppable({ id });
  return <div ref={setNodeRef} className={`${className || ""} ${isOver ? "!bg-blue-100" : ""}`}>{children}</div>;
}

function DragCell({ uid, children }: { uid: string; children: React.ReactNode }) {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({ id: uid });
  return (
    <div ref={setNodeRef} {...listeners} {...attributes} className={`cursor-grab active:cursor-grabbing ${isDragging ? "opacity-20" : ""}`}>
      {children}
    </div>
  );
}

export function ScheduleBoard({ scenario, shifts, report, onCellClick, onShiftClick, onMoveShift }: Props) {
  const [dragUid, setDragUid] = useState<string | null>(null);
  const violatedStaff = useMemo(() => new Set((report.violations || []).map(v => v.staff_id)), [report]);
  const vCounts = useMemo(() => {
    const m = new Map<string, number>();
    for (const v of report.violations || []) m.set(v.staff_id, (m.get(v.staff_id) || 0) + 1);
    return m;
  }, [report]);

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

  const cellData = useMemo(() => {
    const map = new Map<string, { hours: number; onCall: boolean; uids: string[] }>();
    for (const sh of shifts) {
      for (const d of gridDates) {
        if (shiftOverlapsDate(sh.start, sh.end, d)) {
          const dayStartMs = new Date(+d.slice(0,4), +d.slice(5,7)-1, +d.slice(8,10), 0).getTime();
          const dayEndMs = dayStartMs + 86400000;
          const sMs = new Date(+sh.start.slice(0,4), +sh.start.slice(5,7)-1, +sh.start.slice(8,10), +(sh.start.slice(11,13)||0), +(sh.start.slice(14,16)||0)).getTime();
          const eMs = new Date(+sh.end.slice(0,4), +sh.end.slice(5,7)-1, +sh.end.slice(8,10), +(sh.end.slice(11,13)||0), +(sh.end.slice(14,16)||0)).getTime();
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

  const dragShift = dragUid ? shifts.find(s => s._uid === dragUid) : null;
  const dragWorker = dragShift ? scenario.workers.find(w => w.id === dragShift.staff_id) : null;

  function onDragStart(e: DragStartEvent) { setDragUid(e.active.id as string); }
  function onDragEnd(e: DragEndEvent) {
    setDragUid(null);
    if (!e.over) return;
    const [toWorker, toDate] = (e.over.id as string).split("|");
    if (toWorker && toDate && dragUid) onMoveShift(dragUid, toWorker, toDate);
  }

  const cols = gridDates.length;

  return (
    <DndContext onDragStart={onDragStart} onDragEnd={onDragEnd}>
      <div className="border border-neutral-200 rounded-xl overflow-hidden bg-white">
        {/* Header */}
        <div className="grid" style={{ gridTemplateColumns: `140px repeat(${cols}, 1fr) 52px` }}>
          <div className="bg-neutral-50 text-[9px] font-semibold text-neutral-400 uppercase tracking-wider p-2 border-b border-r border-neutral-200">Staff</div>
          {gridDates.map((d, i) => {
            const isWe = dayOfWeek(d) === 0 || dayOfWeek(d) === 6;
            const isMonday = i > 0 && dayOfWeek(d) === 1;
            return (
              <div key={d} className={`text-center text-[9px] font-medium p-1 border-b border-neutral-200 ${isWe ? "bg-neutral-100/60 text-neutral-300" : "bg-neutral-50 text-neutral-400"} ${isMonday ? "border-l-2 border-l-neutral-300" : "border-l border-l-neutral-100"}`}>
                <div className="font-semibold">{dayName(d).slice(0, 2)}</div>
                <div>{dayNum(d)}</div>
              </div>
            );
          })}
          <div className="bg-neutral-50 text-[9px] font-semibold text-neutral-400 uppercase tracking-wider p-1 border-b border-l-2 border-neutral-200 text-center">Tot</div>
        </div>

        {/* Rows */}
        {scenario.workers.map(worker => {
          const flagged = violatedStaff.has(worker.id);
          const vc = vCounts.get(worker.id) || 0;
          const totalH = shifts.filter(s => s.staff_id === worker.id).reduce((sum, s) => sum + hoursBetween(s.start, s.end), 0);

          return (
            <div key={worker.id} className={`grid ${flagged ? "bg-red-50/30" : ""}`} style={{ gridTemplateColumns: `140px repeat(${cols}, 1fr) 52px` }}>
              {/* Name */}
              <div className={`flex items-center gap-0 border-b border-r border-neutral-100 ${flagged ? "bg-red-50" : "bg-white"}`}>
                <div className="w-1 self-stretch shrink-0" style={{ backgroundColor: flagged ? "#ef4444" : worker.color }} />
                <div className="flex items-center gap-1.5 px-2 py-1.5 min-w-0 flex-1">
                  <div className="min-w-0">
                    <div className={`text-[11px] font-semibold leading-tight truncate ${flagged ? "text-red-700" : ""}`}>{worker.name}</div>
                    <div className="text-[8px] text-neutral-400 truncate">{worker.role}</div>
                  </div>
                  {flagged && <Badge variant="destructive" className="text-[7px] px-1 py-0 h-3 shrink-0 ml-auto">{vc}</Badge>}
                </div>
              </div>

              {/* Day cells */}
              {gridDates.map((d, i) => {
                const key = `${worker.id}|${d}`;
                const cell = cellData.get(key);
                const isWe = dayOfWeek(d) === 0 || dayOfWeek(d) === 6;
                const isMonday = i > 0 && dayOfWeek(d) === 1;
                const hasShift = cell && cell.hours > 0;

                return (
                  <DropZone key={d} id={`${worker.id}|${d}`}
                    className={`border-b border-neutral-100 flex items-center justify-center min-h-[40px] ${isMonday ? "border-l-2 border-l-neutral-300" : "border-l border-l-neutral-50"} ${!hasShift && isWe ? "bg-neutral-50/50" : !hasShift ? "bg-white" : ""}`}
                  >
                    {hasShift ? (
                      <DragCell uid={cell.uids[0]}>
                        <button
                          type="button"
                          onClick={e => { e.stopPropagation(); if (cell.uids[0]) onShiftClick(cell.uids[0]); }}
                          className={`w-full min-h-[40px] flex flex-col items-center justify-center text-white text-[11px] font-bold tabular-nums transition-all hover:brightness-110 ${flagged ? "ring-2 ring-inset ring-red-500" : ""}`}
                          style={{ backgroundColor: cell.onCall ? "#f59e0b" : worker.color }}
                        >
                          <span>{Math.round(cell.hours)}h</span>
                          {cell.onCall && <span className="text-[6px] font-bold uppercase opacity-60 leading-none">guard</span>}
                        </button>
                      </DragCell>
                    ) : (
                      <button type="button" onClick={() => onCellClick(worker.id, d)}
                        className="w-full h-full min-h-[40px] flex items-center justify-center text-neutral-300 text-[10px] opacity-0 hover:opacity-100 hover:bg-neutral-50 transition-all">
                        +
                      </button>
                    )}
                  </DropZone>
                );
              })}

              {/* Total */}
              <div className="border-b border-l-2 border-neutral-200 flex items-center justify-center bg-white">
                <span className={`text-[11px] font-bold tabular-nums ${flagged ? "text-red-600" : "text-neutral-600"}`}>{Math.round(totalH)}h</span>
              </div>
            </div>
          );
        })}
      </div>

      <DragOverlay>
        {dragShift && dragWorker && (
          <div className="rounded px-2 py-1 text-white text-[11px] font-bold shadow-xl" style={{ backgroundColor: dragShift.on_call ? "#f59e0b" : dragWorker.color }}>
            {Math.round(hoursBetween(dragShift.start, dragShift.end))}h
          </div>
        )}
      </DragOverlay>
    </DndContext>
  );
}
