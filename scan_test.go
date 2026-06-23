package originality

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCapabilityMethodsUseScanEndpointWithFlags(t *testing.T) {
	paths := map[string]bool{}
	var bodies []ScanTextRequest
	httpClient := roundTripClient(func(r *http.Request) (*http.Response, error) {
		paths[r.URL.Path] = true
		var body ScanTextRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		bodies = append(bodies, body)
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
	}
	for _, call := range calls {
		if err := call(); err != nil {
			t.Fatal(err)
		}
	}
	if len(paths) != 1 || !paths["/scan/ai"] {
		t.Fatalf("capability checks should use /scan/ai only; got %#v", paths)
	}
	if len(bodies) != 5 {
		t.Fatalf("bodies length = %d, want 5", len(bodies))
	}
	if bodies[0].Plagiarism == nil || !*bodies[0].Plagiarism {
		t.Fatalf("plagiarism flag not set: %+v", bodies[0])
	}
	if bodies[1].Readability == nil || !*bodies[1].Readability {
		t.Fatalf("readability flag not set: %+v", bodies[1])
	}
	if bodies[2].Grammar == nil || !*bodies[2].Grammar {
		t.Fatalf("grammar flag not set: %+v", bodies[2])
	}
	if bodies[3].Factuality == nil || !*bodies[3].Factuality {
		t.Fatalf("factuality flag not set: %+v", bodies[3])
	}
	if bodies[4].Optimization == nil || !*bodies[4].Optimization {
		t.Fatalf("optimization flag not set: %+v", bodies[4])
	}
}
