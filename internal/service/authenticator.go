package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

type Authenticator struct {
}

// GenerateBearerToken creates a random 32-byte bearer token.
func (a *Authenticator) GenerateToken(token string) string {
	validTo := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	validToString := validTo.Format(time.RFC3339)
	s := fmt.Sprintf("%s&%s", token, validToString)
	b := []byte(s)
	return base64.StdEncoding.EncodeToString(b)
}

func (a *Authenticator) ValidateToken(token string) (string, error) {
	// Here you would typically check the token against a database or cache.
	// For simplicity, let's assume the token is valid if it starts with "Bearer ".
	if len(token) < 7 || token[:7] != "Bearer " {
		return "", fmt.Errorf("invalid token format")
	}

	hexStr := token[7:]
	b, err := base64.StdEncoding.DecodeString(hexStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode token: %w", err)
	}
	parts := strings.Split(string(b), "&")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid token data")
	}

	validTo, err := time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid time format: %w", err)
	}

	if validTo.Before(time.Now()) {
		return "", fmt.Errorf("token has expired")
	}

	return parts[0], nil
}
