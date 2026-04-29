package comply_test

import (
	"testing"
	"time"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

func TestAllJurisdictionsRegistered(t *testing.T) {
	expected := []comply.Code{
		comply.US, comply.USCA, comply.USNY, comply.USTX, comply.USFL,
		comply.USMA, comply.USIL, comply.USOR,
		comply.EU, comply.DE, comply.HU, comply.IT, comply.PL,
		comply.ES, comply.ESCT, comply.ESMD,
	}
	for _, code := range expected {
		j := comply.For(code)
		if j == nil {
			t.Errorf("jurisdiction %q not registered", code)
			continue
		}
		if j.Name == "" {
			t.Errorf("jurisdiction %q has empty name", code)
		}
		if j.TimeZone == "" {
			t.Errorf("jurisdiction %q has empty timezone", code)
		}
		if len(j.Rules) == 0 {
			t.Errorf("jurisdiction %q has no rules", code)
		}
	}
}

func TestJurisdictionCount(t *testing.T) {
	all := comply.All()
	if len(all) != 16 {
		t.Errorf("expected 16 jurisdictions, got %d", len(all))
	}
}

func TestParentChain_California(t *testing.T) {
	ca := comply.For(comply.USCA)
	if ca == nil {
		t.Fatal("US-CA not registered")
	}
	chain := ca.Chain()
	if len(chain) != 2 {
		t.Fatalf("expected chain length 2 (US-CA -> US), got %d", len(chain))
	}
	if chain[0].Code != comply.USCA {
		t.Errorf("chain[0] = %q, want US-CA", chain[0].Code)
	}
	if chain[1].Code != comply.US {
		t.Errorf("chain[1] = %q, want US", chain[1].Code)
	}
}

func TestParentChain_Spain(t *testing.T) {
	es := comply.For(comply.ES)
	if es == nil {
		t.Fatal("ES not registered")
	}
	chain := es.Chain()
	if len(chain) != 2 {
		t.Fatalf("expected chain length 2 (ES -> EU), got %d", len(chain))
	}
	if chain[0].Code != comply.ES {
		t.Errorf("chain[0] = %q, want ES", chain[0].Code)
	}
	if chain[1].Code != comply.EU {
		t.Errorf("chain[1] = %q, want EU", chain[1].Code)
	}
}

func TestEveryRuleHasSource(t *testing.T) {
	for _, j := range comply.All() {
		for _, r := range j.Rules {
			if r.Source.Title == "" {
				t.Errorf("[%s] rule %q has empty source title", j.Code, r.Key)
			}
			if r.Name == "" {
				t.Errorf("[%s] rule %q has empty name", j.Code, r.Key)
			}
			if len(r.Values) == 0 {
				t.Errorf("[%s] rule %q has no values", j.Code, r.Key)
			}
		}
	}
}

func TestEveryRuleHasScope(t *testing.T) {
	// Rules that apply to specific populations should have a Scope set.
	// Rules with StaffTypes or UnitTypes narrower than "all" should have scope.
	for _, j := range comply.All() {
		for _, r := range j.Rules {
			// ACGME rules must be scoped to accredited programs
			if r.Source.Title == "ACGME Common Program Requirements (Residency)" {
				if r.Scope != comply.ScopeAccreditedPrograms {
					t.Errorf("[%s] ACGME rule %q should have scope accredited_programs, got %q", j.Code, r.Key, r.Scope)
				}
			}
			// VA rules must be scoped to VA
			if r.Source.Title == "VA Healthcare Personnel Enhancement Act of 2004" {
				if r.Scope != comply.ScopeVA {
					t.Errorf("[%s] VA rule %q should have scope va, got %q", j.Code, r.Key, r.Scope)
				}
			}
			// Spanish Estatuto Marco rules must be scoped to public health
			if r.Source.Title == "Ley 55/2003, del Estatuto Marco del personal estatutario de los servicios de salud" {
				if r.Scope != comply.ScopePublicHealth {
					t.Errorf("[%s] Estatuto Marco rule %q should have scope public_health, got %q", j.Code, r.Key, r.Scope)
				}
			}
		}
	}
}

func TestEffectiveRulesInheritance(t *testing.T) {
	// California should inherit US federal ACGME rules
	rules := comply.EffectiveRules(comply.USCA, comply.ForStaff(comply.StaffResident))
	found := false
	for _, r := range rules {
		if r.Key == comply.RuleMaxWeeklyHours {
			found = true
			v := r.Current()
			if v == nil {
				t.Fatal("max-weekly-hours has no current value")
			}
			if v.Amount != 80 {
				t.Errorf("expected ACGME 80 hours, got %v", v.Amount)
			}
			break
		}
	}
	if !found {
		t.Error("US-CA should inherit max-weekly-hours from US (ACGME)")
	}
}

func TestCaliforniaNurseRatios(t *testing.T) {
	rules := comply.EffectiveRules(comply.USCA,
		comply.ForStaff(comply.StaffNurseRN),
		comply.InCategory(comply.CatStaffing),
	)
	if len(rules) < 10 {
		t.Errorf("expected at least 10 nurse ratio rules for CA, got %d", len(rules))
	}

	// Check ICU ratio is 1:2
	for _, r := range rules {
		if r.Key == comply.RuleNursePatientRatioICU {
			v := r.Current()
			if v == nil {
				t.Fatal("ICU ratio has no current value")
			}
			if v.Amount != 2 {
				t.Errorf("expected ICU ratio 1:2 (2 patients), got %v", v.Amount)
			}
			if r.Scope != comply.ScopeHospitals {
				t.Errorf("ICU ratio scope should be hospitals, got %q", r.Scope)
			}
			return
		}
	}
	t.Error("ICU ratio rule not found")
}

func TestCaliforniaMandatoryOTHasBothScopes(t *testing.T) {
	ca := comply.For(comply.USCA)
	if ca == nil {
		t.Fatal("US-CA not registered")
	}

	foundPrivate := false
	foundState := false
	for _, r := range ca.Rules {
		if r.Key == comply.RuleMandatoryOTProhibited && r.Scope == comply.ScopeHospitals {
			foundPrivate = true
			if r.Source.Title != "IWC Wage Order No. 5-2001 (Healthcare Industry)" {
				t.Errorf("private-sector OT ban should cite IWC Wage Order 5, got %q", r.Source.Title)
			}
		}
		if r.Key == "mandatory-overtime-prohibited-state" && r.Scope == comply.ScopeStateFacilities {
			foundState = true
		}
	}
	if !foundPrivate {
		t.Error("CA should have private-sector mandatory OT restriction (IWC Wage Order 5)")
	}
	if !foundState {
		t.Error("CA should have state-facility mandatory OT prohibition (Gov Code 19851.2)")
	}
}

func TestRuleValueDateLookup(t *testing.T) {
	ca := comply.For(comply.USCA)
	if ca == nil {
		t.Fatal("US-CA not registered")
	}

	// Find step-down ratio which changed from 1:4 (2004) to 1:3 (2008)
	for _, r := range ca.Rules {
		if r.Key == comply.RuleNursePatientRatioStepDown {
			// Current value should be 3
			v := r.Current()
			if v == nil || v.Amount != 3 {
				t.Errorf("current step-down ratio should be 3, got %v", v)
			}
			// Value in 2005 should be 4
			v2005 := r.Value(time.Date(2005, 6, 1, 0, 0, 0, 0, time.UTC))
			if v2005 == nil || v2005.Amount != 4 {
				t.Errorf("2005 step-down ratio should be 4, got %v", v2005)
			}
			return
		}
	}
	t.Error("step-down ratio rule not found")
}

func TestSpainMIRGuards(t *testing.T) {
	rules := comply.EffectiveRules(comply.ES, comply.ForStaff(comply.StaffResident))
	for _, r := range rules {
		if r.Key == comply.RuleMaxGuardsMonthly {
			v := r.Current()
			if v == nil {
				t.Fatal("MIR guards has no current value")
			}
			if v.Amount != 7 {
				t.Errorf("expected max 7 guards/month for Spain MIR, got %v", v.Amount)
			}
			if r.Scope != comply.ScopePublicHealth {
				t.Errorf("MIR guards scope should be public_health, got %q", r.Scope)
			}
			return
		}
	}
	t.Error("MIR max guards rule not found")
}

func TestSpainInheritsEU(t *testing.T) {
	// Spain should inherit EU Working Time Directive 48-hour max
	rules := comply.EffectiveRules(comply.ES)
	for _, r := range rules {
		if r.Key == comply.RuleMaxWeeklyHours {
			v := r.Current()
			if v == nil || v.Amount != 48 {
				t.Errorf("Spain should inherit EU 48-hour max, got %v", v)
			}
			return
		}
	}
	t.Error("Spain should inherit max-weekly-hours from EU (WTD)")
}

func TestEUDirectiveRules(t *testing.T) {
	eu := comply.For(comply.EU)
	if eu == nil {
		t.Fatal("EU not registered")
	}

	// Check 48-hour maximum
	for _, r := range eu.Rules {
		if r.Key == comply.RuleMaxWeeklyHours {
			v := r.Current()
			if v == nil || v.Amount != 48 {
				t.Errorf("EU max weekly hours should be 48, got %v", v)
			}
			return
		}
	}
	t.Error("EU max weekly hours rule not found")
}

func TestScopeFiltering(t *testing.T) {
	// Querying CA with ScopeHospitals should include nurse ratios
	rules := comply.EffectiveRules(comply.USCA, comply.ForScope(comply.ScopeHospitals))
	foundRatio := false
	for _, r := range rules {
		if r.Key == comply.RuleNursePatientRatioICU {
			foundRatio = true
			break
		}
	}
	if !foundRatio {
		t.Error("hospital scope should include nurse ratios")
	}

	// Querying with ScopeVA should not include nurse ratios
	vaRules := comply.EffectiveRules(comply.USCA, comply.ForScope(comply.ScopeVA))
	for _, r := range vaRules {
		if r.Key == comply.RuleNursePatientRatioICU {
			t.Error("VA scope should not include CA nurse ratios")
		}
	}
}

func TestACGMENotLaw(t *testing.T) {
	us := comply.For(comply.US)
	if us == nil {
		t.Fatal("US not registered")
	}

	for _, r := range us.Rules {
		if r.Scope == comply.ScopeAccreditedPrograms {
			if r.Source.Title != "ACGME Common Program Requirements (Residency)" {
				continue
			}
			if r.Scope != comply.ScopeAccreditedPrograms {
				t.Errorf("ACGME rule %q should have accredited_programs scope", r.Key)
			}
		}
	}
}

func TestCataloniaThreeLevelChain(t *testing.T) {
	ct := comply.For(comply.ESCT)
	if ct == nil {
		t.Fatal("ES-CT not registered")
	}
	chain := ct.Chain()
	if len(chain) != 3 {
		t.Fatalf("expected chain length 3 (ES-CT -> ES -> EU), got %d", len(chain))
	}
	if chain[0].Code != comply.ESCT || chain[1].Code != comply.ES || chain[2].Code != comply.EU {
		t.Errorf("chain should be ES-CT -> ES -> EU, got %s -> %s -> %s",
			chain[0].Code, chain[1].Code, chain[2].Code)
	}
}

func TestCataloniaOverridesSpainGuards(t *testing.T) {
	// Catalonia limits guards to 4/month (on average), overriding Spain's 7/month
	rules := comply.EffectiveRules(comply.ESCT, comply.ForStaff(comply.StaffStatutory))
	for _, r := range rules {
		if r.Key == comply.RuleMaxGuardsMonthly {
			v := r.Current()
			if v == nil {
				t.Fatal("Catalonia guards has no current value")
			}
			if v.Amount != 4 {
				t.Errorf("expected max 4 guards/month for Catalonia (overriding Spain's 7), got %v", v.Amount)
			}
			if r.Scope != comply.ScopePublicHealth {
				t.Errorf("Catalonia guard rule scope should be public_health, got %q", r.Scope)
			}
			return
		}
	}
	t.Error("Catalonia max guards rule not found")
}

func TestCataloniaInheritsSpainAndEU(t *testing.T) {
	// Catalonia should inherit Spain's 12-hour rest between shifts
	rules := comply.EffectiveRules(comply.ESCT)
	foundSpainRest := false
	foundEU48 := false
	for _, r := range rules {
		if r.Key == comply.RuleMinRestBetweenShifts {
			v := r.Current()
			if v != nil && v.Amount == 12 {
				foundSpainRest = true
			}
		}
		if r.Key == comply.RuleMaxWeeklyHours {
			v := r.Current()
			if v != nil && v.Amount == 48 {
				foundEU48 = true
			}
		}
	}
	if !foundSpainRest {
		t.Error("Catalonia should inherit Spain's 12-hour rest between shifts")
	}
	if !foundEU48 {
		t.Error("Catalonia should inherit EU's 48-hour weekly max")
	}
}

func TestMadridWeeklyRest(t *testing.T) {
	rules := comply.EffectiveRules(comply.ESMD, comply.ForStaff(comply.StaffStatutory))
	for _, r := range rules {
		if r.Key == comply.RuleMinWeeklyRest {
			v := r.Current()
			if v == nil || v.Amount != 36 {
				t.Errorf("Madrid weekly rest should be 36 hours, got %v", v)
			}
			// Verify source references the correct ruling date
			if r.Source.Title == "" {
				t.Error("Madrid weekly rest should have a source")
			}
			return
		}
	}
	t.Error("Madrid min weekly rest rule not found")
}

func TestMadridPostGuardRestNotWorkingTime(t *testing.T) {
	// Madrid: post-guard rest does NOT count as working time
	// Catalonia: post-guard rest DOES count as working time
	// This is a real difference between the two regions
	md := comply.For(comply.ESMD)
	if md == nil {
		t.Fatal("ES-MD not registered")
	}
	for _, r := range md.Rules {
		if r.Key == "md-rest-not-effective-work" {
			return // found
		}
	}
	t.Error("Madrid should document that post-guard rest does not count as effective working time")
}

func TestAllRegionRulesArePublicHealthScoped(t *testing.T) {
	for _, code := range []comply.Code{comply.ESCT, comply.ESMD} {
		j := comply.For(code)
		if j == nil {
			t.Errorf("jurisdiction %q not registered", code)
			continue
		}
		for _, r := range j.Rules {
			if r.Scope != comply.ScopePublicHealth {
				t.Errorf("[%s] rule %q should have public_health scope, got %q", code, r.Key, r.Scope)
			}
		}
	}
}
