package mcp

import (
	"context"

	"github.com/teslashibe/mcptool"
	originality "github.com/teslashibe/originality-go"
)

func checkReadability(ctx context.Context, c *originality.Client, in TextCheckInput) (any, error) {
	return c.CheckReadability(ctx, textRequest(in))
}

var readabilityTools = []mcptool.Tool{
	mcptool.Define[*originality.Client, TextCheckInput](
		"originality_check_readability",
		"Check text for Originality.ai readability signals",
		"CheckReadability",
		checkReadability,
	),
}
