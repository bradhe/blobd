// +build unit

package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bradhe/blobd/blobs"
	"github.com/bradhe/blobd/crypt"
)

func TestValidJWT(t *testing.T) {
	blob := &blobs.Blob{
		Id:        blobs.NewId(),
		ExpiresAt: blobs.DefaultExpirationFromNow(),
	}

	key, err := crypt.NewRandomKey()
	require.NoError(t, err)

	jwt := GenerateJWT(WritableToken, key, blob)

	// We shouldn't be able to parse this JWT...
	_, _, err = ParseJWT(jwt)
	assert.NoError(t, err)
}

func TestExpiredJWT(t *testing.T) {
	blob := &blobs.Blob{
		Id:        blobs.NewId(),
		ExpiresAt: time.Now().UTC().Add(-30 * 24 * time.Hour),
	}

	key, err := crypt.NewRandomKey()
	require.NoError(t, err)

	jwt := GenerateJWT(WritableToken, key, blob)

	// We shouldn't be able to parse this JWT...
	_, _, err = ParseJWT(jwt)
	assert.Error(t, err)
}
