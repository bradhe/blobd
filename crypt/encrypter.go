package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type Encrypter struct {
	closed bool

	dst    io.Reader
	src    io.Reader
	stream io.WriteCloser
}

func (w *Encrypter) Close() error {
	if !w.closed {
		return w.stream.Close()
	}

	return nil
}

func (w *Encrypter) writeMore() error {
	n, err := io.CopyN(w.stream, w.src, 14)

	if err != nil {
		return err
	}

	// If we didn't copy everything we needed...
	if n == 0 {
		return io.EOF
	}

	return nil
}

func (w *Encrypter) Read(dst []byte) (int, error) {
	var cum int

	for {
		n, err := w.dst.Read(dst[cum:])

		if err == io.EOF {
			// try to write more. if we make it to the end of the buffer, we've done
			// all we can here.
			if err := w.writeMore(); err != nil {
				// If this is also an EOF then there's no more work to be done at all.
				if err == io.EOF {
					goto ReadDone
				} else {
					return 0, err
				}
			}
		} else {
			cum += n

			if n == 0 || cum == len(dst) {
				goto ReadDone
			}
		}
	}

ReadDone:
	if cum == 0 {
		return 0, io.EOF
	}

	return cum, nil
}

func NewEncrypter(key *Key, r io.Reader) io.ReadCloser {
	block, _ := aes.NewCipher(key.Bytes())

	var buf bytes.Buffer

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])
	writer := &cipher.StreamWriter{S: stream, W: &buf}

	return &Encrypter{
		dst:    &buf,
		src:    r,
		stream: writer,
	}
}
