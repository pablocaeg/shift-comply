"use client";

import type { ComplianceReport, Violation } from "@/lib/types";
import { Button } from "@/components/ui/button";

export interface FixedViolation {
  ruleKey: string;
  ruleName: string;
  staffId: string;
  fixedAt: number; // timestamp
}

interface Props {
  report: ComplianceReport;
  fixedViolations: FixedViolation[];
  onFix: (violation: Violation) => void;
}

export function ViolationList({ report, fixedViolations, onFix }: Props) {
  const hasFixed = fixedViolations.length > 0;
  const allClear = report.result === "pass";

  return (
    <div className="space-y-3">
      {/* Fixed violations (resolved) */}
      {hasFixed && (
        <div className="space-y-1.5">
          {fixedViolations.map((fv, i) => (
            <div key={`fixed-${i}`}
              className="flex items-center gap-3 p-2.5 rounded-lg bg-emerald-50/60 border border-emerald-200/50 transition-all animate-in fade-in slide-in-from-top-1 duration-300"
            >
              <div className="w-5 h-5 rounded-full bg-emerald-500 flex items-center justify-center shrink-0">
                <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="3">
                  <path d="M20 6L9 17l-5-5" />
                </svg>
              </div>
              <div className="flex-1 min-w-0">
                <span className="text-[12px] font-medium text-emerald-800">{fv.ruleName}</span>
                <span className="text-[11px] text-emerald-600 ml-2">Fixed</span>
              </div>
              <div className="text-[10px] font-medium text-emerald-500 uppercase tracking-wide shrink-0">Resolved</div>
            </div>
          ))}
        </div>
      )}

      {/* Pass state */}
      {allClear && (
        <div className="flex items-center gap-3 p-4 rounded-xl bg-emerald-50 border border-emerald-200">
          <div className="w-8 h-8 rounded-full bg-emerald-100 flex items-center justify-center shrink-0">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" className="text-emerald-600">
              <path d="M20 6L9 17l-5-5" />
            </svg>
          </div>
          <div>
            <div className="text-sm font-semibold text-emerald-800">
              {hasFixed ? "All violations resolved" : "Schedule is fully compliant"}
            </div>
            <div className="text-xs text-emerald-600 mt-0.5">
              {report.constraints_checked} constraints checked.
              {hasFixed && ` ${fixedViolations.length} issue${fixedViolations.length > 1 ? "s" : ""} auto-fixed.`}
            </div>
          </div>
        </div>
      )}

      {/* Active violations */}
      {!allClear && (
        <>
          <div className="flex items-center justify-between">
            <div className="text-xs font-medium text-neutral-500 uppercase tracking-wide">
              {report.violations.length} violation{report.violations.length > 1 ? "s" : ""} remaining
            </div>
            <div className="text-[10px] text-neutral-400">{report.constraints_checked} constraints checked</div>
          </div>

          {report.violations.map((v, i) => (
            <div key={`${v.rule_key}-${v.staff_id}-${i}`}
              className="flex items-start gap-3 p-3 rounded-xl bg-red-50/80 border border-red-200/60 transition-all hover:border-red-300/80"
            >
              <div className="w-6 h-6 rounded-full bg-red-500 text-white flex items-center justify-center text-[10px] font-bold shrink-0 mt-0.5">
                {i + 1}
              </div>
              <div className="flex-1 min-w-0">
                <div className="text-[13px] font-semibold text-neutral-900">{v.rule_name || v.rule_key}</div>
                <div className="text-[11px] text-neutral-600 leading-relaxed mt-0.5">{v.message}</div>
                <div className="font-mono text-[10px] text-neutral-400 mt-1">{v.citation}</div>
              </div>
              <div className="flex items-center gap-3 shrink-0">
                <div className="text-right">
                  <div className="font-mono text-lg font-bold text-red-600 leading-none">{v.actual}</div>
                  <div className="font-mono text-[10px] text-neutral-400 mt-0.5">limit {v.limit}</div>
                </div>
                <Button
                  size="sm"
                  className="h-8 px-3 text-xs font-semibold bg-neutral-900 hover:bg-neutral-800"
                  onClick={() => onFix(v)}
                >
                  Fix
                </Button>
              </div>
            </div>
          ))}
        </>
      )}
    </div>
  );
}
