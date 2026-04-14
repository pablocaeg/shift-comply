// Command shiftcomply is a CLI tool for querying healthcare scheduling
// regulations across jurisdictions.
//
// Usage:
//
//	shiftcomply jurisdictions                          # list all jurisdictions
//	shiftcomply rules US-CA                            # all rules for California
//	shiftcomply rules US-CA --staff nurse-rn --unit icu # filtered rules
//	shiftcomply compare US-CA US-TX                    # diff two jurisdictions
//	shiftcomply constraints US-CA --staff nurse-rn     # optimizer-ready output
//	shiftcomply export US-CA                           # full JSON export
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "jurisdictions", "j":
		cmdJurisdictions()
	case "rules", "r":
		cmdRules(os.Args[2:])
	case "compare", "cmp":
		cmdCompare(os.Args[2:])
	case "constraints", "c":
		cmdConstraints(os.Args[2:])
	case "export", "e":
		cmdExport(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`shift-comply - Healthcare scheduling regulation engine

Usage:
  shiftcomply <command> [arguments]

Commands:
  jurisdictions          List all registered jurisdictions
  rules <code>           Show rules for a jurisdiction
  compare <a> <b>        Compare two jurisdictions
  constraints <code>     Generate optimizer-ready constraints
  export <code>          Export jurisdiction data as JSON

Flags (for rules, constraints):
  --staff <type>         Filter by staff type (e.g., nurse-rn, resident)
  --unit <type>          Filter by hospital unit (e.g., icu, ed)
  --category <cat>       Filter by category (e.g., work_hours, rest, overtime)
  --json                 Output as JSON instead of table`)
}

// jurisdictions

func cmdJurisdictions() {
	all := comply.All()
	sort.Slice(all, func(i, j int) bool { return all[i].Code < all[j].Code })

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CODE\tNAME\tTYPE\tPARENT\tRULES")
	for _, j := range all {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
			j.Code, j.Name, j.Type, j.Parent, len(j.Rules))
	}
	w.Flush()

	total := 0
	for _, j := range all {
		total += len(j.Rules)
	}
	fmt.Printf("\n%d jurisdictions, %d total rules\n", len(all), total)
}

// rules

func cmdRules(args []string) {
	fs := flag.NewFlagSet("rules", flag.ExitOnError)
	staff := fs.String("staff", "", "filter by staff type")
	unit := fs.String("unit", "", "filter by hospital unit")
	category := fs.String("category", "", "filter by category")
	asJSON := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: shiftcomply rules <jurisdiction-code> [--staff X] [--unit X] [--category X] [--json]")
		os.Exit(1)
	}

	code := comply.Code(fs.Arg(0))
	if comply.For(code) == nil {
		fmt.Fprintf(os.Stderr, "unknown jurisdiction: %s\n", code)
		os.Exit(1)
	}

	var opts []comply.QueryOption
	if *staff != "" {
		opts = append(opts, comply.ForStaff(comply.Key(*staff)))
	}
	if *unit != "" {
		opts = append(opts, comply.ForUnit(comply.Key(*unit)))
	}
	if *category != "" {
		opts = append(opts, comply.InCategory(comply.Category(*category)))
	}

	rules := comply.EffectiveRules(code, opts...)

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(rules)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE\tUNIT\tPER\tENFORCEMENT\tSTAFF\tCATEGORY")
	for _, r := range rules {
		v := r.Current()
		if v == nil {
			continue
		}
		staffStr := "all"
		if len(r.StaffTypes) > 0 {
			staffStr = joinKeys(r.StaffTypes)
		}
		perStr := string(v.Per)
		if v.Averaged != nil {
			perStr = fmt.Sprintf("%s (avg %d%s)", v.Per, v.Averaged.Count, v.Averaged.Unit)
		}
		fmt.Fprintf(w, "%s\t%s%.4g\t%s\t%s\t%s\t%s\t%s\n",
			r.Key, opSymbol(r.Operator), v.Amount, v.Unit, perStr, r.Enforcement, staffStr, r.Category)
	}
	w.Flush()
	fmt.Printf("\n%d rules (effective, including inherited)\n", len(rules))
}

// compare

func cmdCompare(args []string) {
	fs := flag.NewFlagSet("compare", flag.ExitOnError)
	staff := fs.String("staff", "", "filter by staff type")
	asJSON := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	if fs.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "usage: shiftcomply compare <code-a> <code-b> [--staff X] [--json]")
		os.Exit(1)
	}

	a, b := comply.Code(fs.Arg(0)), comply.Code(fs.Arg(1))
	var opts []comply.QueryOption
	if *staff != "" {
		opts = append(opts, comply.ForStaff(comply.Key(*staff)))
	}

	comp := comply.Compare(a, b, opts...)

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(comp)
		return
	}

	fmt.Printf("Comparing %s vs %s\n\n", a, b)

	if len(comp.OnlyLeft) > 0 {
		fmt.Printf("--- Only in %s (%d rules) ---\n", a, len(comp.OnlyLeft))
		for _, r := range comp.OnlyLeft {
			v := r.Current()
			if v == nil {
				continue
			}
			fmt.Printf("  %s: %s%.4g %s\n", r.Key, opSymbol(r.Operator), v.Amount, v.Unit)
		}
		fmt.Println()
	}

	if len(comp.OnlyRight) > 0 {
		fmt.Printf("--- Only in %s (%d rules) ---\n", b, len(comp.OnlyRight))
		for _, r := range comp.OnlyRight {
			v := r.Current()
			if v == nil {
				continue
			}
			fmt.Printf("  %s: %s%.4g %s\n", r.Key, opSymbol(r.Operator), v.Amount, v.Unit)
		}
		fmt.Println()
	}

	if len(comp.Different) > 0 {
		fmt.Printf("--- Different values (%d rules) ---\n", len(comp.Different))
		for _, p := range comp.Different {
			lv, rv := p.Left.Current(), p.Right.Current()
			if lv == nil || rv == nil {
				continue
			}
			fmt.Printf("  %s: %s%.4g %s vs %s%.4g %s\n",
				p.Key,
				opSymbol(p.Left.Operator), lv.Amount, lv.Unit,
				opSymbol(p.Right.Operator), rv.Amount, rv.Unit)
		}
		fmt.Println()
	}

	fmt.Printf("Summary: %d only-%s, %d only-%s, %d different, %d same\n",
		len(comp.OnlyLeft), a, len(comp.OnlyRight), b,
		len(comp.Different), len(comp.Same))
}

// constraints

func cmdConstraints(args []string) {
	fs := flag.NewFlagSet("constraints", flag.ExitOnError)
	staff := fs.String("staff", "", "filter by staff type")
	unit := fs.String("unit", "", "filter by hospital unit")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: shiftcomply constraints <code> [--staff X] [--unit X]")
		os.Exit(1)
	}

	code := comply.Code(fs.Arg(0))
	if comply.For(code) == nil {
		fmt.Fprintf(os.Stderr, "unknown jurisdiction: %s\n", code)
		os.Exit(1)
	}

	var opts []comply.QueryOption
	if *staff != "" {
		opts = append(opts, comply.ForStaff(comply.Key(*staff)))
	}
	if *unit != "" {
		opts = append(opts, comply.ForUnit(comply.Key(*unit)))
	}

	constraints := comply.GenerateConstraints(code, opts...)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(constraints)
}

// export

func cmdExport(args []string) {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: shiftcomply export <code>")
		os.Exit(1)
	}

	code := comply.Code(fs.Arg(0))
	j := comply.For(code)
	if j == nil {
		fmt.Fprintf(os.Stderr, "unknown jurisdiction: %s\n", code)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(j)
}

// helpers

func opSymbol(op comply.Operator) string {
	switch op {
	case comply.OpLTE:
		return "<="
	case comply.OpGTE:
		return ">="
	case comply.OpEQ:
		return "=="
	case comply.OpBool:
		return ""
	default:
		return string(op)
	}
}

func joinKeys(keys []comply.Key) string {
	s := make([]string, len(keys))
	for i, k := range keys {
		s[i] = string(k)
	}
	return strings.Join(s, ",")
}
