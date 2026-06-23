// Package originality provides a small, stdlib-first Go client for the
// Originality.ai REST API.
//
// The client is designed for agent/tooling hosts: it accepts an API key from
// ORIGINALITY_API_KEY or an explicit constructor option, exposes typed request
// and response models for content-quality checks, and keeps submitted text
// confined to the outbound API request.
package originality
