import { describe, it, expect } from "vitest";
import { applyFix, type FixResult } from "./autofix";
import type { Shift, Violation, Worker } from "./types";

function makeShift(id: string, staffId: string, staffType: string, start: string, end: string, onCall = false): Shift {
  return { _uid: id, staff_id: staffId, staff_type: staffType, start, end, on_call: onCall };
}

const workers: Worker[] = [
  { id: "w1", name: "Worker One", role: "Resident", type: "resident", color: "#818cf8" },
  { id: "w2", name: "Worker Two", role: "Resident", type: "resident", color: "#f472b6" },
];

describe("applyFix", () => {
  it("shortens shifts for weekly hours violation", () => {
    const shifts: Shift[] = [
      makeShift("s1", "w1", "resident", "2025-03-10T06:00:00", "2025-03-10T20:00:00"), // 14h
      makeShift("s2", "w1", "resident", "2025-03-11T06:00:00", "2025-03-11T20:00:00"), // 14h
    ];
    const violation: Violation = {
      rule_key: "max-weekly-hours", rule_name: "Max Weekly Hours",
      severity: "mandatory", staff_id: "w1", message: "", citation: "",
      actual: 84, limit: 80,
    };
    const result = applyFix(shifts, violation, workers);
    // Should have shortened some shifts, not deleted them
    expect(result.shifts.length).toBe(2);
    expect(result.description).toContain("Reduced");
  });

  it("reassigns guards for guard limit violation", () => {
    const shifts: Shift[] = [
      makeShift("g1", "w1", "resident", "2025-03-10T08:00:00", "2025-03-11T08:00:00", true),
      makeShift("g2", "w1", "resident", "2025-03-13T08:00:00", "2025-03-14T08:00:00", true),
      makeShift("g3", "w1", "resident", "2025-03-17T08:00:00", "2025-03-18T08:00:00", true),
      makeShift("g4", "w2", "resident", "2025-03-11T08:00:00", "2025-03-12T08:00:00", true),
    ];
    const violation: Violation = {
      rule_key: "max-guards-monthly", rule_name: "Max Guards",
      severity: "mandatory", staff_id: "w1", message: "", citation: "",
      actual: 3, limit: 2,
    };
    const result = applyFix(shifts, violation, workers);
    // Should have reassigned 1 guard from w1 to w2
    const w1Guards = result.shifts.filter(s => s.staff_id === "w1" && s.on_call).length;
    const w2Guards = result.shifts.filter(s => s.staff_id === "w2" && s.on_call).length;
    expect(w1Guards).toBe(2);
    expect(w2Guards).toBe(2);
    expect(result.description).toContain("Reassigned");
  });

  it("pushes shifts forward for rest gap violation", () => {
    const shifts: Shift[] = [
      makeShift("s1", "w1", "resident", "2025-03-10T08:00:00", "2025-03-10T20:00:00"),
      makeShift("s2", "w1", "resident", "2025-03-10T23:00:00", "2025-03-11T11:00:00"), // 3h gap
    ];
    const violation: Violation = {
      rule_key: "min-rest-between-shifts", rule_name: "Min Rest",
      severity: "mandatory", staff_id: "w1", message: "", citation: "",
      actual: 3, limit: 12,
    };
    const result = applyFix(shifts, violation, workers);
    expect(result.shifts.length).toBe(2);
    // Second shift should have been pushed forward
    const s2 = result.shifts.find(s => s._uid === "s2")!;
    expect(s2.start).not.toBe("2025-03-10T23:00:00");
    expect(result.description).toContain("rest");
  });

  it("trims shifts for max shift hours violation", () => {
    const shifts: Shift[] = [
      makeShift("s1", "w1", "resident", "2025-03-10T06:00:00", "2025-03-10T23:00:00"), // 17h
    ];
    const violation: Violation = {
      rule_key: "max-shift-hours", rule_name: "Max Shift",
      severity: "mandatory", staff_id: "w1", message: "", citation: "",
      actual: 17, limit: 12,
    };
    const result = applyFix(shifts, violation, workers);
    expect(result.shifts.length).toBe(1);
    expect(result.description).toContain("Trimmed");
  });
});
