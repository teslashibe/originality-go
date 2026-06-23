package originality_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	originality "github.com/teslashibe/originality-go"
)

func ExampleClient_ScanText() {
	client, _ := originality.New(
		originality.WithAPIKey("test-key"),
		originality.WithBaseURL("https://originality.test"),
		originality.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"scan_id":"scan_123","status":"completed"}`)),
			}, nil
		})}),
	)

	res, _ := client.ScanText(context.Background(), &originality.ScanTextRequest{
		Content: "Text to check with Originality.ai.",
		Title:   "example",
	})
	fmt.Printf("scan_id=%s status=%s\n", res.ScanID, res.Status)

	// Output:
	// scan_id=scan_123 status=completed
}

func ExampleNewFromEnv() {
	_ = os.Setenv("ORIGINALITY_API_KEY", "test-key")
	defer os.Unsetenv("ORIGINALITY_API_KEY")

	client, _ := originality.NewFromEnv(
		originality.WithBaseURL("https://originality.test"),
		originality.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"scan_id":"scan_123","status":"completed","credits_used":1.25}`)),
			}, nil
		})}),
	)

	res, _ := client.ScanText(context.Background(), &originality.ScanTextRequest{
		Content: "Text to check with Originality.ai.",
	})
	fmt.Printf("scan_id=%s status=%s credits=%.2f\n", res.ScanID, res.Status, res.CreditsUsed)

	// Output:
	// scan_id=scan_123 status=completed credits=1.25
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
