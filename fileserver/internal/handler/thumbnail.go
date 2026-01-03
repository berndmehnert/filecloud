// internal/handler/thumbnail.go
package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type ThumbnailHandler struct {
	DB *sqlx.DB
}

func (h *ThumbnailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var thumbPath sql.NullString
	var thumbStatus string

	err = h.DB.QueryRowContext(r.Context(), `
		SELECT thumbnail_path, thumbnail_status FROM files WHERE id = ?
	`, id).Scan(&thumbPath, &thumbStatus)

	if err == sql.ErrNoRows {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	switch thumbStatus {
	case "ready":
		if thumbPath.Valid {
			http.ServeFile(w, r, thumbPath.String)
		} else {
			http.Error(w, "thumbnail path missing", http.StatusInternalServerError)
		}
	case "pending", "processing":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status":"` + thumbStatus + `"}`))
	case "failed":
		http.Error(w, "thumbnail generation failed", http.StatusInternalServerError)
	default:
		http.Error(w, "no thumbnail available", http.StatusNotFound)
	}
}
