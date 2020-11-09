package fs

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
)

var (
	wd   string
	root http.Dir
	fs   http.Handler

	files = []string{
		"/css/*",
		"/js/*",
		"/img/*",
		"/favicon.ico",
		"/index.html",
	}
)

func init() {
	wd, _ = os.Getwd()
	root = http.Dir(filepath.Join(wd, "../pictorio-vue3/dist"))
	fs = http.FileServer(root)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, string(root))
}

func handleFolder(r chi.Router, path string) {
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}

func FileServer(r chi.Router) {
	for _, file := range files {
		handleFolder(r, file)
	}
}
