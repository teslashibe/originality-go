package mcp

import (
	"context"

	"github.com/teslashibe/mcptool"
	originality "github.com/teslashibe/originality-go"
)

func checkGrammar(ctx context.Context, c *originality.Client, in TextCheckInput) (any, error) {
	return c.CheckGrammar(ctx, textRequest(in))
}

var grammarTools = []mcptool.Tool{
	mcptool.Define[*originality.Client, TextCheckInput](
		"originality_check_grammar",
		"Check text for Originality.ai grammar and spelling signals",
		"CheckGrammar",
		checkGrammar,
	),
}
