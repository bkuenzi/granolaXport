// Package api provides functionality for interacting with the Granola API.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrWrapperJSON         = errors.New("failed to unmarshal wrapper JSON")
	ErrTokensJSON          = errors.New("failed to unmarshal token JSON")
	ErrAccessTokenNotFound = errors.New("access token not found")
)

// Wrapper holds the data from Granola's supabase.json file.
type Wrapper struct {
	Tokens string `json:"workos_tokens"`
}

// Tokens holds the access token and related information.
type Tokens struct {
	AccessToken string `json:"access_token"`
}

// getAccessToken takes the JSON from supabase.json and returns the Granola access token.
func getAccessToken(file []byte) (string, error) {
	var wrapper Wrapper
	if err := json.Unmarshal(file, &wrapper); err != nil {
		return "", fmt.Errorf("%w: %s", ErrWrapperJSON, err)
	}

	var tokens Tokens
	if err := json.Unmarshal([]byte(wrapper.Tokens), &tokens); err != nil {
		return "", fmt.Errorf("%w: %s", ErrTokensJSON, err)
	}

	if strings.TrimSpace(tokens.AccessToken) == "" {
		return "", ErrAccessTokenNotFound
	}

	// Reject tokens that don't look like JWTs to prevent header injection.
	if parts := strings.Split(tokens.AccessToken, "."); len(parts) != 3 {
		return "", fmt.Errorf("access token has unexpected format (expected 3 JWT segments, got %d)", len(parts))
	}

	return tokens.AccessToken, nil
}
