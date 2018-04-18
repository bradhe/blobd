package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"

	"github.com/bradhe/blobd/blobs"
	"github.com/bradhe/blobd/crypt"
	"github.com/bradhe/blobd/storage"
	"github.com/bradhe/blobd/storage/managers"
)

const (
	DefaultTokenSigningKey = "this is a secret"
)

var (
	TokenSigningKey string = DefaultTokenSigningKey
)

type ServerOptions struct {
	StorageURL *url.URL
}

type Server struct {
	Managers managers.Managers
	Options  ServerOptions
}

type NotFoundHandler struct{}

func (n NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RenderError(w, GetError("not_found", r.Method, r.URL.Path))
}

type MethodNotAllowedHandler struct{}

func (m MethodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RenderError(w, GetError("method_not_allowed", r.URL.Path, r.Method))
}

type BlobClaims struct {
	// The UUID for this blob.
	BlobId uuid.UUID

	// Key used to encrypt the message.
	Key *crypt.Key
}

func (bc *BlobClaims) Valid() error {
	if uuid.Equal(bc.BlobId, uuid.NIL) {
		return ErrInvalidUUID
	}

	if bc.Key == nil {
		return ErrMissingDecryptionKey
	}

	return nil
}

func GetJWT(str string) string {
	return strings.TrimPrefix(str, "Bearer ")
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		str := GetJWT(r.Header.Get("Authorization"))

		token, err := jwt.ParseWithClaims(str, &BlobClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("server: unexpected signing method `%v`", token.Header["alg"])
			}

			return []byte(TokenSigningKey), nil
		})

		if err != nil {
			log.Printf("server: failed to get token. %v", err)
			RenderError(w, GetError("unauthorized"))
		} else {
			// TODO: Do something with these claims. We should put them in to the
			// current context or something.
			if claims, ok := token.Claims.(*BlobClaims); !ok || !token.Valid {
				log.Printf("server: invalid claims. %v", err)
				RenderError(w, GetError("unauthorized"))
			} else {
				ctx := WithDecryptionKey(r.Context(), claims.Key)
				ctx = WithBlobId(ctx, claims.BlobId)

				next(w, r.WithContext(ctx))
			}
		}
	}
}

func GenerateJWT(key *crypt.Key, blob *blobs.Blob) string {
	claims := &BlobClaims{
		BlobId: blob.Id,
		Key:    key,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	str, err := token.SignedString([]byte(TokenSigningKey))

	if err != nil {
		panic(err)
	}

	return str
}

func (s *Server) getMux(ctx context.Context, req *http.Request) http.Handler {
	h := Handler{
		Vars:     make(map[string]string),
		Context:  ctx,
		Managers: s.Managers.WithContext(ctx),
	}

	r := mux.NewRouter()
	r.NotFoundHandler = NotFoundHandler{}
	r.MethodNotAllowedHandler = MethodNotAllowedHandler{}

	// Unauthenticated...
	r.HandleFunc("/", h.PostBlob).Methods("POST")
	r.HandleFunc("/{blob_id}", Authenticate(h.GetBlob)).Methods("GET")
	r.HandleFunc("/{blob_id}", Authenticate(h.PutBlob)).Methods("PUT")

	// Custom walk of the routes to extract the variables we defined
	// in the map here. If we can match a route, we'll populate the
	// handler's variables with the found variables.
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		var match mux.RouteMatch

		if route.Match(req, &match) {
			for k, v := range match.Vars {
				h.Vars[k] = v
			}
		}

		return nil
	})

	// Now we can return this in order to actually service the request.
	return r
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.getMux(context.Background(), r).ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(addr string) {
	http.ListenAndServe(addr, s)
}

func New(opts ServerOptions) *Server {
	return &Server{
		Managers: storage.New(opts.StorageURL),
		Options:  opts,
	}
}
