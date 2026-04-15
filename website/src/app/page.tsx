"use client";

import { useEffect, useState, useCallback, useMemo } from "react";
import { loadWasm, isLoaded } from "@/lib/wasm";
import { SCENARIOS } from "@/lib/scenarios";
import { addDays, hoursBetween, formatDateTime } from "@/lib/dates";
import type { Scenario, Shift, Jurisdiction, Rule, ComplianceReport, Violation } from "@/lib/types";
import { tagShifts, nextUid } from "@/lib/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { ScheduleBoard } from "@/components/schedule-board";
import { ViolationList, type FixRecord } from "@/components/violation-list";

// ---- Helpers ----
function shiftHours(s: Shift): number {
  return hoursBetween(s.start, s.end);
}

// ---- Main page ----
export default function Home() {
  const [loaded, setLoaded] = useState(false);
  const [jurisdictions, setJurisdictions] = useState<Jurisdiction[]>([]);
  const [scenario, setScenario] = useState<Scenario | null>(null);
  const [jurisdiction, setJurisdiction] = useState("");
  const [shifts, setShifts] = useState<Shift[]>([]);
  const [report, setReport] = useState<ComplianceReport | null>(null);
  const [fixes, setFixes] = useState<FixRecord[]>([]);
  const [editUid, setEditUid] = useState<string | null>(null);
  const [addTarget, setAddTarget] = useState<{ workerId: string; date: string } | null>(null);

  // ---- WASM ----
  useEffect(() => {
    loadWasm().then(() => {
      setLoaded(true);
      const jj: Jurisdiction[] = JSON.parse(window.shiftcomply.jurisdictions());
      jj.sort((a, b) => a.code.localeCompare(b.code));
      setJurisdictions(jj);
    });
  }, []);

  const validate = useCallback((jur: string, scope: string, s: Shift[]) => {
    if (!isLoaded() || !s.length) { setReport(null); return; }
    // Strip _uid before sending to WASM (Go doesn't know about it)
    const clean = s.map(({ _uid, ...rest }) => rest);
    const r: ComplianceReport = JSON.parse(
      window.shiftcomply.validate(JSON.stringify({ jurisdiction: jur, facility_scope: scope, shifts: clean }))
    );
    setReport(r);
  }, []);

  // ---- Init ----
  useEffect(() => {
    if (loaded && !scenario) {
      const s = SCENARIOS[0];
      setScenario(s);
      setJurisdiction(s.jurisdiction);
      const tagged = tagShifts(structuredClone(s.shifts));
      setShifts(tagged);
      validate(s.jurisdiction, s.scope, tagged);
    }
  }, [loaded, scenario, validate]);

  // ---- Actions ----
  function pickScenario(s: Scenario) {
    setScenario(s);
    setJurisdiction(s.jurisdiction);
    setFixes([]);
    setEditUid(null);
    setAddTarget(null);
    const tagged = tagShifts(structuredClone(s.shifts));
    setShifts(tagged);
    validate(s.jurisdiction, s.scope, tagged);
  }

  function switchJurisdiction(jur: string) {
    setJurisdiction(jur);
    setFixes([]);
    if (scenario) validate(jur, scenario.scope, shifts);
  }

  function applyShifts(next: Shift[]) {
    setShifts(next);
    if (scenario) validate(jurisdiction, scenario.scope, next);
  }

  function saveShift(uid: string, updated: Shift) {
    applyShifts(shifts.map(s => s._uid === uid ? { ...updated, _uid: uid } : s));
    setEditUid(null);
  }

  function removeShift(uid: string) {
    applyShifts(shifts.filter(s => s._uid !== uid));
    setEditUid(null);
  }

  function createShift(shift: Shift) {
    applyShifts([...shifts, { ...shift, _uid: nextUid() }]);
    setAddTarget(null);
  }

  // ---- Auto-fix (all uid-based) ----
  function autoFix(v: Violation) {
    const sid = v.staff_id;
    const k = v.rule_key;
    let next = [...shifts];
    const workers = scenario?.workers || [];

    const staffShifts = (id: string) => next.filter(s => s.staff_id === id).sort((a, b) => a.start.localeCompare(b.start));
    const patch = (uid: string, u: Partial<Shift>) => { next = next.map(s => s._uid === uid ? { ...s, ...u } : s); };
    let fixDesc = v.rule_name || v.rule_key;

    if (k.includes("max-weekly") || k.includes("max-combined") || k.includes("max-ordinary")) {
      // Shorten the longest shifts to reduce weekly hours
      const excess = Math.ceil(v.actual - v.limit);
      const ss = staffShifts(sid).sort((a, b) => shiftHours(b) - shiftHours(a)); // longest first
      let reduced = 0;
      for (const s of ss) {
        if (reduced >= excess) break;
        const dur = shiftHours(s);
        const cut = Math.min(dur - 4, excess - reduced); // don't go below 4h shifts
        if (cut > 0) {
          const trimEnd = new Date(new Date(s.start).getTime() + (dur - cut) * 3600000);
          patch(s._uid!, { end: formatDateTime(trimEnd) });
          reduced += cut;
        }
      }
      fixDesc = `Shortened shifts by ${reduced}h to meet ${v.limit}h/week limit`;

    } else if (k.includes("days-off") || k.includes("day-of-rest")) {
      // Remove the shortest shift to create a day off
      const ss = staffShifts(sid).sort((a, b) => shiftHours(a) - shiftHours(b));
      if (ss.length) {
        next = next.filter(s => s._uid !== ss[0]._uid);
        fixDesc = `Removed ${Math.round(shiftHours(ss[0]))}h shift to create a day off`;
      }

    } else if (k.includes("rest-between") || k.includes("min-rest")) {
      // Push shifts forward to create enough rest
      const ss = staffShifts(sid);
      for (let i = 1; i < ss.length; i++) {
        const prev = next.find(s => s._uid === ss[i - 1]._uid)!;
        const cur = next.find(s => s._uid === ss[i]._uid)!;
        const gap = (new Date(cur.start).getTime() - new Date(prev.end).getTime()) / 3600000;
        if (gap >= 0 && gap < v.limit) {
          const dur = shiftHours(cur);
          const newStart = new Date(new Date(prev.end).getTime() + v.limit * 3600000);
          const newEnd = new Date(newStart.getTime() + dur * 3600000);
          patch(cur._uid!, { start: formatDateTime(newStart), end: formatDateTime(newEnd) });
        }
      }
      fixDesc = `Adjusted shift times to ensure ${v.limit}h rest between shifts`;

    } else if (k.includes("max-shift")) {
      // Trim to limit
      for (const s of staffShifts(sid)) {
        if (shiftHours(s) > v.limit) {
          const trimEnd = new Date(new Date(s.start).getTime() + v.limit * 3600000);
          patch(s._uid!, { end: formatDateTime(trimEnd) });
        }
      }
      fixDesc = `Trimmed shifts to ${v.limit}h maximum`;

    } else if (k.includes("guards") || k.includes("on-call")) {
      // Reassign excess guards to the worker with fewest guards
      const myGuards = staffShifts(sid).filter(s => s.on_call);
      const excess = myGuards.length - Math.floor(v.limit);
      if (excess > 0) {
        // Find who has capacity
        const guardCounts = new Map<string, number>();
        for (const w of workers) {
          guardCounts.set(w.id, next.filter(s => s.staff_id === w.id && s.on_call).length);
        }
        // Take the last N excess guards and reassign
        const toReassign = myGuards.slice(-excess);
        for (const guard of toReassign) {
          // Find worker with fewest guards (excluding current)
          let minWorker = workers.find(w => w.id !== sid) || workers[0];
          let minCount = Infinity;
          for (const w of workers) {
            if (w.id === sid) continue;
            const count = guardCounts.get(w.id) || 0;
            if (count < minCount) { minCount = count; minWorker = w; }
          }
          patch(guard._uid!, { staff_id: minWorker.id, staff_type: minWorker.type });
          guardCounts.set(minWorker.id, (guardCounts.get(minWorker.id) || 0) + 1);
          guardCounts.set(sid, (guardCounts.get(sid) || 0) - 1);
        }
        fixDesc = `Reassigned ${excess} guard${excess > 1 ? "s" : ""} to ${workers.find(w => w.id !== sid)?.name || "colleague"}`;
      }

    } else {
      // Fallback: shorten the last shift
      const ss = staffShifts(sid);
      if (ss.length) {
        const last = ss[ss.length - 1];
        const trimEnd = new Date(new Date(last.start).getTime() + Math.min(shiftHours(last), 8) * 3600000);
        patch(last._uid!, { end: formatDateTime(trimEnd) });
        fixDesc = `Shortened last shift`;
      }
    }

    setFixes(prev => [...prev, { ruleName: fixDesc, staffId: sid }]);
    applyShifts(next);
  }

  // ---- Derived ----
  const totalRules = useMemo(() => jurisdictions.reduce((sum, j) => sum + j.rules.length, 0), [jurisdictions]);
  const editShift = editUid ? shifts.find(s => s._uid === editUid) || null : null;

  // ---- Render ----
  return (
    <div className="min-h-screen bg-white">
      <header className="sticky top-0 z-50 bg-white/90 backdrop-blur-sm border-b border-neutral-100">
        <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
          <div className="text-[15px] font-bold tracking-tight">shift-comply <span className="text-neutral-400 font-normal text-xs ml-1">v0.1.0</span></div>
          <a href="https://github.com/pablocaeg/shift-comply" target="_blank" rel="noopener noreferrer" className="text-sm text-neutral-500 hover:text-neutral-900 transition-colors">GitHub</a>
        </div>
      </header>

      <section className="py-16 px-6 text-center">
        <h1 className="text-4xl md:text-5xl font-bold tracking-tight leading-tight max-w-2xl mx-auto mb-4">Know if your hospital schedule is legal</h1>
        <p className="text-neutral-500 text-base max-w-lg mx-auto mb-8 leading-relaxed">
          Scheduling systems let you set constraints. Shift Comply tells you what those constraints should be, based on the actual law. Select a jurisdiction, get the legally correct values with citations. No manual research.
        </p>
        {loaded ? (
          <div className="flex justify-center gap-3 flex-wrap">
            {[
              [String(totalRules), "verified regulations"],
              [String(jurisdictions.length), "jurisdictions"],
              ["US, EU, ES", "regions"],
              ["100%", "with legal citations"],
            ].map(([val, label]) => (
              <div key={label} className="flex items-center gap-2 px-4 py-2 rounded-full border border-neutral-200 bg-neutral-50 text-sm">
                <span className="font-mono font-semibold">{val}</span>
                <span className="text-neutral-500">{label}</span>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-neutral-400 text-sm animate-pulse">Loading regulation database...</div>
        )}
      </section>

      {loaded && (
        <main className="max-w-6xl mx-auto px-6 pb-24">
          {/* Scenario picker */}
          <section className="mb-10">
            <div className="text-[11px] font-semibold uppercase tracking-widest text-neutral-400 mb-2">Interactive Demo</div>
            <h2 className="text-xl font-bold tracking-tight mb-1">See it in action</h2>
            <p className="text-sm text-neutral-500 mb-6 max-w-lg">Select a scenario. Click any cell to add a shift, click a shift to edit or delete it, click Fix to auto-correct violations.</p>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
              {SCENARIOS.map(s => (
                <button key={s.id} onClick={() => pickScenario(s)}
                  className={`text-left p-4 rounded-xl border transition-all ${scenario?.id === s.id ? "border-neutral-900 ring-1 ring-neutral-900 shadow-sm" : "border-neutral-200 hover:border-neutral-300 hover:shadow-sm"}`}>
                  <Badge variant={s.badge === "fail" ? "destructive" : "secondary"} className="mb-2 text-[10px]">{s.label}</Badge>
                  <div className="text-sm font-semibold mb-0.5">{s.name}</div>
                  <div className="text-xs text-neutral-500 mb-1.5">{s.who}</div>
                  <div className="text-[11px] text-neutral-400 leading-relaxed">{s.info}</div>
                </button>
              ))}
            </div>
          </section>

          {/* Schedule board + violations */}
          {scenario && report && (
            <section className="mb-16">
              <div className="flex items-center justify-between mb-4 flex-wrap gap-3">
                <div className="flex items-center gap-3">
                  <h3 className="text-base font-semibold">{scenario.who}</h3>
                  <Badge variant={report.result === "pass" ? "secondary" : "destructive"} className="text-xs gap-1.5">
                    <span className={`w-1.5 h-1.5 rounded-full ${report.result === "pass" ? "bg-emerald-500" : "bg-red-500"}`} />
                    {report.result === "pass" ? "Compliant" : `${report.violations.length} violation${report.violations.length !== 1 ? "s" : ""}`}
                  </Badge>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide">Jurisdiction</span>
                  <Select value={jurisdiction} onValueChange={(v) => v && switchJurisdiction(v)}>
                    <SelectTrigger className="w-52 h-8 text-sm"><SelectValue /></SelectTrigger>
                    <SelectContent>{jurisdictions.map(j => <SelectItem key={j.code} value={j.code}>{j.code} - {j.name}</SelectItem>)}</SelectContent>
                  </Select>
                </div>
              </div>

              <ScheduleBoard
                scenario={scenario} shifts={shifts} report={report}
                fixedStaff={new Set(fixes.map(f => f.staffId))}
                onCellClick={(wid, date) => setAddTarget({ workerId: wid, date })}
                onShiftClick={uid => setEditUid(uid)}
              />

              <div className="mt-4">
                <ViolationList report={report} fixes={fixes} onFix={autoFix} />
              </div>
            </section>
          )}

          {/* Rule Explorer */}
          <RuleExplorer jurisdictions={jurisdictions} />
        </main>
      )}

      {/* Edit dialog */}
      <Dialog open={editShift !== null} onOpenChange={() => setEditUid(null)}>
        <DialogContent className="max-w-sm">
          <DialogHeader><DialogTitle>Edit shift</DialogTitle></DialogHeader>
          {editShift && (
            <ShiftForm
              shift={editShift}
              workerName={scenario?.workers.find(w => w.id === editShift.staff_id)?.name || ""}
              workers={scenario?.workers}
              onSave={s => saveShift(editUid!, s)}
              onDelete={() => removeShift(editUid!)}
              onCancel={() => setEditUid(null)}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Add dialog */}
      <Dialog open={addTarget !== null} onOpenChange={() => setAddTarget(null)}>
        <DialogContent className="max-w-sm">
          <DialogHeader><DialogTitle>Add shift</DialogTitle></DialogHeader>
          {addTarget && scenario && (() => {
            const w = scenario.workers.find(x => x.id === addTarget.workerId);
            return (
              <ShiftForm
                shift={{ staff_id: addTarget.workerId, staff_type: w?.type || "", start: `${addTarget.date}T08:00:00`, end: `${addTarget.date}T20:00:00` }}
                workerName={w?.name || ""}
                onSave={createShift}
                onCancel={() => setAddTarget(null)}
              />
            );
          })()}
        </DialogContent>
      </Dialog>
    </div>
  );
}

// ---- Shift Form ----
function ShiftForm({ shift, workerName, workers, onSave, onDelete, onCancel }: {
  shift: Shift; workerName: string; workers?: { id: string; name: string; type: string }[];
  onSave: (s: Shift) => void; onDelete?: () => void; onCancel: () => void;
}) {
  const startTime = shift.start.slice(11, 16) || "08:00";
  const endTime = shift.end.slice(11, 16) || "20:00";
  const dateStr = shift.start.slice(0, 10);
  const dur = hoursBetween(shift.start, shift.end);

  const [start, setStart] = useState(startTime);
  const [end, setEnd] = useState(endTime);
  const [onCall, setOnCall] = useState(shift.on_call || false);
  const [assignTo, setAssignTo] = useState(shift.staff_id);

  function save() {
    const eDate = end <= start ? addDays(dateStr, 1) : dateStr;
    const targetWorker = workers?.find(w => w.id === assignTo);
    onSave({
      ...shift,
      staff_id: assignTo,
      staff_type: targetWorker?.type || shift.staff_type,
      start: `${dateStr}T${start}:00`,
      end: `${eDate}T${end}:00`,
      on_call: onCall,
    });
  }

  return (
    <div className="space-y-5">
      {/* Header info */}
      <div>
        <div className="text-sm font-semibold text-neutral-900">{workerName}</div>
        <div className="text-xs text-neutral-400 mt-0.5">
          {dateStr} {onDelete ? `\u00b7 ${Math.round(dur)}h shift` : ""}
        </div>
      </div>

      {/* Time inputs */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide mb-1.5 block">Start time</label>
          <input type="time" lang="en-GB" value={start} onChange={e => setStart(e.target.value)} step="3600"
            className="w-full font-mono text-base border border-neutral-200 rounded-lg px-3 py-2.5 bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-neutral-900 focus:bg-white focus:border-transparent transition-all" />
        </div>
        <div>
          <label className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide mb-1.5 block">End time</label>
          <input type="time" lang="en-GB" value={end} onChange={e => setEnd(e.target.value)} step="3600"
            className="w-full font-mono text-base border border-neutral-200 rounded-lg px-3 py-2.5 bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-neutral-900 focus:bg-white focus:border-transparent transition-all" />
        </div>
      </div>

      {/* Assign to worker */}
      {workers && workers.length > 1 && (
        <div>
          <label className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide mb-1.5 block">Assigned to</label>
          <select value={assignTo} onChange={e => setAssignTo(e.target.value)}
            className="w-full text-sm border border-neutral-200 rounded-lg px-3 py-2.5 bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-neutral-900 focus:bg-white focus:border-transparent transition-all">
            {workers.map(w => <option key={w.id} value={w.id}>{w.name}</option>)}
          </select>
        </div>
      )}

      {/* On-call toggle */}
      <label className="flex items-center gap-3 p-3 rounded-lg border border-neutral-200 cursor-pointer select-none hover:bg-neutral-50 transition-colors">
        <input type="checkbox" checked={onCall} onChange={e => setOnCall(e.target.checked)}
          className="w-4 h-4 rounded border-neutral-300 text-amber-500 focus:ring-amber-500" />
        <div>
          <div className="text-sm font-medium text-neutral-700">On-call guard</div>
          <div className="text-[11px] text-neutral-400">Mark this shift as a guard duty</div>
        </div>
      </label>

      {/* Actions */}
      <div className="flex gap-2 pt-1">
        <Button onClick={save} className="flex-1 h-10">Save changes</Button>
        {onDelete && <Button variant="destructive" onClick={onDelete} className="h-10">Delete</Button>}
        <Button variant="ghost" onClick={onCancel} className="h-10">Cancel</Button>
      </div>
    </div>
  );
}

// ---- Rule Explorer ----
function RuleExplorer({ jurisdictions }: { jurisdictions: Jurisdiction[] }) {
  const [jur, setJur] = useState("US-CA");
  const [staff, setStaff] = useState("");
  const [cat, setCat] = useState("");
  const [rules, setRules] = useState<Rule[]>([]);

  useEffect(() => {
    if (!isLoaded()) return;
    let r: Rule[] = JSON.parse(window.shiftcomply.rules(jur, staff, "", ""));
    if (cat) r = r.filter(rule => rule.category === cat);
    setRules(r);
  }, [jur, staff, cat]);

  const op = (o: string) => o === "lte" ? "\u2264" : o === "gte" ? "\u2265" : o === "eq" ? "=" : "";

  return (
    <section id="explorer">
      <div className="text-[11px] font-semibold uppercase tracking-widest text-neutral-400 mb-2">Rule Explorer</div>
      <h2 className="text-xl font-bold tracking-tight mb-1">Browse all regulations</h2>
      <p className="text-sm text-neutral-500 mb-6 max-w-lg">Filter by jurisdiction, staff type, and category.</p>
      <div className="flex gap-2 mb-4 flex-wrap">
        <Select value={jur} onValueChange={(v) => v && setJur(v)}>
          <SelectTrigger className="w-52 h-8 text-sm"><SelectValue /></SelectTrigger>
          <SelectContent>{jurisdictions.map(j => <SelectItem key={j.code} value={j.code}>{j.code} - {j.name}</SelectItem>)}</SelectContent>
        </Select>
        <Select value={staff || "all"} onValueChange={(v) => v && setStaff(v === "all" ? "" : v)}>
          <SelectTrigger className="w-40 h-8 text-sm"><SelectValue /></SelectTrigger>
          <SelectContent>
            {[["all","All staff"],["resident","Resident"],["nurse-rn","Nurse (RN)"],["statutory-personnel","Statutory"],["physician","Physician"]].map(([v,l]) => (
              <SelectItem key={v} value={v}>{l}</SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select value={cat || "all"} onValueChange={(v) => v && setCat(v === "all" ? "" : v)}>
          <SelectTrigger className="w-40 h-8 text-sm"><SelectValue /></SelectTrigger>
          <SelectContent>
            {[["all","All categories"],["work_hours","Work Hours"],["rest","Rest"],["overtime","Overtime"],["staffing","Staffing"],["breaks","Breaks"],["on_call","On-Call"],["night_work","Night Work"],["leave","Leave"]].map(([v,l]) => (
              <SelectItem key={v} value={v}>{l}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="text-xs text-neutral-400 mb-3">{rules.length} rules</div>
      <div className="border border-neutral-200 rounded-xl overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="text-[10px]">Rule</TableHead>
              <TableHead className="text-[10px]">Value</TableHead>
              <TableHead className="text-[10px]">Per</TableHead>
              <TableHead className="text-[10px]">Scope</TableHead>
              <TableHead className="text-[10px]">Enforcement</TableHead>
              <TableHead className="text-[10px]">Citation</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rules.map((r, i) => {
              const v = r.values?.[0]; if (!v) return null;
              return (
                <TableRow key={i}>
                  <TableCell><div className="font-mono text-xs font-medium">{r.key}</div><div className="text-[11px] text-neutral-500 mt-0.5">{r.name}</div></TableCell>
                  <TableCell className="font-mono text-xs">{op(r.operator)}{v.amount} {v.unit}</TableCell>
                  <TableCell className="font-mono text-[11px] text-neutral-500">{v.per}{v.averaged ? ` (avg ${v.averaged.count}${v.averaged.unit})` : ""}</TableCell>
                  <TableCell>{r.scope && <Badge variant="outline" className="text-[9px] font-mono">{r.scope}</Badge>}</TableCell>
                  <TableCell><Badge variant={r.enforcement === "mandatory" ? "destructive" : r.enforcement === "recommended" ? "secondary" : "outline"} className="text-[10px]">{r.enforcement}</Badge></TableCell>
                  <TableCell className="text-[11px] text-neutral-500 max-w-[240px]">{r.source.section ? `${r.source.title}, ${r.source.section}` : r.source.title}</TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </div>
    </section>
  );
}
