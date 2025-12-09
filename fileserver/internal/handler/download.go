package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"example.com/filecloud/model"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func HandleDownload(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	id := chi.URLParam(r, "id")
	var f model.FileMeta
	err := db.GetContext(r.Context(), &f, "SELECT * FROM files WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fh, err := os.Open(f.Path)
	if err != nil {
		http.Error(w, "open file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer fh.Close()

	w.Header().Set("Content-Type", f.Mime)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", f.Size))
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(f.Filename)+"\"")
	http.ServeContent(w, r, f.Filename, f.CreatedAt, fh)
}
