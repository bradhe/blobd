package server

import (
	"encoding/json"
	"net/http"

	"github.com/bradhe/blobd/blobs"
)

func Dump(v interface{}) []byte {
	if b, err := json.MarshalIndent(v, "", "  "); err != nil {
		log.Printf("server: dumping object failed: %v", err)
		return []byte("")
	} else {
		return b
	}
}

func DumpString(v interface{}) string {
	return string(Dump(v))
}

func RenderError(w http.ResponseWriter, err Error) {
	http.Error(w, DumpString(err), err.Status)
}

func BlobPath(blob *blobs.Blob) string {
	return "/" + blob.Id.String()
}
