"use client";

import { useEffect, useState } from "react";
import { isLoaded } from "@/lib/wasm";
import type { Jurisdiction, Rule } from "@/lib/types";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

const STAFF_OPTIONS = [["all", "All staff"], ["resident", "Resident"], ["nurse-rn", "Nurse (RN)"], ["statutory-personnel", "Statutory"], ["physician", "Physician"]];
const CATEGORY_OPTIONS = [["all", "All categories"], ["work_hours", "Work Hours"], ["rest", "Rest"], ["overtime", "Overtime"], ["staffing", "Staffing"], ["breaks", "Breaks"], ["on_call", "On-Call"], ["night_work", "Night Work"], ["leave", "Leave"]];

function operatorSymbol(op: string): string {
  return op === "lte" ? "\u2264" : op === "gte" ? "\u2265" : op === "eq" ? "=" : "";
}

export function RuleExplorer({ jurisdictions }: { jurisdictions: Jurisdiction[] }) {
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
          <SelectContent>{STAFF_OPTIONS.map(([v, l]) => <SelectItem key={v} value={v}>{l}</SelectItem>)}</SelectContent>
        </Select>
        <Select value={cat || "all"} onValueChange={(v) => v && setCat(v === "all" ? "" : v)}>
          <SelectTrigger className="w-40 h-8 text-sm"><SelectValue /></SelectTrigger>
          <SelectContent>{CATEGORY_OPTIONS.map(([v, l]) => <SelectItem key={v} value={v}>{l}</SelectItem>)}</SelectContent>
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
              const avgLabel = v.averaged ? ` (avg ${v.averaged.count}${v.averaged.unit})` : "";
              const citation = r.source.section ? `${r.source.title}, ${r.source.section}` : r.source.title;
              const enfVariant = r.enforcement === "mandatory" ? "destructive" as const : r.enforcement === "recommended" ? "secondary" as const : "outline" as const;
              return (
                <TableRow key={i}>
                  <TableCell>
                    <div className="font-mono text-xs font-medium">{r.key}</div>
                    <div className="text-[11px] text-neutral-500 mt-0.5">{r.name}</div>
                  </TableCell>
                  <TableCell className="font-mono text-xs">{operatorSymbol(r.operator)}{v.amount} {v.unit}</TableCell>
                  <TableCell className="font-mono text-[11px] text-neutral-500">{v.per}{avgLabel}</TableCell>
                  <TableCell>{r.scope && <Badge variant="outline" className="text-[9px] font-mono">{r.scope}</Badge>}</TableCell>
                  <TableCell><Badge variant={enfVariant} className="text-[10px]">{r.enforcement}</Badge></TableCell>
                  <TableCell className="text-[11px] text-neutral-500 max-w-[240px]">{citation}</TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </div>
    </section>
  );
}
