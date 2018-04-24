package server

import (
	"encoding/json"
	"net/http"
)

func Dump(v interface{}) []byte {
	if b, err := json.MarshalIndent(v, "", "  "); err != nil {
		log.Printf("server: dumping object failed: %v", err)
		return []byte("")
	} else {
		return b
	}
}

func RenderJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(Dump(v))
}

func DumpString(v interface{}) string {
	return string(Dump(v))
}

func RenderError(w http.ResponseWriter, err Error) {
	w.WriteHeader(err.Status)
	RenderJSON(w, err)
}

func requestContentType(req *http.Request) string {
	// TODO: How to securely destory this content?
	if t := req.Header.Get("Content-Type"); t != "" {
		return t
	}

	return "application/octet-stream"
}

func RedirectHandlerFunc(to string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, http.StatusFound)
	}
}

func CORSHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		// Best we can do for allowing this thing.
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
}
