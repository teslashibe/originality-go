package mcp

// Excluded enumerates exported methods on *originality.Client that are
// intentionally not exposed via MCP. Each entry must have a non-empty reason.
//
// The coverage test in mcp_test.go fails if any exported method on *Client is
// neither wrapped by a Tool nor present in this map.
var Excluded = map[string]string{}
