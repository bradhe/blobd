//go:generate npm run-script build
//go:generate go-bindata -pkg=ui -prefix=build -o=ui.gen.go -ignore=\.swp ./build/...
package ui

import "path"
import "mime"
import "strings"
import "net/http"

const (
	defaultAssetName    = "index.html"
	notFoundAssetName   = "404.html"
	notAllowedAssetName = "405.html"
)

var (
	defaultAsset    = MustAsset(defaultAssetName)
	notFoundAsset   = MustAsset(notFoundAssetName)
	notAllowedAsset = MustAsset(notAllowedAssetName)
)

type assetHandler struct {
	Prefix string
}

func (a assetHandler) AssetFromPath(path string) (string, []byte, error) {
	log.WithFields(map[string]interface{}{
		"path": path,
	}).Debugf("looking up asset")

	path = strings.TrimPrefix(path, a.Prefix)

	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}

	// special-case for serving our default asset
	if path == "" {
		return defaultAssetName, defaultAsset, nil
	}

	if buf, err := Asset(path); err == nil {
		return path, buf, err
	} else {
		return notFoundAssetName, notFoundAsset, ErrAssetNotFound
	}
}

func (a assetHandler) ServeAsset(w http.ResponseWriter, r *http.Request, status int, name string, asset []byte) {
	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(name)))
	w.WriteHeader(status)
	w.Write(asset)
}

func (a assetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path, asset, err := a.AssetFromPath(r.URL.Path)

		if err == ErrAssetNotFound {
			a.ServeAsset(w, r, http.StatusNotFound, path, asset)
		} else {
			a.ServeAsset(w, r, http.StatusOK, path, asset)
		}
	default:
		a.ServeAsset(w, r, http.StatusMethodNotAllowed, notAllowedAssetName, notAllowedAsset)
	}
}

func Handler(prefix string) http.Handler {
	return &assetHandler{prefix}
}
