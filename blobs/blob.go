package blobs

import (
	"io"
)

type Blob struct {
	Id   Id
	Body io.Reader
}
