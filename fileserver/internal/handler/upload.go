package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"example.com/filecloud/model"

	"github.com/jmoiron/sqlx"
)

const (
	maxUploadSize = 1024 * 1024 * 1024 // 1 GiB
	storageDir    = "./storage"
)

func HandleUpload(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32 MiB memory
		http.Error(w, "invalid multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}
	file, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file field: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	meta, err := saveFile(ctx, file, fh)
	if err != nil {
		http.Error(w, "save error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := db.NamedExecContext(ctx, `
		INSERT INTO files (filename, size, mime, checksum, path, created_at)
		VALUES (:filename, :size, :mime, :checksum, :path, :created_at)
	`, meta)
	if err != nil {
		http.Error(w, "db insert: "+err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id":%d,"filename":"%s","size":%d,"checksum":"%s"}`, id, meta.Filename, meta.Size, meta.Checksum)
}

func saveFile(ctx context.Context, src multipart.File, fh *multipart.FileHeader) (*model.FileMeta, error) {
	tmpName := fmt.Sprintf("%d_%s.tmp", time.Now().UnixNano(), filepath.Base(fh.Filename))
	tmpPath := filepath.Join(storageDir, tmpName)
	out, err := os.Create(tmpPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		out.Close()
		// on error, caller could remove file; keep for simplicity
	}()

	hasher := sha256.New()
	written, err := io.Copy(io.MultiWriter(out, hasher), src)
	if err != nil {
		return nil, err
	}

	finalName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fh.Filename))
	finalPath := filepath.Join(storageDir, finalName)
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return nil, err
	}

	return &model.FileMeta{
		Filename:  fh.Filename,
		Size:      written,
		Mime:      fh.Header.Get("Content-Type"),
		Checksum:  hex.EncodeToString(hasher.Sum(nil)),
		Path:      finalPath,
		CreatedAt: time.Now().UTC(),
	}, nil
}
