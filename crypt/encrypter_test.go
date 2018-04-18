// +build unit

package crypt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncrypter(t *testing.T) {
	// we use a large message to make sure that we can encrypt and decrypt
	// sufficiently.
	msg := randomMessage(t, 2048)

	key, err := NewRandomKey()
	require.NoError(t, err)

	enc := NewEncrypter(key, bytes.NewBuffer(msg))

	// We'll just assume that if we can write and the buffer content doesn't
	// equal that the message is encrypted.
	ciphertext := make([]byte, 10)

	n, err := enc.Read(ciphertext)
	require.NoError(t, err)
	assert.Equal(t, n, 10)
	assert.NotEqual(t, ciphertext, msg[:10])
}
