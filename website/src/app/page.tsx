"use client";

import { useEffect, useRef, useState, useMemo } from "react";
import { loadWasm, isLoaded } from "@/lib/wasm";
import type { Jurisdiction, Rule } from "@/lib/types";
import { JURISDICTION_STATS, type JurisdictionInfo } from "@/lib/jurisdiction-data";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { USMap } from "@/components/us-map";
import { WorldMap } from "@/components/world-map";

type MapView = "us" | "world";

export default function Home() {
  const [loaded, setLoaded] = useState(false);
  const [jurisdictions, setJurisdictions] = useState<JurisdictionInfo[]>([]);
  const [selected, setSelected] = useState<string | null>(null);
  const [view, setView] = useState<MapView>("us");
  const [showInherited, setShowInherited] = useState(false);
  const detailRef = useRef<HTMLElement>(null);

  useEffect(() => {
    loadWasm().then(() => {
      setLoaded(true);
      const jj: Jurisdiction[] = JSON.parse(window.shiftcomply.jurisdictions());
      jj.sort((a, b) => a.code.localeCompare(b.code));
      setJurisdictions(jj.map(j => ({
        code: j.code, name: j.name, ruleCount: j.rules.length,
        parent: j.parent, type: j.type,
      })));
    });
  }, []);

  const totalRules = useMemo(() => jurisdictions.reduce((s, j) => s + j.ruleCount, 0), [jurisdictions]);

  // Aggregate healthcare sector stats for current view.
  // Only count top-level jurisdictions to avoid double-counting.
  const viewStats = useMemo(() => {
    const codes = view === "us"
      ? jurisdictions.filter(j => j.code === "US" || j.code.startsWith("US-")).map(j => j.code)
      : jurisdictions.map(j => j.code);
    // Only count top-level to avoid double-counting (US already includes US-CA, etc.)
    const topLevel = new Set(codes.filter(c => {
      const j = jurisdictions.find(x => x.code === c);
      return j && (!j.parent || !codes.includes(j.parent));
    }));
    let hospitals = 0;
    let workers = 0;
    for (const code of topLevel) {
      const s = JURISDICTION_STATS[code];
      if (s) { hospitals += s.hospitals; workers += s.healthcareWorkers; }
    }
    return { hospitals, workers };
  }, [jurisdictions, view]);

  const selectedRules = useMemo(() => {
    if (!selected || !isLoaded()) return [];
    const parsed = JSON.parse(window.shiftcomply.rules(selected, "", "", ""));
    return Array.isArray(parsed) ? parsed as Rule[] : [];
  }, [selected]);

  const selectedInfo = useMemo(() => jurisdictions.find(j => j.code === selected), [selected, jurisdictions]);

  // Keys of rules defined directly on this jurisdiction (not inherited)
  const localRuleKeys = useMemo(() => {
    if (!selected || !isLoaded()) return new Set<string>();
    const jur = JSON.parse(window.shiftcomply.export(selected));
    if (!jur?.rules) return new Set<string>();
    return new Set((jur as Jurisdiction).rules.map(r => r.key));
  }, [selected]);

  const filteredRules = useMemo(() => {
    if (showInherited) return selectedRules;
    return selectedRules.filter(r => localRuleKeys.has(r.key));
  }, [selectedRules, showInherited, localRuleKeys]);

  // Scroll to detail section when a jurisdiction is selected
  useEffect(() => {
    if (selected && detailRef.current) {
      detailRef.current.scrollIntoView({ behavior: "smooth", block: "start" });
    }
  }, [selected]);

  // When switching views, clear selection
  function switchView(v: MapView) {
    setView(v);
    setSelected(null);
  }

  const op = (o: string) => o === "lte" ? "\u2264" : o === "gte" ? "\u2265" : o === "eq" ? "=" : "";

  return (
    <div className="min-h-screen bg-white flex flex-col">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-white/90 backdrop-blur-sm border-b border-neutral-100">
        <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
          <div className="text-[15px] font-bold tracking-tight">
            shift-comply <span className="text-neutral-400 font-normal text-xs ml-1">v0.1.0</span>
          </div>
          <div className="flex items-center gap-4">
            {loaded && (
              <div className="hidden sm:flex items-center gap-1 text-xs text-neutral-400">
                <span className="font-mono font-semibold text-neutral-600">{totalRules}</span> rules
                <span className="mx-1">|</span>
                <span className="font-mono font-semibold text-neutral-600">{jurisdictions.length}</span> jurisdictions
              </div>
            )}
            <a href="https://github.com/pablocaeg/shift-comply" target="_blank" rel="noopener noreferrer"
              className="text-sm text-neutral-500 hover:text-neutral-900 transition-colors">GitHub</a>
          </div>
        </div>
      </header>

      {/* Hero */}
      <section className="pt-10 pb-6 px-6 text-center">
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight leading-tight max-w-xl mx-auto mb-3">
          Healthcare scheduling regulations, mapped
        </h1>
        <p className="text-neutral-500 text-sm max-w-md mx-auto leading-relaxed">
          Every regulation has a real legal citation. Click any jurisdiction to browse its rules.
          Runs entirely in the browser via WebAssembly.
        </p>
      </section>

      {loaded ? (
        <main className="flex-1 max-w-6xl mx-auto px-6 pb-16 w-full">
          {/* View toggle */}
          <div className="flex items-center justify-center gap-1 mb-4">
            <button
              onClick={() => switchView("us")}
              className={`px-4 py-1.5 rounded-l-lg text-sm font-medium transition-all ${view === "us" ? "bg-neutral-900 text-white" : "bg-neutral-100 text-neutral-500 hover:bg-neutral-200"}`}
            >
              United States
            </button>
            <button
              onClick={() => switchView("world")}
              className={`px-4 py-1.5 rounded-r-lg text-sm font-medium transition-all ${view === "world" ? "bg-neutral-900 text-white" : "bg-neutral-100 text-neutral-500 hover:bg-neutral-200"}`}
            >
              World
            </button>
          </div>

          {/* Map */}
          <section className="mb-4">
            {view === "us" ? (
              <USMap jurisdictions={jurisdictions} onSelect={setSelected} selected={selected} />
            ) : (
              <WorldMap jurisdictions={jurisdictions} onSelect={setSelected} selected={selected} />
            )}
          </section>

          {/* Coverage stats */}
          {viewStats.hospitals > 0 && (
            <div className="flex items-center justify-center gap-6 mb-6 text-xs text-neutral-400">
              <div>
                Covering regulations affecting
                <span className="font-mono font-semibold text-neutral-600 mx-1">{viewStats.hospitals.toLocaleString()}+</span>
                hospitals
              </div>
              <span className="text-neutral-200">|</span>
              <div>
                <span className="font-mono font-semibold text-neutral-600 mr-1">{(viewStats.workers / 1_000_000).toFixed(1)}M+</span>
                healthcare workers
              </div>
            </div>
          )}

          {/* Selected jurisdiction detail */}
          {selected && selectedInfo && (
            <section ref={detailRef} className="animate-in fade-in slide-in-from-bottom-2 duration-200 scroll-mt-16">
              <div className="flex items-center justify-between gap-3 mb-3 flex-wrap">
                <div className="flex items-center gap-3 flex-wrap">
                  <h2 className="text-lg font-bold">{selectedInfo.name}</h2>
                  <Badge variant="secondary" className="text-[10px]">{selectedInfo.code}</Badge>
                  <Badge variant="outline" className="text-[10px]">{filteredRules.length} rules</Badge>
                  {JURISDICTION_STATS[selectedInfo.code] && (
                    <span className="text-[11px] text-neutral-400">
                      {JURISDICTION_STATS[selectedInfo.code].hospitals.toLocaleString()} hospitals
                    </span>
                  )}
                </div>
                <button onClick={() => setSelected(null)} className="text-xs text-neutral-400 hover:text-neutral-700 transition-colors">
                  Close
                </button>
              </div>

              {/* Inherited rules toggle */}
              {selectedInfo.parent && (() => {
                const inheritedCount = selectedRules.length - localRuleKeys.size;
                return (
                  <div className="flex items-center gap-3 mb-3 px-3 py-2 bg-neutral-50 rounded-lg border border-neutral-100">
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() => setShowInherited(false)}
                        className={`px-3 py-1 rounded-md text-xs font-medium transition-all ${!showInherited ? "bg-neutral-900 text-white shadow-sm" : "text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100"}`}
                      >
                        {selectedInfo.code} only ({localRuleKeys.size})
                      </button>
                      <button
                        onClick={() => setShowInherited(true)}
                        className={`px-3 py-1 rounded-md text-xs font-medium transition-all ${showInherited ? "bg-neutral-900 text-white shadow-sm" : "text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100"}`}
                      >
                        All rules ({selectedRules.length})
                      </button>
                    </div>
                    {showInherited && inheritedCount > 0 && (
                      <span className="text-[11px] text-neutral-400">
                        {inheritedCount} inherited from {selectedInfo.parent}
                      </span>
                    )}
                  </div>
                );
              })()}

              <div className="border border-neutral-200 rounded-xl overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="text-[10px]">Rule</TableHead>
                      <TableHead className="text-[10px]">Value</TableHead>
                      <TableHead className="text-[10px] hidden sm:table-cell">Per</TableHead>
                      <TableHead className="text-[10px] hidden md:table-cell">Scope</TableHead>
                      <TableHead className="text-[10px]">Enforcement</TableHead>
                      <TableHead className="text-[10px]">Source</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filteredRules.map((r, i) => {
                      const v = r.values?.[0];
                      if (!v) return null;
                      const avg = v.averaged ? ` (avg ${v.averaged.count}${v.averaged.unit})` : "";
                      const cite = r.source.section ? `${r.source.title}, ${r.source.section}` : r.source.title;
                      const enf = r.enforcement === "mandatory" ? "destructive" as const : r.enforcement === "recommended" ? "secondary" as const : "outline" as const;
                      const isInherited = !localRuleKeys.has(r.key);
                      return (
                        <TableRow key={i}>
                          <TableCell>
                            <div className="flex items-center gap-1.5">
                              <span className="font-mono text-xs font-medium">{r.key}</span>
                              {isInherited && (
                                <span className="text-[9px] text-neutral-400 bg-neutral-100 px-1.5 py-0.5 rounded font-medium">inherited</span>
                              )}
                            </div>
                            <div className="text-[11px] text-neutral-500 mt-0.5">{r.name}</div>
                          </TableCell>
                          <TableCell className="font-mono text-xs">{op(r.operator)}{v.amount} {v.unit}</TableCell>
                          <TableCell className="font-mono text-[11px] text-neutral-500 hidden sm:table-cell">{v.per}{avg}</TableCell>
                          <TableCell className="hidden md:table-cell">{r.scope && <Badge variant="outline" className="text-[9px] font-mono">{r.scope}</Badge>}</TableCell>
                          <TableCell><Badge variant={enf} className="text-[10px]">{r.enforcement}</Badge></TableCell>
                          <TableCell className="text-[11px] max-w-[200px]">
                            {r.source.url ? (
                              <a href={r.source.url} target="_blank" rel="noopener noreferrer"
                                className="text-blue-600 hover:text-blue-800 hover:underline line-clamp-2" title={cite}>{cite}</a>
                            ) : (
                              <span className="text-neutral-500 line-clamp-2" title={cite}>{cite}</span>
                            )}
                          </TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              </div>
            </section>
          )}

          {/* Jurisdiction cards */}
          {view === "us" && (
            <section className="mt-10 mb-8">
              <h3 className="text-xs font-semibold text-neutral-400 uppercase tracking-wider mb-3">Other Jurisdictions</h3>
              <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-2">
                {jurisdictions.filter(j => !j.code.startsWith("US")).map(j => {
                  const stats = JURISDICTION_STATS[j.code];
                  return (
                    <button key={j.code} onClick={() => setSelected(selected === j.code ? null : j.code)}
                      className={`text-left px-3 py-2.5 rounded-lg border transition-all text-xs ${selected === j.code ? "border-neutral-900 bg-neutral-50 ring-1 ring-neutral-900" : "border-neutral-200 hover:border-neutral-300"}`}>
                      <div className="font-semibold text-neutral-900">{j.name}</div>
                      <div className="text-neutral-400 mt-0.5">{j.ruleCount} rules</div>
                      {stats && (
                        <div className="text-neutral-300 mt-0.5">{stats.hospitals.toLocaleString()} hospitals</div>
                      )}
                    </button>
                  );
                })}
              </div>
            </section>
          )}

          {view === "world" && (
            <section className="mt-10 mb-8">
              <h3 className="text-xs font-semibold text-neutral-400 uppercase tracking-wider mb-3">Sub-national Jurisdictions</h3>
              <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-2">
                {jurisdictions.filter(j => j.type === "state" || j.type === "region").map(j => {
                  const stats = JURISDICTION_STATS[j.code];
                  return (
                    <button key={j.code} onClick={() => setSelected(selected === j.code ? null : j.code)}
                      className={`text-left px-3 py-2.5 rounded-lg border transition-all text-xs ${selected === j.code ? "border-neutral-900 bg-neutral-50 ring-1 ring-neutral-900" : "border-neutral-200 hover:border-neutral-300"}`}>
                      <div className="font-semibold text-neutral-900">{j.name}</div>
                      <div className="text-neutral-400 mt-0.5">{j.code} | {j.ruleCount} rules</div>
                      {stats && (
                        <div className="text-neutral-300 mt-0.5">{stats.hospitals.toLocaleString()} hospitals</div>
                      )}
                    </button>
                  );
                })}
              </div>
            </section>
          )}
        </main>
      ) : (
        <div className="flex-1 flex items-center justify-center">
          <div className="text-neutral-400 text-sm animate-pulse">Loading regulation database...</div>
        </div>
      )}
    </div>
  );
}
