package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type Decrypter struct {
	// Base reader to get content from for decrypting
	base io.Reader
}

func (d *Decrypter) Read(dst []byte) (int, error) {
	return d.base.Read(dst)
}

func (d *Decrypter) Close() error {
	// NOOP
	return nil
}

func NewDecrypter(key *Key, r io.Reader) io.ReadCloser {
	block, _ := aes.NewCipher(key.Bytes())

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])
	d := &cipher.StreamReader{S: stream, R: r}

	return &Decrypter{
		base: d,
	}
}
