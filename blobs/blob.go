package blobs

import (
	"time"

	"github.com/bradhe/blobd/iox"
)

// Default window of time that an object is valid in the context of our body.
var DefaultExpirationDuration = 7 * 24 * time.Hour

func DefaultExpiration(base time.Time) time.Time {
	return base.Add(DefaultExpirationDuration)
}

func DefaultExpirationFromNow() time.Time {
	return DefaultExpiration(time.Now().UTC())
}

type Blob struct {
	Id        Id
	Body      iox.ContentReader
	MediaType []byte
	ExpiresAt time.Time
}

func (b Blob) Expiration() time.Time {
	if b.ExpiresAt.IsZero() {
		return DefaultExpirationFromNow()
	}

	return b.ExpiresAt
}
