package auth

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
)

const bearerPrefix = "Bearer "

// NewClerkInterceptor returns a ConnectRPC interceptor that verifies Clerk
// session tokens on incoming requests. The verified Clerk user ID is stored in
// the request context and can be retrieved with UserIDFromContext.
func NewClerkInterceptor(verifier *ClerkVerifier) connect.Interceptor {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			auth := req.Header().Get("Authorization")
			if !strings.HasPrefix(auth, bearerPrefix) {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing or invalid authorization header"))
			}

			token := auth[len(bearerPrefix):]
			claims, err := verifier.Verify(ctx, token)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			ctx = WithUserID(ctx, claims.Subject)
			return next(ctx, req)
		}
	})
}
