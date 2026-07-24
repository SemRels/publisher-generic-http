package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestIntegrationPublishRetriesRateLimitWithExactRequestContract(t *testing.T) {
	t.Parallel()

	artifact := filepath.Join(t.TempDir(), "release.tgz")
	if err := os.WriteFile(artifact, []byte("release-bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	var mu sync.Mutex
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
		}
		if r.Method != http.MethodPut || r.URL.Path != "/releases/1.2.3/release.tgz" {
			t.Errorf("request = %s %s, want PUT /releases/1.2.3/release.tgz", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer integration-secret" {
			t.Errorf("authorization = %q", got)
		}
		if got := r.Header.Get("X-Release-Source"); got != "semrel" {
			t.Errorf("X-Release-Source = %q", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/octet-stream" {
			t.Errorf("content type = %q", got)
		}
		if !bytes.Equal(body, []byte("release-bytes")) {
			t.Errorf("body = %q", body)
		}

		mu.Lock()
		attempts++
		attempt := attempts
		mu.Unlock()
		if attempt == 1 {
			w.Header().Set("Retry-After", "0")
			http.Error(w, "slow down", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(server.Close)

	var stdout, stderr bytes.Buffer
	code := run(&stdout, &stderr, env(map[string]string{
		"SEMREL_VERSION":             "v1.2.3",
		"SEMREL_PLUGIN_URL":          server.URL + "/releases/{version}/{artifact}",
		"SEMREL_PLUGIN_ARTIFACT":     artifact,
		"SEMREL_PLUGIN_TOKEN":        "integration-secret",
		"SEMREL_PLUGIN_HEADERS_JSON": `{"X-Release-Source":"semrel"}`,
	}))

	if code != 0 {
		t.Fatalf("run code = %d, stderr = %s", code, stderr.String())
	}
	mu.Lock()
	gotAttempts := attempts
	mu.Unlock()
	if gotAttempts != 2 {
		t.Fatalf("attempts = %d, want 2", gotAttempts)
	}
	if !strings.Contains(stdout.String(), "published 1 artifact(s)") {
		t.Errorf("stdout = %q", stdout.String())
	}
	if strings.Contains(stdout.String(), "integration-secret") || strings.Contains(stderr.String(), "integration-secret") {
		t.Error("publisher leaked its bearer token")
	}
}

func TestIntegrationDryRunAndInvalidConfigurationMakeNoRequest(t *testing.T) {
	t.Parallel()

	artifact := filepath.Join(t.TempDir(), "release.tgz")
	if err := os.WriteFile(artifact, []byte("release-bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		t.Error("dry run made an HTTP request")
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(server.Close)

	var stdout, stderr bytes.Buffer
	code := run(&stdout, &stderr, env(map[string]string{
		"SEMREL_VERSION":         "1.2.3",
		"SEMREL_PLUGIN_URL":      server.URL + "/releases/{artifact}",
		"SEMREL_PLUGIN_ARTIFACT": artifact,
		"SEMREL_DRY_RUN":         "true",
	}))
	if code != 0 || requests != 0 {
		t.Fatalf("dry run code = %d, requests = %d, stderr = %s", code, requests, stderr.String())
	}
	if !strings.Contains(stdout.String(), "[dry-run]") {
		t.Errorf("stdout = %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = run(&stdout, &stderr, env(map[string]string{
		"SEMREL_VERSION":         "1.2.3",
		"SEMREL_PLUGIN_ARTIFACT": artifact,
		"SEMREL_PLUGIN_TOKEN":    "integration-secret",
	}))
	if code != 1 || !strings.Contains(stderr.String(), "SEMREL_PLUGIN_URL is required") {
		t.Fatalf("invalid configuration: code = %d, stderr = %q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "integration-secret") {
		t.Error("invalid configuration leaked its bearer token")
	}
}

func TestIntegrationExternalErrorDoesNotRetryOrLeakToken(t *testing.T) {
	t.Parallel()

	artifact := filepath.Join(t.TempDir(), "release.tgz")
	if err := os.WriteFile(artifact, []byte("release-bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "upstream failed", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	var stdout, stderr bytes.Buffer
	code := run(&stdout, &stderr, env(map[string]string{
		"SEMREL_VERSION":         "1.2.3",
		"SEMREL_PLUGIN_URL":      server.URL + "/releases/{artifact}",
		"SEMREL_PLUGIN_ARTIFACT": artifact,
		"SEMREL_PLUGIN_TOKEN":    "integration-secret",
	}))
	if code != 1 || requests != 1 || !strings.Contains(stderr.String(), "status 500") {
		t.Fatalf("code = %d, requests = %d, stderr = %q", code, requests, stderr.String())
	}
	if strings.Contains(stdout.String(), "integration-secret") || strings.Contains(stderr.String(), "integration-secret") {
		t.Error("external error leaked its bearer token")
	}
}

func env(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}
