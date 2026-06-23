package mcp

import (
	"context"

	"github.com/teslashibe/mcptool"
	originality "github.com/teslashibe/originality-go"
)

// AccountCreditsInput is the typed input for originality_account_credits.
type AccountCreditsInput struct{}

func accountCredits(ctx context.Context, c *originality.Client, _ AccountCreditsInput) (any, error) {
	return c.AccountCredits(ctx)
}

var accountTools = []mcptool.Tool{
	mcptool.Define[*originality.Client, AccountCreditsInput](
		"originality_account_credits",
		"Fetch Originality.ai account credit and balance information",
		"AccountCredits",
		accountCredits,
	),
}
