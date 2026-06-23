package mcp

import (
	"context"

	"github.com/teslashibe/mcptool"
	originality "github.com/teslashibe/originality-go"
)

func checkPlagiarism(ctx context.Context, c *originality.Client, in TextCheckInput) (any, error) {
	return c.CheckPlagiarism(ctx, textRequest(in))
}

var plagiarismTools = []mcptool.Tool{
	mcptool.Define[*originality.Client, TextCheckInput](
		"originality_check_plagiarism",
		"Check text for Originality.ai plagiarism and similarity signals",
		"CheckPlagiarism",
		checkPlagiarism,
	),
}
