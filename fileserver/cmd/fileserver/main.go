package main

import (
	"log"
	"net/http"
	"os"

	"example.com/filecloud/internal/database"
	"example.com/filecloud/internal/handler"
	"example.com/filecloud/internal/middleware"
	"example.com/filecloud/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const (
	databaseFile = "./files.db"
	storageDir   = "./storage"
)

func main() {
	if err := os.MkdirAll(storageDir, 0o755); err != nil {
		log.Fatalf("creating storage dir: %v", err)
	}

	db, err := database.New(databaseFile)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// Migration:
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// Routing:
	r := chi.NewRouter()
	r.Use(middleware.CorsMiddleware)

	// Services:
	thumbService := service.NewThumbnailService(db, "./storage", 3)
	go processPendingThumbnails(db, thumbService)

	// API version 1 group
	r.Route("/v1", func(r chi.Router) {

		uploadHandler := &handler.UploadHandler{
			DB:               db,
			ThumbnailService: thumbService,
		}
		thumbnailHandler := &handler.ThumbnailHandler{DB: db}

		r.Post("/files", uploadHandler.ServeHTTP)

		r.Get("/files/{id}/content", func(w http.ResponseWriter, r *http.Request) {
			handler.HandleDownload(w, r, db)
		})

		r.Get("/files/{id}/thumbnail", thumbnailHandler.ServeHTTP)

		secret := []byte("replace-with-secure-random-secret")
		r.Get("/files", func(w http.ResponseWriter, r *http.Request) {
			handler.HandleListFiles(w, r, db, secret)
		})
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func processPendingThumbnails(db *sqlx.DB, ts *service.ThumbnailService) {
	type pending struct {
		ID       int64  `db:"id"`
		Path     string `db:"path"`
		Filename string `db:"filename"`
	}

	var files []pending
	err := db.Select(&files, `
		SELECT id, path, filename 
		FROM files 
		WHERE thumbnail_status = 'pending'
	`)
	if err != nil {
		log.Printf("Failed to get pending thumbnails: %v", err)
		return
	}

	if len(files) > 0 {
		log.Printf("Queueing %d pending thumbnails", len(files))
		for _, f := range files {
			ts.Queue(f.ID, f.Path, f.Filename)
		}
	}
}
