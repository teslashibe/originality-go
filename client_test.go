package originality

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewUsesEnvAPIKey(t *testing.T) {
	t.Setenv("ORIGINALITY_API_KEY", "env-key")
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if c.apiKey != "env-key" {
		t.Fatalf("apiKey = %q, want env-key", c.apiKey)
	}
}

func TestNewRequiresAPIKey(t *testing.T) {
	t.Setenv("ORIGINALITY_API_KEY", "")
	_, err := New()
	if !errors.Is(err, ErrMissingAPIKey) {
		t.Fatalf("New error = %v, want ErrMissingAPIKey", err)
	}
}

func TestScanTextRequestConstructionAndDecode(t *testing.T) {
	var gotPath, gotMethod, gotKey string
	var gotBody ScanTextRequest
	httpClient := roundTripClient(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotKey = r.Header.Get("X-OAI-API-KEY")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}
		return jsonResponse(http.StatusOK, `{"scan_id":"scan_123","status":"completed","credits_used":1.25}`), nil
	})

	c, err := New(WithAPIKey("test-key"), WithBaseURL("https://originality.test"), WithHTTPClient(httpClient))
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.ScanText(context.Background(), &ScanTextRequest{Content: "hello", Title: "doc"})
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPost || gotPath != "/scan/ai" {
		t.Fatalf("request = %s %s, want POST /scan/ai", gotMethod, gotPath)
	}
	if gotKey != "test-key" {
		t.Fatalf("auth header = %q, want test-key", gotKey)
	}
	if gotBody.Content != "hello" || gotBody.Title != "doc" {
		t.Fatalf("body = %+v", gotBody)
	}
	if res.ScanID != "scan_123" || res.Status != "completed" || res.CreditsUsed != 1.25 {
		t.Fatalf("response = %+v", res)
	}
	if res.Raw["scan_id"] != "scan_123" {
		t.Fatalf("raw response missing scan_id: %#v", res.Raw)
	}
}

func TestEnvelopeDecode(t *testing.T) {
	c, err := New(
		WithAPIKey("test-key"),
		WithBaseURL("https://originality.test"),
		WithHTTPClient(roundTripClient(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, `{"data":{"scan_id":"wrapped","status":"queued"}}`), nil
		})),
	)
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.GetScan(context.Background(), "wrapped")
	if err != nil {
		t.Fatal(err)
	}
	if res.ScanID != "wrapped" || res.Status != "queued" {
		t.Fatalf("response = %+v", res)
	}
}

func TestTransportOverrides(t *testing.T) {
	t.Setenv("ORIGINALITY_API_KEY", "")
	httpClient := &http.Client{Timeout: 250 * time.Millisecond}
	c, err := New(WithAPIKey("test-key"), WithBaseURL("https://example.test/root/"), WithHTTPClient(httpClient))
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL != "https://example.test/root" {
		t.Fatalf("baseURL = %q", c.baseURL)
	}
	if c.httpClient != httpClient {
		t.Fatal("WithHTTPClient did not install injected client")
	}

	c, err = New(WithAPIKey("test-key"), WithTimeout(123*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	if c.httpClient.Timeout != 123*time.Millisecond {
		t.Fatalf("timeout = %s", c.httpClient.Timeout)
	}
}

func TestErrorMapping(t *testing.T) {
	cases := []struct {
		status int
		want   error
	}{
		{http.StatusBadRequest, ErrBadRequest},
		{http.StatusUnauthorized, ErrUnauthorized},
		{http.StatusPaymentRequired, ErrInsufficientCredit},
		{http.StatusTooManyRequests, ErrRateLimited},
		{http.StatusInternalServerError, ErrServerFailure},
	}
	for _, tc := range cases {
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			c, err := New(
				WithAPIKey("test-key"),
				WithBaseURL("https://originality.test"),
				WithHTTPClient(roundTripClient(func(_ *http.Request) (*http.Response, error) {
					resp := jsonResponse(tc.status, `{"message":"sanitized provider message"}`)
					resp.Header.Set("Retry-After", "2")
					return resp, nil
				})),
			)
			if err != nil {
				t.Fatal(err)
			}
			_, err = c.AccountCredits(context.Background())
			if !errors.Is(err, tc.want) {
				t.Fatalf("error = %v, want %v", err, tc.want)
			}
			var httpErr *HTTPError
			if !errors.As(err, &httpErr) {
				t.Fatalf("error type = %T, want *HTTPError", err)
			}
			if tc.status == http.StatusTooManyRequests && httpErr.RetryAfter != 2*time.Second {
				t.Fatalf("RetryAfter = %s", httpErr.RetryAfter)
			}
		})
	}
}

func TestValidation(t *testing.T) {
	c, err := New(WithAPIKey("test-key"), WithHTTPClient(http.DefaultClient))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.ScanText(context.Background(), &ScanTextRequest{}); !errors.Is(err, ErrBadRequest) {
		t.Fatalf("ScanText error = %v, want ErrBadRequest", err)
	}
	if _, err := c.GetScan(context.Background(), ""); !errors.Is(err, ErrBadRequest) {
		t.Fatalf("GetScan error = %v, want ErrBadRequest", err)
	}
}

func TestMain(m *testing.M) {
	_ = os.Unsetenv("ORIGINALITY_API_KEY")
	os.Exit(m.Run())
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
