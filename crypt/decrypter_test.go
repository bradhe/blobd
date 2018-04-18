// +build unit

package crypt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecrypter(t *testing.T) {
	// we use a large message to make sure that we can encrypt and decrypt
	// sufficiently.
	msg := randomMessage(t, 2048)

	key, err := NewRandomKey()
	require.NoError(t, err)

	ciphertext := mustEncypt(t, key, msg)

	dec := NewDecrypter(key, bytes.NewBuffer(ciphertext))

	// read only the first few bytes to ensure we can validly stream this content
	// out.
	prefix := make([]byte, 10)

	n, err := dec.Read(prefix)
	require.NoError(t, err)
	assert.Equal(t, n, 10)
	assert.Equal(t, prefix, msg[:10])

	// make sure we can actually read the rest.
	suffix := make([]byte, 2038)

	n, err = dec.Read(suffix)
	require.NoError(t, err)
	assert.Equal(t, n, 2038)
	assert.Equal(t, suffix, msg[10:])
}
