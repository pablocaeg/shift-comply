export interface Jurisdiction {
  code: string;
  name: string;
  local_name?: string;
  type: string;
  parent?: string;
  currency: string;
  timezone: string;
  rules: Rule[];
}

export interface Rule {
  key: string;
  name: string;
  description: string;
  category: string;
  operator: string;
  staff_types?: string[];
  unit_types?: string[];
  scope?: string;
  enforcement: string;
  values: RuleValue[];
  source: Source;
  notes?: string;
}

export interface RuleValue {
  since: string;
  amount: number;
  unit: string;
  per?: string;
  averaged?: { count: number; unit: string };
  exceptions?: string[];
}

export interface Source {
  title: string;
  section?: string;
  url?: string;
}

export interface Constraint {
  type: string;
  time_scope: string;
  facility_scope?: string;
  limit: number;
  limit_unit: string;
  operator: string;
  averaged_over_days?: number;
  staff_types?: string[];
  unit_types?: string[];
  enforcement: string;
  citation: string;
  jurisdiction: string;
  rule_key: string;
}

export interface Shift {
  _uid?: string;
  staff_id: string;
  staff_type: string;
  unit_type?: string;
  start: string;
  end: string;
  on_call?: boolean;
}

let _counter = 0;
export function tagShifts(shifts: Shift[]): Shift[] {
  return shifts.map(s => ({ ...s, _uid: s._uid || `s${++_counter}` }));
}
export function nextUid(): string { return `s${++_counter}`; }

export interface Schedule {
  jurisdiction: string;
  facility_scope?: string;
  shifts: Shift[];
}

export interface Violation {
  rule_key: string;
  rule_name: string;
  severity: string;
  staff_id: string;
  message: string;
  citation: string;
  actual: number;
  limit: number;
}

export interface ComplianceReport {
  jurisdiction: string;
  result: "pass" | "fail";
  violations: Violation[];
  constraints_checked: number;
}

export interface Comparison {
  left: string;
  right: string;
  only_left?: Rule[];
  only_right?: Rule[];
  different?: RulePair[];
  same?: RulePair[];
}

export interface RulePair {
  key: string;
  left: Rule;
  right: Rule;
}

export interface Worker {
  id: string;
  name: string;
  role: string;
  type: string;
  color: string;
}

export interface Scenario {
  id: string;
  badge: "fail" | "pass";
  label: string;
  name: string;
  who: string;
  info: string;
  jurisdiction: string;
  scope: string;
  flagged: string | null;
  workers: Worker[];
  shifts: Shift[];
}
