package websocket

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthenticator implements JWT-based authentication
type JWTAuthenticator struct {
	secret     []byte
	issuer     string
	expiration time.Duration
}

// NewJWTAuthenticator tạo một JWT authenticator mới
func NewJWTAuthenticator(secret string, issuer string, expiration time.Duration) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:     []byte(secret),
		issuer:     issuer,
		expiration: expiration,
	}
}

// Authenticate xác thực request và trả về AuthInfo
func (a *JWTAuthenticator) Authenticate(req *http.Request) (*AuthInfo, error) {
	// Get token from Authorization header
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}
	
	// Extract token from "Bearer <token>" format
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	
	return a.ValidateToken(parts[1])
}

// ValidateToken xác thực JWT token
func (a *JWTAuthenticator) ValidateToken(tokenString string) (*AuthInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secret, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	
	// Validate issuer
	if iss, ok := claims["iss"].(string); ok && a.issuer != "" {
		if iss != a.issuer {
			return nil, fmt.Errorf("invalid issuer")
		}
	}
	
	// Extract user information
	authInfo := &AuthInfo{
		Token:  tokenString,
		Claims: make(map[string]interface{}),
	}
	
	// Copy all claims
	for key, value := range claims {
		authInfo.Claims[key] = value
	}
	
	// Extract standard claims
	if sub, ok := claims["sub"].(string); ok {
		authInfo.UserID = sub
	}
	
	if username, ok := claims["username"].(string); ok {
		authInfo.Username = username
	}
	
	if roles, ok := claims["roles"].([]interface{}); ok {
		authInfo.Roles = make([]string, len(roles))
		for i, role := range roles {
			if roleStr, ok := role.(string); ok {
				authInfo.Roles[i] = roleStr
			}
		}
	}
	
	return authInfo, nil
}

// RefreshToken tạo token mới từ token cũ
func (a *JWTAuthenticator) RefreshToken(tokenString string) (string, error) {
	// Validate current token
	authInfo, err := a.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token for refresh: %w", err)
	}
	
	// Create new token with same claims but new expiration
	claims := jwt.MapClaims{}
	for key, value := range authInfo.Claims {
		// Skip exp claim, we'll set a new one
		if key != "exp" {
			claims[key] = value
		}
	}
	
	// Set new expiration
	claims["exp"] = time.Now().Add(a.expiration).Unix()
	claims["iat"] = time.Now().Unix()
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret)
}

// RevokeToken thu hồi token (implementation depends on your token storage)
func (a *JWTAuthenticator) RevokeToken(tokenString string) error {
	// In a real implementation, you would add this token to a blacklist
	// stored in Redis, database, or in-memory cache
	return nil
}

// GenerateToken tạo token mới cho user
func (a *JWTAuthenticator) GenerateToken(userID, username string, roles []string, customClaims map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"roles":    roles,
		"iss":      a.issuer,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(a.expiration).Unix(),
	}
	
	// Add custom claims
	for key, value := range customClaims {
		claims[key] = value
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret)
}

// BasicAuthenticator implements basic authentication
type BasicAuthenticator struct {
	users map[string]string // username -> password
}

// NewBasicAuthenticator tạo một basic authenticator mới
func NewBasicAuthenticator(users map[string]string) *BasicAuthenticator {
	return &BasicAuthenticator{
		users: users,
	}
}

// Authenticate xác thực basic auth
func (a *BasicAuthenticator) Authenticate(req *http.Request) (*AuthInfo, error) {
	username, password, ok := req.BasicAuth()
	if !ok {
		return nil, fmt.Errorf("missing basic auth credentials")
	}
	
	expectedPassword, exists := a.users[username]
	if !exists || expectedPassword != password {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	return &AuthInfo{
		UserID:   username,
		Username: username,
		Roles:    []string{"user"},
	}, nil
}

// ValidateToken not implemented for basic auth
func (a *BasicAuthenticator) ValidateToken(token string) (*AuthInfo, error) {
	return nil, fmt.Errorf("token validation not supported for basic auth")
}

// RefreshToken not implemented for basic auth
func (a *BasicAuthenticator) RefreshToken(token string) (string, error) {
	return "", fmt.Errorf("token refresh not supported for basic auth")
}

// RevokeToken not implemented for basic auth
func (a *BasicAuthenticator) RevokeToken(token string) error {
	return fmt.Errorf("token revocation not supported for basic auth")
}

// CustomAuthenticator allows custom authentication logic
type CustomAuthenticator struct {
	authFunc func(*http.Request) (*AuthInfo, error)
}

// NewCustomAuthenticator tạo một custom authenticator
func NewCustomAuthenticator(authFunc func(*http.Request) (*AuthInfo, error)) *CustomAuthenticator {
	return &CustomAuthenticator{
		authFunc: authFunc,
	}
}

// Authenticate sử dụng custom auth function
func (a *CustomAuthenticator) Authenticate(req *http.Request) (*AuthInfo, error) {
	return a.authFunc(req)
}

// ValidateToken not implemented for custom auth
func (a *CustomAuthenticator) ValidateToken(token string) (*AuthInfo, error) {
	return nil, fmt.Errorf("token validation not supported for custom auth")
}

// RefreshToken not implemented for custom auth
func (a *CustomAuthenticator) RefreshToken(token string) (string, error) {
	return "", fmt.Errorf("token refresh not supported for custom auth")
}

// RevokeToken not implemented for custom auth
func (a *CustomAuthenticator) RevokeToken(token string) error {
	return fmt.Errorf("token revocation not supported for custom auth")
}

// NoAuthenticator allows all connections without authentication
type NoAuthenticator struct{}

// NewNoAuthenticator tạo một no-auth authenticator
func NewNoAuthenticator() *NoAuthenticator {
	return &NoAuthenticator{}
}

// Authenticate always returns a default auth info
func (a *NoAuthenticator) Authenticate(req *http.Request) (*AuthInfo, error) {
	return &AuthInfo{
		UserID:   "anonymous",
		Username: "anonymous",
		Roles:    []string{"anonymous"},
	}, nil
}

// ValidateToken not implemented
func (a *NoAuthenticator) ValidateToken(token string) (*AuthInfo, error) {
	return nil, fmt.Errorf("token validation not supported for no-auth")
}

// RefreshToken not implemented
func (a *NoAuthenticator) RefreshToken(token string) (string, error) {
	return "", fmt.Errorf("token refresh not supported for no-auth")
}

// RevokeToken not implemented
func (a *NoAuthenticator) RevokeToken(token string) error {
	return fmt.Errorf("token revocation not supported for no-auth")
}

// Helper functions for creating auth handlers
func CreateJWTAuthHandler(secret, issuer string, expiration time.Duration) func(*http.Request) (*AuthInfo, error) {
	auth := NewJWTAuthenticator(secret, issuer, expiration)
	return auth.Authenticate
}

func CreateBasicAuthHandler(users map[string]string) func(*http.Request) (*AuthInfo, error) {
	auth := NewBasicAuthenticator(users)
	return auth.Authenticate
}

func CreateNoAuthHandler() func(*http.Request) (*AuthInfo, error) {
	auth := NewNoAuthenticator()
	return auth.Authenticate
}
