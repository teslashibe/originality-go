# originality-go

Go client and MCP tools for the Originality.ai API. The root package is a
stdlib-first REST client; the `mcp` package exposes the same useful client
operations as typed MCP tools for agent hosts.

```bash
go get github.com/teslashibe/originality-go
```

## Authentication

Create an API key in Originality.ai and export it before constructing the
client:

```bash
export ORIGINALITY_API_KEY="oai_..."
```

You can also pass the key explicitly:

```go
client, err := originality.New(originality.WithAPIKey("oai_..."))
```

The client sends the key as `X-OAI-API-KEY`. It never logs or returns the key.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    originality "github.com/teslashibe/originality-go"
)

func main() {
    client, err := originality.NewFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    res, err := client.ScanText(context.Background(), &originality.ScanTextRequest{
        Content: "Text to check with Originality.ai.",
        Title:   "example",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("scan_id=%s status=%s credits=%.2f\n", res.ScanID, res.Status, res.CreditsUsed)
}
```

## API Surface

The client models production text-quality flows from the official
Originality.ai REST API:

```go
client.ScanText(ctx, req)          // combined text scan / AI detection flow
client.CheckPlagiarism(ctx, req)   // plagiarism/similarity signals
client.CheckReadability(ctx, req)  // readability signals
client.CheckGrammar(ctx, req)      // grammar and spelling signals
client.CheckFactuality(ctx, req)   // factual-claim signals
client.OptimizeContent(ctx, req)   // content optimization signals
```

The capability-specific helpers call the documented text scan route with the
corresponding option flag. They do not invent separate provider endpoints for
plagiarism, readability, grammar, factuality, optimization, account credits, or
scan polling. Add those as separate exported methods only after verifying the
official API exposes stable routes for them.

Detector and quality scores are API-provided signals. Do not present them as
definitive authorship proof.

## Configuration

```go
client, err := originality.New(
    originality.WithAPIKey(os.Getenv("ORIGINALITY_API_KEY")),
    originality.WithTimeout(30*time.Second),
    originality.WithBaseURL("https://api.originality.ai/api/v1"),
    originality.WithHTTPClient(customHTTPClient),
)
```

`WithBaseURL` and `WithHTTPClient` are intended for tests, private hosts, and
production transport customization. Normal tests use `httptest` and do not
require live credentials.

## Error Handling

HTTP failures wrap sentinel errors so callers can branch with `errors.Is`:

```go
if errors.Is(err, originality.ErrMissingAPIKey)      { /* configure credentials */ }
if errors.Is(err, originality.ErrBadRequest)         { /* fix input */ }
if errors.Is(err, originality.ErrUnauthorized)       { /* rotate API key */ }
if errors.Is(err, originality.ErrRateLimited)        { /* back off */ }
if errors.Is(err, originality.ErrInsufficientCredit) { /* add credits */ }
if errors.Is(err, originality.ErrServerFailure)      { /* retry later */ }
if errors.Is(err, originality.ErrRequestFailed)      { /* transport/decode failure */ }
```

For non-2xx responses, `*originality.HTTPError` includes the status code,
sanitized provider message, and `RetryAfter` when supplied. Submitted text and
API keys are not included in error values.

## MCP Support

The `mcp` package follows the Teslashibe `x-go` / `reddit-go` pattern:
zero-value provider, grouped tool files, `mcptool.Define` typed schemas, and a
coverage test that fails if exported client methods are not wrapped or
explicitly excluded.

```go
import (
    originality "github.com/teslashibe/originality-go"
    originalitymcp "github.com/teslashibe/originality-go/mcp"
)

client, _ := originality.NewFromEnv()
provider := originalitymcp.Provider{}
for _, tool := range provider.Tools() {
    // Register tool with your MCP server and pass client when invoking.
}
```

Tools use the `originality_` prefix:

- `originality_scan_text`
- `originality_check_plagiarism`
- `originality_check_readability`
- `originality_check_grammar`
- `originality_check_factuality`
- `originality_optimize_content`

## Credits And Cost

Originality.ai operations can consume account credits. Inspect response fields
such as `credits_used` where returned by the API, and run large batches behind
your own rate and budget controls.

## Testing

Mocked tests do not require live credentials:

```bash
go test ./...
```

An opt-in smoke test can be run when a real key is available:

```bash
export ORIGINALITY_API_KEY="oai_..."
go test -tags live -run TestLiveScanTextSmoke -v
```
