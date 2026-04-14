// Command wasm compiles the Shift Comply engine to WebAssembly.
//
// Build:
//
//	GOOS=js GOARCH=wasm go build -o shiftcomply.wasm ./cmd/wasm
//
// The resulting .wasm file exposes these functions to JavaScript:
//
//	shiftcomply.jurisdictions()                                     -> JSON array
//	shiftcomply.rules(jurisdiction, staff?, unit?, scope?)          -> JSON array
//	shiftcomply.constraints(jurisdiction, staff?, unit?, scope?)    -> JSON array
//	shiftcomply.compare(left, right, staff?)                        -> JSON object
//	shiftcomply.validate(scheduleJSON)                              -> JSON object
//	shiftcomply.export(jurisdiction)                                -> JSON object
//
//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

func main() {
	sc := js.Global().Get("Object").New()

	sc.Set("jurisdictions", js.FuncOf(jsJurisdictions))
	sc.Set("rules", js.FuncOf(jsRules))
	sc.Set("constraints", js.FuncOf(jsConstraints))
	sc.Set("compare", js.FuncOf(jsCompare))
	sc.Set("validate", js.FuncOf(jsValidate))
	sc.Set("export", js.FuncOf(jsExport))

	js.Global().Set("shiftcomply", sc)

	// Block forever so the WASM module stays alive
	select {}
}

func jsJurisdictions(this js.Value, args []js.Value) any {
	return toJSON(comply.All())
}

func jsRules(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return toJSON(map[string]string{"error": "jurisdiction required"})
	}
	code := comply.Code(args[0].String())
	opts := parseOpts(args)
	return toJSON(comply.EffectiveRules(code, opts...))
}

func jsConstraints(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return toJSON(map[string]string{"error": "jurisdiction required"})
	}
	code := comply.Code(args[0].String())
	opts := parseOpts(args)
	return toJSON(comply.GenerateConstraints(code, opts...))
}

func jsCompare(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return toJSON(map[string]string{"error": "left and right jurisdictions required"})
	}
	left := comply.Code(args[0].String())
	right := comply.Code(args[1].String())
	var opts []comply.QueryOption
	if len(args) > 2 && args[2].String() != "" {
		opts = append(opts, comply.ForStaff(comply.Key(args[2].String())))
	}
	return toJSON(comply.Compare(left, right, opts...))
}

func jsValidate(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return toJSON(map[string]string{"error": "schedule JSON required"})
	}
	var schedule comply.Schedule
	if err := json.Unmarshal([]byte(args[0].String()), &schedule); err != nil {
		return toJSON(map[string]string{"error": "invalid JSON: " + err.Error()})
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		return toJSON(map[string]string{"error": err.Error()})
	}
	return toJSON(report)
}

func jsExport(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return toJSON(map[string]string{"error": "jurisdiction required"})
	}
	code := comply.Code(args[0].String())
	j := comply.For(code)
	if j == nil {
		return toJSON(map[string]string{"error": "unknown jurisdiction: " + string(code)})
	}
	return toJSON(j)
}

func parseOpts(args []js.Value) []comply.QueryOption {
	var opts []comply.QueryOption
	if len(args) > 1 && args[1].String() != "" {
		opts = append(opts, comply.ForStaff(comply.Key(args[1].String())))
	}
	if len(args) > 2 && args[2].String() != "" {
		opts = append(opts, comply.ForUnit(comply.Key(args[2].String())))
	}
	if len(args) > 3 && args[3].String() != "" {
		opts = append(opts, comply.ForScope(comply.Scope(args[3].String())))
	}
	return opts
}

func toJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
