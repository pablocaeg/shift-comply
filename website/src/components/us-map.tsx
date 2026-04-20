"use client";

import { useEffect, useRef, useState, useMemo } from "react";
import { geoAlbersUsa, geoPath } from "d3-geo";
import * as topojson from "topojson-client";
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Topology = any;
import { FIPS_TO_CODE, STATE_NAMES, type JurisdictionInfo } from "@/lib/jurisdiction-data";

interface Props {
  jurisdictions: JurisdictionInfo[];
  onSelect: (code: string | null) => void;
  selected: string | null;
}

export function USMap({ jurisdictions, onSelect, selected }: Props) {
  const svgRef = useRef<SVGSVGElement>(null);
  const [topology, setTopology] = useState<Topology | null>(null);
  const [hovered, setHovered] = useState<string | null>(null);
  const [tooltipPos, setTooltipPos] = useState({ x: 0, y: 0 });

  const covered = useMemo(() => {
    const map = new Map<string, JurisdictionInfo>();
    for (const j of jurisdictions) {
      if (j.code.startsWith("US-")) map.set(j.code, j);
    }
    return map;
  }, [jurisdictions]);

  const federal = useMemo(() => jurisdictions.find(j => j.code === "US"), [jurisdictions]);

  useEffect(() => {
    fetch(`${process.env.NODE_ENV === "production" ? "/shift-comply" : ""}/us-states-10m.json`)
      .then(r => r.json())
      .then(setTopology);
  }, []);

  if (!topology) return <div className="h-[500px] flex items-center justify-center text-neutral-400 text-sm">Loading map...</div>;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const states = topojson.feature(topology, topology.objects.states as any) as any;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const borders = topojson.mesh(topology, topology.objects.states as any, (a: any, b: any) => a !== b);
  const projection = geoAlbersUsa().fitSize([900, 500], states);
  const path = geoPath(projection);

  const hoveredCode = hovered ? FIPS_TO_CODE[hovered] : null;
  const hoveredInfo = hoveredCode ? covered.get(hoveredCode) : null;
  const hoveredName = hoveredCode ? STATE_NAMES[hoveredCode] : null;

  return (
    <div className="relative">
      <svg
        ref={svgRef}
        viewBox="0 0 900 500"
        className="w-full h-auto"
        style={{ background: "#fafafa", borderRadius: 12 }}
      >
        {/* States */}
        {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
        {states.features.map((feature: any) => {
          const fips = String(feature.id).padStart(2, "0");
          const code = FIPS_TO_CODE[fips];
          const info = code ? covered.get(code) : undefined;
          const isCovered = !!info;
          const isHovered = hovered === fips;
          const isSelected = selected === code;

          return (
            <path
              key={fips}
              d={path(feature) || ""}
              fill={isSelected ? "#059669" : isCovered ? "#bbf7d0" : "#f0f0f0"}
              stroke={isSelected ? "#047857" : isHovered ? "#16a34a" : isCovered ? "#86efac" : "#d4d4d4"}
              strokeWidth={isSelected ? 2 : isHovered ? 1.5 : 0.5}
              cursor="pointer"
              onMouseEnter={(e) => {
                setHovered(fips);
                const rect = svgRef.current?.getBoundingClientRect();
                if (rect) {
                  setTooltipPos({
                    x: e.clientX - rect.left,
                    y: e.clientY - rect.top - 10,
                  });
                }
              }}
              onMouseMove={(e) => {
                const rect = svgRef.current?.getBoundingClientRect();
                if (rect) {
                  setTooltipPos({
                    x: e.clientX - rect.left,
                    y: e.clientY - rect.top - 10,
                  });
                }
              }}
              onMouseLeave={() => setHovered(null)}
              onClick={() => onSelect(isSelected ? null : code || null)}
              style={{ transition: "fill 0.15s, stroke 0.15s" }}
            />
          );
        })}

        {/* State borders */}
        <path
          d={path(borders) || ""}
          fill="none"
          stroke="white"
          strokeWidth={0.5}
          strokeLinejoin="round"
          pointerEvents="none"
        />
      </svg>

      {/* Tooltip */}
      {hovered && hoveredName && (
        <div
          className="absolute pointer-events-none z-20 bg-neutral-900 text-white px-3 py-2 rounded-lg text-xs shadow-xl"
          style={{
            left: tooltipPos.x,
            top: tooltipPos.y,
            transform: "translate(-50%, -100%)",
          }}
        >
          <div className="font-semibold text-[13px]">{hoveredName}</div>
          {hoveredInfo ? (
            <div className="text-emerald-400 mt-0.5">{hoveredInfo.ruleCount} state-specific rules</div>
          ) : (
            <div className="text-neutral-400 mt-0.5">Federal rules only ({federal?.ruleCount || 0} rules)</div>
          )}
        </div>
      )}

      {/* Legend */}
      <div className="flex items-center gap-5 mt-3 px-1">
        <div className="flex items-center gap-1.5">
          <div className="w-3.5 h-3.5 rounded bg-[#bbf7d0] border border-[#86efac]" />
          <span className="text-[11px] text-neutral-500">State-specific rules ({covered.size} states)</span>
        </div>
        <div className="flex items-center gap-1.5">
          <div className="w-3.5 h-3.5 rounded bg-[#f0f0f0] border border-[#d4d4d4]" />
          <span className="text-[11px] text-neutral-400">Federal rules only</span>
        </div>
        {federal && (
          <span className="text-[11px] text-neutral-400 ml-auto">{federal.ruleCount} federal rules inherited by all states</span>
        )}
      </div>
    </div>
  );
}
