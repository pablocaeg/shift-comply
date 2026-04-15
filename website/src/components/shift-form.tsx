"use client";

import { useState } from "react";
import type { Shift } from "@/lib/types";
import { addDays, hoursBetween } from "@/lib/dates";
import { Button } from "@/components/ui/button";

interface Props {
  shift: Shift;
  workerName: string;
  workers?: { id: string; name: string; type: string }[];
  onSave: (s: Shift) => void;
  onDelete?: () => void;
  onCancel: () => void;
}

export function ShiftForm({ shift, workerName, workers, onSave, onDelete, onCancel }: Props) {
  const dateStr = shift.start.slice(0, 10);
  const dur = hoursBetween(shift.start, shift.end);

  const [start, setStart] = useState(shift.start.slice(11, 16) || "08:00");
  const [end, setEnd] = useState(shift.end.slice(11, 16) || "20:00");
  const [onCall, setOnCall] = useState(shift.on_call || false);
  const [assignTo, setAssignTo] = useState(shift.staff_id);

  function save() {
    const endDate = end <= start ? addDays(dateStr, 1) : dateStr;
    const targetWorker = workers?.find(w => w.id === assignTo);
    onSave({
      ...shift,
      staff_id: assignTo,
      staff_type: targetWorker?.type || shift.staff_type,
      start: `${dateStr}T${start}:00`,
      end: `${endDate}T${end}:00`,
      on_call: onCall,
    });
  }

  return (
    <div className="space-y-5">
      <div>
        <div className="text-sm font-semibold text-neutral-900">{workerName}</div>
        <div className="text-xs text-neutral-400 mt-0.5">
          {dateStr}{onDelete ? ` \u00b7 ${Math.round(dur)}h shift` : ""}
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide mb-1.5 block">Start time</label>
          <input type="time" lang="en-GB" step="3600" value={start} onChange={e => setStart(e.target.value)}
            className="w-full font-mono text-base border border-neutral-200 rounded-lg px-3 py-2.5 bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-neutral-900 focus:bg-white focus:border-transparent transition-all" />
        </div>
        <div>
          <label className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide mb-1.5 block">End time</label>
          <input type="time" lang="en-GB" step="3600" value={end} onChange={e => setEnd(e.target.value)}
            className="w-full font-mono text-base border border-neutral-200 rounded-lg px-3 py-2.5 bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-neutral-900 focus:bg-white focus:border-transparent transition-all" />
        </div>
      </div>

      {workers && workers.length > 1 && (
        <div>
          <label className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide mb-1.5 block">Assigned to</label>
          <select value={assignTo} onChange={e => setAssignTo(e.target.value)}
            className="w-full text-sm border border-neutral-200 rounded-lg px-3 py-2.5 bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-neutral-900 focus:bg-white focus:border-transparent transition-all">
            {workers.map(w => <option key={w.id} value={w.id}>{w.name}</option>)}
          </select>
        </div>
      )}

      <label className="flex items-center gap-3 p-3 rounded-lg border border-neutral-200 cursor-pointer select-none hover:bg-neutral-50 transition-colors">
        <input type="checkbox" checked={onCall} onChange={e => setOnCall(e.target.checked)}
          className="w-4 h-4 rounded border-neutral-300 text-amber-500 focus:ring-amber-500" />
        <div>
          <div className="text-sm font-medium text-neutral-700">On-call guard</div>
          <div className="text-[11px] text-neutral-400">Mark this shift as a guard duty</div>
        </div>
      </label>

      <div className="flex gap-2 pt-1">
        <Button onClick={save} className="flex-1 h-10">Save changes</Button>
        {onDelete && <Button variant="destructive" onClick={onDelete} className="h-10">Delete</Button>}
        <Button variant="ghost" onClick={onCancel} className="h-10">Cancel</Button>
      </div>
    </div>
  );
}
