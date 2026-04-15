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
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		printUsage(stdout)
		return 1
	}

	switch args[0] {
	case "jurisdictions", "j":
		return runJurisdictions(stdout)
	case "rules", "r":
		return runRules(args[1:], stdout, stderr)
	case "compare", "cmp":
		return runCompare(args[1:], stdout, stderr)
	case "constraints", "c":
		return runConstraints(args[1:], stdout, stderr)
	case "export", "e":
		return runExport(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n\n", args[0])
		printUsage(stderr)
		return 1
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, `shift-comply - Healthcare scheduling regulation engine

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

func runJurisdictions(w io.Writer) int {
	all := comply.All()
	sort.Slice(all, func(i, j int) bool { return all[i].Code < all[j].Code })

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "CODE\tNAME\tTYPE\tPARENT\tRULES")
	for _, j := range all {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\n", j.Code, j.Name, j.Type, j.Parent, len(j.Rules))
	}
	tw.Flush()

	total := 0
	for _, j := range all {
		total += len(j.Rules)
	}
	fmt.Fprintf(w, "\n%d jurisdictions, %d total rules\n", len(all), total)
	return 0
}

func runRules(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("rules", flag.ContinueOnError)
	fs.SetOutput(stderr)
	staff := fs.String("staff", "", "filter by staff type")
	unit := fs.String("unit", "", "filter by hospital unit")
	category := fs.String("category", "", "filter by category")
	asJSON := fs.Bool("json", false, "output as JSON")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(stderr, "usage: shiftcomply rules <jurisdiction-code> [--staff X] [--unit X] [--category X] [--json]")
		return 1
	}

	code := comply.Code(fs.Arg(0))
	if comply.For(code) == nil {
		fmt.Fprintf(stderr, "unknown jurisdiction: %s\n", code)
		return 1
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
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		enc.Encode(rules)
		return 0
	}

	tw := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tVALUE\tUNIT\tPER\tENFORCEMENT\tSTAFF\tCATEGORY")
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
		fmt.Fprintf(tw, "%s\t%s%.4g\t%s\t%s\t%s\t%s\t%s\n",
			r.Key, opSymbol(r.Operator), v.Amount, v.Unit, perStr, r.Enforcement, staffStr, r.Category)
	}
	tw.Flush()
	fmt.Fprintf(stdout, "\n%d rules (effective, including inherited)\n", len(rules))
	return 0
}

func runCompare(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	fs.SetOutput(stderr)
	staff := fs.String("staff", "", "filter by staff type")
	asJSON := fs.Bool("json", false, "output as JSON")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() < 2 {
		fmt.Fprintln(stderr, "usage: shiftcomply compare <code-a> <code-b> [--staff X] [--json]")
		return 1
	}

	a, b := comply.Code(fs.Arg(0)), comply.Code(fs.Arg(1))
	var opts []comply.QueryOption
	if *staff != "" {
		opts = append(opts, comply.ForStaff(comply.Key(*staff)))
	}

	comp := comply.Compare(a, b, opts...)

	if *asJSON {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		enc.Encode(comp)
		return 0
	}

	fmt.Fprintf(stdout, "Comparing %s vs %s\n\n", a, b)

	if len(comp.OnlyLeft) > 0 {
		fmt.Fprintf(stdout, "--- Only in %s (%d rules) ---\n", a, len(comp.OnlyLeft))
		for _, r := range comp.OnlyLeft {
			v := r.Current()
			if v == nil {
				continue
			}
			fmt.Fprintf(stdout, "  %s: %s%.4g %s\n", r.Key, opSymbol(r.Operator), v.Amount, v.Unit)
		}
		fmt.Fprintln(stdout)
	}

	if len(comp.OnlyRight) > 0 {
		fmt.Fprintf(stdout, "--- Only in %s (%d rules) ---\n", b, len(comp.OnlyRight))
		for _, r := range comp.OnlyRight {
			v := r.Current()
			if v == nil {
				continue
			}
			fmt.Fprintf(stdout, "  %s: %s%.4g %s\n", r.Key, opSymbol(r.Operator), v.Amount, v.Unit)
		}
		fmt.Fprintln(stdout)
	}

	if len(comp.Different) > 0 {
		fmt.Fprintf(stdout, "--- Different values (%d rules) ---\n", len(comp.Different))
		for _, p := range comp.Different {
			lv, rv := p.Left.Current(), p.Right.Current()
			if lv == nil || rv == nil {
				continue
			}
			fmt.Fprintf(stdout, "  %s: %s%.4g %s vs %s%.4g %s\n",
				p.Key, opSymbol(p.Left.Operator), lv.Amount, lv.Unit,
				opSymbol(p.Right.Operator), rv.Amount, rv.Unit)
		}
		fmt.Fprintln(stdout)
	}

	fmt.Fprintf(stdout, "Summary: %d only-%s, %d only-%s, %d different, %d same\n",
		len(comp.OnlyLeft), a, len(comp.OnlyRight), b, len(comp.Different), len(comp.Same))
	return 0
}

func runConstraints(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("constraints", flag.ContinueOnError)
	fs.SetOutput(stderr)
	staff := fs.String("staff", "", "filter by staff type")
	unit := fs.String("unit", "", "filter by hospital unit")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(stderr, "usage: shiftcomply constraints <code> [--staff X] [--unit X]")
		return 1
	}

	code := comply.Code(fs.Arg(0))
	if comply.For(code) == nil {
		fmt.Fprintf(stderr, "unknown jurisdiction: %s\n", code)
		return 1
	}

	var opts []comply.QueryOption
	if *staff != "" {
		opts = append(opts, comply.ForStaff(comply.Key(*staff)))
	}
	if *unit != "" {
		opts = append(opts, comply.ForUnit(comply.Key(*unit)))
	}

	constraints := comply.GenerateConstraints(code, opts...)
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	enc.Encode(constraints)
	return 0
}

func runExport(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(stderr, "usage: shiftcomply export <code>")
		return 1
	}

	code := comply.Code(fs.Arg(0))
	j := comply.For(code)
	if j == nil {
		fmt.Fprintf(stderr, "unknown jurisdiction: %s\n", code)
		return 1
	}

	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	enc.Encode(j)
	return 0
}

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
