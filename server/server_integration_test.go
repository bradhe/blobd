// +build integration

package server_test

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bradhe/blobd/server"
)

func randomBlob(size int) []byte {
	blob := make([]byte, size)

	// TODO: This technically has some failure modes we should account for here
	// in that the blob could be not filled all the way!
	rand.Read(blob)

	return blob
}

func TestSecureBlobCreation(t *testing.T) {
	blob := randomBlob(4096)

	// TODO: Default behavior is to use inmem storage...what if that changes?
	s := server.New(server.ServerOptions{})

	r := httptest.NewRequest("POST", "/", bytes.NewBuffer(blob))
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusFound)
	assert.NotEmpty(t, w.HeaderMap["Location"])

	// The response should contain a valid JWT.
	var res server.PostBlobResponse

	dec := json.NewDecoder(w.Body)
	assert.NoError(t, dec.Decode(&res))
	assert.NotEmpty(t, res.WritableJWT)
	assert.NotEmpty(t, res.ReadOnlyJWT)

	// we shouldn't be able to fetch this blob without authentication of some kind.
	loc := w.HeaderMap.Get("Location")

	r = httptest.NewRequest("GET", loc, nil)
	w = httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusUnauthorized)

	// if we use the JWT that we got back, we should be OK.
	r = httptest.NewRequest("GET", loc, nil)
	r.Header.Set("Authorization", "Bearer "+res.ReadOnlyJWT)

	w = httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Body.Bytes(), blob)
}

func TestBlobUpdating(t *testing.T) {
	blob := randomBlob(4096)
	otherBlob := randomBlob(2048)
	require.NotEqual(t, blob, otherBlob)

	// TODO: Default behavior is to use inmem storage...what if that changes?
	s := server.New(server.ServerOptions{})

	r := httptest.NewRequest("POST", "/", bytes.NewBuffer(blob))
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusFound)

	// The response should contain a valid JWT.
	var res server.PostBlobResponse

	dec := json.NewDecoder(w.Body)
	assert.NoError(t, dec.Decode(&res))

	loc := w.HeaderMap.Get("Location")

	r = httptest.NewRequest("PUT", loc, bytes.NewBuffer(otherBlob))
	r.Header.Set("Authorization", "Bearer "+res.WritableJWT)

	w = httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusOK)

	// TODO: Test that we get a valid JWT back again.

	// we should be able to get the blob again, but it should be the updated
	// version of the blob.
	r = httptest.NewRequest("GET", loc, nil)
	r.Header.Set("Authorization", "Bearer "+res.ReadOnlyJWT)

	w = httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Body.Bytes(), otherBlob)
}

func TestInvalidJWTDuringUpdate(t *testing.T) {
	blob := randomBlob(4096)
	otherBlob := randomBlob(2048)
	require.NotEqual(t, blob, otherBlob)

	// TODO: Default behavior is to use inmem storage...what if that changes?
	s := server.New(server.ServerOptions{})

	r := httptest.NewRequest("POST", "/", bytes.NewBuffer(blob))
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusFound)

	// The response should contain a valid JWT.
	var res server.PostBlobResponse

	dec := json.NewDecoder(w.Body)
	assert.NoError(t, dec.Decode(&res))

	// create a second blob which will be secured with a different token.
	r = httptest.NewRequest("POST", "/", bytes.NewBuffer(otherBlob))
	w = httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusFound)

	loc := w.HeaderMap.Get("Location")

	r = httptest.NewRequest("PUT", loc, bytes.NewBuffer(blob))

	// this is the wrong JWT--it's the JWT from the *first* request.
	r.Header.Set("Authorization", "Bearer "+res.WritableJWT)

	w = httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestReadOnlyJWTDuringUpdate(t *testing.T) {
	blob := randomBlob(4096)
	otherBlob := randomBlob(2048)
	require.NotEqual(t, blob, otherBlob)

	// TODO: Default behavior is to use inmem storage...what if that changes?
	s := server.New(server.ServerOptions{})

	r := httptest.NewRequest("POST", "/", bytes.NewBuffer(blob))
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusFound)

	loc := w.HeaderMap.Get("Location")

	// The response should contain a valid JWT.
	var res server.PostBlobResponse

	dec := json.NewDecoder(w.Body)
	assert.NoError(t, dec.Decode(&res))

	r = httptest.NewRequest("PUT", loc, bytes.NewBuffer(blob))

	// this is the wrong JWT--it's the JWT from the *first* request.
	r.Header.Set("Authorization", "Bearer "+res.ReadOnlyJWT)

	w = httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestInvalidJWTDuringGet(t *testing.T) {
	blob := randomBlob(4096)
	otherBlob := randomBlob(2048)
	require.NotEqual(t, blob, otherBlob)

	// TODO: Default behavior is to use inmem storage...what if that changes?
	s := server.New(server.ServerOptions{})

	r := httptest.NewRequest("POST", "/", bytes.NewBuffer(blob))
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusFound)

	// The response should contain a valid JWT.
	var res server.PostBlobResponse

	dec := json.NewDecoder(w.Body)
	assert.NoError(t, dec.Decode(&res))

	// create a second blob which will be secured with a different token.
	r = httptest.NewRequest("POST", "/", bytes.NewBuffer(otherBlob))
	w = httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusFound)

	loc := w.HeaderMap.Get("Location")

	r = httptest.NewRequest("GET", loc, nil)

	// this is the wrong JWT--it's the JWT from the *first* request.
	r.Header.Set("Authorization", "Bearer "+res.ReadOnlyJWT)

	w = httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}
