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

func RenderJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(Dump(v))
}

func RenderRedirectedJSON(w http.ResponseWriter, path string, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", path)
	w.WriteHeader(http.StatusFound)
	w.Write(Dump(v))
}

func DumpString(v interface{}) string {
	return string(Dump(v))
}

func RenderError(w http.ResponseWriter, err Error) {
	w.WriteHeader(err.Status)
	RenderJSON(w, err)
}

func BlobPath(blob *blobs.Blob) string {
	return "/" + blob.Id.String()
}

func requestContentType(req *http.Request) string {
	if t := req.Header.Get("Content-Type"); t != "" {
		return t
	}

	return "application/octet-stream"
}
