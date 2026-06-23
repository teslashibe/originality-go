package mcp

import (
	"context"

	"github.com/teslashibe/mcptool"
	originality "github.com/teslashibe/originality-go"
)

// NoInput is used by read-only account tools with no arguments.
type NoInput struct{}

func creditBalance(ctx context.Context, c *originality.Client, _ NoInput) (any, error) {
	return c.CreditBalance(ctx)
}

func creditUsage(ctx context.Context, c *originality.Client, _ NoInput) (any, error) {
	return c.CreditUsage(ctx)
}

func payments(ctx context.Context, c *originality.Client, _ NoInput) (any, error) {
	return c.Payments(ctx)
}

var accountTools = []mcptool.Tool{
	mcptool.Define[*originality.Client, NoInput](
		"originality_credit_balance",
		"Fetch the current Originality.ai credit balance",
		"CreditBalance",
		creditBalance,
	),
	mcptool.Define[*originality.Client, NoInput](
		"originality_credit_usage",
		"Fetch recent Originality.ai scan credit usage",
		"CreditUsage",
		creditUsage,
	),
	mcptool.Define[*originality.Client, NoInput](
		"originality_payments",
		"Fetch recent Originality.ai credit purchases",
		"Payments",
		payments,
	),
}
