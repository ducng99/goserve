package tmpl

import (
	"html/template"
	"net/http"

	"r.tomng.dev/goserve/internal/files"
)

type DirectoryViewData struct {
	Path  string
	Entries []files.Entry
}

var DirectoryView = template.Must(template.New("directory_view").Parse(`
<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Goserve Indexing</title>
</head>

<body>
    <h1>Indexing - {{ .Path }}</h1>
	<main>
		<ul>
			{{ range .Entries }}
				<li><a href="/{{ .Path }}">{{ .Name }}</a></li>
			{{ end }}
		</ul>
	</main>
	<footer>
		<i>Powered by <a href="https://r.tomng.dev/goserve">goserve</a></i>
	</footer>
</body>
</html>
`))

func RenderDirectoryView(w http.ResponseWriter, rootDir, path string) {
	entries, err := files.GetEntries(rootDir, path)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := DirectoryViewData{
		Path:  path,
		Entries: entries,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := DirectoryView.Execute(w, data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
