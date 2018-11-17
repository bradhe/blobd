package aws

import (
	"context"
	"net/url"

	"github.com/bradhe/blobd/storage/managers"
)

type Managers struct {
	ctx context.Context

	// region for putting content
	region string

	// bucket to read content from
	bucket string

	// prefix to apply to any keys in S3
	prefix string
}

func (m *Managers) Blobs() managers.BlobManager {
	return &BlobManager{
		ctx:    m.ctx,
		region: m.region,
		prefix: m.prefix,

		// There is a documented requirement that the table get it's name from the
		// bucket.
		table:  m.bucket,
		bucket: m.bucket,
	}
}

func (m *Managers) WithContext(ctx context.Context) managers.Managers {
	return &Managers{ctx, m.region, m.bucket, m.prefix}
}

func New(url *url.URL) managers.Managers {
	var prefix string

	bucket := url.Host

	if url.Path != "" {
		prefix = url.Path
	} else {
		prefix = "/"
	}

	// TODO: Parameterize region.
	return &Managers{nil, "us-west-2", bucket, prefix}
}
