//go:generate npm run-script build
//go:generate go-bindata -pkg=ui -o=ui.gen.go -ignore=\.swp ./build/...
package ui

import "path"
import "mime"
import "strings"
import "net/http"

const assetPrefix = "build"

var PathPrefix = "/ui"

func assetWithPrefix(asset string) string {
	asset = strings.TrimPrefix(asset, assetPrefix)

	if len(asset) > 0 && asset[0] != '/' {
		return PathPrefix + "/" + asset
	} else {
		return PathPrefix + asset
	}
}

func assetFromPath(path string) string {
	return assetPrefix + strings.TrimPrefix(path, PathPrefix)
}

func Paths() []string {
	paths := make([]string, 0)

	for key := range _bindata {
		// Strip the prefix for now.
		paths = append(paths, assetWithPrefix(key))
	}

	return paths
}

func ServeAsset(p string) func(w http.ResponseWriter, r *http.Request) {
	val := mime.TypeByExtension(path.Ext(p))

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", val)
		w.WriteHeader(http.StatusOK)
		w.Write(MustAsset(assetFromPath(p)))
	}
}
