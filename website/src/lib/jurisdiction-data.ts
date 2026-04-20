// Maps US state FIPS codes to shift-comply jurisdiction codes
// FIPS codes are used in the topojson file as feature IDs

export const FIPS_TO_CODE: Record<string, string> = {
  "06": "US-CA", "36": "US-NY", "48": "US-TX", "12": "US-FL",
  "01": "US-AL", "02": "US-AK", "04": "US-AZ", "05": "US-AR",
  "08": "US-CO", "09": "US-CT", "10": "US-DE", "11": "US-DC",
  "13": "US-GA", "15": "US-HI", "16": "US-ID", "17": "US-IL",
  "18": "US-IN", "19": "US-IA", "20": "US-KS", "21": "US-KY",
  "22": "US-LA", "23": "US-ME", "24": "US-MD", "25": "US-MA",
  "26": "US-MI", "27": "US-MN", "28": "US-MS", "29": "US-MO",
  "30": "US-MT", "31": "US-NE", "32": "US-NV", "33": "US-NH",
  "34": "US-NJ", "35": "US-NM", "37": "US-NC", "38": "US-ND",
  "39": "US-OH", "40": "US-OK", "41": "US-OR", "42": "US-PA",
  "44": "US-RI", "45": "US-SC", "46": "US-SD", "47": "US-TN",
  "49": "US-UT", "50": "US-VT", "51": "US-VA", "53": "US-WA",
  "54": "US-WV", "55": "US-WI", "56": "US-WY",
};

export const STATE_NAMES: Record<string, string> = {
  "US-CA": "California", "US-NY": "New York", "US-TX": "Texas", "US-FL": "Florida",
  "US-AL": "Alabama", "US-AK": "Alaska", "US-AZ": "Arizona", "US-AR": "Arkansas",
  "US-CO": "Colorado", "US-CT": "Connecticut", "US-DE": "Delaware", "US-DC": "District of Columbia",
  "US-GA": "Georgia", "US-HI": "Hawaii", "US-ID": "Idaho", "US-IL": "Illinois",
  "US-IN": "Indiana", "US-IA": "Iowa", "US-KS": "Kansas", "US-KY": "Kentucky",
  "US-LA": "Louisiana", "US-ME": "Maine", "US-MD": "Maryland", "US-MA": "Massachusetts",
  "US-MI": "Michigan", "US-MN": "Minnesota", "US-MS": "Mississippi", "US-MO": "Missouri",
  "US-MT": "Montana", "US-NE": "Nebraska", "US-NV": "Nevada", "US-NH": "New Hampshire",
  "US-NJ": "New Jersey", "US-NM": "New Mexico", "US-NC": "North Carolina", "US-ND": "North Dakota",
  "US-OH": "Ohio", "US-OK": "Oklahoma", "US-OR": "Oregon", "US-PA": "Pennsylvania",
  "US-RI": "Rhode Island", "US-SC": "South Carolina", "US-SD": "South Dakota", "US-TN": "Tennessee",
  "US-UT": "Utah", "US-VT": "Vermont", "US-VA": "Virginia", "US-WA": "Washington",
  "US-WV": "West Virginia", "US-WI": "Wisconsin", "US-WY": "Wyoming",
};

export interface JurisdictionInfo {
  code: string;
  name: string;
  ruleCount: number;
  parent?: string;
  type: string;
}
