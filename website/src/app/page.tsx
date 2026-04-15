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

    const staff = () => next.filter(s => s.staff_id === sid).sort((a, b) => new Date(a.start).getTime() - new Date(b.start).getTime());
    const drop = (uid: string) => { next = next.filter(s => s._uid !== uid); };
    const patch = (uid: string, u: Partial<Shift>) => { next = next.map(s => s._uid === uid ? { ...s, ...u } : s); };

    if (k.includes("max-weekly") || k.includes("max-combined") || k.includes("max-ordinary") || k.includes("days-off") || k.includes("day-of-rest")) {
      let excess = v.actual - v.limit;
      while (excess > 0) {
        const ss = staff();
        if (!ss.length) break;
        const last = ss[ss.length - 1];
        excess -= shiftHours(last);
        drop(last._uid!);
      }
    } else if (k.includes("rest-between") || k.includes("min-rest")) {
      const ss = staff();
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
    } else if (k.includes("max-shift")) {
      for (const s of staff()) {
        if (shiftHours(s) > v.limit) {
          const trimEnd = new Date(new Date(s.start).getTime() + v.limit * 3600000);
          patch(s._uid!, { end: formatDateTime(trimEnd) });
        }
      }
    } else if (k.includes("guards") || k.includes("on-call")) {
      const oncalls = staff().filter(s => s.on_call);
      if (oncalls.length) drop(oncalls[oncalls.length - 1]._uid!);
    } else {
      const ss = staff();
      if (ss.length) drop(ss[ss.length - 1]._uid!);
    }

    setFixes(prev => [...prev, { ruleName: v.rule_name || v.rule_key, staffId: sid }]);
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
function ShiftForm({ shift, workerName, onSave, onDelete, onCancel }: {
  shift: Shift; workerName: string;
  onSave: (s: Shift) => void; onDelete?: () => void; onCancel: () => void;
}) {
  const [start, setStart] = useState(shift.start.slice(11, 16));
  const [end, setEnd] = useState(shift.end.slice(11, 16));
  const [onCall, setOnCall] = useState(shift.on_call || false);

  function save() {
    const dateStr = shift.start.slice(0, 10);
    const endDate = end <= start ? addDays(dateStr, 1) : dateStr;
    onSave({ ...shift, start: `${dateStr}T${start}:00`, end: `${endDate}T${end}:00`, on_call: onCall });
  }

  return (
    <div className="space-y-4">
      <div className="text-sm text-neutral-500">{workerName}</div>
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide mb-1 block">Start</label>
          <input type="time" value={start} onChange={e => setStart(e.target.value)}
            className="w-full font-mono text-sm border border-neutral-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-neutral-900 focus:border-transparent" />
        </div>
        <div>
          <label className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide mb-1 block">End</label>
          <input type="time" value={end} onChange={e => setEnd(e.target.value)}
            className="w-full font-mono text-sm border border-neutral-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-neutral-900 focus:border-transparent" />
        </div>
      </div>
      <label className="flex items-center gap-2.5 text-sm text-neutral-600 cursor-pointer select-none">
        <input type="checkbox" checked={onCall} onChange={e => setOnCall(e.target.checked)} className="rounded border-neutral-300" />
        On-call guard
      </label>
      <div className="flex gap-2 pt-1">
        <Button onClick={save} className="flex-1">Save</Button>
        {onDelete && <Button variant="destructive" onClick={onDelete}>Delete</Button>}
        <Button variant="ghost" onClick={onCancel}>Cancel</Button>
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
