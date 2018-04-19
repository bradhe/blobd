package managers

import (
	"context"

	"github.com/bradhe/blobd/blobs"
)

type BlobManager interface {
	Get(blobs.Id) (*blobs.Blob, error)
	Create(*blobs.Blob) error
	Update(*blobs.Blob) error
}

type Managers interface {
	Blobs() BlobManager
	WithContext(context.Context) Managers
}
