package originality

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrMissingAPIKey      = errors.New("originality: missing API key")
	ErrBadRequest         = errors.New("originality: bad request")
	ErrUnauthorized       = errors.New("originality: unauthorized")
	ErrRateLimited        = errors.New("originality: rate limited")
	ErrInsufficientCredit = errors.New("originality: insufficient credits")
	ErrServerFailure      = errors.New("originality: server failure")
	ErrRequestFailed      = errors.New("originality: request failed")
)

// HTTPError carries sanitized non-2xx response details.
type HTTPError struct {
	StatusCode int
	Message    string
	Err        error
	RetryAfter time.Duration
}

func (e *HTTPError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		return fmt.Sprintf("%s: HTTP %d", e.Err, e.StatusCode)
	}
	return fmt.Sprintf("%s: HTTP %d: %s", e.Err, e.StatusCode, e.Message)
}

func (e *HTTPError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func errorForStatus(status int, message string, retryAfter time.Duration) error {
	base := ErrRequestFailed
	switch {
	case status == http.StatusBadRequest || status == http.StatusUnprocessableEntity:
		base = ErrBadRequest
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		base = ErrUnauthorized
	case status == http.StatusPaymentRequired:
		base = ErrInsufficientCredit
	case status == http.StatusTooManyRequests:
		base = ErrRateLimited
	case status >= 500:
		base = ErrServerFailure
	}
	return &HTTPError{
		StatusCode: status,
		Message:    message,
		Err:        base,
		RetryAfter: retryAfter,
	}
}

func parseRetryAfter(v string) time.Duration {
	if v == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(v); err == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second
	}
	if when, err := http.ParseTime(v); err == nil {
		if wait := time.Until(when); wait > 0 {
			return wait
		}
	}
	return 0
}
