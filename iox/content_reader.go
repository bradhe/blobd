package iox

import (
	"io"
	"io/ioutil"
)

type ContentReader interface {
	io.Reader

	// Indicates what type of content is inside this reader.
	ContentType() string
}

type ContentReadCloser interface {
	ContentReader
	io.Closer
}

type contentReader struct {
	contentType string
	r           io.ReadCloser
}

func (c *contentReader) ContentType() string {
	return c.contentType
}

func (c *contentReader) Read(buf []byte) (int, error) {
	return c.r.Read(buf)
}

func (c *contentReader) Close() error {
	return c.r.Close()
}

func MakeContentReader(contentType string, r io.Reader) ContentReader {
	return &contentReader{contentType, ioutil.NopCloser(r)}
}

func MakeContentReadCloser(contentType string, r io.ReadCloser) ContentReadCloser {
	return &contentReader{contentType, r}
}
