package server

import (
	"context"
	"net/http"

	"github.com/bradhe/stopwatch"
	"github.com/gorilla/mux"

	"github.com/bradhe/blobd/server/ui"
	"github.com/bradhe/blobd/storage"
	"github.com/bradhe/blobd/storage/managers"
)

type Server struct {
	Managers managers.Managers
	Options  ServerOptions

	ui http.Handler
}

type notFoundHandler struct{}

func (n notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RenderError(w, GetError("not_found", r.Method, r.URL.Path))
}

type methodNotAllowedHandler struct{}

func (m methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RenderError(w, GetError("method_not_allowed", r.URL.Path, r.Method))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CORSHandlerFunc(w, r)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) getMux(ctx context.Context, req *http.Request) http.Handler {
	h := handler{
		Vars:     make(map[string]string),
		Context:  ctx,
		Managers: s.Managers.WithContext(ctx),
	}

	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.NotFoundHandler = notFoundHandler{}
	r.MethodNotAllowedHandler = methodNotAllowedHandler{}

	// Unauthenticated...
	r.HandleFunc("/", h.PostBlob).Methods("POST")
	r.HandleFunc("/", CORSHandlerFunc).Methods("OPTIONS")
	r.Handle("/{blob_id}", &blobHandler{handler: h})
	r.PathPrefix("/ui/").Handler(s.ui)

	// For convenience do a redirect to /ui if there is a GET /
	r.HandleFunc("/", RedirectHandlerFunc("/ui/")).Methods("GET")

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

type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Bytes      int64
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(buf []byte) (int, error) {
	// this is the implicit status code unless one has been explicitly written.
	if w.StatusCode == 0 {
		w.StatusCode = http.StatusOK
	}

	n, err := w.ResponseWriter.Write(buf)
	w.Bytes += int64(n)
	return n, err
}

func newLoggingResponseWriter(base http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{base, 0, 0}
}

func bytes(arr ...int64) uint64 {
	var acc uint64

	for _, b := range arr {
		if b > 0 {
			acc += uint64(b)
		}
	}

	return acc
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrapper := newLoggingResponseWriter(w)

	defer stopwatch.Start().Timer(func(watch stopwatch.Watch) {
		log.WithFields(map[string]interface{}{
			"status": wrapper.StatusCode,
			"bytes":  bytes(wrapper.Bytes, r.ContentLength),
			"time":   watch,
		}).Infof("served %s to %s", r.Method, r.RemoteAddr)
	})

	s.getMux(r.Context(), r).ServeHTTP(wrapper, r)
}

func (s *Server) ListenAndServe(addr string) {
	http.ListenAndServe(addr, s)
}

func New(opts ServerOptions) *Server {
	return &Server{
		Managers: storage.New(opts.StorageURL),
		Options:  opts,
		ui:       ui.Handler(ui.HandlerOptions{Prefix: "/ui/", BlobdHost: "http://localhost:5001"}),
	}
}
