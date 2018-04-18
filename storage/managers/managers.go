package managers

import (
	"context"

	"github.com/pborman/uuid"

	"github.com/bradhe/blobd/blobs"
)

type BlobManager interface {
	Get(uuid.UUID) (*blobs.Blob, error)
	Create(*blobs.Blob) error
	Update(*blobs.Blob) error
}

type Managers interface {
	Blobs() BlobManager
	WithContext(context.Context) Managers
}
