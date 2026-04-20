"use client";

import { useEffect, useRef, useState, useMemo, useCallback } from "react";
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

interface StatePath {
  fips: string;
  code: string | undefined;
  d: string;
}

export function USMap({ jurisdictions, onSelect, selected }: Props) {
  const svgRef = useRef<SVGSVGElement>(null);
  const tooltipRef = useRef<HTMLDivElement>(null);
  const [topology, setTopology] = useState<Topology | null>(null);
  const [hovered, setHovered] = useState<string | null>(null);

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

  // Pre-compute all path strings once when topology loads
  const { statePaths, borderPath } = useMemo((): { statePaths: StatePath[]; borderPath: string } => {
    if (!topology) return { statePaths: [], borderPath: "" };
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const states = topojson.feature(topology, topology.objects.states as any) as any;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const borders = topojson.mesh(topology, topology.objects.states as any, (a: any, b: any) => a !== b);
    const projection = geoAlbersUsa().fitSize([900, 500], states);
    const pathGen = geoPath(projection);

    return {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      statePaths: states.features.map((feature: any) => {
        const fips = String(feature.id).padStart(2, "0");
        return {
          fips,
          code: FIPS_TO_CODE[fips],
          d: pathGen(feature) || "",
        };
      }),
      borderPath: pathGen(borders) || "",
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

  if (!topology) return <div className="h-[500px] flex items-center justify-center text-neutral-400 text-sm">Loading map...</div>;

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
        {statePaths.map(({ fips, code, d }) => {
          const info = code ? covered.get(code) : undefined;
          const isCovered = !!info;
          const isHovered = hovered === fips;
          const isSelected = selected === code;

          return (
            <path
              key={fips}
              d={d}
              fill={isSelected ? "#059669" : isCovered ? "#bbf7d0" : "#f0f0f0"}
              stroke={isSelected ? "#047857" : isHovered ? "#16a34a" : isCovered ? "#86efac" : "#d4d4d4"}
              strokeWidth={isSelected ? 2 : isHovered ? 1.5 : 0.5}
              cursor="pointer"
              onMouseEnter={(e) => {
                setHovered(fips);
                updateTooltip(e);
              }}
              onMouseMove={updateTooltip}
              onMouseLeave={() => setHovered(null)}
              onClick={() => onSelect(isSelected ? null : code || null)}
              style={{ transition: "fill 0.15s, stroke 0.15s" }}
            />
          );
        })}

        {/* State borders */}
        <path
          d={borderPath}
          fill="none"
          stroke="#d4d4d4"
          strokeWidth={0.5}
          strokeLinejoin="round"
          pointerEvents="none"
        />
      </svg>

      {/* Tooltip */}
      <div
        ref={tooltipRef}
        className={`absolute pointer-events-none z-20 bg-neutral-900 text-white px-3 py-2 rounded-lg text-xs shadow-xl transition-opacity duration-100 ${hovered && hoveredName ? "opacity-100" : "opacity-0"}`}
        style={{ transform: "translate(-50%, -100%)" }}
      >
        {hoveredName && (
          <>
            <div className="font-semibold text-[13px]">{hoveredName}</div>
            {hoveredInfo ? (
              <div className="text-emerald-400 mt-0.5">{hoveredInfo.ruleCount} state-specific rules</div>
            ) : (
              <div className="text-neutral-400 mt-0.5">Federal rules only ({federal?.ruleCount || 0} rules)</div>
            )}
          </>
        )}
      </div>

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
