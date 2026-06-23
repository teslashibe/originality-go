package mcp

import (
	"context"

	"github.com/teslashibe/mcptool"
	originality "github.com/teslashibe/originality-go"
)

func checkFactuality(ctx context.Context, c *originality.Client, in TextCheckInput) (any, error) {
	return c.CheckFactuality(ctx, textRequest(in))
}

var factualityTools = []mcptool.Tool{
	mcptool.Define[*originality.Client, TextCheckInput](
		"originality_check_factuality",
		"Check text for Originality.ai factual-claim signals",
		"CheckFactuality",
		checkFactuality,
	),
}
