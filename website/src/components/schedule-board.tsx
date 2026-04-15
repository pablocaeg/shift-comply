"use client";

import { useMemo, useState } from "react";
import { DndContext, DragOverlay, useDraggable, useDroppable, type DragEndEvent, type DragStartEvent, PointerSensor, useSensor, useSensors } from "@dnd-kit/core";
import type { Scenario, Shift, ComplianceReport } from "@/lib/types";
import { dateOf, hoursBetween, mondayOf, addDays, dayName, dayNum, monthShort, dayOfWeek, shiftOverlapsDate } from "@/lib/dates";

const COLS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

interface Props {
  scenario: Scenario;
  shifts: Shift[];
  report: ComplianceReport;
  onCellClick: (workerId: string, date: string) => void;
  onShiftClick: (uid: string) => void;
  onMoveShift: (uid: string, toWorkerId: string) => void;
}

function DraggablePill({ uid, children }: { uid: string; children: React.ReactNode }) {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({ id: uid });
  return (
    <div ref={setNodeRef} {...listeners} {...attributes} className={isDragging ? "opacity-20 scale-95" : ""}>
      {children}
    </div>
  );
}

function DroppableDay({ id, children, className }: { id: string; children: React.ReactNode; className?: string }) {
  const { setNodeRef, isOver } = useDroppable({ id });
  return <div ref={setNodeRef} className={`${className} ${isOver ? "!bg-blue-50 !border-blue-300" : ""}`}>{children}</div>;
}

export function ScheduleBoard({ scenario, shifts, report, onCellClick, onShiftClick, onMoveShift }: Props) {
  const [dragUid, setDragUid] = useState<string | null>(null);
  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));

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
        ex.hours += h;
        if (sh.on_call) ex.onCall = true;
        ex.uids.push(sh._uid || "");
        map.set(key, ex);
      }
    }
    return map;
  }, [shifts, weeks]);

  const dragShift = dragUid ? shifts.find(s => s._uid === dragUid) : null;
  const dragWorker = dragShift ? scenario.workers.find(w => w.id === dragShift.staff_id) : null;

  return (
    <DndContext sensors={sensors} onDragStart={(e: DragStartEvent) => setDragUid(e.active.id as string)} onDragEnd={(e: DragEndEvent) => { setDragUid(null); if (e.over) { const tid = (e.over.id as string).split("|")[0]; if (tid && dragUid) onMoveShift(dragUid, tid); } }}>
      <div>
        {/* Day headers */}
        <div className="grid grid-cols-7 gap-1 mb-1">
          {COLS.map(c => <div key={c} className="text-center text-[10px] font-semibold text-neutral-400 uppercase tracking-wider py-1">{c}</div>)}
        </div>

        {/* Calendar weeks */}
        <div className="space-y-1">
          {weeks.map((week, wi) => (
            <div key={wi} className="grid grid-cols-7 gap-1">
              {week.map(d => {
                const isWe = dayOfWeek(d) === 0 || dayOfWeek(d) === 6;
                const hasPills = scenario.workers.some(w => cellData.has(`${w.id}|${d}`));

                return (
                  <DroppableDay key={d} id={`${scenario.workers[0]?.id || ""}|${d}`}
                    className={`rounded-lg border min-h-[90px] flex flex-col transition-colors ${isWe ? "border-neutral-100 bg-neutral-50/40" : "border-neutral-200 bg-white"}`}>
                    {/* Date */}
                    <div className="px-2 pt-1.5 flex items-baseline justify-between">
                      <span className={`text-[11px] font-semibold tabular-nums ${isWe ? "text-neutral-300" : "text-neutral-500"}`}>{dayNum(d)}</span>
                      {dayNum(d) <= 7 && wi === 0 && <span className="text-[8px] font-medium text-neutral-400 uppercase">{monthShort(d)}</span>}
                    </div>

                    {/* Shift pills */}
                    <div className="px-1 pb-1 flex flex-col gap-[3px] flex-1 mt-0.5">
                      {scenario.workers.map(w => {
                        const cell = cellData.get(`${w.id}|${d}`);
                        if (!cell || cell.hours <= 0) return null;
                        const flagged = violatedStaff.has(w.id);
                        const lastName = w.name.split(" ").pop() || w.name;

                        return (
                          <DraggablePill key={w.id} uid={cell.uids[0]}>
                            <button type="button"
                              onClick={e => { e.stopPropagation(); if (cell.uids[0]) onShiftClick(cell.uids[0]); }}
                              className={`
                                w-full rounded-[5px] px-1.5 py-1 text-[10px] font-semibold
                                flex flex-col cursor-pointer select-none
                                transition-all hover:brightness-110 active:scale-[0.97]
                                ${flagged ? "ring-1.5 ring-red-500 ring-offset-0.5" : ""}
                                ${cell.onCall
                                  ? "bg-amber-500 text-white border border-amber-600/30 bg-[repeating-linear-gradient(135deg,transparent,transparent_3px,rgba(255,255,255,0.12)_3px,rgba(255,255,255,0.12)_6px)]"
                                  : "text-white"}
                              `}
                              style={cell.onCall ? undefined : { backgroundColor: w.color }}
                            >
                              <div className="flex items-center gap-1 w-full">
                                <span className="truncate">{lastName}</span>
                                <span className="tabular-nums opacity-70 ml-auto shrink-0">{Math.round(cell.hours)}h</span>
                                {cell.onCall && <span className="text-[7px] font-bold opacity-60 shrink-0">G</span>}
                              </div>
                              <div className="text-[8px] opacity-50 tabular-nums">
                                {(() => { const sh = shifts.find(s => s._uid === cell.uids[0]); return sh ? `${sh.start.slice(11,16)}\u2013${sh.end.slice(11,16)}` : ""; })()}
                              </div>
                            </button>
                          </DraggablePill>
                        );
                      })}

                      {!hasPills && (
                        <button type="button" onClick={() => onCellClick(scenario.workers[0].id, d)}
                          className="flex-1 min-h-[24px] flex items-center justify-center text-neutral-200 text-[11px] opacity-0 hover:opacity-100 transition-opacity rounded hover:bg-neutral-100">
                          +
                        </button>
                      )}
                    </div>
                  </DroppableDay>
                );
              })}
            </div>
          ))}
        </div>

        {/* Legend */}
        <div className="flex items-center gap-4 flex-wrap pt-4 px-1">
          {scenario.workers.map(w => (
            <div key={w.id} className="flex items-center gap-1.5">
              <div className="w-3 h-3 rounded" style={{ backgroundColor: w.color }} />
              <span className={`text-[11px] ${violatedStaff.has(w.id) ? "text-red-600 font-semibold" : "text-neutral-500"}`}>
                {w.name}
              </span>
            </div>
          ))}
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded bg-amber-500 bg-[repeating-linear-gradient(135deg,transparent,transparent_2px,rgba(255,255,255,0.2)_2px,rgba(255,255,255,0.2)_4px)]" />
            <span className="text-[11px] text-neutral-500">On-call guard</span>
          </div>
        </div>
      </div>

      <DragOverlay>
        {dragShift && dragWorker && (
          <div className="rounded-md px-2 py-1 text-white text-[11px] font-bold shadow-2xl border border-white/20"
            style={{ backgroundColor: dragShift.on_call ? "#f59e0b" : dragWorker.color }}>
            {dragWorker.name.split(" ").pop()} {Math.round(hoursBetween(dragShift.start, dragShift.end))}h
          </div>
        )}
      </DragOverlay>
    </DndContext>
  );
}
