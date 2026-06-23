package originality

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.originality.ai"
	defaultTimeout = 30 * time.Second
	maxBodyBytes   = 10 << 20
)

const (
	creditBalancePath = "/api/v1/account/credits/balance"
	creditUsagePath   = "/api/v1/account/credits/content_scan_usage"
	paymentsPath      = "/api/v1/account/credits/payments"
	scanTextPath      = "/api/v1/scan/ai"
	scanURLPath       = "/api/v1/scan/url"
	scanResultPath    = "/api/v3/scan/%s"
)

// Option configures a Client.
type Option func(*Client)

// Client is an Originality.ai REST client. It is safe for concurrent use when
// the injected HTTP client is safe for concurrent use.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// New constructs a client. If WithAPIKey is not supplied, ORIGINALITY_API_KEY
// is read from the environment.
func New(opts ...Option) (*Client, error) {
	c := &Client{
		apiKey:  strings.TrimSpace(os.Getenv("ORIGINALITY_API_KEY")),
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	c.baseURL = strings.TrimRight(c.baseURL, "/")
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: defaultTimeout}
	}
	return c, nil
}

// NewFromEnv constructs a client using ORIGINALITY_API_KEY.
func NewFromEnv(opts ...Option) (*Client, error) {
	return New(opts...)
}

// WithAPIKey sets the Originality.ai API key explicitly.
func WithAPIKey(apiKey string) Option {
	return func(c *Client) {
		c.apiKey = strings.TrimSpace(apiKey)
	}
}

// WithBaseURL overrides the API base URL for tests or private deployments.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if strings.TrimSpace(baseURL) != "" {
			c.baseURL = strings.TrimSpace(baseURL)
		}
	}
}

// WithHTTPClient injects the HTTP client used for requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithTimeout sets the timeout on the configured HTTP client.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout > 0 {
			if c.httpClient == nil {
				c.httpClient = &http.Client{Timeout: timeout}
				return
			}
			clone := *c.httpClient
			clone.Timeout = timeout
			c.httpClient = &clone
		}
	}
}

// ScanText starts an Originality.ai text scan.
func (c *Client) ScanText(ctx context.Context, req *ScanTextRequest) (*ScanResponse, error) {
	if req == nil || strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("%w: content is required", ErrBadRequest)
	}
	var out ScanResponse
	if err := c.postJSON(ctx, scanTextPath, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ScanURL starts an Originality.ai URL AI-detection scan.
func (c *Client) ScanURL(ctx context.Context, req *ScanURLRequest) (*URLScanResponse, error) {
	if req == nil || strings.TrimSpace(req.URL) == "" {
		return nil, fmt.Errorf("%w: url is required", ErrBadRequest)
	}
	var out URLScanResponse
	if err := c.postJSON(ctx, scanURLPath, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetScan returns a scan result by ID.
func (c *Client) GetScan(ctx context.Context, scanID string) (*ScanResponse, error) {
	if strings.TrimSpace(scanID) == "" {
		return nil, fmt.Errorf("%w: scan id is required", ErrBadRequest)
	}
	var out ScanResponse
	if err := c.getJSON(ctx, fmt.Sprintf(scanResultPath, url.PathEscape(scanID)), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreditBalance returns the current credit balance.
func (c *Client) CreditBalance(ctx context.Context) (*CreditBalanceResponse, error) {
	var out CreditBalanceResponse
	if err := c.getJSON(ctx, creditBalancePath, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreditUsage returns credit usage for recent scans.
func (c *Client) CreditUsage(ctx context.Context) (*CreditUsageResponse, error) {
	var out CreditUsageResponse
	if err := c.getJSON(ctx, creditUsagePath, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Payments returns recent credit purchases.
func (c *Client) Payments(ctx context.Context) (*PaymentsResponse, error) {
	var out PaymentsResponse
	if err := c.getJSON(ctx, paymentsPath, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) getJSON(ctx context.Context, path string, out any) error {
	return c.doJSON(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) postJSON(ctx context.Context, path string, in any, out any) error {
	return c.doJSON(ctx, http.MethodPost, path, in, out)
}

func (c *Client) doJSON(ctx context.Context, method, path string, in any, out any) error {
	if c == nil {
		return fmt.Errorf("%w: nil client", ErrRequestFailed)
	}
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("%w: encoding request: %v", ErrBadRequest, err)
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("%w: building request: %v", ErrRequestFailed, err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-OAI-API-KEY", c.apiKey)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return fmt.Errorf("%w: reading response: %v", ErrRequestFailed, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errorForStatus(resp.StatusCode, responseMessage(raw), parseRetryAfter(resp.Header.Get("Retry-After")))
	}
	if len(raw) == 0 || out == nil {
		return nil
	}
	return decodeResponse(raw, out)
}

func decodeResponse(raw []byte, out any) error {
	var env responseEnvelope
	if err := json.Unmarshal(raw, &env); err == nil {
		for _, candidate := range []json.RawMessage{env.Data, env.Result} {
			if len(candidate) == 0 || string(candidate) == "null" {
				continue
			}
			if err := json.Unmarshal(candidate, out); err != nil {
				return fmt.Errorf("%w: decoding response data: %v", ErrRequestFailed, err)
			}
			attachRaw(candidate, out)
			return nil
		}
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("%w: decoding response: %v", ErrRequestFailed, err)
	}
	attachRaw(raw, out)
	return nil
}

func responseMessage(raw []byte) string {
	var env responseEnvelope
	if err := json.Unmarshal(raw, &env); err == nil {
		if env.Error != "" {
			return truncate(env.Error, 300)
		}
		if env.Message != "" {
			return truncate(env.Message, 300)
		}
	}
	return truncate(string(raw), 300)
}

func attachRaw(raw []byte, out any) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return
	}
	switch v := out.(type) {
	case *ScanResponse:
		v.Raw = m
	case *URLScanResponse:
		v.Raw = m
	case *CreditBalanceResponse:
		v.Raw = m
	case *CreditUsageResponse:
		v.Raw = m
	case *PaymentsResponse:
		v.Raw = m
	}
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
