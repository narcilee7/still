package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
)

// contextKey is the type used for context keys in this package.
type contextKey string

const userIDKey contextKey = "clerkUserID"

// ClerkVerifier verifies Clerk session tokens.
type ClerkVerifier struct {
	jwksClient *jwks.Client
}

// NewClerkVerifier creates a verifier using the given Clerk secret key.
func NewClerkVerifier(secretKey string) *ClerkVerifier {
	cfg := &clerk.ClientConfig{
		BackendConfig: clerk.BackendConfig{Key: clerk.String(secretKey)},
	}
	return &ClerkVerifier{jwksClient: jwks.NewClient(cfg)}
}

// Verify validates a Clerk session token and returns its claims.
func (v *ClerkVerifier) Verify(ctx context.Context, bearerToken string) (*clerk.SessionClaims, error) {
	unsafeClaims, err := jwt.Decode(ctx, &jwt.DecodeParams{Token: bearerToken})
	if err != nil {
		return nil, fmt.Errorf("decode token: %w", err)
	}

	jwk, err := jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
		KeyID:      unsafeClaims.KeyID,
		JWKSClient: v.jwksClient,
	})
	if err != nil {
		return nil, fmt.Errorf("get jwk: %w", err)
	}

	claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
		Token: bearerToken,
		JWK:   jwk,
	})
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	return claims, nil
}

// WithUserID returns a context with the given Clerk user ID attached.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext extracts the Clerk user ID from the context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok && strings.TrimSpace(id) != ""
}
