import { describe, it, expect } from "vitest";
import { dateOf, timeOf, hoursBetween, addDays, mondayOf, dayOfWeek, dayName, dayNum, monthShort, shiftOverlapsDate, clampToDay, formatDate, formatDateTime } from "./dates";

describe("dateOf", () => {
  it("extracts date from ISO string", () => {
    expect(dateOf("2025-03-10T08:00:00")).toBe("2025-03-10");
  });
});

describe("timeOf", () => {
  it("extracts time from ISO string", () => {
    expect(timeOf("2025-03-10T08:30:00")).toBe("08:30");
  });
});

describe("hoursBetween", () => {
  it("calculates hours between two times", () => {
    expect(hoursBetween("2025-03-10T08:00:00", "2025-03-10T20:00:00")).toBe(12);
  });
  it("handles overnight shifts", () => {
    expect(hoursBetween("2025-03-10T19:00:00", "2025-03-11T07:00:00")).toBe(12);
  });
  it("handles 24-hour shifts", () => {
    expect(hoursBetween("2025-03-10T08:00:00", "2025-03-11T08:00:00")).toBe(24);
  });
});

describe("addDays", () => {
  it("adds days correctly", () => {
    expect(addDays("2025-03-10", 1)).toBe("2025-03-11");
    expect(addDays("2025-03-31", 1)).toBe("2025-04-01");
  });
  it("subtracts days", () => {
    expect(addDays("2025-03-01", -1)).toBe("2025-02-28");
  });
});

describe("mondayOf", () => {
  it("returns Monday for a Monday", () => {
    expect(mondayOf("2025-03-10")).toBe("2025-03-10"); // Monday
  });
  it("returns Monday for a Wednesday", () => {
    expect(mondayOf("2025-03-12")).toBe("2025-03-10");
  });
  it("returns Monday for a Sunday", () => {
    expect(mondayOf("2025-03-16")).toBe("2025-03-10");
  });
});

describe("dayOfWeek", () => {
  it("returns correct day", () => {
    expect(dayOfWeek("2025-03-10")).toBe(1); // Monday
    expect(dayOfWeek("2025-03-16")).toBe(0); // Sunday
  });
});

describe("dayName / dayNum / monthShort", () => {
  it("returns day name", () => {
    expect(dayName("2025-03-10")).toBe("Mon");
  });
  it("returns day number", () => {
    expect(dayNum("2025-03-10")).toBe(10);
  });
  it("returns month short", () => {
    expect(monthShort("2025-03-10")).toBe("Mar");
  });
});

describe("shiftOverlapsDate", () => {
  it("detects overlap", () => {
    expect(shiftOverlapsDate("2025-03-10T08:00:00", "2025-03-10T20:00:00", "2025-03-10")).toBe(true);
  });
  it("detects no overlap", () => {
    expect(shiftOverlapsDate("2025-03-10T08:00:00", "2025-03-10T20:00:00", "2025-03-11")).toBe(false);
  });
  it("handles overnight shifts", () => {
    expect(shiftOverlapsDate("2025-03-10T19:00:00", "2025-03-11T07:00:00", "2025-03-10")).toBe(true);
    expect(shiftOverlapsDate("2025-03-10T19:00:00", "2025-03-11T07:00:00", "2025-03-11")).toBe(true);
    expect(shiftOverlapsDate("2025-03-10T19:00:00", "2025-03-11T07:00:00", "2025-03-12")).toBe(false);
  });
});

describe("clampToDay", () => {
  it("clamps a full-day shift", () => {
    const r = clampToDay("2025-03-10T08:00:00", "2025-03-10T20:00:00", "2025-03-10");
    expect(r.visStart).toBe("08:00");
    expect(r.visEnd).toBe("20:00");
    expect(r.visDur).toBe(12);
  });
  it("clamps an overnight shift to first day", () => {
    const r = clampToDay("2025-03-10T19:00:00", "2025-03-11T07:00:00", "2025-03-10");
    expect(r.visStart).toBe("19:00");
    expect(r.visEnd).toBe("24:00");
    expect(r.visDur).toBe(5);
  });
});

describe("formatDate", () => {
  it("formats a date", () => {
    expect(formatDate(new Date(2025, 2, 10))).toBe("2025-03-10");
  });
});

describe("formatDateTime", () => {
  it("formats a datetime", () => {
    expect(formatDateTime(new Date(2025, 2, 10, 8, 30, 0))).toBe("2025-03-10T08:30:00");
  });
});
