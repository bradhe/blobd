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

type BlobHandler struct {
	Handler

	AuthorizedBlobId uuid.UUID
	RequestedBlobId  uuid.UUID
	Key              *crypt.Key
}

func (h BlobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, claims, err := ParseJWT(r.Header.Get("Authorization"))

	if err != nil {
		log.Printf("server: failed to get token. %v", err)
		RenderError(w, GetError("unauthorized"))
	} else {
		h.Key = claims.Key
		h.AuthorizedBlobId = claims.BlobId

		id := uuid.Parse(h.Vars["blob_id"])

		if !uuid.Equal(id, h.AuthorizedBlobId) {
			// Invalid blob, so we can't do anything here.
			RenderError(w, GetError("unauthorized"))
		} else {
			h.RequestedBlobId = id

			switch r.Method {
			case "GET":
				h.GetBlob(w, r)
			case "PUT":
				h.PutBlob(w, r)
			}
		}
	}
}

func (h BlobHandler) GetBlob(w http.ResponseWriter, r *http.Request) {
	manager := h.Managers.Blobs()

	if blob, err := manager.Get(h.RequestedBlobId); err != nil {
		log.Printf("server: error happened %v", err)
		RenderError(w, GetError("internal_server_error"))
	} else {
		io.Copy(w, crypt.NewDecrypter(h.Key, blob.Body))
	}
}

type PutBlobResponse struct {
	JWT string `json:"jwt"`
}

func (h BlobHandler) PutBlob(w http.ResponseWriter, r *http.Request) {
	manager := h.Managers.Blobs()

	blob := blobs.Blob{
		Id:   h.RequestedBlobId,
		Body: crypt.NewEncrypter(h.Key, r.Body),
	}

	if err := manager.Update(&blob); err != nil {
		log.Printf("server: update failed %v", err)
		RenderError(w, GetError("internal_server_error"))
	} else {
		resp := PutBlobResponse{
			JWT: GenerateJWT(h.Key, &blob),
		}

		w.Write(Dump(resp))
	}
}
