package originality

import "encoding/json"

// TextCheckRequest is the shared input shape for single-text quality checks.
type TextCheckRequest struct {
	Content  string         `json:"content"`
	Title    string         `json:"title,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Options  map[string]any `json:"options,omitempty"`
}

// ScanTextRequest starts a combined Originality.ai text scan.
type ScanTextRequest struct {
	Content      string         `json:"content"`
	Title        string         `json:"title,omitempty"`
	AI           *bool          `json:"ai,omitempty"`
	Plagiarism   *bool          `json:"plagiarism,omitempty"`
	Readability  *bool          `json:"readability,omitempty"`
	Grammar      *bool          `json:"grammar,omitempty"`
	Factuality   *bool          `json:"factuality,omitempty"`
	Optimization *bool          `json:"optimization,omitempty"`
	StoreScan    *bool          `json:"store_scan,omitempty"`
	ExcludedURL  []string       `json:"excluded_url,omitempty"`
	ExcludedText []string       `json:"excluded_text,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	Options      map[string]any `json:"options,omitempty"`
}

// ScanResponse is a flexible response model for scan creation and scan reads.
// Fields map to common Originality.ai response concepts while Raw preserves
// operation-specific provider fields for agent consumers.
type ScanResponse struct {
	ID           string            `json:"id,omitempty"`
	ScanID       string            `json:"scan_id,omitempty"`
	Status       string            `json:"status,omitempty"`
	AI           *AIResult         `json:"ai,omitempty"`
	Plagiarism   *PlagiarismResult `json:"plagiarism,omitempty"`
	Readability  map[string]any    `json:"readability,omitempty"`
	Grammar      map[string]any    `json:"grammar,omitempty"`
	Factuality   map[string]any    `json:"factuality,omitempty"`
	Optimization map[string]any    `json:"optimization,omitempty"`
	CreditsUsed  float64           `json:"credits_used,omitempty"`
	CreatedAt    string            `json:"created_at,omitempty"`
	UpdatedAt    string            `json:"updated_at,omitempty"`
	Raw          map[string]any    `json:"raw,omitempty"`
}

// AIResult contains API-provided AI-detection signals, not authorship proof.
type AIResult struct {
	Score      float64        `json:"score,omitempty"`
	AI         float64        `json:"ai,omitempty"`
	Original   float64        `json:"original,omitempty"`
	Confidence string         `json:"confidence,omitempty"`
	Label      string         `json:"label,omitempty"`
	Raw        map[string]any `json:"raw,omitempty"`
}

// PlagiarismResult contains API-provided similarity signals and source hits.
type PlagiarismResult struct {
	Score   float64           `json:"score,omitempty"`
	Percent float64           `json:"percent,omitempty"`
	Matches []PlagiarismMatch `json:"matches,omitempty"`
	Raw     map[string]any    `json:"raw,omitempty"`
}

// PlagiarismMatch is a single matching source returned by Originality.ai.
type PlagiarismMatch struct {
	URL         string  `json:"url,omitempty"`
	Title       string  `json:"title,omitempty"`
	Similarity  float64 `json:"similarity,omitempty"`
	MatchedText string  `json:"matched_text,omitempty"`
}

type responseEnvelope struct {
	Data    json.RawMessage `json:"data"`
	Result  json.RawMessage `json:"result"`
	Message string          `json:"message"`
	Error   string          `json:"error"`
}
