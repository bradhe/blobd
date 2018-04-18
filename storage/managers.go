package storage

import (
	"github.com/bradhe/blobd/storage/inmem"
	"github.com/bradhe/blobd/storage/managers"
	"net/url"
)

func New(url *url.URL) managers.Managers {
	// TODO: Implement backend based on storage URL provided.
	return inmem.New()
}
