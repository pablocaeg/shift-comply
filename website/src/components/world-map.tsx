"use client";

import { useEffect, useRef, useState, useMemo, useCallback } from "react";
import { geoNaturalEarth1, geoPath, geoGraticule } from "d3-geo";
import * as topojson from "topojson-client";
import { COUNTRY_NUMERIC_TO_CODE, EU_MEMBERS, type JurisdictionInfo } from "@/lib/jurisdiction-data";

interface Props {
  jurisdictions: JurisdictionInfo[];
  onSelect: (code: string | null) => void;
  selected: string | null;
}

interface CountryPath {
  numericId: string;
  code: string | undefined;
  d: string;
}

export function WorldMap({ jurisdictions, onSelect, selected }: Props) {
  const svgRef = useRef<SVGSVGElement>(null);
  const tooltipRef = useRef<HTMLDivElement>(null);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [topology, setTopology] = useState<any>(null);
  const [hovered, setHovered] = useState<string | null>(null);

  const covered = useMemo(() => {
    const map = new Map<string, JurisdictionInfo>();
    for (const j of jurisdictions) map.set(j.code, j);
    return map;
  }, [jurisdictions]);

  const hasEU = covered.has("EU");

  useEffect(() => {
    fetch(`${process.env.NODE_ENV === "production" ? "/shift-comply" : ""}/world-110m.json`)
      .then(r => r.json())
      .then(setTopology);
  }, []);

  // Pre-compute all path strings once when topology loads
  const { countryPaths, borderPath, graticulePath } = useMemo((): { countryPaths: CountryPath[]; borderPath: string; graticulePath: string } => {
    if (!topology) return { countryPaths: [], borderPath: "", graticulePath: "" };
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const countries = topojson.feature(topology, topology.objects.countries as any) as any;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const borders = topojson.mesh(topology, topology.objects.countries as any, (a: any, b: any) => a !== b);
    const projection = geoNaturalEarth1().fitSize([900, 440], countries);
    const pathGen = geoPath(projection);
    const graticule = geoGraticule();

    return {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      countryPaths: countries.features.map((feature: any, idx: number) => {
        const numericId = feature.id != null ? String(feature.id) : `unknown-${idx}`;
        return {
          numericId,
          code: COUNTRY_NUMERIC_TO_CODE[numericId],
          d: pathGen(feature) || "",
        };
      }),
      borderPath: pathGen(borders) || "",
      graticulePath: pathGen(graticule()) || "",
    };
  }, [topology]);

  // Update tooltip position via DOM ref (no re-render)
  const updateTooltip = useCallback((e: React.MouseEvent) => {
    const rect = svgRef.current?.getBoundingClientRect();
    if (rect && tooltipRef.current) {
      tooltipRef.current.style.left = `${e.clientX - rect.left}px`;
      tooltipRef.current.style.top = `${e.clientY - rect.top - 10}px`;
    }
  }, []);

  if (!topology) return <div className="h-[400px] flex items-center justify-center text-neutral-400 text-sm">Loading map...</div>;

  const hoveredCode = hovered ? COUNTRY_NUMERIC_TO_CODE[hovered] : null;
  const hoveredInfo = hoveredCode ? covered.get(hoveredCode) : null;
  const isHoveredEU = hovered ? EU_MEMBERS.has(String(hovered).padStart(3, "0")) && hasEU : false;

  return (
    <div className="relative">
      <svg ref={svgRef} viewBox="0 0 900 440" className="w-full h-auto" style={{ background: "#f5f5f4", borderRadius: 12 }}>
        {/* Graticule */}
        <path d={graticulePath} fill="none" stroke="#f0f0f0" strokeWidth={0.3} />

        {/* Countries */}
        {countryPaths.map(({ numericId, code, d }) => {
          const info = code ? covered.get(code) : undefined;
          const isEUMember = EU_MEMBERS.has(numericId.padStart(3, "0"));
          const isCovered = !!info;
          const isEUOnly = !isCovered && isEUMember && hasEU;
          const isClickable = isCovered || isEUOnly;
          const isHovered = hovered === numericId;
          const isSelected = isClickable && selected === code;

          let fill = "#f0f0f0";
          let stroke = "#d4d4d4";
          if (isSelected) { fill = "#059669"; stroke = "#047857"; }
          else if (isCovered) { fill = "#bbf7d0"; stroke = "#86efac"; }
          else if (isEUOnly) { fill = "#e0f2fe"; stroke = "#93c5fd"; }
          if (isHovered && !isSelected) { stroke = isCovered ? "#16a34a" : isEUOnly ? "#3b82f6" : "#a3a3a3"; }

          return (
            <path
              key={numericId}
              d={d}
              fill={fill}
              stroke={stroke}
              strokeWidth={isSelected ? 1.5 : isHovered ? 1 : 0.3}
              cursor={isClickable ? "pointer" : "default"}
              onMouseEnter={(e) => {
                setHovered(numericId);
                updateTooltip(e);
              }}
              onMouseMove={updateTooltip}
              onMouseLeave={() => setHovered(null)}
              onClick={() => {
                if (isSelected) onSelect(null);
                else if (isCovered && code) onSelect(code);
                else if (isEUOnly) onSelect("EU");
              }}
              style={{ transition: "fill 0.15s" }}
            />
          );
        })}

        {/* Borders */}
        <path d={borderPath} fill="none" stroke="rgba(255,255,255,0.85)" strokeWidth={0.6} pointerEvents="none" />
      </svg>

      {/* Tooltip */}
      <div
        ref={tooltipRef}
        className={`absolute pointer-events-none z-20 bg-neutral-900 text-white px-3 py-2 rounded-lg text-xs shadow-xl transition-opacity duration-100 ${hovered && (hoveredInfo || isHoveredEU) ? "opacity-100" : "opacity-0"}`}
        style={{ transform: "translate(-50%, -100%)" }}
      >
        {hoveredInfo ? (
          <>
            <div className="font-semibold text-[13px]">{hoveredInfo.name}</div>
            <div className="text-emerald-400 mt-0.5">{hoveredInfo.ruleCount} rules</div>
          </>
        ) : isHoveredEU ? (
          <>
            <div className="font-semibold text-[13px]">EU Member State</div>
            <div className="text-blue-400 mt-0.5">Inherits EU Working Time Directive ({covered.get("EU")?.ruleCount || 0} rules)</div>
          </>
        ) : null}
      </div>

      {/* Legend */}
      <div className="flex items-center gap-5 mt-3 px-1 flex-wrap">
        <div className="flex items-center gap-1.5">
          <div className="w-3.5 h-3.5 rounded bg-[#bbf7d0] border border-[#86efac]" />
          <span className="text-[11px] text-neutral-500">Country-specific rules</span>
        </div>
        <div className="flex items-center gap-1.5">
          <div className="w-3.5 h-3.5 rounded bg-[#e0f2fe] border border-[#93c5fd]" />
          <span className="text-[11px] text-neutral-500">EU directive (inherited)</span>
        </div>
        <div className="flex items-center gap-1.5">
          <div className="w-3.5 h-3.5 rounded bg-[#f0f0f0] border border-[#d4d4d4]" />
          <span className="text-[11px] text-neutral-400">Not covered</span>
        </div>
      </div>
    </div>
  );
}
