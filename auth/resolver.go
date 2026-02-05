package auth

import (
	"context"

	"github.com/google/uuid"
)

// ResolvedProject is a simplified project structure that the resolver must return.
// It must contain the ID from the auth service.
type ResolvedProject struct {
	AuthServiceProjectID uuid.UUID
}

// ProjectResolver defines the interface required by the middleware to look up a project.
// Any service that uses the middleware must provide an implementation of this interface.
type ProjectResolver interface {
	GetProjectByID(ctx context.Context, id uuid.UUID) (*ResolvedProject, error)
}
