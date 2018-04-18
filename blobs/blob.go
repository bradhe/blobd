package blobs

import (
	"io"

	"github.com/pborman/uuid"
)

type Blob struct {
	Id   uuid.UUID
	Body io.Reader
}
