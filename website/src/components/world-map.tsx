"use client";

import { useEffect, useRef, useState, useMemo } from "react";
import { geoNaturalEarth1, geoPath, geoGraticule } from "d3-geo";
import * as topojson from "topojson-client";
import { COUNTRY_NUMERIC_TO_CODE, EU_MEMBERS, type JurisdictionInfo } from "@/lib/jurisdiction-data";

interface Props {
  jurisdictions: JurisdictionInfo[];
  onSelect: (code: string | null) => void;
  selected: string | null;
}

export function WorldMap({ jurisdictions, onSelect, selected }: Props) {
  const svgRef = useRef<SVGSVGElement>(null);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [topology, setTopology] = useState<any>(null);
  const [hovered, setHovered] = useState<string | null>(null);
  const [tooltipPos, setTooltipPos] = useState({ x: 0, y: 0 });

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

  if (!topology) return <div className="h-[400px] flex items-center justify-center text-neutral-400 text-sm">Loading map...</div>;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const countries = topojson.feature(topology, topology.objects.countries as any) as any;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const borders = topojson.mesh(topology, topology.objects.countries as any, (a: any, b: any) => a !== b);

  const projection = geoNaturalEarth1().fitSize([900, 440], countries);
  const path = geoPath(projection);
  const graticule = geoGraticule();

  const hoveredNumeric = hovered;
  const hoveredCode = hoveredNumeric ? COUNTRY_NUMERIC_TO_CODE[hoveredNumeric] : null;
  const hoveredInfo = hoveredCode ? covered.get(hoveredCode) : null;

  return (
    <div className="relative">
      <svg ref={svgRef} viewBox="0 0 900 440" className="w-full h-auto" style={{ background: "#fafafa", borderRadius: 12 }}>
        {/* Graticule */}
        <path d={path(graticule()) || ""} fill="none" stroke="#f0f0f0" strokeWidth={0.3} />

        {/* Countries */}
        {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
        {countries.features.map((feature: any) => {
          const numericId = String(feature.id);
          const code = COUNTRY_NUMERIC_TO_CODE[numericId];
          const info = code ? covered.get(code) : undefined;
          const isEUMember = EU_MEMBERS.has(numericId.padStart(3, "0"));
          const isCovered = !!info;
          const isEUOnly = !isCovered && isEUMember && hasEU;
          const isHovered = hovered === numericId;
          const isSelected = selected === code;

          let fill = "#f0f0f0";
          let stroke = "#d4d4d4";
          if (isSelected) { fill = "#059669"; stroke = "#047857"; }
          else if (isCovered) { fill = "#bbf7d0"; stroke = "#86efac"; }
          else if (isEUOnly) { fill = "#e0f2fe"; stroke = "#93c5fd"; }
          if (isHovered && !isSelected) { stroke = isCovered ? "#16a34a" : isEUOnly ? "#3b82f6" : "#a3a3a3"; }

          return (
            <path
              key={numericId}
              d={path(feature) || ""}
              fill={fill}
              stroke={stroke}
              strokeWidth={isSelected ? 1.5 : isHovered ? 1 : 0.3}
              cursor={code || isEUOnly ? "pointer" : "default"}
              onMouseEnter={(e) => {
                setHovered(numericId);
                const rect = svgRef.current?.getBoundingClientRect();
                if (rect) setTooltipPos({ x: e.clientX - rect.left, y: e.clientY - rect.top - 10 });
              }}
              onMouseMove={(e) => {
                const rect = svgRef.current?.getBoundingClientRect();
                if (rect) setTooltipPos({ x: e.clientX - rect.left, y: e.clientY - rect.top - 10 });
              }}
              onMouseLeave={() => setHovered(null)}
              onClick={() => {
                if (isSelected) onSelect(null);
                else if (code) onSelect(code);
                else if (isEUOnly) onSelect("EU");
              }}
              style={{ transition: "fill 0.15s" }}
            />
          );
        })}

        {/* Borders */}
        <path d={path(borders) || ""} fill="none" stroke="white" strokeWidth={0.3} pointerEvents="none" />
      </svg>

      {/* Tooltip */}
      {hovered && (
        <div className="absolute pointer-events-none z-20 bg-neutral-900 text-white px-3 py-2 rounded-lg text-xs shadow-xl"
          style={{ left: tooltipPos.x, top: tooltipPos.y, transform: "translate(-50%, -100%)" }}>
          {hoveredInfo ? (
            <>
              <div className="font-semibold text-[13px]">{hoveredInfo.name}</div>
              <div className="text-emerald-400 mt-0.5">{hoveredInfo.ruleCount} rules</div>
            </>
          ) : EU_MEMBERS.has(String(hovered).padStart(3, "0")) && hasEU ? (
            <>
              <div className="font-semibold text-[13px]">EU Member State</div>
              <div className="text-blue-400 mt-0.5">Inherits EU Working Time Directive ({covered.get("EU")?.ruleCount || 0} rules)</div>
            </>
          ) : null}
        </div>
      )}

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
