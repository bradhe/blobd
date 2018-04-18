package server

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/pborman/uuid"

	"github.com/bradhe/blobd/blobs"
	"github.com/bradhe/blobd/crypt"
	"github.com/bradhe/blobd/storage/managers"
)

type Handler struct {
	// The context this request was started in.
	Context context.Context

	// Any URL vars that were passed in to the path.
	Vars map[string]string

	// Managers to use in the current context.
	Managers managers.Managers
}

type PostBlobResponse struct {
	JWT string `json:"jwt"`
}

func (h Handler) PostBlob(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	key, err := crypt.NewRandomKey()

	if err != nil {
		RenderError(w, GetError("internal_server_error"))
	} else {
		manager := h.Managers.Blobs()

		// TODO: Validate this is JSON on the way in.
		blob := blobs.Blob{
			Body: crypt.NewEncrypter(key, r.Body),
		}

		if err := manager.Create(&blob); err != nil {
			RenderError(w, GetError("internal_server_error"))
		} else {
			resp := PostBlobResponse{
				JWT: GenerateJWT(key, &blob),
			}

			http.Redirect(w, r, BlobPath(&blob), http.StatusFound)
			w.Write(Dump(resp))
		}
	}
}

func (h Handler) GetBlob(w http.ResponseWriter, r *http.Request) {
	manager := h.Managers.Blobs()

	// we'll use this to authenticate that we're getting the correct blob
	blobId := BlobId(r.Context())

	if id := uuid.Parse(h.Vars["blob_id"]); uuid.Equal(id, uuid.NIL) {
		log.Printf("server: invalid uuid %v", id)
		RenderError(w, GetError("internal_server_error"))
	} else {
		if !uuid.Equal(blobId, id) {
			RenderError(w, GetError("unauthorized"))
		} else {
			if blob, err := manager.Get(id); err != nil {
				log.Printf("server: error happened %v", err)
				RenderError(w, GetError("internal_server_error"))
			} else {
				io.Copy(w, crypt.NewDecrypter(DecryptionKey(r.Context()), blob.Body))
			}
		}
	}
}

type PutBlobResponse struct {
	JWT string `json:"jwt"`
}

func (h Handler) PutBlob(w http.ResponseWriter, r *http.Request) {
	manager := h.Managers.Blobs()

	// we'll use this to authenticate that we're getting the correct blob
	blobId := BlobId(r.Context())
	key := DecryptionKey(r.Context())

	if id := uuid.Parse(h.Vars["blob_id"]); uuid.Equal(id, uuid.NIL) {
		RenderError(w, GetError("internal_server_error"))
	} else {
		if !uuid.Equal(blobId, id) {
			RenderError(w, GetError("unauthorized"))
		} else {
			blob := blobs.Blob{
				Id:   id,
				Body: crypt.NewEncrypter(key, r.Body),
			}

			if err := manager.Update(&blob); err != nil {
				log.Printf("server: update failed %v", err)
				RenderError(w, GetError("internal_server_error"))
			} else {
				resp := PutBlobResponse{
					JWT: GenerateJWT(key, &blob),
				}

				w.Write(Dump(resp))
			}
		}
	}
}
