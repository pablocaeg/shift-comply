"use client";

import { useEffect, useState, useCallback, useMemo } from "react";
import { loadWasm, isLoaded } from "@/lib/wasm";
import { SCENARIOS } from "@/lib/scenarios";
import { applyFix } from "@/lib/autofix";
import type { Scenario, Shift, Jurisdiction, ComplianceReport, Violation } from "@/lib/types";
import { tagShifts, nextUid } from "@/lib/types";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ScheduleBoard } from "@/components/schedule-board";
import { ViolationList, type FixRecord } from "@/components/violation-list";
import { ShiftForm } from "@/components/shift-form";
import { RuleExplorer } from "@/components/rule-explorer";

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

  // Load WASM on mount
  useEffect(() => {
    loadWasm().then(() => {
      setLoaded(true);
      const jj: Jurisdiction[] = JSON.parse(window.shiftcomply.jurisdictions());
      jj.sort((a, b) => a.code.localeCompare(b.code));
      setJurisdictions(jj);
    });
  }, []);

  // Validate shifts against jurisdiction rules via WASM
  const validate = useCallback((jur: string, scope: string, s: Shift[]) => {
    if (!isLoaded() || !s.length) { setReport(null); return; }
    const clean = s.map(({ _uid, ...rest }) => rest);
    const r: ComplianceReport = JSON.parse(
      window.shiftcomply.validate(JSON.stringify({ jurisdiction: jur, facility_scope: scope, shifts: clean }))
    );
    setReport(r);
  }, []);

  // Auto-select first scenario on load
  useEffect(() => {
    if (loaded && !scenario) {
      pickScenario(SCENARIOS[0]);
    }
  }, [loaded]); // eslint-disable-line react-hooks/exhaustive-deps

  // ---- Schedule mutations ----

  function commit(next: Shift[]) {
    setShifts(next);
    if (scenario) validate(jurisdiction, scenario.scope, next);
  }

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

  function saveShift(uid: string, updated: Shift) {
    commit(shifts.map(s => s._uid === uid ? { ...updated, _uid: uid } : s));
    setEditUid(null);
  }

  function removeShift(uid: string) {
    commit(shifts.filter(s => s._uid !== uid));
    setEditUid(null);
  }

  function createShift(shift: Shift) {
    commit([...shifts, { ...shift, _uid: nextUid() }]);
    setAddTarget(null);
  }

  function handleFix(v: Violation) {
    const result = applyFix(shifts, v, scenario?.workers || []);
    setFixes(prev => [...prev, { ruleName: result.description, staffId: v.staff_id }]);
    commit(result.shifts);
  }

  // ---- Derived state ----

  const totalRules = useMemo(() => jurisdictions.reduce((sum, j) => sum + j.rules.length, 0), [jurisdictions]);
  const editShift = editUid ? shifts.find(s => s._uid === editUid) || null : null;
  const fixedStaff = useMemo(() => new Set(fixes.map(f => f.staffId)), [fixes]);

  // ---- Render ----

  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-white/90 backdrop-blur-sm border-b border-neutral-100">
        <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
          <div className="text-[15px] font-bold tracking-tight">
            shift-comply <span className="text-neutral-400 font-normal text-xs ml-1">v0.1.0</span>
          </div>
          <a href="https://github.com/pablocaeg/shift-comply" target="_blank" rel="noopener noreferrer"
            className="text-sm text-neutral-500 hover:text-neutral-900 transition-colors">GitHub</a>
        </div>
      </header>

      {/* Hero */}
      <section className="py-16 px-6 text-center">
        <h1 className="text-4xl md:text-5xl font-bold tracking-tight leading-tight max-w-2xl mx-auto mb-4">
          Know if your hospital schedule is legal
        </h1>
        <p className="text-neutral-500 text-base max-w-lg mx-auto mb-8 leading-relaxed">
          Scheduling systems let you set constraints. Shift Comply tells you what those constraints should be, based on the actual law. Select a jurisdiction, get the legally correct values with citations. No manual research.
        </p>
        {loaded ? (
          <div className="flex justify-center gap-3 flex-wrap">
            {[[String(totalRules), "verified regulations"], [String(jurisdictions.length), "jurisdictions"], ["US, EU, ES", "regions"], ["100%", "with legal citations"]].map(([val, label]) => (
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
            <p className="text-sm text-neutral-500 mb-6 max-w-lg">
              Select a scenario. Click any cell to add a shift, click a shift to edit or delete it, click Fix to auto-correct violations.
            </p>
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

          {/* Schedule + Violations */}
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
                scenario={scenario} shifts={shifts} report={report} fixedStaff={fixedStaff}
                onCellClick={(wid, date) => setAddTarget({ workerId: wid, date })}
                onShiftClick={uid => setEditUid(uid)}
              />

              <div className="mt-4">
                <ViolationList report={report} fixes={fixes} onFix={handleFix} />
              </div>
            </section>
          )}

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
                workers={scenario.workers}
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
