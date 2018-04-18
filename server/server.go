package server

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/bradhe/blobd/storage"
	"github.com/bradhe/blobd/storage/managers"
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
	r.Handle("/{blob_id}", BlobHandler{Handler: h})

	// Custom walk of the routes to extract the variables we defined
	// in the map here. If we can match a route, we'll populate the
	// handler's variables with the found variables.
	//
	// Note that this is basically just hijacked from gorilla/mux's default
	// implementation but allows us to inject this in intermediate handlers.
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
	s.getMux(r.Context(), r).ServeHTTP(w, r)
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
