package mcp_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/teslashibe/mcptool"
	originality "github.com/teslashibe/originality-go"
	originalitymcp "github.com/teslashibe/originality-go/mcp"
)

func TestEveryClientMethodIsWrappedOrExcluded(t *testing.T) {
	rep := mcptool.Coverage(
		reflect.TypeOf(&originality.Client{}),
		(originalitymcp.Provider{}).Tools(),
		originalitymcp.Excluded,
	)
	if len(rep.Missing) > 0 {
		t.Fatalf("methods missing MCP exposure (add a tool or list in excluded.go): %v", rep.Missing)
	}
	if len(rep.UnknownExclusions) > 0 {
		t.Fatalf("excluded.go references methods that don't exist on *Client: %v", rep.UnknownExclusions)
	}
	if len(rep.Wrapped)+len(rep.Excluded) == 0 {
		t.Fatal("no wrapped or excluded methods detected")
	}
}

func TestToolsValidate(t *testing.T) {
	if err := mcptool.ValidateTools((originalitymcp.Provider{}).Tools()); err != nil {
		t.Fatal(err)
	}
}

func TestPlatformName(t *testing.T) {
	if got := (originalitymcp.Provider{}).Platform(); got != "originality" {
		t.Errorf("Platform() = %q, want originality", got)
	}
}

func TestToolsHaveOriginalityPrefix(t *testing.T) {
	for _, tool := range (originalitymcp.Provider{}).Tools() {
		if !strings.HasPrefix(tool.Name, "originality_") {
			t.Errorf("tool %q lacks originality_ prefix", tool.Name)
		}
	}
}

func TestRepresentativeToolInvocations(t *testing.T) {
	requested := map[string]bool{}
	httpClient := roundTripClient(func(r *http.Request) (*http.Response, error) {
		requested[r.URL.Path] = true
		return jsonResponse(http.StatusOK, `{"scan_id":"scan_1","status":"completed"}`), nil
	})

	client, err := originality.New(
		originality.WithAPIKey("test-key"),
		originality.WithBaseURL("https://originality.test"),
		originality.WithHTTPClient(httpClient),
	)
	if err != nil {
		t.Fatal(err)
	}
	tools := map[string]mcptool.Tool{}
	for _, tool := range (originalitymcp.Provider{}).Tools() {
		tools[tool.Name] = tool
	}

	invoke := func(name string, input any) {
		t.Helper()
		raw, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := tools[name].Invoke(context.Background(), client, raw); err != nil {
			t.Fatalf("%s Invoke error = %v", name, err)
		}
	}

	invoke("originality_scan_text", map[string]any{"content": "sample text"})
	invoke("originality_get_scan", map[string]any{"scan_id": "scan_1"})

	if len(requested) != 2 || !requested["/scan"] || !requested["/scan/scan_1"] {
		t.Fatalf("scan tools should use documented scan paths; got %#v", requested)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func roundTripClient(f roundTripFunc) *http.Client {
	return &http.Client{Transport: f}
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
