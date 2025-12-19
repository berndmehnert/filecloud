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

	meta, err := saveFileWithPipe(ctx, file, fh)
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

// helper: ensure storage dir exists (call from init or main)
func ensureStorageDir() error {
	return os.MkdirAll(storageDir, 0o755)
}

// cleanup helper
func removeIfExists(path string) {
	_ = os.Remove(path)
}

// Pipe-based, cancel-aware save
func saveFileWithPipe(ctx context.Context, src multipart.File, fh *multipart.FileHeader) (*model.FileMeta, error) {
	if err := ensureStorageDir(); err != nil {
		return nil, err
	}

	tmpName := fmt.Sprintf("%d_%s.tmp", time.Now().UnixNano(), filepath.Base(fh.Filename))
	tmpPath := filepath.Join(storageDir, tmpName)
	out, err := os.Create(tmpPath)
	if err != nil {
		return nil, err
	}
	// ensure file closed and removed on any early return
	cleanup := func() {
		out.Close()
		removeIfExists(tmpPath)
	}
	defer func() {
		// if tmpPath moved successfully, os.Stat will fail and we just close
		if _, statErr := os.Stat(tmpPath); statErr == nil {
			cleanup()
		} else {
			out.Close()
		}
	}()

	pr, pw := io.Pipe()
	hasher := sha256.New()

	// copy from pipe reader into file+hasher
	copyErrCh := make(chan error, 1)
	go func() {
		_, err := io.Copy(io.MultiWriter(out, hasher), pr)
		copyErrCh <- err
	}()

	// pump src -> pipe writer with ctx cancellation support
	pumpDone := make(chan struct{})
	go func() {
		defer close(pumpDone)
		buf := make([]byte, 32*1024)
		for {
			// check cancellation before read
			select {
			case <-ctx.Done():
				_ = pw.CloseWithError(ctx.Err())
				return
			default:
			}

			n, rerr := src.Read(buf)
			if n > 0 {
				if _, werr := pw.Write(buf[:n]); werr != nil {
					_ = pw.CloseWithError(werr)
					return
				}
			}
			if rerr != nil {
				if rerr == io.EOF {
					_ = pw.Close()
				} else {
					_ = pw.CloseWithError(rerr)
				}
				return
			}
		}
	}()

	// wait for copy to finish or context cancel
	select {
	case err := <-copyErrCh:
		if err != nil {
			return nil, err
		}
	case <-ctx.Done():
		_ = pw.CloseWithError(ctx.Err())
		<-copyErrCh
		return nil, ctx.Err()
	}

	// finalize
	if err := out.Close(); err != nil {
		removeIfExists(tmpPath)
		return nil, err
	}

	finalName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fh.Filename))
	finalPath := filepath.Join(storageDir, finalName)
	if err := os.Rename(tmpPath, finalPath); err != nil {
		removeIfExists(tmpPath)
		return nil, err
	}

	stat, statErr := os.Stat(finalPath)
	if statErr != nil {
		return nil, statErr
	}

	return &model.FileMeta{
		Filename:  fh.Filename,
		Size:      stat.Size(),
		Mime:      fh.Header.Get("Content-Type"),
		Checksum:  hex.EncodeToString(hasher.Sum(nil)),
		Path:      finalPath,
		CreatedAt: time.Now().UTC(),
	}, nil
}
