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

// Maps ISO 3166-1 numeric codes to shift-comply jurisdiction codes (for world map)
export const COUNTRY_NUMERIC_TO_CODE: Record<string, string> = {
  "840": "US",   // United States
  "724": "ES",   // Spain
};

// EU member state numeric codes (inherit EU rules)
export const EU_MEMBERS = new Set([
  "040", "056", "100", "191", "196", "203", "208", "233", "246", "250",
  "276", "300", "348", "372", "380", "428", "440", "442", "470", "528",
  "616", "620", "642", "703", "705", "724", "752",
]);

export interface JurisdictionInfo {
  code: string;
  name: string;
  ruleCount: number;
  parent?: string;
  type: string;
}

// Healthcare sector stats for covered jurisdictions.
// Sources: AHA Annual Survey 2024, BLS OEWS May 2024, INE 2024,
// Ministerio de Sanidad CNH 2025, Eurostat 2023.
// Sub-national stats are NOT included to avoid double-counting with their parent.
export interface JurisdictionStats {
  hospitals: number;
  healthcareWorkers: number;
}

export const JURISDICTION_STATS: Record<string, JurisdictionStats> = {
  // US Federal: AHA 2024 (6,100 total registered), BLS OEWS May 2024 (18.2M)
  "US":    { hospitals: 6_100,  healthcareWorkers: 18_200_000 },
  // US States: AHA via KFF 2024 (community hospitals), KFF/BLS QCEW 2023 (hospital employees)
  "US-CA": { hospitals: 350,    healthcareWorkers: 610_000 },
  "US-NY": { hospitals: 158,    healthcareWorkers: 465_000 },
  "US-TX": { hospitals: 503,    healthcareWorkers: 500_000 },
  "US-FL": { hospitals: 229,    healthcareWorkers: 425_000 },
  "US-MA": { hospitals: 70,     healthcareWorkers: 214_000 },
  "US-IL": { hospitals: 183,    healthcareWorkers: 274_000 },
  "US-OR": { hospitals: 60,     healthcareWorkers: 68_000 },
  // EU: Eurostat 2023 (physicians + nurses ~ 5.7M), hospital count not published by Eurostat
  "EU":    { hospitals: 15_000, healthcareWorkers: 5_700_000 },
  // Spain: Ministerio de Sanidad CNH 2025 (848 hospitals), INE 2024 (1M professionals)
  "ES":    { hospitals: 848,    healthcareWorkers: 1_009_000 },
  // Catalonia: CNH 2025 (204 hospitals), INE/Idescat 2024 (~115K physicians+nurses)
  "ES-CT": { hospitals: 204,    healthcareWorkers: 115_000 },
  // Madrid: CNH 2025 (91 hospitals), INE 2024 (~110K physicians+nurses)
  "ES-MD": { hospitals: 91,     healthcareWorkers: 110_000 },
};
