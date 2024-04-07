package tmpl

import (
	"net/http"

	"r.tomng.dev/goserve/internal/files"
	"r.tomng.dev/goserve/internal/tmpl/templates"
)

func RenderDirectoryView(w http.ResponseWriter, r *http.Request, relativePath string, entries []files.DirEntry) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")

	templComp := templates.DirectoryView(relativePath, entries)
	if err := templComp.Render(r.Context(), w); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
