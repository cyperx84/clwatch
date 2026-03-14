package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var testBinaryPath string

func buildBinary(t *testing.T) string {
	t.Helper()
	if testBinaryPath != "" {
		if _, err := os.Stat(testBinaryPath); err == nil {
			return testBinaryPath
		}
	}
	// Build in a persistent temp dir
	tmpDir := os.TempDir()
	bin := filepath.Join(tmpDir, "clwatch-test-bin")
	// Find go.mod to determine build dir
	buildDir := "."
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		buildDir = "../.."
	}
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/clwatch")
	cmd.Dir = buildDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	testBinaryPath = bin
	return bin
}

func startTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/refs/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"schema":"clwatch.manifest.v1","generated":"2026-03-12T00:00:00Z","tools":{"claude-code":{"version":"2.1.74","stale_after":"2026-03-19T00:00:00Z"}}}`))
	})
	mux.HandleFunc("/api/refs/status.json", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"schema":"clwatch.status.v1","generated_at":"2026-03-12T00:00:00Z","pipeline":{"last_run_at":"2026-03-12T00:00:00Z","status":"ok","tools_checked":5},"tools":{"claude-code":{"version":"2.1.74","verification_status":"verified","stale":false}}}`))
	})
	return httptest.NewServer(mux)
}

func TestVersionCommand(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "version").CombinedOutput()
	if err != nil {
		t.Fatalf("version failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "clwatch") {
		t.Errorf("expected 'clwatch' in output, got: %s", out)
	}
}

func TestInvalidCommand(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "nonexistent")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected error for invalid command")
	}
	if !strings.Contains(string(out), "unknown command") {
		t.Errorf("expected 'unknown' in error, got: %s", out)
	}
}

func TestDiffCommand(t *testing.T) {
	bin := buildBinary(t)
	srv := startTestServer(t)
	defer srv.Close()

	cmd := exec.Command(bin, "diff")
	cmd.Env = append(os.Environ(), "CLWATCH_MANIFEST_URL="+srv.URL+"/api/refs/manifest.json")
	out, _ := cmd.CombinedOutput()
	if strings.Contains(string(out), "panic") {
		t.Fatalf("diff panicked: %s", out)
	}
}

func TestListCommand(t *testing.T) {
	bin := buildBinary(t)
	srv := startTestServer(t)
	defer srv.Close()

	cmd := exec.Command(bin, "list")
	cmd.Env = append(os.Environ(), "CLWATCH_MANIFEST_URL="+srv.URL+"/api/refs/manifest.json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("list failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "claude-code") {
		t.Errorf("expected tool list, got: %s", out)
	}
}

func TestListJSON(t *testing.T) {
	bin := buildBinary(t)
	srv := startTestServer(t)
	defer srv.Close()

	cmd := exec.Command(bin, "list", "--json")
	cmd.Env = append(os.Environ(), "CLWATCH_MANIFEST_URL="+srv.URL+"/api/refs/manifest.json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("list --json failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "claude-code") {
		t.Errorf("expected JSON with tools, got: %s", out)
	}
}

func TestStatusCommand(t *testing.T) {
	bin := buildBinary(t)
	srv := startTestServer(t)
	defer srv.Close()

	cmd := exec.Command(bin, "status")
	cmd.Env = append(os.Environ(), "CLWATCH_STATUS_URL="+srv.URL+"/api/refs/status.json")
	out, _ := cmd.CombinedOutput()
	if strings.Contains(string(out), "panic") {
		t.Fatalf("status panicked: %s", out)
	}
}
