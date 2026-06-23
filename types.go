package originality

import "encoding/json"

// ScanTextRequest starts an Originality.ai AI-detection text scan.
type ScanTextRequest struct {
	Content string `json:"content"`
	Title   string `json:"title,omitempty"`
}

// ScanResponse is a flexible response model for scan creation. Fields map to
// common Originality.ai response concepts while Raw preserves provider fields
// for agent consumers.
type ScanResponse struct {
	ID          string         `json:"id,omitempty"`
	ScanID      string         `json:"scan_id,omitempty"`
	Status      string         `json:"status,omitempty"`
	Score       float64        `json:"score,omitempty"`
	CreditsUsed float64        `json:"credits_used,omitempty"`
	CreatedAt   string         `json:"created_at,omitempty"`
	UpdatedAt   string         `json:"updated_at,omitempty"`
	Raw         map[string]any `json:"raw,omitempty"`
}

type responseEnvelope struct {
	Data    json.RawMessage `json:"data"`
	Result  json.RawMessage `json:"result"`
	Message string          `json:"message"`
	Error   string          `json:"error"`
}
