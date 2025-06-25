package obs_studio

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// AuthManager handles OBS WebSocket authentication
type AuthManager struct{}

// NewAuthManager creates a new authentication manager
func NewAuthManager() *AuthManager {
	return &AuthManager{}
}

// GenerateAuthString generates authentication string from password and challenge data
// This implements the same authentication mechanism as obs-websocket-js
func (a *AuthManager) GenerateAuthString(password string, challenge *AuthChallenge) (string, error) {
	if challenge == nil {
		return "", fmt.Errorf("challenge data is required for authentication")
	}

	if challenge.Challenge == "" || challenge.Salt == "" {
		return "", fmt.Errorf("invalid challenge data: challenge and salt are required")
	}

	// Step 1: Generate secret string
	// secret = base64_encode(sha256(password + salt))
	passwordSalt := password + challenge.Salt
	passwordSaltHash := sha256.Sum256([]byte(passwordSalt))
	secret := base64.StdEncoding.EncodeToString(passwordSaltHash[:])

	// Step 2: Generate authentication string
	// auth = base64_encode(sha256(secret + challenge))
	secretChallenge := secret + challenge.Challenge
	secretChallengeHash := sha256.Sum256([]byte(secretChallenge))
	authString := base64.StdEncoding.EncodeToString(secretChallengeHash[:])

	return authString, nil
}

// ValidateAuthChallenge validates if the provided challenge data is valid
func (a *AuthManager) ValidateAuthChallenge(challenge *AuthChallenge) error {
	if challenge == nil {
		return fmt.Errorf("challenge data cannot be nil")
	}

	if challenge.Challenge == "" {
		return fmt.Errorf("challenge cannot be empty")
	}

	if challenge.Salt == "" {
		return fmt.Errorf("salt cannot be empty")
	}

	// Validate base64 encoding of challenge and salt
	if _, err := base64.StdEncoding.DecodeString(challenge.Challenge); err != nil {
		return fmt.Errorf("invalid challenge format: must be base64 encoded")
	}

	if _, err := base64.StdEncoding.DecodeString(challenge.Salt); err != nil {
		return fmt.Errorf("invalid salt format: must be base64 encoded")
	}

	return nil
}

// RequiresAuthentication checks if authentication is required based on hello data
func (a *AuthManager) RequiresAuthentication(helloData *HelloData) bool {
	return helloData != nil && helloData.Authentication != nil
}
