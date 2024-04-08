package dirview

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"r.tomng.dev/goserve/internal/files"
	"r.tomng.dev/goserve/internal/logger"
	"r.tomng.dev/goserve/internal/tmpl/dirview/themes"
	"r.tomng.dev/goserve/internal/tmpl/dirview/themes/basic"
	"r.tomng.dev/goserve/internal/tmpl/dirview/themes/pretty"
)

func Render(w http.ResponseWriter, r *http.Request, relativePath string, entries []files.DirEntry, nonce string, theme string) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")

	var templComp templ.Component
	ctx := context.WithValue(r.Context(), "nonce", nonce)

	switch theme {
	case themes.ThemePretty:
		files.Sort(entries)
		templComp = pretty.View(relativePath, entries)
	default:
		templComp = basic.View(relativePath, entries)
	}

	if err := templComp.Render(ctx, w); err != nil {
		http.Error(w, "Cannot render indexing page", http.StatusInternalServerError)
		logger.Printf(logger.LogError, "%v\n", err)
	}
}
