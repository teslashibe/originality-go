package originality

import (
	"context"
	"net/http"
	"testing"
)

func TestCapabilityMethodsUseDocumentedPaths(t *testing.T) {
	paths := map[string]bool{}
	httpClient := roundTripClient(func(r *http.Request) (*http.Response, error) {
		paths[r.URL.Path] = true
		return jsonResponse(http.StatusOK, `{"scan_id":"ok","status":"completed"}`), nil
	})

	c, err := New(WithAPIKey("test-key"), WithBaseURL("https://originality.test"), WithHTTPClient(httpClient))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	req := &TextCheckRequest{Content: "sample text"}
	calls := []func() error{
		func() error { _, err := c.CheckPlagiarism(ctx, req); return err },
		func() error { _, err := c.CheckReadability(ctx, req); return err },
		func() error { _, err := c.CheckGrammar(ctx, req); return err },
		func() error { _, err := c.CheckFactuality(ctx, req); return err },
		func() error { _, err := c.OptimizeContent(ctx, req); return err },
		func() error { _, err := c.ScanStatus(ctx, "scan_1"); return err },
		func() error { _, err := c.AccountCredits(ctx); return err },
	}
	for _, call := range calls {
		if err := call(); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{
		"/scan/plagiarism",
		"/scan/readability",
		"/scan/grammar",
		"/scan/factuality",
		"/scan/optimization",
		"/scan/scan_1/status",
		"/account/credits",
	} {
		if !paths[path] {
			t.Fatalf("path %s was not requested; got %#v", path, paths)
		}
	}
}
