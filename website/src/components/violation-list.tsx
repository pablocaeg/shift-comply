"use client";

import type { ComplianceReport, Violation } from "@/lib/types";
import { Button } from "@/components/ui/button";

export interface FixRecord {
  ruleName: string;
  staffId: string;
}

interface Props {
  report: ComplianceReport;
  fixes: FixRecord[];
  onFix: (violation: Violation) => void;
}

export function ViolationList({ report, fixes, onFix }: Props) {
  const isPass = report.result === "pass";

  return (
    <div className="space-y-2">
      {/* Resolved fixes */}
      {fixes.map((f, i) => (
        <div key={`fix-${i}`} className="flex items-center gap-2.5 px-3 py-2 rounded-lg bg-emerald-50 border border-emerald-200/60 text-[12px]">
          <div className="w-5 h-5 rounded-full bg-emerald-500 flex items-center justify-center shrink-0">
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="3"><path d="M20 6L9 17l-5-5" /></svg>
          </div>
          <span className="font-medium text-emerald-800">{f.ruleName}</span>
          <span className="text-emerald-600 text-[11px] ml-auto">Fixed</span>
        </div>
      ))}

      {/* Pass state */}
      {isPass && (
        <div className="flex items-center gap-3 p-4 rounded-xl bg-emerald-50 border border-emerald-200">
          <div className="w-8 h-8 rounded-full bg-emerald-100 flex items-center justify-center shrink-0">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" className="text-emerald-600"><path d="M20 6L9 17l-5-5" /></svg>
          </div>
          <div>
            <div className="text-sm font-semibold text-emerald-800">
              {fixes.length > 0 ? "All violations resolved" : "Schedule is fully compliant"}
            </div>
            <div className="text-xs text-emerald-600">{report.constraints_checked} constraints checked.{fixes.length > 0 ? ` ${fixes.length} fixed.` : ""}</div>
          </div>
        </div>
      )}

      {/* Active violations */}
      {!isPass && (
        <>
          <div className="flex items-center justify-between py-1">
            <span className="text-[11px] font-medium text-neutral-500 uppercase tracking-wide">
              {report.violations.length} violation{report.violations.length !== 1 ? "s" : ""}
            </span>
            <span className="text-[10px] text-neutral-400">{report.constraints_checked} checked</span>
          </div>

          {report.violations.map((v, i) => (
            <div key={`v-${i}`} className="flex items-start gap-3 p-3 rounded-xl bg-red-50/70 border border-red-200/60">
              <div className="w-6 h-6 rounded-full bg-red-500 text-white flex items-center justify-center text-[10px] font-bold shrink-0 mt-0.5">
                {i + 1}
              </div>
              <div className="flex-1 min-w-0">
                <div className="text-[13px] font-semibold text-neutral-900 leading-tight">{v.rule_name || v.rule_key}</div>
                <div className="text-[11px] text-neutral-600 leading-relaxed mt-0.5">{v.message}</div>
                <div className="font-mono text-[10px] text-neutral-400 mt-1">{v.citation}</div>
              </div>
              <div className="flex items-center gap-2.5 shrink-0 mt-0.5">
                <div className="text-right">
                  <div className="font-mono text-base font-bold text-red-600 leading-none">{v.actual}</div>
                  <div className="font-mono text-[10px] text-neutral-400">limit {v.limit}</div>
                </div>
                <Button size="sm" className="h-7 px-3 text-[11px] font-semibold" onClick={() => onFix(v)}>
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
