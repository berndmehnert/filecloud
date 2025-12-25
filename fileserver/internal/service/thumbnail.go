// internal/service/thumbnail.go
package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/jmoiron/sqlx"
)

type ThumbnailJob struct {
	FileID   int64
	FilePath string
	Filename string
}

type ThumbnailService struct {
	jobs       chan ThumbnailJob
	wg         sync.WaitGroup
	db         *sqlx.DB
	thumbDir   string
	maxWorkers int
}

func NewThumbnailService(db *sqlx.DB, baseDir string, maxWorkers int) *ThumbnailService {
	thumbDir := filepath.Join(baseDir, "thumbnails")
	os.MkdirAll(thumbDir, 0o755)

	ts := &ThumbnailService{
		jobs:       make(chan ThumbnailJob, 100),
		db:         db,
		thumbDir:   thumbDir,
		maxWorkers: maxWorkers,
	}

	ts.startWorkers()
	return ts
}

func (ts *ThumbnailService) startWorkers() {
	for i := 0; i < ts.maxWorkers; i++ {
		ts.wg.Add(1)
		go ts.worker(i)
	}
}

func (ts *ThumbnailService) worker(id int) {
	defer ts.wg.Done()

	for job := range ts.jobs {
		if err := ts.processJob(job); err != nil {
			log.Printf("Worker %d: thumbnail failed for file %d: %v", id, job.FileID, err)
			ts.updateStatus(job.FileID, "failed", nil)
		} else {
			log.Printf("Worker %d: thumbnail created for file %d", id, job.FileID)
		}
	}
}

func (ts *ThumbnailService) processJob(job ThumbnailJob) error {
	ts.updateStatus(job.FileID, "processing", nil)

	src, err := imaging.Open(job.FilePath)
	if err != nil {
		return fmt.Errorf("open image: %w", err)
	}

	// 200px width, maintain aspect ratio
	thumb := imaging.Resize(src, 200, 0, imaging.Lanczos)

	thumbFileName := fmt.Sprintf("thumb_%d.jpg", job.FileID)
	thumbPath := filepath.Join(ts.thumbDir, thumbFileName)

	if err := imaging.Save(thumb, thumbPath, imaging.JPEGQuality(80)); err != nil {
		return fmt.Errorf("save thumbnail: %w", err)
	}

	ts.updateStatus(job.FileID, "ready", &thumbPath)
	return nil
}

func (ts *ThumbnailService) updateStatus(id int64, status string, path *string) {
	_, err := ts.db.ExecContext(context.Background(), `
		UPDATE files 
		SET thumbnail_status = ?, thumbnail_path = ?
		WHERE id = ?
	`, status, path, id)
	if err != nil {
		log.Printf("Failed to update thumbnail status for file %d: %v", id, err)
	}
}

func (ts *ThumbnailService) Queue(fileID int64, filePath, filename string) {
	select {
	case ts.jobs <- ThumbnailJob{FileID: fileID, FilePath: filePath, Filename: filename}:
	default:
		log.Printf("Warning: thumbnail queue full, skipping file %d", fileID)
	}
}

func (ts *ThumbnailService) Shutdown() {
	close(ts.jobs)
	ts.wg.Wait()
}
