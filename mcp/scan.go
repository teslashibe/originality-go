package mcp

import (
	"context"

	"github.com/teslashibe/mcptool"
	originality "github.com/teslashibe/originality-go"
)

// ScanTextInput is the typed input for originality_scan_text.
type ScanTextInput struct {
	Content string `json:"content" jsonschema:"description=text to scan; scores are API-provided signals, not authorship proof,required"`
	Title   string `json:"title,omitempty" jsonschema:"description=optional title or document label"`
}

// ScanIDInput is the typed input for originality_get_scan.
type ScanIDInput struct {
	ScanID string `json:"scan_id" jsonschema:"description=Originality.ai scan identifier,required"`
}

func scanText(ctx context.Context, c *originality.Client, in ScanTextInput) (any, error) {
	return c.ScanText(ctx, &originality.ScanTextRequest{
		Content: in.Content,
		Title:   in.Title,
	})
}

func getScan(ctx context.Context, c *originality.Client, in ScanIDInput) (any, error) {
	return c.GetScan(ctx, in.ScanID)
}

var scanTools = []mcptool.Tool{
	mcptool.Define[*originality.Client, ScanTextInput](
		"originality_scan_text",
		"Run an Originality.ai text scan and return structured provider signals",
		"ScanText",
		scanText,
	),
	mcptool.Define[*originality.Client, ScanIDInput](
		"originality_get_scan",
		"Fetch an Originality.ai scan result by scan ID",
		"GetScan",
		getScan,
	),
}
