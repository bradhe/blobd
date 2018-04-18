package crypt

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func mustEncypt(t *testing.T, key *Key, msg []byte) []byte {
	enc := NewEncrypter(key, bytes.NewBuffer(msg))

	buf, err := ioutil.ReadAll(enc)
	require.NoError(t, err)

	return buf
}

func randomMessage(t *testing.T, size int) []byte {
	msg := make([]byte, size)
	n, err := rand.Read(msg)
	require.NoError(t, err)
	require.Equal(t, n, size)
	return msg
}
