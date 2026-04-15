import type { Scenario } from "./types";

const C = ["#818cf8", "#f472b6", "#38bdf8", "#34d399", "#fb923c", "#a78bfa"];

function d(weekOffset: number, dayOfWeek: number): string {
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
      // Santos: 6 days/week x 14h = 84h/week (exceeds 80h ACGME)
      ...[0,1,2,3].flatMap(w => [0,1,2,3,4,5].map(day => ({ staff_id: "dr-santos", staff_type: "resident", start: d(w,day)+"T06:00:00", end: d(w,day)+"T20:00:00" }))),
      // Park: normal 5 days/week x 12h
      ...[0,1,2,3].flatMap(w => [0,1,2,3,4].map(day => ({ staff_id: "dr-park", staff_type: "resident", start: d(w,day)+"T07:00:00", end: d(w,day)+"T19:00:00" }))),
      // Lopez: 3 days/week
      ...[0,1,2,3].flatMap(w => [0,2,4].map(day => ({ staff_id: "dr-lopez", staff_type: "resident", start: d(w,day)+"T08:00:00", end: d(w,day)+"T20:00:00" }))),
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
      // Elena: weeks of shifts, but week 1 has a back-to-back violation
      ...[0,1,2,3].flatMap(w => {
        if (w === 0) return [
          { staff_id: "elena", staff_type: "statutory-personnel", start: d(0,0)+"T08:00:00", end: d(0,0)+"T20:00:00" },
          { staff_id: "elena", staff_type: "statutory-personnel", start: d(0,0)+"T23:00:00", end: d(0,1)+"T11:00:00" }, // 3h gap!
          { staff_id: "elena", staff_type: "statutory-personnel", start: d(0,3)+"T08:00:00", end: d(0,3)+"T20:00:00" },
        ];
        return [0,2,4].map(day => ({ staff_id: "elena", staff_type: "statutory-personnel", start: d(w,day)+"T08:00:00", end: d(w,day)+"T20:00:00" }));
      }),
      // Maria: normal rotation
      ...[0,1,2,3].flatMap(w => [0,2,4].map(day => ({ staff_id: "maria-g", staff_type: "statutory-personnel", start: d(w,day)+"T08:00:00", end: d(w,day)+"T20:00:00" }))),
      // Pablo: normal rotation
      ...[0,1,2,3].flatMap(w => [1,3].map(day => ({ staff_id: "pablo-f", staff_type: "statutory-personnel", start: d(w,day)+"T08:00:00", end: d(w,day)+"T20:00:00" }))),
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
      // Vega: regular shifts + 5 guards (exceeds 4/month)
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: d(0,0)+"T08:00:00", end: d(0,1)+"T08:00:00", on_call: true },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: d(0,2)+"T08:00:00", end: d(0,2)+"T16:00:00" },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: d(0,4)+"T08:00:00", end: d(0,5)+"T08:00:00", on_call: true },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: d(1,1)+"T08:00:00", end: d(1,1)+"T16:00:00" },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: d(1,3)+"T08:00:00", end: d(1,4)+"T08:00:00", on_call: true },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: d(2,0)+"T08:00:00", end: d(2,1)+"T08:00:00", on_call: true },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: d(2,2)+"T08:00:00", end: d(2,2)+"T16:00:00" },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: d(2,4)+"T08:00:00", end: d(2,5)+"T08:00:00", on_call: true },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: d(3,1)+"T08:00:00", end: d(3,1)+"T16:00:00" },
      { staff_id: "dr-vega", staff_type: "statutory-personnel", start: d(3,3)+"T08:00:00", end: d(3,3)+"T16:00:00" },
      // Ruiz: 3 guards (within limit) + regular shifts
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: d(0,1)+"T08:00:00", end: d(0,2)+"T08:00:00", on_call: true },
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: d(0,3)+"T08:00:00", end: d(0,3)+"T16:00:00" },
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: d(1,0)+"T08:00:00", end: d(1,1)+"T08:00:00", on_call: true },
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: d(1,2)+"T08:00:00", end: d(1,2)+"T16:00:00" },
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: d(2,1)+"T08:00:00", end: d(2,2)+"T08:00:00", on_call: true },
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: d(2,3)+"T08:00:00", end: d(2,3)+"T16:00:00" },
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: d(3,0)+"T08:00:00", end: d(3,0)+"T16:00:00" },
      { staff_id: "dr-ruiz", staff_type: "statutory-personnel", start: d(3,2)+"T08:00:00", end: d(3,2)+"T16:00:00" },
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
      // 4 weeks of compliant 3x12 rotation
      ...[0,1,2,3].flatMap(w => [
        { staff_id: "sarah", staff_type: "nurse-rn", unit_type: "icu", start: d(w,0)+"T07:00:00", end: d(w,0)+"T19:00:00" },
        { staff_id: "sarah", staff_type: "nurse-rn", unit_type: "icu", start: d(w,2)+"T07:00:00", end: d(w,2)+"T19:00:00" },
        { staff_id: "sarah", staff_type: "nurse-rn", unit_type: "icu", start: d(w,4)+"T07:00:00", end: d(w,4)+"T19:00:00" },
        { staff_id: "lisa", staff_type: "nurse-rn", unit_type: "icu", start: d(w,1)+"T07:00:00", end: d(w,1)+"T19:00:00" },
        { staff_id: "lisa", staff_type: "nurse-rn", unit_type: "icu", start: d(w,3)+"T07:00:00", end: d(w,3)+"T19:00:00" },
        { staff_id: "lisa", staff_type: "nurse-rn", unit_type: "icu", start: d(w,5)+"T07:00:00", end: d(w,5)+"T19:00:00" },
        { staff_id: "mark", staff_type: "nurse-rn", unit_type: "icu", start: d(w,0)+"T19:00:00", end: d(w,1)+"T07:00:00" },
        { staff_id: "mark", staff_type: "nurse-rn", unit_type: "icu", start: d(w,2)+"T19:00:00", end: d(w,3)+"T07:00:00" },
        { staff_id: "mark", staff_type: "nurse-rn", unit_type: "icu", start: d(w,4)+"T19:00:00", end: d(w,5)+"T07:00:00" },
      ]),
    ],
  },
];
