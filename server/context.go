package server

import (
	"context"

	"github.com/bradhe/blobd/crypt"

	"github.com/pborman/uuid"
)

type ContextKey int

const (
	ctxDecryptionKey ContextKey = iota
	ctxBlobId
)

func WithBlobId(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxBlobId, id)
}

func BlobId(ctx context.Context) uuid.UUID {
	return ctx.Value(ctxBlobId).(uuid.UUID)
}

func WithDecryptionKey(ctx context.Context, key *crypt.Key) context.Context {
	return context.WithValue(ctx, ctxDecryptionKey, key)
}

func DecryptionKey(ctx context.Context) *crypt.Key {
	return ctx.Value(ctxDecryptionKey).(*crypt.Key)
}
