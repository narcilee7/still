package auth

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
)

// UserReader resolves a Clerk-authenticated identity to an internal user.
type UserReader interface {
	GetOrCreateUserByClerkID(ctx context.Context, clerkUserID, username string) (*stillv1.User, error)
}

// CurrentUser returns the internal user for the Clerk identity stored in ctx.
// It returns an Unauthenticated error when no identity is present.
func CurrentUser(ctx context.Context, reader UserReader) (*stillv1.User, error) {
	clerkUserID, ok := UserIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("unauthenticated"))
	}
	user, err := reader.GetOrCreateUserByClerkID(ctx, clerkUserID, "")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("resolve user failed: %w", err))
	}
	return user, nil
}
