package mcp

import originality "github.com/teslashibe/originality-go"

// TextCheckInput is the typed input for capability-specific text checks.
type TextCheckInput struct {
	Content  string         `json:"content" jsonschema:"description=text to check,required"`
	Title    string         `json:"title,omitempty" jsonschema:"description=optional title or document label"`
	Metadata map[string]any `json:"metadata,omitempty" jsonschema:"description=optional caller metadata sent to Originality.ai"`
	Options  map[string]any `json:"options,omitempty" jsonschema:"description=provider-specific documented options"`
}

func textRequest(in TextCheckInput) *originality.TextCheckRequest {
	return &originality.TextCheckRequest{
		Content:  in.Content,
		Title:    in.Title,
		Metadata: in.Metadata,
		Options:  in.Options,
	}
}
