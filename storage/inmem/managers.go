package inmem

import (
	"bytes"
	"context"

	"github.com/bradhe/blobd/storage/managers"
)

func init() {
	// Initialize our blob store when this guy comes online.
	blobstore = make(map[string]bytes.Buffer)
}

type Managers struct {
	ctx context.Context
}

func (m *Managers) Blobs() managers.BlobManager {
	return &BlobManager{
		ctx: m.ctx,
	}
}

func (m *Managers) WithContext(ctx context.Context) managers.Managers {
	return &Managers{ctx}
}

func New() managers.Managers {
	return &Managers{}
}
