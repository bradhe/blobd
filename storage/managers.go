package storage

import (
	"net/url"

	"github.com/bradhe/blobd/storage/aws"
	"github.com/bradhe/blobd/storage/inmem"
	"github.com/bradhe/blobd/storage/managers"
)

func New(url *url.URL) managers.Managers {
	if url == nil {
		goto Managers_Inmem
	}

	switch url.Scheme {
	case "s3":
		return aws.New(url)
	}

Managers_Inmem:
	return inmem.New()
}
