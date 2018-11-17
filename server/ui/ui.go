//go:generate bash ./build.sh
//go:generate go-bindata -pkg=ui -prefix=build -o=ui.gen.go -ignore=\.swp ./build/...
package ui

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"
	"sync"
	"text/template"
)

const (
	defaultAssetName             = "index.html"
	notFoundAssetName            = "404.html"
	notAllowedAssetName          = "405.html"
	internalServerErrorAssetName = "500.html"
)

var (
	defaultAsset             = MustAsset(defaultAssetName)
	notFoundAsset            = MustAsset(notFoundAssetName)
	notAllowedAsset          = MustAsset(notAllowedAssetName)
	internalServerErrorAsset = MustAsset(internalServerErrorAssetName)
)

type HandlerOptions struct {
	Prefix    string
	BlobdHost string
}

type assetHandler struct {
	Options HandlerOptions

	// Serves as a cache for assets that have been previously requested.
	cache map[string][]byte

	mut sync.RWMutex
}

func getBlobdHostTagTemplate(blobdHost string) string {
	return fmt.Sprintf(`<script>window.BLOBD_HOST = "%s";</script>`, blobdHost)
}

func (a *assetHandler) getAsset(path string) ([]byte, error) {
	a.mut.RLock()

	if asset, ok := a.cache[path]; ok {
		log.WithField("path", path).Debug("asset cache hit")
		a.mut.RUnlock()
		return asset, nil
	}

	a.mut.RUnlock()
	a.mut.Lock()

	// Check again just in case the asset was loaded while we were waiting for
	// the lock.
	if asset, ok := a.cache[path]; ok {
		log.WithField("path", path).Debug("latent asset cache hit")
		a.mut.Unlock()
		return asset, nil
	}

	defer a.mut.Unlock()

	log.WithField("path", path).Debug("compiling asset")

	// Now we can actually load the asset.
	var templateParams struct {
		BlobdHostTag string
	}

	templateParams.BlobdHostTag = getBlobdHostTagTemplate(a.Options.BlobdHost)

	if buf, err := Asset(path); err == nil {
		if tmpl, err := template.New(path).Parse(string(buf)); err != nil {
			log.WithError(err).WithField("path", path).Error("failed to instantiate template")
			return nil, ErrTemplateParsingFailed
		} else {
			var out bytes.Buffer

			if err := tmpl.Execute(&out, templateParams); err != nil {
				log.WithError(err).WithField("path", path).Error("failed to process template")
				return nil, ErrTemplateProcessingFailed
			}

			a.cache[path] = out.Bytes()
			return a.cache[path], nil
		}
	} else {
		log.WithField("path", path).Error("asset not found")
		return nil, ErrAssetNotFound
	}
}

func processAssetError(err error) (string, []byte, error) {
	switch err {
	case ErrTemplateProcessingFailed, ErrTemplateProcessingFailed:
		return internalServerErrorAssetName, internalServerErrorAsset, err
	case ErrAssetNotFound:
		return notFoundAssetName, notFoundAsset, err
	default:
		return internalServerErrorAssetName, internalServerErrorAsset, err
	}
}

func (a *assetHandler) AssetFromPath(path string) (string, []byte, error) {
	log.WithField("path", path).Debug("looking up asset")

	path = strings.TrimPrefix(path, a.Options.Prefix)

	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}

	// special-case for serving our default asset
	if path == "" {
		return defaultAssetName, defaultAsset, nil
	}

	if asset, err := a.getAsset(path); err != nil {
		return processAssetError(err)
	} else {
		return path, asset, nil
	}
}

func (a *assetHandler) ServeAsset(w http.ResponseWriter, r *http.Request, status int, name string, asset []byte) {
	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(name)))
	w.WriteHeader(status)
	w.Write(asset)
}

func (a *assetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func Handler(opts HandlerOptions) http.Handler {
	return &assetHandler{
		Options: opts,
		cache:   map[string][]byte{},
	}
}
