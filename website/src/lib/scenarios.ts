import type { Scenario } from "./types";

const C = ["#818cf8", "#f472b6", "#38bdf8", "#34d399", "#fb923c", "#a78bfa"];

// Helper: generate dates for a 4-week period starting March 3, 2025
function d(weekOffset: number, dayOfWeek: number): string {
  // week 0 = March 3-9, day 0 = Monday
  const base = new Date(2025, 2, 3); // March 3 2025 (Monday)
  base.setDate(base.getDate() + weekOffset * 7 + dayOfWeek);
  return `${base.getFullYear()}-${String(base.getMonth() + 1).padStart(2, "0")}-${String(base.getDate()).padStart(2, "0")}`;
}

export const SCENARIOS: Scenario[] = [
  {
    id: "overworked", badge: "fail", label: "Violations",
    name: "The Overworked Resident", who: "University Hospital, California",
    info: "Dr. Santos averages 84 hours/week over 4 weeks. ACGME limits residents to 80.",
    jurisdiction: "US-CA", scope: "accredited_programs", flagged: "dr-santos",
    workers: [
      { id: "dr-santos", name: "Dr. Maria Santos", role: "PGY-2 Resident", type: "resident", color: C[0] },
      { id: "dr-park", name: "Dr. James Park", role: "PGY-3 Resident", type: "resident", color: C[1] },
      { id: "dr-lopez", name: "Dr. Ana Lopez", role: "PGY-2 Resident", type: "resident", color: C[2] },
    ],
    shifts: [
      // Dr. Santos: 6 days/week x 14h = 84h/week for 4 weeks (exceeds ACGME 80h avg)
      ...[0, 1, 2, 3].flatMap(w => [0, 1, 2, 3, 4, 5].map(day => ({
        staff_id: "dr-santos", staff_type: "resident",
        start: d(w, day) + "T06:00:00", end: d(w, day) + "T20:00:00",
      }))),
      // Dr. Park: normal schedule 5 days/week x 12h = 60h/week
      ...[0, 1, 2, 3].flatMap(w => [0, 1, 2, 3, 4].map(day => ({
        staff_id: "dr-park", staff_type: "resident",
        start: d(w, day) + "T07:00:00", end: d(w, day) + "T19:00:00",
      }))),
      // Dr. Lopez: light schedule 3 days/week
      ...[0, 1, 2, 3].flatMap(w => [0, 2, 4].map(day => ({
        staff_id: "dr-lopez", staff_type: "resident",
        start: d(w, day) + "T08:00:00", end: d(w, day) + "T20:00:00",
      }))),
    ],
  },
  {
    id: "backtoback", badge: "fail", label: "Violations",
    name: "The Back-to-Back Nurse", who: "Hospital Publico, Spain",
    info: "Elena has only 3 hours between shifts. Spain requires 12 hours minimum rest.",
    jurisdiction: "ES", scope: "public_health", flagged: "elena",
    workers: [
      { id: "elena", name: "Elena Rodriguez", role: "ICU Nurse", type: "statutory-personnel", color: C[0] },
      { id: "maria-g", name: "Maria Garcia", role: "ICU Nurse", type: "statutory-personnel", color: C[1] },
      { id: "pablo-f", name: "Pablo Fernandez", role: "ICU Nurse", type: "statutory-personnel", color: C[2] },
    ],
    shifts: [
      // Elena: back-to-back with only 3h rest
      { staff_id: "elena", staff_type: "statutory-personnel", start: "2025-03-10T08:00:00", end: "2025-03-10T20:00:00" },
      { staff_id: "elena", staff_type: "statutory-personnel", start: "2025-03-10T23:00:00", end: "2025-03-11T11:00:00" },
      { staff_id: "elena", staff_type: "statutory-personnel", start: "2025-03-13T08:00:00", end: "2025-03-13T20:00:00" },
      // Maria: normal
      { staff_id: "maria-g", staff_type: "statutory-personnel", start: "2025-03-10T08:00:00", end: "2025-03-10T20:00:00" },
      { staff_id: "maria-g", staff_type: "statutory-personnel", start: "2025-03-12T08:00:00", end: "2025-03-12T20:00:00" },
      { staff_id: "maria-g", staff_type: "statutory-personnel", start: "2025-03-14T08:00:00", end: "2025-03-14T20:00:00" },
      // Pablo: normal
      { staff_id: "pablo-f", staff_type: "statutory-personnel", start: "2025-03-11T08:00:00", end: "2025-03-11T20:00:00" },
      { staff_id: "pablo-f", staff_type: "statutory-personnel", start: "2025-03-13T08:00:00", end: "2025-03-13T20:00:00" },
      { staff_id: "pablo-f", staff_type: "statutory-personnel", start: "2025-03-15T08:00:00", end: "2025-03-15T20:00:00" },
    ],
  },
  {
    id: "exhausted", badge: "fail", label: "Violations",
    name: "The Exhausted MIR", who: "Hospital ICS, Catalonia",
    info: "Dr. Vega has 5 on-call guards in March. The ICS agreement limits guards to 4 per month.",
    jurisdiction: "ES-CT", scope: "public_health", flagged: "dr-vega",
    workers: [
      { id: "dr-vega", name: "Dr. Carlos Vega", role: "MIR Resident", type: "statutory-personnel", color: C[0] },
      { id: "dr-ruiz", name: "Dr. Laura Ruiz", role: "MIR Resident", type: "statutory-personnel", color: C[1] },
    ],
    shifts: [
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: "2025-03-10T08:00:00", end: "2025-03-11T08:00:00", on_call: true },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: "2025-03-12T08:00:00", end: "2025-03-12T16:00:00" },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: "2025-03-13T08:00:00", end: "2025-03-14T08:00:00", on_call: true },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: "2025-03-15T08:00:00", end: "2025-03-16T08:00:00", on_call: true },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: "2025-03-17T08:00:00", end: "2025-03-18T08:00:00", on_call: true },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: "2025-03-20T08:00:00", end: "2025-03-21T08:00:00", on_call: true },
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: "2025-03-11T08:00:00", end: "2025-03-12T08:00:00", on_call: true },
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: "2025-03-14T08:00:00", end: "2025-03-15T08:00:00", on_call: true },
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: "2025-03-18T08:00:00", end: "2025-03-19T08:00:00", on_call: true },
    ],
  },
  {
    id: "compliant", badge: "pass", label: "Compliant",
    name: "The Compliant Week", who: "Memorial Hospital, California",
    info: "Three ICU nurses with 12-hour shifts and proper rest between rotations.",
    jurisdiction: "US-CA", scope: "hospitals", flagged: null,
    workers: [
      { id: "sarah", name: "Sarah Chen", role: "RN, ICU", type: "nurse-rn", color: C[3] },
      { id: "lisa", name: "Lisa Wang", role: "RN, ICU", type: "nurse-rn", color: C[4] },
      { id: "mark", name: "Mark Davis", role: "RN, ICU", type: "nurse-rn", color: C[5] },
    ],
    shifts: [
      { staff_id: "sarah", staff_type: "nurse-rn", unit_type: "icu", start: "2025-03-10T07:00:00", end: "2025-03-10T19:00:00" },
      { staff_id: "sarah", staff_type: "nurse-rn", unit_type: "icu", start: "2025-03-12T07:00:00", end: "2025-03-12T19:00:00" },
      { staff_id: "sarah", staff_type: "nurse-rn", unit_type: "icu", start: "2025-03-14T07:00:00", end: "2025-03-14T19:00:00" },
      { staff_id: "lisa", staff_type: "nurse-rn", unit_type: "icu", start: "2025-03-11T07:00:00", end: "2025-03-11T19:00:00" },
      { staff_id: "lisa", staff_type: "nurse-rn", unit_type: "icu", start: "2025-03-13T07:00:00", end: "2025-03-13T19:00:00" },
      { staff_id: "lisa", staff_type: "nurse-rn", unit_type: "icu", start: "2025-03-15T07:00:00", end: "2025-03-15T19:00:00" },
      { staff_id: "mark", staff_type: "nurse-rn", unit_type: "icu", start: "2025-03-10T19:00:00", end: "2025-03-11T07:00:00" },
      { staff_id: "mark", staff_type: "nurse-rn", unit_type: "icu", start: "2025-03-12T19:00:00", end: "2025-03-13T07:00:00" },
      { staff_id: "mark", staff_type: "nurse-rn", unit_type: "icu", start: "2025-03-14T19:00:00", end: "2025-03-15T07:00:00" },
    ],
  },
];
