package originality_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
