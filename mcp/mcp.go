// Package mcp exposes the originality-go client as MCP tools.
//
// All tools wrap exported methods on *originality.Client through
// mcptool.Define so JSON input schemas are derived from typed structs.
package mcp

import "github.com/teslashibe/mcptool"

// Provider implements mcptool.Provider for Originality.ai. The zero value is
// ready to use.
type Provider struct{}

// Platform returns "originality".
func (Provider) Platform() string { return "originality" }

// Tools returns every Originality.ai MCP tool in registration order.
func (Provider) Tools() []mcptool.Tool {
	out := make([]mcptool.Tool, 0,
		len(scanTools)+
			len(plagiarismTools)+
			len(readabilityTools)+
			len(grammarTools)+
			len(factualityTools)+
			len(optimizationTools),
	)
	out = append(out, scanTools...)
	out = append(out, plagiarismTools...)
	out = append(out, readabilityTools...)
	out = append(out, grammarTools...)
	out = append(out, factualityTools...)
	out = append(out, optimizationTools...)
	return out
}
