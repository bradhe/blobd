package server

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/bradhe/blobd/iox"
)

type multipartReader struct {
	rdr  *multipart.Reader
	part *multipart.Part
}

func (r *multipartReader) nextPart() error {
	if r.part == nil {
		part, err := r.rdr.NextPart()

		if err != nil {
			log.WithError(err).Debugf("failed to read next part in multipart stream")
			return err
		}

		log.Debugf("got multipart %v", part.Header)
		r.part = part
	}

	return nil
}

func (r *multipartReader) Read(buf []byte) (int, error) {
	r.nextPart()
	return r.part.Read(buf)
}

func (r *multipartReader) ContentType() string {
	r.nextPart()
	return r.part.Header.Get("Content-Type")
}

func newMultipartReader(mediatype string, body io.Reader) iox.ContentReader {
	// If there was a problem, there's nothing to do here.
	if _, props, err := mime.ParseMediaType(mediatype); err != nil {
		log.WithError(err).Debugf("creating multipart reader failed")
		return iox.MakeContentReader(mediatype, body)
	} else if boundary, ok := props["boundary"]; ok {
		return &multipartReader{multipart.NewReader(body, boundary), nil}
	}

	return iox.MakeContentReader(mediatype, body)
}

func isMultipart(str string) bool {
	if mediatype, _, err := mime.ParseMediaType(str); err != nil {
		return false
	} else {
		return mediatype == "multipart/form-data"
	}
}

func newRequestReader(r *http.Request) iox.ContentReader {
	mediatype := requestContentType(r)

	if isMultipart(mediatype) {
		return newMultipartReader(mediatype, r.Body)
	} else {
		return iox.MakeContentReader(mediatype, r.Body)
	}
}
