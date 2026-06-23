//go:build live

package originality

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestLiveScanTextSmoke(t *testing.T) {
	if os.Getenv("ORIGINALITY_API_KEY") == "" {
		t.Skip("ORIGINALITY_API_KEY is not set")
	}
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	res, err := c.ScanText(ctx, &ScanTextRequest{
		Content: "This is a short smoke-test paragraph for the originality-go client.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.ScanID == "" && res.ID == "" && res.Status == "" {
		t.Fatalf("live response did not include an id or status: %+v", res)
	}
	t.Logf("scan_id=%q id=%q status=%q credits_used=%.2f", res.ScanID, res.ID, res.Status, res.CreditsUsed)
}
