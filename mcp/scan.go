package mcp

import (
	"context"

	"github.com/teslashibe/mcptool"
	originality "github.com/teslashibe/originality-go"
)

// ScanTextInput is the typed input for originality_scan_text.
type ScanTextInput struct {
	Content      string         `json:"content" jsonschema:"description=text to scan; scores are API-provided signals, not authorship proof,required"`
	Title        string         `json:"title,omitempty" jsonschema:"description=optional title or document label"`
	AI           *bool          `json:"ai,omitempty" jsonschema:"description=include AI detection when supported by the API"`
	Plagiarism   *bool          `json:"plagiarism,omitempty" jsonschema:"description=include plagiarism/similarity checking when supported by the API"`
	Readability  *bool          `json:"readability,omitempty" jsonschema:"description=include readability checking when supported by the API"`
	Grammar      *bool          `json:"grammar,omitempty" jsonschema:"description=include grammar and spelling checking when supported by the API"`
	Factuality   *bool          `json:"factuality,omitempty" jsonschema:"description=include factual-claim checking when supported by the API"`
	Optimization *bool          `json:"optimization,omitempty" jsonschema:"description=include content optimization when supported by the API"`
	StoreScan    *bool          `json:"store_scan,omitempty" jsonschema:"description=ask Originality.ai to store the scan when supported"`
	ExcludedURL  []string       `json:"excluded_url,omitempty" jsonschema:"description=URLs to exclude from similarity matching when supported"`
	ExcludedText []string       `json:"excluded_text,omitempty" jsonschema:"description=text snippets to exclude from matching when supported"`
	Metadata     map[string]any `json:"metadata,omitempty" jsonschema:"description=optional caller metadata sent to Originality.ai"`
	Options      map[string]any `json:"options,omitempty" jsonschema:"description=provider-specific documented options"`
}

func scanText(ctx context.Context, c *originality.Client, in ScanTextInput) (any, error) {
	return c.ScanText(ctx, &originality.ScanTextRequest{
		Content:      in.Content,
		Title:        in.Title,
		AI:           in.AI,
		Plagiarism:   in.Plagiarism,
		Readability:  in.Readability,
		Grammar:      in.Grammar,
		Factuality:   in.Factuality,
		Optimization: in.Optimization,
		StoreScan:    in.StoreScan,
		ExcludedURL:  in.ExcludedURL,
		ExcludedText: in.ExcludedText,
		Metadata:     in.Metadata,
		Options:      in.Options,
	})
}

var scanTools = []mcptool.Tool{
	mcptool.Define[*originality.Client, ScanTextInput](
		"originality_scan_text",
		"Run an Originality.ai text scan and return structured provider signals",
		"ScanText",
		scanText,
	),
}
