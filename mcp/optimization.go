package mcp

import (
	"context"

	"github.com/teslashibe/mcptool"
	originality "github.com/teslashibe/originality-go"
)

func optimizeContent(ctx context.Context, c *originality.Client, in TextCheckInput) (any, error) {
	return c.OptimizeContent(ctx, textRequest(in))
}

var optimizationTools = []mcptool.Tool{
	mcptool.Define[*originality.Client, TextCheckInput](
		"originality_optimize_content",
		"Run Originality.ai content optimization checks for text",
		"OptimizeContent",
		optimizeContent,
	),
}
