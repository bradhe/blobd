package server

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bradhe/blobd/blobs"
	"github.com/bradhe/blobd/crypt"
	"github.com/bradhe/blobd/storage/managers"
)

type handler struct {
	// The context this request was started in.
	Context context.Context

	// Any URL vars that were passed in to the path.
	Vars map[string]string

	// Managers to use in the current context.
	Managers managers.Managers
}

type PostBlobResponse struct {
	Id          string `json:"blob_id"`
	ExpiresAt   string `json:"expires_at"`
	WritableJWT string `json:"write_jwt"`
	ReadOnlyJWT string `json:"read_jwt"`
}

func (h handler) upload(key *crypt.Key, r *http.Request) (*blobs.Blob, error) {
	reader := newRequestReader(r)

	log.WithFields(map[string]interface{}{
		"mediatype": reader.ContentType(),
	}).Debugf("uploading file")

	manager := h.Managers.Blobs()

	// TODO: Validate this is JSON on the way in.
	blob := blobs.Blob{
		Body:      crypt.NewEncrypter(key, reader),
		ExpiresAt: blobs.DefaultExpirationFromNow(),
		MediaType: reader.ContentType(),
	}

	if err := manager.Create(&blob); err != nil {
		log.WithError(err).Error("failed to create blob")
		return nil, err
	} else {
		return &blob, nil
	}
}

func (h handler) PostBlob(w http.ResponseWriter, r *http.Request) {
	// always close the body when we're done with it otherwise we end up with a
	// bunch of open handles over time.
	defer r.Body.Close()

	key, err := crypt.NewRandomKey()

	if err != nil {
		log.WithError(err).Error("key generation failed")
		RenderError(w, GetError("internal_server_error"))
	} else {
		defer key.Destroy()

		if blob, err := h.upload(key, r); err != nil {
			log.WithError(err).Error("failed to upload request content")
			RenderError(w, GetError("internal_server_error"))
		} else {
			resp := PostBlobResponse{
				Id:          blob.Id.String(),
				ExpiresAt:   blob.ExpiresAt.UTC().Format(time.RFC3339),
				ReadOnlyJWT: GenerateJWT(ReadOnlyToken, key, blob),
				WritableJWT: GenerateJWT(WritableToken, key, blob),
			}

			RenderJSON(w, resp)
		}
	}
}

type blobHandler struct {
	handler

	AuthorizedBlobId blobs.Id
	RequestedBlobId  blobs.Id

	Claims *BlobClaims
}

type authenticatedHandlerFunc = func(claims *BlobClaims, w http.ResponseWriter, r *http.Request)

func (h *blobHandler) withAuthenticatedRequest(token string, w http.ResponseWriter, r *http.Request, fn authenticatedHandlerFunc) {
	_, claims, err := ParseJWT(token)

	if claims != nil {
		defer claims.Destroy()
	}

	if err != nil {
		log.WithError(err).Error("failed to parse token")
		RenderError(w, GetError("unauthorized"))
	} else {
		h.AuthorizedBlobId = claims.BlobId
		h.Claims = claims

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

func (h *blobHandler) withAuthenticatedHeaderRequest(w http.ResponseWriter, r *http.Request, fn authenticatedHandlerFunc) {
	h.withAuthenticatedRequest(r.Header.Get("Authorization"), w, r, fn)
}

func (h *blobHandler) withAuthenticatedDownloadRequest(w http.ResponseWriter, r *http.Request, fn authenticatedHandlerFunc) {
	vals := r.URL.Query()
	h.withAuthenticatedRequest(vals.Get("token"), w, r, fn)
}

func isDownloadRequest(r *http.Request) bool {
	vals := r.URL.Query()
	return strings.EqualFold(r.Method, "get") && vals.Get("dl") == "1"
}

func (h *blobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if isDownloadRequest(r) {
		h.withAuthenticatedDownloadRequest(w, r, func(claims *BlobClaims, w http.ResponseWriter, r *http.Request) {
			if claims.IsReadable() {
				h.DownloadBlob(w, r)
			} else {
				RenderError(w, GetError("unauthorized"))
			}
		})
	} else {
		h.withAuthenticatedHeaderRequest(w, r, func(claims *BlobClaims, w http.ResponseWriter, r *http.Request) {
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
			default:
				RenderError(w, GetError("method_not_allowed"))
			}
		})
	}
}

func (h *blobHandler) GetBlob(w http.ResponseWriter, r *http.Request) {
	manager := h.Managers.Blobs()

	if blob, err := manager.Get(h.RequestedBlobId); err != nil {
		log.WithError(err).Error("failed to get blob")
		RenderError(w, GetError("internal_server_error"))
	} else {
		w.Header().Set("Content-Type", string(h.Claims.MediaType))
		io.Copy(w, crypt.NewDecrypter(h.Claims.Key, blob.Body))
	}
}

func (h *blobHandler) DownloadBlob(w http.ResponseWriter, r *http.Request) {
	manager := h.Managers.Blobs()

	if blob, err := manager.Get(h.RequestedBlobId); err != nil {
		log.WithError(err).Error("failed to get blob")
		RenderError(w, GetError("internal_server_error"))
	} else {
		w.Header().Set("Content-Type", string(h.Claims.MediaType))
		w.Header().Set("Content-Disposition", "attachment")
		io.Copy(w, crypt.NewDecrypter(h.Claims.Key, blob.Body))
	}
}

type PutBlobResponse struct {
	ExpiresAt   string `json:"expires_at"`
	ReadOnlyJWT string `json:"read_jwt"`
	WritableJWT string `json:"write_jwt"`
}

func (h *blobHandler) PutBlob(w http.ResponseWriter, r *http.Request) {
	manager := h.Managers.Blobs()

	blob := blobs.Blob{
		Id:        h.RequestedBlobId,
		ExpiresAt: h.Claims.ExpiresAt(),
		Body:      crypt.NewEncrypter(h.Claims.Key, r.Body),
	}

	if err := manager.Update(&blob); err != nil {
		log.WithError(err).Error("update failed")
		RenderError(w, GetError("internal_server_error"))
	} else {
		resp := PutBlobResponse{
			ExpiresAt:   blob.ExpiresAt.UTC().Format(time.RFC3339),
			ReadOnlyJWT: GenerateJWT(ReadOnlyToken, h.Claims.Key, &blob),
			WritableJWT: GenerateJWT(WritableToken, h.Claims.Key, &blob),
		}

		RenderJSON(w, resp)
	}
}
