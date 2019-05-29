// +build integration

package server_test

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func AssertGet(t *testing.T, s *server.Server, path string) io.ReadCloser {
	r := httptest.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	return ioutil.NopCloser(w.Body)
}

func AssertOptions(t *testing.T, s *server.Server, path string) {
	r := httptest.NewRequest("OPTIONS", path, nil)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func AssertGetNotFound(t *testing.T, s *server.Server, path string) {
	r := httptest.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func AssertPutNotAllowed(t *testing.T, s *server.Server, path string) {
	r := httptest.NewRequest("PUT", path, nil)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
}

func AssertCreateBlob(t *testing.T, s *server.Server, blob []byte) server.PostBlobResponse {
	t.Helper()

	r := httptest.NewRequest("POST", "/", bytes.NewBuffer(blob))
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusOK)

	// The response should contain a valid JWT.
	var res server.PostBlobResponse

	dec := json.NewDecoder(w.Body)
	assert.NoError(t, dec.Decode(&res))
	assert.NotEmpty(t, res.Id)
	assert.NotEmpty(t, res.WritableJWT)
	assert.NotEmpty(t, res.ReadOnlyJWT)

	return res
}

func AssertUnauthenticatedReadFails(t *testing.T, s *server.Server, id string) {
	t.Helper()

	r := httptest.NewRequest("GET", "/"+id, nil)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func AssertAuthenticatedReadFails(t *testing.T, s *server.Server, id, jwt string) {
	t.Helper()

	r := httptest.NewRequest("GET", "/"+id, nil)
	r.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func AssertAuthenticatedReadSucceeds(t *testing.T, s *server.Server, id, jwt string) []byte {
	t.Helper()

	r := httptest.NewRequest("GET", "/"+id, nil)
	r.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusOK)

	return w.Body.Bytes()
}

func AssertAuthenticatedUpdateSucceeds(t *testing.T, s *server.Server, id, jwt string, blob []byte) server.PutBlobResponse {
	t.Helper()

	r := httptest.NewRequest("PUT", "/"+id, bytes.NewBuffer(blob))
	r.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusOK)

	// The response should contain a valid JWT.
	var res server.PutBlobResponse

	dec := json.NewDecoder(w.Body)
	assert.NoError(t, dec.Decode(&res))
	assert.NotEmpty(t, res.WritableJWT)
	assert.NotEmpty(t, res.ReadOnlyJWT)

	return res
}

func AssertAuthenticatedUpdateFails(t *testing.T, s *server.Server, id, jwt string, blob []byte) {
	t.Helper()

	r := httptest.NewRequest("PUT", "/"+id, bytes.NewBuffer(blob))
	r.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestSecureBlobCreation(t *testing.T) {
	blob := randomBlob(4096)

	// TODO: Default behavior is to use inmem storage...what if that changes?
	s := server.New(server.ServerOptions{})

	res := AssertCreateBlob(t, s, blob)

	AssertUnauthenticatedReadFails(t, s, res.Id)
	body := AssertAuthenticatedReadSucceeds(t, s, res.Id, res.ReadOnlyJWT)
	assert.Equal(t, body, blob)
}

func TestBlobUpdating(t *testing.T) {
	blob := randomBlob(4096)
	otherBlob := randomBlob(2048)
	require.NotEqual(t, blob, otherBlob)

	// TODO: Default behavior is to use inmem storage...what if that changes?
	s := server.New(server.ServerOptions{})

	res := AssertCreateBlob(t, s, blob)
	AssertAuthenticatedUpdateSucceeds(t, s, res.Id, res.WritableJWT, otherBlob)
	body := AssertAuthenticatedReadSucceeds(t, s, res.Id, res.ReadOnlyJWT)
	assert.Equal(t, body, otherBlob)
}

func TestInvalidJWTDuringUpdate(t *testing.T) {
	blob := randomBlob(4096)
	otherBlob := randomBlob(2048)
	require.NotEqual(t, blob, otherBlob)

	// TODO: Default behavior is to use inmem storage...what if that changes?
	s := server.New(server.ServerOptions{})

	// first blob is what we'll use on second update.
	res := AssertCreateBlob(t, s, blob)

	// create a second blob which will be secured with a different token.
	res2 := AssertCreateBlob(t, s, otherBlob)

	AssertAuthenticatedUpdateFails(t, s, res2.Id, res.WritableJWT, blob)
}

func TestReadOnlyJWTDuringUpdate(t *testing.T) {
	blob := randomBlob(4096)
	otherBlob := randomBlob(2048)
	require.NotEqual(t, blob, otherBlob)

	// TODO: Default behavior is to use inmem storage...what if that changes?
	s := server.New(server.ServerOptions{})

	res := AssertCreateBlob(t, s, blob)
	AssertAuthenticatedUpdateFails(t, s, res.Id, res.ReadOnlyJWT, blob)
}

func TestInvalidJWTDuringGet(t *testing.T) {
	blob := randomBlob(4096)
	otherBlob := randomBlob(2048)
	require.NotEqual(t, blob, otherBlob)

	// TODO: Default behavior is to use inmem storage...what if that changes?
	s := server.New(server.ServerOptions{})

	firstBlob := AssertCreateBlob(t, s, blob)
	secondBlob := AssertCreateBlob(t, s, otherBlob)

	AssertAuthenticatedReadFails(t, s, secondBlob.Id, firstBlob.ReadOnlyJWT)
}

func TestServingUI(t *testing.T) {
	s := server.New(server.ServerOptions{})
	AssertGet(t, s, "/ui/")
	AssertGet(t, s, "/ui/index.html")
	AssertGetNotFound(t, s, "/ui/something-is-not-right.html")
	AssertPutNotAllowed(t, s, "/ui/")
	AssertOptions(t, s, "/")
}

func TestBrowserDownloading(t *testing.T) {
	s := server.New(server.ServerOptions{})
	content := randomBlob(2048)

	blob := AssertCreateBlob(t, s, content)
	resp := AssertGet(t, s, fmt.Sprintf("/%s?token=%s&dl=1", blob.Id, blob.ReadOnlyJWT))

	// makes sure we actually get the content of our original blob back.
	body, _ := ioutil.ReadAll(resp)
	assert.Equal(t, body, content)
}
