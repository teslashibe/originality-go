package originality

import "encoding/json"

// ScanTextRequest starts an Originality.ai AI-detection text scan.
type ScanTextRequest struct {
	Content string `json:"content"`
	Title   string `json:"title,omitempty"`
}

// ScanURLRequest starts an Originality.ai AI-detection scan for a webpage.
type ScanURLRequest struct {
	URL string `json:"url"`
}

// ScanResponse is a flexible response model for text scan creation/read
// operations. Raw preserves provider fields that are not yet promoted.
type ScanResponse struct {
	Success     bool           `json:"success,omitempty"`
	ID          string         `json:"id,omitempty"`
	ScanID      string         `json:"scan_id,omitempty"`
	Status      string         `json:"status,omitempty"`
	Score       float64        `json:"score,omitempty"`
	ScoreDetail *ScoreDetail   `json:"score_detail,omitempty"`
	CreditsUsed float64        `json:"credits_used,omitempty"`
	Credits     float64        `json:"credits,omitempty"`
	Content     string         `json:"content,omitempty"`
	CreatedAt   string         `json:"created_at,omitempty"`
	UpdatedAt   string         `json:"updated_at,omitempty"`
	Raw         map[string]any `json:"raw,omitempty"`
}

// URLScanResponse is the documented URL AI-detection response shape with Raw
// preserving nested result fields.
type URLScanResponse struct {
	Success        bool             `json:"success,omitempty"`
	URL            string           `json:"url,omitempty"`
	URLCode        int              `json:"url_code,omitempty"`
	CreditsUsed    float64          `json:"credits_used,omitempty"`
	Credits        float64          `json:"credits,omitempty"`
	WordCount      int              `json:"word_count,omitempty"`
	Score          *ScoreDetail     `json:"score,omitempty"`
	ScoreBreakdown []ScoreBreakdown `json:"score_breakdown,omitempty"`
	Results        map[string]any   `json:"results,omitempty"`
	Raw            map[string]any   `json:"raw,omitempty"`
}

// ScoreDetail contains API-provided original/AI percentages.
type ScoreDetail struct {
	Original float64 `json:"original,omitempty"`
	AI       float64 `json:"ai,omitempty"`
}

// ScoreBreakdown describes one text block score.
type ScoreBreakdown struct {
	Original float64 `json:"original,omitempty"`
	AI       float64 `json:"ai,omitempty"`
	Text     string  `json:"text,omitempty"`
}

// CreditBalanceResponse contains account balance information.
type CreditBalanceResponse struct {
	Balance int            `json:"balance,omitempty"`
	Raw     map[string]any `json:"raw,omitempty"`
}

// CreditUsageResponse contains recent scan credit usage.
type CreditUsageResponse struct {
	Usage []CreditUsage  `json:"usage,omitempty"`
	Raw   map[string]any `json:"raw,omitempty"`
}

// CreditUsage is one credit usage record.
type CreditUsage struct {
	ContentID   string  `json:"contentID,omitempty"`
	CreditsUsed float64 `json:"credits_used,omitempty"`
	Date        string  `json:"date,omitempty"`
}

// PaymentsResponse contains recent credit purchases.
type PaymentsResponse struct {
	Payments []Payment      `json:"payments,omitempty"`
	Raw      map[string]any `json:"raw,omitempty"`
}

// Payment is one credit purchase record.
type Payment struct {
	Credits float64 `json:"credits,omitempty"`
	Price   string  `json:"price,omitempty"`
	Receipt string  `json:"receipt,omitempty"`
	Date    string  `json:"date,omitempty"`
}

type responseEnvelope struct {
	Data    json.RawMessage `json:"data"`
	Result  json.RawMessage `json:"result"`
	Message string          `json:"message"`
	Error   string          `json:"error"`
}
