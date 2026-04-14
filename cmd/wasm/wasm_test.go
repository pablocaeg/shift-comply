//go:build !js

// These tests verify the WASM build works by building the binary
// and checking that it produces a valid .wasm file.
// The actual JS interop (syscall/js) can only run in a browser/Node
// environment, so we test the build artifact instead.

package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestWASM_Builds(t *testing.T) {
	tmpFile := t.TempDir() + "/shiftcomply.wasm"
	cmd := exec.Command("go", "build", "-o", tmpFile, ".")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("WASM build failed: %v\n%s", err, out)
	}

	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("WASM file not created: %v", err)
	}
	if info.Size() < 1_000_000 {
		t.Errorf("WASM file suspiciously small: %d bytes", info.Size())
	}
	if info.Size() > 20_000_000 {
		t.Errorf("WASM file suspiciously large: %d bytes", info.Size())
	}
	t.Logf("WASM binary size: %.1f MB", float64(info.Size())/1_000_000)
}

func TestWASM_WasmExecExists(t *testing.T) {
	// Verify the Go WASM support file exists (needed by consumers)
	cmd := exec.Command("go", "env", "GOROOT")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("cannot get GOROOT: %v", err)
	}
	goroot := string(out[:len(out)-1]) // trim newline
	wasmExec := goroot + "/lib/wasm/wasm_exec.js"
	if _, err := os.Stat(wasmExec); err != nil {
		// Older Go versions have it in misc/wasm
		wasmExec = goroot + "/misc/wasm/wasm_exec.js"
		if _, err := os.Stat(wasmExec); err != nil {
			t.Errorf("wasm_exec.js not found at expected locations")
		}
	}
}
