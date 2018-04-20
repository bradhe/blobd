package server

import (
	"context"
	"io"
	"net/http"

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
	WritableJWT string `json:"write_jwt"`
	ReadOnlyJWT string `json:"read_jwt"`
}

func (h Handler) PostBlob(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	key, err := crypt.NewRandomKey()

	if err != nil {
		log.WithError(err).Error("failed to generate key")
		RenderError(w, GetError("internal_server_error"))
	} else {
		defer key.Destroy()

		manager := h.Managers.Blobs()

		// TODO: Validate this is JSON on the way in.
		blob := blobs.Blob{
			Body:      crypt.NewEncrypter(key, r.Body),
			ExpiresAt: blobs.DefaultExpirationFromNow(),
		}

		if err := manager.Create(&blob); err != nil {
			log.WithError(err).Error("failed to create blob")
			RenderError(w, GetError("internal_server_error"))
		} else {
			resp := PostBlobResponse{
				ReadOnlyJWT: GenerateJWT(ReadOnlyToken, key, &blob),
				WritableJWT: GenerateJWT(WritableToken, key, &blob),
			}

			http.Redirect(w, r, BlobPath(&blob), http.StatusFound)
			w.Write(Dump(resp))
		}
	}
}

type BlobHandler struct {
	Handler

	AuthorizedBlobId blobs.Id
	RequestedBlobId  blobs.Id
	Key              *crypt.Key
}

type authenticatedHandlerFunc = func(claims *BlobClaims, w http.ResponseWriter, r *http.Request)

func (h *BlobHandler) withAuthenticatedRequest(w http.ResponseWriter, r *http.Request, fn authenticatedHandlerFunc) {
	_, claims, err := ParseJWT(r.Header.Get("Authorization"))

	if err != nil {
		log.WithError(err).Error("failed to parse token")
		RenderError(w, GetError("unauthorized"))
	} else {
		defer claims.Key.Destroy()

		h.Key = claims.Key
		h.AuthorizedBlobId = claims.BlobId

		if id, err := blobs.ParseId(h.Vars["blob_id"]); err != nil {
			// Invalid blob, so we can't do anything here.
			RenderError(w, GetError("unauthorized"))
		} else if !h.AuthorizedBlobId.Equal(id) {
			// Invalid blob, so we can't do anything here.
			RenderError(w, GetError("unauthorized"))
		} else {
			h.RequestedBlobId = id

			fn(claims, w, r)
		}
	}
}

func (h *BlobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.withAuthenticatedRequest(w, r, func(claims *BlobClaims, w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if claims.IsReadable() {
				h.GetBlob(w, r)
			} else {
				RenderError(w, GetError("unauthorized"))
			}
		case "PUT":
			if claims.IsWritable() {
				h.PutBlob(w, r)
			} else {
				RenderError(w, GetError("unauthorized"))
			}
		}
	})
}

func (h *BlobHandler) GetBlob(w http.ResponseWriter, r *http.Request) {
	manager := h.Managers.Blobs()

	if blob, err := manager.Get(h.RequestedBlobId); err != nil {
		log.WithError(err).Error("failed to get blob")
		RenderError(w, GetError("internal_server_error"))
	} else {
		io.Copy(w, crypt.NewDecrypter(h.Key, blob.Body))
	}
}

type PutBlobResponse struct {
	ReadOnlyJWT string `json:"read_jwt"`
	WritableJWT string `json:"write_jwt"`
}

func (h *BlobHandler) PutBlob(w http.ResponseWriter, r *http.Request) {
	manager := h.Managers.Blobs()

	blob := blobs.Blob{
		Id:   h.RequestedBlobId,
		Body: crypt.NewEncrypter(h.Key, r.Body),
	}

	if err := manager.Update(&blob); err != nil {
		log.WithError(err).Error("update failed")
		RenderError(w, GetError("internal_server_error"))
	} else {
		resp := PutBlobResponse{
			ReadOnlyJWT: GenerateJWT(ReadOnlyToken, h.Key, &blob),
			WritableJWT: GenerateJWT(WritableToken, h.Key, &blob),
		}

		w.Write(Dump(resp))
	}
}
