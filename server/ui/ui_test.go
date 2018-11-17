// +build unit

package ui

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServingAssets(t *testing.T) {
	t.Run("with user defined prefix", func(t *testing.T) {
		h := Handler(HandlerOptions{Prefix: "/my-prefix/"})

		r := httptest.NewRequest("GET", "/my-prefix/index.html", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusOK)
	})

	t.Run("with alternate prefix", func(t *testing.T) {
		h := Handler(HandlerOptions{Prefix: "/my-other-prefix"})

		r := httptest.NewRequest("GET", "/my-other-prefix/index.html", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusOK)
	})

	t.Run("with default path", func(t *testing.T) {
		h := Handler(HandlerOptions{Prefix: "/my-other-prefix"})

		r := httptest.NewRequest("GET", "/my-other-prefix/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusOK)
	})

	t.Run("with missing asset", func(t *testing.T) {
		h := Handler(HandlerOptions{Prefix: "/my-other-prefix"})

		r := httptest.NewRequest("GET", "/my-other-prefix/some-dumb-asset.html", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusNotFound)
	})
}

func TestMethodsNotAllowed(t *testing.T) {
	t.Run("with POST", func(t *testing.T) {
		h := Handler(HandlerOptions{Prefix: "/my-prefix/"})

		r := httptest.NewRequest("POST", "/my-prefix/index.html", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
	})

	t.Run("with PUT", func(t *testing.T) {
		h := Handler(HandlerOptions{Prefix: "/my-prefix/"})

		r := httptest.NewRequest("PUT", "/my-prefix/index.html", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
	})

	t.Run("with DELETE", func(t *testing.T) {
		h := Handler(HandlerOptions{Prefix: "/my-prefix/"})

		r := httptest.NewRequest("DELETE", "/my-prefix/index.html", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
	})
}
