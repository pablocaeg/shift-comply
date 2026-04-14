"use client";

import { useEffect, useState, useCallback, useMemo } from "react";
import { loadWasm, isLoaded } from "@/lib/wasm";
import { SCENARIOS } from "@/lib/scenarios";
import type { Scenario, Shift, Jurisdiction, Rule, ComplianceReport, Violation } from "@/lib/types";
import { tagShifts, nextUid } from "@/lib/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { ScheduleBoard } from "@/components/schedule-board";
import { ViolationList, type FixedViolation } from "@/components/violation-list";

function shiftDuration(s: Shift) {
  return (new Date(s.end).getTime() - new Date(s.start).getTime()) / 3600000;
}

// Format Date to local ISO string (avoids toISOString UTC conversion)
function toLocalISO(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}

export default function Home() {
  const [loaded, setLoaded] = useState(false);
  const [scenario, setScenario] = useState<Scenario | null>(null);
  const [jurisdiction, setJurisdiction] = useState("");
  const [shifts, setShifts] = useState<Shift[]>([]);
  const [report, setReport] = useState<ComplianceReport | null>(null);
  const [jurisdictions, setJurisdictions] = useState<Jurisdiction[]>([]);
  const [editDialog, setEditDialog] = useState<{ uid: string; shift: Shift } | null>(null);
  const [addDialog, setAddDialog] = useState<{ workerId: string; date: string } | null>(null);
  const [fixedViolations, setFixedViolations] = useState<FixedViolation[]>([]);

  // Load WASM
  useEffect(() => {
    loadWasm().then(() => {
      setLoaded(true);
      const jj: Jurisdiction[] = JSON.parse(window.shiftcomply.jurisdictions());
      jj.sort((a, b) => a.code.localeCompare(b.code));
      setJurisdictions(jj);
    });
  }, []);

  // Validate schedule
  const validate = useCallback((jur: string, scope: string, s: Shift[]) => {
    if (!isLoaded() || s.length === 0) { setReport(null); return; }
    const r: ComplianceReport = JSON.parse(window.shiftcomply.validate(JSON.stringify({
      jurisdiction: jur, facility_scope: scope, shifts: s,
    })));
    setReport(r);
  }, []);

  // Initialize first scenario
  useEffect(() => {
    if (loaded && !scenario) {
      const s = SCENARIOS[0];
      setScenario(s);
      setJurisdiction(s.jurisdiction);
      const initial = tagShifts(JSON.parse(JSON.stringify(s.shifts)) as Shift[]);
      setShifts(initial);
      validate(s.jurisdiction, s.scope, initial);
    }
  }, [loaded, scenario, validate]);

  const pickScenario = (s: Scenario) => {
    setScenario(s);
    setJurisdiction(s.jurisdiction);
    setFixedViolations([]);
    const next = tagShifts(JSON.parse(JSON.stringify(s.shifts)) as Shift[]);
    setShifts(next);
    validate(s.jurisdiction, s.scope, next);
  };

  const switchJurisdiction = (jur: string) => {
    setJurisdiction(jur);
    if (scenario) validate(jur, scenario.scope, shifts);
  };

  // Shift mutations (all uid-based)
  const updateShift = (uid: string, updated: Shift) => {
    const next = shifts.map(s => s._uid === uid ? { ...updated, _uid: uid } : s);
    setShifts(next);
    if (scenario) validate(jurisdiction, scenario.scope, next);
    setEditDialog(null);
  };

  const deleteShift = (uid: string) => {
    const next = shifts.filter(s => s._uid !== uid);
    setShifts(next);
    if (scenario) validate(jurisdiction, scenario.scope, next);
    setEditDialog(null);
  };

  const addShift = (shift: Shift) => {
    const next = [...shifts, { ...shift, _uid: nextUid() }];
    setShifts(next);
    if (scenario) validate(jurisdiction, scenario.scope, next);
    setAddDialog(null);
  };

  const moveShift = (uid: string, newWorkerId: string) => {
    const worker = scenario?.workers.find(w => w.id === newWorkerId);
    if (!worker) return;
    const next = shifts.map(s => s._uid === uid ? { ...s, staff_id: newWorkerId, staff_type: worker.type } : s);
    setShifts(next);
    if (scenario) validate(jurisdiction, scenario.scope, next);
  };

  // Auto-fix
  const autoFix = (violation: Violation) => {
    const k = violation.rule_key;
    const sid = violation.staff_id;
    const next = [...shifts];

    if (k.includes("max-weekly") || k.includes("max-combined") || k.includes("max-ordinary") || k.includes("days-off") || k.includes("day-of-rest")) {
      // Keep removing the last shift until the violation would be resolved
      const excessHours = violation.actual - violation.limit;
      let removed = 0;
      while (removed < excessHours + 14) { // safety margin
        const staffShifts = next.map((s, i) => ({ s, i })).filter(x => x.s.staff_id === sid).sort((a, b) => new Date(b.s.start).getTime() - new Date(a.s.start).getTime());
        if (!staffShifts.length) break;
        const dur = shiftDuration(staffShifts[0].s);
        next.splice(staffShifts[0].i, 1);
        removed += dur;
      }

    } else if (k.includes("rest-between") || k.includes("min-rest")) {
      // Fix ALL rest gaps for this worker by cascading forward
      const staffIndices = next.map((s, i) => ({ s, i })).filter(x => x.s.staff_id === sid).sort((a, b) => new Date(a.s.start).getTime() - new Date(b.s.start).getTime());

      for (let j = 1; j < staffIndices.length; j++) {
        const prevEnd = new Date(next[staffIndices[j - 1].i].end);
        const curStart = new Date(next[staffIndices[j].i].start);
        const gap = (curStart.getTime() - prevEnd.getTime()) / 3600000;

        if (gap >= 0 && gap < violation.limit) {
          const need = Math.ceil(violation.limit - gap);
          const dur = shiftDuration(next[staffIndices[j].i]);
          const newStart = new Date(prevEnd.getTime() + violation.limit * 3600000);
          const newEnd = new Date(newStart.getTime() + dur * 3600000);
          next[staffIndices[j].i] = {
            ...next[staffIndices[j].i],
            start: toLocalISO(newStart),
            end: toLocalISO(newEnd),
          };
        }
      }

    } else if (k.includes("max-shift")) {
      // Trim ALL over-long shifts for this worker
      for (let i = 0; i < next.length; i++) {
        if (next[i].staff_id !== sid) continue;
        const dur = shiftDuration(next[i]);
        if (dur > violation.limit) {
          const nE = new Date(new Date(next[i].start).getTime() + violation.limit * 3600000);
          next[i] = { ...next[i], end: toLocalISO(nE) };
        }
      }

    } else if (k.includes("guards") || k.includes("on-call")) {
      // Remove the last on-call shift
      const idx = next.map((s, i) => ({ s, i })).filter(x => x.s.staff_id === sid && x.s.on_call).sort((a, b) => new Date(b.s.start).getTime() - new Date(a.s.start).getTime())[0]?.i;
      if (idx !== undefined) next.splice(idx, 1);

    } else {
      // Fallback: remove the last shift
      const idx = next.map((s, i) => ({ s, i })).filter(x => x.s.staff_id === sid).sort((a, b) => new Date(b.s.start).getTime() - new Date(a.s.start).getTime())[0]?.i;
      if (idx !== undefined) next.splice(idx, 1);
    }

    // Record the fix
    setFixedViolations(prev => [...prev, {
      ruleKey: violation.rule_key,
      ruleName: violation.rule_name || violation.rule_key,
      staffId: violation.staff_id,
      fixedAt: Date.now(),
    }]);

    setShifts(next);
    if (scenario) validate(jurisdiction, scenario.scope, next);
  };

  const totalRules = useMemo(() => jurisdictions.reduce((sum, j) => sum + j.rules.length, 0), [jurisdictions]);

  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-white/90 backdrop-blur-sm border-b border-neutral-100">
        <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
          <div className="text-[15px] font-bold tracking-tight">
            shift-comply <span className="text-neutral-400 font-normal text-xs ml-1">v0.1.0</span>
          </div>
          <a href="https://github.com/pablocaeg/shift-comply" target="_blank" rel="noopener noreferrer" className="text-sm text-neutral-500 hover:text-neutral-900 transition-colors">
            GitHub
          </a>
        </div>
      </header>

      {/* Hero */}
      <section className="py-16 px-6 text-center">
        <h1 className="text-4xl md:text-5xl font-bold tracking-tight leading-tight max-w-2xl mx-auto mb-4">
          Know if your hospital schedule is legal
        </h1>
        <p className="text-neutral-500 text-base max-w-lg mx-auto mb-8 leading-relaxed">
          Healthcare labor laws vary across jurisdictions. ACGME caps residents at 80 hrs/week, California mandates nurse ratios, Spain limits guard duty. Shift Comply validates any schedule against the actual regulations, with legal citations.
        </p>
        {loaded ? (
          <div className="flex justify-center gap-3 flex-wrap">
            <div className="flex items-center gap-2 px-4 py-2 rounded-full border border-neutral-200 bg-neutral-50 text-sm">
              <span className="font-mono font-semibold">{totalRules}</span>
              <span className="text-neutral-500">verified regulations</span>
            </div>
            <div className="flex items-center gap-2 px-4 py-2 rounded-full border border-neutral-200 bg-neutral-50 text-sm">
              <span className="font-mono font-semibold">{jurisdictions.length}</span>
              <span className="text-neutral-500">jurisdictions covered</span>
            </div>
            <div className="flex items-center gap-2 px-4 py-2 rounded-full border border-neutral-200 bg-neutral-50 text-sm">
              <span className="font-mono font-semibold">US, EU, ES</span>
              <span className="text-neutral-500">regions</span>
            </div>
            <div className="flex items-center gap-2 px-4 py-2 rounded-full border border-neutral-200 bg-neutral-50 text-sm">
              <span className="font-mono font-semibold">100%</span>
              <span className="text-neutral-500">with legal citations</span>
            </div>
          </div>
        ) : (
          <div className="text-neutral-400 text-sm animate-pulse">Loading regulation database...</div>
        )}
      </section>

      {loaded && (
        <main className="max-w-6xl mx-auto px-6 pb-24">
          {/* Scenarios */}
          <section className="mb-10">
            <div className="text-[11px] font-semibold uppercase tracking-widest text-neutral-400 mb-2">Interactive Demo</div>
            <h2 className="text-xl font-bold tracking-tight mb-1">See it in action</h2>
            <p className="text-sm text-neutral-500 mb-6 max-w-lg">
              Select a scenario. Click any cell to add a shift, click a shift to edit it, or click Fix to auto-correct violations.
            </p>

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
              {SCENARIOS.map(s => (
                <button key={s.id} onClick={() => pickScenario(s)}
                  className={`text-left p-4 rounded-xl border transition-all duration-150 ${
                    scenario?.id === s.id
                      ? "border-neutral-900 shadow-sm ring-1 ring-neutral-900"
                      : "border-neutral-200 hover:border-neutral-300 hover:shadow-sm"
                  }`}
                >
                  <Badge variant={s.badge === "fail" ? "destructive" : "secondary"} className="mb-2 text-[10px]">{s.label}</Badge>
                  <div className="text-sm font-semibold mb-0.5">{s.name}</div>
                  <div className="text-xs text-neutral-500 mb-1.5">{s.who}</div>
                  <div className="text-[11px] text-neutral-400 leading-relaxed">{s.info}</div>
                </button>
              ))}
            </div>
          </section>

          {/* Schedule + Violations */}
          {scenario && report && (
            <section className="mb-16">
              {/* Board header */}
              <div className="flex items-center justify-between mb-4 flex-wrap gap-3">
                <div className="flex items-center gap-3">
                  <h3 className="text-base font-semibold">{scenario.who}</h3>
                  <Badge variant={report.result === "pass" ? "secondary" : "destructive"} className="text-xs gap-1.5">
                    <span className={`w-1.5 h-1.5 rounded-full ${report.result === "pass" ? "bg-emerald-500" : "bg-red-500"}`} />
                    {report.result === "pass" ? "Compliant" : `${report.violations.length} violation${report.violations.length > 1 ? "s" : ""}`}
                  </Badge>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide">Jurisdiction</span>
                  <Select value={jurisdiction} onValueChange={(v) => v && switchJurisdiction(v)}>
                    <SelectTrigger className="w-52 h-8 text-sm">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {jurisdictions.map(j => (
                        <SelectItem key={j.code} value={j.code}>{j.code} - {j.name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>

              {/* Board */}
              <ScheduleBoard
                scenario={scenario}
                shifts={shifts}
                report={report}
                onCellClick={(wid, date) => setAddDialog({ workerId: wid, date })}
                onShiftClick={(uid) => {
                  const s = shifts.find(x => x._uid === uid);
                  if (s) setEditDialog({ uid, shift: { ...s } });
                }}
                onMoveShift={moveShift}
              />

              {/* Violations */}
              <div className="mt-5">
                <ViolationList report={report} fixedViolations={fixedViolations} onFix={autoFix} />
              </div>
            </section>
          )}

          {/* Rule Explorer */}
          <RuleExplorer jurisdictions={jurisdictions} />
        </main>
      )}

      {/* Edit Shift Dialog */}
      <Dialog open={editDialog !== null} onOpenChange={() => setEditDialog(null)}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Edit shift</DialogTitle>
          </DialogHeader>
          {editDialog && (
            <ShiftForm
              shift={editDialog.shift}
              workerName={scenario?.workers.find(w => w.id === editDialog.shift.staff_id)?.name || ""}
              onSave={(s) => updateShift(editDialog.uid, s)}
              onDelete={() => deleteShift(editDialog.uid)}
              onCancel={() => setEditDialog(null)}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Add Shift Dialog */}
      <Dialog open={addDialog !== null} onOpenChange={() => setAddDialog(null)}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Add shift</DialogTitle>
          </DialogHeader>
          {addDialog && scenario && (
            <ShiftForm
              shift={{
                staff_id: addDialog.workerId,
                staff_type: scenario.workers.find(w => w.id === addDialog.workerId)?.type || "",
                start: addDialog.date + "T08:00:00",
                end: addDialog.date + "T20:00:00",
              }}
              workerName={scenario.workers.find(w => w.id === addDialog.workerId)?.name || ""}
              onSave={addShift}
              onCancel={() => setAddDialog(null)}
            />
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}

/* ---- Shift Form ---- */
function ShiftForm({ shift, workerName, onSave, onDelete, onCancel }: {
  shift: Shift; workerName: string;
  onSave: (s: Shift) => void; onDelete?: () => void; onCancel: () => void;
}) {
  const [startTime, setStartTime] = useState(shift.start.slice(11, 16));
  const [endTime, setEndTime] = useState(shift.end.slice(11, 16));
  const [onCall, setOnCall] = useState(shift.on_call || false);

  const handleSave = () => {
    const dateStr = shift.start.slice(0, 10);
    let endDate = dateStr;
    if (endTime <= startTime) {
      // Shift crosses midnight
      const d = new Date(dateStr + "T00:00:00");
      d.setDate(d.getDate() + 1);
      endDate = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
    }
    onSave({ ...shift, start: `${dateStr}T${startTime}:00`, end: `${endDate}T${endTime}:00`, on_call: onCall });
  };

  return (
    <div className="space-y-4">
      <div className="text-sm text-neutral-500">{workerName}</div>
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide mb-1 block">Start</label>
          <input type="time" value={startTime} onChange={e => setStartTime(e.target.value)}
            className="w-full font-mono text-sm border border-neutral-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-neutral-900 focus:border-transparent" />
        </div>
        <div>
          <label className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide mb-1 block">End</label>
          <input type="time" value={endTime} onChange={e => setEndTime(e.target.value)}
            className="w-full font-mono text-sm border border-neutral-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-neutral-900 focus:border-transparent" />
        </div>
      </div>
      <label className="flex items-center gap-2.5 text-sm text-neutral-600 cursor-pointer select-none">
        <input type="checkbox" checked={onCall} onChange={e => setOnCall(e.target.checked)} className="rounded border-neutral-300" />
        On-call guard
      </label>
      <div className="flex gap-2 pt-1">
        <Button onClick={handleSave} className="flex-1">Save</Button>
        {onDelete && <Button variant="destructive" onClick={onDelete}>Delete</Button>}
        <Button variant="ghost" onClick={onCancel}>Cancel</Button>
      </div>
    </div>
  );
}

/* ---- Rule Explorer ---- */
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

  const opStr = (op: string) => op === "lte" ? "\u2264" : op === "gte" ? "\u2265" : op === "eq" ? "=" : "";

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
            <SelectItem value="all">All staff</SelectItem>
            <SelectItem value="resident">Resident</SelectItem>
            <SelectItem value="nurse-rn">Nurse (RN)</SelectItem>
            <SelectItem value="statutory-personnel">Statutory</SelectItem>
            <SelectItem value="physician">Physician</SelectItem>
          </SelectContent>
        </Select>
        <Select value={cat || "all"} onValueChange={(v) => v && setCat(v === "all" ? "" : v)}>
          <SelectTrigger className="w-40 h-8 text-sm"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All categories</SelectItem>
            <SelectItem value="work_hours">Work Hours</SelectItem>
            <SelectItem value="rest">Rest</SelectItem>
            <SelectItem value="overtime">Overtime</SelectItem>
            <SelectItem value="staffing">Staffing</SelectItem>
            <SelectItem value="breaks">Breaks</SelectItem>
            <SelectItem value="on_call">On-Call</SelectItem>
            <SelectItem value="night_work">Night Work</SelectItem>
            <SelectItem value="leave">Leave</SelectItem>
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
              const v = r.values?.[0];
              if (!v) return null;
              const avg = v.averaged ? ` (avg ${v.averaged.count}${v.averaged.unit})` : "";
              return (
                <TableRow key={i}>
                  <TableCell>
                    <div className="font-mono text-xs font-medium">{r.key}</div>
                    <div className="text-[11px] text-neutral-500 mt-0.5">{r.name}</div>
                  </TableCell>
                  <TableCell className="font-mono text-xs">{opStr(r.operator)}{v.amount} {v.unit}</TableCell>
                  <TableCell className="font-mono text-[11px] text-neutral-500">{v.per}{avg}</TableCell>
                  <TableCell>{r.scope && <Badge variant="outline" className="text-[9px] font-mono">{r.scope}</Badge>}</TableCell>
                  <TableCell>
                    <Badge variant={r.enforcement === "mandatory" ? "destructive" : r.enforcement === "recommended" ? "secondary" : "outline"} className="text-[10px]">
                      {r.enforcement}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-[11px] text-neutral-500 max-w-[240px]">
                    {r.source.section ? `${r.source.title}, ${r.source.section}` : r.source.title}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </div>
    </section>
  );
}
