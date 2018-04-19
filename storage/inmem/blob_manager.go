package inmem

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/bradhe/blobd/blobs"
	"github.com/bradhe/blobd/storage/managers"
)

var (
	// TODO: can't remember if you can concurrently read from a golang maps?
	blobmut   sync.RWMutex
	blobstore map[string]bytes.Buffer
)

type BlobManager struct {
	ctx context.Context
}

func (bm *BlobManager) Get(id blobs.Id) (*blobs.Blob, error) {
	blobmut.RLock()
	defer blobmut.RUnlock()

	if buf, ok := blobstore[id.String()]; ok {
		return &blobs.Blob{
			Id:   id,
			Body: &buf,
		}, nil
	}

	return nil, managers.ErrNotFound
}

func (bm *BlobManager) Create(blob *blobs.Blob) error {
	blobmut.Lock()
	defer blobmut.Unlock()

	// TODO: Validate this doesn't already exist.
	blob.Id = blobs.NewId()

	var buf bytes.Buffer
	io.Copy(&buf, blob.Body)

	blobstore[blob.Id.String()] = buf

	return nil
}

func (bm *BlobManager) Update(blob *blobs.Blob) error {
	blobmut.Lock()
	defer blobmut.Unlock()

	var buf bytes.Buffer
	io.Copy(&buf, blob.Body)

	blobstore[blob.Id.String()] = buf

	return nil
}
