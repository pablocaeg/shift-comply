"use client";

import { useEffect, useState, useMemo } from "react";
import { loadWasm, isLoaded } from "@/lib/wasm";
import type { Jurisdiction, Rule } from "@/lib/types";
import type { JurisdictionInfo } from "@/lib/jurisdiction-data";
import { STATE_NAMES } from "@/lib/jurisdiction-data";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { USMap } from "@/components/us-map";

export default function Home() {
  const [loaded, setLoaded] = useState(false);
  const [jurisdictions, setJurisdictions] = useState<JurisdictionInfo[]>([]);
  const [allJurisdictions, setAllJurisdictions] = useState<Jurisdiction[]>([]);
  const [selected, setSelected] = useState<string | null>(null);

  useEffect(() => {
    loadWasm().then(() => {
      setLoaded(true);
      const jj: Jurisdiction[] = JSON.parse(window.shiftcomply.jurisdictions());
      jj.sort((a, b) => a.code.localeCompare(b.code));
      setAllJurisdictions(jj);
      setJurisdictions(jj.map(j => ({
        code: j.code, name: j.name, ruleCount: j.rules.length,
        parent: j.parent, type: j.type,
      })));
    });
  }, []);

  const totalRules = useMemo(() => jurisdictions.reduce((s, j) => s + j.ruleCount, 0), [jurisdictions]);

  // Rules for selected jurisdiction
  const selectedRules = useMemo(() => {
    if (!selected || !isLoaded()) return [];
    return JSON.parse(window.shiftcomply.rules(selected, "", "", "")) as Rule[];
  }, [selected]);

  const selectedInfo = useMemo(() => {
    return jurisdictions.find(j => j.code === selected);
  }, [selected, jurisdictions]);

  const op = (o: string) => o === "lte" ? "\u2264" : o === "gte" ? "\u2265" : o === "eq" ? "=" : "";

  return (
    <div className="min-h-screen bg-white">
      <header className="sticky top-0 z-50 bg-white/90 backdrop-blur-sm border-b border-neutral-100">
        <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
          <div className="text-[15px] font-bold tracking-tight">
            shift-comply <span className="text-neutral-400 font-normal text-xs ml-1">v0.1.0</span>
          </div>
          <a href="https://github.com/pablocaeg/shift-comply" target="_blank" rel="noopener noreferrer"
            className="text-sm text-neutral-500 hover:text-neutral-900 transition-colors">GitHub</a>
        </div>
      </header>

      <section className="py-12 px-6 text-center">
        <h1 className="text-4xl md:text-5xl font-bold tracking-tight leading-tight max-w-2xl mx-auto mb-4">
          Healthcare scheduling regulations, mapped
        </h1>
        <p className="text-neutral-500 text-base max-w-lg mx-auto mb-6 leading-relaxed">
          Every rule has a real legal citation. Click a state to see its regulations. Green states have state-specific rules beyond federal law.
        </p>
        {loaded ? (
          <div className="flex justify-center gap-3 flex-wrap mb-2">
            {[[String(totalRules), "verified regulations"], [String(jurisdictions.length), "jurisdictions"], ["100%", "with legal citations"]].map(([val, label]) => (
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
          {/* US Map */}
          <section className="mb-12">
            <USMap jurisdictions={jurisdictions} onSelect={setSelected} selected={selected} />
          </section>

          {/* Selected jurisdiction detail */}
          {selected && selectedInfo && (
            <section className="mb-12 animate-in fade-in slide-in-from-bottom-2 duration-300">
              <div className="flex items-center gap-3 mb-4">
                <h2 className="text-xl font-bold">{selectedInfo.name}</h2>
                <Badge variant="secondary" className="text-xs">{selectedInfo.code}</Badge>
                <Badge variant="outline" className="text-xs">{selectedInfo.ruleCount} rules</Badge>
                {selectedInfo.parent && (
                  <span className="text-xs text-neutral-400">inherits from {selectedInfo.parent}</span>
                )}
                <button onClick={() => setSelected(null)} className="ml-auto text-xs text-neutral-400 hover:text-neutral-700">Close</button>
              </div>

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
                    {selectedRules.map((r, i) => {
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
                          <TableCell className="font-mono text-xs">{op(r.operator)}{v.amount} {v.unit}</TableCell>
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
          )}

          {/* Non-US jurisdictions */}
          <section>
            <h2 className="text-lg font-bold mb-3">Other Jurisdictions</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              {jurisdictions.filter(j => !j.code.startsWith("US")).map(j => (
                <button key={j.code} onClick={() => setSelected(j.code)}
                  className={`text-left p-4 rounded-xl border transition-all ${selected === j.code ? "border-neutral-900 ring-1 ring-neutral-900" : "border-neutral-200 hover:border-neutral-300"}`}>
                  <div className="flex items-center gap-2 mb-1">
                    <span className="text-sm font-semibold">{j.name}</span>
                    <Badge variant="outline" className="text-[9px]">{j.code}</Badge>
                  </div>
                  <div className="text-xs text-neutral-400">{j.ruleCount} rules{j.parent ? ` (inherits from ${j.parent})` : ""}</div>
                </button>
              ))}
            </div>
          </section>
        </main>
      )}
    </div>
  );
}
