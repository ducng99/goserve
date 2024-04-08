package assets

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type Asset struct {
	Name    string
	Type    string
	Content []byte
}

const PrefixPath = "/_goserveass/"

// Key is the file name to be served
// Value is the file content in bytes
var assets = make(map[string]Asset, 5)

func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	assetName := r.PathValue("asset")

	if asset, ok := assets[assetName]; ok {
		asset.Serve(w)
		return
	}

	http.NotFound(w, r)
}

// Add asset to the pool and return it's URI path
func (a Asset) AddAsset() (string, error) {
	prefix := make([]byte, 4)
	_, err := rand.Read(prefix)
	if err != nil {
		return "", err
	}

	newName := hex.EncodeToString(prefix) + a.Name
	assets[newName] = a

	return PrefixPath + newName, nil
}

func (a Asset) Serve(w http.ResponseWriter) {
	w.Header().Set("Content-Type", a.Type)
	// Cache for 1 hour
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(a.Content)
}
