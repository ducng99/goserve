package tmpl

import (
	"context"
	"net/http"

	"r.tomng.dev/goserve/internal/files"
	"r.tomng.dev/goserve/internal/tmpl/templates"
)

func RenderDirectoryView(w http.ResponseWriter, rootDir, path string) {
	entries, err := files.GetEntries(rootDir, path)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	templComp := templates.DirectoryView(path, entries)
	if err := templComp.Render(context.TODO(), w); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
